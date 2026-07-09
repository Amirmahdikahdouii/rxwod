package wod

import (
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type WOD struct {
	id            WODID
	gymID         gym.GymID
	createdBy     user.UserID
	name          WODName
	description   WODDescription
	stages        []Stage
	status        WODStatus
	scheduledDate *time.Time
	publishedAt   *time.Time
	createdAt     time.Time
	updatedAt     time.Time
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
	scheduledDate *time.Time,
	publishedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (WOD, error) {
	w, err := NewWOD(id, gymID, createdBy, name, description, stages, createdAt)
	if err != nil {
		return WOD{}, err
	}
	w.status = status
	w.scheduledDate = cloneDate(scheduledDate)
	w.publishedAt = cloneDate(publishedAt)
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

func cloneDate(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := value.UTC()
	return &cloned
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

func (w WOD) ScheduledDate() *time.Time {
	return cloneDate(w.scheduledDate)
}

func (w WOD) PublishedAt() *time.Time {
	return cloneDate(w.publishedAt)
}

func (w WOD) CreatedAt() time.Time {
	return w.createdAt
}

func (w WOD) UpdatedAt() time.Time {
	return w.updatedAt
}

func (w *WOD) SetScheduledDate(date *time.Time, now time.Time) {
	w.scheduledDate = cloneDate(date)
	w.updatedAt = now
}

func (w *WOD) Publish(now time.Time) error {
	if w.status != WODStatusDraft {
		return ErrInvalidStatusTransition
	}
	if w.scheduledDate == nil {
		return ErrScheduledDateRequired
	}
	w.status = WODStatusPublished
	publishedAt := now.UTC()
	w.publishedAt = &publishedAt
	w.updatedAt = now
	return nil
}

func (w *WOD) Archive(now time.Time) error {
	if w.status != WODStatusPublished {
		return ErrInvalidStatusTransition
	}
	w.status = WODStatusArchived
	w.updatedAt = now
	return nil
}

func (w WOD) CanDelete() bool {
	return w.status == WODStatusDraft || w.status == WODStatusArchived
}
