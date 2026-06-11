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

type ConfigDTO struct {
	TimeCapSeconds  *int
	Rounds          *int
	WorkSeconds     *int
	RestSeconds     *int
	Cycles          *int
	IntervalSeconds *int
}

type StageDTO struct {
	ID          string
	Kind        domainwod.StageKind
	Position    int
	Type        domainwod.WODType
	ScoringKind domainwod.ScoringKind
	Config      ConfigDTO
	Movements   []MovementDTO
}

type StageSummaryDTO struct {
	Kind        domainwod.StageKind
	Position    int
	Type        domainwod.WODType
	ScoringKind domainwod.ScoringKind
}

type WODSummaryDTO struct {
	ID         string
	Name       string
	Status     domainwod.WODStatus
	StageCount int
	Stages     []StageSummaryDTO
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type WODDetailDTO struct {
	ID          string
	Name        string
	Description string
	Status      domainwod.WODStatus
	Stages      []StageDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateWODResultDTO struct {
	ID         string
	Name       string
	Status     domainwod.WODStatus
	StageCount int
	Stages     []StageSummaryDTO
}
