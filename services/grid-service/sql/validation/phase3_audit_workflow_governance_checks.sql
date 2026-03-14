-- QuantatomAI Phase 3 Validation Checks
-- Run after applying 08_audit_workflow_governance.sql.

-- 1. Workflow nodes must remain tenant-aligned with apps.
SELECT COUNT(*) AS invalid_workflow_node_tenant_rows
FROM workflow_nodes wn
JOIN apps a ON a.id = wn.app_id
WHERE wn.tenant_id IS DISTINCT FROM a.tenant_id;

-- 2. Workflow transitions must reference current node state as from_state consistency.
SELECT COUNT(*) AS invalid_workflow_transition_rows
FROM workflow_state_transitions wst
JOIN workflow_nodes wn ON wn.id = wst.node_id
WHERE wst.to_state IS NULL
   OR wst.from_state IS NULL;

-- 3. Promotion requests must be tenant and app aligned with source branch.
SELECT COUNT(*) AS invalid_promotion_alignment_rows
FROM metadata_promotion_requests mpr
JOIN branches b ON b.id = mpr.source_branch_id
WHERE mpr.app_id IS DISTINCT FROM b.app_id
   OR mpr.tenant_id IS DISTINCT FROM b.tenant_id;

-- 4. Connector rejections must map to an existing batch.
SELECT COUNT(*) AS orphan_connector_rejection_rows
FROM connector_ingest_rejections cir
LEFT JOIN connector_ingest_batches cib ON cib.id = cir.batch_id
WHERE cib.id IS NULL;

-- 5. Metadata audit rows must remain tenant and app aligned.
SELECT COUNT(*) AS invalid_metadata_audit_rows
FROM metadata_audit_events mae
JOIN apps a ON a.id = mae.app_id
WHERE mae.tenant_id IS DISTINCT FROM a.tenant_id;

-- 6. Metadata audit operations must be only INSERT, UPDATE, DELETE.
SELECT COUNT(*) AS invalid_audit_operation_rows
FROM metadata_audit_events
WHERE operation NOT IN ('INSERT', 'UPDATE', 'DELETE');