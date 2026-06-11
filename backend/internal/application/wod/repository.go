package wod

import (
	"context"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type Repository interface {
	Save(ctx context.Context, variant domainwod.Variant) error
	FindByID(ctx context.Context, id domainwod.WODID) (domainwod.Variant, error)
	List(ctx context.Context) ([]domainwod.Variant, error)
}
