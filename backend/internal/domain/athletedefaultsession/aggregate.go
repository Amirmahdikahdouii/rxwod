package athletedefaultsession

import (
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
)

type AthleteDefaultSession struct {
	id              AthleteDefaultSessionID
	gymMembershipID gym.MembershipID
	dayOfWeek       DayOfWeek
	timeSlot        TimeSlot
	createdAt       time.Time
	updatedAt       time.Time
}

func NewAthleteDefaultSession(
	id AthleteDefaultSessionID,
	gymMembershipID gym.MembershipID,
	dayOfWeek DayOfWeek,
	timeSlot TimeSlot,
	now time.Time,
) (AthleteDefaultSession, error) {
	if gymMembershipID == "" {
		return AthleteDefaultSession{}, ErrMembershipRequired
	}
	if err := validateDayOfWeek(dayOfWeek); err != nil {
		return AthleteDefaultSession{}, err
	}
	if err := validateTimeSlot(timeSlot); err != nil {
		return AthleteDefaultSession{}, err
	}

	return AthleteDefaultSession{
		id:              id,
		gymMembershipID: gymMembershipID,
		dayOfWeek:       dayOfWeek,
		timeSlot:        timeSlot,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructAthleteDefaultSession(
	id AthleteDefaultSessionID,
	gymMembershipID gym.MembershipID,
	dayOfWeek DayOfWeek,
	timeSlot TimeSlot,
	createdAt time.Time,
	updatedAt time.Time,
) (AthleteDefaultSession, error) {
	session, err := NewAthleteDefaultSession(id, gymMembershipID, dayOfWeek, timeSlot, createdAt)
	if err != nil {
		return AthleteDefaultSession{}, err
	}
	session.updatedAt = updatedAt
	return session, nil
}

func (s AthleteDefaultSession) ID() AthleteDefaultSessionID {
	return s.id
}

func (s AthleteDefaultSession) GymMembershipID() gym.MembershipID {
	return s.gymMembershipID
}

func (s AthleteDefaultSession) DayOfWeek() DayOfWeek {
	return s.dayOfWeek
}

func (s AthleteDefaultSession) TimeSlot() TimeSlot {
	return s.timeSlot
}

func (s AthleteDefaultSession) CreatedAt() time.Time {
	return s.createdAt
}

func (s AthleteDefaultSession) UpdatedAt() time.Time {
	return s.updatedAt
}
