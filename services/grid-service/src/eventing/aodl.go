package eventing

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "quantatomai/grid-service/domain"
)

// AODLClient defines the contract for emitting atom-level events into the Atom-Oriented Data Log.
type AODLClient interface {
    EmitCellUpdated(ctx context.Context, atom domain.AtomWrite, correlationID string) error
    EmitCellUpdatedBatch(ctx context.Context, atoms []domain.AtomWrite, correlationID string) error
}

// KinesisLikeClient abstracts the streaming client (Kinesis, Kafka, Pub/Sub, etc.).
type KinesisLikeClient interface {
    PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) error
    PutRecords(ctx context.Context, streamName string, records []StreamRecord) error
}

// StreamRecord is a generic batch record for any streaming backend.
type StreamRecord struct {
    PartitionKey string
    Data         []byte
}

// KinesisAODLClient is a concrete implementation that writes CellUpdated events into an AODL stream.
type KinesisAODLClient struct {
    client     KinesisLikeClient
    streamName string
}

// NewKinesisAODLClient constructs a new AODL client.
func NewKinesisAODLClient(client KinesisLikeClient, streamName string) *KinesisAODLClient {
    return &KinesisAODLClient{
        client:     client,
        streamName: streamName,
    }
}

// cellUpdatedEvent is the canonical AODL event payload for a single cell update.
type cellUpdatedEvent struct {
    Type          string        `json:"type"`          // "CellUpdated"
    AtomKey       atomKeyWire   `json:"atomKey"`
    Value         float64       `json:"value"`
    User          string        `json:"user"`
    Timestamp     time.Time     `json:"timestamp"`
    Source        string        `json:"source"`        // e.g., "grid-service"
    Version       int           `json:"version"`       // schema version
    CorrelationID string        `json:"correlationId"` // traceability across layers
}

type atomKeyWire struct {
    DimIDs     []int64 `json:"dimIds"`
    MeasureID  int64   `json:"measureId"`
    ScenarioID int64   `json:"scenarioId"`
}

// EmitCellUpdated serializes and emits a single CellUpdated event into the AODL stream.
func (c *KinesisAODLClient) EmitCellUpdated(ctx context.Context, atom domain.AtomWrite, correlationID string) error {
    atom.Key.EnsureCanonical()

    evt := cellUpdatedEvent{
        Type: "CellUpdated",
        AtomKey: atomKeyWire{
            DimIDs:     atom.Key.DimIDs,
            MeasureID:  atom.Key.MeasureID,
            ScenarioID: atom.Key.ScenarioID,
        },
        Value:         atom.Value,
        User:          atom.User,
        Timestamp:     time.Now().UTC(),
        Source:        "grid-service",
        Version:       1,
        CorrelationID: correlationID,
    }

    payload, err := json.Marshal(evt)
    if err != nil {
        return fmt.Errorf("failed to marshal CellUpdated event: %w", err)
    }

    partitionKey := fmt.Sprintf("%d", atom.Key.HashKey())

    if err := c.client.PutRecord(ctx, c.streamName, partitionKey, payload); err != nil {
        return fmt.Errorf("failed to emit CellUpdated event: %w", err)
    }

    return nil
}

// EmitCellUpdatedBatch emits a batch of CellUpdated events.
// This is critical for spreads, allocations, and range edits.
func (c *KinesisAODLClient) EmitCellUpdatedBatch(ctx context.Context, atoms []domain.AtomWrite, correlationID string) error {
    if len(atoms) == 0 {
        return nil
    }

    records := make([]StreamRecord, 0, len(atoms))

    for _, atom := range atoms {
        atom.Key.EnsureCanonical()

        evt := cellUpdatedEvent{
            Type: "CellUpdated",
            AtomKey: atomKeyWire{
                DimIDs:     atom.Key.DimIDs,
                MeasureID:  atom.Key.MeasureID,
                ScenarioID: atom.Key.ScenarioID,
            },
            Value:         atom.Value,
            User:          atom.User,
            Timestamp:     time.Now().UTC(),
            Source:        "grid-service",
            Version:       1,
            CorrelationID: correlationID,
        }

        payload, err := json.Marshal(evt)
        if err != nil {
            return fmt.Errorf("failed to marshal batch CellUpdated event: %w", err)
        }

        partitionKey := fmt.Sprintf("%d", atom.Key.HashKey())

        records = append(records, StreamRecord{
            PartitionKey: partitionKey,
            Data:         payload,
        })
    }

    if err := c.client.PutRecords(ctx, c.streamName, records); err != nil {
        return fmt.Errorf("failed to emit batch CellUpdated events: %w", err)
    }

    return nil
}
