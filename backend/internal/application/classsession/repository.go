package classsession

import (
	"context"
	"time"

	domainclassbooking "github.com/rxwod/backend/internal/domain/classbooking"
	domainclasssession "github.com/rxwod/backend/internal/domain/classsession"
	"github.com/rxwod/backend/internal/domain/gym"
)

type Repository interface {
	Save(ctx context.Context, session domainclasssession.ClassSession) error
	SaveWithDefaultBookings(ctx context.Context, session domainclasssession.ClassSession, bookings []domainclassbooking.ClassBooking) error
	FindByID(ctx context.Context, gymID gym.GymID, id domainclasssession.ClassSessionID) (domainclasssession.ClassSession, error)
	ListByGymAndDate(ctx context.Context, gymID gym.GymID, from, to time.Time) ([]domainclasssession.ClassSession, error)
	Delete(ctx context.Context, gymID gym.GymID, id domainclasssession.ClassSessionID) error
}
