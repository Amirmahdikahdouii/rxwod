TRUNCATE TABLE wod_movements, wod_stages, wods CASCADE;

ALTER TABLE wods
ADD COLUMN gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT;

CREATE INDEX idx_wods_gym_id ON wods(gym_id);
CREATE INDEX idx_wods_gym_created_at ON wods(gym_id, created_at DESC);
