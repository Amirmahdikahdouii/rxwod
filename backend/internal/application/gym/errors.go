package gym

import "errors"

var (
	ErrRoleNotAssignable = errors.New("role cannot be assigned")
	ErrInviteeNotFound   = errors.New("invitee is not registered")
)
