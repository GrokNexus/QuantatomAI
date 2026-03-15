-- QuantatomAI Phase 7 AI Inference Governance Validation

SELECT COUNT(*) AS missing_ai_tenant_alignment
FROM ai_inference_log ail
JOIN apps a ON a.id = ail.app_id
WHERE ail.app_id IS NOT NULL
  AND ail.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS invalid_confidence_rows
FROM ai_inference_log
WHERE confidence_score < 0 OR confidence_score > 1;

SELECT COUNT(*) AS invalid_override_reason_rows
FROM ai_inference_log
WHERE human_override = TRUE
  AND (override_reason IS NULL OR length(trim(override_reason)) = 0);
