-- Phase 4 tenant-aware benchmark fixtures.
-- Optional runtime setting: phase4.seed_prefix (defaults to 'phase4').
DO $$
DECLARE
    v_seed_prefix TEXT := COALESCE(NULLIF(current_setting('phase4.seed_prefix', true), ''), 'phase4');
    v_tenant_labels TEXT[] := ARRAY['alpha', 'beta', 'gamma'];
    v_tenant_label TEXT;
    v_tenant_name TEXT;
    v_app_name TEXT;
    v_primary_region TEXT;
    v_secondary_region TEXT;
    v_tenant_id UUID;
    v_admin_user_id UUID;
    v_planner_user_id UUID;
    v_viewer_user_id UUID;
    v_app_id UUID;
    v_main_branch_id UUID;
    v_sandbox_branch_id UUID;
    v_account_dimension_id UUID;
    v_region_dimension_id UUID;
    v_scenario_dimension_id UUID;
    v_account_root_id UUID;
    v_region_root_id UUID;
    v_scenario_root_id UUID;
    v_revenue_member_id UUID;
    v_cogs_member_id UUID;
    v_workflow_publish_id UUID;
    v_workflow_reject_id UUID;
    v_validated_batch_id UUID;
    v_partial_batch_id UUID;
BEGIN
    FOREACH v_tenant_label IN ARRAY v_tenant_labels LOOP
        v_tenant_name := format('%s-tenant-%s', v_seed_prefix, v_tenant_label);
        v_app_name := format('%s-planning-%s', v_seed_prefix, v_tenant_label);
        v_primary_region := CASE v_tenant_label
            WHEN 'alpha' THEN 'us-east-1'
            WHEN 'beta' THEN 'eu-west-1'
            ELSE 'ap-southeast-1'
        END;
        v_secondary_region := CASE v_tenant_label
            WHEN 'alpha' THEN 'us-west-2'
            WHEN 'beta' THEN 'eu-central-1'
            ELSE 'ap-northeast-1'
        END;

        SELECT id INTO v_tenant_id
        FROM tenants
        WHERE name = v_tenant_name
        ORDER BY created_at
        LIMIT 1;

        IF v_tenant_id IS NULL THEN
            INSERT INTO tenants (
                name,
                plan_tier,
                status,
                residency_mode,
                primary_region,
                isolation_tier,
                ai_learning_mode,
                cost_center_code,
                created_at,
                updated_at
            )
            VALUES (
                v_tenant_name,
                CASE WHEN v_tenant_label = 'gamma' THEN 'ultra' ELSE 'enterprise' END,
                'active',
                CASE WHEN v_tenant_label = 'beta' THEN 'geo-fenced' ELSE 'single-region' END,
                v_primary_region,
                CASE WHEN v_tenant_label = 'gamma' THEN 'dedicated-data-plane' ELSE 'logical' END,
                'tenant-only',
                upper(replace(v_tenant_name, '-', '_')),
                NOW(),
                NOW()
            )
            RETURNING id INTO v_tenant_id;
        ELSE
            UPDATE tenants
            SET plan_tier = CASE WHEN v_tenant_label = 'gamma' THEN 'ultra' ELSE 'enterprise' END,
                status = 'active',
                residency_mode = CASE WHEN v_tenant_label = 'beta' THEN 'geo-fenced' ELSE 'single-region' END,
                primary_region = v_primary_region,
                isolation_tier = CASE WHEN v_tenant_label = 'gamma' THEN 'dedicated-data-plane' ELSE 'logical' END,
                ai_learning_mode = 'tenant-only',
                cost_center_code = upper(replace(v_tenant_name, '-', '_')),
                updated_at = NOW()
            WHERE id = v_tenant_id;
        END IF;

        INSERT INTO tenant_regions (
            tenant_id,
            region_code,
            region_role,
            residency_class,
            is_write_region,
            is_read_region,
            is_failover_region,
            metadata
        )
        VALUES (
            v_tenant_id,
            v_primary_region,
            'primary',
            'customer-data',
            TRUE,
            TRUE,
            FALSE,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label)
        )
        ON CONFLICT (tenant_id, region_code)
        DO UPDATE SET
            region_role = EXCLUDED.region_role,
            residency_class = EXCLUDED.residency_class,
            is_write_region = EXCLUDED.is_write_region,
            is_read_region = EXCLUDED.is_read_region,
            is_failover_region = EXCLUDED.is_failover_region,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO tenant_regions (
            tenant_id,
            region_code,
            region_role,
            residency_class,
            is_write_region,
            is_read_region,
            is_failover_region,
            metadata
        )
        VALUES (
            v_tenant_id,
            v_secondary_region,
            'failover',
            'audit',
            FALSE,
            TRUE,
            TRUE,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label, 'purpose', 'failover')
        )
        ON CONFLICT (tenant_id, region_code)
        DO UPDATE SET
            region_role = EXCLUDED.region_role,
            residency_class = EXCLUDED.residency_class,
            is_write_region = EXCLUDED.is_write_region,
            is_read_region = EXCLUDED.is_read_region,
            is_failover_region = EXCLUDED.is_failover_region,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO tenant_key_domains (
            tenant_id,
            region_code,
            purpose,
            kms_provider,
            key_uri,
            rotation_interval_days,
            is_active,
            metadata
        )
        VALUES (
            v_tenant_id,
            v_primary_region,
            'app-data',
            'aws-kms',
            format('kms://%s/%s/app-data', v_primary_region, v_tenant_name),
            45,
            TRUE,
            jsonb_build_object('seed', v_seed_prefix)
        )
        ON CONFLICT (tenant_id, region_code, purpose)
        DO UPDATE SET
            kms_provider = EXCLUDED.kms_provider,
            key_uri = EXCLUDED.key_uri,
            rotation_interval_days = EXCLUDED.rotation_interval_days,
            is_active = EXCLUDED.is_active,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO tenant_key_domains (
            tenant_id,
            region_code,
            purpose,
            kms_provider,
            key_uri,
            rotation_interval_days,
            is_active,
            metadata
        )
        VALUES (
            v_tenant_id,
            v_secondary_region,
            'audit',
            'aws-kms',
            format('kms://%s/%s/audit', v_secondary_region, v_tenant_name),
            30,
            TRUE,
            jsonb_build_object('seed', v_seed_prefix)
        )
        ON CONFLICT (tenant_id, region_code, purpose)
        DO UPDATE SET
            kms_provider = EXCLUDED.kms_provider,
            key_uri = EXCLUDED.key_uri,
            rotation_interval_days = EXCLUDED.rotation_interval_days,
            is_active = EXCLUDED.is_active,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO tenant_quota_policies (
            tenant_id,
            max_users,
            max_apps,
            max_storage_gb,
            max_hot_working_set_gb,
            max_events_per_sec,
            max_api_rps,
            max_concurrent_jobs,
            max_vector_bytes,
            chargeback_model,
            monthly_spend_limit_usd,
            overage_behavior,
            metadata
        )
        VALUES (
            v_tenant_id,
            200,
            24,
            4096,
            256,
            15000,
            1200,
            48,
            2147483648,
            'showback',
            25000.00,
            'throttle',
            jsonb_build_object('seed', v_seed_prefix, 'profile', 'phase4')
        )
        ON CONFLICT (tenant_id)
        DO UPDATE SET
            max_users = EXCLUDED.max_users,
            max_apps = EXCLUDED.max_apps,
            max_storage_gb = EXCLUDED.max_storage_gb,
            max_hot_working_set_gb = EXCLUDED.max_hot_working_set_gb,
            max_events_per_sec = EXCLUDED.max_events_per_sec,
            max_api_rps = EXCLUDED.max_api_rps,
            max_concurrent_jobs = EXCLUDED.max_concurrent_jobs,
            max_vector_bytes = EXCLUDED.max_vector_bytes,
            chargeback_model = EXCLUDED.chargeback_model,
            monthly_spend_limit_usd = EXCLUDED.monthly_spend_limit_usd,
            overage_behavior = EXCLUDED.overage_behavior,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO tenant_ai_policies (
            tenant_id,
            retrieval_scope,
            allow_cross_tenant_learning,
            allow_external_inference,
            require_prompt_audit,
            require_human_approval_for_generative_write,
            max_context_rows,
            vector_namespace_strategy,
            metadata
        )
        VALUES (
            v_tenant_id,
            'tenant-only',
            FALSE,
            FALSE,
            TRUE,
            TRUE,
            2500,
            'tenant-segregated',
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label)
        )
        ON CONFLICT (tenant_id)
        DO UPDATE SET
            retrieval_scope = EXCLUDED.retrieval_scope,
            allow_cross_tenant_learning = EXCLUDED.allow_cross_tenant_learning,
            allow_external_inference = EXCLUDED.allow_external_inference,
            require_prompt_audit = EXCLUDED.require_prompt_audit,
            require_human_approval_for_generative_write = EXCLUDED.require_human_approval_for_generative_write,
            max_context_rows = EXCLUDED.max_context_rows,
            vector_namespace_strategy = EXCLUDED.vector_namespace_strategy,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO users (tenant_id, email, password_hash, role)
        VALUES (v_tenant_id, format('admin+%s@phase4.quantatom.ai', v_tenant_label), 'phase4-admin-hash', 'admin')
        ON CONFLICT (tenant_id, email)
        DO UPDATE SET
            password_hash = EXCLUDED.password_hash,
            role = EXCLUDED.role;

        INSERT INTO users (tenant_id, email, password_hash, role)
        VALUES (v_tenant_id, format('planner+%s@phase4.quantatom.ai', v_tenant_label), 'phase4-planner-hash', 'planner')
        ON CONFLICT (tenant_id, email)
        DO UPDATE SET
            password_hash = EXCLUDED.password_hash,
            role = EXCLUDED.role;

        INSERT INTO users (tenant_id, email, password_hash, role)
        VALUES (v_tenant_id, format('viewer+%s@phase4.quantatom.ai', v_tenant_label), 'phase4-viewer-hash', 'viewer')
        ON CONFLICT (tenant_id, email)
        DO UPDATE SET
            password_hash = EXCLUDED.password_hash,
            role = EXCLUDED.role;

        SELECT id INTO v_admin_user_id FROM users WHERE tenant_id = v_tenant_id AND email = format('admin+%s@phase4.quantatom.ai', v_tenant_label);
        SELECT id INTO v_planner_user_id FROM users WHERE tenant_id = v_tenant_id AND email = format('planner+%s@phase4.quantatom.ai', v_tenant_label);
        SELECT id INTO v_viewer_user_id FROM users WHERE tenant_id = v_tenant_id AND email = format('viewer+%s@phase4.quantatom.ai', v_tenant_label);

        INSERT INTO apps (
            tenant_id,
            name,
            description,
            planning_mode,
            default_currency,
            created_by
        )
        VALUES (
            v_tenant_id,
            v_app_name,
            format('Phase 4 governed planning workspace for %s.', v_tenant_label),
            'scenario_based',
            'USD',
            v_admin_user_id
        )
        ON CONFLICT (tenant_id, name)
        DO UPDATE SET
            description = EXCLUDED.description,
            planning_mode = EXCLUDED.planning_mode,
            default_currency = EXCLUDED.default_currency,
            created_by = EXCLUDED.created_by,
            updated_at = NOW()
        RETURNING id INTO v_app_id;

        INSERT INTO branches (app_id, tenant_id, name, base_branch_id)
        VALUES (v_app_id, v_tenant_id, 'main', NULL)
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            base_branch_id = EXCLUDED.base_branch_id,
            updated_at = NOW()
        RETURNING id INTO v_main_branch_id;

        INSERT INTO branches (app_id, tenant_id, name, base_branch_id)
        VALUES (v_app_id, v_tenant_id, 'phase4_sandbox', v_main_branch_id)
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            base_branch_id = EXCLUDED.base_branch_id,
            updated_at = NOW()
        RETURNING id INTO v_sandbox_branch_id;

        INSERT INTO app_partitions (
            app_id,
            tenant_id,
            write_region,
            hot_namespace,
            warm_partition_template,
            cold_object_prefix,
            event_topic_prefix,
            cache_namespace,
            metadata
        )
        VALUES (
            v_app_id,
            v_tenant_id,
            v_primary_region,
            format('%s.%s.hot', v_seed_prefix, v_tenant_label),
            format('warm_%s_%s_{yyyy_mm}', v_seed_prefix, v_tenant_label),
            format('s3://quantatomai-phase4/%s/%s', v_seed_prefix, v_tenant_label),
            format('%s.%s.events', v_seed_prefix, v_tenant_label),
            format('%s.%s.cache', v_seed_prefix, v_tenant_label),
            jsonb_build_object('seed', v_seed_prefix, 'app', v_app_name)
        )
        ON CONFLICT (app_id)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            write_region = EXCLUDED.write_region,
            hot_namespace = EXCLUDED.hot_namespace,
            warm_partition_template = EXCLUDED.warm_partition_template,
            cold_object_prefix = EXCLUDED.cold_object_prefix,
            event_topic_prefix = EXCLUDED.event_topic_prefix,
            cache_namespace = EXCLUDED.cache_namespace,
            metadata = EXCLUDED.metadata,
            updated_at = NOW();

        INSERT INTO dimensions (app_id, tenant_id, name, type, is_core, sort_order, properties)
        VALUES (v_app_id, v_tenant_id, 'Account', 'standard', TRUE, 10, jsonb_build_object('seed', v_seed_prefix, 'axis', 'account'))
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            type = EXCLUDED.type,
            is_core = EXCLUDED.is_core,
            sort_order = EXCLUDED.sort_order,
            properties = EXCLUDED.properties
        RETURNING id INTO v_account_dimension_id;

        INSERT INTO dimensions (app_id, tenant_id, name, type, is_core, sort_order, properties)
        VALUES (v_app_id, v_tenant_id, 'Region', 'standard', TRUE, 20, jsonb_build_object('seed', v_seed_prefix, 'axis', 'region'))
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            type = EXCLUDED.type,
            is_core = EXCLUDED.is_core,
            sort_order = EXCLUDED.sort_order,
            properties = EXCLUDED.properties
        RETURNING id INTO v_region_dimension_id;

        INSERT INTO dimensions (app_id, tenant_id, name, type, is_core, sort_order, properties)
        VALUES (v_app_id, v_tenant_id, 'Scenario', 'scenario', TRUE, 30, jsonb_build_object('seed', v_seed_prefix, 'axis', 'scenario'))
        ON CONFLICT (app_id, name)
        DO UPDATE SET
            tenant_id = EXCLUDED.tenant_id,
            type = EXCLUDED.type,
            is_core = EXCLUDED.is_core,
            sort_order = EXCLUDED.sort_order,
            properties = EXCLUDED.properties
        RETURNING id INTO v_scenario_dimension_id;

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
            v_main_branch_id,
            'AccountRoot',
            'account_root'::ltree,
            NULL,
            0,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label),
            NULL,
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
        RETURNING id INTO v_account_root_id;

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
            v_main_branch_id,
            'Revenue',
            'account_root.revenue'::ltree,
            v_account_root_id,
            10,
            jsonb_build_object('seed', v_seed_prefix, 'classification', 'income'),
            'SUM(children)',
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
        RETURNING id INTO v_revenue_member_id;

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
            v_main_branch_id,
            'COGS',
            'account_root.cogs'::ltree,
            v_account_root_id,
            20,
            jsonb_build_object('seed', v_seed_prefix, 'classification', 'expense'),
            'SUM(children)',
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
        RETURNING id INTO v_cogs_member_id;

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
            v_region_dimension_id,
            v_app_id,
            v_tenant_id,
            v_main_branch_id,
            'RegionRoot',
            'region_root'::ltree,
            NULL,
            0,
            jsonb_build_object('seed', v_seed_prefix),
            NULL,
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
        RETURNING id INTO v_region_root_id;

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
            v_region_dimension_id,
            v_app_id,
            v_tenant_id,
            v_main_branch_id,
            'PrimaryRegion',
            'region_root.primary_region'::ltree,
            v_region_root_id,
            10,
            jsonb_build_object('seed', v_seed_prefix, 'region_code', v_primary_region),
            NULL,
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
            is_deleted = EXCLUDED.is_deleted;

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
            v_scenario_dimension_id,
            v_app_id,
            v_tenant_id,
            v_main_branch_id,
            'ScenarioRoot',
            'scenario_root'::ltree,
            NULL,
            0,
            jsonb_build_object('seed', v_seed_prefix),
            NULL,
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
        RETURNING id INTO v_scenario_root_id;

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
            v_scenario_dimension_id,
            v_app_id,
            v_tenant_id,
            v_main_branch_id,
            'Forecast',
            'scenario_root.forecast'::ltree,
            v_scenario_root_id,
            10,
            jsonb_build_object('seed', v_seed_prefix, 'version', 'baseline'),
            NULL,
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
            is_deleted = EXCLUDED.is_deleted;

        DELETE FROM security_policies
        WHERE app_id = v_app_id
          AND name IN ('Phase4 Planner Policy', 'Phase4 Viewer Policy');

        INSERT INTO security_policies (app_id, tenant_id, name, rules, user_id, permission_level)
        VALUES (
            v_app_id,
            v_tenant_id,
            'Phase4 Planner Policy',
            jsonb_build_object('Region', jsonb_build_array('PrimaryRegion'), 'Scenario', jsonb_build_array('Forecast')),
            v_planner_user_id,
            'write'
        );

        INSERT INTO security_policies (app_id, tenant_id, name, rules, user_id, permission_level)
        VALUES (
            v_app_id,
            v_tenant_id,
            'Phase4 Viewer Policy',
            jsonb_build_object('Region', jsonb_build_array('PrimaryRegion')),
            v_viewer_user_id,
            'read'
        );

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
            'phase4-budget-publish',
            'draft',
            v_planner_user_id,
            v_admin_user_id,
            'editable',
            FALSE,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label)
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
        RETURNING id INTO v_workflow_publish_id;

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
            'phase4-budget-reject',
            'draft',
            v_planner_user_id,
            v_admin_user_id,
            'editable',
            FALSE,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label)
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
        RETURNING id INTO v_workflow_reject_id;

        DELETE FROM workflow_state_transitions
        WHERE node_id IN (v_workflow_publish_id, v_workflow_reject_id);

        UPDATE workflow_nodes
        SET state = 'draft',
            lock_mode = 'editable',
            is_locked = FALSE,
            updated_at = NOW()
        WHERE id IN (v_workflow_publish_id, v_workflow_reject_id);

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_publish_id, 'draft', 'in_review', v_planner_user_id, 'Phase 4 seeded review flow', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_publish_id, 'in_review', 'approved', v_admin_user_id, 'Phase 4 seeded approval flow', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_publish_id, 'approved', 'published', v_admin_user_id, 'Phase 4 seeded publish flow', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_reject_id, 'draft', 'in_review', v_planner_user_id, 'Phase 4 seeded rejection flow', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_reject_id, 'in_review', 'rejected', v_admin_user_id, 'Phase 4 seeded rejection flow', jsonb_build_object('seed', v_seed_prefix));

        INSERT INTO workflow_state_transitions (node_id, from_state, to_state, actor_user_id, reason, metadata)
        VALUES (v_workflow_reject_id, 'rejected', 'draft', v_planner_user_id, 'Phase 4 seeded resubmission flow', jsonb_build_object('seed', v_seed_prefix));

        DELETE FROM metadata_promotion_requests
        WHERE app_id = v_app_id
          AND summary LIKE format('%s%%', initcap(v_seed_prefix));

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
            v_sandbox_branch_id,
            'main',
            v_planner_user_id,
            v_admin_user_id,
            'approved',
            'medium',
            format('%s baseline promotion for %s', initcap(v_seed_prefix), v_tenant_label),
            md5(format('%s:%s:baseline', v_seed_prefix, v_tenant_label)),
            jsonb_build_object('seed', v_seed_prefix, 'kind', 'baseline-promotion'),
            NOW() - INTERVAL '2 hours',
            NOW() - INTERVAL '90 minutes',
            NOW()
        );

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
            v_sandbox_branch_id,
            'main',
            v_planner_user_id,
            NULL,
            'pending',
            'high',
            format('%s pending promotion for %s', initcap(v_seed_prefix), v_tenant_label),
            md5(format('%s:%s:pending', v_seed_prefix, v_tenant_label)),
            jsonb_build_object('seed', v_seed_prefix, 'kind', 'pending-promotion'),
            NOW() - INTERVAL '20 minutes',
            NULL,
            NOW()
        );

        DELETE FROM connector_ingest_batches
        WHERE app_id = v_app_id
          AND connector_name LIKE format('%s-%%', v_seed_prefix);

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
            format('%s-erp', v_seed_prefix),
            format('s3://quantatomai-phase4/%s/%s/erp/batch-001.csv', v_seed_prefix, v_tenant_label),
            'applied',
            NOW() - INTERVAL '3 hours',
            NOW() - INTERVAL '170 minutes',
            NOW() - INTERVAL '160 minutes',
            120000,
            119600,
            400,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label, 'profile', 'C'),
            NOW()
        )
        RETURNING id INTO v_validated_batch_id;

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
            format('%s-recon', v_seed_prefix),
            format('s3://quantatomai-phase4/%s/%s/recon/batch-002.csv', v_seed_prefix, v_tenant_label),
            'partial',
            NOW() - INTERVAL '50 minutes',
            NOW() - INTERVAL '40 minutes',
            NOW() - INTERVAL '35 minutes',
            40000,
            39750,
            250,
            jsonb_build_object('seed', v_seed_prefix, 'tenant', v_tenant_label, 'profile', 'C'),
            NOW()
        )
        RETURNING id INTO v_partial_batch_id;

        INSERT INTO connector_ingest_rejections (batch_id, record_key, reason_code, reason_detail, raw_record)
        VALUES (
            v_validated_batch_id,
            format('%s-%s-%s', v_seed_prefix, v_tenant_label, 'erp-400'),
            'dimension-mismatch',
            'Record references a non-governed account member.',
            jsonb_build_object('account', 'UNKNOWN', 'amount', 100.25)
        );

        INSERT INTO connector_ingest_rejections (batch_id, record_key, reason_code, reason_detail, raw_record)
        VALUES (
            v_partial_batch_id,
            format('%s-%s-%s', v_seed_prefix, v_tenant_label, 'recon-250'),
            'workflow-locked',
            'Rejected because the target workflow node is published and locked.',
            jsonb_build_object('workflowNode', 'phase4-budget-publish', 'status', 'locked')
        );
    END LOOP;
END $$;