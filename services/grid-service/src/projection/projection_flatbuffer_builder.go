package projection

import (
	"github.com/google/flatbuffers/go"
	"quantatomai/grid-service/src/gen/fbgrid"
)

// BuildFlatBufferFromGridResult streams GridResult into a FlatBuffer using the low-level Builder API.
// This is the absolute peak performance path, skipping the Object API's intermediate allocations.
func BuildFlatBufferFromGridResult(gr *GridResult) []byte {
	if gr == nil {
		return nil
	}

	b := flatbuffers.NewBuilder(1024)

	// --- cells vector (Peak: Unsafe Zero-Copy) ---
	floats := CellsToFloat64Slice(gr.Cells)
	floatBytes := Float64Bytes(floats)
	cellsOff := b.CreateByteVector(floatBytes)

	// --- row_offsets vector ---
	rowOffsetsCount := len(gr.RowOffsets)
	fbgrid.GridStartRowOffsetsVector(b, rowOffsetsCount)
	for i := rowOffsetsCount - 1; i >= 0; i-- {
		b.PrependInt32(gr.RowOffsets[i])
	}
	rowOffsetsOff := b.EndVector(rowOffsetsCount)

	// --- row_masks vector ---
	rowMasksCount := len(gr.RowMasks)
	fbgrid.GridStartRowMasksVector(b, rowMasksCount)
	for i := rowMasksCount - 1; i >= 0; i-- {
		b.PrependUint64(gr.RowMasks[i])
	}
	rowMasksOff := b.EndVector(rowMasksCount)

	// --- string_arena vector ---
	stringArenaOff := flatbuffers.UOffsetT(0)
	if len(gr.StringArena) > 0 {
		stringArenaOff = b.CreateByteVector(gr.StringArena)
	}

	// --- string_offsets vector ---
	stringOffsetsCount := len(gr.StringOffsets)
	fbgrid.GridStartStringOffsetsVector(b, stringOffsetsCount)
	for i := stringOffsetsCount - 1; i >= 0; i-- {
		b.PrependUint32(gr.StringOffsets[i])
	}
	stringOffsetsOff := b.EndVector(stringOffsetsCount)

	// --- root table ---
	planIDOff := b.CreateString(gr.PlanID)
	viewIDOff := b.CreateString(gr.ViewID)
	windowHashOff := b.CreateString(gr.WindowHash)

	fbgrid.GridStart(b)

	fbgrid.GridAddPlanId(b, planIDOff)
	fbgrid.GridAddViewId(b, viewIDOff)
	fbgrid.GridAddWindowHash(b, windowHashOff)
	fbgrid.GridAddAtomRevision(b, gr.AtomRevision)

	fbgrid.GridAddCells(b, cellsOff)
	fbgrid.GridAddRowOffsets(b, rowOffsetsOff)
	fbgrid.GridAddRowMasks(b, rowMasksOff)
	fbgrid.GridAddStringArena(b, stringArenaOff)
	fbgrid.GridAddStringOffsets(b, stringOffsetsOff)

	root := fbgrid.GridEnd(b)
	b.Finish(root)

	// Return a copy of the finished bytes to avoid builder state dependency
	out := make([]byte, len(b.FinishedBytes()))
	copy(out, b.FinishedBytes())
	return out
}
