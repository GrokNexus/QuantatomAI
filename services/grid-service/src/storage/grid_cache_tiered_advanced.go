package storage

import (
    "context"
    "encoding/json"
    "log"
    "math"
    "math/rand"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"

    "github.com/dgraph-io/ristretto"
    "github.com/redis/go-redis/v9"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

var (
    cacheMeter           = otel.Meter("quantatomai/grid-cache")
    l1HitCounter         metric.Int64Counter
    l2HitCounter         metric.Int64Counter
    cacheMissCounter     metric.Int64Counter
    cacheInitMetricsOnce sync.Once
)

func initCacheMetrics() {
    cacheInitMetricsOnce.Do(func() {
        var err error
        l1HitCounter, err = cacheMeter.Int64Counter("grid_cache_l1_hit")
        if err != nil {
            log.Printf("grid-cache: failed to create l1_hit counter: %v", err)
        }
        l2HitCounter, err = cacheMeter.Int64Counter("grid_cache_l2_hit")
        if err != nil {
            log.Printf("grid-cache: failed to create l2_hit counter: %v", err)
        }
        cacheMissCounter, err = cacheMeter.Int64Counter("grid_cache_miss")
        if err != nil {
            log.Printf("grid-cache: failed to create miss counter: %v", err)
        }
    })
}

// GoldenWindowPrefetcher is called after invalidation to refill hot windows.
type GoldenWindowPrefetcher interface {
    PrefetchHotWindows(ctx context.Context, planID, viewID string, version int64)
}

// PrefetchCoordinator decides whether this node should perform prefetch work.
// This enables distributed coordination (e.g., only one node per plan).
type PrefetchCoordinator interface {
    ShouldPrefetch(planID string) bool
}

// XFetchRefresher is used by X-Fetch to refresh keys in the background.
type XFetchRefresher interface {
    Refresh(ctx context.Context, key GridCacheKey)
}

// hotKeyStats tracks approximate frequency per (planID|viewID|windowHash|atomRevision).
type hotKeyStats struct {
    mu   sync.RWMutex
    data map[string]int64
}

func newHotKeyStats() *hotKeyStats {
    return &hotKeyStats{
        data: make(map[string]int64),
    }
}

func (h *hotKeyStats) Inc(key string) {
    h.mu.Lock()
    h.data[key]++
    h.mu.Unlock()
}

func (h *hotKeyStats) Frequency(key string) int64 {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.data[key]
}

func (h *hotKeyStats) TopNForPlanView(planID, viewID string, n int) []string {
    prefix := planID + "|" + viewID + "|"
    h.mu.RLock()
    defer h.mu.RUnlock()

    type kv struct {
        k string
        v int64
    }
    var arr []kv
    for k, v := range h.data {
        if strings.HasPrefix(k, prefix) {
            arr = append(arr, kv{k: k, v: v})
        }
    }
    if len(arr) == 0 {
        return nil
    }
    if len(arr) > n {
        for i := 0; i < n; i++ {
            maxIdx := i
            for j := i + 1; j < len(arr); j++ {
                if arr[j].v > arr[maxIdx].v {
                    maxIdx = j
                }
            }
            arr[i], arr[maxIdx] = arr[maxIdx], arr[i]
        }
        arr = arr[:n]
    }
    out := make([]string, 0, len(arr))
    for _, kv := range arr {
        out = append(out, kv.k)
    }
    return out
}

// DefaultGoldenWindowPrefetcher uses hotKeyStats and a recompute callback.
type DefaultGoldenWindowPrefetcher struct {
    stats       *hotKeyStats
    prefetchN   int
    recomputeFn func(ctx context.Context, planID, viewID, windowHash string, version int64)
}

func NewDefaultGoldenWindowPrefetcher(
    stats *hotKeyStats,
    prefetchN int,
    recomputeFn func(ctx context.Context, planID, viewID, windowHash string, version int64),
) *DefaultGoldenWindowPrefetcher {
    if prefetchN <= 0 {
        prefetchN = 5
    }
    return &DefaultGoldenWindowPrefetcher{
        stats:       stats,
        prefetchN:   prefetchN,
        recomputeFn: recomputeFn,
    }
}

func (p *DefaultGoldenWindowPrefetcher) PrefetchHotWindows(ctx context.Context, planID, viewID string, version int64) {
    if p.recomputeFn == nil {
        return
    }
    keys := p.stats.TopNForPlanView(planID, viewID, p.prefetchN)
    for _, k := range keys {
        parts := strings.SplitN(k, "|", 4)
        if len(parts) != 4 {
            continue
        }
        p.recomputeFn(ctx, parts[0], parts[1], parts[2], version)
    }
}

// TieredGridCache is the advanced tiered cache engine:
// - L1 (Ristretto, memory-capped)
// - L2 (GridCache, e.g. Redis)
// - Pub/Sub invalidation with per-plan/view versioning
// - Lock-free version map
// - Golden Window prefetching
// - X-Fetch (probabilistic early refresh)
// - Adaptive TTLs
// - Version rotation / maintenance
// - Warm-up siphoning
// - Telemetry
type TieredGridCache struct {
    l1         *ristretto.Cache
    l2         GridCache
    redis      *redis.Client
    pubsubPref string
    nodeID     string

    versionMap atomic.Value // map[string]int64

    subscriber *redis.PubSub
    cancelFn   context.CancelFunc
    subscribed bool

    rngMu sync.Mutex
    rng   *rand.Rand

    prefetcher          GoldenWindowPrefetcher
    prefetchCoordinator PrefetchCoordinator
    hotStats            *hotKeyStats

    xfetchRefresher XFetchRefresher
    xfetchProb      float64 // 0.0–1.0

    maintenanceOnce sync.Once
}

// InvalidationEvent is broadcast over Pub/Sub.
type InvalidationEvent struct {
    PlanID       string `json:"plan_id"`
    ViewID       string `json:"view_id,omitempty"`
    AtomRevision int64  `json:"atom_revision"`
    SourceNodeID string `json:"source_node_id"`
    Version      int64  `json:"version"`
}

// NewTieredGridCache constructs the advanced tiered cache.
//
// maxMemory should be ~25% of pod RAM.
// xfetchProb is the probability (0–1) of triggering X-Fetch on a hit.
// xfetchRefresher may be nil if X-Fetch is not yet wired.
func NewTieredGridCache(
    l2 GridCache,
    redis *redis.Client,
    maxMemory int64,
    pubsubPref string,
    nodeID string,
    prefetcher GoldenWindowPrefetcher,
    prefetchCoordinator PrefetchCoordinator,
    xfetchRefresher XFetchRefresher,
    xfetchProb float64,
) (*TieredGridCache, error) {
    if maxMemory <= 0 {
        maxMemory = 256 * 1024 * 1024
    }
    if pubsubPref == "" {
        pubsubPref = "grid-cache:invalidate:"
    }
    if nodeID == "" {
        nodeID = "node-unknown"
    }
    if xfetchProb < 0 {
        xfetchProb = 0
    }
    if xfetchProb > 1 {
        xfetchProb = 1
    }

    l1, err := ristretto.NewCache(&ristretto.Config{
        NumCounters: 1_000_000,
        MaxCost:     maxMemory,
        BufferItems: 64,
        Cost: func(value interface{}) int64 {
            entry, ok := value.(*versionedEntry)
            if !ok || entry == nil || entry.Entry == nil {
                return 1
            }
            return int64(len(entry.Entry.Payload) + 256)
        },
    })
    if err != nil {
        return nil, err
    }

    initCacheMetrics()

    src := rand.NewSource(time.Now().UnixNano())

    t := &TieredGridCache{
        l1:                  l1,
        l2:                  l2,
        redis:               redis,
        pubsubPref:          pubsubPref,
        nodeID:              nodeID,
        rng:                 rand.New(src),
        prefetcher:          prefetcher,
        prefetchCoordinator: prefetchCoordinator,
        hotStats:            newHotKeyStats(),
        xfetchRefresher:     xfetchRefresher,
        xfetchProb:          xfetchProb,
    }
    t.versionMap.Store(make(map[string]int64))
    return t, nil
}

type versionedEntry struct {
    VersionKey string
    Version    int64
    Entry      *GridCacheEntry
}

// StartInvalidationSubscriber starts Pub/Sub listener and maintenance loop.
func (t *TieredGridCache) StartInvalidationSubscriber(parent context.Context) error {
    if t.redis == nil {
        return nil
    }
    if t.subscribed {
        return nil
    }

    ctx, cancel := context.WithCancel(parent)
    t.cancelFn = cancel

    pattern := t.pubsubPref + "*"
    ps := t.redis.PSubscribe(ctx, pattern)
    t.subscriber = ps
    t.subscribed = true

    go t.runSubscriberLoop(ctx, ps)
    t.startMaintenanceLoop(ctx)

    return nil
}

func (t *TieredGridCache) runSubscriberLoop(ctx context.Context, ps *redis.PubSub) {
    backoff := time.Second
    maxBackoff := 30 * time.Second

    for {
        msg, err := ps.ReceiveMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            log.Printf("tiered-cache: pubsub receive error: %v", err)

            t.rngMu.Lock()
            jitter := time.Duration(t.rng.Int63n(int64(backoff / 2)))
            t.rngMu.Unlock()
            sleep := backoff + jitter
            if sleep > maxBackoff {
                sleep = maxBackoff
            }
            time.Sleep(sleep)

            backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
            continue
        }

        backoff = time.Second

        var ev InvalidationEvent
        if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
            log.Printf("tiered-cache: invalidation event unmarshal error: %v", err)
            continue
        }

        versionKey := t.versionKey(ev.PlanID, ev.ViewID)
        t.setVersion(versionKey, ev.Version)

        if t.prefetcher != nil && t.shouldPrefetch(ev.PlanID) {
            go t.prefetcher.PrefetchHotWindows(ctx, ev.PlanID, ev.ViewID, ev.Version)
        }
    }
}

