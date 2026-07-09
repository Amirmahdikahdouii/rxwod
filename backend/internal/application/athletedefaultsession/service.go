package athletedefaultsession

import (
	"context"
	"fmt"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainathletedefaultsession "github.com/rxwod/backend/internal/domain/athletedefaultsession"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo    Repository
	gymRepo appgym.Repository
	clock   clock.Clock
	idgen   idgen.Generator
}

func NewService(
	repo Repository,
	gymRepo appgym.Repository,
	clock clock.Clock,
	idgen idgen.Generator,
) *Service {
	return &Service{
		repo:    repo,
		gymRepo: gymRepo,
		clock:   clock,
		idgen:   idgen,
	}
}

func (s *Service) SetDefaultSession(ctx context.Context, cmd SetDefaultSessionCommand) (DefaultSessionDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionAthleteDefaultSessionManage)
	if err != nil {
		return DefaultSessionDTO{}, err
	}

	membership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return DefaultSessionDTO{}, err
	}

	now := s.clock.Now()
	pref, err := domainathletedefaultsession.NewAthleteDefaultSession(
		domainathletedefaultsession.AthleteDefaultSessionID(s.idgen.NewID()),
		membership.ID(),
		domainathletedefaultsession.DayOfWeek(cmd.DayOfWeek),
		domainathletedefaultsession.TimeSlot(cmd.TimeSlot),
		now,
	)
	if err != nil {
		return DefaultSessionDTO{}, err
	}

	if err := s.repo.Save(ctx, pref); err != nil {
		return DefaultSessionDTO{}, fmt.Errorf("save athlete default session: %w", err)
	}

	return toDefaultSessionDTO(pref), nil
}

func toDefaultSessionDTO(pref domainathletedefaultsession.AthleteDefaultSession) DefaultSessionDTO {
	return DefaultSessionDTO{
		ID:              pref.ID().String(),
		GymMembershipID: pref.GymMembershipID().String(),
		DayOfWeek:       int(pref.DayOfWeek()),
		TimeSlot:        string(pref.TimeSlot()),
		CreatedAt:       pref.CreatedAt(),
		UpdatedAt:       pref.UpdatedAt(),
	}
}
