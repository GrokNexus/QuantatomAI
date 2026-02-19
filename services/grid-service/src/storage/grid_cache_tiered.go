package storage

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
)

// TieredGridCache implements L1 (in-process) + L2 (Redis) caching with
// Pub/Sub-based invalidation and plan-scoped channels (cluster-slot friendly).
//
// L1:
//   - In-memory LRU cache (fastest path)
//   - Keyed by a deterministic logical key derived from GridCacheKey
//
// L2:
//   - Any GridCache implementation (typically RedisGridCache)
//
// Invalidation:
//   - InvalidateByAtomRevision:
//       * Calls L2.InvalidateByAtomRevision
//       * Publishes an invalidation event on Redis Pub/Sub
//   - All nodes subscribe to the invalidation channel and purge their L1
//     (coarse-grained but safe and fast).
type TieredGridCache struct {
	l1         *lru.Cache[string, *GridCacheEntry]
	l2         GridCache
	redis      *redis.Client
	pubsubPref string
	nodeID     string
	subscriber *redis.PubSub
	subscribed bool
}

// InvalidationEvent is broadcast over Pub/Sub.
type InvalidationEvent struct {
	PlanID       string `json:"plan_id"`
	AtomRevision int64  `json:"atom_revision"`
	SourceNodeID string `json:"source_node_id"`
}

// NewTieredGridCache constructs a TieredGridCache.
func NewTieredGridCache(
	l2 GridCache,
	redis *redis.Client,
	l1Size int,
	pubsubPref string,
	nodeID string,
) (*TieredGridCache, error) {
	if l1Size <= 0 {
		l1Size = 1024
	}
	if pubsubPref == "" {
		pubsubPref = "grid-cache:invalidate:"
	}
	if nodeID == "" {
		nodeID = "node-unknown"
	}

	l1, err := lru.New[string, *GridCacheEntry](l1Size)
	if err != nil {
		return nil, err
	}

	return &TieredGridCache{
		l1:         l1,
		l2:         l2,
		redis:      redis,
		pubsubPref: pubsubPref,
		nodeID:     nodeID,
	}, nil
}

// StartInvalidationSubscriber starts a background goroutine that listens
// for invalidation events and purges the L1 cache when they arrive.
func (t *TieredGridCache) StartInvalidationSubscriber(ctx context.Context) error {
	if t.redis == nil {
		return nil
	}

	// Pattern-based subscription so we can use plan-scoped channels.
	pattern := t.pubsubPref + "*"

	ps := t.redis.PSubscribe(ctx, pattern)
	t.subscriber = ps
	t.subscribed = true

	go func() {
		for {
			msg, err := ps.ReceiveMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("tiered-cache: pubsub receive error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var ev InvalidationEvent
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				log.Printf("tiered-cache: invalidation event unmarshal error: %v", err)
				continue
			}

			// Coarse but safe: purge entire L1 on any invalidation.
			t.l1.Purge()
		}
	}()

	return nil
}

// Get implements GridCache.Get with L1/L2 tiering.
func (t *TieredGridCache) Get(ctx context.Context, key GridCacheKey) (*GridCacheEntry, bool, error) {
	l1Key := t.l1Key(key)

	if entry, ok := t.l1.Get(l1Key); ok && entry != nil {
		return entry, true, nil
	}

	entry, found, err := t.l2.Get(ctx, key)
	if err != nil || !found || entry == nil {
		return nil, false, err
	}

	t.l1.Add(l1Key, entry)
	return entry, true, nil
}

// Set implements GridCache.Set with L1/L2 tiering.
func (t *TieredGridCache) Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error {
	if err := t.l2.Set(ctx, key, entry, ttl); err != nil {
		return err
	}
	l1Key := t.l1Key(key)
	t.l1.Add(l1Key, entry)
	return nil
}

// InvalidateByAtomRevision invalidates L2 and broadcasts an invalidation
// event so all nodes purge their L1 caches.
func (t *TieredGridCache) InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error {
	if err := t.l2.InvalidateByAtomRevision(ctx, planID, atomRevision); err != nil {
		return err
	}

	ev := InvalidationEvent{
		PlanID:       planID,
		AtomRevision: atomRevision,
		SourceNodeID: t.nodeID,
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	channel := t.invalidationChannel(planID, atomRevision)
	if t.redis != nil {
		if err := t.redis.Publish(ctx, channel, b).Err(); err != nil {
			return err
		}
	}

	// Local L1 purge as well.
	t.l1.Purge()
	return nil
}

// l1Key builds a deterministic logical key for L1 cache.
func (t *TieredGridCache) l1Key(key GridCacheKey) string {
	var sb strings.Builder
	sb.Grow(128)
	sb.WriteString(key.PlanID)
	sb.WriteByte('|')
	sb.WriteString(key.ViewID)
	sb.WriteByte('|')
	sb.WriteString(key.UserID)
	sb.WriteByte('|')
	sb.WriteString(key.ScenarioID)
	sb.WriteByte('|')
	sb.WriteString(strconv.FormatInt(key.AtomRevision, 10))
	sb.WriteByte('|')
	sb.WriteString(key.WindowHash)
	if key.AdditionalTag != "" {
		sb.WriteByte('|')
		sb.WriteString(key.AdditionalTag)
	}
	return sb.String()
}

// invalidationChannel builds a plan-scoped Pub/Sub channel name.
func (t *TieredGridCache) invalidationChannel(planID string, atomRevision int64) string {
	var sb strings.Builder
	sb.Grow(len(t.pubsubPref) + len(planID) + 32)
	sb.WriteString(t.pubsubPref)
	sb.WriteByte('{')
	sb.WriteString(planID)
	sb.WriteByte('}')
	sb.WriteByte(':')
	sb.WriteString(strconv.FormatInt(atomRevision, 10))
	return sb.String()
}
