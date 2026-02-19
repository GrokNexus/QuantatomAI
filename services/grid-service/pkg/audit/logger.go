package audit

import (
	"context"
	"fmt"
	"sync" // Ultra Diamond: Graceful Shutdown
	"time"

	"github.com/google/uuid"
)

// EventType categorizes the audit action.
type EventType string

const (
	TypeLogin      EventType = "LOGIN"
	TypeWriteCell  EventType = "WRITE_CELL"
	TypeExecute    EventType = "EXECUTE_CALC"
	TypeSchemaEdit EventType = "UPDATE_METADATA"
)

// AuditEvent represents a single immutable log entry.
type AuditEvent struct {
	EventID   uuid.UUID
	Timestamp time.Time
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Action    EventType
	Details   string // JSON payload
	Meta      map[string]string
}

// Logger is the interface for recording audit events.
type Logger interface {
	Log(ctx context.Context, tenantID, userID uuid.UUID, action EventType, details string)
	Close() error
}

// AsyncClickHouseLogger buffers events and flushes them to ClickHouse.
type AsyncClickHouseLogger struct {
	eventCh chan *AuditEvent
	doneCh  chan struct{}
	wg      sync.WaitGroup // Ultra Diamond: Wait for worker to finish
}

// NewAsyncLogger creates a logger with a buffered channel (e.g., size 10,000).
func NewAsyncLogger() *AsyncClickHouseLogger {
	l := &AsyncClickHouseLogger{
		eventCh: make(chan *AuditEvent, 10000),
		doneCh:  make(chan struct{}),
	}
	l.wg.Add(1) // Track the worker
	go l.worker()
	return l
}

// Log pushes an event to the channel. It is non-blocking unless the buffer is full.
func (l *AsyncClickHouseLogger) Log(ctx context.Context, tenantID, userID uuid.UUID, action EventType, details string) {
	event := &AuditEvent{
		EventID:   uuid.New(),
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	select {
	case l.eventCh <- event:
		// Success
	default:
		// Buffer full: Drop or log to stderr to avoid blocking main thread
		fmt.Printf("AUDIT_DROP: Buffer full. Dropped event %s\n", event.EventID)
	}
}

// worker consumes events and writes them in batches (Simulated for Layer 2.3).
func (l *AsyncClickHouseLogger) worker() {
	defer l.wg.Done() // Signal completion

	batch := make([]*AuditEvent, 0, 100)
	ticker := time.NewTicker(time.Second * 1) // Flush every second
	defer ticker.Stop()

	for {
		select {
		case event := <-l.eventCh:
			batch = append(batch, event)
			if len(batch) >= 100 {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-l.doneCh:
			// Drain remaining events
			if len(batch) > 0 {
				l.flush(batch)
			}
			// Optional: drain channel if needed, but for now we flush batch and exit
			return
		}
	}
}

// flush simulates writing to ClickHouse. In production, this uses clickhouse-go.
func (l *AsyncClickHouseLogger) flush(events []*AuditEvent) {
	// TODO: Replace with actual ClickHouse batch insert
	// For now, we just print to stdout to verify hook functionality
	fmt.Printf("[AUDIT_FLUSH] Writing %d events to Entropy Ledger (ClickHouse). First ID: %s\n", len(events), events[0].EventID)
}

func (l *AsyncClickHouseLogger) Close() error {
	close(l.doneCh)
	l.wg.Wait() // Block until worker flushes pending logs
	return nil
}

// AsyncClickHouseLogger buffers events and flushes them to ClickHouse.
type AsyncClickHouseLogger struct {
	eventCh chan *AuditEvent
	doneCh  chan struct{}
}

// NewAsyncLogger creates a logger with a buffered channel (e.g., size 10,000).
func NewAsyncLogger() *AsyncClickHouseLogger {
	l := &AsyncClickHouseLogger{
		eventCh: make(chan *AuditEvent, 10000),
		doneCh:  make(chan struct{}),
	}
	go l.worker()
	return l
}

// Log pushes an event to the channel. It is non-blocking unless the buffer is full.
func (l *AsyncClickHouseLogger) Log(ctx context.Context, tenantID, userID uuid.UUID, action EventType, details string) {
	event := &AuditEvent{
		EventID:   uuid.New(),
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	select {
	case l.eventCh <- event:
		// Success
	default:
		// Buffer full: Drop or log to stderr to avoid blocking main thread
		fmt.Printf("AUDIT_DROP: Buffer full. Dropped event %s\n", event.EventID)
	}
}

// worker consumes events and writes them in batches (Simulated for Layer 2.3).
func (l *AsyncClickHouseLogger) worker() {
	batch := make([]*AuditEvent, 0, 100)
	ticker := time.NewTicker(time.Second * 1) // Flush every second
	defer ticker.Stop()

	for {
		select {
		case event := <-l.eventCh:
			batch = append(batch, event)
			if len(batch) >= 100 {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				l.flush(batch)
				batch = batch[:0]
			}
		case <-l.doneCh:
			if len(batch) > 0 {
				l.flush(batch)
			}
			return
		}
	}
}

// flush simulates writing to ClickHouse. In production, this uses clickhouse-go.
func (l *AsyncClickHouseLogger) flush(events []*AuditEvent) {
	// TODO: Replace with actual ClickHouse batch insert
	// For now, we just print to stdout to verify hook functionality
	fmt.Printf("[AUDIT_FLUSH] Writing %d events to Entropy Ledger (ClickHouse). First ID: %s\n", len(events), events[0].EventID)
}

func (l *AsyncClickHouseLogger) Close() error {
	close(l.doneCh)
	return nil
}
