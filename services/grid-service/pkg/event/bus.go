package event

import (
	"context"
)

// EventType matches the FlatBuffers enum.
type EventType byte

const (
	TypeUnknown      EventType = 0
	TypeMoleculeWrite EventType = 1
	TypeCalcRequest   EventType = 2
	TypeSchemaChange  EventType = 3
	TypeHeartbeat     EventType = 4
)

// AtomEventGo is the high-level Go struct for an event.
// We use this before serializing to FlatBuffers.
type AtomEventGo struct {
	EventID   string
	TraceID   string
	Timestamp int64
	Type      EventType
	TenantID  string
	UserID    string
	Payload   []byte // The serialized FlatBuffer payload (molecules or schema)
}

// Bus is the interface for the Nervous System.
// It abstracts the underlying Redpanda/Kafka implementation.
type Bus interface {
	// Publish sends an event to the "atom-events" topic.
	// It should return nil if the event is successfully buffered (async).
	Publish(ctx context.Context, event *AtomEventGo) error

	// Subscribe listens for events.
	// In a real implementation, this would return a channel or take a callback.
	Subscribe(ctx context.Context, topic string) (<-chan *AtomEventGo, error)

	// Close flushes buffers and closes connections.
	Close() error
}
