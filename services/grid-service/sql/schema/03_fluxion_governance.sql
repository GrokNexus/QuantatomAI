-- =========================================================================================
-- QuantatomAI Layer 8.1: Fluxion AI Engine Governance
-- Description: Injects Strict AI Governance ("Kill-Switch") into the existing Schema.
-- =========================================================================================

-- 1. TENANT LEVEL GOVERNANCE
-- If this is false, NO ONE in the tenant can use Fluxion (Enterprise Data Sovereignty)
ALTER TABLE tenants 
ADD COLUMN fluxion_ai_enabled BOOLEAN DEFAULT FALSE;

-- 2. APPLICATION LEVEL GOVERNANCE
-- If tenant allows it, specific high-security apps can still disable it.
ALTER TABLE apps 
ADD COLUMN fluxion_ai_enabled BOOLEAN DEFAULT FALSE;

-- 3. USER LEVEL ROLES
-- Only Admins can execute Deep-Seek Solvers or alter the Hierarchy via Fluxion. Planners can only Auto-Forecast.
-- (Leveraging existing 'role' column in Users table: 'admin', 'planner', 'viewer')

-- 4. NEW: AUDIT LOGS FOR AI
-- Track exactly what the AI generated and who approved it.
CREATE TABLE fluxion_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID REFERENCES apps(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    action_type VARCHAR(100) NOT NULL, -- e.g. 'AUTO_FORECAST', 'GOAL_SEEK', 'GENERATIVE_VARIANCE'
    prompt_used TEXT,
    ast_generated TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_fluxion_audit_tenant ON fluxion_audit_logs(tenant_id);
CREATE INDEX idx_fluxion_audit_app ON fluxion_audit_logs(app_id);
