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

type PasswordResetTokenRepository struct {
	db *DB
}

func NewPasswordResetTokenRepository(db *DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

func (r *PasswordResetTokenRepository) Save(ctx context.Context, token appauth.PasswordResetToken) error {
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, used_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, token.ID, token.UserID.String(), token.TokenHash, token.ExpiresAt, token.UsedAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("save password reset token: %w", err)
	}
	return nil
}

func (r *PasswordResetTokenRepository) FindByHash(ctx context.Context, tokenHash string) (appauth.PasswordResetToken, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`, tokenHash)

	var (
		token     appauth.PasswordResetToken
		userID    string
		usedAt    sql.NullTime
		createdAt time.Time
	)
	if err := row.Scan(&token.ID, &userID, &token.TokenHash, &token.ExpiresAt, &usedAt, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appauth.PasswordResetToken{}, ErrNotFound
		}
		return appauth.PasswordResetToken{}, fmt.Errorf("scan password reset token: %w", err)
	}
	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}
	token.UserID = user.UserID(userID)
	token.CreatedAt = createdAt
	return token, nil
}

func (r *PasswordResetTokenRepository) InvalidateUnusedForUser(ctx context.Context, userID user.UserID, now time.Time) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE password_reset_tokens
		SET used_at = $2
		WHERE user_id = $1 AND used_at IS NULL
	`, userID.String(), now)
	if err != nil {
		return fmt.Errorf("invalidate password reset tokens: %w", err)
	}
	return nil
}

func (r *PasswordResetTokenRepository) MarkUsed(ctx context.Context, id string, now time.Time) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE password_reset_tokens
		SET used_at = $2
		WHERE id = $1
	`, id, now)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}
	return nil
}
