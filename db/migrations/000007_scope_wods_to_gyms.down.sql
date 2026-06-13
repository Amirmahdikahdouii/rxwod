DROP INDEX IF EXISTS idx_wods_gym_created_at;
DROP INDEX IF EXISTS idx_wods_gym_id;

ALTER TABLE wods
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS gym_id;