func (t *TieredGridCache) shouldPrefetch(planID string) bool {
    if t.prefetchCoordinator == nil {
        return true
    }
    return t.prefetchCoordinator.ShouldPrefetch(planID)
}

// startMaintenanceLoop periodically compacts versionMap.
func (t *TieredGridCache) startMaintenanceLoop(ctx context.Context) {
    t.maintenanceOnce.Do(func() {
        go func() {
            ticker := time.NewTicker(30 * time.Minute)
            defer ticker.Stop()

            for {
                select {
                case <-ctx.Done():
                    return
                case <-ticker.C:
                    t.pruneVersionMap()
                }
            }
        }()
    })
}

func (t *TieredGridCache) pruneVersionMap() {
    for {
        old := t.versionMap.Load().(map[string]int64)
        if len(old) == 0 {
            return
        }
        newMap := make(map[string]int64, len(old))
        for k, v := range old {
            newMap[k] = v
        }
        if t.versionMap.CompareAndSwap(old, newMap) {
            return
        }
    }
}

// Close shuts down subscriber.
func (t *TieredGridCache) Close() error {
    if t.cancelFn != nil {
        t.cancelFn()
    }
    if t.subscriber != nil {
        if err := t.subscriber.Close(); err != nil {
            return err
        }
    }
    t.subscribed = false
    return nil
}

