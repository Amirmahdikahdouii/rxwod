package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	appwodresult "github.com/rxwod/backend/internal/application/wodresult"
	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	domainwodresult "github.com/rxwod/backend/internal/domain/wodresult"
)

var ErrWODResultNotFound = errors.New("wod result not found")

type WODResultRepository struct {
	db *DB
}

func NewWODResultRepository(db *DB) *WODResultRepository {
	return &WODResultRepository{db: db}
}

func (r *WODResultRepository) Save(ctx context.Context, result domainwodresult.WODResult) error {
	record := wodResultToRecord(result)

	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO wod_results (id, wod_id, gym_membership_id, score_value, is_rx, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (wod_id, gym_membership_id) DO UPDATE SET
			score_value = EXCLUDED.score_value,
			is_rx = EXCLUDED.is_rx,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.WODID, record.GymMembershipID, record.ScoreValue, record.IsRx, record.Notes, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert wod result: %w", err)
	}
	return nil
}

func (r *WODResultRepository) FindByWODAndMembership(
	ctx context.Context,
	wodID domainwod.WODID,
	gymMembershipID gym.MembershipID,
) (domainwodresult.WODResult, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, wod_id, gym_membership_id, score_value, is_rx, notes, created_at, updated_at
		FROM wod_results
		WHERE wod_id = $1 AND gym_membership_id = $2
	`, wodID.String(), gymMembershipID.String())

	record, err := scanWODResult(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainwodresult.WODResult{}, ErrWODResultNotFound
		}
		return domainwodresult.WODResult{}, err
	}
	return recordToWODResult(record)
}

func (r *WODResultRepository) ListByWOD(ctx context.Context, wodID domainwod.WODID) ([]domainwodresult.WODResult, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, wod_id, gym_membership_id, score_value, is_rx, notes, created_at, updated_at
		FROM wod_results
		WHERE wod_id = $1
		ORDER BY created_at ASC
	`, wodID.String())
	if err != nil {
		return nil, fmt.Errorf("list wod results: %w", err)
	}
	defer rows.Close()

	var results []domainwodresult.WODResult
	for rows.Next() {
		record, err := scanWODResult(rows)
		if err != nil {
			return nil, err
		}
		result, err := recordToWODResult(record)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wod results: %w", err)
	}
	return results, nil
}

func (r *WODResultRepository) ListLeaderboardByWOD(ctx context.Context, wodID domainwod.WODID) ([]appwodresult.LeaderboardRow, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT wr.id, wr.wod_id, wr.gym_membership_id, wr.score_value, wr.is_rx, wr.notes,
		       wr.created_at, wr.updated_at, u.display_name
		FROM wod_results wr
		JOIN gym_memberships gm ON gm.id = wr.gym_membership_id
		JOIN users u ON u.id = gm.user_id
		WHERE wr.wod_id = $1
	`, wodID.String())
	if err != nil {
		return nil, fmt.Errorf("list leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []appwodresult.LeaderboardRow
	for rows.Next() {
		var record wodResultRecord
		var displayName string
		if err := rows.Scan(
			&record.ID,
			&record.WODID,
			&record.GymMembershipID,
			&record.ScoreValue,
			&record.IsRx,
			&record.Notes,
			&record.CreatedAt,
			&record.UpdatedAt,
			&displayName,
		); err != nil {
			return nil, fmt.Errorf("scan leaderboard row: %w", err)
		}
		result, err := recordToWODResult(record)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, appwodresult.LeaderboardRow{
			Result:      result,
			DisplayName: displayName,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate leaderboard rows: %w", err)
	}
	return leaderboard, nil
}

func scanWODResult(scanner rowScanner) (wodResultRecord, error) {
	var record wodResultRecord
	if err := scanner.Scan(
		&record.ID,
		&record.WODID,
		&record.GymMembershipID,
		&record.ScoreValue,
		&record.IsRx,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return wodResultRecord{}, fmt.Errorf("scan wod result: %w", err)
	}
	return record, nil
}
