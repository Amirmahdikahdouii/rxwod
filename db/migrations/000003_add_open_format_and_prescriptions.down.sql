ALTER TABLE wod_movements DROP CONSTRAINT IF EXISTS wod_movements_name_or_prescription;

ALTER TABLE wod_movements DROP COLUMN IF EXISTS prescription;
ALTER TABLE wod_movements DROP COLUMN IF EXISTS label;

ALTER TABLE wod_movements ALTER COLUMN name SET NOT NULL;

ALTER TABLE wod_stages DROP COLUMN IF EXISTS instructions;

-- PostgreSQL does not support removing enum values; OPEN and NONE remain in wod_type/scoring_kind.
