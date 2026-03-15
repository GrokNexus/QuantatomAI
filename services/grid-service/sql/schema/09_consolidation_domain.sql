-- =========================================================================================
-- QuantatomAI Phase 6: Consolidation and External Reporting Domain Pack
-- Description: Adds close calendar, intercompany ownership, journals, FX policies,
--              elimination rules, and disclosure mappings with full tenant propagation.
-- =========================================================================================

CREATE TABLE IF NOT EXISTS entity_close_calendar (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    period_code VARCHAR(32) NOT NULL,
    period_start_date DATE NOT NULL,
    period_end_date DATE NOT NULL,
    close_deadline TIMESTAMPTZ NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'open',
    owner_user_id UUID REFERENCES users(id),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, period_code)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'entity_close_calendar_status_chk') THEN
        ALTER TABLE entity_close_calendar
        ADD CONSTRAINT entity_close_calendar_status_chk
        CHECK (status IN ('open', 'in_review', 'approved', 'locked'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_close_calendar_tenant_app_status
ON entity_close_calendar (tenant_id, app_id, status);

CREATE TABLE IF NOT EXISTS intercompany_ownership (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    parent_entity_member_id UUID NOT NULL REFERENCES dimension_members(id) ON DELETE CASCADE,
    child_entity_member_id UUID NOT NULL REFERENCES dimension_members(id) ON DELETE CASCADE,
    ownership_pct NUMERIC(9,6) NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, parent_entity_member_id, child_entity_member_id, effective_from)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'intercompany_ownership_pct_chk') THEN
        ALTER TABLE intercompany_ownership
        ADD CONSTRAINT intercompany_ownership_pct_chk
        CHECK (ownership_pct >= 0 AND ownership_pct <= 1);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_intercompany_ownership_tenant_app
ON intercompany_ownership (tenant_id, app_id, effective_from DESC);

CREATE TABLE IF NOT EXISTS journal_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    close_calendar_id UUID REFERENCES entity_close_calendar(id) ON DELETE SET NULL,
    journal_type VARCHAR(32) NOT NULL DEFAULT 'adjustment',
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    description TEXT,
    source_system VARCHAR(64) NOT NULL DEFAULT 'manual',
    total_amount NUMERIC(20,4) NOT NULL DEFAULT 0,
    currency_code VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    posted_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'journal_entries_type_chk') THEN
        ALTER TABLE journal_entries
        ADD CONSTRAINT journal_entries_type_chk
        CHECK (journal_type IN ('adjustment', 'reclass', 'elimination', 'fx-translation'));
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'journal_entries_status_chk') THEN
        ALTER TABLE journal_entries
        ADD CONSTRAINT journal_entries_status_chk
        CHECK (status IN ('draft', 'submitted', 'approved', 'posted', 'rejected'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_journal_entries_tenant_app_status
ON journal_entries (tenant_id, app_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS fx_translation_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    account_member_id UUID REFERENCES dimension_members(id) ON DELETE SET NULL,
    policy_name VARCHAR(128) NOT NULL,
    translation_method VARCHAR(32) NOT NULL,
    rate_source VARCHAR(64) NOT NULL DEFAULT 'enterprise_rate_table',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, policy_name)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fx_translation_method_chk') THEN
        ALTER TABLE fx_translation_policies
        ADD CONSTRAINT fx_translation_method_chk
        CHECK (translation_method IN ('closing_rate', 'average_rate', 'historical_rate'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_fx_translation_policies_tenant_app
ON fx_translation_policies (tenant_id, app_id, translation_method);

CREATE TABLE IF NOT EXISTS elimination_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    rule_name VARCHAR(128) NOT NULL,
    debit_account_member_id UUID REFERENCES dimension_members(id) ON DELETE SET NULL,
    credit_account_member_id UUID REFERENCES dimension_members(id) ON DELETE SET NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, rule_name)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'elimination_rules_status_chk') THEN
        ALTER TABLE elimination_rules
        ADD CONSTRAINT elimination_rules_status_chk
        CHECK (status IN ('active', 'disabled'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_elimination_rules_tenant_app_status
ON elimination_rules (tenant_id, app_id, status);

CREATE TABLE IF NOT EXISTS disclosure_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    disclosure_code VARCHAR(128) NOT NULL,
    disclosure_name VARCHAR(255) NOT NULL,
    account_member_id UUID REFERENCES dimension_members(id) ON DELETE SET NULL,
    statement_type VARCHAR(32) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (app_id, disclosure_code)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'disclosure_mappings_statement_type_chk') THEN
        ALTER TABLE disclosure_mappings
        ADD CONSTRAINT disclosure_mappings_statement_type_chk
        CHECK (statement_type IN ('income-statement', 'balance-sheet', 'cash-flow', 'esg'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_disclosure_mappings_tenant_app
ON disclosure_mappings (tenant_id, app_id, statement_type);

CREATE OR REPLACE FUNCTION enforce_consolidation_tenant_alignment()
RETURNS TRIGGER AS $$
DECLARE
    app_tenant UUID;
BEGIN
    SELECT tenant_id INTO app_tenant
    FROM apps
    WHERE id = NEW.app_id;

    IF app_tenant IS NULL THEN
        RAISE EXCEPTION 'app % not found for consolidation record', NEW.app_id;
    END IF;

    IF NEW.tenant_id IS DISTINCT FROM app_tenant THEN
        RAISE EXCEPTION 'consolidation tenant mismatch: app tenant % != record tenant %', app_tenant, NEW.tenant_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_close_calendar_tenant') THEN
        CREATE TRIGGER trigger_enforce_close_calendar_tenant
        BEFORE INSERT OR UPDATE ON entity_close_calendar
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_intercompany_tenant') THEN
        CREATE TRIGGER trigger_enforce_intercompany_tenant
        BEFORE INSERT OR UPDATE ON intercompany_ownership
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_journal_tenant') THEN
        CREATE TRIGGER trigger_enforce_journal_tenant
        BEFORE INSERT OR UPDATE ON journal_entries
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_fx_policy_tenant') THEN
        CREATE TRIGGER trigger_enforce_fx_policy_tenant
        BEFORE INSERT OR UPDATE ON fx_translation_policies
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_elimination_tenant') THEN
        CREATE TRIGGER trigger_enforce_elimination_tenant
        BEFORE INSERT OR UPDATE ON elimination_rules
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_enforce_disclosure_tenant') THEN
        CREATE TRIGGER trigger_enforce_disclosure_tenant
        BEFORE INSERT OR UPDATE ON disclosure_mappings
        FOR EACH ROW EXECUTE FUNCTION enforce_consolidation_tenant_alignment();
    END IF;
END $$;
