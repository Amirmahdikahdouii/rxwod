package gym

import (
	"strings"
	"unicode/utf8"

	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/user"
)

type GymID string
type MembershipID string
type InvitationID string
type GymName string

type MembershipStatus string

const (
	MembershipStatusPending MembershipStatus = "pending"
	MembershipStatusActive  MembershipStatus = "active"
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

func (id GymID) String() string {
	return string(id)
}

func (id MembershipID) String() string {
	return string(id)
}

func (id InvitationID) String() string {
	return string(id)
}

func validateGymName(name GymName) error {
	length := utf8.RuneCountInString(strings.TrimSpace(string(name)))
	if length < 2 || length > 120 {
		return ErrInvalidGymName
	}
	return nil
}

func validateMembershipStatus(status MembershipStatus) error {
	switch status {
	case MembershipStatusPending, MembershipStatusActive:
		return nil
	default:
		return ErrInvalidMembershipStatus
	}
}

func validateInvitationStatus(status InvitationStatus) error {
	switch status {
	case InvitationStatusPending, InvitationStatusAccepted, InvitationStatusExpired, InvitationStatusRevoked:
		return nil
	default:
		return ErrInvalidInvitationStatus
	}
}

func validateAssignableRole(role domainauthz.Role) error {
	switch role {
	case domainauthz.RoleCoach, domainauthz.RoleAthlete:
		return nil
	default:
		return ErrInvalidInvitationRole
	}
}

type MemberUser struct {
	ID          user.UserID
	Email       user.Email
	DisplayName user.DisplayName
}
