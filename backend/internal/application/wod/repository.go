package wod

import (
	"context"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type Repository interface {
	Save(ctx context.Context, wod domainwod.WOD) error
	FindByID(ctx context.Context, id domainwod.WODID) (domainwod.WOD, error)
	List(ctx context.Context) ([]domainwod.WOD, error)
}
