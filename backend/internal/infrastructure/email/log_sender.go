package email

import (
	"context"
	"log"
)

type LogSender struct{}

func NewLogSender() LogSender {
	return LogSender{}
}

func (LogSender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	log.Printf("password reset for %s: %s", toEmail, resetURL)
	return nil
}
