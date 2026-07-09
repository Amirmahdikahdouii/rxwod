package postgres

import (
	"database/sql"
	"time"

	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type classSessionRecord struct {
	ID        string
	GymID     string
	WodID     sql.NullString
	StartTime time.Time
	EndTime   time.Time
	Capacity  int
	CoachID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func classSessionToRecord(session domainclasssession.ClassSession) classSessionRecord {
	record := classSessionRecord{
		ID:        session.ID().String(),
		GymID:     session.GymID().String(),
		StartTime: session.StartTime(),
		EndTime:   session.EndTime(),
		Capacity:  int(session.Capacity()),
		CoachID:   session.CoachID().String(),
		CreatedAt: session.CreatedAt(),
		UpdatedAt: session.UpdatedAt(),
	}
	if wodID := session.WODID(); wodID != nil {
		record.WodID = sql.NullString{String: wodID.String(), Valid: true}
	}
	return record
}

func recordToClassSession(record classSessionRecord) (domainclasssession.ClassSession, error) {
	var wodID *domainwod.WODID
	if record.WodID.Valid {
		value := domainwod.WODID(record.WodID.String)
		wodID = &value
	}
	return domainclasssession.ReconstructClassSession(
		domainclasssession.ClassSessionID(record.ID),
		gym.GymID(record.GymID),
		wodID,
		record.StartTime,
		record.EndTime,
		domainclasssession.Capacity(record.Capacity),
		gym.MembershipID(record.CoachID),
		record.CreatedAt,
		record.UpdatedAt,
	)
}
