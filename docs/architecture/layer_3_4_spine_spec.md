# QuantatomAI Layer 3 & 4 Specification: The Nervous System

## Layer 3: Resonance Eventing (The "Nerves")
This layer handles the high-speed propagation of writes, signals, and triggers.

### 3.1 The Log Backbone
*   **Technology:** **Redpanda** (C++ Kafka API).
*   **Performance:** 10x lower latency than JVM Kafka (no GC pauses).
*   **Topics:**
    *   `atom.write.raw`: Ingest stream (Protobuf Molecule).
    *   `atom.calc.trigger`: Signals to the Calc Engine.
    *   `audit.log`: Async stream to ClickHouse.

### 3.2 The Message Protocol
*   **Format:** **FlatBuffers v2**.
*   **Why:** Zero-parsing overhead. We can read a "Region ID" from the byte stream without decoding the whole message.
*   **Schema Registry:** **Buf** (Protobuf/RPC management).

### 3.3 The Stream Processors (Rust)
*   **Technology:** **Rust** + **Arroyo** / **Bytewax**.
*   **Role:**
    1.  **Enrichment:** Adds "Project Name" to a raw atom ID.
    2.  **Windowing:** Batches writes for 5ms before flushing to S3 (MDF).

---

## Layer 4: The Lattice Spine (The "Structure")
This layer provides the fast read-path and the "Shape" of the data for the engines.

### 4.1 The Warm Store (ClickHouse)
*   **Technology:** **ClickHouse** (MergeTree).
*   **Role:** Stores the "Flattened" version of the Lattice for aggregation.
*   **Moat Feature:** **Z-Order Curve Indexing**. We spatially index data by `(Time, Scenario, Account)` so fetching a "Quarterly Report" is a single sequential HDD read.

### 4.2 The IPC Layer (Arrow Flight)
*   **Technology:** **Apache Arrow Flight**.
*   **Function:** Delivers 100MB of data from ClickHouse to the Rust Engine (Layer 5) in **<50ms**.
*   **Mechanism:** Zero-Copy Shared Memory. The database writes to RAM; the engine reads the same RAM pointer.

### 4.3 The Metadata Spine (Postgres + Ltree)
*   **Technology:** **Postgres 16**.
*   **Optimizations:**
    *   **Materialized Paths:** Using `ltree` to store `Global.USA.East.NY` allows `WHERE path ~ 'Global.USA.*'` queries to execute in O(1) time.
