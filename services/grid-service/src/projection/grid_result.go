// File: services/grid-service/src/projection/grid_result.go
package projection

type GridResult struct {
	PlanID       string
	ViewID       string
	WindowHash   string
	AtomRevision int64

	Cells         []Cell
	RowOffsets    []int32
	RowMasks      []uint64
	StringArena   []byte
	StringOffsets []uint32
	ColumnStats   []ColumnStats

	OffHeapCells       *OffHeapArena
	OffHeapStringArena *OffHeapArena
}

func (gr *GridResult) Reset() {
	gr.PlanID = ""
	gr.ViewID = ""
	gr.WindowHash = ""
	gr.AtomRevision = 0

	gr.Cells = gr.Cells[:0]
	gr.RowOffsets = gr.RowOffsets[:0]
	gr.RowMasks = gr.RowMasks[:0]
	gr.StringArena = gr.StringArena[:0]
	gr.StringOffsets = gr.StringOffsets[:0]
	gr.ColumnStats = gr.ColumnStats[:0]

	gr.OffHeapCells = nil
	gr.OffHeapStringArena = nil
}

func (gr *GridResult) HasOffHeap() bool {
	return gr.OffHeapCells != nil || gr.OffHeapStringArena != nil
}

func (gr *GridResult) CellsFloat64() []float64 {
	if gr.OffHeapCells != nil {
		return gr.OffHeapCells.Float64Slice()
	}
	return CellsToFloat64Slice(gr.Cells)
}

func (gr *GridResult) StringArenaBytes() []byte {
	if gr.OffHeapStringArena != nil {
		return gr.OffHeapStringArena.Bytes()
	}
	return gr.StringArena
}
