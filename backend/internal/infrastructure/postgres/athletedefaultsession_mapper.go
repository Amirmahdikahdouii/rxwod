package postgres

import (
	"time"

	domainathletedefaultsession "github.com/rxwod/backend/internal/domain/athletedefaultsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type athleteDefaultSessionRecord struct {
	ID              string
	GymMembershipID string
	DayOfWeek       int
	TimeSlot        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func athleteDefaultSessionToRecord(pref domainathletedefaultsession.AthleteDefaultSession) athleteDefaultSessionRecord {
	return athleteDefaultSessionRecord{
		ID:              pref.ID().String(),
		GymMembershipID: pref.GymMembershipID().String(),
		DayOfWeek:       int(pref.DayOfWeek()),
		TimeSlot:        string(pref.TimeSlot()),
		CreatedAt:       pref.CreatedAt(),
		UpdatedAt:       pref.UpdatedAt(),
	}
}

func recordToAthleteDefaultSession(record athleteDefaultSessionRecord) (domainathletedefaultsession.AthleteDefaultSession, error) {
	return domainathletedefaultsession.ReconstructAthleteDefaultSession(
		domainathletedefaultsession.AthleteDefaultSessionID(record.ID),
		gym.MembershipID(record.GymMembershipID),
		domainathletedefaultsession.DayOfWeek(record.DayOfWeek),
		domainathletedefaultsession.TimeSlot(record.TimeSlot),
		record.CreatedAt,
		record.UpdatedAt,
	)
}
