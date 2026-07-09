package athletedefaultsession

import (
	"context"

	domainathletedefaultsession "github.com/rxwod/backend/internal/domain/athletedefaultsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type Repository interface {
	Save(ctx context.Context, pref domainathletedefaultsession.AthleteDefaultSession) error
	FindByGymMembership(ctx context.Context, membershipID gym.MembershipID) ([]domainathletedefaultsession.AthleteDefaultSession, error)
	FindMatchingMemberships(ctx context.Context, gymID gym.GymID, dayOfWeek int, timeSlot string) ([]gym.MembershipID, error)
}
