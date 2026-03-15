from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
import pyarrow as pa
import pyarrow.flight as fl
import time
from rag import generate_variance_narrative, SynthesisRequest, SynthesisResponse
from db import log_ai_inference_event

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
    response = {
        "status": "queued", 
        "message": f"Auto-baselining initiated for scope {req.scoping_id}"
    }

    tenant_id = "00000000-0000-0000-0000-000000000000"
    log_ai_inference_event(
        tenant_id=tenant_id,
        request_type="auto-baseline",
        model_provider="stub",
        model_id="auto-baseline-queue-v1",
        confidence_score=0.5,
        request_payload=req.model_dump(),
        response_payload=response,
        grounding_atoms=[],
    )

    return response

@app.post("/api/v1/narrative/variance", response_model=SynthesisResponse)
async def generate_variance_narrative_endpoint(req: SynthesisRequest, http_req: Request):
    """
    Phase 8.3: Generative Interface (RAG)
    Ingests Rust-computed variance drivers and synthesizes an executive narrative.
    """
    tenant_id = http_req.headers.get("X-Tenant-ID", "unknown")
    started = time.time()
    response = generate_variance_narrative(req, tenant_id=tenant_id)
    latency_ms = int((time.time() - started) * 1000)

    log_ai_inference_event(
        tenant_id=tenant_id,
        request_type="variance-narrative",
        model_provider="openai" if response.model_id != "deterministic-fallback-v1" else "stub",
        model_id=response.model_id,
        confidence_score=response.ai_confidence_score,
        request_payload=req.model_dump(),
        response_payload=response.model_dump(),
        grounding_atoms=[d.member for d in req.drivers[:5]],
        inference_latency_ms=latency_ms,
    )

    return response

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

    response = {
        "tenant_id": tenant_id,
        "app_id": req.app_id,
        "scenario_id": req.scenario_id,
        "results": results,
        "model_id": "vector-sim-v1"
    }

    log_ai_inference_event(
        tenant_id=tenant_id,
        app_id=req.app_id,
        request_type="metadata-suggestion",
        model_provider="pgvector-stub",
        model_id="vector-sim-v1",
        confidence_score=0.8,
        request_payload=req.model_dump(),
        response_payload=response,
        grounding_atoms=[r["atom_id"] for r in results[:5]],
    )

    return response

if __name__ == "__main__":
    import uvicorn
    # Run locally on 8081 so it doesn't conflict with GridService (8080)
    uvicorn.run(app, host="0.0.0.0", port=8081)
