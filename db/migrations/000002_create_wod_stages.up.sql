DROP INDEX IF EXISTS idx_wods_config_gin;
DROP INDEX IF EXISTS idx_wods_status;
DROP INDEX IF EXISTS idx_wods_type;
DROP TABLE IF EXISTS wod_movements;
DROP TABLE IF EXISTS wods;

CREATE TYPE stage_kind AS ENUM ('WARMUP', 'STRENGTH', 'CORE', 'METCON', 'COOLDOWN');

CREATE TABLE wods (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL CHECK (char_length(name) BETWEEN 3 AND 120),
  status wod_status NOT NULL DEFAULT 'DRAFT',
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE wod_stages (
  id UUID PRIMARY KEY,
  wod_id UUID NOT NULL REFERENCES wods(id) ON DELETE CASCADE,
  position INTEGER NOT NULL CHECK (position > 0),
  stage_kind stage_kind NOT NULL,
  wod_type wod_type NOT NULL,
  config JSONB NOT NULL,
  scoring_kind scoring_kind NOT NULL,
  UNIQUE (wod_id, position),
  CHECK (jsonb_typeof(config) = 'object')
);

CREATE TABLE wod_movements (
  id UUID PRIMARY KEY,
  stage_id UUID NOT NULL REFERENCES wod_stages(id) ON DELETE CASCADE,
  position INTEGER NOT NULL CHECK (position > 0),
  name TEXT NOT NULL,
  reps INTEGER CHECK (reps IS NULL OR reps > 0),
  load_value NUMERIC(8,2),
  load_unit TEXT CHECK (load_unit IS NULL OR load_unit IN ('kg', 'lb', 'bodyweight')),
  notes TEXT,
  UNIQUE (stage_id, position)
);

CREATE INDEX idx_wods_status ON wods(status);
CREATE INDEX idx_wod_stages_wod_id ON wod_stages(wod_id);
CREATE INDEX idx_wod_stages_kind ON wod_stages(stage_kind);
CREATE INDEX idx_wod_stages_type ON wod_stages(wod_type);
CREATE INDEX idx_wod_stages_config_gin ON wod_stages USING GIN (config jsonb_path_ops);
CREATE INDEX idx_wod_movements_stage_id ON wod_movements(stage_id);
