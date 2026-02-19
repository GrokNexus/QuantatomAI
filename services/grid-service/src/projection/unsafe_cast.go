package projection

import (
	"unsafe"
)

// CellsToFloat64Slice reinterprets []Cell as []float64 without copying.
func CellsToFloat64Slice(cells []Cell) []float64 {
	if len(cells) == 0 {
		return nil
	}
	return unsafe.Slice((*float64)(unsafe.Pointer(&cells[0])), len(cells))
}

// Float64Bytes returns the raw byte representation of a []float64 slice.
func Float64Bytes(floats []float64) []byte {
	if len(floats) == 0 {
		return nil
	}
	// Note: We use unsafe.Slice for a more modern and safer approach than the array pointer cast
	return unsafe.Slice((*byte)(unsafe.Pointer(&floats[0])), len(floats)*8)
}

// BytesToCellSlice reinterprets []byte as []Cell without copying.
func BytesToCellSlice(b []byte) []Cell {
	if len(b) == 0 {
		return nil
	}
	return unsafe.Slice((*Cell)(unsafe.Pointer(&b[0])), len(b)/8)
}
