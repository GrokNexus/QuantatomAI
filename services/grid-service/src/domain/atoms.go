package domain

import (
    "encoding/binary"
	"hash/fnv"
	"sort"
)

// GridQuery encapsulates the full requested grid shape and members.
type GridQuery struct {
	Dimensions struct {
		Rows    []string            `json:"rows"`
		Columns []string            `json:"columns"`
		Pages   []string            `json:"pages"`
		Filters map[string][]string `json:"filters"`
	} `json:"dimensions"`
	Members map[string][]string `json:"members"`
	Stream  bool                `json:"stream,omitempty"`
}

// ProjectedCell represents a single calculated cell for grid display or streaming.
type ProjectedCell struct {
	RowIndex int     `json:"rowIndex"`
	ColIndex int     `json:"colIndex"`
	Value    float64 `json:"value"`
}

// ProjectedCellResolved represents a cell with pre-resolved metadata labels.
type ProjectedCellResolved struct {
	RowIndex int     `json:"rowIndex"`
	ColIndex int     `json:"colIndex"`
	Value    float64 `json:"value"`

	Currency string `json:"currency,omitempty"`
	Scenario string `json:"scenario,omitempty"`
	Measure  string `json:"measure,omitempty"`
}

type AtomKey struct {
	// Canonically ordered dimension IDs (excluding Measure and Scenario)
	DimIDs   [8]int64
	DimCount int

	// First-class fields
	MeasureID  int64
	ScenarioID int64
}

// EnsureCanonical enforces deterministic ordering of DimIDs.
func (k *AtomKey) EnsureCanonical() {
	if k.DimCount <= 1 {
		return
	}
	dims := k.DimIDs[:k.DimCount]
	sort.Slice(dims, func(i, j int) bool {
		return dims[i] < dims[j]
	})
}

// HashKey returns a compact uint64 hash suitable for Redis/Scylla primary keys.
func (k *AtomKey) HashKey() uint64 {
	k.EnsureCanonical()

	h := fnv.New64a()

	// Hash dimension IDs
	for i := 0; i < k.DimCount; i++ {
		id := k.DimIDs[i]
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(id))
		h.Write(buf[:])
	}

	// Hash measure
	{
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(k.MeasureID))
		h.Write(buf[:])
	}

	// Hash scenario
	{
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(k.ScenarioID))
		h.Write(buf[:])
	}

	return h.Sum64()
}

// AtomWrite represents a writeback operation.
type AtomWrite struct {
    Key   AtomKey
    Value float64
    User  string
}
