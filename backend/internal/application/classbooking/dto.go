package classbooking

import "time"

type BookingDTO struct {
	ID              string
	SessionID       string
	GymMembershipID string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
