package fetcher

import (
	"context"
	"errors"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"quantatomai/grid-service/domain"
	"quantatomai/grid-service/planner"
)

// CircuitBreakerState represents the current state of the fetcher's resiliency breaker.
type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

// RedisAtomFetcher is a resilient, pipelined, chunked, and circuit-breaker-protected hot-store fetcher.
// It dynamically tunes itself based on system resources and query magnitude.
type RedisAtomFetcher struct {
	client *redis.Client
	prefix string

	// Dynamic tuning configuration
	baseChunkSize int
	maxChunkSize  int
	minChunkSize  int

	baseConcurrency int
	maxConcurrency  int
	minConcurrency  int

	timeout time.Duration

	// Circuit breaker state
	breakerMu        sync.Mutex
	breakerState     CircuitBreakerState
	failureCount     int
	failureThreshold int
	openUntil        time.Time
	openDuration     time.Duration
}

// NewRedisAtomFetcher constructs a production-ready atom fetcher with adaptive defaults.
func NewRedisAtomFetcher(client *redis.Client, prefix string, timeout time.Duration) *RedisAtomFetcher {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if prefix == "" {
		prefix = "atom:"
	}

	cpu := runtime.NumCPU()

	return &RedisAtomFetcher{
		client: client,
		prefix: prefix,

		// Dynamic chunk sizing bounds
		baseChunkSize: 5000,
		maxChunkSize:  15000,
		minChunkSize:  1000,

		// Dynamic concurrency based on available CPU
		baseConcurrency: int(math.Max(2, float64(cpu/2))),
		maxConcurrency:  int(math.Max(4, float64(cpu))),
		minConcurrency:  1,

		timeout: timeout,

		// Circuit breaker defaults (5 failures to open, 5s recovery window)
		breakerState:     Closed,
		failureThreshold: 5,
		openDuration:     5 * time.Second,
	}
}

// -----------------------------
// Circuit Breaker Logic
// -----------------------------

func (f *RedisAtomFetcher) allowRequest() bool {
	f.breakerMu.Lock()
	defer f.breakerMu.Unlock()

	switch f.breakerState {
	case Open:
		if time.Now().After(f.openUntil) {
			f.breakerState = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	default:
		return true
	}
}

func (f *RedisAtomFetcher) recordSuccess() {
	f.breakerMu.Lock()
	defer f.breakerMu.Unlock()

	f.failureCount = 0
	if f.breakerState == HalfOpen {
		f.breakerState = Closed
	}
}

func (f *RedisAtomFetcher) recordFailure() {
	f.breakerMu.Lock()
	defer f.breakerMu.Unlock()

	f.failureCount++
	if f.failureCount >= f.failureThreshold {
		f.breakerState = Open
		f.openUntil = time.Now().Add(f.openDuration)
	}
}

// -----------------------------
// FetchAtoms Implementation
// -----------------------------

// FetchAtoms retrieves atoms from Redis with pipelining, concurrency, and circuit breaking.
func (f *RedisAtomFetcher) FetchAtoms(
	ctx context.Context,
	block planner.SparsedBlockRef,
) (map[domain.AtomKey]float64, error) {

	// 1. Guard with Circuit Breaker
	if !f.allowRequest() {
		return nil, errors.New("redis hot-store temporarily degraded (circuit open)")
	}

	// 2. Expand Block to Keys
	keys := expandBlockToAtomKeys(block)
	if len(keys) == 0 {
		return map[domain.AtomKey]float64{}, nil
	}

	// 3. Prepare optimized Redis keys
	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = f.redisKeyFor(k)
	}

	// 4. Determine Dynamic Tuning Parameters
	chunkSize := f.dynamicChunkSize(len(keys))
	concurrency := f.dynamicConcurrency()

	chunks := chunkStrings(redisKeys, chunkSize)

	out := make(map[domain.AtomKey]float64, len(keys))
	var mu sync.Mutex

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	var firstErr error
	var errMu sync.Mutex

	// 5. Concurrent Fetcher Loop
	for chunkIdx, chunk := range chunks {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int, c []string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Per-chunk timeout
			batchCtx, cancel := context.WithTimeout(ctx, f.timeout)
			defer cancel()

			// Pipeline Execution
			pipe := f.client.Pipeline()
			cmd := pipe.MGet(batchCtx, c...)
			_, err := pipe.Exec(batchCtx)

			if err != nil {
				f.recordFailure()
				errMu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("pipeline chunk %d failed: %w", idx, err)
				}
				errMu.Unlock()
				return
			}

			vals := cmd.Val()
			if len(vals) != len(c) {
				f.recordFailure()
				return
			}

			// Local result mapping to reduce lock contention
			localMap := make(map[domain.AtomKey]float64, len(vals))
			for i, raw := range vals {
				if raw == nil {
					continue
				}
				strVal, ok := raw.(string)
				if !ok {
					continue
				}
				
				// High-performance float parsing
				v, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					continue
				}

				key := keys[idx*chunkSize+i]
				localMap[key] = v
			}

			// Bulk merge into output
			mu.Lock()
			for k, v := range localMap {
				out[k] = v
			}
			mu.Unlock()

		}(chunkIdx, chunk)
	}

	wg.Wait()

	if firstErr != nil {
		return out, firstErr
	}

	f.recordSuccess()
	return out, nil
}

