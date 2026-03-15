-- QuantatomAI Phase 4 Fixture Smoke Validation
-- Run after fixture preparation for profile C or D.

-- 1. Report alignment error counters for fast troubleshooting.
SELECT COUNT(*) AS invalid_workflow_node_tenant_rows
FROM workflow_nodes wn
JOIN apps a ON a.id = wn.app_id
WHERE wn.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS invalid_promotion_alignment_rows
FROM metadata_promotion_requests mpr
JOIN branches b ON b.id = mpr.source_branch_id
WHERE mpr.app_id IS DISTINCT FROM b.app_id
   OR mpr.tenant_id IS DISTINCT FROM b.tenant_id;

SELECT COUNT(*) AS orphan_connector_rejection_rows
FROM connector_ingest_rejections cir
LEFT JOIN connector_ingest_batches cib ON cib.id = cir.batch_id
WHERE cib.id IS NULL;

SELECT COUNT(*) AS invalid_metadata_audit_rows
FROM metadata_audit_events mae
JOIN apps a ON a.id = mae.app_id
WHERE mae.tenant_id IS DISTINCT FROM a.tenant_id;

-- 2. Enforce minimum seeded fixture footprint.
DO $$
DECLARE
    tenant_count INTEGER;
    app_count INTEGER;
    workflow_publish_count INTEGER;
    workflow_reject_count INTEGER;
    workflow_replay_count INTEGER;
    promotion_count INTEGER;
    ingest_batch_count INTEGER;
    rejection_count INTEGER;
    replay_audit_count INTEGER;
    invalid_workflow_count INTEGER;
    invalid_promotion_count INTEGER;
    orphan_rejection_count INTEGER;
    invalid_audit_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO tenant_count
    FROM tenants
    WHERE name LIKE 'phase4-tenant-%';

    SELECT COUNT(*) INTO app_count
    FROM apps
    WHERE name LIKE 'phase4-planning-%';

    SELECT COUNT(*) INTO workflow_publish_count
    FROM workflow_nodes
    WHERE node_key = 'phase4-budget-publish';

    SELECT COUNT(*) INTO workflow_reject_count
    FROM workflow_nodes
    WHERE node_key = 'phase4-budget-reject';

    SELECT COUNT(*) INTO workflow_replay_count
    FROM workflow_nodes
    WHERE node_key = 'phase4-replay-window';

    SELECT COUNT(*) INTO promotion_count
    FROM metadata_promotion_requests
    WHERE summary ILIKE 'Phase4 %';

    SELECT COUNT(*) INTO ingest_batch_count
    FROM connector_ingest_batches
    WHERE connector_name LIKE 'phase4-%';

    SELECT COUNT(*) INTO rejection_count
    FROM connector_ingest_rejections cir
    JOIN connector_ingest_batches cib ON cib.id = cir.batch_id
    WHERE cib.connector_name LIKE 'phase4-%';

    SELECT COUNT(*) INTO replay_audit_count
    FROM metadata_audit_events
    WHERE entity_type = 'replay-window';

    SELECT COUNT(*) INTO invalid_workflow_count
    FROM workflow_nodes wn
    JOIN apps a ON a.id = wn.app_id
    WHERE wn.tenant_id IS DISTINCT FROM a.tenant_id;

    SELECT COUNT(*) INTO invalid_promotion_count
    FROM metadata_promotion_requests mpr
    JOIN branches b ON b.id = mpr.source_branch_id
    WHERE mpr.app_id IS DISTINCT FROM b.app_id
       OR mpr.tenant_id IS DISTINCT FROM b.tenant_id;

    SELECT COUNT(*) INTO orphan_rejection_count
    FROM connector_ingest_rejections cir
    LEFT JOIN connector_ingest_batches cib ON cib.id = cir.batch_id
    WHERE cib.id IS NULL;

    SELECT COUNT(*) INTO invalid_audit_count
    FROM metadata_audit_events mae
    JOIN apps a ON a.id = mae.app_id
    WHERE mae.tenant_id IS DISTINCT FROM a.tenant_id;

    IF tenant_count < 3 THEN
        RAISE EXCEPTION 'Expected at least 3 phase4 tenants, found %', tenant_count;
    END IF;

    IF app_count < 3 THEN
        RAISE EXCEPTION 'Expected at least 3 phase4 apps, found %', app_count;
    END IF;

    IF workflow_publish_count < 3 OR workflow_reject_count < 3 THEN
        RAISE EXCEPTION 'Expected workflow publish/reject nodes per tenant, found publish=% reject=%', workflow_publish_count, workflow_reject_count;
    END IF;

    IF workflow_replay_count < 3 THEN
        RAISE EXCEPTION 'Expected replay workflow nodes per tenant, found %', workflow_replay_count;
    END IF;

    IF promotion_count < 9 THEN
        RAISE EXCEPTION 'Expected at least 9 phase4 promotion requests, found %', promotion_count;
    END IF;

    IF ingest_batch_count < 9 THEN
        RAISE EXCEPTION 'Expected at least 9 phase4 ingest batches, found %', ingest_batch_count;
    END IF;

    IF rejection_count < 9 THEN
        RAISE EXCEPTION 'Expected at least 9 phase4 connector rejections, found %', rejection_count;
    END IF;

    IF replay_audit_count < 3 THEN
        RAISE EXCEPTION 'Expected replay audit checkpoints for each tenant, found %', replay_audit_count;
    END IF;

    IF invalid_workflow_count <> 0 THEN
        RAISE EXCEPTION 'Workflow tenant alignment validation failed with % invalid rows', invalid_workflow_count;
    END IF;

    IF invalid_promotion_count <> 0 THEN
        RAISE EXCEPTION 'Promotion alignment validation failed with % invalid rows', invalid_promotion_count;
    END IF;

    IF orphan_rejection_count <> 0 THEN
        RAISE EXCEPTION 'Connector rejection validation failed with % orphan rows', orphan_rejection_count;
    END IF;

    IF invalid_audit_count <> 0 THEN
        RAISE EXCEPTION 'Metadata audit alignment validation failed with % invalid rows', invalid_audit_count;
    END IF;

    RAISE NOTICE 'Phase 4 fixture smoke validation passed: tenants=% apps=% promotions=% ingest_batches=% rejections=% replay_audits=%',
        tenant_count,
        app_count,
        promotion_count,
        ingest_batch_count,
        rejection_count,
        replay_audit_count;
END $$;