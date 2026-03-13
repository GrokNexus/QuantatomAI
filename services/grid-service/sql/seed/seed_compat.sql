-- Seed 15 dimensions and ~100,000 members for UI/grid discovery
-- Model is configurable via PG runtime setting gridseed.model_id (optional), defaults to 'default_model'.
DO $$
DECLARE
    v_model TEXT := COALESCE(current_setting('gridseed.model_id', true), 'default_model');
    v_dim_idx INT;
    v_dim_id BIGINT;
    v_target INT;
BEGIN
    FOR v_dim_idx IN 1..15 LOOP
        INSERT INTO dimensions_compat (model_id, name, sort_order, is_active)
        VALUES (v_model, 'Dim_' || v_dim_idx, v_dim_idx, TRUE)
        ON CONFLICT (model_id, name)
        DO UPDATE SET sort_order = EXCLUDED.sort_order, is_active = EXCLUDED.is_active
        RETURNING id INTO v_dim_id;

        IF NOT FOUND THEN
            SELECT id INTO v_dim_id FROM dimensions_compat WHERE model_id = v_model AND name = 'Dim_' || v_dim_idx;
        END IF;

        v_target := CASE WHEN v_dim_idx = 1 THEN 100000 - (15 - 1) * 6667 ELSE 6667 END;

        INSERT INTO members_compat (
            dimension_id,
            code,
            name,
            sequence,
            is_active,
            is_deleted,
            effective_start,
            path,
            parent_code
        )
        SELECT
            v_dim_id,
            format('M%s_%s', v_dim_idx, gs),
            format('Member %s %s', v_dim_idx, gs),
            gs,
            TRUE,
            FALSE,
            NOW(),
            format('Dim_%s.%s', v_dim_idx, gs),
            NULL
        FROM generate_series(1, v_target) AS gs;
    END LOOP;
END$$;
