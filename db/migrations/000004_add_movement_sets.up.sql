ALTER TABLE wod_movements
ADD COLUMN sets INTEGER CHECK (sets IS NULL OR sets > 0);
