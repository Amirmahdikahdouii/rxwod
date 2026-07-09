package classsession

import (
	"context"

	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type SessionBookingReader interface {
	CountBookedBySession(ctx context.Context, sessionID domainclasssession.ClassSessionID) (int, error)
	FindBySessionAndMembership(ctx context.Context, sessionID domainclasssession.ClassSessionID, membershipID gym.MembershipID) (domainclassbooking.ClassBooking, error)
}
