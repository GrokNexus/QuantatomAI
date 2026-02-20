package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CellLineageRecord represents a single state change for a specific grid coordinate.
type CellLineageRecord struct {
	EventID   uuid.UUID
	Timestamp time.Time
	UserID    uuid.UUID
	Action    EventType
	OldValue  float64
	NewValue  float64
	Comment   string
}

// LineageStore defines the interface for pulling audit history of a specific cell.
type LineageStore interface {
	// GetCellLineage retrieves the history of edits for a specific 128-bit coordinate hash.
	GetCellLineage(ctx context.Context, coordinateHash []byte, limit int) ([]*CellLineageRecord, error)
}

// ClickHouseLineageStore connects to the Entropy Ledger to retrieve column-oriented histories.
type ClickHouseLineageStore struct {
	// db *clickhouse.Conn
}

func NewClickHouseLineageStore() *ClickHouseLineageStore {
	return &ClickHouseLineageStore{}
}

// GetCellLineage queries ClickHouse for all WriteCell events targeting a specific coordinate.
func (s *ClickHouseLineageStore) GetCellLineage(ctx context.Context, coordinateHash []byte, limit int) ([]*CellLineageRecord, error) {
	// In production, this executes a ClickHouse `SELECT ... FROM audit_log WHERE hash = ? ORDER BY timestamp DESC`

	hashStr := fmt.Sprintf("%x", coordinateHash)
	fmt.Printf("[AUDIT_LINEAGE] Pulling cell lineage for hash: %s\n", hashStr)

	// Simulate returning historical data to prove the interface
	now := time.Now()
	history := []*CellLineageRecord{
		{
			EventID:   uuid.New(),
			Timestamp: now,
			UserID:    uuid.New(),
			Action:    TypeWriteCell,
			OldValue:  500.0,
			NewValue:  650.0,
			Comment:   "Updated Q3 forecast based on new guidance.",
		},
		{
			EventID:   uuid.New(),
			Timestamp: now.Add(-2 * time.Hour),
			UserID:    uuid.New(),
			Action:    TypeWriteCell,
			OldValue:  0.0,
			NewValue:  500.0,
			Comment:   "Initial data load.",
		},
	}

	if limit > 0 && len(history) > limit {
		return history[:limit], nil
	}
	return history, nil
}
