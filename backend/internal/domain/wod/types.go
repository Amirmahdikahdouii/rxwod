package wod

import (
	"strings"
	"unicode/utf8"
)

type WODID string
type StageID string
type MovementID string
type WODName string
type WODDescription string
type TimeCapSeconds int
type RepCount int
type RoundCount int
type WorkSeconds int
type RestSeconds int
type IntervalSeconds int
type CycleCount int
type LoadValue float64

type WODType string

const (
	WODTypeAMRAP   WODType = "AMRAP"
	WODTypeForTime WODType = "FORTIME"
	WODTypeTabata  WODType = "TABATA"
	WODTypeEMOM    WODType = "EMOM"
)

type WODStatus string

const (
	WODStatusDraft     WODStatus = "DRAFT"
	WODStatusPublished WODStatus = "PUBLISHED"
	WODStatusArchived  WODStatus = "ARCHIVED"
)

type LoadUnit string

const (
	LoadUnitKG         LoadUnit = "kg"
	LoadUnitLB         LoadUnit = "lb"
	LoadUnitBodyweight LoadUnit = "bodyweight"
)

type StageKind string

const (
	StageWarmup   StageKind = "WARMUP"
	StageStrength StageKind = "STRENGTH"
	StageCore     StageKind = "CORE"
	StageMetcon   StageKind = "METCON"
	StageCooldown StageKind = "COOLDOWN"
)

func (id WODID) String() string {
	return string(id)
}

func (id StageID) String() string {
	return string(id)
}

func (id MovementID) String() string {
	return string(id)
}

func validateStageKind(kind StageKind) error {
	switch kind {
	case StageWarmup, StageStrength, StageCore, StageMetcon, StageCooldown:
		return nil
	default:
		return ErrInvalidStageKind
	}
}

func validateName(name WODName) error {
	trimmed := strings.TrimSpace(string(name))
	if utf8.RuneCountInString(trimmed) < 3 {
		return ErrInvalidName
	}
	if utf8.RuneCountInString(trimmed) > 120 {
		return ErrInvalidName
	}
	return nil
}

func validateLoadUnit(unit LoadUnit) error {
	switch unit {
	case LoadUnitKG, LoadUnitLB, LoadUnitBodyweight:
		return nil
	default:
		return ErrInvalidLoadUnit
	}
}
