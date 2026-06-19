package gym

import "errors"

var (
	ErrRoleNotAssignable          = errors.New("role cannot be assigned")
	ErrInviteeNotFound            = errors.New("invitee is not registered")
	ErrMemberNotFound             = errors.New("member not found")
	ErrOwnerMembershipProtected   = errors.New("owner membership cannot be changed or removed")
)
