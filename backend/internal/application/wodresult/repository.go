package wodresult

import (
	"context"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	domainwodresult "github.com/rxwod/backend/internal/domain/wodresult"
)

type Repository interface {
	Save(ctx context.Context, result domainwodresult.WODResult) error
	FindByWODAndMembership(ctx context.Context, wodID domainwod.WODID, gymMembershipID gym.MembershipID) (domainwodresult.WODResult, error)
	ListByWOD(ctx context.Context, wodID domainwod.WODID) ([]domainwodresult.WODResult, error)
}
