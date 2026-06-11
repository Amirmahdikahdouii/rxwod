package wod

import "time"

type WOD[C Config] struct {
	id          WODID
	name        WODName
	description WODDescription
	config      C
	movements   []Movement
	status      WODStatus
	scoring     ScoringConfig
	createdAt   time.Time
	updatedAt   time.Time
}

func NewWOD[C Config](
	id WODID,
	name WODName,
	description WODDescription,
	cfg C,
	movements []Movement,
	now time.Time,
) (WOD[C], error) {
	if err := cfg.Validate(); err != nil {
		return WOD[C]{}, err
	}
	if err := validateName(name); err != nil {
		return WOD[C]{}, err
	}
	if len(movements) == 0 {
		return WOD[C]{}, ErrMovementRequired
	}
	for _, movement := range movements {
		if err := movement.Validate(); err != nil {
			return WOD[C]{}, err
		}
	}

	return WOD[C]{
		id:          id,
		name:        WODName(string(name)),
		description: description,
		config:      cfg,
		movements:   cloneMovements(movements),
		status:      WODStatusDraft,
		scoring:     NewScoringConfig(cfg.ScoringKind()),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructWOD[C Config](
	id WODID,
	name WODName,
	description WODDescription,
	cfg C,
	movements []Movement,
	status WODStatus,
	scoring ScoringConfig,
	createdAt time.Time,
	updatedAt time.Time,
) (WOD[C], error) {
	w, err := NewWOD(id, name, description, cfg, movements, createdAt)
	if err != nil {
		return WOD[C]{}, err
	}
	w.status = status
	w.scoring = scoring
	w.updatedAt = updatedAt
	return w, nil
}

func (w WOD[C]) ID() WODID {
	return w.id
}

func (w WOD[C]) Name() WODName {
	return w.name
}

func (w WOD[C]) Description() WODDescription {
	return w.description
}

func (w WOD[C]) Config() C {
	return w.config
}

func (w WOD[C]) Movements() []Movement {
	return cloneMovements(w.movements)
}

func (w WOD[C]) Status() WODStatus {
	return w.status
}

func (w WOD[C]) Scoring() ScoringConfig {
	return w.scoring
}

func (w WOD[C]) CreatedAt() time.Time {
	return w.createdAt
}

func (w WOD[C]) UpdatedAt() time.Time {
	return w.updatedAt
}

func (w *WOD[C]) Publish(now time.Time) error {
	if w.status != WODStatusDraft {
		return ErrInvalidStatusTransition
	}
	w.status = WODStatusPublished
	w.updatedAt = now
	return nil
}
