from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
import pyarrow as pa
import pyarrow.flight as fl
from rag import generate_variance_narrative, SynthesisRequest, SynthesisResponse

app = FastAPI(
    title="QuantatomAI Cortex",
    description="Layer 8: Probabilistic Intelligence & Inference Engine",
    version="0.1.0"
)

# === Domain Models ===
class ForecastRequest(BaseModel):
    scoping_id: str
    dimension_filters: dict
    horizon_periods: int

class GeneralNarrativeRequest(BaseModel):
    grid_view_id: str
    target_kpi: str

class SimilaritySearchRequest(BaseModel):
    app_id: str
    scenario_id: str
    top_k: int = 5

# === API Endpoints ===

@app.get("/health")
async def health_check():
    """
    K8s Liveness/Readiness Probe.
    """
    return {"status": "ok", "layer": "8 (Cortex)"}

@app.post("/api/v1/forecast/auto-baseline")
async def generate_baseline(req: ForecastRequest):
    """
    [STUB] Layer 8.2: The Auto-Forecast (Zero-Draft)
    Reads historical MDF, applies lag-models, returns vector.
    """
    return {
        "status": "queued", 
        "message": f"Auto-baselining initiated for scope {req.scoping_id}"
    }

@app.post("/api/v1/narrative/variance", response_model=SynthesisResponse)
async def generate_variance_narrative_endpoint(req: SynthesisRequest, http_req: Request):
    """
    Phase 8.3: Generative Interface (RAG)
    Ingests Rust-computed variance drivers and synthesizes an executive narrative.
    """
    tenant_id = http_req.headers.get("X-Tenant-ID", "unknown")
    return generate_variance_narrative(req, tenant_id=tenant_id)

@app.post("/api/v1/vector/similarity")
async def vector_similarity_search(req: SimilaritySearchRequest, http_req: Request):
    """
    Phase 7: Tenant-scoped vector retrieval contract.
    For now, this endpoint returns a deterministic tenant-scoped sample while the
    full pgvector query path is wired into grid-service storage adapters.
    """
    tenant_id = http_req.headers.get("X-Tenant-ID")
    if not tenant_id:
        raise HTTPException(status_code=401, detail="Missing X-Tenant-ID header")

    top_k = max(1, min(req.top_k, 20))
    results = []
    for idx in range(top_k):
        results.append({
            "tenant_id": tenant_id,
            "atom_id": f"{tenant_id}:{req.app_id}:{req.scenario_id}:atom_{idx}",
            "score": round(0.99 - (idx * 0.05), 4)
        })

    return {
        "tenant_id": tenant_id,
        "app_id": req.app_id,
        "scenario_id": req.scenario_id,
        "results": results,
        "model_id": "vector-sim-v1"
    }

if __name__ == "__main__":
    import uvicorn
    # Run locally on 8081 so it doesn't conflict with GridService (8080)
    uvicorn.run(app, host="0.0.0.0", port=8081)
