CREATE TYPE class_booking_status AS ENUM ('BOOKED', 'CANCELLED', 'ATTENDED');

CREATE TABLE class_sessions (
  id UUID PRIMARY KEY,
  gym_id UUID NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
  wod_id UUID REFERENCES wods(id) ON DELETE SET NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  capacity INTEGER NOT NULL CHECK (capacity > 0),
  coach_id UUID NOT NULL REFERENCES gym_memberships(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (end_time > start_time)
);

CREATE TABLE class_bookings (
  id UUID PRIMARY KEY,
  session_id UUID NOT NULL REFERENCES class_sessions(id) ON DELETE CASCADE,
  gym_membership_id UUID NOT NULL REFERENCES gym_memberships(id) ON DELETE CASCADE,
  status class_booking_status NOT NULL DEFAULT 'BOOKED',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (session_id, gym_membership_id)
);

CREATE INDEX idx_class_sessions_gym_start_time ON class_sessions (gym_id, start_time);
CREATE INDEX idx_class_sessions_wod_id ON class_sessions (wod_id);
CREATE INDEX idx_class_sessions_coach_id ON class_sessions (coach_id);

CREATE INDEX idx_class_bookings_session_id ON class_bookings (session_id);
CREATE INDEX idx_class_bookings_gym_membership_id ON class_bookings (gym_membership_id);
CREATE INDEX idx_class_bookings_session_booked ON class_bookings (session_id) WHERE status = 'BOOKED';
