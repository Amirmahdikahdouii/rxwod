package gym

import "errors"

var (
	ErrInvalidGymName          = errors.New("gym name must be between 2 and 120 characters")
	ErrInvalidMembershipStatus = errors.New("membership status is invalid")
	ErrInvalidInvitationStatus = errors.New("invitation status is invalid")
	ErrInvalidInvitationRole   = errors.New("invitation role is invalid")
	ErrInvitationNotPending    = errors.New("invitation is not pending")
	ErrInvitationExpired       = errors.New("invitation has expired")
)
