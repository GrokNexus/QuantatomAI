package storage

import (
	"fmt"
	"io"
	"sync" // Ultra Diamond: Thread Safety

	"github.com/segmentio/parquet-go"
	mdfv1 "quantatomai/grid-service/pkg/mdf/v1"
)

// MdfWriter handles the serialization of Molecules into Parquet streams.
// It is designed for high-throughput appending of data.
type MdfWriter struct {
	writer *parquet.GenericWriter[mdfv1.Molecule]
	mu     sync.Mutex // Protects concurrent writes
}

// NewMdfWriter creates a new writer that streams Parquet data to the provided io.Writer.
// The writer uses the schema defined in mdfv1.Molecule.
func NewMdfWriter(w io.Writer) *MdfWriter {
	// Configure writer for high throughput (zstd compression, row group size)
	config := parquet.WriterConfig{
		Compression: &parquet.Zstd,
		PageSize:    8 * 1024, // 8KB pages
	}

	pw := parquet.NewGenericWriter[mdfv1.Molecule](w, config)
	
	return &MdfWriter{
		writer: pw,
	}
}

// Write appends a single molecule to the current row group.
// This is thread-safe.
func (mw *MdfWriter) Write(mol *mdfv1.Molecule) error {
	if mol == nil {
		return fmt.Errorf("cannot write nil molecule")
	}

	mw.mu.Lock()
	defer mw.mu.Unlock()

	// Write the molecule directly. Parquet-go handles the struct mapping.
	// Since we are writing *Molecule, generic writer expects []Molecule or explicit Write call
	// For GenericWriter[T], Write takes a slice of T.
	n, err := mw.writer.Write([]mdfv1.Molecule{*mol})
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("expected to write 1 record, wrote %d", n)
	}

	return nil
}

// WriteBatch writes multiple molecules at once for better throughput.
func (mw *MdfWriter) WriteBatch(mols []mdfv1.Molecule) error {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	
	n, err := mw.writer.Write(mols)
	if err != nil {
		return err
	}
	if n != len(mols) {
		return fmt.Errorf("partial write: expected %d, wrote %d", len(mols), n)
	}
	return nil
}

// Close flushes any remaining data and closes the underlying parquet writer.
// It does NOT close the underlying io.Writer.
func (mw *MdfWriter) Close() error {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return mw.writer.Close()
}

// NewMdfWriter creates a new writer that streams Parquet data to the provided io.Writer.
// The writer uses the schema defined in mdfv1.Molecule.
func NewMdfWriter(w io.Writer) *MdfWriter {
	// Configure writer for high throughput (zstd compression, row group size)
	config := parquet.WriterConfig{
		Compression: &parquet.Zstd,
		PageSize:    8 * 1024, // 8KB pages
	}

	pw := parquet.NewGenericWriter[mdfv1.Molecule](w, config)
	
	return &MdfWriter{
		writer: pw,
	}
}

// Write appends a single molecule to the current row group.
// This is thread-safe if the underlying parquet library supports it, but typical usage is single-producer.
func (mw *MdfWriter) Write(mol *mdfv1.Molecule) error {
	if mol == nil {
		return fmt.Errorf("cannot write nil molecule")
	}

	// Write the molecule directly. Parquet-go handles the struct mapping.
	// Since we are writing *Molecule, generic writer expects []Molecule or explicit Write call
	// For GenericWriter[T], Write takes a slice of T.
	
	// We wrap in a slice of 1 for now, but batching in the caller is recommended.
	n, err := mw.writer.Write([]mdfv1.Molecule{*mol})
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("expected to write 1 record, wrote %d", n)
	}

	return nil
}

// WriteBatch writes multiple molecules at once for better throughput.
func (mw *MdfWriter) WriteBatch(mols []mdfv1.Molecule) error {
	n, err := mw.writer.Write(mols)
	if err != nil {
		return err
	}
	if n != len(mols) {
		return fmt.Errorf("partial write: expected %d, wrote %d", len(mols), n)
	}
	return nil
}

// Close flushes any remaining data and closes the underlying parquet writer.
// It does NOT close the underlying io.Writer.
func (mw *MdfWriter) Close() error {
	return mw.writer.Close()
}
