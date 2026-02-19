package compute

import (
	"context"
	"fmt"
	"math"

	"quantatomai/grid-service/domain"
)

//
// ─────────────────────────────────────────────────────────────
//   CONFIG
// ─────────────────────────────────────────────────────────────
//

// MissingDataPolicy controls how missing base/compare values are treated during variance calculation.
type MissingDataPolicy int

const (
	// MissingDataSkip: if either base or compare is missing, skip emitting variance.
	MissingDataSkip MissingDataPolicy = iota

	// MissingDataZero: treat missing values as zero.
	MissingDataZero

	// MissingDataEmitNull: reserved for future null-emission semantics.
	MissingDataEmitNull
)

// VarianceConfig defines which measures participate in variance and how.
type VarianceConfig struct {
	// Key: measure ID that should get variance computed.
	Measures map[int64]VarianceMeasureConfig

	// How to treat missing base/compare values.
	MissingData MissingDataPolicy

	// Whether to emit audit records for each variance calculation.
	EnableAudit bool

	// Suppress variance emission if absolute delta is below this threshold.
	NoiseThreshold float64
}

// VarianceMeasureConfig defines base/compare scenarios and which outputs to emit.
type VarianceMeasureConfig struct {
	BaseScenarioID    int64
	CompareScenarioID int64

	// If true, emit absolute delta: (compare - base) * Sign.
	OutputDelta bool

	// If true, emit percent delta: ((compare - base) / base) * Sign.
	OutputPercent bool

	// Financial directionality: +1 for revenue (growth is good), -1 for expense (growth is bad).
	// If left as 0, it defaults to +1.
	Sign int

	// Optional: separate measure IDs for outputs.
	// If zero, the output is written back to the original measure ID but on the compare scenario.
	DeltaMeasureID   int64
	PercentMeasureID int64
}

//
// ─────────────────────────────────────────────────────────────
//   AUDIT STRUCT
// ─────────────────────────────────────────────────────────────
//

// VarianceAudit captures the inputs and result of a specific variance computation.
type VarianceAudit struct {
	BaseScenarioID    int64
	CompareScenarioID int64
	BaseValue         float64
	CompareValue      float64
	Delta             float64
	Percent           float64
}

//
// ─────────────────────────────────────────────────────────────
//   TRANSFORMER STRUCT
// ─────────────────────────────────────────────────────────────
//

// VarianceTransformer computes scenario-based variances for a set of atoms.
type VarianceTransformer struct {
	cfg   VarianceConfig
	audit map[domain.AtomKey]VarianceAudit
}

//
// ─────────────────────────────────────────────────────────────
//   CONSTRUCTOR
// ─────────────────────────────────────────────────────────────
//

// NewVarianceTransformer constructs a variance compute engine with the provided config.
func NewVarianceTransformer(cfg VarianceConfig) *VarianceTransformer {
	// Normalize config: ensure Sign defaults to +1.
	for measureID, m := range cfg.Measures {
		if m.Sign == 0 {
			m.Sign = 1
		}
		cfg.Measures[measureID] = m
	}

	// Validate config before initialization.
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("invalid VarianceConfig: %v", err))
	}

	t := &VarianceTransformer{
		cfg: cfg,
	}
	if cfg.EnableAudit {
		t.audit = make(map[domain.AtomKey]VarianceAudit)
	}
	return t
}

//
// ─────────────────────────────────────────────────────────────
//   CONFIG VALIDATION
// ─────────────────────────────────────────────────────────────
//

// Validate ensures the variance configuration is logically sound.
func (c VarianceConfig) Validate() error {
	for measureID, m := range c.Measures {
		if m.BaseScenarioID == 0 || m.CompareScenarioID == 0 {
			return fmt.Errorf("variance: measure %d has zero base/compare scenario", measureID)
		}
		if !m.OutputDelta && !m.OutputPercent {
			return fmt.Errorf("variance: measure %d has neither delta nor percent enabled", measureID)
		}
	}
	return nil
}

//
// ─────────────────────────────────────────────────────────────
//   NAME
// ─────────────────────────────────────────────────────────────
//

func (t *VarianceTransformer) Name() string {
	return "variance_v1"
}

//
// ─────────────────────────────────────────────────────────────
//   GET AUDIT
// ─────────────────────────────────────────────────────────────
//

// GetAudit returns the audit trail of variance calculations (if enabled).
func (t *VarianceTransformer) GetAudit() map[domain.AtomKey]VarianceAudit {
	return t.audit
}

