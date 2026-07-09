package postgres

import (
	"time"

	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type classBookingRecord struct {
	ID              string
	SessionID       string
	GymMembershipID string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func classBookingToRecord(booking domainclassbooking.ClassBooking) classBookingRecord {
	return classBookingRecord{
		ID:              booking.ID().String(),
		SessionID:       booking.SessionID().String(),
		GymMembershipID: booking.GymMembershipID().String(),
		Status:          string(booking.Status()),
		CreatedAt:       booking.CreatedAt(),
		UpdatedAt:       booking.UpdatedAt(),
	}
}

func recordToClassBooking(record classBookingRecord) (domainclassbooking.ClassBooking, error) {
	return domainclassbooking.ReconstructClassBooking(
		domainclassbooking.ClassBookingID(record.ID),
		domainclasssession.ClassSessionID(record.SessionID),
		gym.MembershipID(record.GymMembershipID),
		domainclassbooking.BookingStatus(record.Status),
		record.CreatedAt,
		record.UpdatedAt,
	)
}
