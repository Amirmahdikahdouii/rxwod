package user

import "errors"

var (
	ErrInvalidEmail       = errors.New("email is invalid")
	ErrInvalidDisplayName = errors.New("display name must be at most 120 characters")
	ErrPasswordHashEmpty  = errors.New("password hash is required")
)
