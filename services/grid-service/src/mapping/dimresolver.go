package mapping

import (
    "fmt"
    "sort"

    "quantatomai/grid-service/domain"
)

type DimResolver interface {
    ResolveDimID(dimName, memberCode string) (int64, error)
    ResolveMeasureID(measure string) (int64, error)
    ResolveScenarioID(scenario string) (int64, error)
}

// StaticDimResolver is a placeholder implementation.
// In production, replace this with a resolver backed by metadata service / Postgres.
type StaticDimResolver struct{}

func NewStaticDimResolver() *StaticDimResolver {
    return &StaticDimResolver{}
}

func (r *StaticDimResolver) ResolveDimID(dimName, memberCode string) (int64, error) {
    // TODO: look up from metadata service / Postgres
    // This is intentionally a stub; the interface is production-ready.
    return 1, nil
}

func (r *StaticDimResolver) ResolveMeasureID(measure string) (int64, error) {
    // TODO: look up from metadata
    return 10, nil
}

func (r *StaticDimResolver) ResolveScenarioID(scenario string) (int64, error) {
    // TODO: look up from metadata
    return 20, nil
}

// BuildAtomKeyFromEdit converts a UI-level cell edit into a canonical AtomKey.
// It guarantees deterministic dimension ordering and respects Measure/Scenario
// as first-class fields.
func BuildAtomKeyFromEdit(
    resolver DimResolver,
    dims map[string]string,
    measure string,
    scenario string,
) (domain.AtomKey, error) {
    if len(dims) == 0 {
        return domain.AtomKey{}, fmt.Errorf("no dimensions provided for edit")
    }

    // 1) Canonicalize dimension order by sorting on dimension name.
    type dimPair struct {
        Name  string
        Value string
    }

    pairs := make([]dimPair, 0, len(dims))
    for dimName, memberCode := range dims {
        pairs = append(pairs, dimPair{Name: dimName, Value: memberCode})
    }

    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].Name < pairs[j].Name
    })

    // 2) Resolve dimension member IDs in canonical order.
    dimIDs := make([]int64, 0, len(pairs))
    for _, p := range pairs {
        id, err := resolver.ResolveDimID(p.Name, p.Value)
        if err != nil {
            return domain.AtomKey{}, fmt.Errorf("unknown member %s:%s", p.Name, p.Value)
        }
        dimIDs = append(dimIDs, id)
    }

    // 3) Resolve measure and scenario IDs.
    measureID, err := resolver.ResolveMeasureID(measure)
    if err != nil {
        return domain.AtomKey{}, fmt.Errorf("unknown measure %s", measure)
    }

    scenarioID, err := resolver.ResolveScenarioID(scenario)
    if err != nil {
        return domain.AtomKey{}, fmt.Errorf("unknown scenario %s", scenario)
    }

	// 4) Build AtomKey.
	if len(dimIDs) > 8 {
		return domain.AtomKey{}, fmt.Errorf("too many dimensions: %d (max 8)", len(dimIDs))
	}

	key := domain.AtomKey{
		DimCount:   len(dimIDs),
		MeasureID:  measureID,
		ScenarioID: scenarioID,
	}
	copy(key.DimIDs[:], dimIDs)

    // Optional: enforce canonical ordering immediately for safety.
    key.EnsureCanonical()

    return key, nil
}
