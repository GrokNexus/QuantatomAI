package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewGridSubsystem is the single entry point for initializing the Ultra-Diamond performance layer.
// It wires the L1/L2 tiered cache stack and the adaptive cache writer.
func NewGridSubsystem(
	ctx context.Context,
	redisClient *redis.Client,
	nodeID string,
	cfg *CacheConfig,
	prefetchRecomputeFn func(ctx context.Context, planID, viewID, windowHash string, version int64),
	coordinator PrefetchCoordinator,
	xfetchRefresher XFetchRefresher,
) (*TieredGridCache, *GridCacheWriter, error) {
	if cfg == nil {
		cfg = NewDefaultCacheConfig()
	}
	codec := cfg.BuildWireCodec()

	// 1. Initialize L2 Cache (Redis)
	l2 := NewRedisGridCache(redisClient, codec)

	// 2. Initialize Analytics & Prefetching
	stats := newHotKeyStats()
	prefetcher := NewDefaultGoldenWindowPrefetcher(stats, 5, prefetchRecomputeFn)

	// 3. Initialize Tiered Engine (L1 + L2 + Pub/Sub)
	tiered, err := NewTieredGridCache(
		l2,
		redisClient,
		512*1024*1024,          // maxMemory
		"grid-cache:invalidate:", // pubsub prefix
		nodeID,
		prefetcher,
		coordinator,
		xfetchRefresher,
		cfg.XFetchProbability,
	)
	if err != nil {
		return nil, nil, err
	}

	// 4. Activate Distributed Invalidation
	if err := tiered.StartInvalidationSubscriber(ctx); err != nil {
		return nil, nil, err
	}

	// 5. Initialize Adaptive Hydrators
	var simpleHydrator ChunkedHydrator
	var streamingHydrator ChunkedHydrator
	if cfg.EnableChunkedHydration {
		simpleHydrator = NewSimpleChunkedHydrator(l2)
		streamingHydrator = NewStreamingRedisHydrator(redisClient)
	}

	// 6. Assemble Adaptive Writer
	writer := NewGridCacheWriter(
		tiered,
		codec,
		cfg,
		simpleHydrator,
		streamingHydrator,
	)

	return tiered, writer, nil
}
