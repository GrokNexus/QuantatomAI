# 🧠 QuantatomAI Cortex: The Intelligence Layer (Layer 8)

You asked: *"How are AI capabilities included? Go Red Team and innovate."*

**The Brutal Truth:**
Our current 7-Layer Architecture is a **Deterministic Compute Engine**. It is brilliant at calculating *what* happened (Aggregation) and *what if* (Formula Logic), but it lacks the **Probabilistic Intelligence** to tell you *why* it happened or *what might* happen.

To solve this, we introduce **Layer 8: The QuantatomAI Cortex**.

## 1. The Core Innovation: "The Dual-Brain Architecture"
*   **Left Brain (AtomEngine - Rust):** Deterministic, precision accounting, 100% accuracy. Handles the "Grid."
*   **Right Brain (Cortex - Python/Mojo):** Probabilistic, creative, pattern-matching. Handles the "Insight."

### The Protocol: "Vector Resonance"
Since we switched to the **Molecular Data Format (MDF)**, every data point already has an embedding slot (`embedding_vector`). The Cortex reads these vectors directly from S3/ClickHouse.

## 2. Capabilities: The "AI Moat"

### A. Auto-Baseline (The "Zero-Draft" Forecast)
*   **Problem:** Planning starts with a blank sheet. Users hate manual entry.
*   **Cortex Solution:**
    *   **Mechanism:** When a user opens a "Budget 2026" version, the Cortex automatically runs a **Transformer-based Time-Series Model (e.g., TimeGPT or Lag-Llama)** on the 5-year history in MDF.
    *   **Result:** The grid pre-populates with a high-confidence baseline. The user just "adjusts by exception."
    *   **Tech:** Python Service (FastAPI) + PyTorch + Arrow Flight (for data ingest).

### B. Generative Scenarios ("Commander Mode")
*   **Problem:** Creating a scenario requires clicking "New Version," copying data, and applying % uplifts.
*   **Cortex Solution:**
    *   **Interface:** A Chat Sidebar (LLM).
    *   **Prompt:** *"Create a 'High Inflation' scenario where COGS increases by 8% and Demand drops by 2% in EMEA."*
    *   **Action:** The LLM translates this NLP into **AtomScript**:
        ```typescript
        Scenario.Create("High Inflation");
        Scope({Region: "EMEA"}, () => {
             [COGS] = [COGS] * 1.08;
             [Volume] = [Volume] * 0.98;
        });
        ```
    *   **Result:** The scenario is created instantly in the Grid.

### C. Explainable Variance ("The Narrative Engine")
*   **Problem:** Variance is -10%. The CFO asks "Why?" The answer is buried in transaction logs.
*   **Cortex Solution:**
    *   **Mechanism:**
        1.  **Diff:** Cortex scans the `audit_trail` (Entropy Ledger) for the largest contributors to the variance.
        2.  **RAG:** It retrieves the "Commentary Molecules" and "External News" (e.g., "Strike in Port of LA").
        3.  **Synthesis:** It generates a narrative: *"Variance is driven by a $2M drop in Product X due to supply chain delays (Cluster A), despite a 5% price increase."*
    *   **Output:** A "Narrative Card" floats above the cell.

### D. Anomaly Detection (" The Watchtower")
*   **Problem:** Fat-finger errors (User types 10M instead of 1M).
*   **Cortex Solution:**
    *   **Mechanism:** An unsupervised **Isolation Forest** model runs on the ingestion stream (Redpanda).
    *   **Action:** If a write has a Z-Score > 3 (3 standard deviations from the norm), the cell glows **Red** and asks: *"This is 400% higher than historical average. Are you sure?"*

## 3. Implementation Check (Feasibility)
*   **We do NOT build LLMs from scratch.** We use hosted inference (OpenAI/Anthropic) or local SLMs (Llama 3 8B).
*   **We leverage the MDF:** The data is already in Parquet/Arrow. The AI models can read it natively without ETL.
*   **Stack:**
    *   **Inference:** vLLM (Python).
    *   **Orchestration:** LangChain (Go Port) or Semantic Kernel.
    *   **Storage:** We reuse the Vector capabilities of Postgres (`pgvector`) defined in Layer 2.

**Verdict:**
By adding the **Cortex Layer**, we move from a "Spreadsheet Replacement" to an "Autonomous FP&A Analyst."
