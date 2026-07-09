package email

import (
	"context"
	"log/slog"
)

type LogSender struct{}

func NewLogSender() LogSender {
	return LogSender{}
}

func (LogSender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	slog.Info("password reset email sent", "event", "email.password_reset_sent", "toEmail", toEmail, "resetURL", resetURL)
	return nil
}

func (LogSender) SendEmailVerification(_ context.Context, toEmail, verifyURL string) error {
	slog.Info("email verification sent", "event", "email.verification_sent", "toEmail", toEmail, "verifyURL", verifyURL)
	return nil
}
