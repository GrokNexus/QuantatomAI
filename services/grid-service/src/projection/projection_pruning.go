package projection

// pruneEmptyColumns removes columns whose ColumnStats show no non-empty cells.
// It rewrites:
//   - Cells (flat array)
//   - RowOffsets
//   - RowMasks (64-bit fast path)
//   - ColumnStats
//
// This function assumes:
//   - ColumnStats has length = projected column count
//   - RowOffsets and RowMasks are already sized to rowOut
//   - Cells is a flat row-major array
//
// This function does NOT modify:
//   - StringArena
//   - StringOffsets
//
// This is a projection-time transformation and must NOT live in grid_model.go.
func pruneEmptyColumns(g *GridResult) {
	colCount := len(g.ColumnStats)
	if colCount == 0 {
		return
	}

	// Determine which columns to keep
	keep := make([]bool, colCount)
	keptCols := 0
	for i := 0; i < colCount; i++ {
		if g.ColumnStats[i].NonEmptyCount > 0 {
			keep[i] = true
			keptCols++
		}
	}

	// Nothing to prune
	if keptCols == colCount {
		return
	}

	// Rebuild ColumnStats
	newStats := make([]ColumnStats, 0, keptCols)
	for i := 0; i < colCount; i++ {
		if keep[i] {
			newStats = append(newStats, g.ColumnStats[i])
		}
	}
	g.ColumnStats = newStats

	// Rebuild Cells + RowMasks
	newCells := make([]Cell, 0, len(g.Cells))
	newMasks := make([]uint64, len(g.RowMasks))

	for r := 0; r < len(g.RowOffsets); r++ {
		oldStart := g.RowOffsets[r]
		// Determine old row end
		oldEnd := int32(len(g.Cells))
		if r+1 < len(g.RowOffsets) {
			oldEnd = g.RowOffsets[r+1]
		}
		_ = oldEnd // Explicitly ignore if not used in current loop logic

		newStart := len(newCells)
		newMask := uint64(0)
		bit := uint(0)

		// Iterate old columns, copy only kept ones
		for oldCol := 0; oldCol < colCount; oldCol++ {
			if !keep[oldCol] {
				continue
			}

			// Boundary check for safety in slice-based access
			idx := int(oldStart) + oldCol
			if idx < len(g.Cells) {
				cell := g.Cells[idx]
				newCells = append(newCells, cell)

				if !cell.IsNull() {
					newMask |= 1 << bit
				}
			}
			bit++
		}

		g.RowOffsets[r] = int32(newStart)
		newMasks[r] = newMask
	}

	g.Cells = newCells
	g.RowMasks = newMasks
}
