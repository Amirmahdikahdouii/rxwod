package wod

import (
	"context"
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type Repository interface {
	Save(ctx context.Context, wod domainwod.WOD) error
	FindByID(ctx context.Context, gymID gym.GymID, id domainwod.WODID) (domainwod.WOD, error)
	List(ctx context.Context, gymID gym.GymID, filter ListFilter) (ListResult, error)
	ListCalendar(ctx context.Context, gymID gym.GymID, from, to time.Time, includeDrafts bool) ([]CalendarEntry, error)
	Delete(ctx context.Context, gymID gym.GymID, id domainwod.WODID) error
}

type CalendarEntry struct {
	ID            string
	Name          string
	Status        domainwod.WODStatus
	ScheduledDate time.Time
}

type ListFilter struct {
	Page          int
	Limit         int
	Status        *domainwod.WODStatus
	PublishedOnly bool
}

type ListResult struct {
	Items []domainwod.WOD
	Total int
}
