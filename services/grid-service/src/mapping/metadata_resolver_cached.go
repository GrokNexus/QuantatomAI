package mapping

import (
    "context"
    "sort"
    "strings"
    "sync"
    "time"

    "quantatomai/grid-service/planner"
)

// CachedMetadataResolver wraps a MetadataResolver with TTL memoization.
// It uses a decorator pattern to add performance without modifying the underlying storage logic.
type CachedMetadataResolver struct {
    inner   planner.MetadataResolver
    ttl     time.Duration
    maxSize int

    mu    sync.Mutex
    cache map[string]metaCacheEntry
}

type metaCacheEntry struct {
    value any
    exp   time.Time
}

func NewCachedMetadataResolver(inner planner.MetadataResolver, ttl time.Duration, maxSize int) *CachedMetadataResolver {
    if ttl <= 0 {
        ttl = 5 * time.Minute
    }
    if maxSize <= 0 {
        maxSize = 10000
    }
    return &CachedMetadataResolver{
        inner:   inner,
        ttl:     ttl,
        maxSize: maxSize,
        cache:   make(map[string]metaCacheEntry),
    }
}

func (c *CachedMetadataResolver) ResolveMembers(ctx context.Context, dim string, codes []string) ([]planner.MemberInfo, error) {
    key := c.keyFor("members", dim, codes)
    if v, ok := c.get(key); ok {
        return v.([]planner.MemberInfo), nil
    }
    res, err := c.inner.ResolveMembers(ctx, dim, codes)
    if err != nil {
        return nil, err
    }
    c.set(key, res)
    return res, nil
}

func (c *CachedMetadataResolver) ResolveMeasureIDs(ctx context.Context, measures []string) ([]int64, error) {
    key := c.keyFor("measures", "", measures)
    if v, ok := c.get(key); ok {
        return v.([]int64), nil
    }
    res, err := c.inner.ResolveMeasureIDs(ctx, measures)
    if err != nil {
        return nil, err
    }
    c.set(key, res)
    return res, nil
}

func (c *CachedMetadataResolver) ResolveScenarioIDs(ctx context.Context, scenarios []string) ([]int64, error) {
    key := c.keyFor("scenarios", "", scenarios)
    if v, ok := c.get(key); ok {
        return v.([]int64), nil
    }
    res, err := c.inner.ResolveScenarioIDs(ctx, scenarios)
    if err != nil {
        return nil, err
    }
    c.set(key, res)
    return res, nil
}

func (c *CachedMetadataResolver) keyFor(prefix, dim string, codes []string) string {
    // Sort codes to ensure deterministic keys regardless of request order.
    sortedCodes := make([]string, len(codes))
    copy(sortedCodes, codes)
    sort.Strings(sortedCodes)

    var sb strings.Builder
    sb.WriteString(prefix)
    sb.WriteString("|")
    sb.WriteString(dim)
    sb.WriteString("|")
    for i, code := range sortedCodes {
        if i > 0 {
            sb.WriteString(",")
        }
        sb.WriteString(code)
    }
    return sb.String()
}

func (c *CachedMetadataResolver) get(key string) (any, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    entry, ok := c.cache[key]
    if !ok {
        return nil, false
    }
    if time.Now().After(entry.exp) {
        delete(c.cache, key)
        return nil, false
    }
    return entry.value, true
}

func (c *CachedMetadataResolver) set(key string, value any) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if len(c.cache) >= c.maxSize {
        // Simple random eviction; for heavy load, consider an LRU implementation.
        for k := range c.cache {
            delete(c.cache, k)
            break
        }
    }
    c.cache[key] = metaCacheEntry{
        value: value,
        exp:   time.Now().Add(c.ttl),
    }
}
