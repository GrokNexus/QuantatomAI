-- QuantatomAI Phase 4 Planning Workload Smoke Validation (Profile B)
-- Run after fixture preparation for profile B (mixed planning workload with tenant controls).
-- Validates Phase 2 tenant control plane + Phase 3 workflow governance are populated correctly.

-- 1. Report alignment error counters for fast troubleshooting.
SELECT COUNT(*) AS invalid_workflow_node_tenant_rows
FROM workflow_nodes wn
JOIN apps a ON a.id = wn.app_id
WHERE wn.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS invalid_branch_tenant_rows
FROM branches b
JOIN apps a ON a.id = b.app_id
WHERE b.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS invalid_promotion_alignment_rows
FROM metadata_promotion_requests mpr
JOIN branches b ON b.id = mpr.source_branch_id
WHERE mpr.app_id IS DISTINCT FROM b.app_id
   OR mpr.tenant_id IS DISTINCT FROM b.tenant_id;

SELECT COUNT(*) AS invalid_metadata_audit_rows
FROM metadata_audit_events mae
JOIN apps a ON a.id = mae.app_id
WHERE mae.tenant_id IS DISTINCT FROM a.tenant_id;

-- 2. Enforce minimum seeded fixture footprint for mixed planning workload.
DO $$
DECLARE
    tenant_count           INTEGER;
    app_count              INTEGER;
    branch_count           INTEGER;
    workflow_node_count    INTEGER;
    workflow_publish_count INTEGER;
    promotion_count        INTEGER;
    audit_event_count      INTEGER;
    invalid_wn_count       INTEGER;
    invalid_branch_count   INTEGER;
    invalid_promo_count    INTEGER;
    invalid_audit_count    INTEGER;
BEGIN
    SELECT COUNT(*) INTO tenant_count
    FROM tenants
    WHERE name LIKE 'phase4-tenant-%';

    SELECT COUNT(*) INTO app_count
    FROM apps
    WHERE name LIKE 'phase4-planning-%';

    SELECT COUNT(*) INTO branch_count
    FROM branches b
    JOIN apps a ON a.id = b.app_id
    WHERE a.name LIKE 'phase4-planning-%';

    SELECT COUNT(*) INTO workflow_node_count
    FROM workflow_nodes wn
    JOIN apps a ON a.id = wn.app_id
    WHERE a.name LIKE 'phase4-planning-%';

    SELECT COUNT(*) INTO workflow_publish_count
    FROM workflow_nodes wn
    JOIN apps a ON a.id = wn.app_id
    WHERE a.name LIKE 'phase4-planning-%'
      AND wn.status IN ('approved', 'locked');

    SELECT COUNT(*) INTO promotion_count
    FROM metadata_promotion_requests mpr
    JOIN apps a ON a.id = mpr.app_id
    WHERE a.name LIKE 'phase4-planning-%';

    SELECT COUNT(*) INTO audit_event_count
    FROM metadata_audit_events mae
    JOIN apps a ON a.id = mae.app_id
    WHERE a.name LIKE 'phase4-planning-%';

    -- Alignment invariant checks
    SELECT COUNT(*) INTO invalid_wn_count
    FROM workflow_nodes wn
    JOIN apps a ON a.id = wn.app_id
    WHERE wn.tenant_id IS DISTINCT FROM a.tenant_id;

    SELECT COUNT(*) INTO invalid_branch_count
    FROM branches b
    JOIN apps a ON a.id = b.app_id
    WHERE b.tenant_id IS DISTINCT FROM a.tenant_id;

    SELECT COUNT(*) INTO invalid_promo_count
    FROM metadata_promotion_requests mpr
    JOIN branches b ON b.id = mpr.source_branch_id
    WHERE mpr.app_id IS DISTINCT FROM b.app_id
       OR mpr.tenant_id IS DISTINCT FROM b.tenant_id;

    SELECT COUNT(*) INTO invalid_audit_count
    FROM metadata_audit_events mae
    JOIN apps a ON a.id = mae.app_id
    WHERE mae.tenant_id IS DISTINCT FROM a.tenant_id;

    -- Alignment invariants must be zero
    IF invalid_wn_count > 0 THEN
        RAISE EXCEPTION 'Phase 4 profile B: % workflow_nodes have mismatched tenant_id', invalid_wn_count;
    END IF;

    IF invalid_branch_count > 0 THEN
        RAISE EXCEPTION 'Phase 4 profile B: % branches have mismatched tenant_id', invalid_branch_count;
    END IF;

    IF invalid_promo_count > 0 THEN
        RAISE EXCEPTION 'Phase 4 profile B: % promotion_requests have mismatched app_id or tenant_id', invalid_promo_count;
    END IF;

    IF invalid_audit_count > 0 THEN
        RAISE EXCEPTION 'Phase 4 profile B: % metadata_audit_events have mismatched tenant_id', invalid_audit_count;
    END IF;

    -- Minimum footprint enforcement
    IF tenant_count < 3 THEN
        RAISE EXCEPTION 'Phase 4 profile B: insufficient tenant fixture footprint (found %, need >=3)', tenant_count;
    END IF;

    IF app_count < 3 THEN
        RAISE EXCEPTION 'Phase 4 profile B: insufficient app fixture footprint (found %, need >=3)', app_count;
    END IF;

    IF branch_count < 6 THEN
        RAISE EXCEPTION 'Phase 4 profile B: insufficient branch footprint for mixed-workload (found %, need >=6)', branch_count;
    END IF;

    IF workflow_node_count < 3 THEN
        RAISE EXCEPTION 'Phase 4 profile B: insufficient workflow_node footprint (found %, need >=3)', workflow_node_count;
    END IF;

    IF promotion_count < 3 THEN
        RAISE EXCEPTION 'Phase 4 profile B: insufficient promotion_request footprint (found %, need >=3)', promotion_count;
    END IF;

    RAISE NOTICE 'Phase 4 profile B planning-workload smoke validation passed: tenants=% apps=% branches=% workflow_nodes=% promotions=% audit_events=%',
        tenant_count, app_count, branch_count, workflow_node_count, promotion_count, audit_event_count;
END
$$;
