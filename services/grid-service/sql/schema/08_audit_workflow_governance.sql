-- =========================================================================================
-- QuantatomAI Phase 3: Audit, Lineage, and Workflow Governance
-- Description: Adds immutable metadata audit events, workflow state controls,
--              metadata promotion governance, and connector staging governance.
-- =========================================================================================

-- 1. IMMUTABLE METADATA AUDIT EVENTS
CREATE TABLE IF NOT EXISTS metadata_audit_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    entity_type VARCHAR(64) NOT NULL,
    entity_id UUID NOT NULL,
    operation VARCHAR(16) NOT NULL,
    actor_user_id UUID REFERENCES users(id),
    source_channel VARCHAR(32) NOT NULL DEFAULT 'system',
    trace_id VARCHAR(128),
    old_data JSONB,
    new_data JSONB,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'metadata_audit_operation_chk'
    ) THEN
        ALTER TABLE metadata_audit_events
        ADD CONSTRAINT metadata_audit_operation_chk
        CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'metadata_audit_source_channel_chk'
    ) THEN
        ALTER TABLE metadata_audit_events
        ADD CONSTRAINT metadata_audit_source_channel_chk
        CHECK (source_channel IN ('system', 'api', 'ui', 'ingestion', 'migration'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_metadata_audit_tenant_app_time
ON metadata_audit_events(tenant_id, app_id, occurred_at DESC);

CREATE INDEX IF NOT EXISTS idx_metadata_audit_entity
ON metadata_audit_events(entity_type, entity_id, occurred_at DESC);

CREATE INDEX IF NOT EXISTS idx_metadata_audit_operation
ON metadata_audit_events(operation, occurred_at DESC);

CREATE OR REPLACE FUNCTION emit_metadata_audit_event()
RETURNS TRIGGER AS $$
DECLARE
    resolved_tenant_id UUID;
    resolved_app_id UUID;
    resolved_entity_id UUID;
    old_payload JSONB;
    new_payload JSONB;
BEGIN
    IF TG_OP = 'DELETE' THEN
        old_payload := to_jsonb(OLD);
        new_payload := NULL;
    ELSE
        old_payload := CASE WHEN TG_OP = 'UPDATE' THEN to_jsonb(OLD) ELSE NULL END;
        new_payload := to_jsonb(NEW);
    END IF;

    IF TG_TABLE_NAME = 'dimensions' THEN
        resolved_tenant_id := COALESCE(NEW.tenant_id, OLD.tenant_id);
        resolved_app_id := COALESCE(NEW.app_id, OLD.app_id);
        resolved_entity_id := COALESCE(NEW.id, OLD.id);
    ELSIF TG_TABLE_NAME = 'dimension_members' THEN
        resolved_tenant_id := COALESCE(NEW.tenant_id, OLD.tenant_id);
        resolved_app_id := COALESCE(NEW.app_id, OLD.app_id);
        resolved_entity_id := COALESCE(NEW.id, OLD.id);
    ELSIF TG_TABLE_NAME = 'security_policies' THEN
        resolved_tenant_id := COALESCE(NEW.tenant_id, OLD.tenant_id);
        resolved_app_id := COALESCE(NEW.app_id, OLD.app_id);
        resolved_entity_id := COALESCE(NEW.id, OLD.id);
    ELSIF TG_TABLE_NAME = 'branches' THEN
        resolved_tenant_id := COALESCE(NEW.tenant_id, OLD.tenant_id);
        resolved_app_id := COALESCE(NEW.app_id, OLD.app_id);
        resolved_entity_id := COALESCE(NEW.id, OLD.id);
    ELSE
        RAISE EXCEPTION 'audit trigger unsupported for table %', TG_TABLE_NAME;
    END IF;

    INSERT INTO metadata_audit_events (
        tenant_id,
        app_id,
        entity_type,
        entity_id,
        operation,
        old_data,
        new_data,
        source_channel,
        metadata
    ) VALUES (
        resolved_tenant_id,
        resolved_app_id,
        TG_TABLE_NAME,
        resolved_entity_id,
        TG_OP,
        old_payload,
        new_payload,
        'system',
        jsonb_build_object('trigger', TG_NAME)
    );

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_audit_dimensions') THEN
        CREATE TRIGGER trigger_audit_dimensions
        AFTER INSERT OR UPDATE OR DELETE ON dimensions
        FOR EACH ROW EXECUTE FUNCTION emit_metadata_audit_event();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_audit_dimension_members') THEN
        CREATE TRIGGER trigger_audit_dimension_members
        AFTER INSERT OR UPDATE OR DELETE ON dimension_members
        FOR EACH ROW EXECUTE FUNCTION emit_metadata_audit_event();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_audit_security_policies') THEN
        CREATE TRIGGER trigger_audit_security_policies
        AFTER INSERT OR UPDATE OR DELETE ON security_policies
        FOR EACH ROW EXECUTE FUNCTION emit_metadata_audit_event();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_audit_branches') THEN
        CREATE TRIGGER trigger_audit_branches
        AFTER INSERT OR UPDATE OR DELETE ON branches
        FOR EACH ROW EXECUTE FUNCTION emit_metadata_audit_event();
    END IF;
END $$;

-- 2. WORKFLOW GOVERNANCE
CREATE TABLE IF NOT EXISTS workflow_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    node_key VARCHAR(255) NOT NULL,
    state VARCHAR(32) NOT NULL DEFAULT 'draft',
    owner_user_id UUID REFERENCES users(id),
    approver_user_id UUID REFERENCES users(id),
    lock_mode VARCHAR(32) NOT NULL DEFAULT 'editable',
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, node_key)
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'workflow_nodes_state_chk'
    ) THEN
        ALTER TABLE workflow_nodes
        ADD CONSTRAINT workflow_nodes_state_chk
        CHECK (state IN ('draft', 'in_review', 'rejected', 'approved', 'published'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'workflow_nodes_lock_mode_chk'
    ) THEN
        ALTER TABLE workflow_nodes
        ADD CONSTRAINT workflow_nodes_lock_mode_chk
        CHECK (lock_mode IN ('editable', 'review-locked', 'approved-locked', 'published-locked'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_workflow_nodes_tenant_app_state
ON workflow_nodes(tenant_id, app_id, state);

CREATE TABLE IF NOT EXISTS workflow_state_transitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    node_id UUID NOT NULL REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    from_state VARCHAR(32) NOT NULL,
    to_state VARCHAR(32) NOT NULL,
    actor_user_id UUID REFERENCES users(id),
    reason TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'workflow_transitions_from_state_chk'
    ) THEN
        ALTER TABLE workflow_state_transitions
        ADD CONSTRAINT workflow_transitions_from_state_chk
        CHECK (from_state IN ('draft', 'in_review', 'rejected', 'approved', 'published'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'workflow_transitions_to_state_chk'
    ) THEN
        ALTER TABLE workflow_state_transitions
        ADD CONSTRAINT workflow_transitions_to_state_chk
        CHECK (to_state IN ('draft', 'in_review', 'rejected', 'approved', 'published'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_workflow_state_transitions_node_time
ON workflow_state_transitions(node_id, created_at DESC);

CREATE OR REPLACE FUNCTION enforce_workflow_transition_rules()
RETURNS TRIGGER AS $$
DECLARE
    current_state VARCHAR(32);
BEGIN
    SELECT state INTO current_state
    FROM workflow_nodes
    WHERE id = NEW.node_id;

    IF current_state IS NULL THEN
        RAISE EXCEPTION 'workflow node % not found', NEW.node_id;
    END IF;

    IF current_state <> NEW.from_state THEN
        RAISE EXCEPTION 'workflow state mismatch for node %: expected %, got %', NEW.node_id, current_state, NEW.from_state;
    END IF;

    IF NOT (
        (NEW.from_state = 'draft' AND NEW.to_state = 'in_review') OR
        (NEW.from_state = 'in_review' AND NEW.to_state IN ('approved', 'rejected')) OR
        (NEW.from_state = 'rejected' AND NEW.to_state = 'draft') OR
        (NEW.from_state = 'approved' AND NEW.to_state = 'published') OR
        (NEW.from_state = 'published' AND NEW.to_state = 'draft')
    ) THEN
        RAISE EXCEPTION 'invalid workflow transition: % -> %', NEW.from_state, NEW.to_state;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION apply_workflow_transition()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE workflow_nodes
    SET state = NEW.to_state,
        is_locked = CASE WHEN NEW.to_state IN ('approved', 'published') THEN TRUE ELSE FALSE END,
        lock_mode = CASE
            WHEN NEW.to_state = 'draft' THEN 'editable'
            WHEN NEW.to_state = 'in_review' THEN 'review-locked'
            WHEN NEW.to_state = 'rejected' THEN 'editable'
            WHEN NEW.to_state = 'approved' THEN 'approved-locked'
            WHEN NEW.to_state = 'published' THEN 'published-locked'
            ELSE lock_mode
        END,
        updated_at = NOW()
    WHERE id = NEW.node_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_workflow_transition') THEN
        CREATE TRIGGER trigger_enforce_workflow_transition
        BEFORE INSERT ON workflow_state_transitions
        FOR EACH ROW EXECUTE FUNCTION enforce_workflow_transition_rules();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_apply_workflow_transition') THEN
        CREATE TRIGGER trigger_apply_workflow_transition
        AFTER INSERT ON workflow_state_transitions
        FOR EACH ROW EXECUTE FUNCTION apply_workflow_transition();
    END IF;
END $$;

-- 3. METADATA PROMOTION GOVERNANCE
CREATE TABLE IF NOT EXISTS metadata_promotion_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    source_branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    target_branch_name VARCHAR(255) NOT NULL DEFAULT 'main',
    requested_by_user_id UUID REFERENCES users(id),
    approved_by_user_id UUID REFERENCES users(id),
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    risk_level VARCHAR(32) NOT NULL DEFAULT 'medium',
    summary TEXT,
    diff_fingerprint VARCHAR(128),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'metadata_promotion_status_chk'
    ) THEN
        ALTER TABLE metadata_promotion_requests
        ADD CONSTRAINT metadata_promotion_status_chk
        CHECK (status IN ('pending', 'approved', 'rejected', 'applied', 'cancelled'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'metadata_promotion_risk_chk'
    ) THEN
        ALTER TABLE metadata_promotion_requests
        ADD CONSTRAINT metadata_promotion_risk_chk
        CHECK (risk_level IN ('low', 'medium', 'high', 'critical'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_metadata_promotion_tenant_app_status
ON metadata_promotion_requests(tenant_id, app_id, status, requested_at DESC);

-- 4. CONNECTOR STAGING GOVERNANCE
CREATE TABLE IF NOT EXISTS connector_ingest_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    connector_name VARCHAR(128) NOT NULL,
    source_uri TEXT,
    ingest_status VARCHAR(32) NOT NULL DEFAULT 'received',
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    validated_at TIMESTAMPTZ,
    applied_at TIMESTAMPTZ,
    records_received BIGINT NOT NULL DEFAULT 0,
    records_applied BIGINT NOT NULL DEFAULT 0,
    records_rejected BIGINT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'connector_ingest_status_chk'
    ) THEN
        ALTER TABLE connector_ingest_batches
        ADD CONSTRAINT connector_ingest_status_chk
        CHECK (ingest_status IN ('received', 'validated', 'applied', 'partial', 'failed', 'quarantined'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_connector_ingest_batches_tenant_app_time
ON connector_ingest_batches(tenant_id, app_id, received_at DESC);

CREATE TABLE IF NOT EXISTS connector_ingest_rejections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES connector_ingest_batches(id) ON DELETE CASCADE,
    record_key VARCHAR(255),
    reason_code VARCHAR(64) NOT NULL,
    reason_detail TEXT,
    raw_record JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_connector_ingest_rejections_batch
ON connector_ingest_rejections(batch_id, created_at DESC);

-- 5. UPDATED_AT TRIGGERS
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_workflow_nodes_timestamp') THEN
        CREATE TRIGGER trigger_update_workflow_nodes_timestamp
        BEFORE UPDATE ON workflow_nodes
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_metadata_promotion_requests_timestamp') THEN
        CREATE TRIGGER trigger_update_metadata_promotion_requests_timestamp
        BEFORE UPDATE ON metadata_promotion_requests
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_connector_ingest_batches_timestamp') THEN
        CREATE TRIGGER trigger_update_connector_ingest_batches_timestamp
        BEFORE UPDATE ON connector_ingest_batches
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;
END $$;