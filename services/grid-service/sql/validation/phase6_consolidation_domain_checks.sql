-- QuantatomAI Phase 6 Consolidation Domain Validation

SELECT COUNT(*) AS missing_close_calendar_tenant_alignment
FROM entity_close_calendar ecc
JOIN apps a ON a.id = ecc.app_id
WHERE ecc.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS missing_journal_tenant_alignment
FROM journal_entries je
JOIN apps a ON a.id = je.app_id
WHERE je.tenant_id IS DISTINCT FROM a.tenant_id;

SELECT COUNT(*) AS invalid_ownership_pct_rows
FROM intercompany_ownership
WHERE ownership_pct < 0 OR ownership_pct > 1;

SELECT COUNT(*) AS invalid_disclosure_type_rows
FROM disclosure_mappings
WHERE statement_type NOT IN ('income-statement', 'balance-sheet', 'cash-flow', 'esg');
