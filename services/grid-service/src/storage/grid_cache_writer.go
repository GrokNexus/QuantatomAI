package storage

import (
    "bytes"
    "context"
    "time"
)

type GridProjector interface {
    Project(ctx context.Context, planID, viewID, windowHash string) (interface{}, error)
}

type TieredGridCacheWriter struct {
	cache     *TieredGridCache
	wireCodec WireFormatCodec
	cfg       *CacheConfig
	simple    ChunkedHydrator
	streaming ChunkedHydrator
}

func NewTieredGridCacheWriter(
	cache *TieredGridCache,
	wireCodec WireFormatCodec,
	cfg *CacheConfig,
	simpleHydrator ChunkedHydrator,
	streamingHydrator ChunkedHydrator,
) *TieredGridCacheWriter {
    if wireCodec == nil {
        wireCodec = NewJSONWireFormatCodec()
    }
    if cfg == nil {
        cfg = NewDefaultCacheConfig()
    }
	return &TieredGridCacheWriter{
		cache:     cache,
		wireCodec: wireCodec,
		cfg:       cfg,
		simple:    simpleHydrator,
		streaming: streamingHydrator,
	}
}

func (w *TieredGridCacheWriter) WriteProjectedGrid(
    ctx context.Context,
    key GridCacheKey,
    projected interface{},
) error {
    data, err := w.wireCodec.EncodeGrid(projected)
    if err != nil {
        return err
    }

    ttl := w.cfg.DefaultTTL
    size := len(data)

    // No chunking → direct write
    if !w.cfg.EnableChunkedHydration {
        entry := &GridCacheEntry{
            PlanID:       key.PlanID,
            ViewID:       key.ViewID,
            AtomRevision: key.AtomRevision,
            WindowHash:   key.WindowHash,
            Payload:      data,
        }
        return w.cache.Set(ctx, key, entry, ttl)
    }

    switch {
    case size < 512*1024:
        entry := &GridCacheEntry{
            PlanID:       key.PlanID,
            ViewID:       key.ViewID,
            AtomRevision: key.AtomRevision,
            WindowHash:   key.WindowHash,
            Payload:      data,
        }
        return w.cache.Set(ctx, key, entry, ttl)

    case size < 5*1024*1024:
        cr := NewChunkedReader(256 * 1024)
        chunks, total, err := cr.ReadIntoChunks(bytes.NewReader(data))
        if err != nil {
            return err
        }
        payload := &ChunkedPayload{
            PlanID:       key.PlanID,
            ViewID:       key.ViewID,
            AtomRevision: key.AtomRevision,
            WindowHash:   key.WindowHash,
            TotalSize:    total,
            Chunks:       chunks,
        }
        return w.simple.Hydrate(ctx, key, payload, ttl)

    default:
        cr := NewChunkedReader(512 * 1024)
        chunks, total, err := cr.ReadIntoChunks(bytes.NewReader(data))
        if err != nil {
            return err
        }
        payload := &ChunkedPayload{
            PlanID:       key.PlanID,
            ViewID:       key.ViewID,
            AtomRevision: key.AtomRevision,
            WindowHash:   key.WindowHash,
            TotalSize:    total,
            Chunks:       chunks,
        }
        return w.streaming.Hydrate(ctx, key, payload, ttl)
    }
}
