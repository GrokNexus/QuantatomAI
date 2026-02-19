package projection

type ColumnStats struct {
	NonEmptyCount int
	NumericCount  int
	StringCount   int
	NullCount     int

	// Optional: min/max for numeric columns
	MinNumeric float64
	MaxNumeric float64
	HasMinMax  bool
}
