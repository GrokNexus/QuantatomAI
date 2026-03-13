-- Add sort_order to compatibility dimensions for stable ordering
ALTER TABLE dimensions_compat
    ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;
