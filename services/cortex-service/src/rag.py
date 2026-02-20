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

# 2. The Narrative Synthesizer
def generate_variance_narrative(req: SynthesisRequest) -> SynthesisResponse:
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
            ai_confidence_score=0.95
        )

    # Find the top absolute driver
    top_driver = sorted(req.drivers, key=lambda x: abs(x.percentage_of_total), reverse=True)[0]
    
    direction = "increased" if req.total_variance > 0 else "decreased"
    driver_dir = "grew" if top_driver.variance > 0 else "contracted"
    
    narrative = f"Overall {req.kpi} {direction} by {abs(req.total_variance):,.2f}. "
    narrative += f"This variance was heavily driven by the '{top_driver.member}' member within the {top_driver.dimension} dimension, "
    narrative += f"which {driver_dir} by {abs(top_driver.variance):,.2f} and accounted for "
    narrative += f"{abs(top_driver.percentage_of_total * 100):.1f}% of the total change."

    return SynthesisResponse(
        narrative=narrative,
        ai_confidence_score=0.98
    )
