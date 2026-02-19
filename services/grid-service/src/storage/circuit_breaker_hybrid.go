package storage

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var breakerTracer = otel.Tracer("quantatomai/circuit-breaker")

type BreakerState string

const (
	BreakerClosed   BreakerState = "CLOSED"
	BreakerOpen     BreakerState = "OPEN"
	BreakerHalfOpen BreakerState = "HALF_OPEN"
)

type FallbackMode string

const (
	FallbackNone    FallbackMode = "NONE"
	FallbackCached  FallbackMode = "CACHED_ONLY"
	FallbackStaleOK FallbackMode = "STALE_OK"
)

type HybridBreakerConfig struct {
	FailureThreshold   int
	ResetTimeout       time.Duration
	WindowSize         time.Duration

	DistributedKey     string
	DistributedTTL     time.Duration
	DistributedEnabled bool

	BaseBackoff        time.Duration
	MaxBackoff         time.Duration

	EnableCachedOnly   bool
	EnableStaleOK      bool

	HalfOpenMaxTrials  int

	OnStateChange      func(oldState, newState BreakerState, reason string)
}

type DistributedStateStore interface {
	GetBreakerState(ctx context.Context, key string) (BreakerState, time.Time, string, error)
	SetBreakerState(ctx context.Context, key string, state BreakerState, retryAfter time.Time, reason string, ttl time.Duration) error
}

// snapshot for atomic fast-path
type breakerSnapshot struct {
	state      BreakerState
	retryAfter time.Time
	reason     string
}

// power-of-two bucketed rolling window (Titanium Tier: Atomic & Zero-Allocation)
type failureBuckets struct {
	buckets     []int32
	bucketStart int64 // UnixNano
	bucketSize  int64 // Nano
	mask        int32
	total       atomic.Int32
}

func newFailureBuckets(windowSize time.Duration, bucketCount int) *failureBuckets {
	if bucketCount <= 0 {
		bucketCount = 16
	}
	n := 1
	for n < bucketCount {
		n <<= 1
	}
	size := int64(windowSize) / int64(n)
	if size <= 0 {
		size = int64(time.Second)
	}

	return &failureBuckets{
		buckets:     make([]int32, n),
		bucketStart: time.Now().UnixNano(),
		bucketSize:  size,
		mask:        int32(n - 1),
	}
}

func (fb *failureBuckets) Reset(now time.Time) {
	for i := range fb.buckets {
		atomic.StoreInt32(&fb.buckets[i], 0)
	}
	atomic.StoreInt64(&fb.bucketStart, now.UnixNano())
	fb.total.Store(0)
}

func (fb *failureBuckets) addFailure(now time.Time) {
	t := now.UnixNano()
	fb.rotate(t)
	idx := fb.index(t)
	atomic.AddInt32(&fb.buckets[idx], 1)
	fb.total.Add(1)
}

func (fb *failureBuckets) count(now time.Time) int {
	fb.rotate(now.UnixNano())
	return int(fb.total.Load())
}

func (fb *failureBuckets) index(now int64) int32 {
	start := atomic.LoadInt64(&fb.bucketStart)
	delta := now - start
	if delta < 0 {
		return 0
	}
	steps := int32(delta / fb.bucketSize)
	return steps & fb.mask
}

func (fb *failureBuckets) rotate(now int64) {
	start := atomic.LoadInt64(&fb.bucketStart)
	if now < start+fb.bucketSize {
		return
	}

	steps := int32((now - start) / fb.bucketSize)
	if steps >= int32(len(fb.buckets)) {
		fb.Reset(time.Unix(0, now))
		return
	}

	for i := int32(0); i < steps; i++ {
		curStart := atomic.LoadInt64(&fb.bucketStart)
		idx := fb.index(curStart + int64(i)*fb.bucketSize)
		v := atomic.SwapInt32(&fb.buckets[idx], 0)
		fb.total.Add(-v)
	}
	atomic.AddInt64(&fb.bucketStart, int64(steps)*fb.bucketSize)
}

type HybridCircuitBreaker struct {
	cfg   HybridBreakerConfig
	store DistributedStateStore

	mu           sync.Mutex
	state        BreakerState
	openedAt     time.Time
	retryAfter   time.Time
	reason       string

	halfOpenTrials atomic.Int32
	rng            *rand.Rand
	rngMu          sync.Mutex

	lastDistributedOpen time.Time
	failures            *failureBuckets
	snap                atomic.Value // breakerSnapshot
}

