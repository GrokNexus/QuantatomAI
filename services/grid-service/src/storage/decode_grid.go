// File: services/grid-service/src/storage/decode_grid.go
package storage

import (
	"unsafe"

	"github.com/google/flatbuffers/go"

	"quantatomai/grid-service/src/gen/fbgrid"
	"quantatomai/grid-service/src/projection"
)

const offHeapDecodeThreshold = 1_000_000

var GlobalArenaManager *projection.ArenaManager

func DecodeGrid(data []byte, target *projection.GridResult) {
	buf := flatbuffers.NewByteBuffer(data)
	fbGrid := fbgrid.GetRootAsGrid(buf.Bytes, buf.Pos)

	target.PlanID = string(fbGrid.PlanId())
	target.ViewID = string(fbGrid.ViewId())
	target.WindowHash = string(fbGrid.WindowHash())
	target.AtomRevision = fbGrid.AtomRevision()

	rowOffsetsLen := fbGrid.RowOffsetsLength()
	if cap(target.RowOffsets) < rowOffsetsLen {
		target.RowOffsets = make([]int32, rowOffsetsLen)
	} else {
		target.RowOffsets = target.RowOffsets[:rowOffsetsLen]
	}
	for i := 0; i < rowOffsetsLen; i++ {
		target.RowOffsets[i] = fbGrid.RowOffsets(i)
	}

	rowMasksLen := fbGrid.RowMasksLength()
	if cap(target.RowMasks) < rowMasksLen {
		target.RowMasks = make([]uint64, rowMasksLen)
	} else {
		target.RowMasks = target.RowMasks[:rowMasksLen]
	}
	for i := 0; i < rowMasksLen; i++ {
		target.RowMasks[i] = fbGrid.RowMasks(i)
	}

	cellsLen := fbGrid.CellsLength()
	raw := fbGrid.CellsBytes()
	byteLen := len(raw)

	useOffHeap := GlobalArenaManager != nil && byteLen >= offHeapDecodeThreshold

	if useOffHeap {
		arena, err := GlobalArenaManager.Acquire(byteLen)
		if err == nil {
			arena.WriteBytes(raw)
			target.OffHeapCells = arena
			target.Cells = target.Cells[:0]
		} else {
			useOffHeap = false
		}
	}

	if !useOffHeap {
		if cap(target.Cells) < cellsLen {
			target.Cells = make([]projection.Cell, cellsLen)
		} else {
			target.Cells = target.Cells[:cellsLen]
		}

		if byteLen > 0 {
			floatSlice := unsafe.Slice((*float64)(unsafe.Pointer(&raw[0])), cellsLen)
			for i := 0; i < cellsLen; i++ {
				target.Cells[i].Value = floatSlice[i]
			}
		}
	}

	stringArenaLen := fbGrid.StringArenaLength()
	if stringArenaLen > 0 {
		useOffHeapSA := GlobalArenaManager != nil && stringArenaLen >= offHeapDecodeThreshold

		if useOffHeapSA {
			arena, err := GlobalArenaManager.Acquire(stringArenaLen)
			if err == nil {
				tmp := make([]byte, stringArenaLen)
				for i := 0; i < stringArenaLen; i++ {
					tmp[i] = fbGrid.StringArena(i)
				}
				arena.WriteBytes(tmp)
				target.OffHeapStringArena = arena
				target.StringArena = target.StringArena[:0]
			} else {
				useOffHeapSA = false
			}
		}

		if !useOffHeapSA {
			if cap(target.StringArena) < stringArenaLen {
				target.StringArena = make([]byte, stringArenaLen)
			} else {
				target.StringArena = target.StringArena[:stringArenaLen]
			}
			for i := 0; i < stringArenaLen; i++ {
				target.StringArena[i] = fbGrid.StringArena(i)
			}
		}
	} else {
		target.StringArena = target.StringArena[:0]
	}

	stringOffsetsLen := fbGrid.StringOffsetsLength()
	if cap(target.StringOffsets) < stringOffsetsLen {
		target.StringOffsets = make([]uint32, stringOffsetsLen)
	} else {
		target.StringOffsets = target.StringOffsets[:stringOffsetsLen]
	}
	for i := 0; i < stringOffsetsLen; i++ {
		target.StringOffsets[i] = fbGrid.StringOffsets(i)
	}

	target.ColumnStats = target.ColumnStats[:0]
}
