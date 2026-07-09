package classbooking

import (
	"time"

	"github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type ClassBooking struct {
	id              ClassBookingID
	sessionID       classsession.ClassSessionID
	gymMembershipID gym.MembershipID
	status          BookingStatus
	createdAt       time.Time
	updatedAt       time.Time
}

func NewClassBooking(
	id ClassBookingID,
	sessionID classsession.ClassSessionID,
	gymMembershipID gym.MembershipID,
	now time.Time,
) (ClassBooking, error) {
	return ClassBooking{
		id:              id,
		sessionID:       sessionID,
		gymMembershipID: gymMembershipID,
		status:          BookingStatusBooked,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructClassBooking(
	id ClassBookingID,
	sessionID classsession.ClassSessionID,
	gymMembershipID gym.MembershipID,
	status BookingStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (ClassBooking, error) {
	if err := validateBookingStatus(status); err != nil {
		return ClassBooking{}, err
	}
	booking, err := NewClassBooking(id, sessionID, gymMembershipID, createdAt)
	if err != nil {
		return ClassBooking{}, err
	}
	booking.status = status
	booking.updatedAt = updatedAt
	return booking, nil
}

func (b ClassBooking) ID() ClassBookingID {
	return b.id
}

func (b ClassBooking) SessionID() classsession.ClassSessionID {
	return b.sessionID
}

func (b ClassBooking) GymMembershipID() gym.MembershipID {
	return b.gymMembershipID
}

func (b ClassBooking) Status() BookingStatus {
	return b.status
}

func (b ClassBooking) CreatedAt() time.Time {
	return b.createdAt
}

func (b ClassBooking) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *ClassBooking) Cancel(now time.Time) error {
	if b.status != BookingStatusBooked {
		return ErrInvalidStatusTransition
	}
	b.status = BookingStatusCancelled
	b.updatedAt = now
	return nil
}

func (b *ClassBooking) Rebook(now time.Time) error {
	if b.status != BookingStatusCancelled {
		return ErrInvalidStatusTransition
	}
	b.status = BookingStatusBooked
	b.updatedAt = now
	return nil
}

func (b *ClassBooking) MarkAttended(now time.Time) error {
	if b.status != BookingStatusBooked {
		return ErrInvalidStatusTransition
	}
	b.status = BookingStatusAttended
	b.updatedAt = now
	return nil
}
