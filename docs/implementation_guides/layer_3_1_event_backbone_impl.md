# 📖 Layer 3.1 Implementation Guide: The Event Backbone (Nervous System)

**Status:** ✅ Implemented
**Location:**
*   **Schema (FlatBuffers):** `services/grid-service/schema/event/v1/atom_event.fbs`
*   **Interface (Go):** `services/grid-service/pkg/event/bus.go`
*   **Producer (Kafka/Redpanda):** `services/grid-service/pkg/event/kafka_producer.go`
**Technology:** Redpanda (Kafka Protocol) + FlatBuffers + Go

---

## 1. The "Why": The Speed of Thought
In a distributed system, latency kills the "Interactive Feel."
**The Problem:** Using HTTP/REST between services means serializing JSON, waiting for TCP, parsing JSON. Too slow (>50ms).
**The Solution:** We use **Redpanda** (C++ Kafka) with **FlatBuffers**.
*   **Redpanda:** Delivers events in <2ms.
*   **FlatBuffers:** The consumer (Rust Engine) can read the "Region" field from the message byte buffer *without* parsing the whole object.

## 2. The Implementation (Code Deep Dive)

### A. The Schema (`atom_event.fbs`)
We define `AtomEvent` as the standard envelope.
*   **`molecules` ([FlatMolecule]):** A vector of changed data.
*   **`schema_payload` (string):** For structural changes (added dimensions).
*   **`type` (Enum):** `MOLECULE_WRITE`, `CALC_REQUEST`, `HEARTBEAT`.

### B. The Producer (`kafka_producer.go`)
We use `segmentio/kafka-go` for high-throughput async writing.
*   **`BatchSize: 100`:** We don't send every keystroke individually. We buffer slightly (10ms) to group writes into packets.
*   **`Async: true`:** The API handler returns *immediately* after pushing to the local buffer. The user never waits for Redpanda.

## 3. Usage Instructions

### Initialization
Ensure Redpanda is running:
```bash
rpk container start
```

### Publishing an Event (Go)
```go
bus := event.NewKafkaBus([]string{"localhost:9092"}, "atom-events")
// Create event struct
evt := &event.AtomEventGo{
    TenantID: "tenant-1",
    Type:     event.TypeMoleculeWrite,
    Payload:  flatBufferBytes,
}
// Fire and forget (Async)
bus.Publish(ctx, evt)
```

## 4. Key Design Decisions related to "The Moat"
1.  **Zero-Copy Serialization:** By using FlatBuffers, the Rust Engine can map the incoming Kafka message directly into memory and read the `coordinate_hash` to find the cell in the Lattice. 
2.  **Partitioning Strategy:** We partition by `TenantID`. This ensures that all events for "Company A" land on the same Redpanda shard and are processed in strict order, guaranteeing causal consistency.
