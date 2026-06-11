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

func (r *WODRepository) Save(ctx context.Context, w domainwod.WOD) error {
	record, stages, movements, err := wodToRecords(w)
	if err != nil {
		return err
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO wods (id, name, status, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			status = EXCLUDED.status,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.Name, record.Status, record.Description, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert wod: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM wod_stages WHERE wod_id = $1`, record.ID); err != nil {
		return fmt.Errorf("delete stages: %w", err)
	}

	for _, stage := range stages {
		_, err = tx.Exec(ctx, `
			INSERT INTO wod_stages (id, wod_id, position, stage_kind, wod_type, instructions, config, scoring_kind)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, stage.ID, stage.WODID, stage.Position, stage.StageKind, stage.WODType, stage.Instructions, stage.Config, stage.ScoringKind)
		if err != nil {
			return fmt.Errorf("insert stage: %w", err)
		}
	}

	for _, movement := range movements {
		_, err = tx.Exec(ctx, `
			INSERT INTO wod_movements (id, stage_id, position, label, name, prescription, reps, load_value, load_unit, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, movement.ID, movement.StageID, movement.Position, movement.Label, movement.Name, movement.Prescription, movement.Reps, movement.LoadValue, movement.LoadUnit, movement.Notes)
		if err != nil {
			return fmt.Errorf("insert movement: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *WODRepository) FindByID(ctx context.Context, id domainwod.WODID) (domainwod.WOD, error) {
	record, err := r.fetchRecord(ctx, id.String())
	if err != nil {
		return domainwod.WOD{}, err
	}
	stages, err := r.fetchStages(ctx, id.String())
	if err != nil {
		return domainwod.WOD{}, err
	}
	movementsByStage, err := r.fetchMovements(ctx, stageIDs(stages))
	if err != nil {
		return domainwod.WOD{}, err
	}
	return recordsToWOD(record, stages, movementsByStage)
}

func (r *WODRepository) List(ctx context.Context) ([]domainwod.WOD, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, name, status, description, created_at, updated_at
		FROM wods
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list wods: %w", err)
	}
	defer rows.Close()

	var records []wodRecord
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wods: %w", err)
	}

	wods := make([]domainwod.WOD, 0, len(records))
	for _, record := range records {
		stages, err := r.fetchStages(ctx, record.ID)
		if err != nil {
			return nil, err
		}
		movementsByStage, err := r.fetchMovements(ctx, stageIDs(stages))
		if err != nil {
			return nil, err
		}
		w, err := recordsToWOD(record, stages, movementsByStage)
		if err != nil {
			return nil, err
		}
		wods = append(wods, w)
	}
	return wods, nil
}

func (r *WODRepository) fetchRecord(ctx context.Context, id string) (wodRecord, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, name, status, description, created_at, updated_at
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

func (r *WODRepository) fetchStages(ctx context.Context, wodID string) ([]stageRecord, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, wod_id, position, stage_kind, wod_type, instructions, config, scoring_kind
		FROM wod_stages
		WHERE wod_id = $1
		ORDER BY position ASC
	`, wodID)
	if err != nil {
		return nil, fmt.Errorf("fetch stages: %w", err)
	}
	defer rows.Close()

	var stages []stageRecord
	for rows.Next() {
		var stage stageRecord
		if err := rows.Scan(
			&stage.ID,
			&stage.WODID,
			&stage.Position,
			&stage.StageKind,
			&stage.WODType,
			&stage.Instructions,
			&stage.Config,
			&stage.ScoringKind,
		); err != nil {
			return nil, fmt.Errorf("scan stage: %w", err)
		}
		stages = append(stages, stage)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stages: %w", err)
	}
	return stages, nil
}

func (r *WODRepository) fetchMovements(ctx context.Context, stageIDs []string) (map[string][]movementRecord, error) {
	movementsByStage := make(map[string][]movementRecord)
	if len(stageIDs) == 0 {
		return movementsByStage, nil
	}

	rows, err := r.db.pool.Query(ctx, `
		SELECT id, stage_id, position, label, name, prescription, reps, load_value, load_unit, notes
		FROM wod_movements
		WHERE stage_id = ANY($1)
		ORDER BY position ASC
	`, stageIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch movements: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movement movementRecord
		if err := rows.Scan(
			&movement.ID,
			&movement.StageID,
			&movement.Position,
			&movement.Label,
			&movement.Name,
			&movement.Prescription,
			&movement.Reps,
			&movement.LoadValue,
			&movement.LoadUnit,
			&movement.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan movement: %w", err)
		}
		movementsByStage[movement.StageID] = append(movementsByStage[movement.StageID], movement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movements: %w", err)
	}
	return movementsByStage, nil
}

func stageIDs(stages []stageRecord) []string {
	ids := make([]string, 0, len(stages))
	for _, stage := range stages {
		ids = append(ids, stage.ID)
	}
	return ids
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRecord(scanner rowScanner) (wodRecord, error) {
	var record wodRecord
	if err := scanner.Scan(
		&record.ID,
		&record.Name,
		&record.Status,
		&record.Description,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return wodRecord{}, fmt.Errorf("scan wod: %w", err)
	}
	return record, nil
}
