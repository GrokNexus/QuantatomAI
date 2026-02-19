package projection

import (
	"context"
	"fmt"
	"math"
)

type AtomEngine interface {
	EvalCell(ctx context.Context, planID, viewID, windowHash string, row, col int) (float64, string, bool)
}

type DefaultProjectionEngine struct {
	atomEngine AtomEngine
}

func NewDefaultProjectionEngine(atomEngine AtomEngine) *DefaultProjectionEngine {
	return &DefaultProjectionEngine{atomEngine: atomEngine}
}

func (p *DefaultProjectionEngine) ProjectGrid(
	ctx context.Context,
	planID, viewID, windowHash string,
	target *GridResult,
	opts ProjectionOptions,
) error {
	// TODO: derive from plan/metadata
	totalRows := 1000
	totalCols := 100

	// 1. Calculate ROI Bounds
	startRow := opts.OffsetRows
	if startRow < 0 {
		startRow = 0
	}
	startCol := opts.OffsetCols
	if startCol < 0 {
		startCol = 0
	}

	rowCount := totalRows - startRow
	if opts.MaxRows > 0 && opts.MaxRows < rowCount {
		rowCount = opts.MaxRows
	}
	if rowCount < 0 {
		rowCount = 0
	}

	// 2. Column Mapping (ROI + Selection)
	var cols []int
	if len(opts.Columns) > 0 {
		// Use specific columns, but apply OffsetCols
		for i := startCol; i < len(opts.Columns); i++ {
			cols = append(cols, opts.Columns[i])
		}
	} else {
		// Use sequential columns starting from OffsetCols
		endCol := totalCols
		if opts.MaxCols > 0 && startCol+opts.MaxCols < endCol {
			endCol = startCol + opts.MaxCols
		}
		for i := startCol; i < endCol; i++ {
			cols = append(cols, i)
		}
	}
	colCount := len(cols)

	// 3. Metadata
	target.PlanID = planID
	target.ViewID = viewID
	target.WindowHash = windowHash

	// 4. Capacity & Reset
	neededCells := rowCount * colCount
	if cap(target.Cells) < neededCells {
		target.Cells = make([]Cell, 0, neededCells)
	} else {
		target.Cells = target.Cells[:0]
	}

	if cap(target.RowOffsets) < rowCount {
		target.RowOffsets = make([]int32, rowCount)
	} else {
		target.RowOffsets = target.RowOffsets[:rowCount]
	}

	if cap(target.RowMasks) < rowCount {
		target.RowMasks = make([]uint64, rowCount)
	} else {
		target.RowMasks = target.RowMasks[:rowCount]
	}

	// Allocate/Reset ColumnStats
	if cap(target.ColumnStats) < colCount {
		target.ColumnStats = make([]ColumnStats, colCount)
	} else {
		target.ColumnStats = target.ColumnStats[:colCount]
		for i := range target.ColumnStats {
			target.ColumnStats[i] = ColumnStats{}
		}
	}

	target.StringArena = target.StringArena[:0]
	target.StringOffsets = target.StringOffsets[:0]

	// String arena helper
	stringIndex := func(s string) uint32 {
		if s == "" {
			return 0
		}
		offset := uint32(len(target.StringArena))
		target.StringArena = append(target.StringArena, s...)
		target.StringOffsets = append(target.StringOffsets, offset)
		return uint32(len(target.StringOffsets) - 1)
	}

	// 5. Projection Loop
	rowOut := 0
	for r := 0; r < rowCount; r++ {
		globalRow := startRow + r
		rowStart := len(target.Cells)
		emptyRow := true
		var rowMask uint64 = 0
		trackMask := colCount <= 64

		for ci := 0; ci < colCount; ci++ {
			globalCol := cols[ci]

			num, str, ok := p.atomEngine.EvalCell(ctx, planID, viewID, windowHash, globalRow, globalCol)
			if !ok {
				target.Cells = append(target.Cells, EncodeNull())
				target.ColumnStats[ci].NullCount++
				continue
			}

			emptyRow = false

			// Update ColumnStats
			cs := &target.ColumnStats[ci]
			cs.NonEmptyCount++

			if str != "" {
				cs.StringCount++
			} else {
				cs.NumericCount++
				if !cs.HasMinMax {
					cs.MinNumeric = num
					cs.MaxNumeric = num
					cs.HasMinMax = true
				} else {
					if num < cs.MinNumeric {
						cs.MinNumeric = num
					}
					if num > cs.MaxNumeric {
						cs.MaxNumeric = num
					}
				}
			}

			// Regex filter (On-the-fly)
			if opts.CellFilter != nil {
				var candidate string
				if str != "" {
					candidate = str
				} else {
					candidate = formatNumericForFilter(num, opts.Precision)
				}
				if !opts.CellFilter.MatchString(candidate) {
					target.Cells = append(target.Cells, EncodeNull())
					continue
				}
			}

			// Precision rounding
			if opts.Precision >= 0 && str == "" {
				f := math.Pow10(opts.Precision)
				num = math.Round(num*f) / f
			}

			// Bitmask
			if trackMask {
				rowMask |= 1 << uint(ci)
			}

			// Encode
			if str != "" {
				idx := stringIndex(str)
				target.Cells = append(target.Cells, EncodeStringIndex(idx))
			} else {
				target.Cells = append(target.Cells, EncodeNumeric(num))
			}
		}

		// Sparsity Logic (Zero-Alloc Rollback)
		if opts.SkipEmptyRows() && emptyRow {
			target.Cells = target.Cells[:rowStart]
			continue
		}

		target.RowOffsets[rowOut] = int32(rowStart)
		if trackMask {
			target.RowMasks[rowOut] = rowMask
		}
		rowOut++
	}

	// Shrink buffers to actual row count
	target.RowOffsets = target.RowOffsets[:rowOut]
	target.RowMasks = target.RowMasks[:rowOut]

	// 6. SkipEmptyCols logic
	if opts.SkipEmptyCols() {
		pruneEmptyColumns(target)
	}

	return nil
}

// formatNumericForFilter provides a string representation of numbers for regex matching.
func formatNumericForFilter(num float64, precision int) string {
	if precision >= 0 {
		format := fmt.Sprintf("%%.%df", precision)
		return fmt.Sprintf(format, num)
	}
	return fmt.Sprintf("%v", num)
}