// -----------------------------
// Dynamic Tuning Logic
// -----------------------------

func (f *RedisAtomFetcher) dynamicChunkSize(total int) int {
	if total < f.minChunkSize {
		return f.minChunkSize
	}
	if total > 50000 {
		return f.maxChunkSize
	}
	return f.baseChunkSize
}

func (f *RedisAtomFetcher) dynamicConcurrency() int {
	cpu := runtime.NumCPU()
	c := int(math.Min(float64(f.maxConcurrency), math.Max(float64(f.minConcurrency), float64(cpu/2))))
	return c
}

// -----------------------------
// Internal Helpers
// -----------------------------

// redisKeyFor encodes an AtomKey into a stable Redis key using strings.Builder.
// Optimized for zero-allocation hotspots.
func (f *RedisAtomFetcher) redisKeyFor(k domain.AtomKey) string {
	var sb strings.Builder
	sb.Grow(len(f.prefix) + 32 + (k.DimCount * 11))
	
	sb.WriteString(f.prefix)
	sb.WriteString(strconv.FormatInt(k.MeasureID, 10))
	sb.WriteByte(':')
	sb.WriteString(strconv.FormatInt(k.ScenarioID, 10))
	
	for i := 0; i < k.DimCount; i++ {
		sb.WriteByte(':')
		sb.WriteString(strconv.FormatInt(k.DimIDs[i], 10))
	}
	return sb.String()
}

// chunkStrings splits a slice into smaller chunks.
func chunkStrings(arr []string, size int) [][]string {
	if size <= 0 {
		size = 5000
	}
	var chunks [][]string
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// expandBlockToAtomKeys expands a SparsedBlockRef into concrete AtomKeys.
// Handles multi-axial blocks while ensuring canonical ordering.
func expandBlockToAtomKeys(block planner.SparsedBlockRef) []domain.AtomKey {
	// 1. Identify constant filter IDs to include in coordinates
	var baseDimIDs []int64
	filterDims := make([]string, 0, len(block.Filters))
	for k := range block.Filters {
		filterDims = append(filterDims, k)
	}
	sort.Strings(filterDims)
	
	for _, dim := range filterDims {
		ids := block.Filters[dim]
		if len(ids) == 1 {
			baseDimIDs = append(baseDimIDs, ids[0])
		}
	}

	// 2. Cartesian Product of Block IDs
	var out []domain.AtomKey
	for _, r := range block.RowIDs {
		for _, c := range block.ColIDs {
			for _, m := range block.MeasureIDs {
				for _, s := range block.ScenarioIDs {
					ak := domain.AtomKey{
						DimCount:   int(math.Min(8, float64(2+len(baseDimIDs)))),
						MeasureID:  m,
						ScenarioID: s,
					}
					if ak.DimCount >= 2 {
						ak.DimIDs[0] = r
						ak.DimIDs[1] = c
						copy(ak.DimIDs[2:], baseDimIDs)
					}
					ak.EnsureCanonical()
					out = append(out, ak)
				}
			}
		}
	}
	return out
}
