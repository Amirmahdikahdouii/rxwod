package gym

import (
	"time"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/user"
)

type Invitation struct {
	id        InvitationID
	gymID     GymID
	email     user.Email
	role      domainauthz.Role
	status    InvitationStatus
	invitedBy user.UserID
	expiresAt time.Time
	createdAt time.Time
	updatedAt time.Time
}

func NewInvitation(
	id InvitationID,
	gymID GymID,
	email user.Email,
	role domainauthz.Role,
	invitedBy user.UserID,
	expiresAt time.Time,
	now time.Time,
) (Invitation, error) {
	return ReconstructInvitation(id, gymID, email, role, InvitationStatusPending, invitedBy, expiresAt, now, now)
}

func ReconstructInvitation(
	id InvitationID,
	gymID GymID,
	email user.Email,
	role domainauthz.Role,
	status InvitationStatus,
	invitedBy user.UserID,
	expiresAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (Invitation, error) {
	if err := validateAssignableRole(role); err != nil {
		return Invitation{}, err
	}
	if err := validateInvitationStatus(status); err != nil {
		return Invitation{}, err
	}
	return Invitation{
		id:        id,
		gymID:     gymID,
		email:     email,
		role:      role,
		status:    status,
		invitedBy: invitedBy,
		expiresAt: expiresAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (i Invitation) ID() InvitationID {
	return i.id
}

func (i Invitation) GymID() GymID {
	return i.gymID
}

func (i Invitation) Email() user.Email {
	return i.email
}

func (i Invitation) Role() domainauthz.Role {
	return i.role
}

func (i Invitation) Status() InvitationStatus {
	return i.status
}

func (i Invitation) InvitedBy() user.UserID {
	return i.invitedBy
}

func (i Invitation) ExpiresAt() time.Time {
	return i.expiresAt
}

func (i Invitation) CreatedAt() time.Time {
	return i.createdAt
}

func (i Invitation) UpdatedAt() time.Time {
	return i.updatedAt
}
