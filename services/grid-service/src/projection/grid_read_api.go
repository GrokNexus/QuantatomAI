package projection

// CellAt returns the Cell at logical (row, col) using RowOffsets.
// It assumes a dense row-major layout after projection/pruning.
func (g *GridResult) CellAt(row, col int) (Cell, bool) {
	if row < 0 || row >= len(g.RowOffsets) {
		return Cell{}, false
	}
	start := int(g.RowOffsets[row])
	// Determine old row end
	end := len(g.Cells)
	if row+1 < len(g.RowOffsets) {
		end = int(g.RowOffsets[row+1])
	}
	if col < 0 || start+col >= end {
		return Cell{}, false
	}
	return g.Cells[start+col], true
}

// NumericAt returns the numeric value at (row, col) if the cell is numeric.
func (g *GridResult) NumericAt(row, col int) (float64, bool) {
	cell, ok := g.CellAt(row, col)
	if !ok || cell.IsNull() || cell.IsString() {
		return 0, false
	}
	return cell.AsNumeric(), true
}

// StringAt returns the string value at (row, col) if the cell is string-backed.
func (g *GridResult) StringAt(row, col int) (string, bool) {
	cell, ok := g.CellAt(row, col)
	if !ok || !cell.IsString() {
		return "", false
	}
	idx := cell.StringIndex()
	if int(idx) < 0 || int(idx) >= len(g.StringOffsets) {
		return "", false
	}
	offset := g.StringOffsets[idx]
	var end uint32
	if int(idx)+1 < len(g.StringOffsets) {
		end = g.StringOffsets[int(idx)+1]
	} else {
		end = uint32(len(g.StringArena))
	}
	if int(end) > len(g.StringArena) {
		return "", false
	}
	return string(g.StringArena[offset:end]), true
}

// HasDataFast uses RowMasks (if available) to quickly check if (row, col) is non-empty.
// Falls back to CellAt when masks are not applicable.
func (g *GridResult) HasDataFast(row, col int) bool {
	if col >= 0 && col < 64 && row >= 0 && row < len(g.RowMasks) {
		return (g.RowMasks[row] & (1 << uint(col))) != 0
	}
	cell, ok := g.CellAt(row, col)
	if !ok {
		return false
	}
	return !cell.IsNull()
}
