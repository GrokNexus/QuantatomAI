-- =========================================================================================
-- QuantatomAI Phase 7: AI Inference Governance
-- Description: Persists AI inference provenance, confidence, grounding atoms,
--              and human override metadata with strict tenant boundaries.
-- =========================================================================================

CREATE TABLE IF NOT EXISTS ai_inference_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID REFERENCES apps(id) ON DELETE SET NULL,
    request_type VARCHAR(64) NOT NULL,
    model_provider VARCHAR(64) NOT NULL,
    model_id VARCHAR(128) NOT NULL,
    confidence_score NUMERIC(5,4) NOT NULL,
    prompt_hash VARCHAR(128),
    request_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    response_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    grounding_atoms JSONB NOT NULL DEFAULT '[]'::jsonb,
    human_override BOOLEAN NOT NULL DEFAULT FALSE,
    override_reason TEXT,
    inference_latency_ms INTEGER,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ai_inference_confidence_chk') THEN
        ALTER TABLE ai_inference_log
        ADD CONSTRAINT ai_inference_confidence_chk
        CHECK (confidence_score >= 0 AND confidence_score <= 1);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ai_inference_request_type_chk') THEN
        ALTER TABLE ai_inference_log
        ADD CONSTRAINT ai_inference_request_type_chk
        CHECK (request_type IN ('variance-narrative', 'auto-baseline', 'anomaly-detection', 'metadata-suggestion', 'scenario-generation'));
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ai_inference_override_reason_chk') THEN
        ALTER TABLE ai_inference_log
        ADD CONSTRAINT ai_inference_override_reason_chk
        CHECK (NOT (human_override = TRUE AND (override_reason IS NULL OR length(trim(override_reason)) = 0)));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_ai_inference_tenant_app_time
ON ai_inference_log (tenant_id, app_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ai_inference_request_type
ON ai_inference_log (request_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ai_inference_model
ON ai_inference_log (model_provider, model_id, created_at DESC);

CREATE OR REPLACE FUNCTION enforce_ai_inference_tenant_alignment()
RETURNS TRIGGER AS $$
DECLARE
    app_tenant UUID;
BEGIN
    IF NEW.app_id IS NULL THEN
        RETURN NEW;
    END IF;

    SELECT tenant_id INTO app_tenant
    FROM apps
    WHERE id = NEW.app_id;

    IF app_tenant IS NULL THEN
        RAISE EXCEPTION 'app % not found for ai inference log', NEW.app_id;
    END IF;

    IF NEW.tenant_id IS DISTINCT FROM app_tenant THEN
        RAISE EXCEPTION 'ai inference tenant mismatch: app tenant % != record tenant %', app_tenant, NEW.tenant_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_ai_inference_tenant') THEN
        CREATE TRIGGER trigger_enforce_ai_inference_tenant
        BEFORE INSERT OR UPDATE ON ai_inference_log
        FOR EACH ROW EXECUTE FUNCTION enforce_ai_inference_tenant_alignment();
    END IF;
END $$;
