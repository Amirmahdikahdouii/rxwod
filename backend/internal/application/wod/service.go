package wod

import (
	"context"
	"fmt"
	"strings"
	"time"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo  Repository
	clock clock.Clock
	idgen idgen.Generator
}

func NewService(repo Repository, clock clock.Clock, idgen idgen.Generator) *Service {
	return &Service{
		repo:  repo,
		clock: clock,
		idgen: idgen,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateWODCommand) (CreateWODResultDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODCreate)
	if err != nil {
		return CreateWODResultDTO{}, err
	}

	aggregate, err := s.buildWOD(cmd, principal)
	if err != nil {
		return CreateWODResultDTO{}, err
	}
	if cmd.ScheduledDate != nil {
		aggregate.SetScheduledDate(cmd.ScheduledDate, s.clock.Now())
	}

	if err := s.repo.Save(ctx, aggregate); err != nil {
		return CreateWODResultDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toCreateResultDTO(aggregate), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (WODDetailDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODRead)
	if err != nil {
		return WODDetailDTO{}, err
	}

	aggregate, err := s.repo.FindByID(ctx, principal.GymID, wod.WODID(id))
	if err != nil {
		return WODDetailDTO{}, err
	}
	if !canViewWOD(principal, aggregate) {
		return WODDetailDTO{}, appauthz.ErrForbidden
	}
	return toDetailDTO(aggregate), nil
}

func (s *Service) Update(ctx context.Context, id string, cmd CreateWODCommand) (WODDetailDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODUpdate)
	if err != nil {
		return WODDetailDTO{}, err
	}

	existing, err := s.repo.FindByID(ctx, principal.GymID, wod.WODID(id))
	if err != nil {
		return WODDetailDTO{}, err
	}
	if err := ensureCanModifyWOD(principal, existing); err != nil {
		return WODDetailDTO{}, err
	}

	stages, err := s.buildStages(cmd.Stages)
	if err != nil {
		return WODDetailDTO{}, err
	}

	now := s.clock.Now()
	updated, err := wod.ReconstructWOD(
		existing.ID(),
		existing.GymID(),
		existing.CreatedBy(),
		wod.WODName(cmd.Name),
		wod.WODDescription(cmd.Description),
		stages,
		existing.Status(),
		cmd.ScheduledDate,
		existing.PublishedAt(),
		existing.CreatedAt(),
		now,
	)
	if err != nil {
		return WODDetailDTO{}, err
	}
	if cmd.ScheduledDate == nil {
		updated.SetScheduledDate(nil, now)
	}

	if err := s.repo.Save(ctx, updated); err != nil {
		return WODDetailDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toDetailDTO(updated), nil
}

func (s *Service) Publish(ctx context.Context, id string) (WODDetailDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODPublish)
	if err != nil {
		return WODDetailDTO{}, err
	}

	existing, err := s.repo.FindByID(ctx, principal.GymID, wod.WODID(id))
	if err != nil {
		return WODDetailDTO{}, err
	}
	if err := ensureCanModifyWOD(principal, existing); err != nil {
		return WODDetailDTO{}, err
	}

	aggregate := existing
	if err := aggregate.Publish(s.clock.Now()); err != nil {
		return WODDetailDTO{}, err
	}

	if err := s.repo.Save(ctx, aggregate); err != nil {
		return WODDetailDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toDetailDTO(aggregate), nil
}

func (s *Service) Archive(ctx context.Context, id string) (WODDetailDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODDelete)
	if err != nil {
		return WODDetailDTO{}, err
	}

	existing, err := s.repo.FindByID(ctx, principal.GymID, wod.WODID(id))
	if err != nil {
		return WODDetailDTO{}, err
	}

	aggregate := existing
	if err := aggregate.Archive(s.clock.Now()); err != nil {
		return WODDetailDTO{}, err
	}

	if err := s.repo.Save(ctx, aggregate); err != nil {
		return WODDetailDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toDetailDTO(aggregate), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODDelete)
	if err != nil {
		return err
	}

	existing, err := s.repo.FindByID(ctx, principal.GymID, wod.WODID(id))
	if err != nil {
		return err
	}
	if !existing.CanDelete() {
		return wod.ErrCannotDeletePublished
	}

	if err := s.repo.Delete(ctx, principal.GymID, wod.WODID(id)); err != nil {
		return fmt.Errorf("delete wod: %w", err)
	}
	return nil
}

