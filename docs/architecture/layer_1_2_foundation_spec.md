# QuantatomAI Layer 1 & 2 Specification: The Bedrock

## Layer 1: The Cloud Infrastructure (The "Ground")
This layer provides the raw compute, networking, and security substrate.

### 1.1 Kubernetes & Orchestration
*   **Engine:** **Amazon EKS / Azure AKS (v1.29+)** or Bare Metal K8s.
*   **Distribution:** **Talos Linux** (Immutable OS) for zero-overhead security.
*   **Networking (CNI):** **Cilium (eBPF)**. Replaces kube-proxy for 10x faster networking and transparent encryption.
*   **Service Mesh:** **Istio Ambient Mesh**. Sidecar-less architecture to reduce RAM usage by 40%.

### 1.2 Infrastructure as Code (IaC)
*   **Tool:** **Crossplane** (Go). We manage AWS RDS and S3 buckets as Kubernetes Custom Resources (CRDs). No Terraform state files to lose.
*   **Secrets:** **External Secrets Operator** syncing with AWS Secrets Manager / HashiCorp Vault.

### 1.3 Observability (The "Eyes")
*   **Metrics:** **VictoriaMetrics** (High-performance Prometheus alternative).
*   **Logs:** **Grafana Loki** (Log aggregation).
*   **Traces:** **honeycomb.io** or **Jaeger** (OpenTelemetry native).

---

## Layer 2: The Data Foundation (The "Vault")
This layer manages the persistence, transactional integrity, and atomic storage.

### 2.1 The Molecular Store (MDF)
*   **Technology:** **Object Storage (S3 / MinIO)** + **Apache Parquet**.
*   **Format:** The "MDF" binary format defined in `quantatomai_data_structure_innovation.md`.
*   **Strategy:** Data is immutable. We write new Parquet files (Delta Lake style) rather than updating rows.

### 2.2 The Metadata Registry (Postgres)
*   **Technology:** **Postgres 16 (Citus Extension)** on AWS Aurora / Crunchy Data.
*   **Role:** Stores Dimensions, Hierarchies, Users, and Permissions. **NOT** the cell data.
*   **Extensions:**
    *   `pgvector`: For embedding search.
    *   `ltree`: For high-speed hierarchy traversal.

### 2.3 The High-Velocity Cache (Redis)
*   **Technology:** **DragonflyDB** (Redis-compatible, stored on SSDs/NVMe).
*   **Role:** Stores the "Hot Tiles" (Active User Grids) and Session State.
*   **Pattern:** Write-Through Cache for the Calculation Engine.

### 2.4 The Entropy Ledger (Audit)
*   **Technology:** **ClickHouse**.
*   **Role:** Append-only log of every single formula change and cell update for SOX compliance.
*   **Retention:** Infinite.
