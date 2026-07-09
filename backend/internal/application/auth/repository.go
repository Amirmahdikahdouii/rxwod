package auth

import (
	"context"
	"errors"
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Save(ctx context.Context, user user.User) error
	FindByEmail(ctx context.Context, email user.Email) (user.User, error)
	FindByID(ctx context.Context, id user.UserID) (user.User, error)
}

type RefreshToken struct {
	ID        string
	UserID    user.UserID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeAllForUser(ctx context.Context, userID user.UserID, now time.Time) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password string, passwordHash string) error
}

type AccessTokenIssuer interface {
	Issue(userID user.UserID, now time.Time) (string, time.Time, error)
	Verify(token string) (user.UserID, error)
}

type InviteAccepter interface {
	AcceptInvitesForEmail(ctx context.Context, email user.Email, userID user.UserID) error
}

type PasswordResetToken struct {
	ID        string
	UserID    user.UserID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PasswordResetTokenRepository interface {
	Save(ctx context.Context, token PasswordResetToken) error
	FindByHash(ctx context.Context, tokenHash string) (PasswordResetToken, error)
	InvalidateUnusedForUser(ctx context.Context, userID user.UserID, now time.Time) error
	MarkUsed(ctx context.Context, id string, now time.Time) error
}
