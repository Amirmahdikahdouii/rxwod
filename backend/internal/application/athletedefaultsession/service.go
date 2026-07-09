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

func (s *Service) ListMyDefaultSessions(ctx context.Context) ([]DefaultSessionDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionAthleteDefaultSessionManage)
	if err != nil {
		return nil, err
	}

	membership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return nil, err
	}

	prefs, err := s.repo.FindByGymMembership(ctx, membership.ID())
	if err != nil {
		return nil, fmt.Errorf("list athlete default sessions: %w", err)
	}

	result := make([]DefaultSessionDTO, 0, len(prefs))
	for _, pref := range prefs {
		result = append(result, toDefaultSessionDTO(pref))
	}
	return result, nil
}

func (s *Service) RemoveDefaultSession(ctx context.Context, id string) error {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionAthleteDefaultSessionManage)
	if err != nil {
		return err
	}

	membership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return err
	}

	pref, err := s.repo.FindByID(ctx, domainathletedefaultsession.AthleteDefaultSessionID(id))
	if err != nil {
		return err
	}

	if pref.GymMembershipID() != membership.ID() {
		return appauthz.ErrForbidden
	}

	if err := s.repo.Delete(ctx, pref.ID()); err != nil {
		return fmt.Errorf("delete athlete default session: %w", err)
	}
	return nil
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
