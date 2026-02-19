// File: services/grid-service/src/storage/grid_cache.go
package storage

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/google/flatbuffers/go"

	"quantatomai/grid-service/src/gen/fbgrid"
	"quantatomai/grid-service/src/projection"
)

const (
	WireVersionV2 = 2
)

type GridCacheKey struct {
	PlanID       string
	ViewID       string
	WindowHash   string
	AtomRevision int64
}

func (k GridCacheKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%d", k.PlanID, k.ViewID, k.WindowHash, k.AtomRevision)
}

type RedisClient interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type GridCache interface {
	Read(ctx context.Context, key GridCacheKey, target *projection.GridResult) error
	Write(ctx context.Context, key GridCacheKey, gr *projection.GridResult) error
}

type RedisGridCache struct {
	client  RedisClient
	builder *flatbuffers.Builder
}

func NewRedisGridCache(client RedisClient) *RedisGridCache {
	return &RedisGridCache{
		client:  client,
		builder: flatbuffers.NewBuilder(1024),
	}
}

func (c *RedisGridCache) Write(
	ctx context.Context,
	key GridCacheKey,
	gr *projection.GridResult,
) error {

	payload := c.encodeGridPayload(gr)
	checksum := crc32.ChecksumIEEE(payload)

	b := c.builder
	b.Reset()

	payloadOff := b.CreateByteVector(payload)

	fbgrid.GridWireEnvelopeStart(b)
	fbgrid.GridWireEnvelopeAddWireVersion(b, WireVersionV2)
	fbgrid.GridWireEnvelopeAddCompression(b, fbgrid.CompressionNONE)
	fbgrid.GridWireEnvelopeAddContentLength(b, uint32(len(payload)))
	fbgrid.GridWireEnvelopeAddContentChecksum(b, checksum)
	fbgrid.GridWireEnvelopeAddPayload(b, payloadOff)
	root := fbgrid.GridWireEnvelopeEnd(b)
	b.Finish(root)

	return c.client.Set(ctx, key.String(), b.FinishedBytes())
}

func (c *RedisGridCache) Read(
	ctx context.Context,
	key GridCacheKey,
	target *projection.GridResult,
) error {

	data, err := c.client.Get(ctx, key.String())
	if err != nil {
		return err
	}

	buf := flatbuffers.NewByteBuffer(data)
	env := fbgrid.GetRootAsGridWireEnvelope(buf.Bytes, buf.Pos)

	if env.WireVersion() != WireVersionV2 {
		return ErrIncompatibleWireVersion
	}

	payload := env.PayloadBytes()
	if len(payload) == 0 {
		return ErrEmptyPayload
	}

	if crc32.ChecksumIEEE(payload) != env.ContentChecksum() {
		return ErrChecksumMismatch
	}

	DecodeGrid(payload, target)
	return nil
}

func (c *RedisGridCache) encodeGridPayload(gr *projection.GridResult) []byte {
	b := c.builder
	b.Reset()

	planIDOff := b.CreateString(gr.PlanID)
	viewIDOff := b.CreateString(gr.ViewID)
	windowHashOff := b.CreateString(gr.WindowHash)

	fbgrid.GridStartRowOffsetsVector(b, len(gr.RowOffsets))
	for i := len(gr.RowOffsets) - 1; i >= 0; i-- {
		b.PrependInt32(gr.RowOffsets[i])
	}
	rowOffsetsOff := b.EndVector(len(gr.RowOffsets))

	fbgrid.GridStartRowMasksVector(b, len(gr.RowMasks))
	for i := len(gr.RowMasks) - 1; i >= 0; i-- {
		b.PrependUint64(gr.RowMasks[i])
	}
	rowMasksOff := b.EndVector(len(gr.RowMasks))

	var cellsOff flatbuffers.UOffsetT
	if gr.OffHeapCells != nil {
		raw := gr.OffHeapCells.Bytes()
		cellsOff = b.CreateByteVector(raw)
	} else {
		floatSlice := projection.CellsToFloat64Slice(gr.Cells)
		byteSlice := projection.Float64Bytes(floatSlice)
		cellsOff = b.CreateByteVector(byteSlice)
	}

	var stringArenaOff flatbuffers.UOffsetT
	if gr.OffHeapStringArena != nil {
		raw := gr.OffHeapStringArena.Bytes()
		stringArenaOff = b.CreateByteVector(raw)
	} else if len(gr.StringArena) > 0 {
		stringArenaOff = b.CreateByteVector(gr.StringArena)
	}

	fbgrid.GridStartStringOffsetsVector(b, len(gr.StringOffsets))
	for i := len(gr.StringOffsets) - 1; i >= 0; i-- {
		b.PrependUint32(gr.StringOffsets[i])
	}
	stringOffsetsOff := b.EndVector(len(gr.StringOffsets))

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

	return b.FinishedBytes()
}

var (
	ErrIncompatibleWireVersion = fmt.Errorf("incompatible wire version")
	ErrChecksumMismatch        = fmt.Errorf("payload checksum mismatch")
	ErrEmptyPayload            = fmt.Errorf("empty payload in envelope")
)
