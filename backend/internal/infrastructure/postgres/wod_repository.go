package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

var ErrNotFound = errors.New("wod not found")

type WODRepository struct {
	db *DB
}

func NewWODRepository(db *DB) *WODRepository {
	return &WODRepository{db: db}
}

func (r *WODRepository) Save(ctx context.Context, variant domainwod.Variant) error {
	record, movements, err := variantToRecord(variant)
	if err != nil {
		return err
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO wods (id, name, wod_type, status, description, config, scoring_kind, scoring_config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			wod_type = EXCLUDED.wod_type,
			status = EXCLUDED.status,
			description = EXCLUDED.description,
			config = EXCLUDED.config,
			scoring_kind = EXCLUDED.scoring_kind,
			scoring_config = EXCLUDED.scoring_config,
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.Name, record.WODType, record.Status, record.Description, record.Config, record.ScoringKind, record.ScoringConfig, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert wod: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM wod_movements WHERE wod_id = $1`, record.ID); err != nil {
		return fmt.Errorf("delete movements: %w", err)
	}

	for _, movement := range movements {
		_, err = tx.Exec(ctx, `
			INSERT INTO wod_movements (id, wod_id, position, name, reps, load_value, load_unit, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, movement.ID, movement.WODID, movement.Position, movement.Name, movement.Reps, movement.LoadValue, movement.LoadUnit, movement.Notes)
		if err != nil {
			return fmt.Errorf("insert movement: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *WODRepository) FindByID(ctx context.Context, id domainwod.WODID) (domainwod.Variant, error) {
	record, err := r.fetchRecord(ctx, id.String())
	if err != nil {
		return nil, err
	}
	movements, err := r.fetchMovements(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return recordToVariant(record, movements)
}

func (r *WODRepository) List(ctx context.Context) ([]domainwod.Variant, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, name, wod_type, status, description, config, scoring_kind, scoring_config, created_at, updated_at
		FROM wods
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list wods: %w", err)
	}
	defer rows.Close()

	var variants []domainwod.Variant
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		movements, err := r.fetchMovements(ctx, record.ID)
		if err != nil {
			return nil, err
		}
		variant, err := recordToVariant(record, movements)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wods: %w", err)
	}
	return variants, nil
}

func (r *WODRepository) fetchRecord(ctx context.Context, id string) (wodRecord, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, name, wod_type, status, description, config, scoring_kind, scoring_config, created_at, updated_at
		FROM wods
		WHERE id = $1
	`, id)
	record, err := scanRecord(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return wodRecord{}, ErrNotFound
		}
		return wodRecord{}, err
	}
	return record, nil
}

func (r *WODRepository) fetchMovements(ctx context.Context, wodID string) ([]movementRecord, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, wod_id, position, name, reps, load_value, load_unit, notes
		FROM wod_movements
		WHERE wod_id = $1
		ORDER BY position ASC
	`, wodID)
	if err != nil {
		return nil, fmt.Errorf("fetch movements: %w", err)
	}
	defer rows.Close()

	var movements []movementRecord
	for rows.Next() {
		var movement movementRecord
		if err := rows.Scan(
			&movement.ID,
			&movement.WODID,
			&movement.Position,
			&movement.Name,
			&movement.Reps,
			&movement.LoadValue,
			&movement.LoadUnit,
			&movement.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan movement: %w", err)
		}
		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movements: %w", err)
	}
	return movements, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRecord(scanner rowScanner) (wodRecord, error) {
	var record wodRecord
	if err := scanner.Scan(
		&record.ID,
		&record.Name,
		&record.WODType,
		&record.Status,
		&record.Description,
		&record.Config,
		&record.ScoringKind,
		&record.ScoringConfig,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return wodRecord{}, fmt.Errorf("scan wod: %w", err)
	}
	return record, nil
}
