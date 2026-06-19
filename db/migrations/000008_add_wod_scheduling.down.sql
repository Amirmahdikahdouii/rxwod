DROP INDEX IF EXISTS idx_wods_gym_status_scheduled_date;
DROP INDEX IF EXISTS idx_wods_gym_scheduled_date;

ALTER TABLE wods DROP COLUMN IF EXISTS published_at;
ALTER TABLE wods DROP COLUMN IF EXISTS scheduled_date;
