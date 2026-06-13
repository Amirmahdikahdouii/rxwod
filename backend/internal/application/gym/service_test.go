package gym

import (
	"context"
	"testing"
	"time"

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

func (f *fakeGymRepo) ListMembers(context.Context, domaingym.GymID) ([]MemberDTO, error) {
	return nil, nil
}

func (f *fakeGymRepo) FindUserByEmail(context.Context, user.Email) (user.User, error) {
	return user.User{}, ErrInviteeNotFound
}

func (f *fakeGymRepo) UpsertMembership(context.Context, domaingym.Membership) error {
	return nil
}

func (f *fakeGymRepo) SaveInvitation(context.Context, domaingym.Invitation) error {
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

func (f *fakeGymRepo) AcceptInvitationWithMembership(_ context.Context, invitation domaingym.Invitation, membership domaingym.Membership) error {
	f.invitations = []domaingym.Invitation{invitation}
	f.memberships = append(f.memberships, membership)
	return nil
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