// Get with L1/L2 tiering, version-awareness, telemetry, hot-key tracking, and X-Fetch trigger.
func (t *TieredGridCache) Get(ctx context.Context, key GridCacheKey) (*GridCacheEntry, bool, error) {
    l1Key := t.l1Key(key)
    vKey := t.versionKey(key.PlanID, key.ViewID)
    version := t.getVersion(vKey)

    if v, ok := t.l1.Get(l1Key); ok {
        if ve, ok2 := v.(*versionedEntry); ok2 && ve != nil && ve.Entry != nil {
            if ve.VersionKey == vKey && ve.Version == version {
                if l1HitCounter != (metric.Int64Counter{}) {
                    l1HitCounter.Add(ctx, 1, metric.WithAttributes(
                        attribute.String("plan_id", key.PlanID),
                        attribute.String("view_id", key.ViewID),
                    ))
                }
                t.recordHotKey(key)
                t.maybeTriggerXFetch(ctx, key)
                return ve.Entry, true, nil
            }
        }
        t.l1.Del(l1Key)
    }

    entry, found, err := t.l2.Get(ctx, key)
    if err != nil || !found || entry == nil {
        if cacheMissCounter != (metric.Int64Counter{}) {
            cacheMissCounter.Add(ctx, 1, metric.WithAttributes(
                attribute.String("plan_id", key.PlanID),
                attribute.String("view_id", key.ViewID),
            ))
        }
        return nil, false, err
    }

    if l2HitCounter != (metric.Int64Counter{}) {
        l2HitCounter.Add(ctx, 1, metric.WithAttributes(
            attribute.String("plan_id", key.PlanID),
            attribute.String("view_id", key.ViewID),
        ))
    }

    ve := &versionedEntry{
        VersionKey: vKey,
        Version:    version,
        Entry:      entry,
    }
    t.l1.Set(l1Key, ve, 0)
    t.recordHotKey(key)
    t.maybeTriggerXFetch(ctx, key)
    return entry, true, nil
}

// Set with L1/L2 tiering and adaptive TTL.
func (t *TieredGridCache) Set(ctx context.Context, key GridCacheKey, entry *GridCacheEntry, ttl time.Duration) error {
    ttl = t.adaptTTL(key, ttl)

    if err := t.l2.Set(ctx, key, entry, ttl); err != nil {
        return err
    }
    vKey := t.versionKey(key.PlanID, key.ViewID)
    version := t.getVersion(vKey)
    ve := &versionedEntry{
        VersionKey: vKey,
        Version:    version,
        Entry:      entry,
    }
    t.l1.Set(t.l1Key(key), ve, 0)
    return nil
}

