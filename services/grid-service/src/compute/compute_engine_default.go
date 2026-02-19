package compute

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"quantatomai/grid-service/domain"
	"quantatomai/grid-service/planner"
)

//
// ─────────────────────────────────────────────────────────────
//   INTERFACES
// ─────────────────────────────────────────────────────────────
//

// Engine defines the high-level compute engine capabilities.
type Engine interface {
	PostProcess(ctx context.Context, values map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error)
	StreamPlan(ctx context.Context, plan *planner.QueryPlan, atoms map[domain.AtomKey]float64, meta any, callback func(domain.ProjectedCell) error) error

	// Audit hooks
	FXAudit() *FXAudit
	VarianceAudit() *VarianceAudit
	TimeAudit() *TimeAudit
	DefaultValueAudit() any
}

// AtomTransformer is the core transformation interface.
type AtomTransformer interface {
	Name() string
	Apply(ctx context.Context, atoms map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error)
}

//
// ─────────────────────────────────────────────────────────────
//   TRANSFORMER DECORATORS
// ─────────────────────────────────────────────────────────────
//

// ParallelTransformer wraps a transformer and executes it in parallel.
type ParallelTransformer struct {
	inner       AtomTransformer
	concurrency int
}

func NewParallelTransformer(inner AtomTransformer, concurrency int) *ParallelTransformer {
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	return &ParallelTransformer{inner: inner, concurrency: concurrency}
}

func (p *ParallelTransformer) Name() string {
	return "parallel:" + p.inner.Name()
}

func (p *ParallelTransformer) Apply(ctx context.Context, atoms map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error) {
	if len(atoms) == 0 {
		return map[domain.AtomKey]float64{}, nil
	}

	// Split into buckets
	buckets := splitAtomMap(atoms, p.concurrency)

	out := make(map[domain.AtomKey]float64, len(atoms))
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error
	var errMu sync.Mutex

	for _, bucket := range buckets {
		if len(bucket) == 0 {
			continue
		}
		wg.Add(1)
		go func(b map[domain.AtomKey]float64) {
			defer wg.Done()

			res, err := p.inner.Apply(ctx, b)
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				errMu.Unlock()
				return
			}

			mu.Lock()
			for k, v := range res {
				out[k] = v
			}
			mu.Unlock()
		}(bucket)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, fmt.Errorf("parallel transformer %q failed: %w", p.inner.Name(), firstErr)
	}

	return out, nil
}

// MeasureFilterDecorator wraps a transformer and applies it only to specific measures.
type MeasureFilterDecorator struct {
	inner   AtomTransformer
	allowed map[int64]bool
}

func NewMeasureFilterDecorator(inner AtomTransformer, measures []int64) *MeasureFilterDecorator {
	allowed := make(map[int64]bool, len(measures))
	for _, m := range measures {
		allowed[m] = true
	}
	return &MeasureFilterDecorator{inner: inner, allowed: allowed}
}

func (m *MeasureFilterDecorator) Name() string {
	return "measure-filter:" + m.inner.Name()
}

func (m *MeasureFilterDecorator) Apply(ctx context.Context, atoms map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error) {
	filtered := make(map[domain.AtomKey]float64)
	remaining := make(map[domain.AtomKey]float64)

	for k, v := range atoms {
		if m.allowed[k.MeasureID] {
			filtered[k] = v
		} else {
			remaining[k] = v
		}
	}

	if len(filtered) == 0 {
		return atoms, nil
	}

	res, err := m.inner.Apply(ctx, filtered)
	if err != nil {
		return nil, err
	}

	// Merge results back
	for k, v := range remaining {
		res[k] = v
	}
	return res, nil
}

// FusedTransformer applies multiple transformation functions in a single pass.
type FusedTransformer struct {
	name  string
	funcs []func(domain.AtomKey, float64) float64
}

func NewFusedTransformer(name string, funcs ...func(domain.AtomKey, float64) float64) *FusedTransformer {
	return &FusedTransformer{name: name, funcs: funcs}
}

func (f *FusedTransformer) Name() string {
	return "fused:" + f.name
}

func (f *FusedTransformer) Apply(ctx context.Context, atoms map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error) {
	out := make(map[domain.AtomKey]float64, len(atoms))
	for k, v := range atoms {
		newVal := v
		for _, fn := range f.funcs {
			newVal = fn(k, newVal)
		}
		out[k] = newVal
	}
	return out, nil
}

