package classbooking

import "errors"

var (
	ErrClassFull               = errors.New("class session has reached capacity")
	ErrInvalidStatus           = errors.New("booking status is invalid")
	ErrInvalidStatusTransition = errors.New("invalid booking status transition")
	ErrAlreadyBooked           = errors.New("athlete already has a booking for this session")
	ErrNotFound                = errors.New("class booking not found")
)
