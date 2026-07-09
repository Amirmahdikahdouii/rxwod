package classsession

import (
	"context"
	"fmt"
	"time"

	appathletedefaultsession "github.com/rxwod/backend/internal/application/athletedefaultsession"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo        Repository
	defaultRepo appathletedefaultsession.Repository
	gymRepo     appgym.Repository
	clock       clock.Clock
	idgen       idgen.Generator
}

func NewService(
	repo Repository,
	defaultRepo appathletedefaultsession.Repository,
	gymRepo appgym.Repository,
	clock clock.Clock,
	idgen idgen.Generator,
) *Service {
	return &Service{
		repo:        repo,
		defaultRepo: defaultRepo,
		gymRepo:     gymRepo,
		clock:       clock,
		idgen:       idgen,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateClassSessionCommand) (CreateClassSessionResultDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassSessionCreate)
	if err != nil {
		return CreateClassSessionResultDTO{}, err
	}

	coachMembership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return CreateClassSessionResultDTO{}, err
	}

	var wodID *domainwod.WODID
	if cmd.WodID != nil && *cmd.WodID != "" {
		value := domainwod.WODID(*cmd.WodID)
		wodID = &value
	}

	now := s.clock.Now()
	session, err := domainclasssession.NewClassSession(
		domainclasssession.ClassSessionID(s.idgen.NewID()),
		principal.GymID,
		wodID,
		cmd.StartTime,
		cmd.EndTime,
		domainclasssession.Capacity(cmd.Capacity),
		coachMembership.ID(),
		now,
	)
	if err != nil {
		return CreateClassSessionResultDTO{}, err
	}

	dayOfWeek, timeSlot := deriveSlot(cmd.StartTime)
	matches, err := s.defaultRepo.FindMatchingMemberships(ctx, principal.GymID, dayOfWeek, timeSlot)
	if err != nil {
		return CreateClassSessionResultDTO{}, fmt.Errorf("find matching default sessions: %w", err)
	}

	bookings := make([]domainclassbooking.ClassBooking, 0, len(matches))
	membershipIDs := make([]string, 0, len(matches))
	for _, membershipID := range matches {
		booking, err := domainclassbooking.NewClassBooking(
			domainclassbooking.ClassBookingID(s.idgen.NewID()),
			session.ID(),
			membershipID,
			now,
		)
		if err != nil {
			return CreateClassSessionResultDTO{}, err
		}
		bookings = append(bookings, booking)
		membershipIDs = append(membershipIDs, membershipID.String())
	}

	if err := s.repo.SaveWithDefaultBookings(ctx, session, bookings); err != nil {
		return CreateClassSessionResultDTO{}, fmt.Errorf("save class session with default bookings: %w", err)
	}

	return CreateClassSessionResultDTO{
		Session:                 toClassSessionDTO(session),
		AutoBookedCount:         len(membershipIDs),
		AutoBookedMembershipIDs: membershipIDs,
	}, nil
}

func (s *Service) List(ctx context.Context, cmd ListClassSessionsCommand) ([]ClassSessionDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassSessionRead)
	if err != nil {
		return nil, err
	}

	sessions, err := s.repo.ListByGymAndDate(ctx, principal.GymID, cmd.From, cmd.To)
	if err != nil {
		return nil, fmt.Errorf("list class sessions: %w", err)
	}

	result := make([]ClassSessionDTO, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, toClassSessionDTO(session))
	}
	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (ClassSessionDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassSessionRead)
	if err != nil {
		return ClassSessionDTO{}, err
	}

	session, err := s.repo.FindByID(ctx, principal.GymID, domainclasssession.ClassSessionID(id))
	if err != nil {
		return ClassSessionDTO{}, err
	}
	return toClassSessionDTO(session), nil
}

func deriveSlot(startTime time.Time) (int, string) {
	return int(startTime.Weekday()), startTime.Format("15:04")
}

func toClassSessionDTO(session domainclasssession.ClassSession) ClassSessionDTO {
	dto := ClassSessionDTO{
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
		value := wodID.String()
		dto.WodID = &value
	}
	return dto
}
