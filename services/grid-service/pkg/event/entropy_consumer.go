package event

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
)

const NumShards = 64

type EntropyShard struct {
	mu          sync.RWMutex
	rollingMean map[string]float64
	rollingM2   map[string]float64
	counts      map[string]int
}

// EntropyStreamer taps into the Redpanda 'grid_events' backbone to detect
// real-time statistical anomalies (Z-Score > 3) instantly as users make edits.
// This is Phase 8.2 of the Fluxion AI Engine.
type EntropyStreamer struct {
	bus   *KafkaBus
	topic string

	shards [NumShards]*EntropyShard
}

func NewEntropyStreamer(brokers []string, topic string) *EntropyStreamer {
	e := &EntropyStreamer{
		bus:   NewKafkaBus(brokers, topic),
		topic: topic,
	}
	for i := 0; i < NumShards; i++ {
		e.shards[i] = &EntropyShard{
			rollingMean: make(map[string]float64),
			rollingM2:   make(map[string]float64),
			counts:      make(map[string]int),
		}
	}
	return e
}

func (e *EntropyStreamer) Start(ctx context.Context) error {
	log.Println("[ENTROPY] Initializing Maxwell Daemon on Redpanda Consumer...")

	eventsChan, err := e.bus.Subscribe(ctx, e.topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to entropy topic: %w", err)
	}

	go func() {
		for evt := range eventsChan {
			// In production, we decode the FlatBuffer Payload.
			// e.g. val := DecodeFlatMolecule(evt.Payload).Value()

			// Simulate parsed cell value and coordinate hash
			cellHash := evt.TraceID // using TraceID as a dummy coordinate identifier

			// For simulation, let's assume we can parse 'Payload' into a float64
			// if it's an EDIT event. We'll mock the value here based on payload length.
			simulatedVal := float64(len(evt.Payload)) * 1000.0

			e.observeAndDetect(cellHash, simulatedVal)
		}
	}()

	return nil
}

// simpleHash provides a fast basic hash for sharding
func simpleHash(s string) uint32 {
	h := uint32(0)
	for i := 0; i < len(s); i++ {
		h = 31*h + uint32(s[i])
	}
	return h
}

// observeAndDetect updates Welford's online algorithm for variance
// utilizing 64-way Mutex Sharding (Ultra-Diamond)
func (e *EntropyStreamer) observeAndDetect(cellHash string, value float64) {
	shardIdx := simpleHash(cellHash) % NumShards
	shard := e.shards[shardIdx]

	shard.mu.Lock()
	defer shard.mu.Unlock()

	count := shard.counts[cellHash]
	mean := shard.rollingMean[cellHash]
	m2 := shard.rollingM2[cellHash]

	count++
	delta := value - mean
	mean += delta / float64(count)
	delta2 := value - mean
	m2 += delta * delta2

	shard.counts[cellHash] = count
	shard.rollingMean[cellHash] = mean
	shard.rollingM2[cellHash] = m2

	// Need at least a small sample size to compute reliable Z-scores
	if count > 5 {
		variance := m2 / float64(count-1)
		stddev := math.Sqrt(variance)

		if stddev > 0 {
			zScore := math.Abs((value - mean) / stddev)

			// 3-Sigma Anomaly Threshold
			if zScore > 3.0 {
				log.Printf("🚨 [ENTROPY ANOMALY DETECTED] Cell %s: Value %f deviates by %.2f standard deviations (Mean: %f, StdDev: %f)\n",
					cellHash, value, zScore, mean, stddev)

				// Here we would push a WebSockets CRDT comment forcing the user to justify the variance.
			}
		}
	}
}
