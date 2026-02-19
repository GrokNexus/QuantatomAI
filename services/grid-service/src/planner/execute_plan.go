package planner

import (
    "context"
    "fmt"

    "quantatomai/grid-service/domain"
)

// -----------------------------
// Experimental/Binary Types (Stubs)
// -----------------------------

// grid package stub for binary/protobuf serialization.
type gridPackage struct{}

type GridResultBinary struct {
    Cells []*Cell
}

type Cell struct {
    RowIndex uint32
    ColIndex uint32
    Value    float64
}

var grid gridPackage

// -----------------------------
// Execution Orchestrator
// -----------------------------

// -----------------------------
// Execution Orchestrator
// -----------------------------

// Planner extensions for caching
func (p *Planner) SetCache(cache ProjectionCache) {
    p.cache = cache
}

var _ Planner // allow planner to have cache field

// ExecutePlan executes a QueryPlan with advanced optimizations:
// 1) Axis Caching to avoid redundant Cartesian joins.
// 2) Dual Output (JSON & Binary) for heterogeneous clients.
// 3) Windowed projection and pre-allocated memory pools.
func (p *Planner) ExecutePlan(
    ctx context.Context,
    plan *QueryPlan,
    fetcher AtomFetcher,
    compute ComputeEngine,
    window ProjectionWindow,
    defaults DefaultValuePolicy,
) (*QueryResult, *GridResultBinary, error) {
    if len(plan.BlockRef.MeasureIDs) == 0 {
        return nil, nil, fmt.Errorf("plan has no measures")
    }
    if len(plan.BlockRef.ScenarioIDs) == 0 {
        return nil, nil, fmt.Errorf("plan has no scenarios")
    }

    primaryMeasure := plan.BlockRef.MeasureIDs[0]
    primaryScenario := plan.BlockRef.ScenarioIDs[0]

    // 1) Fetch raw atom values for the SparsedBlockRef.
    rawValues, err := fetcher.FetchAtoms(ctx, plan.BlockRef)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to fetch atoms: %w", err)
    }

    // 2) Post-process via compute engine (FX, variances, etc.).
    processedValues, err := compute.PostProcess(ctx, rawValues)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to post-process atoms: %w", err)
    }

    // 3) Build or Retrieve Row/Col Axes combinations.
    var rowCombos, colCombos [][]MemberInfo
    var rowIDCombos, colIDCombos [][]int64

    ak := AxisKey{ModelID: "default", ViewID: "default"} // In production, derive from Query context.

    // Use cache if available to avoid expensive Cartesian products.
    if p.cache != nil {
        if r, rIDs, c, cIDs, ok := p.cache.GetAxes(ak); ok {
            rowCombos, rowIDCombos, colCombos, colIDCombos = r, rIDs, c, cIDs
        } else {
            rowCombos, rowIDCombos, _ = MaterializeAxisCombos(plan.RowAxes, 64, true)
            colCombos, colIDCombos, _ = MaterializeAxisCombos(plan.ColAxes, 64, true)
            p.cache.SetAxes(ak, rowCombos, rowIDCombos, colCombos, colIDCombos)
        }
    } else {
        rowCombos, rowIDCombos, _ = MaterializeAxisCombos(plan.RowAxes, 64, true)
        colCombos, colIDCombos, _ = MaterializeAxisCombos(plan.ColAxes, 64, true)
    }

    totalRows := len(rowCombos)
    totalCols := len(colCombos)

    // 4) Normalize projection window.
    if window.RowStart < 0 {
        window.RowStart = 0
    }
    if window.ColStart < 0 {
        window.ColStart = 0
    }
    if window.RowEnd <= 0 || window.RowEnd > totalRows {
        window.RowEnd = totalRows
    }
    if window.ColEnd <= 0 || window.ColEnd > totalCols {
        window.ColEnd = totalCols
    }
    
    if window.RowStart >= window.RowEnd || window.ColStart >= window.ColEnd {
        return &QueryResult{
            Rows:    rowCombos,
            Columns: colCombos,
            Cells:   nil,
        }, &GridResultBinary{}, nil
    }

    // 5) Project atom values into both JSON and Binary formats.
    cellCount := (window.RowEnd - window.RowStart) * (window.ColEnd - window.ColStart)
    cells := make([]CellValue, 0, cellCount)
    binaryCells := make([]*Cell, 0, cellCount)

    rowDimCount := len(plan.RowAxes)
    colDimCount := len(plan.ColAxes)

    for rIdx := window.RowStart; rIdx < window.RowEnd; rIdx++ {
        rowIDs := rowIDCombos[rIdx]

        for cIdx := window.ColStart; cIdx < window.ColEnd; cIdx++ {
            colIDs := colIDCombos[cIdx]

			key := domain.AtomKey{
				DimCount:   len(rowIDs) + len(colIDs),
				MeasureID:  primaryMeasure,
				ScenarioID: primaryScenario,
			}
			copy(key.DimIDs[:], rowIDs)
			copy(key.DimIDs[len(rowIDs):], colIDs)
            key.EnsureCanonical()

            val, ok := processedValues[key]
            if !ok {
                if defaults != nil {
                    val = defaults.DefaultFor(primaryMeasure)
                } else {
                    val = 0.0
                }
            }

            // Populate JSON slice.
            cells = append(cells, CellValue{
                RowIndex: rIdx,
                ColIndex: cIdx,
                Value:    val,
            })

            // Populate Binary slice.
            binaryCells = append(binaryCells, &Cell{
                RowIndex: uint32(rIdx),
                ColIndex: uint32(cIdx),
                Value:    val,
            })
        }
    }

    result := &QueryResult{
        Rows:    rowCombos,
        Columns: colCombos,
        Cells:   cells,
    }

    bin := &GridResultBinary{
        Cells: binaryCells,
    }

    return result, bin, nil
}