//
// ─────────────────────────────────────────────────────────────
//   APPLY
// ─────────────────────────────────────────────────────────────
//

// Apply computes variance between scenarios. It is additive; it clones the input
// map and injects the new variance atoms into the result.
func (t *VarianceTransformer) Apply(
	ctx context.Context,
	atoms map[domain.AtomKey]float64,
) (map[domain.AtomKey]float64, error) {

	if len(t.cfg.Measures) == 0 {
		return atoms, nil
	}

	// 1. Clone input (immutability)
	out := cloneAtomMap(atoms)

	// 2. Group values by "non-scenario" coordinate + measure.
	// We use this to find matching pairs across the Base and Compare scenarios.
	type scenarioKey struct {
		BaseKey   domain.AtomKey
		MeasureID int64
	}

	grouped := make(map[scenarioKey]map[int64]float64)

	for key, value := range atoms {
		measureID := key.MeasureID
		_, ok := t.cfg.Measures[measureID]
		if !ok {
			continue
		}

		scenarioID := extractScenario(key)
		if scenarioID == 0 {
			continue
		}

		// Strip scenario to create a stable grouping key.
		baseKey := stripScenario(key)
		sk := scenarioKey{
			BaseKey:   baseKey,
			MeasureID: measureID,
		}

		bucket, exists := grouped[sk]
		if !exists {
			bucket = make(map[int64]float64)
			grouped[sk] = bucket
		}
		bucket[scenarioID] = value
	}

	// 3. Compute Variances
	for sk, byScenario := range grouped {
		cfg := t.cfg.Measures[sk.MeasureID]

		baseVal, hasBase := byScenario[cfg.BaseScenarioID]
		compareVal, hasCompare := byScenario[cfg.CompareScenarioID]

		// Handle missing data according to policy.
		switch t.cfg.MissingData {
		case MissingDataSkip:
			if !hasBase || !hasCompare {
				continue
			}
		case MissingDataZero:
			if !hasBase {
				baseVal = 0
				hasBase = true
			}
			if !hasCompare {
				compareVal = 0
				hasCompare = true
			}
		case MissingDataEmitNull:
			if !hasBase || !hasCompare {
				continue
			}
		}

		// Double-check if we have valid inputs after policy application.
		if !hasBase && !hasCompare {
			continue
		}

		// A. Absolute Delta: (Compare - Base) * Sign
		rawDelta := compareVal - baseVal
		delta := rawDelta * float64(cfg.Sign)

		// Noise suppression.
		if t.cfg.NoiseThreshold > 0 && math.Abs(delta) < t.cfg.NoiseThreshold {
			continue
		}

		// B. Percent Delta: ((Compare - Base) / Base) * Sign
		var percent float64
		hasPercent := false
		if cfg.OutputPercent && hasBase && baseVal != 0 {
			percent = (rawDelta / baseVal) * float64(cfg.Sign)
			hasPercent = true
		}

		// C. Emit Delta Cell
		if cfg.OutputDelta {
			targetMeasureID := cfg.DeltaMeasureID
			if targetMeasureID == 0 {
				targetMeasureID = sk.MeasureID
			}

			deltaKey := withScenario(sk.BaseKey, cfg.CompareScenarioID)
			deltaKey.MeasureID = targetMeasureID
			deltaKey.EnsureCanonical()

			out[deltaKey] = delta

			if t.cfg.EnableAudit {
				t.audit[deltaKey] = VarianceAudit{
					BaseScenarioID:    cfg.BaseScenarioID,
					CompareScenarioID: cfg.CompareScenarioID,
					BaseValue:         baseVal,
					CompareValue:      compareVal,
					Delta:             delta,
					Percent:           0, // will be updated if percent is also written to this key
				}
			}
		}

		// D. Emit Percent Cell
		if cfg.OutputPercent && hasPercent {
			targetMeasureID := cfg.PercentMeasureID
			if targetMeasureID == 0 {
				targetMeasureID = sk.MeasureID
			}

			pctKey := withScenario(sk.BaseKey, cfg.CompareScenarioID)
			pctKey.MeasureID = targetMeasureID
			pctKey.EnsureCanonical()

			out[pctKey] = percent

			if t.cfg.EnableAudit {
				// Update or create audit record for this key.
				audit := t.audit[pctKey]
				audit.BaseScenarioID = cfg.BaseScenarioID
				audit.CompareScenarioID = cfg.CompareScenarioID
				audit.BaseValue = baseVal
				audit.CompareValue = compareVal
				audit.Percent = percent
				t.audit[pctKey] = audit
			}
		}
	}

	return out, nil
}
