package compute

import (
	"context"
	"fmt"
	"time"

	"quantatomai/grid-service/domain"
)

//
// ─────────────────────────────────────────────────────────────
//   CORE INTERFACES
// ─────────────────────────────────────────────────────────────
//

// CurrencyResolver resolves the "local" currency for a given atom coordinate.
type CurrencyResolver interface {
	ResolveCurrency(ctx context.Context, key domain.AtomKey) (string, error)
}

// FXRateProvider abstracts how FX rates are obtained and composed.
type FXRateProvider interface {
	GetRate(ctx context.Context, sourceCurrency, targetCurrency string, asOf time.Time) (float64, error)
}

//
// ─────────────────────────────────────────────────────────────
//   CONFIG
// ─────────────────────────────────────────────────────────────
//

// FXConfig defines how FX should be applied at the transformer level.
type FXConfig struct {
	// TargetCurrency is the reporting currency (e.g., "USD").
	TargetCurrency string

	// AsOf is the effective date/time for FX rates.
	AsOf time.Time

	// Optional pivot currency hint (e.g., "USD") for triangulation.
	PivotCurrency string

	// MeasuresToTranslate is an optional whitelist of measures that should be FX-translated.
	MeasuresToTranslate map[int64]bool

	// EnableAudit controls whether FX audit information is recorded per cell.
	EnableAudit bool
}

//
// ─────────────────────────────────────────────────────────────
//   AUDIT STRUCT (Option B additions included)
// ─────────────────────────────────────────────────────────────
//

// FXAudit captures the FX details used for a specific cell.
type FXAudit struct {
	SourceCurrency string
	TargetCurrency string
	RateUsed       float64
	// Path represents the currency path used, e.g. ["BRL", "USD", "JPY"].
	Path []string

	// NEW: Currency resolution metadata
	ResolvedViaRole   DimensionRole
	ResolvedViaDimID  int64
	ResolvedViaMember int64
}

//
// ─────────────────────────────────────────────────────────────
//   TRANSFORMER STRUCT
// ─────────────────────────────────────────────────────────────
//

// FXTransformer v3 applies FX translation using metadata-backed resolution and pivot-aware rates.
type FXTransformer struct {
	cfg              FXConfig
	rateProvider     FXRateProvider
	currencyResolver CurrencyResolver

	// audit is optional; only populated when cfg.EnableAudit is true.
	audit map[domain.AtomKey]FXAudit
}

//
// ─────────────────────────────────────────────────────────────
//   CONSTRUCTOR
// ─────────────────────────────────────────────────────────────
//

// NewFXTransformer constructs a new multi-axial FX transformer with auditing support.
func NewFXTransformer(cfg FXConfig, provider FXRateProvider, resolver CurrencyResolver) *FXTransformer {
	if cfg.AsOf.IsZero() {
		cfg.AsOf = time.Now().UTC()
	}
	return &FXTransformer{
		cfg:              cfg,
		rateProvider:     provider,
		currencyResolver: resolver,
	}
}

//
// ─────────────────────────────────────────────────────────────
//   NAME
// ─────────────────────────────────────────────────────────────
//

func (t *FXTransformer) Name() string {
	return "fx_v3"
}

//
// ─────────────────────────────────────────────────────────────
//   GET AUDIT
// ─────────────────────────────────────────────────────────────
//

// GetAudit returns the FX audit map (if auditing is enabled).
func (t *FXTransformer) GetAudit() map[domain.AtomKey]FXAudit {
	return t.audit
}

//
// ─────────────────────────────────────────────────────────────
//   APPLY (Option B logic inside)
// ─────────────────────────────────────────────────────────────
//

// Apply performs the FX translation with optional audit trail recording.
func (t *FXTransformer) Apply(
	ctx context.Context,
	atoms map[domain.AtomKey]float64,
) (map[domain.AtomKey]float64, error) {

	if t.cfg.TargetCurrency == "" {
		// No target currency → no-op.
		return atoms, nil
	}
	if t.currencyResolver == nil || t.rateProvider == nil {
		return nil, fmt.Errorf("fx_v3: missing currencyResolver or rateProvider")
	}

	// Initialize audit map only if enabled.
	if t.cfg.EnableAudit && t.audit == nil {
		t.audit = make(map[domain.AtomKey]FXAudit, len(atoms))
	}

	// Copy-on-write: only clone when we actually change something.
	out := atoms

	// Cache for FX rates: "SRC->TGT" → rate to avoid redundant provider calls in this pass.
	rateCache := make(map[string]float64)

	for key, value := range atoms {
		// 1. Measure Whitelist Check
		if len(t.cfg.MeasuresToTranslate) > 0 && !t.cfg.MeasuresToTranslate[key.MeasureID] {
			continue
		}

		// 2. Resolve Local Currency via Coordinate
		// Detect if resolver supports audit context
		resolverCtx, hasCtx := t.currencyResolver.(CurrencyResolverWithContext)

		var srcCurrency string
		var resCtx CurrencyResolutionContext
		var err error

		if hasCtx {
			resCtx, err = resolverCtx.ResolveCurrencyWithContext(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("fx_v3: resolve currency with context failed: %w", err)
			}
			srcCurrency = resCtx.Currency
		} else {
			srcCurrency, err = t.currencyResolver.ResolveCurrency(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("fx_v3: resolve currency failed: %w", err)
			}
		}

		if srcCurrency == "" || srcCurrency == t.cfg.TargetCurrency {
			// No currency found or already in target → skip.
			continue
		}

		// 3. Resolve Rate (with pass-through caching)
		cacheKey := srcCurrency + "->" + t.cfg.TargetCurrency
		rate, ok := rateCache[cacheKey]
		if !ok {
			rate, err = t.rateProvider.GetRate(ctx, srcCurrency, t.cfg.TargetCurrency, t.cfg.AsOf)
			if err != nil {
				return nil, fmt.Errorf("fx_v3: get rate %s failed: %w", cacheKey, err)
			}
			rateCache[cacheKey] = rate
		}

		// 4. Transform Value
		converted := value * rate

		// First mutation → clone map to maintain immutability of the source.
		if out == atoms {
			out = cloneAtomMap(atoms)
		}
		out[key] = converted

		// 5. Audit Trail (Optional)
		if t.cfg.EnableAudit {
			audit := FXAudit{
				SourceCurrency: srcCurrency,
				TargetCurrency: t.cfg.TargetCurrency,
				RateUsed:       rate,
				Path:           []string{srcCurrency, t.cfg.TargetCurrency},
			}

			// NEW: attach resolution metadata
			if hasCtx {
				audit.ResolvedViaRole = resCtx.Role
				audit.ResolvedViaDimID = resCtx.DimID
				audit.ResolvedViaMember = resCtx.MemberID
			}

			t.audit[key] = audit
		}
	}

	return out, nil
}

