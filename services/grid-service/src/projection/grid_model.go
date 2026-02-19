package projection

import (
	"math"
)

const (
	valueTagMask   uint64 = 0x7FF8000000000000
	valuePayloadMask      = 0x0007FFFFFFFFFFFF
	nullPayload    uint64 = valuePayloadMask
)

func (g *GridResult) HasData(row, col int) bool {
	if row < 0 || row >= len(g.RowMasks) {
		return false
	}
	if col < 0 || col >= 64 { // fast-path only
		return false
	}
	return (g.RowMasks[row] & (1 << uint(col))) != 0
}

func EncodeNumeric(v float64) Cell {
	return Cell{Value: v}
}

func EncodeNull() Cell {
	payload := valueTagMask | nullPayload
	return Cell{Value: math.Float64frombits(payload)}
}

func EncodeStringIndex(idx uint32) Cell {
	payload := uint64(idx) & valuePayloadMask
	bits := valueTagMask | payload
	return Cell{Value: math.Float64frombits(bits)}
}

func (c Cell) IsNull() bool {
	bits := math.Float64bits(c.Value)
	return (bits&valueTagMask) == valueTagMask && (bits&valuePayloadMask) == nullPayload
}

func (c Cell) IsString() bool {
	bits := math.Float64bits(c.Value)
	return (bits&valueTagMask) == valueTagMask && (bits&valuePayloadMask) != nullPayload
}

func (c Cell) StringIndex() uint32 {
	bits := math.Float64bits(c.Value)
	return uint32(bits & valuePayloadMask)
}

func (c Cell) AsNumeric() float64 {
	return c.Value
}
