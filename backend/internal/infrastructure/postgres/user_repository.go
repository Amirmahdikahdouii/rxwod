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

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, aggregate user.User) error {
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, display_name, email_verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			display_name = EXCLUDED.display_name,
			email_verified_at = EXCLUDED.email_verified_at,
			updated_at = EXCLUDED.updated_at
	`, aggregate.ID().String(), aggregate.Email(), aggregate.PasswordHash(), aggregate.DisplayName(), aggregate.EmailVerifiedAt(), aggregate.CreatedAt(), aggregate.UpdatedAt())
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id.String())
	return scanUser(row)
}

func scanUser(scanner rowScanner) (user.User, error) {
	var (
		id              string
		email           string
		passwordHash    string
		displayName     string
		emailVerifiedAt sql.NullTime
		createdAt       time.Time
		updatedAt       time.Time
	)
	if err := scanner.Scan(&id, &email, &passwordHash, &displayName, &emailVerifiedAt, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, appauth.ErrUserNotFound
		}
		return user.User{}, fmt.Errorf("scan user: %w", err)
	}
	var verifiedAt *time.Time
	if emailVerifiedAt.Valid {
		verifiedAt = &emailVerifiedAt.Time
	}
	return user.ReconstructUser(
		user.UserID(id),
		user.Email(email),
		user.PasswordHash(passwordHash),
		user.DisplayName(displayName),
		verifiedAt,
		createdAt,
		updatedAt,
	)
}
