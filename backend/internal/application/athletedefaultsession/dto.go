package athletedefaultsession

import "time"

type DefaultSessionDTO struct {
	ID              string
	GymMembershipID string
	DayOfWeek       int
	TimeSlot        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
