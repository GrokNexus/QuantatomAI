package compute

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"

	"quantatomai/grid-service/domain"
	"quantatomai/grid-service/planner"
)

// FastMetadata contains pre-serialized JSON fragments for metadata labels.
type FastMetadata struct {
	MeasureJSON  map[int64][]byte
	ScenarioJSON map[int64][]byte
	CurrencyJSON map[int64][]byte
}

// ProjectionEngine implements the Diamond-tier high-performance streaming grid projector.
type ProjectionEngine struct {
	workerCount      int
	densityThreshold float64
	bufPool          *sync.Pool
}

func NewProjectionEngine(workerCount int, densityThreshold float64) *ProjectionEngine {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	if densityThreshold <= 0 {
		densityThreshold = 0.2
	}
	return &ProjectionEngine{
		workerCount:      workerCount,
		densityThreshold: densityThreshold,
		bufPool: &sync.Pool{
			New: func() any {
				// Allocate 2MB buffers to minimize reallocations for large chunks
				b := make([]byte, 0, 2*1024*1024)
				return &b
			},
		},
	}
}

// Stream performs zero-materialization hierarchical parallel streaming.
func (p *ProjectionEngine) Stream(
	ctx context.Context,
	plan *planner.QueryPlan,
	atoms map[domain.AtomKey]float64,
	meta *FastMetadata,
	writer io.Writer,
) error {
	// 1. Pre-calculate grid dimensions and resolve axes once
	rowCombos, rowIDCombos, _ := planner.MaterializeAxisCombos(plan.RowAxes, 64, true)
	_, colIDCombos, _ := planner.MaterializeAxisCombos(plan.ColAxes, 64, true)
	
	totalRows := len(rowCombos)
	if totalRows == 0 {
		return nil
	}

	// 2. Hierarchical Dispatch: Parallelize at Row-Range Level
	rowsPerWorker := (totalRows + p.workerCount - 1) / p.workerCount
	if rowsPerWorker < 1 { rowsPerWorker = 1 }

	type chunkResult struct {
		index int
		data  *[]byte
		err   error
	}

	resultChan := make(chan chunkResult, p.workerCount)
	var wg sync.WaitGroup

	for i := 0; i < p.workerCount; i++ {
		startRow := i * rowsPerWorker
		if startRow >= totalRows {
			break
		}
		endRow := startRow + rowsPerWorker
		if endRow > totalRows {
			endRow = totalRows
		}

		wg.Add(1)
		go func(idx int, s, e int) {
			defer wg.Done()
			
			bufPtr := p.bufPool.Get().(*[]byte)
			*bufPtr = (*bufPtr)[:0] // Reset but keep capacity

			err := p.streamRowRange(ctx, plan, atoms, meta, s, e, rowIDCombos, colIDCombos, bufPtr)
			resultChan <- chunkResult{index: idx, data: bufPtr, err: err}
		}(i, startRow, endRow)
	}

	// 3. Serialized Emitter: Collect and emit in strict order
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	pending := make(map[int]chunkResult)
	nextIndex := 0

	for res := range resultChan {
		pending[res.index] = res
		for {
			pNext, ok := pending[nextIndex]
			if !ok {
				break
			}
			if pNext.err != nil {
				return pNext.err
			}
			
			if len(*pNext.data) > 0 {
				if _, err := writer.Write(*pNext.data); err != nil {
					return err
				}
			}
			
			// Return buffer to pool
			p.bufPool.Put(pNext.data)
			
			delete(pending, nextIndex)
			nextIndex++
		}
	}

	return nil
}

// streamRowRange is the hot loop for workers. Zero allocations occur in the inner column loop.
func (p *ProjectionEngine) streamRowRange(
	ctx context.Context,
	plan *planner.QueryPlan,
	atoms map[domain.AtomKey]float64,
	meta *FastMetadata,
	rowStart, rowEnd int,
	rowIDCombos [][]int64,
	colIDCombos [][]int64,
	buf *[]byte,
) error {
	primaryMeasure := plan.BlockRef.MeasureIDs[0]
	primaryScenario := plan.BlockRef.ScenarioIDs[0]

	first := true
	for r := rowStart; r < rowEnd; r++ {
		rIDs := rowIDCombos[r]
		for c, cIDs := range colIDCombos {
			// Construct AtomKey - Zero Allocations (Value Type)
			key := domain.AtomKey{
				DimCount:   len(rIDs) + len(cIDs),
				MeasureID:  primaryMeasure,
				ScenarioID: primaryScenario,
			}
			copy(key.DimIDs[:], rIDs)
			copy(key.DimIDs[len(rIDs):], cIDs)
			key.EnsureCanonical()

			val, ok := atoms[key]
			if !ok {
				continue
			}

			if !first {
				*buf = append(*buf, ',')
			}
			first = false

			p.writeJSONCell(buf, r, c, val, meta, primaryMeasure, primaryScenario)
		}
	}
	return nil
}

// writeJSONCell performs ultra-fast byte splicing for JSON emission.
func (p *ProjectionEngine) writeJSONCell(
	buf *[]byte,
	r, c int,
	val float64,
	meta *FastMetadata,
	measureID, scenarioID int64,
) {
	*buf = append(*buf, `{"r":`...)
	*buf = strconv.AppendInt(*buf, int64(r), 10)
	*buf = append(*buf, `,"c":`...)
	*buf = strconv.AppendInt(*buf, int64(c), 10)
	*buf = append(*buf, `,"v":`...)
	*buf = strconv.AppendFloat(*buf, val, 'f', -1, 64)

	if meta != nil {
		if mJSON, ok := meta.MeasureJSON[measureID]; ok {
			*buf = append(*buf, `,"m":`...)
			*buf = append(*buf, mJSON...)
		}
		if sJSON, ok := meta.ScenarioJSON[scenarioID]; ok {
			*buf = append(*buf, `,"s":`...)
			*buf = append(*buf, sJSON...)
		}
	}
	*buf = append(*buf, '}')
}
