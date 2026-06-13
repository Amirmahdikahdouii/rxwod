package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

type fakeUsers struct{}

func (fakeUsers) Save(context.Context, user.User) error { return nil }
func (fakeUsers) FindByEmail(context.Context, user.Email) (user.User, error) {
	return user.User{}, errors.New("not found")
}
func (fakeUsers) FindByID(context.Context, user.UserID) (user.User, error) {
	return user.User{}, errors.New("not found")
}

type fakeRefreshTokens struct{}

func (fakeRefreshTokens) Save(context.Context, RefreshToken) error { return nil }
func (fakeRefreshTokens) FindByHash(context.Context, string) (RefreshToken, error) {
	return RefreshToken{}, errors.New("not found")
}

type fakeHasher struct{}

func (fakeHasher) Hash(string) (string, error) { return "hash", nil }
func (fakeHasher) Verify(string, string) error { return nil }

type fakeIssuer struct{}

func (fakeIssuer) Issue(user.UserID, time.Time) (string, time.Time, error) {
	return "access-token", time.Now().Add(time.Minute), nil
}
func (fakeIssuer) Verify(string) (user.UserID, error) { return "user-1", nil }

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time { return f.now }

type fakeIDGen struct{}

func (fakeIDGen) NewID() string { return "generated-id" }

func TestRegisterRejectsShortPassword(t *testing.T) {
	service := NewService(fakeUsers{}, fakeRefreshTokens{}, fakeHasher{}, fakeIssuer{}, nil, fakeClock{now: time.Now().UTC()}, fakeIDGen{}, time.Hour)

	_, err := service.Register(context.Background(), RegisterCommand{
		Email:       "owner@example.com",
		Password:    "short",
		DisplayName: "Owner",
	})
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}
