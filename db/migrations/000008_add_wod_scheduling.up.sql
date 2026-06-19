ALTER TABLE wods ADD COLUMN scheduled_date DATE;
ALTER TABLE wods ADD COLUMN published_at TIMESTAMPTZ;

CREATE INDEX idx_wods_gym_scheduled_date ON wods (gym_id, scheduled_date);
CREATE INDEX idx_wods_gym_status_scheduled_date ON wods (gym_id, status, scheduled_date);
