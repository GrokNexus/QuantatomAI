import os
from pydantic import BaseModel
from typing import List

# 1. Models for Ingestion from Rust/Go
class VarianceDriverItem(BaseModel):
    dimension: str
    member: str
    scenario_a_val: float
    scenario_b_val: float
    variance: float
    percentage_of_total: float

class SynthesisRequest(BaseModel):
    kpi: str
    total_variance: float
    drivers: List[VarianceDriverItem]

class SynthesisResponse(BaseModel):
    narrative: str
    ai_confidence_score: float
    model_id: str
    tenant_id: str

# 2. The Narrative Synthesizer
def _build_prompt(req: SynthesisRequest) -> str:
    driver_lines = []
    for d in req.drivers[:10]:
        driver_lines.append(
            f"- Dimension={d.dimension}, Member={d.member}, Variance={d.variance:.2f}, Share={d.percentage_of_total * 100:.1f}%"
        )

    return (
        "You are a CFO-grade planning narrative assistant. "
        "Write a concise, factual variance explanation grounded ONLY in provided drivers.\n\n"
        f"KPI: {req.kpi}\n"
        f"Total Variance: {req.total_variance:.2f}\n"
        "Top Drivers:\n"
        + "\n".join(driver_lines)
        + "\n\nConstraints:\n"
        "1) Mention top 1-2 drivers by contribution.\n"
        "2) Use absolute numbers and percentages.\n"
        "3) No speculative claims.\n"
        "4) Max 4 sentences."
    )


def _synthesize_with_litellm(req: SynthesisRequest) -> str:
    from litellm import completion

    model_id = os.getenv("CORTEX_LLM_MODEL", "gpt-4o-mini")
    response = completion(
        model=model_id,
        messages=[
            {"role": "system", "content": "You produce grounded financial variance narratives."},
            {"role": "user", "content": _build_prompt(req)},
        ],
        temperature=0.2,
        max_tokens=220,
    )

    return response.choices[0].message.content.strip()


def generate_variance_narrative(req: SynthesisRequest, tenant_id: str = "unknown") -> SynthesisResponse:
    """
    Acts as the LLM RAG interface. 
    It ingests the mathematically proven drivers from the Rust Attribution Engine
    and generates a plain-English executive summary.
    """
    
    # In a production environment, this would call OpenAI/Llama via litellm
    # e.g., litellm.completion(model="gpt-4", messages=[...])
    
    # For Phase 8.3 execution, we generate a high-quality deterministic narrative
    if not req.drivers:
        return SynthesisResponse(
            narrative=f"There was a {"increase" if req.total_variance > 0 else "decrease"} of {req.total_variance} in {req.kpi}, but no specific dimensional drivers were identified.",
            ai_confidence_score=0.95,
            model_id="deterministic-fallback-v1",
            tenant_id=tenant_id,
        )

    # Find the top absolute driver
    top_driver = sorted(req.drivers, key=lambda x: abs(x.percentage_of_total), reverse=True)[0]
    
    direction = "increased" if req.total_variance > 0 else "decreased"
    driver_dir = "grew" if top_driver.variance > 0 else "contracted"
    
    narrative = f"Overall {req.kpi} {direction} by {abs(req.total_variance):,.2f}. "
    narrative += f"This variance was heavily driven by the '{top_driver.member}' member within the {top_driver.dimension} dimension, "
    narrative += f"which {driver_dir} by {abs(top_driver.variance):,.2f} and accounted for "
    narrative += f"{abs(top_driver.percentage_of_total * 100):.1f}% of the total change."

    model_id = "deterministic-fallback-v1"
    confidence = 0.98

    api_key = os.getenv("CORTEX_LLM_API_KEY", "")
    if api_key:
        try:
            os.environ["OPENAI_API_KEY"] = api_key
            narrative = _synthesize_with_litellm(req)
            model_id = os.getenv("CORTEX_LLM_MODEL", "gpt-4o-mini")
            confidence = 0.90
        except Exception:
            # Keep deterministic fallback behavior if provider call fails.
            pass

    return SynthesisResponse(
        narrative=narrative,
        ai_confidence_score=confidence,
        model_id=model_id,
        tenant_id=tenant_id,
    )
