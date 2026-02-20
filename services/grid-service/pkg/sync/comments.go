package sync

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Comment represents a single message in a cell-level threaded discussion.
type Comment struct {
	CommentID  string    `json:"comment_id"`
	AuthorID   string    `json:"author_id"`
	Author     string    `json:"author"`
	Text       string    `json:"text"`
	Timestamp  time.Time `json:"timestamp"`
	IsResolved bool      `json:"is_resolved"`
}

// Thread represents a collection of comments pinned to a specific Grid Coordinate.
type Thread struct {
	ThreadID       string    `json:"thread_id"`
	CoordinateHash string    `json:"coordinate_hash"` // Hex encoded 128-bit coordinate
	GridID         string    `json:"grid_id"`
	Comments       []Comment `json:"comments"`
	Status         string    `json:"status"` // "Open", "Resolved"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ThreadStore defines the interface for persisting and retrieving cell comments.
// In production, this would be backed by Postgres (Metadata layer) or MongoDB.
type ThreadStore interface {
	// CreateThread initializes a new discussion thread on a specific cell.
	CreateThread(gridID string, coordinateHash []byte) (*Thread, error)

	// AddComment appends a new message to an existing thread.
	AddComment(threadID string, comment Comment) error

	// GetThreadForCell retrieves the active discussion for a specific cell, if any.
	GetThreadForCell(gridID string, coordinateHash []byte) (*Thread, error)

	// ResolveThread marks the entire discussion as resolved.
	ResolveThread(threadID string) error
}

// InMemoryThreadStore is a transient implementation for the Go Orchestrator.
type InMemoryThreadStore struct {
	mu      sync.RWMutex
	threads map[string]*Thread // Keyed by ThreadID

	// Index to quickly find a thread by its cell coordinate
	cellIndex map[string]string // Key: "GridID:CoordHash", Value: ThreadID
}

func NewInMemoryThreadStore() *InMemoryThreadStore {
	return &InMemoryThreadStore{
		threads:   make(map[string]*Thread),
		cellIndex: make(map[string]string),
	}
}

// stringifyHash converts the 16-byte coordinate to a map key. (Reused from crdt.go logic)
func encodeCoordMapKey(gridID string, hash []byte) string {
	return gridID + ":" + hex.EncodeToString(hash)
}

func (s *InMemoryThreadStore) CreateThread(gridID string, coordinateHash []byte) (*Thread, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := encodeCoordMapKey(gridID, coordinateHash)
	if existingID, exists := s.cellIndex[key]; exists {
		return s.threads[existingID], nil // Thread already exists
	}

	thread := &Thread{
		ThreadID:       uuid.New().String(),
		CoordinateHash: hex.EncodeToString(coordinateHash),
		GridID:         gridID,
		Comments:       []Comment{},
		Status:         "Open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	s.threads[thread.ThreadID] = thread
	s.cellIndex[key] = thread.ThreadID

	return thread, nil
}

func (s *InMemoryThreadStore) AddComment(threadID string, comment Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	thread, exists := s.threads[threadID]
	if !exists {
		return fmt.Errorf("thread not found")
	}

	comment.CommentID = uuid.New().String()
	if comment.Timestamp.IsZero() {
		comment.Timestamp = time.Now()
	}

	thread.Comments = append(thread.Comments, comment)
	thread.UpdatedAt = time.Now()

	return nil
}

func (s *InMemoryThreadStore) GetThreadForCell(gridID string, coordinateHash []byte) (*Thread, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := encodeCoordMapKey(gridID, coordinateHash)
	threadID, exists := s.cellIndex[key]
	if !exists {
		return nil, nil // No active thread for this cell
	}

	return s.threads[threadID], nil
}

func (s *InMemoryThreadStore) ResolveThread(threadID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	thread, exists := s.threads[threadID]
	if !exists {
		return fmt.Errorf("thread not found")
	}

	thread.Status = "Resolved"
	thread.UpdatedAt = time.Now()

	// Mark all comments as resolved
	for i := range thread.Comments {
		thread.Comments[i].IsResolved = true
	}

	return nil
}
