package wod

import (
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type WOD struct {
	id          WODID
	gymID       gym.GymID
	createdBy   user.UserID
	name        WODName
	description WODDescription
	stages      []Stage
	status      WODStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewWOD(
	id WODID,
	gymID gym.GymID,
	createdBy user.UserID,
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
		gymID:       gymID,
		createdBy:   createdBy,
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
	gymID gym.GymID,
	createdBy user.UserID,
	name WODName,
	description WODDescription,
	stages []Stage,
	status WODStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (WOD, error) {
	w, err := NewWOD(id, gymID, createdBy, name, description, stages, createdAt)
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

func (w WOD) GymID() gym.GymID {
	return w.gymID
}

func (w WOD) CreatedBy() user.UserID {
	return w.createdBy
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
