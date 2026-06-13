package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
		INSERT INTO users (id, email, password_hash, display_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			display_name = EXCLUDED.display_name,
			updated_at = EXCLUDED.updated_at
	`, aggregate.ID().String(), aggregate.Email(), aggregate.PasswordHash(), aggregate.DisplayName(), aggregate.CreatedAt(), aggregate.UpdatedAt())
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id.String())
	return scanUser(row)
}

func scanUser(scanner rowScanner) (user.User, error) {
	var (
		id           string
		email        string
		passwordHash string
		displayName  string
		createdAt    time.Time
		updatedAt    time.Time
	)
	if err := scanner.Scan(&id, &email, &passwordHash, &displayName, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, ErrNotFound
		}
		return user.User{}, fmt.Errorf("scan user: %w", err)
	}
	return user.ReconstructUser(
		user.UserID(id),
		user.Email(email),
		user.PasswordHash(passwordHash),
		user.DisplayName(displayName),
		createdAt,
		updatedAt,
	)
}
