-- Phase 6/7 seed: consolidation close-cycle and AI inference governance bootstrap.
-- Optional runtime setting: phase4.seed_prefix (defaults to 'phase4').
DO $$
DECLARE
    v_seed_prefix TEXT := COALESCE(NULLIF(current_setting('phase4.seed_prefix', true), ''), 'phase4');
    v_app RECORD;
    v_period_code TEXT := 'FY2026M03';
    v_close_id UUID;
    v_policy_name TEXT;
    v_rule_name TEXT;
    v_disclosure_code TEXT;
    v_parent_member_id UUID;
    v_child_member_id UUID;
    v_account_member_id UUID;
BEGIN
    FOR v_app IN
        SELECT a.id AS app_id, a.tenant_id, a.name, u.id AS admin_user_id
        FROM apps a
        LEFT JOIN users u
          ON u.tenant_id = a.tenant_id
         AND u.role = 'admin'
        WHERE a.name LIKE format('%s-planning-%%', v_seed_prefix)
    LOOP
        INSERT INTO entity_close_calendar (
            tenant_id,
            app_id,
            period_code,
            period_start_date,
            period_end_date,
            close_deadline,
            status,
            owner_user_id,
            metadata
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            v_period_code,
            DATE '2026-03-01',
            DATE '2026-03-31',
            NOW() + INTERVAL '5 days',
            'in_review',
            v_app.admin_user_id,
            jsonb_build_object('source', 'phase6_seed', 'appName', v_app.name)
        )
        ON CONFLICT (app_id, period_code)
        DO UPDATE SET
            status = EXCLUDED.status,
            close_deadline = EXCLUDED.close_deadline,
            owner_user_id = EXCLUDED.owner_user_id,
            metadata = EXCLUDED.metadata,
            updated_at = NOW()
        RETURNING id INTO v_close_id;

        SELECT dm.id INTO v_parent_member_id
        FROM dimension_members dm
        WHERE dm.app_id = v_app.app_id
          AND COALESCE(dm.is_deleted, false) = false
        ORDER BY dm.created_at
        LIMIT 1;

        SELECT dm.id INTO v_child_member_id
        FROM dimension_members dm
        WHERE dm.app_id = v_app.app_id
          AND dm.id <> v_parent_member_id
          AND COALESCE(dm.is_deleted, false) = false
        ORDER BY dm.created_at
        LIMIT 1;

        SELECT dm.id INTO v_account_member_id
        FROM dimension_members dm
        WHERE dm.app_id = v_app.app_id
          AND COALESCE(dm.is_deleted, false) = false
        ORDER BY dm.created_at DESC
        LIMIT 1;

        IF v_parent_member_id IS NOT NULL AND v_child_member_id IS NOT NULL THEN
            INSERT INTO intercompany_ownership (
                tenant_id,
                app_id,
                parent_entity_member_id,
                child_entity_member_id,
                ownership_pct,
                effective_from,
                metadata
            ) VALUES (
                v_app.tenant_id,
                v_app.app_id,
                v_parent_member_id,
                v_child_member_id,
                1.0,
                DATE '2026-01-01',
                jsonb_build_object('source', 'phase6_seed')
            )
            ON CONFLICT (app_id, parent_entity_member_id, child_entity_member_id, effective_from)
            DO NOTHING;
        END IF;

        INSERT INTO journal_entries (
            tenant_id,
            app_id,
            close_calendar_id,
            journal_type,
            status,
            description,
            source_system,
            total_amount,
            currency_code,
            created_by,
            approved_by,
            posted_at,
            metadata
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            v_close_id,
            'adjustment',
            'approved',
            'Phase6 seeded adjustment journal',
            'phase6-seed',
            25000.00,
            'USD',
            v_app.admin_user_id,
            v_app.admin_user_id,
            NOW(),
            jsonb_build_object('seed', true)
        );

        v_policy_name := format('%s-fx-policy', replace(v_app.name, ' ', '-'));
        INSERT INTO fx_translation_policies (
            tenant_id,
            app_id,
            account_member_id,
            policy_name,
            translation_method,
            rate_source,
            metadata
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            v_account_member_id,
            v_policy_name,
            'closing_rate',
            'enterprise_rate_table',
            jsonb_build_object('seed', true)
        )
        ON CONFLICT (app_id, policy_name)
        DO NOTHING;

        v_rule_name := format('%s-elimination-rule', replace(v_app.name, ' ', '-'));
        INSERT INTO elimination_rules (
            tenant_id,
            app_id,
            rule_name,
            debit_account_member_id,
            credit_account_member_id,
            status,
            metadata
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            v_rule_name,
            v_parent_member_id,
            v_child_member_id,
            'active',
            jsonb_build_object('seed', true)
        )
        ON CONFLICT (app_id, rule_name)
        DO NOTHING;

        v_disclosure_code := format('%s-IS-001', upper(replace(v_app.name, '-', '_')));
        INSERT INTO disclosure_mappings (
            tenant_id,
            app_id,
            disclosure_code,
            disclosure_name,
            account_member_id,
            statement_type,
            metadata
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            v_disclosure_code,
            'Seeded Revenue Disclosure',
            v_account_member_id,
            'income-statement',
            jsonb_build_object('seed', true)
        )
        ON CONFLICT (app_id, disclosure_code)
        DO NOTHING;

        INSERT INTO ai_inference_log (
            tenant_id,
            app_id,
            request_type,
            model_provider,
            model_id,
            confidence_score,
            request_payload,
            response_payload,
            grounding_atoms,
            human_override,
            override_reason,
            inference_latency_ms,
            created_by
        ) VALUES (
            v_app.tenant_id,
            v_app.app_id,
            'variance-narrative',
            'seed-generator',
            'phase7-seed-model-v1',
            0.92,
            jsonb_build_object('source', 'phase7_seed', 'app', v_app.name),
            jsonb_build_object('narrative', 'Seeded AI inference record for governance validation.'),
            '[]'::jsonb,
            FALSE,
            NULL,
            25,
            v_app.admin_user_id
        );
    END LOOP;

    RAISE NOTICE 'Phase 6/7 close-cycle and ai-governance seed applied for apps matching prefix %', v_seed_prefix;
END
$$;
