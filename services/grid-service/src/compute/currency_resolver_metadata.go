package compute

import (
	"context"
	"fmt"
	"sync"
	"time"

	"quantatomai/grid-service/domain"
	"quantatomai/grid-service/planner"
)

//
// 1. TYPE DECLARATIONS
//

// CurrencyMetadataProvider abstracts metadata lookups for currency attributes.
type CurrencyMetadataProvider interface {
	GetCurrencyCode(ctx context.Context, dimensionID, memberID int64) (string, error)
	BulkGetCurrencyCodes(ctx context.Context, dimensionID int64, memberIDs []int64) (map[int64]string, error)
}

// DimensionRole describes which dimension carries which semantic role.
type DimensionRole string

const (
	DimensionRoleEntity    DimensionRole = "entity"
	DimensionRoleGeography DimensionRole = "geography"
	DimensionRoleCustom    DimensionRole = "custom"
)

// CurrencyDimensionBinding binds a semantic role to a dimension ID and
// its index in AtomKey.DimIDs.
type CurrencyDimensionBinding struct {
	Role        DimensionRole
	DimensionID int64
	Index       int
}

// CurrencyResolverConfig defines settings for the metadata-backed resolver.
type CurrencyResolverConfig struct {
	Bindings         []CurrencyDimensionBinding
	DefaultCurrency  string
	FastPathEnabled  bool
	FastPathIndex    int
	FastPathDimID    int64
	FailureThreshold int
	OpenDuration     time.Duration
	LookupTimeout    time.Duration
}

// CurrencyResolutionContext provides detailed audit information about how a currency was resolved.
type CurrencyResolutionContext struct {
	Currency string
	Role     DimensionRole
	DimID    int64
	MemberID int64
}

// CurrencyResolverWithContext is an interface for resolvers that support comprehensive audit trails.
type CurrencyResolverWithContext interface {
	ResolveCurrencyWithContext(ctx context.Context, key domain.AtomKey) (CurrencyResolutionContext, error)
}

//
// 2. STRUCT & CONSTRUCTOR
//

// CurrencyResolverMetadata is a production-grade metadata-backed currency resolver.
type CurrencyResolverMetadata struct {
	provider     CurrencyMetadataProvider
	cfg          CurrencyResolverConfig
	memo         map[int64]string
	memoLock     sync.RWMutex
	breakerMu    sync.Mutex
	breakerState breakerState
	failureCount int
	openUntil    time.Time
}

// NewCurrencyResolverMetadata constructs a production-ready currency resolver.
func NewCurrencyResolverMetadata(
	provider CurrencyMetadataProvider,
	cfg CurrencyResolverConfig,
) *CurrencyResolverMetadata {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.OpenDuration <= 0 {
		cfg.OpenDuration = 5 * time.Second
	}
	if cfg.LookupTimeout <= 0 {
		cfg.LookupTimeout = 200 * time.Millisecond
	}
	if cfg.DefaultCurrency == "" {
		cfg.DefaultCurrency = "USD"
	}

	return &CurrencyResolverMetadata{
		provider:     provider,
		cfg:          cfg,
		memo:         make(map[int64]string),
		breakerState: breakerClosed,
	}
}

//
// 3. PUBLIC METHODS
//

// Prefetch extracts all unique member IDs from the SparsedBlockRef and resolves
// their currencies in bulk, drastically reducing per-cell IO overhead.
func (r *CurrencyResolverMetadata) Prefetch(
	ctx context.Context,
	block planner.SparsedBlockRef,
) error {
	dimMembers := make(map[int64]map[int64]struct{})
	for _, b := range r.cfg.Bindings {
		dimMembers[b.DimensionID] = make(map[int64]struct{})
	}

	for _, row := range block.RowIDs {
		for _, b := range r.cfg.Bindings {
			if b.Index == 0 {
				dimMembers[b.DimensionID][row] = struct{}{}
			}
		}
	}
	for _, col := range block.ColIDs {
		for _, b := range r.cfg.Bindings {
			if b.Index == 1 {
				dimMembers[b.DimensionID][col] = struct{}{}
			}
		}
	}

	for dimID, members := range dimMembers {
		if len(members) == 0 {
			continue
		}

		memberIDs := make([]int64, 0, len(members))
		for m := range members {
			r.memoLock.RLock()
			_, ok := r.memo[m]
			r.memoLock.RUnlock()
			if !ok {
				memberIDs = append(memberIDs, m)
			}
		}

		if len(memberIDs) == 0 {
			continue
		}

		cctx, cancel := context.WithTimeout(ctx, r.cfg.LookupTimeout)
		results, err := r.provider.BulkGetCurrencyCodes(cctx, dimID, memberIDs)
		cancel()

		if err != nil {
			r.recordFailure()
			return fmt.Errorf("currency_resolver_metadata: bulk prefetch failed: %w", err)
		}

		r.memoLock.Lock()
		for memberID, code := range results {
			r.memo[memberID] = code
		}
		r.memoLock.Unlock()
	}

	r.recordSuccess()
	return nil
}

