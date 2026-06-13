CREATE TYPE gym_role AS ENUM ('owner', 'coach', 'athlete');
CREATE TYPE membership_status AS ENUM ('pending', 'active');
CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'revoked');

CREATE TABLE gyms (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL CHECK (char_length(name) BETWEEN 2 AND 120),
  owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE gym_memberships (
  id UUID PRIMARY KEY,
  gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role gym_role NOT NULL,
  status membership_status NOT NULL DEFAULT 'pending',
  invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (gym_id, user_id)
);

CREATE TABLE gym_invitations (
  id UUID PRIMARY KEY,
  gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
  email TEXT NOT NULL,
  role gym_role NOT NULL CHECK (role IN ('coach', 'athlete')),
  status invitation_status NOT NULL DEFAULT 'pending',
  invited_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (email = lower(trim(email)))
);

CREATE UNIQUE INDEX idx_gym_invitations_pending_unique
  ON gym_invitations(gym_id, email, role)
  WHERE status = 'pending';

CREATE INDEX idx_gyms_owner_id ON gyms(owner_id);
CREATE INDEX idx_gym_memberships_user_id ON gym_memberships(user_id);
CREATE INDEX idx_gym_memberships_gym_id ON gym_memberships(gym_id);
CREATE INDEX idx_gym_invitations_email_status ON gym_invitations(email, status);