func NewHybridCircuitBreaker(cfg HybridBreakerConfig, store DistributedStateStore) *HybridCircuitBreaker {
	if cfg.HalfOpenMaxTrials <= 0 {
		cfg.HalfOpenMaxTrials = 1
	}
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = 30 * time.Second
	}

	src := rand.NewSource(time.Now().UnixNano())
	fb := newFailureBuckets(cfg.WindowSize, 16)

	b := &HybridCircuitBreaker{
		cfg:      cfg,
		store:    store,
		state:    BreakerClosed,
		rng:      rand.New(src),
		failures: fb,
	}
	b.snap.Store(breakerSnapshot{
		state: BreakerClosed,
	})
	b.halfOpenTrials.Store(0)
	return b
}

func [T any] (b *HybridCircuitBreaker) Execute(
	ctx context.Context,
	fn func(context.Context) (T, error),
	cachedFn func(context.Context) (T, bool, error),
	staleFn func(context.Context) (T, bool, error),
) (T, error) {
	var zero T

	ctx, span := breakerTracer.Start(ctx, "breaker.execute")
	defer span.End()

	now := time.Now()

	// 1. Titanium Fast-Path (Atomic Check)
	snap := b.snap.Load().(breakerSnapshot)
	if snap.state == BreakerClosed {
		res, err := fn(ctx)
		if err == nil {
			// Success path: Lazy status update
			if b.failures.total.Load() > 0 {
				b.recordSuccess(now)
			}
			span.SetAttributes(attribute.String("breaker.outcome", "success"))
			return res, nil
		}
		// Failure path
		span.RecordError(err)
		b.recordFailure(now, err)
		if b.cfg.DistributedEnabled && b.store != nil && b.justOpened(now) {
			go b.propagateOpenDistributed(context.Background(), err)
		}
		return b.handleFailureWithFallback(ctx, span, now, err, cachedFn, staleFn)
	}

	// 2. Recovery Path (Full Status Check)
	state, retryAfter, reason := b.currentState(now)
	span.SetAttributes(
		attribute.String("breaker.state.local", string(state)),
		attribute.String("breaker.reason.local", reason),
	)

	if state == BreakerOpen {
		if b.cfg.DistributedEnabled && b.store != nil {
			dState, dRetry, dReason, err := b.store.GetBreakerState(ctx, b.cfg.DistributedKey)
			if err == nil && dState == BreakerOpen {
				span.SetAttributes(
					attribute.String("breaker.state.distributed", string(dState)),
					attribute.String("breaker.reason.distributed", dReason),
				)
				return b.handleOpenWithFallback(ctx, span, now, dRetry, dReason, cachedFn, staleFn)
			}
		}
		return b.handleOpenWithFallback(ctx, span, now, retryAfter, reason, cachedFn, staleFn)
	}

	if state == BreakerHalfOpen {
		if !b.allowHalfOpenTrial() {
			return b.handleOpenWithFallback(ctx, span, now, retryAfter, reason, cachedFn, staleFn)
		}
	}

	res, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		b.recordFailure(now, err)
		if b.cfg.DistributedEnabled && b.store != nil && b.justOpened(now) {
			go b.propagateOpenDistributed(context.Background(), err)
		}
		return b.handleFailureWithFallback(ctx, span, now, err, cachedFn, staleFn)
	}

	b.recordSuccess(now)
	span.SetAttributes(attribute.String("breaker.outcome", "success"))
	return res, nil
}

func (b *HybridCircuitBreaker) currentState(now time.Time) (BreakerState, time.Time, string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == BreakerOpen && !b.retryAfter.IsZero() && now.After(b.retryAfter) {
		b.transitionLocked(BreakerHalfOpen, b.reason)
		b.halfOpenTrials.Store(0)
	}

	b.updateSnapshotLocked()
	return b.state, b.retryAfter, b.reason
}

func (b *HybridCircuitBreaker) recordFailure(now time.Time, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.reason = err.Error()
	b.failures.addFailure(now)
	count := b.failures.count(now)

	if b.state == BreakerOpen {
		b.updateSnapshotLocked()
		return
	}

	if count >= b.cfg.FailureThreshold {
		b.state = BreakerOpen
		b.openedAt = now
		backoff := b.computeBackoffLocked()
		b.retryAfter = now.Add(backoff)
		b.transitionLocked(BreakerOpen, err.Error())
	}
	b.updateSnapshotLocked()
}

func (b *HybridCircuitBreaker) recordSuccess(now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures.Reset(now)
	b.reason = ""

	if b.state == BreakerHalfOpen || b.state == BreakerOpen {
		b.transitionLocked(BreakerClosed, "recovered")
	}
	b.updateSnapshotLocked()
}

func (b *HybridCircuitBreaker) updateSnapshotLocked() {
	b.snap.Store(breakerSnapshot{
		state:      b.state,
		retryAfter: b.retryAfter,
		reason:     b.reason,
	})
}

func (b *HybridCircuitBreaker) justOpened(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state == BreakerOpen && now.Sub(b.openedAt) < time.Second
}

