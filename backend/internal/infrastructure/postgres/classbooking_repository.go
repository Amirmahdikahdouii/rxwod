package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type ClassBookingRepository struct {
	db *DB
}

func NewClassBookingRepository(db *DB) *ClassBookingRepository {
	return &ClassBookingRepository{db: db}
}

func (r *ClassBookingRepository) Save(ctx context.Context, booking domainclassbooking.ClassBooking) error {
	record := classBookingToRecord(booking)
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO class_bookings (id, session_id, gym_membership_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (session_id, gym_membership_id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`, record.ID, record.SessionID, record.GymMembershipID, record.Status, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert class booking: %w", err)
	}
	return nil
}

func (r *ClassBookingRepository) FindByID(
	ctx context.Context,
	id domainclassbooking.ClassBookingID,
) (domainclassbooking.ClassBooking, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, session_id, gym_membership_id, status, created_at, updated_at
		FROM class_bookings
		WHERE id = $1
	`, id.String())

	record, err := scanClassBooking(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainclassbooking.ClassBooking{}, domainclassbooking.ErrNotFound
		}
		return domainclassbooking.ClassBooking{}, err
	}
	return recordToClassBooking(record)
}

func (r *ClassBookingRepository) FindBySessionAndMembership(
	ctx context.Context,
	sessionID domainclasssession.ClassSessionID,
	membershipID gym.MembershipID,
) (domainclassbooking.ClassBooking, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, session_id, gym_membership_id, status, created_at, updated_at
		FROM class_bookings
		WHERE session_id = $1 AND gym_membership_id = $2
	`, sessionID.String(), membershipID.String())

	record, err := scanClassBooking(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainclassbooking.ClassBooking{}, domainclassbooking.ErrNotFound
		}
		return domainclassbooking.ClassBooking{}, err
	}
	return recordToClassBooking(record)
}

func (r *ClassBookingRepository) ListBookingsBySession(
	ctx context.Context,
	sessionID domainclasssession.ClassSessionID,
) ([]domainclassbooking.ClassBooking, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, session_id, gym_membership_id, status, created_at, updated_at
		FROM class_bookings
		WHERE session_id = $1
		ORDER BY created_at ASC
	`, sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("list class bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domainclassbooking.ClassBooking
	for rows.Next() {
		record, err := scanClassBooking(rows)
		if err != nil {
			return nil, err
		}
		booking, err := recordToClassBooking(record)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate class bookings: %w", err)
	}
	return bookings, nil
}

func (r *ClassBookingRepository) CountBookedBySession(
	ctx context.Context,
	sessionID domainclasssession.ClassSessionID,
) (int, error) {
	var count int
	err := r.db.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM class_bookings
		WHERE session_id = $1 AND status = 'BOOKED'
	`, sessionID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count booked class bookings: %w", err)
	}
	return count, nil
}

func scanClassBooking(scanner rowScanner) (classBookingRecord, error) {
	var record classBookingRecord
	if err := scanner.Scan(
		&record.ID,
		&record.SessionID,
		&record.GymMembershipID,
		&record.Status,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return classBookingRecord{}, err
	}
	return record, nil
}
