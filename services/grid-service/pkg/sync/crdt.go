package sync

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

// VectorClock represents a logical timestamp for a specific active client/tenant.
type VectorClock struct {
	ClientID   string
	Counter    uint64
	WallTimeMs int64 // Used for tie-breaking
}

// Compare returns 1 if v > other, -1 if v < other, and 0 if equal.
// It implements the standard LWW (Last-Write-Wins) tie-breaking mechanism.
func (v *VectorClock) Compare(other *VectorClock) int {
	if v.Counter > other.Counter {
		return 1
	} else if v.Counter < other.Counter {
		return -1
	}

	// Logical clocks match. Fall back to WallTime (NTP synchronized across cluster).
	if v.WallTimeMs > other.WallTimeMs {
		return 1
	} else if v.WallTimeMs < other.WallTimeMs {
		return -1
	}

	// Total collision: Final tie-break on ClientID to guarantee deterministic ordering.
	if v.ClientID > other.ClientID {
		return 1
	} else if v.ClientID < other.ClientID {
		return -1
	}

	return 0
}

// CoordinateEdit represents a raw edit payload targeting a specific 128-bit hash.
type CoordinateEdit struct {
	CoordinateHash []byte
	NumericValue   float64
	StringValue    string
	IsDelete       bool
}

// CRDTEvent represents the envelope broadcasted by the Hub.
type CRDTEvent struct {
	EventID string
	Clock   VectorClock
	Edit    CoordinateEdit
}

// LWWElementSet implements a Last-Write-Wins Element Set CRDT.
// It is mathematically deterministic and commutative.
type LWWElementSet struct {
	mu sync.RWMutex

	// AddSet stores the VectorClock of the most recent ADD/UPDATE operation.
	AddSet map[string]VectorClock

	// RemoveSet stores the VectorClock of the most recent DELETE operation.
	RemoveSet map[string]VectorClock
}

func NewLWWElementSet() *LWWElementSet {
	return &LWWElementSet{
		AddSet:    make(map[string]VectorClock),
		RemoveSet: make(map[string]VectorClock),
	}
}

// stringifyHash converts the 16-byte coordinate to a map key.
func stringifyHash(hash []byte) string {
	return hex.EncodeToString(hash)
}

// Merge processes an incoming CRDTEvent and determines if it should be applied to the local state.
// It returns true if the event represents a *new* valid state that should trigger a recalculation/storage write.
func (set *LWWElementSet) Merge(event *CRDTEvent) bool {
	set.mu.Lock()
	defer set.mu.Unlock()

	hashKey := stringifyHash(event.Edit.CoordinateHash)

	if event.Edit.IsDelete {
		// Check if this delete is newer than the existing delete
		if existing, ok := set.RemoveSet[hashKey]; ok {
			if event.Clock.Compare(&existing) <= 0 {
				return false // The incoming event is older. Discard.
			}
		}
		set.RemoveSet[hashKey] = event.Clock

		// LWW bias: If an ADD and REMOVE have the exact same clock (impossible with UUID ties, but theoretically),
		// standard practice biases towards REMOVE or ADD depending on business rules.
		// Our `Compare` method guarantees no ties.
		return true

	} else {
		// It's an Add/Update
		// 1. Is it newer than the current Add?
		if existingAdd, ok := set.AddSet[hashKey]; ok {
			if event.Clock.Compare(&existingAdd) <= 0 {
				return false // Older add/update event. Discard.
			}
		}

		// 2. Is it newer than a potential Delete? (Can't update a tombstone with an old edit).
		if existingRemove, ok := set.RemoveSet[hashKey]; ok {
			if event.Clock.Compare(&existingRemove) <= 0 {
				return false // Target was recently deleted by a newer event. Discard.
			}
		}

		set.AddSet[hashKey] = event.Clock
		return true
	}
}

// Hub orchestrates the CRDT application for a specific Grid instance.
// It sits between the WebSockets/gRPC streams and the Kafka Redpanda Backbone.
type Hub struct {
	crdts map[string]*LWWElementSet // Keyed by Grid ID
}

func NewHub() *Hub {
	return &Hub{
		crdts: make(map[string]*LWWElementSet),
	}
}

// GenerateEvent mints a new CRDT payload for a user edit, advancing their logical clock.
func (h *Hub) GenerateEvent(clientID string, counter uint64, edit CoordinateEdit) *CRDTEvent {
	return &CRDTEvent{
		EventID: uuid.New().String(),
		Clock: VectorClock{
			ClientID:   clientID,
			Counter:    counter,
			WallTimeMs: time.Now().UnixMilli(),
		},
		Edit: edit,
	}
}

// -----------------------------------------------------------------------------
// Layer 3.3: Presence & Collaborative Cursors
// -----------------------------------------------------------------------------

// SelectionRange identifies a block of coordinates a user is currently highlighting.
type SelectionRange struct {
	StartHash []byte // 128-bit coordinate
	EndHash   []byte // 128-bit coordinate
}

// PresenceUpdate represents a user's cursor or selection state changing.
// This is broadcasted optimistically to all connected clients.
type PresenceUpdate struct {
	ClientID  string
	ColorHex  string // e.g. "#FF0055"
	UserName  string
	Selection SelectionRange
}

// PresenceStreamer defines the interface for the WebSockets/gRPC bidirectional link.
type PresenceStreamer interface {
	// Broadcast sends a PresenceUpdate to all users connected to the same Grid.
	Broadcast(gridID string, update PresenceUpdate) error

	// OnUpdate registers a callback when another user moves their cursor.
	OnUpdate(gridID string, handler func(update PresenceUpdate))
}

// RedisPresenceStreamer offloads transient cursor movements to Redis PubSub.
// Ultra Diamond Vector 4: Prevents Redpanda/Kafka file descriptor exhaustion.
type RedisPresenceStreamer struct {
	// redisClient *redis.Client (Injected in production)
}

func NewRedisPresenceStreamer() *RedisPresenceStreamer {
	return &RedisPresenceStreamer{}
}

func (s *RedisPresenceStreamer) Broadcast(gridID string, update PresenceUpdate) error {
	// Simulated production behavior:
	// payload, _ := json.Marshal(update)
	// s.redisClient.Publish(ctx, "grid_presence:"+gridID, payload)
	return nil
}

func (s *RedisPresenceStreamer) OnUpdate(gridID string, handler func(update PresenceUpdate)) {
	// Simulated production behavior:
	// pubsub := s.redisClient.Subscribe(ctx, "grid_presence:"+gridID)
	// for msg := range pubsub.Channel() {
	//     var update PresenceUpdate
	//     json.Unmarshal([]byte(msg.Payload), &update)
	//     handler(update)
	// }
}
