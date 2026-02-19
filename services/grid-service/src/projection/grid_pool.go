// File: services/grid-service/src/projection/grid_pool.go
package projection

import "sync"

var gridResultPool = sync.Pool{
	New: func() any {
		return &GridResult{
			Cells:         make([]Cell, 0, 1024),
			RowOffsets:    make([]int32, 0, 256),
			RowMasks:      make([]uint64, 0, 256),
			ColumnStats:   make([]ColumnStats, 0, 64),
			StringArena:   make([]byte, 0, 4096),
			StringOffsets: make([]uint32, 0, 256),
		}
	},
}

func GetGridResult() *GridResult {
	gr := gridResultPool.Get().(*GridResult)
	gr.Reset()
	return gr
}

func PutGridResult(gr *GridResult) {
	if gr == nil {
		return
	}

	gr.OffHeapCells = nil
	gr.OffHeapStringArena = nil

	if cap(gr.Cells) > 1_000_000 {
		gr.Cells = nil
	}
	if cap(gr.StringArena) > 8_000_000 {
		gr.StringArena = nil
	}

	gridResultPool.Put(gr)
}
