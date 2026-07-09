package gym

import (
	"context"
	"errors"
	"testing"
	"time"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time { return f.now }

type sequentialIDGen struct {
	next string
}

func (s *sequentialIDGen) NewID() string {
	if s.next == "" {
		s.next = "generated-id"
	}
	return s.next
}

type fakeGymRepo struct {
	invitations []domaingym.Invitation
	memberships []domaingym.Membership
	memberDTOs  map[user.UserID]MemberDTO
	users       map[user.UserID]user.User
	tokenHashes map[string]string
	previews    map[string]InvitationPreviewDTO
}

func (f *fakeGymRepo) findMembership(gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error) {
	for _, membership := range f.memberships {
		if membership.GymID() == gymID && membership.UserID() == userID {
			return membership, nil
		}
	}
	return domaingym.Membership{}, ErrMemberNotFound
}

func (f *fakeGymRepo) CreateGymWithOwner(context.Context, domaingym.Gym, domaingym.Membership) error {
	return nil
}

func (f *fakeGymRepo) ListForUser(context.Context, user.UserID) ([]WorkspaceDTO, error) {
	return nil, nil
}

func (f *fakeGymRepo) FindByID(context.Context, domaingym.GymID) (domaingym.Gym, error) {
	return domaingym.Gym{}, nil
}

func (f *fakeGymRepo) FindActiveMembership(context.Context, domaingym.GymID, user.UserID) (domaingym.Membership, error) {
	return domaingym.Membership{}, nil
}

func (f *fakeGymRepo) FindMembership(_ context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error) {
	return f.findMembership(gymID, userID)
}

func (f *fakeGymRepo) FindMember(_ context.Context, gymID domaingym.GymID, userID user.UserID) (MemberDTO, error) {
	membership, err := f.findMembership(gymID, userID)
	if err != nil {
		return MemberDTO{}, err
	}
	dto, ok := f.memberDTOs[userID]
	if !ok {
		return MemberDTO{}, ErrMemberNotFound
	}
	dto.Role = membership.Role()
	dto.Status = membership.Status()
	return dto, nil
}

func (f *fakeGymRepo) DeleteMembership(_ context.Context, gymID domaingym.GymID, userID user.UserID) error {
	for i, membership := range f.memberships {
		if membership.GymID() == gymID && membership.UserID() == userID {
			f.memberships = append(f.memberships[:i], f.memberships[i+1:]...)
			return nil
		}
	}
	return ErrMemberNotFound
}

func (f *fakeGymRepo) ListMembers(context.Context, domaingym.GymID) ([]MemberDTO, error) {
	return nil, nil
}

func (f *fakeGymRepo) FindUserByEmail(context.Context, user.Email) (user.User, error) {
	return user.User{}, ErrInviteeNotFound
}

func (f *fakeGymRepo) FindUserByID(_ context.Context, userID user.UserID) (user.User, error) {
	found, ok := f.users[userID]
	if !ok {
		return user.User{}, ErrInviteeNotFound
	}
	return found, nil
}

func (f *fakeGymRepo) UpsertMembership(_ context.Context, membership domaingym.Membership) error {
	for i, existing := range f.memberships {
		if existing.GymID() == membership.GymID() && existing.UserID() == membership.UserID() {
			f.memberships[i] = membership
			return nil
		}
	}
	f.memberships = append(f.memberships, membership)
	return nil
}

func (f *fakeGymRepo) SaveInvitation(_ context.Context, invitation domaingym.Invitation, tokenHash string) error {
	if f.tokenHashes == nil {
		f.tokenHashes = map[string]string{}
	}
	f.tokenHashes[invitation.ID().String()] = tokenHash
	f.invitations = append(f.invitations, invitation)
	return nil
}

func (f *fakeGymRepo) FindPendingInvitationsByEmail(_ context.Context, email user.Email) ([]domaingym.Invitation, error) {
	var matches []domaingym.Invitation
	for _, invitation := range f.invitations {
		if invitation.Email() == email {
			matches = append(matches, invitation)
		}
	}
	return matches, nil
}

func (f *fakeGymRepo) FindInvitationPreviewByTokenHash(_ context.Context, tokenHash string) (InvitationPreviewDTO, error) {
	if f.previews != nil {
		if preview, ok := f.previews[tokenHash]; ok {
			return preview, nil
		}
	}
	for _, invitation := range f.invitations {
		if f.tokenHashes[invitation.ID().String()] == tokenHash {
			return InvitationPreviewDTO{
				GymID:     invitation.GymID().String(),
				GymName:   "Test Gym",
				Email:     string(invitation.Email()),
				Role:      invitation.Role(),
				Status:    invitation.Status(),
				ExpiresAt: invitation.ExpiresAt(),
			}, nil
		}
	}
	return InvitationPreviewDTO{}, ErrInvitationNotFound
}

func (f *fakeGymRepo) FindPendingInvitationByTokenHash(_ context.Context, gymID domaingym.GymID, tokenHash string) (domaingym.Invitation, error) {
	for _, invitation := range f.invitations {
		if invitation.GymID() != gymID || invitation.Status() != domaingym.InvitationStatusPending {
			continue
		}
		if f.tokenHashes[invitation.ID().String()] == tokenHash {
			return invitation, nil
		}
	}
	return domaingym.Invitation{}, ErrInvitationNotFound
}

func (f *fakeGymRepo) AcceptInvitationWithMembership(_ context.Context, invitation domaingym.Invitation, membership domaingym.Membership) error {
	f.invitations = []domaingym.Invitation{invitation}
	f.memberships = append(f.memberships, membership)
	return nil
}

func TestGetInvitationPreviewPending(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	token := "secret-token"
	repo := &fakeGymRepo{
		previews: map[string]InvitationPreviewDTO{
			hashInvitationToken(token): {
				GymID:     "gym-1",
				GymName:   "CrossFit Downtown",
				Email:     "athlete@example.com",
				Role:      domainauthz.RoleAthlete,
				Status:    domaingym.InvitationStatusPending,
				ExpiresAt: now.Add(time.Hour),
			},
		},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	preview, err := service.GetInvitationPreview(context.Background(), token)
	if err != nil {
		t.Fatalf("get invitation preview: %v", err)
	}
	if preview.GymName != "CrossFit Downtown" || preview.Status != domaingym.InvitationStatusPending {
		t.Fatalf("unexpected preview: %+v", preview)
	}
}

func TestGetInvitationPreviewExpiredByClock(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	token := "secret-token"
	repo := &fakeGymRepo{
		previews: map[string]InvitationPreviewDTO{
			hashInvitationToken(token): {
				GymID:     "gym-1",
				GymName:   "CrossFit Downtown",
				Email:     "athlete@example.com",
				Role:      domainauthz.RoleAthlete,
				Status:    domaingym.InvitationStatusPending,
				ExpiresAt: now.Add(-time.Hour),
			},
		},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	preview, err := service.GetInvitationPreview(context.Background(), token)
	if err != nil {
		t.Fatalf("get invitation preview: %v", err)
	}
	if preview.Status != domaingym.InvitationStatusExpired {
		t.Fatalf("expected expired status, got %s", preview.Status)
	}
}

func TestGetInvitationPreviewNotFound(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &fakeGymRepo{}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err := service.GetInvitationPreview(context.Background(), "unknown-token")
	if !errors.Is(err, ErrInvitationNotFound) {
		t.Fatalf("expected ErrInvitationNotFound, got %v", err)
	}
}

func TestAcceptInvitesForEmailCreatesActiveMembership(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	invitation, err := domaingym.NewInvitation(
		domaingym.InvitationID("invite-1"),
		domaingym.GymID("gym-1"),
		user.Email("athlete@example.com"),
		domainauthz.RoleAthlete,
		user.UserID("owner-1"),
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("new invitation: %v", err)
	}

	repo := &fakeGymRepo{invitations: []domaingym.Invitation{invitation}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{next: "membership-1"}, time.Hour)

	if err := service.AcceptInvitesForEmail(context.Background(), user.Email("athlete@example.com"), user.UserID("athlete-1")); err != nil {
		t.Fatalf("accept invites: %v", err)
	}
	if len(repo.memberships) != 1 {
		t.Fatalf("expected one membership, got %d", len(repo.memberships))
	}
	if repo.memberships[0].Role() != domainauthz.RoleAthlete || repo.memberships[0].Status() != domaingym.MembershipStatusActive {
		t.Fatalf("unexpected membership: %+v", repo.memberships[0])
	}
	if repo.invitations[0].Status() != domaingym.InvitationStatusAccepted {
		t.Fatalf("expected invitation accepted, got %s", repo.invitations[0].Status())
	}
}

func acceptorContext() context.Context {
	return appauthz.WithPrincipal(context.Background(), appauthz.Principal{
		UserID: user.UserID("athlete-1"),
	})
}

func TestAcceptInvitationSuccess(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	invitation, err := domaingym.NewInvitation(
		domaingym.InvitationID("invite-1"),
		domaingym.GymID("gym-1"),
		user.Email("athlete@example.com"),
		domainauthz.RoleAthlete,
		user.UserID("owner-1"),
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("new invitation: %v", err)
	}
	athlete, err := user.NewUser(user.UserID("athlete-1"), user.Email("athlete@example.com"), "hash", "Athlete", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}

	repo := &fakeGymRepo{
		invitations: []domaingym.Invitation{invitation},
		tokenHashes: map[string]string{"invite-1": hashInvitationToken("secret-token")},
		users:       map[user.UserID]user.User{"athlete-1": athlete},
		memberDTOs: map[user.UserID]MemberDTO{
			user.UserID("athlete-1"): {
				UserID:      "athlete-1",
				Email:       "athlete@example.com",
				DisplayName: "Athlete",
			},
		},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{next: "membership-1"}, time.Hour)

	result, err := service.AcceptInvitation(acceptorContext(), "gym-1", "secret-token")
	if err != nil {
		t.Fatalf("accept invitation: %v", err)
	}
	if result.Role != domainauthz.RoleAthlete || result.Status != domaingym.MembershipStatusActive {
		t.Fatalf("unexpected member result: %+v", result)
	}
	if len(repo.memberships) != 1 || repo.memberships[0].Status() != domaingym.MembershipStatusActive {
		t.Fatalf("expected one active membership, got %+v", repo.memberships)
	}
	if repo.invitations[0].Status() != domaingym.InvitationStatusAccepted {
		t.Fatalf("expected invitation accepted, got %s", repo.invitations[0].Status())
	}
}

func TestAcceptInvitationInvalidToken(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	athlete, err := user.NewUser(user.UserID("athlete-1"), user.Email("athlete@example.com"), "hash", "Athlete", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	repo := &fakeGymRepo{users: map[user.UserID]user.User{"athlete-1": athlete}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err = service.AcceptInvitation(acceptorContext(), "gym-1", "unknown-token")
	if !errors.Is(err, ErrInvitationNotFound) {
		t.Fatalf("expected ErrInvitationNotFound, got %v", err)
	}
}

func TestAcceptInvitationExpired(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	invitation, err := domaingym.NewInvitation(
		domaingym.InvitationID("invite-1"),
		domaingym.GymID("gym-1"),
		user.Email("athlete@example.com"),
		domainauthz.RoleAthlete,
		user.UserID("owner-1"),
		now.Add(-time.Hour),
		now.Add(-2*time.Hour),
	)
	if err != nil {
		t.Fatalf("new invitation: %v", err)
	}
	athlete, err := user.NewUser(user.UserID("athlete-1"), user.Email("athlete@example.com"), "hash", "Athlete", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	repo := &fakeGymRepo{
		invitations: []domaingym.Invitation{invitation},
		tokenHashes: map[string]string{"invite-1": hashInvitationToken("secret-token")},
		users:       map[user.UserID]user.User{"athlete-1": athlete},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err = service.AcceptInvitation(acceptorContext(), "gym-1", "secret-token")
	if !errors.Is(err, domaingym.ErrInvitationExpired) {
		t.Fatalf("expected ErrInvitationExpired, got %v", err)
	}
}

func TestAcceptInvitationEmailMismatch(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	invitation, err := domaingym.NewInvitation(
		domaingym.InvitationID("invite-1"),
		domaingym.GymID("gym-1"),
		user.Email("someone-else@example.com"),
		domainauthz.RoleAthlete,
		user.UserID("owner-1"),
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("new invitation: %v", err)
	}
	athlete, err := user.NewUser(user.UserID("athlete-1"), user.Email("athlete@example.com"), "hash", "Athlete", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	repo := &fakeGymRepo{
		invitations: []domaingym.Invitation{invitation},
		tokenHashes: map[string]string{"invite-1": hashInvitationToken("secret-token")},
		users:       map[user.UserID]user.User{"athlete-1": athlete},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err = service.AcceptInvitation(acceptorContext(), "gym-1", "secret-token")
	if !errors.Is(err, ErrInvitationEmailMismatch) {
		t.Fatalf("expected ErrInvitationEmailMismatch, got %v", err)
	}
}

func ownerContext(gymID string) context.Context {
	return appauthz.WithPrincipal(context.Background(), appauthz.Principal{
		UserID: user.UserID("owner-1"),
		GymID:  domaingym.GymID(gymID),
		Role:   domainauthz.RoleOwner,
	})
}

func TestUpdateMemberRoleCoachToAthlete(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	membership, err := domaingym.NewMembership(
		domaingym.MembershipID("membership-1"),
		domaingym.GymID("gym-1"),
		user.UserID("coach-1"),
		domainauthz.RoleCoach,
		domaingym.MembershipStatusActive,
		nil,
		now,
	)
	if err != nil {
		t.Fatalf("new membership: %v", err)
	}

	repo := &fakeGymRepo{
		memberships: []domaingym.Membership{membership},
		memberDTOs: map[user.UserID]MemberDTO{
			user.UserID("coach-1"): {
				UserID:      "coach-1",
				Email:       "coach@example.com",
				DisplayName: "Coach",
				Role:        domainauthz.RoleCoach,
				Status:      domaingym.MembershipStatusActive,
			},
		},
	}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	result, err := service.UpdateMemberRole(ownerContext("gym-1"), "gym-1", "coach-1", domainauthz.RoleAthlete)
	if err != nil {
		t.Fatalf("update member role: %v", err)
	}
	if result.Role != domainauthz.RoleAthlete {
		t.Fatalf("expected athlete role, got %s", result.Role)
	}
	if repo.memberships[0].Role() != domainauthz.RoleAthlete {
		t.Fatalf("expected persisted athlete role, got %s", repo.memberships[0].Role())
	}
}

func TestUpdateMemberRoleRejectsOwner(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	membership, err := domaingym.NewOwnerMembership(
		domaingym.MembershipID("membership-owner"),
		domaingym.GymID("gym-1"),
		user.UserID("owner-1"),
		now,
	)
	if err != nil {
		t.Fatalf("new owner membership: %v", err)
	}

	repo := &fakeGymRepo{memberships: []domaingym.Membership{membership}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err = service.UpdateMemberRole(ownerContext("gym-1"), "gym-1", "owner-1", domainauthz.RoleCoach)
	if !errors.Is(err, ErrOwnerMembershipProtected) {
		t.Fatalf("expected ErrOwnerMembershipProtected, got %v", err)
	}
}

func TestUpdateMemberRoleRejectsInvalidRole(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	membership, err := domaingym.NewMembership(
		domaingym.MembershipID("membership-1"),
		domaingym.GymID("gym-1"),
		user.UserID("coach-1"),
		domainauthz.RoleCoach,
		domaingym.MembershipStatusActive,
		nil,
		now,
	)
	if err != nil {
		t.Fatalf("new membership: %v", err)
	}

	repo := &fakeGymRepo{memberships: []domaingym.Membership{membership}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	_, err = service.UpdateMemberRole(ownerContext("gym-1"), "gym-1", "coach-1", domainauthz.RoleOwner)
	if !errors.Is(err, ErrRoleNotAssignable) {
		t.Fatalf("expected ErrRoleNotAssignable, got %v", err)
	}
}

func TestRemoveMemberRemovesCoach(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	membership, err := domaingym.NewMembership(
		domaingym.MembershipID("membership-1"),
		domaingym.GymID("gym-1"),
		user.UserID("coach-1"),
		domainauthz.RoleCoach,
		domaingym.MembershipStatusActive,
		nil,
		now,
	)
	if err != nil {
		t.Fatalf("new membership: %v", err)
	}

	repo := &fakeGymRepo{memberships: []domaingym.Membership{membership}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	if err := service.RemoveMember(ownerContext("gym-1"), "gym-1", "coach-1"); err != nil {
		t.Fatalf("remove member: %v", err)
	}
	if len(repo.memberships) != 0 {
		t.Fatalf("expected membership removed, got %d", len(repo.memberships))
	}
}

func TestRemoveMemberRejectsOwner(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	membership, err := domaingym.NewOwnerMembership(
		domaingym.MembershipID("membership-owner"),
		domaingym.GymID("gym-1"),
		user.UserID("owner-1"),
		now,
	)
	if err != nil {
		t.Fatalf("new owner membership: %v", err)
	}

	repo := &fakeGymRepo{memberships: []domaingym.Membership{membership}}
	service := NewService(repo, fixedClock{now: now}, &sequentialIDGen{}, time.Hour)

	err = service.RemoveMember(ownerContext("gym-1"), "gym-1", "owner-1")
	if !errors.Is(err, ErrOwnerMembershipProtected) {
		t.Fatalf("expected ErrOwnerMembershipProtected, got %v", err)
	}
}
