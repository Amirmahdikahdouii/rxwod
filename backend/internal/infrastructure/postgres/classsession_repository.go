package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type ClassSessionRepository struct {
	db *DB
}

func NewClassSessionRepository(db *DB) *ClassSessionRepository {
	return &ClassSessionRepository{db: db}
}

func (r *ClassSessionRepository) Save(ctx context.Context, session domainclasssession.ClassSession) error {
	record := classSessionToRecord(session)
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO class_sessions (id, gym_id, wod_id, start_time, end_time, capacity, coach_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			wod_id = EXCLUDED.wod_id,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			capacity = EXCLUDED.capacity,
			coach_id = EXCLUDED.coach_id,
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.GymID, nullStringArg(record.WodID), record.StartTime, record.EndTime, record.Capacity, record.CoachID, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert class session: %w", err)
	}
	return nil
}

func (r *ClassSessionRepository) SaveWithDefaultBookings(
	ctx context.Context,
	session domainclasssession.ClassSession,
	bookings []domainclassbooking.ClassBooking,
) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	record := classSessionToRecord(session)
	_, err = tx.Exec(ctx, `
		INSERT INTO class_sessions (id, gym_id, wod_id, start_time, end_time, capacity, coach_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, record.ID, record.GymID, nullStringArg(record.WodID), record.StartTime, record.EndTime, record.Capacity, record.CoachID, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert class session: %w", err)
	}

	for _, booking := range bookings {
		bookingRecord := classBookingToRecord(booking)
		_, err = tx.Exec(ctx, `
			INSERT INTO class_bookings (id, session_id, gym_membership_id, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (session_id, gym_membership_id) DO NOTHING
		`, bookingRecord.ID, bookingRecord.SessionID, bookingRecord.GymMembershipID, bookingRecord.Status, bookingRecord.CreatedAt, bookingRecord.UpdatedAt)
		if err != nil {
			return fmt.Errorf("insert default booking: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *ClassSessionRepository) FindByID(
	ctx context.Context,
	gymID gym.GymID,
	id domainclasssession.ClassSessionID,
) (domainclasssession.ClassSession, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, gym_id, wod_id, start_time, end_time, capacity, coach_id, created_at, updated_at
		FROM class_sessions
		WHERE gym_id = $1 AND id = $2
	`, gymID.String(), id.String())

	record, err := scanClassSession(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainclasssession.ClassSession{}, domainclasssession.ErrNotFound
		}
		return domainclasssession.ClassSession{}, err
	}
	return recordToClassSession(record)
}

func (r *ClassSessionRepository) ListByGymAndDate(
	ctx context.Context,
	gymID gym.GymID,
	from, to time.Time,
) ([]domainclasssession.ClassSession, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, gym_id, wod_id, start_time, end_time, capacity, coach_id, created_at, updated_at
		FROM class_sessions
		WHERE gym_id = $1 AND start_time >= $2 AND start_time < $3
		ORDER BY start_time ASC
	`, gymID.String(), from, to)
	if err != nil {
		return nil, fmt.Errorf("list class sessions: %w", err)
	}
	defer rows.Close()

	var sessions []domainclasssession.ClassSession
	for rows.Next() {
		record, err := scanClassSession(rows)
		if err != nil {
			return nil, err
		}
		session, err := recordToClassSession(record)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate class sessions: %w", err)
	}
	return sessions, nil
}

func (r *ClassSessionRepository) Delete(ctx context.Context, gymID gym.GymID, id domainclasssession.ClassSessionID) error {
	tag, err := r.db.pool.Exec(ctx, `
		DELETE FROM class_sessions
		WHERE gym_id = $1 AND id = $2
	`, gymID.String(), id.String())
	if err != nil {
		return fmt.Errorf("delete class session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainclasssession.ErrNotFound
	}
	return nil
}

func scanClassSession(scanner rowScanner) (classSessionRecord, error) {
	var record classSessionRecord
	if err := scanner.Scan(
		&record.ID,
		&record.GymID,
		&record.WodID,
		&record.StartTime,
		&record.EndTime,
		&record.Capacity,
		&record.CoachID,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return classSessionRecord{}, err
	}
	return record, nil
}

func nullStringArg(value sql.NullString) any {
	if !value.Valid {
		return nil
	}
	return value.String
}
