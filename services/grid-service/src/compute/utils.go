package compute

import (
	"quantatomai/grid-service/domain"
)

//
// ─────────────────────────────────────────────────────────────
//   MAP UTILITIES
// ─────────────────────────────────────────────────────────────
//

// cloneAtomMap performs a shallow copy of the atom map.
// This is used for "Copy-on-Write" semantics in transformers.
func cloneAtomMap(src map[domain.AtomKey]float64) map[domain.AtomKey]float64 {
	dst := make(map[domain.AtomKey]float64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

//
// ─────────────────────────────────────────────────────────────
//   SCENARIO HELPERS
// ─────────────────────────────────────────────────────────────
//

// extractScenario returns the scenario ID from an AtomKey.
func extractScenario(key domain.AtomKey) int64 {
	return key.ScenarioID
}

// stripScenario returns a copy of the key with the scenario ID zeroed out.
// This is useful for grouping values across different scenarios.
func stripScenario(key domain.AtomKey) domain.AtomKey {
	k := key
	k.ScenarioID = 0
	return k
}

// withScenario returns a copy of the key with a specific scenario ID applied.
func withScenario(key domain.AtomKey, scenarioID int64) domain.AtomKey {
	k := key
	k.ScenarioID = scenarioID
	return k
}

//
// ─────────────────────────────────────────────────────────────
//   TIME HELPERS
// ─────────────────────────────────────────────────────────────
//

// extractTimeKey returns the period ID from an AtomKey.
// By convention, we assume the last dimension ID is the time period if not otherwise specified.
// In a prod system, this index would be provided by metadata.
func extractTimeKey(key domain.AtomKey, timeIndex int) int64 {
	if timeIndex < 0 || timeIndex >= len(key.DimIDs) {
		return 0
	}
	return key.DimIDs[timeIndex]
}

// stripTimeKey returns a copy of the key with the time dimension zeroed out.
func stripTimeKey(key domain.AtomKey, timeIndex int) domain.AtomKey {
	k := key
	if timeIndex >= 0 && timeIndex < len(k.DimIDs) {
		k.DimIDs = make([]int64, len(key.DimIDs))
		copy(k.DimIDs, key.DimIDs)
		k.DimIDs[timeIndex] = 0
	}
	return k
}

// withTimeKey returns a copy of the key with a specific period ID applied.
func withTimeKey(key domain.AtomKey, timeIndex int, periodID int64) domain.AtomKey {
	k := key
	if timeIndex >= 0 && timeIndex < len(k.DimIDs) {
		k.DimIDs = make([]int64, len(key.DimIDs))
		copy(k.DimIDs, key.DimIDs)
		k.DimIDs[timeIndex] = periodID
	}
	return k
}
