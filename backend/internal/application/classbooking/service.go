package classbooking

import (
	"context"
	"errors"
	"fmt"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	appclasssession "github.com/rxwod/backend/internal/application/classsession"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo        Repository
	sessionRepo appclasssession.Repository
	gymRepo     appgym.Repository
	clock       clock.Clock
	idgen       idgen.Generator
}

func NewService(
	repo Repository,
	sessionRepo appclasssession.Repository,
	gymRepo appgym.Repository,
	clock clock.Clock,
	idgen idgen.Generator,
) *Service {
	return &Service{
		repo:        repo,
		sessionRepo: sessionRepo,
		gymRepo:     gymRepo,
		clock:       clock,
		idgen:       idgen,
	}
}

func (s *Service) Book(ctx context.Context, cmd BookCommand) (BookingDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassBookingBook)
	if err != nil {
		return BookingDTO{}, err
	}

	membership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return BookingDTO{}, err
	}

	session, err := s.sessionRepo.FindByID(ctx, principal.GymID, domainclasssession.ClassSessionID(cmd.SessionID))
	if err != nil {
		return BookingDTO{}, err
	}

	return s.bookForMembership(ctx, session, membership.ID(), false)
}

func (s *Service) Overbook(ctx context.Context, cmd OverbookCommand) (BookingDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassBookingOverbook)
	if err != nil {
		return BookingDTO{}, err
	}

	session, err := s.sessionRepo.FindByID(ctx, principal.GymID, domainclasssession.ClassSessionID(cmd.SessionID))
	if err != nil {
		return BookingDTO{}, err
	}

	targetMembership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, user.UserID(cmd.AthleteUserID))
	if err != nil {
		return BookingDTO{}, err
	}

	return s.bookForMembership(ctx, session, targetMembership.ID(), true)
}

func (s *Service) ListBookings(ctx context.Context, sessionID string) ([]BookingDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassBookingRoster)
	if err != nil {
		return nil, err
	}

	if _, err := s.sessionRepo.FindByID(ctx, principal.GymID, domainclasssession.ClassSessionID(sessionID)); err != nil {
		return nil, err
	}

	bookings, err := s.repo.ListBookingsBySession(ctx, domainclasssession.ClassSessionID(sessionID))
	if err != nil {
		return nil, fmt.Errorf("list bookings by session: %w", err)
	}

	result := make([]BookingDTO, 0, len(bookings))
	for _, booking := range bookings {
		result = append(result, toBookingDTO(booking))
	}
	return result, nil
}

func (s *Service) Cancel(ctx context.Context, cmd CancelCommand) (BookingDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionClassBookingCancel)
	if err != nil {
		return BookingDTO{}, err
	}

	if _, err := s.sessionRepo.FindByID(ctx, principal.GymID, domainclasssession.ClassSessionID(cmd.SessionID)); err != nil {
		return BookingDTO{}, err
	}

	var targetMembership gym.Membership
	if cmd.AthleteUserID == nil || *cmd.AthleteUserID == "" {
		targetMembership, err = s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
		if err != nil {
			return BookingDTO{}, err
		}
	} else {
		if principal.Role != domainauthz.RoleOwner && principal.Role != domainauthz.RoleCoach {
			return BookingDTO{}, appauthz.ErrForbidden
		}
		targetMembership, err = s.gymRepo.FindActiveMembership(ctx, principal.GymID, user.UserID(*cmd.AthleteUserID))
		if err != nil {
			return BookingDTO{}, err
		}
	}

	booking, err := s.repo.FindBySessionAndMembership(ctx, domainclasssession.ClassSessionID(cmd.SessionID), targetMembership.ID())
	if err != nil {
		return BookingDTO{}, err
	}

	now := s.clock.Now()
	if err := booking.Cancel(now); err != nil {
		return BookingDTO{}, err
	}

	if err := s.repo.Save(ctx, booking); err != nil {
		return BookingDTO{}, fmt.Errorf("save cancelled booking: %w", err)
	}

	return toBookingDTO(booking), nil
}

func (s *Service) bookForMembership(
	ctx context.Context,
	session domainclasssession.ClassSession,
	membershipID gym.MembershipID,
	bypassCapacity bool,
) (BookingDTO, error) {
	now := s.clock.Now()
	existing, err := s.repo.FindBySessionAndMembership(ctx, session.ID(), membershipID)
	if err != nil && !errors.Is(err, domainclassbooking.ErrNotFound) {
		return BookingDTO{}, err
	}

	if err == nil {
		switch existing.Status() {
		case domainclassbooking.BookingStatusBooked, domainclassbooking.BookingStatusAttended:
			return BookingDTO{}, domainclassbooking.ErrAlreadyBooked
		case domainclassbooking.BookingStatusCancelled:
			if !bypassCapacity {
				if err := s.ensureCapacity(ctx, session); err != nil {
					return BookingDTO{}, err
				}
			}
			if err := existing.Rebook(now); err != nil {
				return BookingDTO{}, err
			}
			if err := s.repo.Save(ctx, existing); err != nil {
				return BookingDTO{}, fmt.Errorf("save rebooked booking: %w", err)
			}
			return toBookingDTO(existing), nil
		default:
			return BookingDTO{}, domainclassbooking.ErrInvalidStatus
		}
	}

	if !bypassCapacity {
		if err := s.ensureCapacity(ctx, session); err != nil {
			return BookingDTO{}, err
		}
	}

	booking, err := domainclassbooking.NewClassBooking(
		domainclassbooking.ClassBookingID(s.idgen.NewID()),
		session.ID(),
		membershipID,
		now,
	)
	if err != nil {
		return BookingDTO{}, err
	}

	if err := s.repo.Save(ctx, booking); err != nil {
		return BookingDTO{}, fmt.Errorf("save booking: %w", err)
	}

	return toBookingDTO(booking), nil
}

func (s *Service) ensureCapacity(ctx context.Context, session domainclasssession.ClassSession) error {
	count, err := s.repo.CountBookedBySession(ctx, session.ID())
	if err != nil {
		return fmt.Errorf("count booked: %w", err)
	}
	if count >= int(session.Capacity()) {
		return domainclassbooking.ErrClassFull
	}
	return nil
}

func toBookingDTO(booking domainclassbooking.ClassBooking) BookingDTO {
	return BookingDTO{
		ID:              booking.ID().String(),
		SessionID:       booking.SessionID().String(),
		GymMembershipID: booking.GymMembershipID().String(),
		Status:          string(booking.Status()),
		CreatedAt:       booking.CreatedAt(),
		UpdatedAt:       booking.UpdatedAt(),
	}
}
