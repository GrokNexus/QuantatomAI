package planner

import (
    "sync"
    "time"
)

//
// ===============================
//   IN-MEMORY LRU + TTL CACHE
// ===============================
//

// cachedAxes stores the actual projection.
type cachedAxes struct {
    Rows    [][]MemberInfo
    RowIDs  [][]int64
    Cols    [][]MemberInfo
    ColIDs  [][]int64
    Expiry  time.Time
}

// lruNode is a doubly-linked list node for LRU eviction.
type lruNode struct {
    key  AxisKey
    prev *lruNode
    next *lruNode
}

// InMemoryProjectionCache is a TTL + LRU bounded cache.
type InMemoryProjectionCache struct {
    mu sync.Mutex

    // Config
    ttl       time.Duration
    maxSize   int

    // Storage
    items map[AxisKey]*cacheEntry

    // LRU pointers
    head *lruNode
    tail *lruNode
}

type cacheEntry struct {
    data cachedAxes
    node *lruNode
}

// NewInMemoryProjectionCache creates a bounded LRU cache with TTL.
func NewInMemoryProjectionCache(ttl time.Duration, maxSize int) *InMemoryProjectionCache {
    if maxSize <= 0 {
        maxSize = 1000 // sane default
    }
    if ttl <= 0 {
        ttl = 5 * time.Minute
    }

    return &InMemoryProjectionCache{
        ttl:     ttl,
        maxSize: maxSize,
        items:   make(map[AxisKey]*cacheEntry),
    }
}

//
// ===============================
//   PUBLIC METHODS
// ===============================
//

func (c *InMemoryProjectionCache) GetAxes(
    key AxisKey,
) ([][]MemberInfo, [][]int64, [][]MemberInfo, [][]int64, bool) {

    c.mu.Lock()
    defer c.mu.Unlock()

    entry, ok := c.items[key]
    if !ok {
        return nil, nil, nil, nil, false
    }

    // TTL check
    if time.Now().After(entry.data.Expiry) {
        c.removeEntry(key)
        return nil, nil, nil, nil, false
    }

    // Move to front (LRU)
    c.moveToFront(entry.node)

    return entry.data.Rows, entry.data.RowIDs, entry.data.Cols, entry.data.ColIDs, true
}

func (c *InMemoryProjectionCache) SetAxes(
    key AxisKey,
    rows [][]MemberInfo,
    rowIDs [][]int64,
    cols [][]MemberInfo,
    colIDs [][]int64,
) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // If exists, update + move to front
    if entry, ok := c.items[key]; ok {
        entry.data = cachedAxes{
            Rows:   rows,
            RowIDs: rowIDs,
            Cols:   cols,
            ColIDs: colIDs,
            Expiry: time.Now().Add(c.ttl),
        }
        c.moveToFront(entry.node)
        return
    }

    // Evict if full
    if len(c.items) >= c.maxSize {
        c.evictLRU()
    }

    // Insert new
    node := &lruNode{key: key}
    c.insertFront(node)

    c.items[key] = &cacheEntry{
        data: cachedAxes{
            Rows:   rows,
            RowIDs: rowIDs,
            Cols:   cols,
            ColIDs: colIDs,
            Expiry: time.Now().Add(c.ttl),
        },
        node: node,
    }
}

//
// ===============================
//   INTERNAL LRU OPERATIONS
// ===============================
//

func (c *InMemoryProjectionCache) insertFront(n *lruNode) {
    n.prev = nil
    n.next = c.head
    if c.head != nil {
        c.head.prev = n
    }
    c.head = n
    if c.tail == nil {
        c.tail = n
    }
}

func (c *InMemoryProjectionCache) moveToFront(n *lruNode) {
    if c.head == n {
        return
    }
    // unlink
    if n.prev != nil {
        n.prev.next = n.next
    }
    if n.next != nil {
        n.next.prev = n.prev
    }
    if c.tail == n {
        c.tail = n.prev
    }
    // insert at front
    n.prev = nil
    n.next = c.head
    if c.head != nil {
        c.head.prev = n
    }
    c.head = n
}

func (c *InMemoryProjectionCache) evictLRU() {
    if c.tail == nil {
        return
    }
    key := c.tail.key
    c.removeEntry(key)
}

func (c *InMemoryProjectionCache) removeEntry(key AxisKey) {
    entry, ok := c.items[key]
    if !ok {
        return
    }
    n := entry.node

    // unlink
    if n.prev != nil {
        n.prev.next = n.next
    }
    if n.next != nil {
        n.next.prev = n.prev
    }
    if c.head == n {
        c.head = n.next
    }
    if c.tail == n {
        c.tail = n.prev
    }

    delete(c.items, key)
}

//
// ===============================
//   REDIS-BASED CROSS-INSTANCE CACHE (stub)
// ===============================
//

// RedisProjectionCache is a placeholder for a distributed cache.
type RedisProjectionCache struct {
    ttl time.Duration
}

func NewRedisProjectionCache(ttl time.Duration) *RedisProjectionCache {
    return &RedisProjectionCache{ttl: ttl}
}

func (r *RedisProjectionCache) GetAxes(key AxisKey) ([][]MemberInfo, [][]int64, [][]MemberInfo, [][]int64, bool) {
    // TODO: Implement Redis lookup
    return nil, nil, nil, nil, false
}

func (r *RedisProjectionCache) SetAxes(key AxisKey, rows [][]MemberInfo, rowIDs [][]int64, cols [][]MemberInfo, colIDs [][]int64) {
    // TODO: Implement Redis write
}
