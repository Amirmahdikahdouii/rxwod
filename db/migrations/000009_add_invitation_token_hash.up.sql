ALTER TABLE gym_invitations ADD COLUMN token_hash TEXT;

CREATE UNIQUE INDEX idx_gym_invitations_token_hash
  ON gym_invitations (token_hash)
  WHERE token_hash IS NOT NULL;
