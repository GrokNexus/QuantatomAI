-- Phase 4 replay and recovery fixtures.
-- Optional runtime setting: phase4.seed_prefix (defaults to 'phase4').
DO $$
DECLARE
    v_seed_prefix TEXT := COALESCE(NULLIF(current_setting('phase4.seed_prefix', true), ''), 'phase4');
    v_tenant_labels TEXT[] := ARRAY['alpha', 'beta', 'gamma'];
    v_tenant_label TEXT;
    v_tenant_id UUID;
    v_admin_user_id UUID;
    v_app_id UUID;
    v_main_branch_id UUID;
    v_recovery_branch_id UUID;
    v_account_dimension_id UUID;
    v_account_root_id UUID;
    v_revenue_override_id UUID;
    v_cogs_tombstone_id UUID;
    v_replay_node_id UUID;
    v_quarantined_batch_id UUID;
BEGIN
    FOREACH v_tenant_label IN ARRAY v_tenant_labels LOOP
        SELECT id INTO v_tenant_id
        FROM tenants
        WHERE name = format('%s-tenant-%s', v_seed_prefix, v_tenant_label)
        ORDER BY created_at
        LIMIT 1;

        IF v_tenant_id IS NULL THEN
            RAISE NOTICE 'Skipping replay fixture for tenant label %, tenant baseline not found.', v_tenant_label;
            CONTINUE;
        END IF;

        SELECT id INTO v_admin_user_id
        FROM users
        WHERE tenant_id = v_tenant_id
          AND email = format('admin+%s@phase4.quantatom.ai', v_tenant_label)
        LIMIT 1;

        SELECT id INTO v_app_id
        FROM apps
        WHERE tenant_id = v_tenant_id
          AND name = format('%s-planning-%s', v_seed_prefix, v_tenant_label)
        LIMIT 1;

        SELECT id INTO v_main_branch_id
        FROM branches
        WHERE app_id = v_app_id
          AND name = 'main'
        LIMIT 1;

        INSERT INTO branches (app_id, tenant_id, name, base_branch_id)
        VALUES (v_app_id, v_tenant_id, 'phase4_recovery', v_main_branch_id)
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            base_branch_id = EXCLUDED.base_branch_id,
            updated_at = NOW()
        RETURNING id INTO v_recovery_branch_id;

        SELECT id INTO v_account_dimension_id
        FROM dimensions
        WHERE app_id = v_app_id
          AND name = 'Account'
        LIMIT 1;

        SELECT id INTO v_account_root_id
        FROM dimension_members
        WHERE dimension_id = v_account_dimension_id
          AND branch_id = v_main_branch_id
          AND name = 'AccountRoot'
        LIMIT 1;

        INSERT INTO dimension_members (
            dimension_id,
            app_id,
            tenant_id,
            branch_id,
            name,
            path,
            parent_id,
            weight,
            attributes,
            formula,
            is_deleted
        )
        VALUES (
            v_account_dimension_id,
            v_app_id,
            v_tenant_id,
            v_recovery_branch_id,
            'Revenue',
            'account_root.revenue'::ltree,
            v_account_root_id,
            10,
            jsonb_build_object('seed', v_seed_prefix, 'replay_window', '2024Q4', 'recovered', TRUE),
            'SUM(children) * 1.01',
            FALSE
        )
        ON CONFLICT (dimension_id, name, branch_id)
        DO UPDATE SET
            app_id = EXCLUDED.app_id,
            tenant_id = EXCLUDED.tenant_id,
            path = EXCLUDED.path,
            parent_id = EXCLUDED.parent_id,
            weight = EXCLUDED.weight,
            attributes = EXCLUDED.attributes,
            formula = EXCLUDED.formula,
            is_deleted = EXCLUDED.is_deleted
        RETURNING id INTO v_revenue_override_id;

        INSERT INTO dimension_members (
            dimension_id,
            app_id,
            tenant_id,
            branch_id,
            name,
            path,
            parent_id,
            weight,
            attributes,
            formula,
            is_deleted
        )
        VALUES (
            v_account_dimension_id,
            v_app_id,
            v_tenant_id,
            v_recovery_branch_id,
            'COGS',
            'account_root.cogs'::ltree,
            v_account_root_id,
            20,
            jsonb_build_object('seed', v_seed_prefix, 'replay_window', '2024Q4', 'tombstone', TRUE),
            NULL,
            TRUE
        )
        ON CONFLICT (dimension_id, name, branch_id)
        DO UPDATE SET
            app_id = EXCLUDED.app_id,
            tenant_id = EXCLUDED.tenant_id,
            path = EXCLUDED.path,
            parent_id = EXCLUDED.parent_id,
            weight = EXCLUDED.weight,
            attributes = EXCLUDED.attributes,
            formula = EXCLUDED.formula,
            is_deleted = EXCLUDED.is_deleted
        RETURNING id INTO v_cogs_tombstone_id;

        INSERT INTO workflow_nodes (
            tenant_id,
            app_id,
            node_key,
            state,
            owner_user_id,
            approver_user_id,
            lock_mode,
            is_locked,
            metadata
        )
        VALUES (
            v_tenant_id,
            v_app_id,
            'phase4-replay-window',
            'draft',
            v_admin_user_id,
            v_admin_user_id,
            'editable',
            FALSE,
            jsonb_build_object('seed', v_seed_prefix, 'window', '2024Q4')
        )
        ON CONFLICT (app_id, node_key)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            owner_user_id = EXCLUDED.owner_user_id,
            approver_user_id = EXCLUDED.approver_user_id,
            state = 'draft',
            lock_mode = 'editable',
            is_locked = FALSE,
            metadata = EXCLUDED.metadata,
            updated_at = NOW()
        RETURNING id INTO v_replay_node_id;

        DELETE FROM workflow_state_transitions
        WHERE node_id = v_replay_node_id;

        UPDATE workflow_nodes
        SET state = 'draft',
            lock_mode = 'editable',
            is_locked = FALSE,
            updated_at = NOW()
        WHERE id = v_replay_node_id;

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_replay_node_id, 'draft', 'in_review', v_admin_user_id, 'Seeded replay verification review', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_replay_node_id, 'in_review', 'approved', v_admin_user_id, 'Seeded replay verification approval', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_replay_node_id, 'approved', 'published', v_admin_user_id, 'Seeded replay verification publication', jsonb_build_object('seed', v_seed_prefix));

        DELETE FROM metadata_promotion_requests
        WHERE app_id = v_app_id
          AND summary = format('%s replay recovery promotion for %s', initcap(v_seed_prefix), v_tenant_label);

        INSERT INTO metadata_promotion_requests (
            tenant_id,
            app_id,
            source_branch_id,
            target_branch_name,
            requested_by_user_id,
            approved_by_user_id,
            status,
            risk_level,
            summary,
            diff_fingerprint,
            metadata,
            requested_at,
            approved_at,
            updated_at
        )
        VALUES (
            v_tenant_id,
            v_app_id,
            v_recovery_branch_id,
            'main',
            v_admin_user_id,
            v_admin_user_id,
            'applied',
            'critical',
            format('%s replay recovery promotion for %s', initcap(v_seed_prefix), v_tenant_label),
            md5(format('%s:%s:replay-recovery', v_seed_prefix, v_tenant_label)),
            jsonb_build_object('seed', v_seed_prefix, 'window', '2024Q4', 'source_member_id', v_revenue_override_id),
            NOW() - INTERVAL '15 minutes',
            NOW() - INTERVAL '10 minutes',
            NOW()
        );

        DELETE FROM connector_ingest_batches
        WHERE app_id = v_app_id
          AND connector_name = format('%s-replay', v_seed_prefix);

        INSERT INTO connector_ingest_batches (
            tenant_id,
            app_id,
            connector_name,
            source_uri,
            ingest_status,
            received_at,
            validated_at,
            applied_at,
            records_received,
            records_applied,
            records_rejected,
            metadata,
            updated_at
        )
        VALUES (
            v_tenant_id,
            v_app_id,
            format('%s-replay', v_seed_prefix),
            format('s3://quantatomai-phase4/%s/%s/replay/window-2024Q4.ndjson', v_seed_prefix, v_tenant_label),
            'quarantined',
            NOW() - INTERVAL '12 minutes',
            NOW() - INTERVAL '11 minutes',
            NULL,
            18000,
            0,
            18000,
            jsonb_build_object('seed', v_seed_prefix, 'window', '2024Q4', 'mode', 'replay-drill'),
            NOW()
        )
        RETURNING id INTO v_quarantined_batch_id;

        INSERT INTO connector_ingest_rejections (batch_id, record_key, reason_code, reason_detail, raw_record)
        VALUES (
            v_quarantined_batch_id,
            format('%s-%s-replay-window-2024Q4', v_seed_prefix, v_tenant_label),
            'quarantined-replay-window',
            'Replay batch intentionally quarantined for operator recovery drill.',
            jsonb_build_object('branch', 'phase4_recovery', 'member', 'Revenue', 'tombstone_member_id', v_cogs_tombstone_id)
        );

        DELETE FROM metadata_audit_events
        WHERE entity_type = 'replay-window'
          AND entity_id = v_recovery_branch_id;

        INSERT INTO metadata_audit_events (
            tenant_id,
            app_id,
            entity_type,
            entity_id,
            operation,
            actor_user_id,
            source_channel,
            trace_id,
            old_data,
            new_data,
            metadata,
            occurred_at,
            created_at
        )
        VALUES (
            v_tenant_id,
            v_app_id,
            'replay-window',
            v_recovery_branch_id,
            'INSERT',
            v_admin_user_id,
            'migration',
            format('%s-%s-replay-2024Q4', v_seed_prefix, v_tenant_label),
            NULL,
            jsonb_build_object('branch', 'phase4_recovery', 'window', '2024Q4', 'recovered_member_id', v_revenue_override_id),
            jsonb_build_object('seed', v_seed_prefix, 'drill', 'replay-recovery'),
            NOW() - INTERVAL '12 minutes',
            NOW()
        );
    END LOOP;
END $$;