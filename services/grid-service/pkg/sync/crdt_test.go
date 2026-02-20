package sync

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVectorClock_Compare(t *testing.T) {
	clock1 := VectorClock{ClientID: "A", Counter: 1, WallTimeMs: 1000}
	clock2 := VectorClock{ClientID: "B", Counter: 2, WallTimeMs: 1000}
	clock3 := VectorClock{ClientID: "A", Counter: 1, WallTimeMs: 1001}

	// Tie breaker identical except specific ClientID
	clockTie1 := VectorClock{ClientID: "A", Counter: 1, WallTimeMs: 1000}
	clockTie2 := VectorClock{ClientID: "B", Counter: 1, WallTimeMs: 1000}

	assert.Equal(t, -1, clock1.Compare(&clock2), "clock1 < clock2 due to counter")
	assert.Equal(t, 1, clock2.Compare(&clock1), "clock2 > clock1 due to counter")

	assert.Equal(t, -1, clock1.Compare(&clock3), "clock1 < clock3 due to walltime")
	assert.Equal(t, 1, clock3.Compare(&clock1), "clock3 > clock1 due to walltime")

	assert.Equal(t, -1, clockTie1.Compare(&clockTie2), "clockTie1 < clockTie2 due to ClientID tie-break")
	assert.Equal(t, 1, clockTie2.Compare(&clockTie1), "clockTie2 > clockTie1 due to ClientID tie-break")
	assert.Equal(t, 0, clock1.Compare(&clock1), "clock1 == clock1")
}

func TestLWWElementSet_Merge(t *testing.T) {
	crdt := NewLWWElementSet()

	hash := sha256.Sum256([]byte("coord_A"))

	// 1. Initial Edit from Client A
	editA1 := CRDTEvent{
		EventID: "evt1",
		Clock:   VectorClock{ClientID: "ClientA", Counter: 1, WallTimeMs: time.Now().UnixMilli()},
		Edit:    CoordinateEdit{CoordinateHash: hash[:], NumericValue: 100.0, IsDelete: false},
	}

	accepted := crdt.Merge(&editA1)
	assert.True(t, accepted, "Initial edit should be accepted")
	assert.Equal(t, crdt.AddSet[stringifyHash(hash[:])].Counter, uint64(1))

	// 2. Conflicting Older Edit from Client B (Delayed by network)
	editB0 := CRDTEvent{
		EventID: "evt2",
		Clock:   VectorClock{ClientID: "ClientB", Counter: 0, WallTimeMs: time.Now().UnixMilli() - 1000},
		Edit:    CoordinateEdit{CoordinateHash: hash[:], NumericValue: 50.0, IsDelete: false},
	}

	acceptedDelay := crdt.Merge(&editB0)
	assert.False(t, acceptedDelay, "Older edit should be discarded by LWW")
	assert.Equal(t, crdt.AddSet[stringifyHash(hash[:])].Counter, uint64(1), "State should not move backward")

	// 3. Delete Operation from Client C
	editC2 := CRDTEvent{
		EventID: "evt3",
		Clock:   VectorClock{ClientID: "ClientC", Counter: 2, WallTimeMs: time.Now().UnixMilli()},
		Edit:    CoordinateEdit{CoordinateHash: hash[:], NumericValue: 0.0, IsDelete: true},
	}

	acceptedDelete := crdt.Merge(&editC2)
	assert.True(t, acceptedDelete, "Delete with newer clock should be accepted")

	// 4. Update Operation from A attempting to revive deleted cell (Older clock)
	editA1_Revive := CRDTEvent{
		EventID: "evt4",
		Clock:   VectorClock{ClientID: "ClientA", Counter: 1, WallTimeMs: time.Now().UnixMilli() + 5000}, // Even with newer time, Counter 1 is older than Counter 2
		Edit:    CoordinateEdit{CoordinateHash: hash[:], NumericValue: 200.0, IsDelete: false},
	}

	acceptedRevive := crdt.Merge(&editA1_Revive)
	assert.False(t, acceptedRevive, "Cannot update a tombstone with an older vector clock")
}

func TestPresenceUpdate_StructTypes(t *testing.T) {
	// Simple struct initialization test to ensure fields are valid and exported
	p := PresenceUpdate{
		ClientID: "User123",
		ColorHex: "#ff0000",
		UserName: "John Doe",
		Selection: SelectionRange{
			StartHash: []byte{1, 2, 3},
			EndHash:   []byte{4, 5, 6},
		},
	}

	assert.Equal(t, "User123", p.ClientID)
	assert.Equal(t, "#ff0000", p.ColorHex)
	assert.Equal(t, "John Doe", p.UserName)
}
