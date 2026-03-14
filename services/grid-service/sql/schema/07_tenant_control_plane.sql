-- =========================================================================================
-- QuantatomAI Phase 2: Multi-Tenant Control Plane
-- Description: Hardens tenant identity, residency, quota, key-domain, AI-boundary,
--              and partition metadata so the database layer can evolve toward
--              enterprise-grade planning, consolidation, and reporting.
-- =========================================================================================

-- 1. TENANT HARDENING
ALTER TABLE tenants
ADD COLUMN IF NOT EXISTS status VARCHAR(32) NOT NULL DEFAULT 'active',
ADD COLUMN IF NOT EXISTS residency_mode VARCHAR(32) NOT NULL DEFAULT 'single-region',
ADD COLUMN IF NOT EXISTS primary_region VARCHAR(32) NOT NULL DEFAULT 'us-east-1',
ADD COLUMN IF NOT EXISTS isolation_tier VARCHAR(32) NOT NULL DEFAULT 'logical',
ADD COLUMN IF NOT EXISTS ai_learning_mode VARCHAR(32) NOT NULL DEFAULT 'tenant-only',
ADD COLUMN IF NOT EXISTS cost_center_code VARCHAR(64),
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_plan_tier_chk'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT tenants_plan_tier_chk
        CHECK (plan_tier IN ('standard', 'enterprise', 'ultra'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_status_chk'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT tenants_status_chk
        CHECK (status IN ('active', 'suspended', 'archived'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_residency_mode_chk'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT tenants_residency_mode_chk
        CHECK (residency_mode IN ('single-region', 'geo-fenced', 'sovereign'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_isolation_tier_chk'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT tenants_isolation_tier_chk
        CHECK (isolation_tier IN ('logical', 'dedicated-cache', 'dedicated-data-plane', 'dedicated-cluster'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_ai_learning_mode_chk'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT tenants_ai_learning_mode_chk
        CHECK (ai_learning_mode IN ('disabled', 'tenant-only', 'federated-opt-in'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_tenants_primary_region ON tenants(primary_region);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);

-- 2. TENANT REGION REGISTRY
CREATE TABLE IF NOT EXISTS tenant_regions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    region_code VARCHAR(32) NOT NULL,
    region_role VARCHAR(32) NOT NULL DEFAULT 'primary',
    residency_class VARCHAR(32) NOT NULL DEFAULT 'customer-data',
    is_write_region BOOLEAN NOT NULL DEFAULT FALSE,
    is_read_region BOOLEAN NOT NULL DEFAULT TRUE,
    is_failover_region BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, region_code)
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_regions_role_chk'
    ) THEN
        ALTER TABLE tenant_regions
        ADD CONSTRAINT tenant_regions_role_chk
        CHECK (region_role IN ('primary', 'secondary', 'failover', 'archive', 'analytics'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_regions_residency_class_chk'
    ) THEN
        ALTER TABLE tenant_regions
        ADD CONSTRAINT tenant_regions_residency_class_chk
        CHECK (residency_class IN ('customer-data', 'audit', 'metadata', 'ai-features', 'archive'));
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_regions_single_write_region
ON tenant_regions(tenant_id)
WHERE is_write_region = TRUE;

-- 3. TENANT KEY DOMAINS
CREATE TABLE IF NOT EXISTS tenant_key_domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    region_code VARCHAR(32) NOT NULL,
    purpose VARCHAR(32) NOT NULL,
    kms_provider VARCHAR(32) NOT NULL,
    key_uri TEXT NOT NULL,
    rotation_interval_days INTEGER NOT NULL DEFAULT 90,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, region_code, purpose),
    FOREIGN KEY (tenant_id, region_code)
        REFERENCES tenant_regions(tenant_id, region_code)
        ON DELETE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_key_domains_purpose_chk'
    ) THEN
        ALTER TABLE tenant_key_domains
        ADD CONSTRAINT tenant_key_domains_purpose_chk
        CHECK (purpose IN ('app-data', 'audit', 'embedding', 'export', 'backup'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_key_domains_kms_provider_chk'
    ) THEN
        ALTER TABLE tenant_key_domains
        ADD CONSTRAINT tenant_key_domains_kms_provider_chk
        CHECK (kms_provider IN ('aws-kms', 'azure-key-vault', 'gcp-kms', 'oci-vault', 'hashicorp-vault'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_key_domains_rotation_days_chk'
    ) THEN
        ALTER TABLE tenant_key_domains
        ADD CONSTRAINT tenant_key_domains_rotation_days_chk
        CHECK (rotation_interval_days > 0);
    END IF;
END $$;

-- 4. TENANT QUOTA AND CHARGEBACK POLICY
CREATE TABLE IF NOT EXISTS tenant_quota_policies (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    max_users INTEGER NOT NULL DEFAULT 500,
    max_apps INTEGER NOT NULL DEFAULT 100,
    max_storage_gb INTEGER NOT NULL DEFAULT 10240,
    max_hot_working_set_gb INTEGER NOT NULL DEFAULT 512,
    max_events_per_sec INTEGER NOT NULL DEFAULT 50000,
    max_api_rps INTEGER NOT NULL DEFAULT 5000,
    max_concurrent_jobs INTEGER NOT NULL DEFAULT 100,
    max_vector_bytes BIGINT NOT NULL DEFAULT 10737418240,
    chargeback_model VARCHAR(32) NOT NULL DEFAULT 'showback',
    monthly_spend_limit_usd NUMERIC(18,2),
    overage_behavior VARCHAR(32) NOT NULL DEFAULT 'throttle',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_quota_chargeback_chk'
    ) THEN
        ALTER TABLE tenant_quota_policies
        ADD CONSTRAINT tenant_quota_chargeback_chk
        CHECK (chargeback_model IN ('showback', 'hard-chargeback', 'none'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_quota_overage_behavior_chk'
    ) THEN
        ALTER TABLE tenant_quota_policies
        ADD CONSTRAINT tenant_quota_overage_behavior_chk
        CHECK (overage_behavior IN ('throttle', 'alert', 'block'));
    END IF;
END $$;

-- 5. TENANT AI BOUNDARY POLICY
CREATE TABLE IF NOT EXISTS tenant_ai_policies (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    retrieval_scope VARCHAR(32) NOT NULL DEFAULT 'tenant-only',
    allow_cross_tenant_learning BOOLEAN NOT NULL DEFAULT FALSE,
    allow_external_inference BOOLEAN NOT NULL DEFAULT FALSE,
    require_prompt_audit BOOLEAN NOT NULL DEFAULT TRUE,
    require_human_approval_for_generative_write BOOLEAN NOT NULL DEFAULT TRUE,
    max_context_rows INTEGER NOT NULL DEFAULT 5000,
    vector_namespace_strategy VARCHAR(32) NOT NULL DEFAULT 'tenant-segregated',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_ai_policies_retrieval_scope_chk'
    ) THEN
        ALTER TABLE tenant_ai_policies
        ADD CONSTRAINT tenant_ai_policies_retrieval_scope_chk
        CHECK (retrieval_scope IN ('tenant-only', 'tenant-region-only'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_ai_policies_namespace_strategy_chk'
    ) THEN
        ALTER TABLE tenant_ai_policies
        ADD CONSTRAINT tenant_ai_policies_namespace_strategy_chk
        CHECK (vector_namespace_strategy IN ('tenant-segregated', 'tenant-region-segregated'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenant_ai_policies_context_rows_chk'
    ) THEN
        ALTER TABLE tenant_ai_policies
        ADD CONSTRAINT tenant_ai_policies_context_rows_chk
        CHECK (max_context_rows > 0);
    END IF;
END $$;

-- 6. APP PARTITION REGISTRY
CREATE TABLE IF NOT EXISTS app_partitions (
    app_id UUID PRIMARY KEY REFERENCES apps(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    write_region VARCHAR(32) NOT NULL,
    hot_namespace VARCHAR(128) NOT NULL,
    warm_partition_template TEXT NOT NULL,
    cold_object_prefix TEXT NOT NULL,
    event_topic_prefix VARCHAR(128) NOT NULL,
    cache_namespace VARCHAR(128) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (tenant_id, write_region)
        REFERENCES tenant_regions(tenant_id, region_code)
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_app_partitions_cache_namespace
ON app_partitions(cache_namespace);

CREATE UNIQUE INDEX IF NOT EXISTS idx_app_partitions_event_topic_prefix
ON app_partitions(event_topic_prefix);

-- 7. PROPAGATE TENANT CONTEXT TO CORE TABLES
ALTER TABLE dimensions ADD COLUMN IF NOT EXISTS tenant_id UUID;
UPDATE dimensions d
SET tenant_id = a.tenant_id
FROM apps a
WHERE d.app_id = a.id
  AND d.tenant_id IS NULL;
ALTER TABLE dimensions ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE dimension_members ADD COLUMN IF NOT EXISTS app_id UUID;
ALTER TABLE dimension_members ADD COLUMN IF NOT EXISTS tenant_id UUID;
UPDATE dimension_members dm
SET app_id = d.app_id,
    tenant_id = d.tenant_id
FROM dimensions d
WHERE dm.dimension_id = d.id
  AND (dm.app_id IS NULL OR dm.tenant_id IS NULL);
ALTER TABLE dimension_members ALTER COLUMN app_id SET NOT NULL;
ALTER TABLE dimension_members ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE security_policies ADD COLUMN IF NOT EXISTS tenant_id UUID;
UPDATE security_policies sp
SET tenant_id = a.tenant_id
FROM apps a
WHERE sp.app_id = a.id
  AND sp.tenant_id IS NULL;
ALTER TABLE security_policies ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE branches ADD COLUMN IF NOT EXISTS tenant_id UUID;
UPDATE branches b
SET tenant_id = a.tenant_id
FROM apps a
WHERE b.app_id = a.id
  AND b.tenant_id IS NULL;
ALTER TABLE branches ALTER COLUMN tenant_id SET NOT NULL;

-- 8. FOREIGN KEYS AND INDEXES FOR PROPAGATED CONTEXT
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'dimensions_tenant_fk'
    ) THEN
        ALTER TABLE dimensions
        ADD CONSTRAINT dimensions_tenant_fk
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'dimension_members_app_fk'
    ) THEN
        ALTER TABLE dimension_members
        ADD CONSTRAINT dimension_members_app_fk
        FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'dimension_members_tenant_fk'
    ) THEN
        ALTER TABLE dimension_members
        ADD CONSTRAINT dimension_members_tenant_fk
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'security_policies_tenant_fk'
    ) THEN
        ALTER TABLE security_policies
        ADD CONSTRAINT security_policies_tenant_fk
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'branches_tenant_fk'
    ) THEN
        ALTER TABLE branches
        ADD CONSTRAINT branches_tenant_fk
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_dimensions_tenant_id ON dimensions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dimension_members_tenant_app ON dimension_members(tenant_id, app_id, dimension_id);
CREATE INDEX IF NOT EXISTS idx_security_policies_tenant_user ON security_policies(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_branches_tenant_app ON branches(tenant_id, app_id);

-- 9. TENANT CONTEXT SYNCHRONIZATION TRIGGERS
CREATE OR REPLACE FUNCTION sync_dimension_tenant_context()
RETURNS TRIGGER AS $$
DECLARE
    source_tenant_id UUID;
BEGIN
    SELECT tenant_id INTO source_tenant_id
    FROM apps
    WHERE id = NEW.app_id;

    IF source_tenant_id IS NULL THEN
        RAISE EXCEPTION 'app % does not resolve to a tenant', NEW.app_id;
    END IF;

    NEW.tenant_id := source_tenant_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_member_tenant_context()
RETURNS TRIGGER AS $$
DECLARE
    source_app_id UUID;
    source_tenant_id UUID;
BEGIN
    SELECT app_id, tenant_id INTO source_app_id, source_tenant_id
    FROM dimensions
    WHERE id = NEW.dimension_id;

    IF source_app_id IS NULL OR source_tenant_id IS NULL THEN
        RAISE EXCEPTION 'dimension % does not resolve to app and tenant context', NEW.dimension_id;
    END IF;

    NEW.app_id := source_app_id;
    NEW.tenant_id := source_tenant_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_policy_tenant_context()
RETURNS TRIGGER AS $$
DECLARE
    source_tenant_id UUID;
BEGIN
    SELECT tenant_id INTO source_tenant_id
    FROM apps
    WHERE id = NEW.app_id;

    IF source_tenant_id IS NULL THEN
        RAISE EXCEPTION 'policy app % does not resolve to tenant', NEW.app_id;
    END IF;

    NEW.tenant_id := source_tenant_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_branch_tenant_context()
RETURNS TRIGGER AS $$
DECLARE
    source_tenant_id UUID;
BEGIN
    SELECT tenant_id INTO source_tenant_id
    FROM apps
    WHERE id = NEW.app_id;

    IF source_tenant_id IS NULL THEN
        RAISE EXCEPTION 'branch app % does not resolve to tenant', NEW.app_id;
    END IF;

    NEW.tenant_id := source_tenant_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_app_partition_tenant_context()
RETURNS TRIGGER AS $$
DECLARE
    source_tenant_id UUID;
BEGIN
    SELECT tenant_id INTO source_tenant_id
    FROM apps
    WHERE id = NEW.app_id;

    IF source_tenant_id IS NULL THEN
        RAISE EXCEPTION 'app partition app % does not resolve to tenant', NEW.app_id;
    END IF;

    NEW.tenant_id := source_tenant_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_sync_dimension_tenant_context') THEN
        CREATE TRIGGER trigger_sync_dimension_tenant_context
        BEFORE INSERT OR UPDATE OF app_id ON dimensions
        FOR EACH ROW EXECUTE FUNCTION sync_dimension_tenant_context();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_sync_member_tenant_context') THEN
        CREATE TRIGGER trigger_sync_member_tenant_context
        BEFORE INSERT OR UPDATE OF dimension_id ON dimension_members
        FOR EACH ROW EXECUTE FUNCTION sync_member_tenant_context();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_sync_policy_tenant_context') THEN
        CREATE TRIGGER trigger_sync_policy_tenant_context
        BEFORE INSERT OR UPDATE OF app_id ON security_policies
        FOR EACH ROW EXECUTE FUNCTION sync_policy_tenant_context();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_sync_branch_tenant_context') THEN
        CREATE TRIGGER trigger_sync_branch_tenant_context
        BEFORE INSERT OR UPDATE OF app_id ON branches
        FOR EACH ROW EXECUTE FUNCTION sync_branch_tenant_context();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_sync_app_partition_tenant_context') THEN
        CREATE TRIGGER trigger_sync_app_partition_tenant_context
        BEFORE INSERT OR UPDATE OF app_id ON app_partitions
        FOR EACH ROW EXECUTE FUNCTION sync_app_partition_tenant_context();
    END IF;
END $$;

-- 10. UPDATED_AT TRIGGERS FOR CONTROL-PLANE TABLES
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_tenant_timestamp') THEN
        CREATE TRIGGER trigger_update_tenant_timestamp
        BEFORE UPDATE ON tenants
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_tenant_regions_timestamp') THEN
        CREATE TRIGGER trigger_update_tenant_regions_timestamp
        BEFORE UPDATE ON tenant_regions
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_tenant_key_domains_timestamp') THEN
        CREATE TRIGGER trigger_update_tenant_key_domains_timestamp
        BEFORE UPDATE ON tenant_key_domains
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_tenant_quota_policies_timestamp') THEN
        CREATE TRIGGER trigger_update_tenant_quota_policies_timestamp
        BEFORE UPDATE ON tenant_quota_policies
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_tenant_ai_policies_timestamp') THEN
        CREATE TRIGGER trigger_update_tenant_ai_policies_timestamp
        BEFORE UPDATE ON tenant_ai_policies
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_app_partitions_timestamp') THEN
        CREATE TRIGGER trigger_update_app_partitions_timestamp
        BEFORE UPDATE ON app_partitions
        FOR EACH ROW EXECUTE FUNCTION update_timestamp();
    END IF;
END $$;