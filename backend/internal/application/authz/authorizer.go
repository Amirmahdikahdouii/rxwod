package authz

import (
	"context"
	"errors"
	"fmt"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type MembershipReader interface {
	FindActiveMembership(ctx context.Context, gymID gym.GymID, userID user.UserID) (gym.Membership, error)
}

type Authorizer struct {
	memberships MembershipReader
}

func NewAuthorizer(memberships MembershipReader) *Authorizer {
	return &Authorizer{memberships: memberships}
}

func (a *Authorizer) ResolveGym(ctx context.Context, userID user.UserID, gymID gym.GymID) (Principal, error) {
	if gymID == "" {
		return Principal{}, ErrGymRequired
	}
	membership, err := a.memberships.FindActiveMembership(ctx, gymID, userID)
	if err != nil {
		if errors.Is(err, ErrActiveMembershipMissing) {
			return Principal{}, err
		}
		return Principal{}, fmt.Errorf("find active membership: %w", err)
	}
	return Principal{
		UserID: userID,
		GymID:  membership.GymID(),
		Role:   membership.Role(),
	}, nil
}

func (a *Authorizer) Require(ctx context.Context, permission domainauthz.Permission) (Principal, error) {
	return Require(ctx, permission)
}
