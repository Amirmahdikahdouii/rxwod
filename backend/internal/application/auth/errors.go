package auth

import "errors"

var (
	ErrInvalidCredentials       = errors.New("invalid email or password")
	ErrPasswordTooShort           = errors.New("password must be at least 8 characters")
	ErrRefreshTokenInvalid        = errors.New("refresh token is invalid")
	ErrResetTokenInvalid          = errors.New("password reset token is invalid or expired")
	ErrVerificationTokenInvalid   = errors.New("email verification token is invalid or expired")
	ErrEmailNotVerified           = errors.New("email verification is required")
)
