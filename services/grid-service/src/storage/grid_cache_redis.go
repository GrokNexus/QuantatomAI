package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

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

// GridCacheEntry is what we store in Redis.
type GridCacheEntry struct {
	ETag          string          `json:"etag"`
	CreatedAtUnix int64           `json:"created_at_unix"`
	Stale         bool            `json:"stale"`
	Payload       json.RawMessage `json:"payload"`
}

// GridCache defines the contract used by the handler + breaker.
type GridCache interface {
	// Get returns (entry, found, err).
	Get(ctx context.Context, key GridCacheKey) (*GridCacheEntry, bool, error)
	// Set stores an entry with TTL.
	Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error
	// InvalidateByAtomRevision invalidates all entries for a given plan + atom revision.
	InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error
}

// RedisGridCache is a Redis-backed implementation of GridCache.
type RedisGridCache struct {
	client *redis.Client

	// key prefixes
	dataPrefix  string
	indexPrefix string
	defaultTTL  time.Duration
}

// NewRedisGridCache constructs a RedisGridCache.
func NewRedisGridCache(client *redis.Client, dataPrefix, indexPrefix string, defaultTTL time.Duration) *RedisGridCache {
	if dataPrefix == "" {
		dataPrefix = "grid-cache:data:"
	}
	if indexPrefix == "" {
		indexPrefix = "grid-cache:index:"
	}
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}
	return &RedisGridCache{
		client:      client,
		dataPrefix:  dataPrefix,
		indexPrefix: indexPrefix,
		defaultTTL:  defaultTTL,
	}
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

	var entry GridCacheEntry
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return nil, false, err
	}
	return &entry, true, nil
}

// Set implements GridCache.Set.
func (c *RedisGridCache) Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	dataKey := c.dataKey(key)
	indexKey := c.indexKey(key.PlanID, key.AtomRevision)

	b, err := json.Marshal(entry)
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

// InvalidateByAtomRevision implements GridCache.InvalidateByAtomRevision.
func (c *RedisGridCache) InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error {
	indexKey := c.indexKey(planID, atomRevision)

	keys, err := c.client.SMembers(ctx, indexKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(keys) == 0 {
		_ = c.client.Del(ctx, indexKey).Err()
		return nil
	}

	pipe := c.client.TxPipeline()
	for _, k := range keys {
		pipe.Del(ctx, k)
	}
	pipe.Del(ctx, indexKey)
	_, err = pipe.Exec(ctx)
	return err
}

// Helper: build Redis keys.

func (c *RedisGridCache) dataKey(key GridCacheKey) string {
	raw, _ := json.Marshal(key)
	sum := sha256.Sum256(raw)
	return c.dataPrefix + hex.EncodeToString(sum[:])
}

func (c *RedisGridCache) indexKey(planID string, atomRevision int64) string {
	return c.indexPrefix + planID + ":" + int64ToString(atomRevision)
}

func int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}
