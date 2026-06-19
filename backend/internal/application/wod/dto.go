package wod

import (
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type MovementDTO struct {
	ID           string
	Position     int
	Label        string
	Name         string
	Prescription string
	Sets         *int
	Reps         *int
	LoadValue    *float64
	LoadUnit     *string
	Notes        string
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
	ID           string
	Kind         domainwod.StageKind
	Position     int
	Instructions string
	Type         domainwod.WODType
	ScoringKind  domainwod.ScoringKind
	Config       ConfigDTO
	Movements    []MovementDTO
}

type StageSummaryDTO struct {
	Kind        domainwod.StageKind
	Position    int
	Type        domainwod.WODType
	ScoringKind domainwod.ScoringKind
}

type WODSummaryDTO struct {
	ID            string
	Name          string
	Status        domainwod.WODStatus
	StageCount    int
	Stages        []StageSummaryDTO
	CreatedBy     string
	ScheduledDate *time.Time
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type WODDetailDTO struct {
	ID            string
	Name          string
	Description   string
	Status        domainwod.WODStatus
	Stages        []StageDTO
	CreatedBy     string
	ScheduledDate *time.Time
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateWODResultDTO struct {
	ID            string
	Name          string
	Status        domainwod.WODStatus
	StageCount    int
	Stages        []StageSummaryDTO
	ScheduledDate *time.Time
}

type CalendarPlanDTO struct {
	ID     string
	Name   string
	Status domainwod.WODStatus
}

type CalendarDayDTO struct {
	Date           string
	PublishedCount int
	DraftCount     int
	Plans          []CalendarPlanDTO
}
