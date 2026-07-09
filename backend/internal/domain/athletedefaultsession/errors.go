package athletedefaultsession

import "errors"

var (
	ErrMembershipRequired = errors.New("gym membership is required")
	ErrInvalidDayOfWeek   = errors.New("day of week must be between 0 and 6")
	ErrInvalidTimeSlot    = errors.New("time slot must be in HH:MM format")
	ErrNotFound           = errors.New("athlete default session not found")
)
