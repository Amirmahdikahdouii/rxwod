package gym

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo      Repository
	users     appauth.UserRepository
	clock     clock.Clock
	idgen     idgen.Generator
	inviteTTL time.Duration
}

func NewService(repo Repository, users appauth.UserRepository, clock clock.Clock, idgen idgen.Generator, inviteTTL time.Duration) *Service {
	return &Service{repo: repo, users: users, clock: clock, idgen: idgen, inviteTTL: inviteTTL}
}

func (s *Service) requireVerifiedEmail(ctx context.Context) error {
	userID, err := appauthz.CurrentUserID(ctx)
	if err != nil {
		return err
	}
	aggregate, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !aggregate.IsEmailVerified() {
		return appauth.ErrEmailNotVerified
	}
	return nil
}

func (s *Service) Create(ctx context.Context, cmd CreateGymCommand) (GymDTO, error) {
	if err := s.requireVerifiedEmail(ctx); err != nil {
		return GymDTO{}, err
	}
	userID, err := appauthz.CurrentUserID(ctx)
	if err != nil {
		return GymDTO{}, err
	}
	now := s.clock.Now()
	aggregate, err := domaingym.NewGym(domaingym.GymID(s.idgen.NewID()), domaingym.GymName(cmd.Name), userID, now)
	if err != nil {
		return GymDTO{}, err
	}
	membership, err := domaingym.NewOwnerMembership(domaingym.MembershipID(s.idgen.NewID()), aggregate.ID(), userID, now)
	if err != nil {
		return GymDTO{}, err
	}
	if err := s.repo.CreateGymWithOwner(ctx, aggregate, membership); err != nil {
		return GymDTO{}, fmt.Errorf("create gym with owner: %w", err)
	}
	return toGymDTO(aggregate), nil
}

func (s *Service) ListForCurrentUser(ctx context.Context) ([]WorkspaceDTO, error) {
	userID, err := appauthz.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListForUser(ctx, userID)
}

func (s *Service) Get(ctx context.Context, gymID string) (GymDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionGymRead)
	if err != nil {
		return GymDTO{}, err
	}
	if principal.GymID.String() != gymID {
		return GymDTO{}, appauthz.ErrGymMismatch
	}
	aggregate, err := s.repo.FindByID(ctx, domaingym.GymID(gymID))
	if err != nil {
		return GymDTO{}, err
	}
	return toGymDTO(aggregate), nil
}

func (s *Service) ListMembers(ctx context.Context, gymID string) ([]MemberDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionMemberList)
	if err != nil {
		return nil, err
	}
	if principal.GymID.String() != gymID {
		return nil, appauthz.ErrGymMismatch
	}
	return s.repo.ListMembers(ctx, principal.GymID)
}

func (s *Service) UpdateMemberRole(ctx context.Context, gymID string, userID string, role domainauthz.Role) (MemberDTO, error) {
	if err := s.requireVerifiedEmail(ctx); err != nil {
		return MemberDTO{}, err
	}
	principal, err := appauthz.Require(ctx, domainauthz.PermissionMemberUpdateRole)
	if err != nil {
		return MemberDTO{}, err
	}
	if principal.GymID.String() != gymID {
		return MemberDTO{}, appauthz.ErrGymMismatch
	}
	if role != domainauthz.RoleCoach && role != domainauthz.RoleAthlete {
		return MemberDTO{}, ErrRoleNotAssignable
	}

	membership, err := s.guardManageableMembership(ctx, principal.GymID, user.UserID(userID))
	if err != nil {
		return MemberDTO{}, err
	}

	now := s.clock.Now()
	updated, err := domaingym.ReconstructMembership(
		membership.ID(),
		membership.GymID(),
		membership.UserID(),
		role,
		membership.Status(),
		membership.InvitedBy(),
		membership.CreatedAt(),
		now,
	)
	if err != nil {
		return MemberDTO{}, err
	}
	if err := s.repo.UpsertMembership(ctx, updated); err != nil {
		return MemberDTO{}, fmt.Errorf("update member role: %w", err)
	}
	return s.repo.FindMember(ctx, principal.GymID, user.UserID(userID))
}

func (s *Service) RemoveMember(ctx context.Context, gymID string, userID string) error {
	if err := s.requireVerifiedEmail(ctx); err != nil {
		return err
	}
	principal, err := appauthz.Require(ctx, domainauthz.PermissionMemberRemove)
	if err != nil {
		return err
	}
	if principal.GymID.String() != gymID {
		return appauthz.ErrGymMismatch
	}

	if _, err := s.guardManageableMembership(ctx, principal.GymID, user.UserID(userID)); err != nil {
		return err
	}
	if err := s.repo.DeleteMembership(ctx, principal.GymID, user.UserID(userID)); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (s *Service) guardManageableMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error) {
	membership, err := s.repo.FindMembership(ctx, gymID, userID)
	if err != nil {
		return domaingym.Membership{}, err
	}
	if membership.Role() == domainauthz.RoleOwner {
		return domaingym.Membership{}, ErrOwnerMembershipProtected
	}
	return membership, nil
}

