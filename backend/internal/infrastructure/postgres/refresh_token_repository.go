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

type RefreshTokenRepository struct {
	db *DB
}

func NewRefreshTokenRepository(db *DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token appauth.RefreshToken) error {
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, token.ID, token.UserID.String(), token.TokenHash, token.ExpiresAt, token.RevokedAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (appauth.RefreshToken, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash)

	var (
		token     appauth.RefreshToken
		userID    string
		revokedAt sql.NullTime
		createdAt time.Time
	)
	if err := row.Scan(&token.ID, &userID, &token.TokenHash, &token.ExpiresAt, &revokedAt, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appauth.RefreshToken{}, ErrNotFound
		}
		return appauth.RefreshToken{}, fmt.Errorf("scan refresh token: %w", err)
	}
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	token.UserID = user.UserID(userID)
	token.CreatedAt = createdAt
	return token, nil
}
