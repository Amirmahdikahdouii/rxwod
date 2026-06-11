package wod

import (
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type MovementInput struct {
	Position     int
	Label        string
	Name         string
	Prescription string
	Reps         *int
	LoadValue    *float64
	LoadUnit     *string
	Notes        string
}

type StageConfigInput struct {
	Type            domainwod.WODType
	TimeCapSeconds  *int
	Rounds          *int
	WorkSeconds     *int
	RestSeconds     *int
	Cycles          *int
	IntervalSeconds *int
}

type StageInput struct {
	Kind         domainwod.StageKind
	Instructions string
	Config       StageConfigInput
	Movements    []MovementInput
}

type CreateWODCommand struct {
	Name        string
	Description string
	Stages      []StageInput
}
