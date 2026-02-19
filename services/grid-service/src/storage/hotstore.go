package storage

import (
    "context"
    "fmt"
    "time"

    "quantatomai/grid-service/domain"
)

type HotStoreClient interface {
	WriteAtom(ctx context.Context, atom domain.AtomWrite) error
}

// CircuitBreaker defines a protection interface for external dependencies.
type CircuitBreaker interface {
	// Execute wraps a call in a circuit breaker.
	Execute(
		ctx context.Context,
		fn func(context.Context) (map[domain.AtomKey]float64, error),
		cachedFn func(context.Context) (map[domain.AtomKey]float64, bool, error),
		staleFn func(context.Context) (map[domain.AtomKey]float64, bool, error),
	) (map[domain.AtomKey]float64, error)
}

// RedisHotStoreClient is a concrete implementation backed by Redis (or Scylla via a Redis-compatible API).
type RedisHotStoreClient struct {
    client RedisLikeClient
    ttl    time.Duration
}

// RedisLikeClient abstracts the Redis driver so you can swap implementations.
type RedisLikeClient interface {
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

// NewRedisHotStoreClient constructs a new hot store client.
func NewRedisHotStoreClient(client RedisLikeClient, ttl time.Duration) *RedisHotStoreClient {
    return &RedisHotStoreClient{
        client: client,
        ttl:    ttl,
    }
}

// redisKeyFromAtomKey builds a stable, canonical Redis key from an AtomKey.
func redisKeyFromAtomKey(key domain.AtomKey) string {
    // Use the canonical hash as the primary key component.
    hash := key.HashKey()
    // You can evolve this format later without changing the hash semantics.
    return fmt.Sprintf("atom:%d", hash)
}

// WriteAtom writes a single atom into the hot store with a stable, canonical key.
func (c *RedisHotStoreClient) WriteAtom(ctx context.Context, atom domain.AtomWrite) error {
    // Ensure canonical ordering before hashing.
    atom.Key.EnsureCanonical()

    redisKey := redisKeyFromAtomKey(atom.Key)

    // For now we store just the numeric value; you can later evolve this to a struct (JSON/msgpack).
    if err := c.client.Set(ctx, redisKey, atom.Value, c.ttl); err != nil {
        return fmt.Errorf("hotstore write failed for key %s: %w", redisKey, err)
    }

    return nil
}
