package wod

import (
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type MovementInput struct {
	Position  int
	Name      string
	Reps      *int
	LoadValue *float64
	LoadUnit  *string
	Notes     string
}

type CreateAMRAPCommand struct {
	Name        string
	Description string
	TimeCap     int
	Movements   []MovementInput
}

type CreateForTimeCommand struct {
	Name        string
	Description string
	Rounds      int
	TimeCap     *int
	Movements   []MovementInput
}

type CreateTabataCommand struct {
	Name        string
	Description string
	WorkSeconds int
	RestSeconds int
	Rounds      int
	Cycles      int
	Movements   []MovementInput
}

type CreateEMOMCommand struct {
	Name             string
	Description      string
	IntervalSeconds  int
	Rounds           int
	Movements        []MovementInput
}

type CreateCommand struct {
	Type domainwod.WODType
	AMRAP   *CreateAMRAPCommand
	ForTime *CreateForTimeCommand
	Tabata  *CreateTabataCommand
	EMOM    *CreateEMOMCommand
}
