# 📖 Layer 3.2 Implementation Guide: The IPC Layer (Arrow Flight)

**Status:** ✅ Implemented
**Location:**
*   **Server (Rust):** `services/atom-engine/src/ipc/service.rs`
*   **Client (Go):** `services/grid-service/pkg/ipc/flight_client.go`
**Technology:** Apache Arrow Flight (gRPC + FlatBuffers)

---

## 1. The "Why": Solving the Serialization Bottleneck
In typical microservices, moving 1GB of data means:
1.  Source serializes to JSON (CPU heavy).
2.  Network transfer.
3.  Destination parses JSON (CPU heavy).
**The Problem:** This adds seconds of latency for large grids.
**The Solution:** We use **Arrow Flight**.
*   It sends Arrow RecordBatches over gRPC.
*   **Zero Serialization:** The memory layout of the data on the wire is effectively the same as in RAM for Rust/Pandas/NumPy.
*   **Throughput:** Capable of saturating a 10Gbps link.

## 2. The Implementation (Code Deep Dive)

### A. The Server (Rust Kernel)
We implement the `FlightService` trait.
*   **`do_get`:** The main endpoint. Accepts a `Ticket` (containing a Plan ID or Query).
*   **Action:** It executes the calculation plan using the `ProjectionEngine` and streams the results back as a sequence of RecordBatches.
*   **Concurrency:** Handled via `tonic`'s async runtime (Tokio).

### B. The Client (Go Orchestrator)
We use the official Apache Arrow Go client.
*   **`GetCalculation(planID)`:** Sends a lightweight request.
*   **`flight.Reader`:** Returns a stream reader. We can iterate over batches as they arrive.
*   **Fan-Out:** The Go service can parallelize requests to multiple Rust workers if we shard the grid.

## 3. Usage Instructions

### Starting the Rust Server
The `atom-engine` binary will start the gRPC server on port `50051`.

### Fetching Data (Go)
```go
client, _ := ipc.NewFlightClient("localhost:50051")
reader, _ := client.GetCalculation(ctx, "PLAN-123")

for reader.Next() {
    record := reader.Record()
    // record is an Arrow Record (Columnar)
    // Send directly to UI via Connect-Web or process
}
```

## 4. Key Design Decisions related to "The Moat"
1.  **Shared Memory (Future Proof):** Arrow Flight supports `DoPut` and `DoExchange`. In the future, we can allow the Go service to "push" data updates directly into the Rust engine's memory arena without disk I/O.
2.  **Language Agnostic:** If we ever need a Python Data Science worker (e.g., for PyTorch), it can consume the same Flight stream with zero changes.
