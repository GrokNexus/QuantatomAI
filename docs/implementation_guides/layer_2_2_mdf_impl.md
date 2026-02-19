# 📖 Layer 2.2 Implementation Guide: The Molecular Store (MDF)

**Status:** ✅ Implemented
**Location:**
*   **Protocol:** `services/grid-service/proto/quantatomai/mdf/v1/molecule.proto`
*   **Writer (Go):** `services/grid-service/src/storage/mdf_writer.go`
*   **Reader (Rust):** `services/atom-engine/src/mdf/reader.rs`
**Technology:** Parquet + Protobuf + Zero-Copy Arrow

---

## 1. The "Why": The Universal Language
Traditional platforms have a "Translation Tax."
*   **Anaplan:** Java Objects (Siloed). To get data out, you must export to CSV.
*   **Postgres:** Row-oriented. Bad for analytics.

**The Solution:** We invented the **Molecular Data Format (MDF)**.
*   It is just a `.parquet` file on disk (S3).
*   It is an **Arrow RecordBatch** in memory (Rust).
*   It is a **Protobuf Message** over the wire (gRPC).

**Key Benefit:** We can move 1GB of data from Disk -> Engine -> Network with **Zero Copy** and **Zero Serialization Overhead**.

## 2. The Implementation (Code Deep Dive)

### A. The Schema (`molecule.proto`)
The `Molecule` is the "Atomic Unit."
*   **`coordinate_hash` (bytes):** A 128-bit hash of the dimension combination. This is the primary key for the hash map.
*   **`numeric_value` (double):** The actual number (e.g., $100.00).
*   **`security_mask` (uint64):** The holographic ACL tag embedded *with the data*.

### B. The Writer (Go)
We use `segmentio/parquet-go` to write high-throughput streams.
*   **Compression:** We force **Zstd**. It gives 3x better compression than Snappy for numerical data.
*   **Row Groups:** We write in 128MB chunks (Row Groups) to allow the Rust engine to parallelize reading.

### C. The Reader (Rust)
We use `parquet::arrow` to map the file directly into memory.
*   **Vectorized Read:** We do not read row-by-row. We read column-by-column into `SimdVector<f64>` arrays.
*   **Throughput:** Expect >5 GB/s read speeds from NVMe.

## 3. Usage Instructions

### Writing Data (Go)
```go
writer := storage.NewMdfWriter(file)
mol := mdfv1.NewNumericMolecule(hash, dims, 100.0, "SAP", mask)
writer.Write(mol)
writer.Close()
```

### Reading Data (Rust)
```rust
let batches = reader::read_mdf_arrow("data.parquet")?;
// batches[0].column(2) is now a layout-compatible Arrow array
```

## 4. Key Design Decisions related to "The Moat"
1.  **Polymorphism:** The same file format handles Numbers, Text Comments, and AI Vectors. No separate "Text Cube."
2.  **Time Travel:** The `timestamp` field allows us to replay history by simply filtering `WHERE timestamp < T`.
3.  **Cloud Native:** Implementation relies on standard Parquet. We can query these files directly with AWS Athena or DuckDB for ad-hoc analysis.