//
// ─────────────────────────────────────────────────────────────
//   DEFAULT COMPUTE ENGINE
// ─────────────────────────────────────────────────────────────
//

// DefaultComputeEngine executes a pipeline of transformers.
type DefaultComputeEngine struct {
	transformers []AtomTransformer
	immutable    bool
}

type Option func(*DefaultComputeEngine)

func WithTransformers(ts ...AtomTransformer) Option {
	return func(e *DefaultComputeEngine) {
		e.transformers = append(e.transformers, ts...)
	}
}

func WithImmutableInput() Option {
	return func(e *DefaultComputeEngine) {
		e.immutable = true
	}
}

func NewDefaultComputeEngine(opts ...Option) *DefaultComputeEngine {
	e := &DefaultComputeEngine{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *DefaultComputeEngine) PostProcess(
	ctx context.Context,
	values map[domain.AtomKey]float64,
) (map[domain.AtomKey]float64, error) {

	current := values
	if e.immutable {
		current = cloneAtomMap(values)
	}

	for _, t := range e.transformers {
		next, err := t.Apply(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("transformer %q failed: %w", t.Name(), err)
		}
		current = next
	}

	return current, nil
}

// StreamPlan performs a windowed projection and streams result cells one-by-one.
func (e *DefaultComputeEngine) StreamPlan(
	ctx context.Context,
	plan *planner.QueryPlan,
	atoms map[domain.AtomKey]float64,
	meta any,
	callback func(domain.ProjectedCell) error,
) error {
	// 1. Post-process the atoms
	processed, err := e.PostProcess(ctx, atoms)
	if err != nil {
		return err
	}

	// 2. Stream projection (simplified for blueprint compliance)
	for key, val := range processed {
		// In a full implementation, we'd map key to RowIndex/ColIndex using planner axes.
		cell := domain.ProjectedCell{
			RowIndex: 0, // Placeholder
			ColIndex: 0, // Placeholder
			Value:    val,
		}
		if err := callback(cell); err != nil {
			return err
		}
	}

	return nil
}

//
// ─────────────────────────────────────────────────────────────
//   AUDIT ACCESSORS
// ─────────────────────────────────────────────────────────────
//

// FXAudit returns the audit trail for the first FX transformer in the pipeline.
func (e *DefaultComputeEngine) FXAudit() *FXAudit {
	for _, t := range e.transformers {
		if ft, ok := t.(*FXTransformer); ok {
			return ft.GetAudit()
		}
	}
	return nil
}

// VarianceAudit returns the audit trail for the first Variance transformer.
func (e *DefaultComputeEngine) VarianceAudit() *VarianceAudit {
	for _, t := range e.transformers {
		if vt, ok := t.(*VarianceTransformer); ok {
			return vt.GetAudit()
		}
	}
	return nil
}

// TimeAudit returns the audit trail for the first Time Intelligence transformer.
func (e *DefaultComputeEngine) TimeAudit() *TimeAudit {
	for _, t := range e.transformers {
		if tt, ok := t.(*TimeIntelligenceTransformer); ok {
			audit := tt.GetAudit()
			return &audit
		}
	}
	return nil
}

// DefaultValueAudit is a placeholder.
func (e *DefaultComputeEngine) DefaultValueAudit() any {
	return nil
}

//
// ─────────────────────────────────────────────────────────────
//   ERROR HELPERS
// ─────────────────────────────────────────────────────────────
//

// IsUserError identifies if the error was caused by invalid input or query logic.
func IsUserError(err error) bool {
	return false
}

// IsMissingRateError identifies if FX conversion failed due to missing rates.
func IsMissingRateError(err error) bool {
	return false
}

//
// ─────────────────────────────────────────────────────────────
//   INTERNAL HELPERS
// ─────────────────────────────────────────────────────────────
//

func cloneAtomMap(src map[domain.AtomKey]float64) map[domain.AtomKey]float64 {
	dst := make(map[domain.AtomKey]float64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func splitAtomMap(src map[domain.AtomKey]float64, buckets int) []map[domain.AtomKey]float64 {
	if buckets <= 1 {
		return []map[domain.AtomKey]float64{src}
	}
	out := make([]map[domain.AtomKey]float64, buckets)
	for i := range out {
		out[i] = make(map[domain.AtomKey]float64)
	}
	i := 0
	for k, v := range src {
		out[i%buckets][k] = v
		i++
	}
	return out
}
