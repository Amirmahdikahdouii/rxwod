package wod

import (
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type MovementDTO struct {
	ID        string
	Position  int
	Name      string
	Reps      *int
	LoadValue *float64
	LoadUnit  *string
	Notes     string
}

type WODSummaryDTO struct {
	ID          string
	Name        string
	Type        domainwod.WODType
	Status      domainwod.WODStatus
	ScoringKind domainwod.ScoringKind
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WODDetailDTO struct {
	ID          string
	Name        string
	Description string
	Type        domainwod.WODType
	Status      domainwod.WODStatus
	ScoringKind domainwod.ScoringKind
	Config      ConfigDTO
	Movements   []MovementDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ConfigDTO struct {
	TimeCapSeconds    *int
	Rounds            *int
	WorkSeconds       *int
	RestSeconds       *int
	Cycles            *int
	IntervalSeconds   *int
}

type CreateWODResultDTO struct {
	ID          string
	Name        string
	Type        domainwod.WODType
	Status      domainwod.WODStatus
	ScoringKind domainwod.ScoringKind
}
