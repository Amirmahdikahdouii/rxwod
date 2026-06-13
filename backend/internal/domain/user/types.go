package user

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

type UserID string
type Email string
type DisplayName string
type PasswordHash string

func (id UserID) String() string {
	return string(id)
}

func NormalizeEmail(value string) Email {
	return Email(strings.ToLower(strings.TrimSpace(value)))
}

func validateEmail(email Email) error {
	if strings.TrimSpace(string(email)) == "" {
		return ErrInvalidEmail
	}
	parsed, err := mail.ParseAddress(string(email))
	if err != nil || parsed.Address != string(email) {
		return ErrInvalidEmail
	}
	return nil
}

func validateDisplayName(name DisplayName) error {
	if utf8.RuneCountInString(strings.TrimSpace(string(name))) > 120 {
		return ErrInvalidDisplayName
	}
	return nil
}
