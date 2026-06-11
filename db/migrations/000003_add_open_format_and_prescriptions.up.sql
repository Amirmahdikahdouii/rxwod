ALTER TYPE wod_type ADD VALUE IF NOT EXISTS 'OPEN';
ALTER TYPE scoring_kind ADD VALUE IF NOT EXISTS 'NONE';

ALTER TABLE wod_stages ADD COLUMN instructions TEXT;

ALTER TABLE wod_movements ADD COLUMN label TEXT;
ALTER TABLE wod_movements ADD COLUMN prescription TEXT;

ALTER TABLE wod_movements ALTER COLUMN name DROP NOT NULL;

ALTER TABLE wod_movements ADD CONSTRAINT wod_movements_name_or_prescription
  CHECK (
    char_length(trim(coalesce(name, ''))) > 0
    OR char_length(trim(coalesce(prescription, ''))) > 0
  );
