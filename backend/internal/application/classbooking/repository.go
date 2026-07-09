package classbooking

import (
	"context"

	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type Repository interface {
	Save(ctx context.Context, booking domainclassbooking.ClassBooking) error
	FindByID(ctx context.Context, id domainclassbooking.ClassBookingID) (domainclassbooking.ClassBooking, error)
	FindBySessionAndMembership(ctx context.Context, sessionID domainclasssession.ClassSessionID, membershipID gym.MembershipID) (domainclassbooking.ClassBooking, error)
	ListBookingsBySession(ctx context.Context, sessionID domainclasssession.ClassSessionID) ([]domainclassbooking.ClassBooking, error)
	CountBookedBySession(ctx context.Context, sessionID domainclasssession.ClassSessionID) (int, error)
}
