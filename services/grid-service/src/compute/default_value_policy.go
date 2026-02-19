package compute

import (
	"fmt"

	"quantatomai/grid-service/domain"
)

//
// ─────────────────────────────────────────────────────────────
//   INTERFACE
// ─────────────────────────────────────────────────────────────
//

// DefaultValuePolicy defines how missing atoms are treated at read time.
type DefaultValuePolicy interface {
	// ResolveDefault returns the value to use when an atom is missing.
	// ok = false means "no default, treat as missing".
	ResolveDefault(key domain.AtomKey) (value float64, ok bool)

	// Name is useful for logging / diagnostics.
	Name() string
}

// DefaultValueAudit captures when a value was synthesized by a policy.
type DefaultValueAudit struct {
	Key       domain.AtomKey
	Policy    string
	Value     float64
	IsDefault bool
}

//
// ─────────────────────────────────────────────────────────────
//   CONFIG
// ─────────────────────────────────────────────────────────────
//

type DefaultValuePolicyKind int

const (
	DefaultValueNone             DefaultValuePolicyKind = iota // no default, missing stays missing
	DefaultValueZero                                           // missing → 0
	DefaultValueMeasureSpecific                                // per-measure defaults
	DefaultValueScenarioSpecific                               // per-scenario defaults
	DefaultValueComposite                                      // chain of policies
	DefaultValueIdentity                                       // missing → 1.0 (for Multiplicative measures)
)

type DefaultValueConfig struct {
	Kind DefaultValuePolicyKind

	// Optional: per-measure defaults (used when Kind == DefaultValueMeasureSpecific).
	MeasureDefaults map[int64]float64

	// Optional: per-scenario defaults (used when Kind == DefaultValueScenarioSpecific).
	ScenarioDefaults map[int64]float64

	// Optional: composite chain (used when Kind == DefaultValueComposite).
	// Example: [MeasureSpecific, ScenarioSpecific, Zero]
	CompositeKinds []DefaultValuePolicyKind
}

//
// ─────────────────────────────────────────────────────────────
//   CONCRETE POLICIES
// ─────────────────────────────────────────────────────────────
//

// noneDefaultValuePolicy: no default, missing stays missing.
type noneDefaultValuePolicy struct{}

func (p *noneDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	return 0, false
}
func (p *noneDefaultValuePolicy) Name() string { return "none" }

// zeroDefaultValuePolicy: missing → 0 (singleton).
type zeroDefaultValuePolicy struct{}

var zeroPolicyInstance = &zeroDefaultValuePolicy{}

func (p *zeroDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	return 0, true
}
func (p *zeroDefaultValuePolicy) Name() string { return "zero" }

// identityDefaultValuePolicy: missing → 1.0 (singleton).
type identityDefaultValuePolicy struct{}

var identityPolicyInstance = &identityDefaultValuePolicy{}

func (p *identityDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	return 1.0, true
}
func (p *identityDefaultValuePolicy) Name() string { return "identity" }

// measureDefaultValuePolicy: per-measure defaults.
type measureDefaultValuePolicy struct {
	defaults map[int64]float64
}

func (p *measureDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	v, ok := p.defaults[key.MeasureID]
	return v, ok
}
func (p *measureDefaultValuePolicy) Name() string { return "measure_specific" }

// scenarioDefaultValuePolicy: per-scenario defaults.
type scenarioDefaultValuePolicy struct {
	defaults map[int64]float64
}

func (p *scenarioDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	v, ok := p.defaults[key.ScenarioID]
	return v, ok
}
func (p *scenarioDefaultValuePolicy) Name() string { return "scenario_specific" }

// compositeDefaultValuePolicy: chains multiple policies.
type compositeDefaultValuePolicy struct {
	policies []DefaultValuePolicy
}

func (p *compositeDefaultValuePolicy) ResolveDefault(key domain.AtomKey) (float64, bool) {
	for _, policy := range p.policies {
		if v, ok := policy.ResolveDefault(key); ok {
			return v, true
		}
	}
	return 0, false
}
func (p *compositeDefaultValuePolicy) Name() string { return "composite" }

//
// ─────────────────────────────────────────────────────────────
//   FACTORY + VALIDATION
// ─────────────────────────────────────────────────────────────
//

func NewDefaultValuePolicy(cfg DefaultValueConfig) (DefaultValuePolicy, error) {
	switch cfg.Kind {
	case DefaultValueNone:
		return &noneDefaultValuePolicy{}, nil

	case DefaultValueZero:
		return zeroPolicyInstance, nil

	case DefaultValueIdentity:
		return identityPolicyInstance, nil

	case DefaultValueMeasureSpecific:
		if cfg.MeasureDefaults == nil {
			return nil, fmt.Errorf("default value policy: MeasureDefaults is nil")
		}
		return &measureDefaultValuePolicy{
			defaults: cfg.MeasureDefaults,
		}, nil

	case DefaultValueScenarioSpecific:
		if cfg.ScenarioDefaults == nil {
			return nil, fmt.Errorf("default value policy: ScenarioDefaults is nil")
		}
		return &scenarioDefaultValuePolicy{
			defaults: cfg.ScenarioDefaults,
		}, nil

	case DefaultValueComposite:
		if len(cfg.CompositeKinds) == 0 {
			return nil, fmt.Errorf("default value policy: CompositeKinds is empty")
		}
		policies := make([]DefaultValuePolicy, 0, len(cfg.CompositeKinds))
		for _, kind := range cfg.CompositeKinds {
			subCfg := DefaultValueConfig{
				Kind:             kind,
				MeasureDefaults:  cfg.MeasureDefaults,
				ScenarioDefaults: cfg.ScenarioDefaults,
			}
			p, err := NewDefaultValuePolicy(subCfg)
			if err != nil {
				return nil, fmt.Errorf("default value policy: composite sub-policy error: %w", err)
			}
			policies = append(policies, p)
		}
		return &compositeDefaultValuePolicy{policies: policies}, nil

	default:
		return nil, fmt.Errorf("default value policy: unknown kind %d", cfg.Kind)
	}
}

//
// ─────────────────────────────────────────────────────────────
//   AUDIT HOOK (OPTIONAL INTEGRATION PATTERN)
// ─────────────────────────────────────────────────────────────
//

// ApplyDefaultWithAudit is a helper pattern you can use at read time:
// 1) Try to read from atoms.
// 2) If missing, consult policy.
// 3) If defaulted, record audit.
func ApplyDefaultWithAudit(
	atoms map[domain.AtomKey]float64,
	key domain.AtomKey,
	policy DefaultValuePolicy,
	auditSink func(DefaultValueAudit),
) (float64, bool) {

	if v, ok := atoms[key]; ok {
		return v, true
	}

	if policy == nil {
		return 0, false
	}

	v, ok := policy.ResolveDefault(key)
	if !ok {
		return 0, false
	}

	if auditSink != nil {
		auditSink(DefaultValueAudit{
			Key:       key,
			Policy:    policy.Name(),
			Value:     v,
			IsDefault: true,
		})
	}

	return v, true
}
