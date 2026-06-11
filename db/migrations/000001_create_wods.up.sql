CREATE TYPE wod_type AS ENUM ('AMRAP', 'FORTIME', 'TABATA', 'EMOM');
CREATE TYPE wod_status AS ENUM ('DRAFT', 'PUBLISHED', 'ARCHIVED');
CREATE TYPE scoring_kind AS ENUM ('ROUNDS_REPS', 'TIME_TO_COMPLETE', 'TOTAL_REPS');

CREATE TABLE wods (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL CHECK (char_length(name) BETWEEN 3 AND 120),
  wod_type wod_type NOT NULL,
  status wod_status NOT NULL DEFAULT 'DRAFT',
  description TEXT,
  config JSONB NOT NULL,
  scoring_kind scoring_kind NOT NULL,
  scoring_config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (jsonb_typeof(config) = 'object'),
  CHECK (jsonb_typeof(scoring_config) = 'object')
);

CREATE TABLE wod_movements (
  id UUID PRIMARY KEY,
  wod_id UUID NOT NULL REFERENCES wods(id) ON DELETE CASCADE,
  position INTEGER NOT NULL CHECK (position > 0),
  name TEXT NOT NULL,
  reps INTEGER CHECK (reps IS NULL OR reps > 0),
  load_value NUMERIC(8,2),
  load_unit TEXT CHECK (load_unit IS NULL OR load_unit IN ('kg', 'lb', 'bodyweight')),
  notes TEXT,
  UNIQUE (wod_id, position)
);

CREATE INDEX idx_wods_type ON wods(wod_type);
CREATE INDEX idx_wods_status ON wods(status);
CREATE INDEX idx_wods_config_gin ON wods USING GIN (config jsonb_path_ops);
