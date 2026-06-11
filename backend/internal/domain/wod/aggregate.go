package wod

import "time"

type WOD struct {
	id          WODID
	name        WODName
	description WODDescription
	stages      []Stage
	status      WODStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewWOD(
	id WODID,
	name WODName,
	description WODDescription,
	stages []Stage,
	now time.Time,
) (WOD, error) {
	if err := validateName(name); err != nil {
		return WOD{}, err
	}
	if err := validateStages(stages); err != nil {
		return WOD{}, err
	}

	return WOD{
		id:          id,
		name:        name,
		description: description,
		stages:      cloneStages(stages),
		status:      WODStatusDraft,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructWOD(
	id WODID,
	name WODName,
	description WODDescription,
	stages []Stage,
	status WODStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (WOD, error) {
	w, err := NewWOD(id, name, description, stages, createdAt)
	if err != nil {
		return WOD{}, err
	}
	w.status = status
	w.updatedAt = updatedAt
	return w, nil
}

func validateStages(stages []Stage) error {
	if len(stages) == 0 {
		return ErrStageRequired
	}
	for i, stage := range stages {
		if stage.position != i+1 {
			return ErrInvalidStagePosition
		}
	}
	return nil
}

func (w WOD) ID() WODID {
	return w.id
}

func (w WOD) Name() WODName {
	return w.name
}

func (w WOD) Description() WODDescription {
	return w.description
}

func (w WOD) Stages() []Stage {
	return cloneStages(w.stages)
}

func (w WOD) Status() WODStatus {
	return w.status
}

func (w WOD) CreatedAt() time.Time {
	return w.createdAt
}

func (w WOD) UpdatedAt() time.Time {
	return w.updatedAt
}

func (w *WOD) Publish(now time.Time) error {
	if w.status != WODStatusDraft {
		return ErrInvalidStatusTransition
	}
	w.status = WODStatusPublished
	w.updatedAt = now
	return nil
}
