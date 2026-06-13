package wod

import (
	"context"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type Repository interface {
	Save(ctx context.Context, wod domainwod.WOD) error
	FindByID(ctx context.Context, gymID gym.GymID, id domainwod.WODID) (domainwod.WOD, error)
	List(ctx context.Context, gymID gym.GymID) ([]domainwod.WOD, error)
}
