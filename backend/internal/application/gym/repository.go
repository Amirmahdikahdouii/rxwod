package gym

import (
	"context"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type Repository interface {
	CreateGymWithOwner(ctx context.Context, gym domaingym.Gym, ownerMembership domaingym.Membership) error
	ListForUser(ctx context.Context, userID user.UserID) ([]WorkspaceDTO, error)
	FindByID(ctx context.Context, gymID domaingym.GymID) (domaingym.Gym, error)
	FindActiveMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error)
	FindMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error)
	FindMember(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (MemberDTO, error)
	ListMembers(ctx context.Context, gymID domaingym.GymID) ([]MemberDTO, error)
	DeleteMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) error
	FindUserByEmail(ctx context.Context, email user.Email) (user.User, error)
	FindUserByID(ctx context.Context, userID user.UserID) (user.User, error)
	UpsertMembership(ctx context.Context, membership domaingym.Membership) error
	SaveInvitation(ctx context.Context, invitation domaingym.Invitation, tokenHash string) error
	FindPendingInvitationsByEmail(ctx context.Context, email user.Email) ([]domaingym.Invitation, error)
	FindPendingInvitationByTokenHash(ctx context.Context, gymID domaingym.GymID, tokenHash string) (domaingym.Invitation, error)
	AcceptInvitationWithMembership(ctx context.Context, invitation domaingym.Invitation, membership domaingym.Membership) error
}

type InviteCommand struct {
	Email string
	Role  domainauthz.Role
}

type CreateGymCommand struct {
	Name string
}

type WorkspaceDTO struct {
	ID   string
	Name string
	Role domainauthz.Role
}

type MemberDTO struct {
	UserID      string
	Email       string
	DisplayName string
	Role        domainauthz.Role
	Status      domaingym.MembershipStatus
}

type GymDTO struct {
	ID      string
	Name    string
	OwnerID string
}

type InvitationDTO struct {
	ID    string
	GymID string
	Email string
	Role  domainauthz.Role
	Token string
}
