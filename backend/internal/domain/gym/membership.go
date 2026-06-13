package gym

import (
	"time"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/user"
)

type Membership struct {
	id        MembershipID
	gymID     GymID
	userID    user.UserID
	role      domainauthz.Role
	status    MembershipStatus
	invitedBy *user.UserID
	createdAt time.Time
	updatedAt time.Time
}

func NewMembership(
	id MembershipID,
	gymID GymID,
	userID user.UserID,
	role domainauthz.Role,
	status MembershipStatus,
	invitedBy *user.UserID,
	now time.Time,
) (Membership, error) {
	if !domainauthz.IsValidRole(role) {
		return Membership{}, ErrInvalidInvitationRole
	}
	if err := validateMembershipStatus(status); err != nil {
		return Membership{}, err
	}
	return Membership{
		id:        id,
		gymID:     gymID,
		userID:    userID,
		role:      role,
		status:    status,
		invitedBy: cloneUserID(invitedBy),
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructMembership(
	id MembershipID,
	gymID GymID,
	userID user.UserID,
	role domainauthz.Role,
	status MembershipStatus,
	invitedBy *user.UserID,
	createdAt time.Time,
	updatedAt time.Time,
) (Membership, error) {
	membership, err := NewMembership(id, gymID, userID, role, status, invitedBy, createdAt)
	if err != nil {
		return Membership{}, err
	}
	membership.updatedAt = updatedAt
	return membership, nil
}

func NewOwnerMembership(id MembershipID, gymID GymID, ownerID user.UserID, now time.Time) (Membership, error) {
	return NewMembership(id, gymID, ownerID, domainauthz.RoleOwner, MembershipStatusActive, nil, now)
}

func (m Membership) ID() MembershipID {
	return m.id
}

func (m Membership) GymID() GymID {
	return m.gymID
}

func (m Membership) UserID() user.UserID {
	return m.userID
}

func (m Membership) Role() domainauthz.Role {
	return m.role
}

func (m Membership) Status() MembershipStatus {
	return m.status
}

func (m Membership) InvitedBy() *user.UserID {
	return cloneUserID(m.invitedBy)
}

func (m Membership) CreatedAt() time.Time {
	return m.createdAt
}

func (m Membership) UpdatedAt() time.Time {
	return m.updatedAt
}

func cloneUserID(id *user.UserID) *user.UserID {
	if id == nil {
		return nil
	}
	value := *id
	return &value
}
