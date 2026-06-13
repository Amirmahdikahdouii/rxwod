package authz

import (
	"context"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type contextKey struct{}

type Principal struct {
	UserID user.UserID
	GymID  gym.GymID
	Role   domainauthz.Role
}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, contextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(contextKey{}).(Principal)
	return principal, ok
}

func CurrentUserID(ctx context.Context) (user.UserID, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal.UserID == "" {
		return "", ErrUnauthenticated
	}
	return principal.UserID, nil
}

func CurrentGymID(ctx context.Context) (gym.GymID, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal.GymID == "" {
		return "", ErrGymRequired
	}
	return principal.GymID, nil
}

func Require(ctx context.Context, permission domainauthz.Permission) (Principal, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal.UserID == "" {
		return Principal{}, ErrUnauthenticated
	}
	if principal.GymID == "" {
		return Principal{}, ErrGymRequired
	}
	if !domainauthz.HasPermission(principal.Role, permission) {
		return Principal{}, ErrForbidden
	}
	return principal, nil
}