// InvalidateByAtomRevision (plan-level).
func (t *TieredGridCache) InvalidateByAtomRevision(ctx context.Context, planID string, atomRevision int64) error {
    return t.invalidateInternal(ctx, planID, "", atomRevision)
}

// InvalidateByPlanView (view-level).
func (t *TieredGridCache) InvalidateByPlanView(ctx context.Context, planID, viewID string, atomRevision int64) error {
    return t.invalidateInternal(ctx, planID, viewID, atomRevision)
}

func (t *TieredGridCache) invalidateInternal(ctx context.Context, planID, viewID string, atomRevision int64) error {
    vKey := t.versionKey(planID, viewID)
    newVersion := t.bumpVersion(vKey)

    if err := t.l2.InvalidateByAtomRevision(ctx, planID, atomRevision); err != nil {
        return err
    }

    ev := InvalidationEvent{
        PlanID:       planID,
        ViewID:       viewID,
        AtomRevision: atomRevision,
        SourceNodeID: t.nodeID,
        Version:      newVersion,
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

    return nil
}

// PrewarmKey supports warm-up siphoning.
func (t *TieredGridCache) PrewarmKey(ctx context.Context, key GridCacheKey) {
    _, _, _ = t.Get(ctx, key)
}

func (t *TieredGridCache) recordHotKey(key GridCacheKey) {
    if t.hotStats == nil {
        return
    }
    var sb strings.Builder
    sb.Grow(128)
    sb.WriteString(key.PlanID)
    sb.WriteByte('|')
    sb.WriteString(key.ViewID)
    sb.WriteByte('|')
    sb.WriteString(key.WindowHash)
    sb.WriteByte('|')
    sb.WriteString(strconv.FormatInt(key.AtomRevision, 10))
    t.hotStats.Inc(sb.String())
}

// maybeTriggerXFetch triggers background refresh with probability xfetchProb.
func (t *TieredGridCache) maybeTriggerXFetch(ctx context.Context, key GridCacheKey) {
    if t.xfetchRefresher == nil || t.xfetchProb <= 0 {
        return
    }
    t.rngMu.Lock()
    p := t.rng.Float64()
    t.rngMu.Unlock()
    if p <= t.xfetchProb {
        go t.xfetchRefresher.Refresh(ctx, key)
    }
}

// adaptTTL adjusts TTL based on hotness (simple adaptive TTL).
func (t *TieredGridCache) adaptTTL(key GridCacheKey, baseTTL time.Duration) time.Duration {
    if t.hotStats == nil || baseTTL <= 0 {
        return baseTTL
    }
    var sb strings.Builder
    sb.Grow(128)
    sb.WriteString(key.PlanID)
    sb.WriteByte('|')
    sb.WriteString(key.ViewID)
    sb.WriteByte('|')
    sb.WriteString(key.WindowHash)
    sb.WriteByte('|')
    sb.WriteString(strconv.FormatInt(key.AtomRevision, 10))
    freq := t.hotStats.Frequency(sb.String())
    if freq <= 0 {
        return baseTTL
    }
    switch {
    case freq >= 100:
        return time.Duration(float64(baseTTL) * 2.0)
    case freq >= 20:
        return time.Duration(float64(baseTTL) * 1.5)
    default:
        return baseTTL
    }
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

// invalidationChannel uses a Redis Cluster hash tag {planID}.
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

// versionKey composes the logical version key (plan-level or plan+view).
func (t *TieredGridCache) versionKey(planID, viewID string) string {
    if viewID == "" {
        return planID
    }
    return planID + "|" + viewID
}

func (t *TieredGridCache) getVersion(key string) int64 {
    m := t.versionMap.Load().(map[string]int64)
    if v, ok := m[key]; ok {
        return v
    }
    return 0
}

func (t *TieredGridCache) setVersion(key string, version int64) {
    for {
        old := t.versionMap.Load().(map[string]int64)
        newMap := make(map[string]int64, len(old)+1)
        for k, v := range old {
            newMap[k] = v
        }
        newMap[key] = version
        if t.versionMap.CompareAndSwap(old, newMap) {
            return
        }
    }
}

func (t *TieredGridCache) bumpVersion(key string) int64 {
    for {
        old := t.versionMap.Load().(map[string]int64)
        newMap := make(map[string]int64, len(old)+1)
        var cur int64
        for k, v := range old {
            newMap[k] = v
            if k == key {
                cur = v
            }
        }
        cur++
        newMap[key] = cur
        if t.versionMap.CompareAndSwap(old, newMap) {
            return cur
        }
    }
}
