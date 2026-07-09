package postgres

import (
	"database/sql"
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	domainwodresult "github.com/rxwod/backend/internal/domain/wodresult"
)

type wodResultRecord struct {
	ID              string
	WODID           string
	GymMembershipID string
	ScoreValue      int
	IsRx            bool
	Notes           sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func wodResultToRecord(r domainwodresult.WODResult) wodResultRecord {
	return wodResultRecord{
		ID:              r.ID().String(),
		WODID:           r.WODID().String(),
		GymMembershipID: r.GymMembershipID().String(),
		ScoreValue:      int(r.ScoreValue()),
		IsRx:            r.IsRx(),
		Notes:           sql.NullString{String: string(r.Notes()), Valid: r.Notes() != ""},
		CreatedAt:       r.CreatedAt(),
		UpdatedAt:       r.UpdatedAt(),
	}
}

func recordToWODResult(record wodResultRecord) (domainwodresult.WODResult, error) {
	return domainwodresult.ReconstructWODResult(
		domainwodresult.WODResultID(record.ID),
		domainwod.WODID(record.WODID),
		gym.MembershipID(record.GymMembershipID),
		domainwodresult.ScoreValue(record.ScoreValue),
		record.IsRx,
		domainwodresult.Notes(record.Notes.String),
		record.CreatedAt,
		record.UpdatedAt,
	)
}
