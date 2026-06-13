package authz

import "errors"

var (
	ErrUnauthenticated         = errors.New("authentication is required")
	ErrGymRequired             = errors.New("x-gym-id header is required")
	ErrForbidden               = errors.New("permission denied")
	ErrActiveMembershipMissing = errors.New("active gym membership is required")
	ErrGymMismatch             = errors.New("requested gym does not match active workspace")
)
