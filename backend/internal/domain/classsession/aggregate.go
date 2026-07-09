package classsession

import (
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type ClassSession struct {
	id        ClassSessionID
	gymID     gym.GymID
	wodID     *domainwod.WODID
	startTime time.Time
	endTime   time.Time
	capacity  Capacity
	coachID   gym.MembershipID
	createdAt time.Time
	updatedAt time.Time
}

func NewClassSession(
	id ClassSessionID,
	gymID gym.GymID,
	wodID *domainwod.WODID,
	startTime time.Time,
	endTime time.Time,
	capacity Capacity,
	coachID gym.MembershipID,
	now time.Time,
) (ClassSession, error) {
	if coachID == "" {
		return ClassSession{}, ErrCoachRequired
	}
	if err := validateTimeRange(startTime, endTime); err != nil {
		return ClassSession{}, err
	}
	if err := validateCapacity(capacity); err != nil {
		return ClassSession{}, err
	}

	return ClassSession{
		id:        id,
		gymID:     gymID,
		wodID:     cloneWODID(wodID),
		startTime: startTime,
		endTime:   endTime,
		capacity:  capacity,
		coachID:   coachID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructClassSession(
	id ClassSessionID,
	gymID gym.GymID,
	wodID *domainwod.WODID,
	startTime time.Time,
	endTime time.Time,
	capacity Capacity,
	coachID gym.MembershipID,
	createdAt time.Time,
	updatedAt time.Time,
) (ClassSession, error) {
	session, err := NewClassSession(id, gymID, wodID, startTime, endTime, capacity, coachID, createdAt)
	if err != nil {
		return ClassSession{}, err
	}
	session.updatedAt = updatedAt
	return session, nil
}

func validateTimeRange(startTime, endTime time.Time) error {
	if !endTime.After(startTime) {
		return ErrInvalidTimeRange
	}
	return nil
}

func cloneWODID(id *domainwod.WODID) *domainwod.WODID {
	if id == nil {
		return nil
	}
	value := *id
	return &value
}

func (s ClassSession) ID() ClassSessionID {
	return s.id
}

func (s ClassSession) GymID() gym.GymID {
	return s.gymID
}

func (s ClassSession) WODID() *domainwod.WODID {
	return cloneWODID(s.wodID)
}

func (s ClassSession) StartTime() time.Time {
	return s.startTime
}

func (s ClassSession) EndTime() time.Time {
	return s.endTime
}

func (s ClassSession) Capacity() Capacity {
	return s.capacity
}

func (s ClassSession) CoachID() gym.MembershipID {
	return s.coachID
}

func (s ClassSession) CreatedAt() time.Time {
	return s.createdAt
}

func (s ClassSession) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *ClassSession) Reschedule(startTime, endTime time.Time, now time.Time) error {
	if err := validateTimeRange(startTime, endTime); err != nil {
		return err
	}
	s.startTime = startTime
	s.endTime = endTime
	s.updatedAt = now
	return nil
}

func (s *ClassSession) UpdateCapacity(capacity Capacity, now time.Time) error {
	if err := validateCapacity(capacity); err != nil {
		return err
	}
	s.capacity = capacity
	s.updatedAt = now
	return nil
}

func (s *ClassSession) AssignWOD(wodID domainwod.WODID, now time.Time) {
	value := wodID
	s.wodID = &value
	s.updatedAt = now
}

func (s *ClassSession) ClearWOD(now time.Time) {
	s.wodID = nil
	s.updatedAt = now
}
