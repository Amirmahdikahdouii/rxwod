package classsession

import "errors"

var (
	ErrInvalidTimeRange = errors.New("end time must be after start time")
	ErrInvalidCapacity  = errors.New("capacity must be greater than zero")
	ErrCoachRequired    = errors.New("coach is required")
	ErrNotFound         = errors.New("class session not found")
)
