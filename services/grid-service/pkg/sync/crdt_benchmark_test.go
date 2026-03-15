package sync

import (
	"crypto/sha256"
	"testing"
	"time"
)

func BenchmarkLWWElementSetMergeSequential(b *testing.B) {
	set := NewLWWElementSet()
	hash := sha256.Sum256([]byte("benchmark-coordinate"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		event := &CRDTEvent{
			EventID: "bench-seq",
			Clock: VectorClock{
				ClientID:   "tenant-a",
				Counter:    uint64(i + 1),
				WallTimeMs: time.Now().UnixMilli(),
			},
			Edit: CoordinateEdit{
				CoordinateHash: hash[:],
				NumericValue:   float64(i),
			},
		}

		_ = set.Merge(event)
	}
}

func BenchmarkLWWElementSetMergeConflictHeavy(b *testing.B) {
	set := NewLWWElementSet()
	hash := sha256.Sum256([]byte("benchmark-conflict-coordinate"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter := uint64((i / 2) + 1)
		event := &CRDTEvent{
			EventID: "bench-conflict",
			Clock: VectorClock{
				ClientID:   []string{"tenant-a", "tenant-b"}[i%2],
				Counter:    counter,
				WallTimeMs: time.Now().UnixMilli() + int64(i%2),
			},
			Edit: CoordinateEdit{
				CoordinateHash: hash[:],
				NumericValue:   float64(i),
				IsDelete:       i%5 == 0,
			},
		}

		_ = set.Merge(event)
	}
}
