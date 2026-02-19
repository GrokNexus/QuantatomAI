package storage

import (
    "bytes"
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "strconv"
    "strings"
    "time"

    "github.com/klauspost/compress/zstd"
    "github.com/redis/go-redis/v9"
)

// GridCacheKey captures the logical identity of a grid response.
type GridCacheKey struct {
    PlanID        string `json:"plan_id"`
    ViewID        string `json:"view_id"`
    UserID        string `json:"user_id"`
    ScenarioID    string `json:"scenario_id"`
    AtomRevision  int64  `json:"atom_revision"`
    WindowHash    string `json:"window_hash"`
    AdditionalTag string `json:"additional_tag,omitempty"`
}

// GridCacheEntry is what we store logically (before compression).
type GridCacheEntry struct {
    ETag          string `json:"etag"`
    CreatedAtUnix int64  `json:"created_at_unix"`
    Stale         bool   `json:"stale"`
    Payload       []byte `json:"payload"`
}

// internal wire format stored in Redis (compressed payload)
type gridCacheWire struct {
    ETag          string `json:"etag"`
    CreatedAtUnix int64  `json:"created_at_unix"`
    Stale         bool   `json:"stale"`
    // compressed payload bytes (zstd)
    PayloadCompressed []byte `json:"payload_compressed"`
}

// GridCache defines the contract used by the handler + breaker.
type GridCache interface {
    Get(ctx context.Context, key GridCacheKey) (*GridCacheEntry, bool, error)
    Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error
    InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error
}

// RedisGridCache is a Redis-backed implementation of GridCache with:
// - deterministic keys
// - SSCAN-based invalidation
// - zstd compression for payloads
type RedisGridCache struct {
    client *redis.Client

    dataPrefix  string
    indexPrefix string
    defaultTTL  time.Duration

    zstdEnc *zstd.Encoder
    zstdDec *zstd.Decoder
}

func NewRedisGridCache(client *redis.Client, dataPrefix, indexPrefix string, defaultTTL time.Duration) (*RedisGridCache, error) {
    if dataPrefix == "" {
        dataPrefix = "grid-cache:data:"
    }
    if indexPrefix == "" {
        indexPrefix = "grid-cache:index:"
    }
    if defaultTTL <= 0 {
        defaultTTL = 5 * time.Minute
    }

    enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
    if err != nil {
        return nil, err
    }
    dec, err := zstd.NewReader(nil)
    if err != nil {
        return nil, err
    }

    return &RedisGridCache{
        client:      client,
        dataPrefix:  dataPrefix,
        indexPrefix: indexPrefix,
        defaultTTL:  defaultTTL,
        zstdEnc:     enc,
        zstdDec:     dec,
    }, nil
}

// Get implements GridCache.Get.
func (c *RedisGridCache) Get(ctx context.Context, key GridCacheKey) (*GridCacheEntry, bool, error) {
    dataKey := c.dataKey(key)
    val, err := c.client.Get(ctx, dataKey).Result()
    if err == redis.Nil {
        return nil, false, nil
    }
    if err != nil {
        return nil, false, err
    }

    var wire gridCacheWire
    if err := json.Unmarshal([]byte(val), &wire); err != nil {
        return nil, false, err
    }

    // decompress payload
    payload, err := c.zstdDec.DecodeAll(wire.PayloadCompressed, nil)
    if err != nil {
        return nil, false, err
    }

    entry := &GridCacheEntry{
        ETag:          wire.ETag,
        CreatedAtUnix: wire.CreatedAtUnix,
        Stale:         wire.Stale,
        Payload:       json.RawMessage(payload),
    }
    return entry, true, nil
}

// Set implements GridCache.Set.
func (c *RedisGridCache) Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error {
    if ttl <= 0 {
        ttl = c.defaultTTL
    }

    dataKey := c.dataKey(key)
    indexKey := c.indexKey(key.PlanID, key.AtomRevision)

    // compress payload
    compressed, err := c.zstdEnc.EncodeAll([]byte(entry.Payload), nil)
    if err != nil {
        return err
    }

    wire := gridCacheWire{
        ETag:             entry.ETag,
        CreatedAtUnix:    entry.CreatedAtUnix,
        Stale:            entry.Stale,
        PayloadCompressed: compressed,
    }
    b, err := json.Marshal(wire)
    if err != nil {
        return err
    }

    pipe := c.client.TxPipeline()
    pipe.Set(ctx, dataKey, b, ttl)
    pipe.SAdd(ctx, indexKey, dataKey)
    pipe.Expire(ctx, indexKey, ttl*2)
    _, err = pipe.Exec(ctx)
    return err
}

// InvalidateByAtomRevision implements GridCache.InvalidateByAtomRevision using SSCAN.
func (c *RedisGridCache) InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error {
    indexKey := c.indexKey(planID, atomRevision)

    var cursor uint64
    for {
        keys, nextCursor, err := c.client.SScan(ctx, indexKey, cursor, "*", 500).Result()
        if err != nil && err != redis.Nil {
            return err
        }
        if len(keys) > 0 {
            pipe := c.client.TxPipeline()
            for _, k := range keys {
                pipe.Del(ctx, k)
                pipe.SRem(ctx, indexKey, k)
            }
            if _, err := pipe.Exec(ctx); err != nil {
                return err
            }
        }
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }

    _ = c.client.Del(ctx, indexKey).Err()
    return nil
}

// Helper: deterministic Redis keys.

func (c *RedisGridCache) dataKey(key GridCacheKey) string {
    // Deterministic, stable key composition.
    // We still hash the final string to keep keys compact.
    var sb strings.Builder
    sb.Grow(128)
    sb.WriteString(key.PlanID)
    sb.WriteByte(':')
    sb.WriteString(key.ViewID)
    sb.WriteByte(':')
    sb.WriteString(key.UserID)
    sb.WriteByte(':')
    sb.WriteString(key.ScenarioID)
    sb.WriteByte(':')
    sb.WriteString(strconv.FormatInt(key.AtomRevision, 10))
    sb.WriteByte(':')
    sb.WriteString(key.WindowHash)
    if key.AdditionalTag != "" {
        sb.WriteByte(':')
        sb.WriteString(key.AdditionalTag)
    }

    raw := sb.String()
    sum := sha256.Sum256([]byte(raw))
    return c.dataPrefix + hex.EncodeToString(sum[:])
}

func (c *RedisGridCache) indexKey(planID string, atomRevision int64) string {
    var buf bytes.Buffer
    buf.Grow(len(c.indexPrefix) + len(planID) + 24)
    buf.WriteString(c.indexPrefix)
    buf.WriteString(planID)
    buf.WriteByte(':')
    buf.WriteString(strconv.FormatInt(atomRevision, 10))
    return buf.String()
}
