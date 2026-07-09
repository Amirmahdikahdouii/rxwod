CREATE TABLE wod_results (
  id UUID PRIMARY KEY,
  wod_id UUID NOT NULL REFERENCES wods(id) ON DELETE CASCADE,
  gym_membership_id UUID NOT NULL REFERENCES gym_memberships(id) ON DELETE CASCADE,
  score_value INTEGER NOT NULL CHECK (score_value >= 0),
  is_rx BOOLEAN NOT NULL DEFAULT true,
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (wod_id, gym_membership_id)
);

CREATE INDEX idx_wod_results_wod_id ON wod_results (wod_id);
CREATE INDEX idx_wod_results_gym_membership_id ON wod_results (gym_membership_id);
