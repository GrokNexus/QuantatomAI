package projection

import "regexp"

// Bit flags for boolean options.
const (
	FlagSkipEmptyRows uint8 = 1 << iota
	FlagSkipEmptyCols
	FlagIncludeMetadata
	FlagIncludeColumnStats
)

type ProjectionOptions struct {
	// Column selection
	Columns []int // nil => all columns

	// Row/Column limits
	MaxRows int
	MaxCols int

	// Region of Interest (ROI) slicing
	OffsetRows int // starting row index
	OffsetCols int // starting column index

	// Precision control
	Precision int // -1 => no rounding

	// Packed flags
	Flags uint8

	// Optional regex filter (precompiled)
	CellFilter *regexp.Regexp

	// Optional aggregation overrides (future use)
	Aggregation map[int]AggregationType
}

type AggregationType uint8

const (
	AggNone AggregationType = iota
	AggSum
	AggAvg
	AggMin
	AggMax
)

// DefaultOptions returns a safe, optimized default configuration.
func DefaultOptions() ProjectionOptions {
	return ProjectionOptions{
		Columns:     nil,
		MaxRows:     0,
		MaxCols:     0,
		OffsetRows:  0,
		OffsetCols:  0,
		Precision:   -1,
		Flags:       0,
		CellFilter:  nil,
		Aggregation: nil,
	}
}

// Helper methods for flags
func (o ProjectionOptions) SkipEmptyRows() bool {
	return o.Flags&FlagSkipEmptyRows != 0
}

func (o ProjectionOptions) SkipEmptyCols() bool {
	return o.Flags&FlagSkipEmptyCols != 0
}

func (o ProjectionOptions) IncludeMetadata() bool {
	return o.Flags&FlagIncludeMetadata != 0
}

func (o ProjectionOptions) IncludeColumnStats() bool {
	return o.Flags&FlagIncludeColumnStats != 0
}
