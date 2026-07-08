DROP INDEX IF EXISTS idx_gym_invitations_token_hash;

ALTER TABLE gym_invitations DROP COLUMN IF EXISTS token_hash;
