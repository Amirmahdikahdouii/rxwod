package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	domainathletedefaultsession "github.com/rxwod/backend/internal/domain/athletedefaultsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type AthleteDefaultSessionRepository struct {
	db *DB
}

func NewAthleteDefaultSessionRepository(db *DB) *AthleteDefaultSessionRepository {
	return &AthleteDefaultSessionRepository{db: db}
}

func (r *AthleteDefaultSessionRepository) Save(
	ctx context.Context,
	pref domainathletedefaultsession.AthleteDefaultSession,
) error {
	record := athleteDefaultSessionToRecord(pref)
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO athlete_default_sessions (id, gym_membership_id, day_of_week, time_slot, created_at, updated_at)
		VALUES ($1, $2, $3, $4::time, $5, $6)
		ON CONFLICT (gym_membership_id, day_of_week, time_slot) DO UPDATE SET
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.GymMembershipID, record.DayOfWeek, record.TimeSlot, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert athlete default session: %w", err)
	}
	return nil
}

func (r *AthleteDefaultSessionRepository) FindByID(
	ctx context.Context,
	id domainathletedefaultsession.AthleteDefaultSessionID,
) (domainathletedefaultsession.AthleteDefaultSession, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, gym_membership_id, day_of_week, to_char(time_slot, 'HH24:MI') AS time_slot, created_at, updated_at
		FROM athlete_default_sessions
		WHERE id = $1
	`, id.String())

	record, err := scanAthleteDefaultSession(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainathletedefaultsession.AthleteDefaultSession{}, domainathletedefaultsession.ErrNotFound
		}
		return domainathletedefaultsession.AthleteDefaultSession{}, err
	}
	return recordToAthleteDefaultSession(record)
}

func (r *AthleteDefaultSessionRepository) Delete(
	ctx context.Context,
	id domainathletedefaultsession.AthleteDefaultSessionID,
) error {
	tag, err := r.db.pool.Exec(ctx, `
		DELETE FROM athlete_default_sessions
		WHERE id = $1
	`, id.String())
	if err != nil {
		return fmt.Errorf("delete athlete default session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainathletedefaultsession.ErrNotFound
	}
	return nil
}

func (r *AthleteDefaultSessionRepository) FindByGymMembership(
	ctx context.Context,
	membershipID gym.MembershipID,
) ([]domainathletedefaultsession.AthleteDefaultSession, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, gym_membership_id, day_of_week, to_char(time_slot, 'HH24:MI') AS time_slot, created_at, updated_at
		FROM athlete_default_sessions
		WHERE gym_membership_id = $1
		ORDER BY day_of_week ASC, time_slot ASC
	`, membershipID.String())
	if err != nil {
		return nil, fmt.Errorf("list athlete default sessions: %w", err)
	}
	defer rows.Close()

	var prefs []domainathletedefaultsession.AthleteDefaultSession
	for rows.Next() {
		record, err := scanAthleteDefaultSession(rows)
		if err != nil {
			return nil, err
		}
		pref, err := recordToAthleteDefaultSession(record)
		if err != nil {
			return nil, err
		}
		prefs = append(prefs, pref)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate athlete default sessions: %w", err)
	}
	return prefs, nil
}

func (r *AthleteDefaultSessionRepository) FindMatchingMemberships(
	ctx context.Context,
	gymID gym.GymID,
	dayOfWeek int,
	timeSlot string,
) ([]gym.MembershipID, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT ads.gym_membership_id
		FROM athlete_default_sessions ads
		JOIN gym_memberships gm ON gm.id = ads.gym_membership_id
		WHERE gm.gym_id = $1
		  AND gm.status = 'active'
		  AND ads.day_of_week = $2
		  AND ads.time_slot = $3::time
	`, gymID.String(), dayOfWeek, timeSlot)
	if err != nil {
		return nil, fmt.Errorf("find matching default sessions: %w", err)
	}
	defer rows.Close()

	var membershipIDs []gym.MembershipID
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan matching membership: %w", err)
		}
		membershipIDs = append(membershipIDs, gym.MembershipID(id))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate matching memberships: %w", err)
	}
	return membershipIDs, nil
}

func scanAthleteDefaultSession(scanner rowScanner) (athleteDefaultSessionRecord, error) {
	var record athleteDefaultSessionRecord
	var createdAt time.Time
	var updatedAt time.Time
	if err := scanner.Scan(
		&record.ID,
		&record.GymMembershipID,
		&record.DayOfWeek,
		&record.TimeSlot,
		&createdAt,
		&updatedAt,
	); err != nil {
		return athleteDefaultSessionRecord{}, err
	}
	record.CreatedAt = createdAt
	record.UpdatedAt = updatedAt
	return record, nil
}
