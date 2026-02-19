package storage

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// ChunkedPayload represents a large grid payload broken into chunks.
type ChunkedPayload struct {
    PlanID       string
    ViewID       string
    AtomRevision int64
    WindowHash   string
    TotalSize    int64
    Chunks       [][]byte
}

// ChunkedHydrator defines the contract for streaming large grid payloads.
type ChunkedHydrator interface {
    Hydrate(ctx context.Context, key GridCacheKey, payload *ChunkedPayload, ttl time.Duration) error
}

// -----------------------------------------------------------------------------
// SimpleChunkedHydrator (baseline reference)
// -----------------------------------------------------------------------------

type SimpleChunkedHydrator struct {
    l2 GridCache
}

func NewSimpleChunkedHydrator(l2 GridCache) *SimpleChunkedHydrator {
    return &SimpleChunkedHydrator{l2: l2}
}

func (h *SimpleChunkedHydrator) Hydrate(ctx context.Context, key GridCacheKey, payload *ChunkedPayload, ttl time.Duration) error {
    if payload == nil {
        return nil
    }

    // Concatenate chunks (baseline behavior)
    buf := make([]byte, 0, payload.TotalSize)
    for _, c := range payload.Chunks {
        buf = append(buf, c...)
    }

    entry := &GridCacheEntry{
        PlanID:       key.PlanID,
        ViewID:       key.ViewID,
        AtomRevision: key.AtomRevision,
        WindowHash:   key.WindowHash,
        Payload:      buf,
    }
    return h.l2.Set(ctx, key, entry, ttl)
}

// -----------------------------------------------------------------------------
// StreamingRedisHydrator (Titanium++ version)
// -----------------------------------------------------------------------------

// StreamingRedisHydrator streams chunks to Redis using pipelining.
// It writes:
//   grid:{plan}:{view}:{rev}:{hash}:chunk:0
//   grid:{plan}:{view}:{rev}:{hash}:chunk:1
//   ...
// And a manifest:
//   grid:{plan}:{view}:{rev}:{hash}:manifest
//
// The Tiered Cache can read the manifest and reassemble if needed.
type StreamingRedisHydrator struct {
    redis *redis.Client
}

func NewStreamingRedisHydrator(redis *redis.Client) *StreamingRedisHydrator {
    return &StreamingRedisHydrator{redis: redis}
}

func (h *StreamingRedisHydrator) Hydrate(ctx context.Context, key GridCacheKey, payload *ChunkedPayload, ttl time.Duration) error {
    if payload == nil {
        return nil
    }

    pipe := h.redis.Pipeline()

    baseKey := fmt.Sprintf(
        "grid:{%s}:%s:%d:%s",
        key.PlanID, key.ViewID, key.AtomRevision, key.WindowHash,
    )

    // Write each chunk independently
    for i, chunk := range payload.Chunks {
        chunkKey := fmt.Sprintf("%s:chunk:%d", baseKey, i)
        pipe.Set(ctx, chunkKey, chunk, ttl)
    }

    // Write manifest
    manifestKey := fmt.Sprintf("%s:manifest", baseKey)
    manifest := map[string]interface{}{
        "plan_id":       key.PlanID,
        "view_id":       key.ViewID,
        "atom_revision": key.AtomRevision,
        "window_hash":   key.WindowHash,
        "total_size":    payload.TotalSize,
        "chunk_count":   len(payload.Chunks),
    }
    pipe.HSet(ctx, manifestKey, manifest)
    pipe.Expire(ctx, manifestKey, ttl)

    _, err := pipe.Exec(ctx)
    return err
}