func (s *Service) InviteCoach(ctx context.Context, gymID string, email string) (InvitationDTO, error) {
	return s.invite(ctx, gymID, InviteCommand{Email: email, Role: domainauthz.RoleCoach}, domainauthz.PermissionMemberInviteCoach)
}

func (s *Service) InviteAthlete(ctx context.Context, gymID string, email string) (InvitationDTO, error) {
	return s.invite(ctx, gymID, InviteCommand{Email: email, Role: domainauthz.RoleAthlete}, domainauthz.PermissionMemberInviteAthlete)
}

func (s *Service) AcceptInvitesForEmail(ctx context.Context, email user.Email, userID user.UserID) error {
	invitations, err := s.repo.FindPendingInvitationsByEmail(ctx, email)
	if err != nil {
		return err
	}
	now := s.clock.Now()
	for _, invitation := range invitations {
		if _, err := s.acceptInvitation(ctx, invitation, userID, now); err != nil {
			if errors.Is(err, domaingym.ErrInvitationExpired) {
				continue
			}
			return err
		}
	}
	return nil
}

func (s *Service) GetInvitationPreview(ctx context.Context, token string) (InvitationPreviewDTO, error) {
	preview, err := s.repo.FindInvitationPreviewByTokenHash(ctx, hashInvitationToken(token))
	if err != nil {
		return InvitationPreviewDTO{}, err
	}
	if preview.Status == domaingym.InvitationStatusPending && !preview.ExpiresAt.After(s.clock.Now()) {
		preview.Status = domaingym.InvitationStatusExpired
	}
	return preview, nil
}

func (s *Service) AcceptInvitation(ctx context.Context, gymID string, token string) (MemberDTO, error) {
	userID, err := appauthz.CurrentUserID(ctx)
	if err != nil {
		return MemberDTO{}, err
	}
	currentUser, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return MemberDTO{}, err
	}
	invitation, err := s.repo.FindPendingInvitationByTokenHash(ctx, domaingym.GymID(gymID), hashInvitationToken(token))
	if err != nil {
		return MemberDTO{}, err
	}
	if invitation.Email() != currentUser.Email() {
		return MemberDTO{}, ErrInvitationEmailMismatch
	}
	now := s.clock.Now()
	if _, err := s.acceptInvitation(ctx, invitation, userID, now); err != nil {
		return MemberDTO{}, err
	}
	return s.repo.FindMember(ctx, invitation.GymID(), userID)
}

func (s *Service) acceptInvitation(ctx context.Context, invitation domaingym.Invitation, userID user.UserID, now time.Time) (domaingym.Membership, error) {
	accepted, err := invitation.Accept(now)
	if err != nil {
		return domaingym.Membership{}, err
	}
	membership, err := domaingym.NewMembership(
		domaingym.MembershipID(s.idgen.NewID()),
		accepted.GymID(),
		userID,
		accepted.Role(),
		domaingym.MembershipStatusActive,
		userIDPtr(accepted.InvitedBy()),
		now,
	)
	if err != nil {
		return domaingym.Membership{}, err
	}
	if err := s.repo.AcceptInvitationWithMembership(ctx, accepted, membership); err != nil {
		return domaingym.Membership{}, fmt.Errorf("accept invitation: %w", err)
	}
	return membership, nil
}

func (s *Service) invite(ctx context.Context, gymID string, cmd InviteCommand, permission domainauthz.Permission) (InvitationDTO, error) {
	if err := s.requireVerifiedEmail(ctx); err != nil {
		return InvitationDTO{}, err
	}
	principal, err := appauthz.Require(ctx, permission)
	if err != nil {
		return InvitationDTO{}, err
	}
	if principal.GymID.String() != gymID {
		return InvitationDTO{}, appauthz.ErrGymMismatch
	}
	if cmd.Role != domainauthz.RoleCoach && cmd.Role != domainauthz.RoleAthlete {
		return InvitationDTO{}, ErrRoleNotAssignable
	}

	email := user.NormalizeEmail(cmd.Email)
	now := s.clock.Now()
	invitation, err := domaingym.NewInvitation(
		domaingym.InvitationID(s.idgen.NewID()),
		principal.GymID,
		email,
		cmd.Role,
		principal.UserID,
		now.Add(s.inviteTTL),
		now,
	)
	if err != nil {
		return InvitationDTO{}, err
	}

	token := s.idgen.NewID()
	if err := s.repo.SaveInvitation(ctx, invitation, hashInvitationToken(token)); err != nil {
		return InvitationDTO{}, fmt.Errorf("save invitation: %w", err)
	}
	return InvitationDTO{
		ID:    invitation.ID().String(),
		GymID: invitation.GymID().String(),
		Email: string(invitation.Email()),
		Role:  invitation.Role(),
		Token: token,
	}, nil
}

func hashInvitationToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func toGymDTO(aggregate domaingym.Gym) GymDTO {
	return GymDTO{
		ID:      aggregate.ID().String(),
		Name:    string(aggregate.Name()),
		OwnerID: aggregate.OwnerID().String(),
	}
}

func userIDPtr(id user.UserID) *user.UserID {
	return &id
}
