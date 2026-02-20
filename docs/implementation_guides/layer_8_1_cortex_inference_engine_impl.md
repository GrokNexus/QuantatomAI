ommit# QuantatomAI Layer 8.1: The Intelligence Cortex (Inference Engine) Implementation Guide

## 1. Executive Summary
Layer 8 introduces the "Cortex," the probabilistic "Right Brain" of the QuantatomAI architecture. While the deterministic Go/Rust engines (Layers 1-7) handle precision mathematics and grid rendering, the Cortex is responsible for autonomous FP&A capabilities: Auto-Baselining, Generative Scenarios, and Anomaly Detection.

Layer 8.1 establishes the foundational Inference Engine capable of executing these AI workloads.

## 2. Architecture & Design Decisions
To ensure the Cortex does not degrade the "Diamond-Tier" latency of the Grid Service, we implemented it as an independent microservice.

*   **Runtime:** Python 3.12 with **FastAPI**. Python is the undisputed leader in the AI/ML ecosystem (PyTorch, Langchain, vLLM). FastAPI provides async, high-throughput HTTP endpoints that match Go's performance profile well enough for AI workloads where inference is the bottleneck.
*   **Data Link:** The most critical innovation is the **Zero-Copy MDF Ingestion**.
    *   *Anti-Pattern:* Querying the Go Grid Service via REST, forcing the engine to serialize 10M rows to JSON, only for Python to parse it back into Pandas.
    *   *QuantatomAI Pattern:* We implemented `MDFVectorReader` using `pyarrow`. The Cortex reads the Molecular Data Format (Parquet) files directly from the underlying object storage (S3/MinIO), entirely bypassing the transactional Grid Engine.

## 3. Implementation Details

### 3.1 The FastAPI Service (`src/main.py`)
Provides the RESTful API contract for the UI and the Go Orchestrator to trigger AI workloads.
*   **Endpoints Established:**
    *   `/health`: Standard K8s probe.
    *   `/api/v1/forecast/auto-baseline`: Stub for Layer 8.2 (Transformer forecasting).
    *   `/api/v1/narrative/variance`: Stub for Layer 8.3 (Generative explainability).

### 3.2 The MDF Reader (`src/data/vector_reader.py`)
*   `MDFVectorReader` class utilizes `pyarrow.parquet.read_table` for native integration with the Data Spine.
*   Includes fallback simulation logic to allow independent development of the UI while the Rust pipeline is finalizing.

### 3.3 Infrastructure Deployment (`infra/k8s/base/cortex-service/`)
*   **Containerization:** A multi-stage `Dockerfile` ensures a lightweight production image by building Python wheels in an isolated stage.
*   **Kubernetes Manifest:** `deployment.yaml` with explicit CPU/Memory limits designed for memory-intensive inference operations.
*   **Integration:** Registered in `kustomization.yaml` to deploy alongside `grid-service` and `web-ui`.

## 4. Operational Readiness
The Cortex is now "alive" and ready to host specialized machine learning models (e.g., TimeGPT, Lag-Llama, Isolation Forests).