func (s *Service) List(ctx context.Context, statusFilter *wod.WODStatus) ([]WODSummaryDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODRead)
	if err != nil {
		return nil, err
	}

	aggregates, err := s.repo.List(ctx, principal.GymID)
	if err != nil {
		return nil, err
	}

	summaries := make([]WODSummaryDTO, 0, len(aggregates))
	for _, aggregate := range aggregates {
		if !canViewWOD(principal, aggregate) {
			continue
		}
		if statusFilter != nil && aggregate.Status() != *statusFilter {
			continue
		}
		summaries = append(summaries, toSummaryDTO(aggregate))
	}
	return summaries, nil
}

func (s *Service) Calendar(ctx context.Context, from, to time.Time) ([]CalendarDayDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODRead)
	if err != nil {
		return nil, err
	}

	includeDrafts := principal.Role != domainauthz.RoleAthlete
	entries, err := s.repo.ListCalendar(ctx, principal.GymID, from, to, includeDrafts)
	if err != nil {
		return nil, err
	}

	byDate := make(map[string]*CalendarDayDTO)
	for _, entry := range entries {
		dateKey := entry.ScheduledDate.Format("2006-01-02")
		day, ok := byDate[dateKey]
		if !ok {
			day = &CalendarDayDTO{Date: dateKey, Plans: []CalendarPlanDTO{}}
			byDate[dateKey] = day
		}
		if entry.Status == wod.WODStatusPublished {
			day.PublishedCount++
		} else if entry.Status == wod.WODStatusDraft && includeDrafts {
			day.DraftCount++
		}
		if entry.Status == wod.WODStatusPublished || includeDrafts {
			day.Plans = append(day.Plans, CalendarPlanDTO{
				ID:     entry.ID,
				Name:   entry.Name,
				Status: entry.Status,
			})
		}
	}

	days := make([]CalendarDayDTO, 0, len(byDate))
	for _, day := range byDate {
		days = append(days, *day)
	}
	return days, nil
}

func canViewWOD(principal appauthz.Principal, aggregate wod.WOD) bool {
	if principal.Role == domainauthz.RoleAthlete && aggregate.Status() != wod.WODStatusPublished {
		return false
	}
	return true
}

func ensureCanModifyWOD(principal appauthz.Principal, aggregate wod.WOD) error {
	if principal.Role == domainauthz.RoleOwner {
		return nil
	}
	if principal.Role == domainauthz.RoleCoach &&
		aggregate.CreatedBy() == principal.UserID &&
		(aggregate.Status() == wod.WODStatusDraft || aggregate.Status() == wod.WODStatusPublished) {
		return nil
	}
	return appauthz.ErrForbidden
}

func (s *Service) buildWOD(cmd CreateWODCommand, principal appauthz.Principal) (wod.WOD, error) {
	now := s.clock.Now()
	id := wod.WODID(s.idgen.NewID())

	stages, err := s.buildStages(cmd.Stages)
	if err != nil {
		return wod.WOD{}, err
	}

	return wod.NewWOD(id, principal.GymID, principal.UserID, wod.WODName(cmd.Name), wod.WODDescription(cmd.Description), stages, now)
}