// ResolveCurrency implements CurrencyResolver with multi-layered optimizations.
func (r *CurrencyResolverMetadata) ResolveCurrency(
	ctx context.Context,
	key domain.AtomKey,
) (string, error) {
	if !r.allowRequest() {
		return r.cfg.DefaultCurrency, nil
	}

	if r.cfg.FastPathEnabled {
		if r.cfg.FastPathIndex < len(key.DimIDs) {
			memberID := key.DimIDs[r.cfg.FastPathIndex]
			if memberID != 0 {
				if code := r.lookupMemo(memberID); code != "" {
					return code, nil
				}
				code, err := r.lookupProvider(ctx, r.cfg.FastPathDimID, memberID)
				if err == nil && code != "" {
					return code, nil
				}
			}
		}
	}

	for _, b := range r.cfg.Bindings {
		if b.Index < 0 || b.Index >= len(key.DimIDs) {
			continue
		}
		memberID := key.DimIDs[b.Index]
		if memberID == 0 {
			continue
		}

		if code := r.lookupMemo(memberID); code != "" {
			return code, nil
		}

		code, err := r.lookupProvider(ctx, b.DimensionID, memberID)
		if err != nil {
			r.recordFailure()
			continue
		}
		if code != "" {
			return code, nil
		}
	}

	return r.cfg.DefaultCurrency, nil
}

// ResolveCurrencyWithContext returns both the currency and the binding metadata
// that produced it (role, dimension, member). This is used by FXAudit.
func (r *CurrencyResolverMetadata) ResolveCurrencyWithContext(
	ctx context.Context,
	key domain.AtomKey,
) (CurrencyResolutionContext, error) {
	// First resolve currency normally.
	currency, err := r.ResolveCurrency(ctx, key)
	if err != nil {
		return CurrencyResolutionContext{}, err
	}

	if currency == "" {
		return CurrencyResolutionContext{
			Currency: r.cfg.DefaultCurrency,
		}, nil
	}

	// Re-run binding logic to determine which binding produced the currency.
	for _, b := range r.cfg.Bindings {
		if b.Index < 0 || b.Index >= len(key.DimIDs) {
			continue
		}
		memberID := key.DimIDs[b.Index]
		if memberID == 0 {
			continue
		}

		if code := r.lookupMemo(memberID); code == currency {
			return CurrencyResolutionContext{
				Currency: currency,
				Role:     b.Role,
				DimID:    b.DimensionID,
				MemberID: memberID,
			}, nil
		}
	}

	return CurrencyResolutionContext{
		Currency: currency,
	}, nil
}

//
// 4. INTERNAL HELPERS
//

func (r *CurrencyResolverMetadata) lookupMemo(memberID int64) string {
	r.memoLock.RLock()
	defer r.memoLock.RUnlock()
	return r.memo[memberID]
}

func (r *CurrencyResolverMetadata) lookupProvider(
	ctx context.Context,
	dimensionID, memberID int64,
) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, r.cfg.LookupTimeout)
	defer cancel()

	code, err := r.provider.GetCurrencyCode(cctx, dimensionID, memberID)
	if err != nil {
		return r.cfg.DefaultCurrency, err
	}

	r.memoLock.Lock()
	r.memo[memberID] = code
	r.memoLock.Unlock()

	r.recordSuccess()
	return code, nil
}

func (r *CurrencyResolverMetadata) Name() string {
	return "currency_resolver_v3"
}

// Breaker State

type breakerState int

const (
	breakerClosed breakerState = iota
	breakerOpen
	breakerHalfOpen
)

func (r *CurrencyResolverMetadata) allowRequest() bool {
	r.breakerMu.Lock()
	defer r.breakerMu.Unlock()

	switch r.breakerState {
	case breakerOpen:
		if time.Now().After(r.openUntil) {
			r.breakerState = breakerHalfOpen
			return true
		}
		return false
	case breakerHalfOpen:
		return true
	default:
		return true
	}
}

func (r *CurrencyResolverMetadata) recordSuccess() {
	r.breakerMu.Lock()
	defer r.breakerMu.Unlock()

	r.failureCount = 0
	if r.breakerState == breakerHalfOpen {
		r.breakerState = breakerClosed
	}
}

func (r *CurrencyResolverMetadata) recordFailure() {
	r.breakerMu.Lock()
	defer r.breakerMu.Unlock()

	r.failureCount++
	if r.failureCount >= r.cfg.FailureThreshold {
		r.breakerState = breakerOpen
		r.openUntil = time.Now().Add(r.cfg.OpenDuration)
	}
}
