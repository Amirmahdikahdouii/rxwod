package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	appauth "github.com/rxwod/backend/internal/application/auth"
	"github.com/rxwod/backend/internal/domain/user"
)

type EmailVerificationTokenRepository struct {
	db *DB
}

func NewEmailVerificationTokenRepository(db *DB) *EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{db: db}
}

func (r *EmailVerificationTokenRepository) Save(ctx context.Context, token appauth.EmailVerificationToken) error {
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO email_verification_tokens (id, user_id, token_hash, expires_at, used_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, token.ID, token.UserID.String(), token.TokenHash, token.ExpiresAt, token.UsedAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("save email verification token: %w", err)
	}
	return nil
}

func (r *EmailVerificationTokenRepository) FindByHash(ctx context.Context, tokenHash string) (appauth.EmailVerificationToken, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM email_verification_tokens
		WHERE token_hash = $1
	`, tokenHash)

	var (
		token     appauth.EmailVerificationToken
		userID    string
		usedAt    sql.NullTime
		createdAt time.Time
	)
	if err := row.Scan(&token.ID, &userID, &token.TokenHash, &token.ExpiresAt, &usedAt, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appauth.EmailVerificationToken{}, ErrNotFound
		}
		return appauth.EmailVerificationToken{}, fmt.Errorf("scan email verification token: %w", err)
	}
	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}
	token.UserID = user.UserID(userID)
	token.CreatedAt = createdAt
	return token, nil
}

func (r *EmailVerificationTokenRepository) InvalidateUnusedForUser(ctx context.Context, userID user.UserID, now time.Time) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE email_verification_tokens
		SET used_at = $2
		WHERE user_id = $1 AND used_at IS NULL
	`, userID.String(), now)
	if err != nil {
		return fmt.Errorf("invalidate email verification tokens: %w", err)
	}
	return nil
}

func (r *EmailVerificationTokenRepository) MarkUsed(ctx context.Context, id string, now time.Time) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE email_verification_tokens
		SET used_at = $2
		WHERE id = $1
	`, id, now)
	if err != nil {
		return fmt.Errorf("mark email verification token used: %w", err)
	}
	return nil
}
