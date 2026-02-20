package mdfv1

import (
	"time"
)

// Molecule represents the atomic unit of data in QuantatomAI (MDF).
// This struct is designed to map cleanly to Parquet columns.
type Molecule struct {
	// 1. The Coordinates (The "Where")
	CoordinateHash   []byte            `parquet:"coordinate_hash"`   // 128-bit MurmurHash3
	CustomDimensions map[string]string `parquet:"custom_dimensions"` // Map for sparse dims

	// 2. The Payload (The "What") - Polymorphic
	// Parquet handles these as separate columns, usually null if not set.
	NumericValue    *float64 `parquet:"numeric_value,optional"`
	TextCommentary  *string  `parquet:"text_commentary,optional"`
	EmbeddingVector []byte   `parquet:"embedding_vector,optional"` // Blob for simplicity in Parquet

	// Ultra Diamond Upgrade: Rich Types
	DateValue    *int64  `parquet:"date_value,optional"` // Unix Millis
	BooleanValue *bool   `parquet:"boolean_value,optional"`
	ErrorValue   *string `parquet:"error_value,optional"`

	// 3. The Context (The "Why")
	Timestamp      int64  `parquet:"timestamp"`
	SourceSystem   string `parquet:"source_system"`
	SecurityMask   uint64 `parquet:"security_mask"`
	CausalityClock []byte `parquet:"causality_clock"`
	IsLocked       bool   `parquet:"is_locked"` // Prevent Top-Down Overwrites
}

// NewNumericMolecule creates a standard numeric data point.
func NewNumericMolecule(hash []byte, dims map[string]string, val float64, source string, mask uint64) *Molecule {
	v := val
	return &Molecule{
		CoordinateHash:   hash,
		CustomDimensions: dims,
		NumericValue:     &v,
		Timestamp:        time.Now().UnixNano(),
		SourceSystem:     source,
		SecurityMask:     mask,
	}
}

// NewCommentaryMolecule creates a text annotation.
func NewCommentaryMolecule(hash []byte, dims map[string]string, comment string, source string, mask uint64) *Molecule {
	c := comment
	return &Molecule{
		CoordinateHash:   hash,
		CustomDimensions: dims,
		TextCommentary:   &c,
		Timestamp:        time.Now().UnixNano(),
		SourceSystem:     source,
		SecurityMask:     mask,
	}
}
