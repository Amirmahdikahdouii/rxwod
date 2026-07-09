package classsession

import "time"

type ClassSessionDTO struct {
	ID              string
	GymID           string
	WodID           *string
	StartTime       time.Time
	EndTime         time.Time
	Capacity        int
	CoachID         string
	BookedCount     int
	MyBookingStatus *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateClassSessionResultDTO struct {
	Session                 ClassSessionDTO
	AutoBookedCount         int
	AutoBookedMembershipIDs []string
}
