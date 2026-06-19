package wod

import "errors"

var (
	ErrInvalidName             = errors.New("wod name must be between 3 and 120 characters")
	ErrInvalidTimeCap          = errors.New("time cap must be greater than zero")
	ErrInvalidRounds           = errors.New("rounds must be greater than zero")
	ErrInvalidWorkSeconds      = errors.New("work seconds must be greater than zero")
	ErrInvalidRestSeconds      = errors.New("rest seconds must be zero or greater")
	ErrInvalidInterval         = errors.New("interval seconds must be greater than zero")
	ErrInvalidCycles           = errors.New("cycles must be greater than zero")
	ErrMovementRequired        = errors.New("at least one movement is required")
	ErrMovementTextRequired    = errors.New("movement name or prescription is required")
	ErrLoadValueRequiresUnit   = errors.New("load value requires a load unit")
	ErrInvalidMovementLabel    = errors.New("movement label must be at most 4 characters")
	ErrInvalidLoadUnit         = errors.New("load unit is invalid")
	ErrInvalidSets             = errors.New("sets must be greater than zero when provided")
	ErrInvalidReps             = errors.New("reps must be greater than zero when provided")
	ErrInvalidPosition         = errors.New("movement position must be greater than zero")
	ErrUnknownWODType          = errors.New("unknown wod type")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrScheduledDateRequired = errors.New("scheduled date is required to publish")
	ErrStageRequired         = errors.New("at least one stage is required")
	ErrInvalidStageKind        = errors.New("stage kind is invalid")
	ErrInvalidStagePosition    = errors.New("stage position must be greater than zero")
	ErrNilConfig               = errors.New("stage config is required")
)
