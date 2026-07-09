CREATE TABLE athlete_default_sessions (
  id UUID PRIMARY KEY,
  gym_membership_id UUID NOT NULL REFERENCES gym_memberships(id) ON DELETE CASCADE,
  day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
  time_slot TIME NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (gym_membership_id, day_of_week, time_slot)
);

CREATE INDEX idx_athlete_default_sessions_membership ON athlete_default_sessions (gym_membership_id);
CREATE INDEX idx_athlete_default_sessions_day_time ON athlete_default_sessions (day_of_week, time_slot);
