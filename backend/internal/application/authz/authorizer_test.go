package authz

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type fakeMembershipReader struct {
	membership gym.Membership
	err        error
}

func (f fakeMembershipReader) FindActiveMembership(context.Context, gym.GymID, user.UserID) (gym.Membership, error) {
	return f.membership, f.err
}

func TestAuthorizerResolveGymUsesActiveMembershipRole(t *testing.T) {
	membership, err := gym.NewMembership("membership-1", "gym-1", "user-1", domainauthz.RoleCoach, gym.MembershipStatusActive, nil, time.Now().UTC())
	if err != nil {
		t.Fatalf("new membership: %v", err)
	}
	authorizer := NewAuthorizer(fakeMembershipReader{membership: membership})

	principal, err := authorizer.ResolveGym(context.Background(), "user-1", "gym-1")
	if err != nil {
		t.Fatalf("resolve gym: %v", err)
	}
	if principal.Role != domainauthz.RoleCoach {
		t.Fatalf("expected coach role, got %s", principal.Role)
	}
}

func TestRequireRejectsMissingPermission(t *testing.T) {
	ctx := WithPrincipal(context.Background(), Principal{
		UserID: "user-1",
		GymID:  "gym-1",
		Role:   domainauthz.RoleAthlete,
	})

	_, err := Require(ctx, domainauthz.PermissionWODCreate)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