func (s *Service) buildStages(inputs []StageInput) ([]wod.Stage, error) {
	stages := make([]wod.Stage, 0, len(inputs))
	for i, stageInput := range inputs {
		stage, err := s.buildStage(stageInput, i+1)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	return stages, nil
}

func (s *Service) buildStage(input StageInput, position int) (wod.Stage, error) {
	cfg, err := buildConfig(input.Config)
	if err != nil {
		return wod.Stage{}, formatStageError(input.Kind, err)
	}
	movements, err := buildMovements(s.idgen, input.Kind, input.Movements)
	if err != nil {
		return wod.Stage{}, err
	}
	stage, err := wod.NewStage(wod.StageID(s.idgen.NewID()), input.Kind, position, input.Instructions, cfg, movements)
	if err != nil {
		return wod.Stage{}, formatStageError(input.Kind, err)
	}
	return stage, nil
}

func buildConfig(input StageConfigInput) (wod.Config, error) {
	switch input.Type {
	case wod.WODTypeOpen:
		return wod.NewOpenConfig()
	case wod.WODTypeAMRAP:
		if input.TimeCapSeconds == nil {
			return nil, ErrMissingConfigField
		}
		return wod.NewAMRAPConfig(wod.TimeCapSeconds(*input.TimeCapSeconds))
	case wod.WODTypeForTime:
		if input.Rounds == nil {
			return nil, ErrMissingConfigField
		}
		var timeCap *wod.TimeCapSeconds
		if input.TimeCapSeconds != nil {
			value := wod.TimeCapSeconds(*input.TimeCapSeconds)
			timeCap = &value
		}
		return wod.NewForTimeConfig(wod.RoundCount(*input.Rounds), timeCap)
	case wod.WODTypeTabata:
		if input.WorkSeconds == nil || input.RestSeconds == nil || input.Rounds == nil || input.Cycles == nil {
			return nil, ErrMissingConfigField
		}
		return wod.NewTabataConfig(
			wod.WorkSeconds(*input.WorkSeconds),
			wod.RestSeconds(*input.RestSeconds),
			wod.RoundCount(*input.Rounds),
			wod.CycleCount(*input.Cycles),
		)
	case wod.WODTypeEMOM:
		if input.IntervalSeconds == nil || input.Rounds == nil {
			return nil, ErrMissingConfigField
		}
		return wod.NewEMOMConfig(wod.IntervalSeconds(*input.IntervalSeconds), wod.RoundCount(*input.Rounds))
	default:
		return nil, wod.ErrUnknownWODType
	}
}

func buildMovements(generator idgen.Generator, kind wod.StageKind, inputs []MovementInput) ([]wod.Movement, error) {
	movements := make([]wod.Movement, 0, len(inputs))
	for i, input := range inputs {
		var sets *wod.SetCount
		if input.Sets != nil {
			value := wod.SetCount(*input.Sets)
			sets = &value
		}
		var reps *wod.RepCount
		if input.Reps != nil {
			value := wod.RepCount(*input.Reps)
			reps = &value
		}
		var loadValue *wod.LoadValue
		var loadUnit *wod.LoadUnit
		if input.LoadValue != nil && *input.LoadValue > 0 {
			value := wod.LoadValue(*input.LoadValue)
			loadValue = &value
			if input.LoadUnit != nil {
				unit := wod.LoadUnit(*input.LoadUnit)
				loadUnit = &unit
			}
		}

		position := input.Position
		if position == 0 {
			position = i + 1
		}

		movement, err := wod.NewMovement(
			wod.MovementID(generator.NewID()),
			position,
			input.Label,
			input.Name,
			input.Prescription,
			sets,
			reps,
			loadValue,
			loadUnit,
			input.Notes,
		)
		if err != nil {
			return nil, formatMovementError(kind, input.Label, position, input.Name, err)
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func formatStageError(kind wod.StageKind, err error) error {
	return fmt.Errorf("%s stage: %w", kind, err)
}

func formatMovementError(kind wod.StageKind, label string, position int, name string, err error) error {
	return fmt.Errorf("%s stage, %s: %w", kind, movementItemRef(label, position, name), err)
}

func movementItemRef(label string, position int, name string) string {
	if trimmed := strings.TrimSpace(label); trimmed != "" {
		return "item " + trimmed
	}
	if trimmed := strings.TrimSpace(name); trimmed != "" {
		return "\"" + trimmed + "\""
	}
	return fmt.Sprintf("item %d", position)
}

func toCreateResultDTO(aggregate wod.WOD) CreateWODResultDTO {
	stages := aggregate.Stages()
	return CreateWODResultDTO{
		ID:            aggregate.ID().String(),
		Name:          string(aggregate.Name()),
		Status:        aggregate.Status(),
		StageCount:    len(stages),
		Stages:        stagesToSummaryDTO(stages),
		ScheduledDate: aggregate.ScheduledDate(),
	}
}

func toSummaryDTO(aggregate wod.WOD) WODSummaryDTO {
	stages := aggregate.Stages()
	return WODSummaryDTO{
		ID:            aggregate.ID().String(),
		Name:          string(aggregate.Name()),
		Status:        aggregate.Status(),
		StageCount:    len(stages),
		Stages:        stagesToSummaryDTO(stages),
		CreatedBy:     aggregate.CreatedBy().String(),
		ScheduledDate: aggregate.ScheduledDate(),
		PublishedAt:   aggregate.PublishedAt(),
		CreatedAt:     aggregate.CreatedAt(),
		UpdatedAt:     aggregate.UpdatedAt(),
	}
}

func toDetailDTO(aggregate wod.WOD) WODDetailDTO {
	stages := aggregate.Stages()
	stageDTOs := make([]StageDTO, 0, len(stages))
	for _, stage := range stages {
		stageDTOs = append(stageDTOs, stageToDTO(stage))
	}

	return WODDetailDTO{
		ID:            aggregate.ID().String(),
		Name:          string(aggregate.Name()),
		Description:   string(aggregate.Description()),
		Status:        aggregate.Status(),
		Stages:        stageDTOs,
		CreatedBy:     aggregate.CreatedBy().String(),
		ScheduledDate: aggregate.ScheduledDate(),
		PublishedAt:   aggregate.PublishedAt(),
		CreatedAt:     aggregate.CreatedAt(),
		UpdatedAt:     aggregate.UpdatedAt(),
	}
}

func stagesToSummaryDTO(stages []wod.Stage) []StageSummaryDTO {
	summaries := make([]StageSummaryDTO, 0, len(stages))
	for _, stage := range stages {
		summaries = append(summaries, StageSummaryDTO{
			Kind:        stage.Kind(),
			Position:    stage.Position(),
			Type:        stage.Type(),
			ScoringKind: stage.ScoringKind(),
		})
	}
	return summaries
}

func stageToDTO(stage wod.Stage) StageDTO {
	movements := make([]MovementDTO, 0, len(stage.Movements()))
	for _, movement := range stage.Movements() {
		movements = append(movements, movementToDTO(movement))
	}

	return StageDTO{
		ID:           stage.ID().String(),
		Kind:         stage.Kind(),
		Position:     stage.Position(),
		Instructions: stage.Instructions(),
		Type:         stage.Type(),
		ScoringKind:  stage.ScoringKind(),
		Config:       configToDTO(stage.Config()),
		Movements:    movements,
	}
}

func configToDTO(cfg wod.Config) ConfigDTO {
	switch c := cfg.(type) {
	case wod.OpenConfig:
		return ConfigDTO{}
	case wod.AMRAPConfig:
		timeCap := int(c.TimeCap())
		return ConfigDTO{TimeCapSeconds: &timeCap}
	case wod.ForTimeConfig:
		rounds := int(c.Rounds())
		dto := ConfigDTO{Rounds: &rounds}
		if capValue := c.TimeCap(); capValue != nil {
			value := int(*capValue)
			dto.TimeCapSeconds = &value
		}
		return dto
	case wod.TabataConfig:
		work := int(c.WorkSeconds())
		rest := int(c.RestSeconds())
		rounds := int(c.Rounds())
		cycles := int(c.Cycles())
		return ConfigDTO{
			WorkSeconds: &work,
			RestSeconds: &rest,
			Rounds:      &rounds,
			Cycles:      &cycles,
		}
	case wod.EMOMConfig:
		interval := int(c.IntervalSeconds())
		rounds := int(c.Rounds())
		return ConfigDTO{
			IntervalSeconds: &interval,
			Rounds:          &rounds,
		}
	default:
		return ConfigDTO{}
	}
}

func movementToDTO(movement wod.Movement) MovementDTO {
	dto := MovementDTO{
		ID:           movement.ID().String(),
		Position:     movement.Position(),
		Label:        movement.Label(),
		Name:         movement.Name(),
		Prescription: movement.Prescription(),
		Notes:        movement.Notes(),
	}
	if reps := movement.Reps(); reps != nil {
		value := int(*reps)
		dto.Reps = &value
	}
	if sets := movement.Sets(); sets != nil {
		value := int(*sets)
		dto.Sets = &value
	}
	if loadValue := movement.LoadValue(); loadValue != nil {
		value := float64(*loadValue)
		dto.LoadValue = &value
	}
	if loadUnit := movement.LoadUnit(); loadUnit != nil {
		value := string(*loadUnit)
		dto.LoadUnit = &value
	}
	return dto
}
