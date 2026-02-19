package storage

import (
	"fmt"

	"github.com/google/flatbuffers/go"
	"quantatomai/grid-service/src/projection"
)

type FlatBufferWireFormatCodec struct{}

func NewFlatBufferWireFormatCodec() *FlatBufferWireFormatCodec {
	return &FlatBufferWireFormatCodec{}
}

// EncodeGrid expects *projection.GridResult and uses the low-level Builder API
// to achieve theoretical peak serialization performance.
func (c *FlatBufferWireFormatCodec) EncodeGrid(grid interface{}) ([]byte, error) {
	gr, ok := grid.(*projection.GridResult)
	if !ok {
		return nil, fmt.Errorf("FlatBufferWireFormatCodec: expected *projection.GridResult, got %T", grid)
	}
	// Call the optimized low-level builder
	return projection.BuildFlatBufferFromGridResult(gr), nil
}

// DecodeGrid unpacks the FlatBuffer data back into a projection.GridResult.
// It supports object reuse by utilizing the capacity of the target slices.
func (c *FlatBufferWireFormatCodec) DecodeGrid(data []byte, target interface{}) error {
	gr, ok := target.(*projection.GridResult)
	if !ok {
		return fmt.Errorf("FlatBufferWireFormatCodec: expected *projection.GridResult as target, got %T", target)
	}

	DecodeGrid(data, gr)
	return nil
}

func (c *FlatBufferWireFormatCodec) Format() WireFormatType {
	return WireFormatFlatBuffer
}