func (b *HybridCircuitBreaker) allowHalfOpenTrial() bool {
	max := int32(b.cfg.HalfOpenMaxTrials)
	for {
		cur := b.halfOpenTrials.Load()
		if cur >= max {
			return false
		}
		if b.halfOpenTrials.CompareAndSwap(cur, cur+1) {
			return true
		}
	}
}

func (b *HybridCircuitBreaker) computeBackoffLocked() time.Duration {
	base := b.cfg.BaseBackoff
	if base <= 0 {
		base = time.Second
	}
	max := b.cfg.MaxBackoff
	if max <= 0 {
		max = 30 * time.Second
	}

	count := int(b.failures.total.Load())
	over := count - b.cfg.FailureThreshold
	if over < 0 {
		over = 0
	}
	factor := 1 << over
	d := time.Duration(factor) * base
	if d > max {
		d = max
	}

	b.rngMu.Lock()
	jitter := time.Duration(b.rng.Int63n(int64(base)))
	b.rngMu.Unlock()
	return d + jitter
}

func (b *HybridCircuitBreaker) propagateOpenDistributed(ctx context.Context, err error) {
	if b.store == nil {
		return
	}
	now := time.Now()

	b.mu.Lock()
	if !b.lastDistributedOpen.IsZero() && now.Sub(b.lastDistributedOpen) < 2*time.Second {
		b.mu.Unlock()
		return
	}
	b.lastDistributedOpen = now
	backoff := b.computeBackoffLocked()
	b.mu.Unlock()

	retry := now.Add(backoff)
	_ = b.store.SetBreakerState(ctx, b.cfg.DistributedKey, BreakerOpen, retry, err.Error(), b.cfg.DistributedTTL)
}

func (b *HybridCircuitBreaker) transitionLocked(newState BreakerState, reason string) {
	if b.state == newState && newState != BreakerOpen {
		return
	}
	old := b.state
	b.state = newState
	if newState == BreakerOpen {
		b.openedAt = time.Now()
	}
	if b.cfg.OnStateChange != nil {
		go b.cfg.OnStateChange(old, newState, reason)
	}
}

func [T any] (b *HybridCircuitBreaker) handleOpenWithFallback(
	ctx context.Context,
	span trace.Span,
	now time.Time,
	retryAfter time.Time,
	reason string,
	cachedFn func(context.Context) (T, bool, error),
	staleFn func(context.Context) (T, bool, error),
) (T, error) {
	var zero T

	span.SetAttributes(
		attribute.String("breaker.outcome", "open"),
		attribute.String("breaker.reason", reason),
	)

	if b.cfg.EnableCachedOnly && cachedFn != nil {
		res, ok, err := cachedFn(ctx)
		if err == nil && ok {
			span.SetAttributes(attribute.String("breaker.fallback", string(FallbackCached)))
			return res, nil
		}
	}

	if b.cfg.EnableStaleOK && staleFn != nil {
		res, ok, err := staleFn(ctx)
		if err == nil && ok {
			span.SetAttributes(attribute.String("breaker.fallback", string(FallbackStaleOK)))
			return res, nil
		}
	}

	span.SetAttributes(attribute.String("breaker.fallback", string(FallbackNone)))
	return zero, &BreakerOpenError{
		State:      BreakerOpen,
		Reason:     reason,
		RetryAfter: retryAfter,
	}
}

func [T any] (b *HybridCircuitBreaker) handleFailureWithFallback(
	ctx context.Context,
	span trace.Span,
	now time.Time,
	cause error,
	cachedFn func(context.Context) (T, bool, error),
	staleFn func(context.Context) (T, bool, error),
) (T, error) {
	var zero T

	span.SetAttributes(
		attribute.String("breaker.outcome", "failure"),
		attribute.String("breaker.failure.cause", cause.Error()),
	)

	if b.cfg.EnableCachedOnly && cachedFn != nil {
		res, ok, err := cachedFn(ctx)
		if err == nil && ok {
			span.SetAttributes(attribute.String("breaker.fallback", string(FallbackCached)))
			return res, nil
		}
	}

	if b.cfg.EnableStaleOK && staleFn != nil {
		res, ok, err := staleFn(ctx)
		if err == nil && ok {
			span.SetAttributes(attribute.String("breaker.fallback", string(FallbackStaleOK)))
			return res, nil
		}
	}

	span.SetAttributes(attribute.String("breaker.fallback", string(FallbackNone)))
	return zero, cause
}

type BreakerOpenError struct {
	State      BreakerState
	Reason     string
	RetryAfter time.Time
}

func (e *BreakerOpenError) Error() string {
	return "breaker open: " + e.Reason
}

func IsBreakerOpen(err error) bool {
	var boe *BreakerOpenError
	return errors.As(err, &boe)
}
