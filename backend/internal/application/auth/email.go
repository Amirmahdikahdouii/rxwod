package auth

import "context"

type EmailSender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetURL string) error
	SendEmailVerification(ctx context.Context, toEmail, verifyURL string) error
}
