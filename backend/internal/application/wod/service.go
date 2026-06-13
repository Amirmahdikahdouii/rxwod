package wod

import (
	"context"
	"fmt"

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

	stages, err := s.buildStages(cmd.Stages)
	if err != nil {
		return WODDetailDTO{}, err
	}

	updated, err := wod.ReconstructWOD(
		existing.ID(),
		existing.GymID(),
		existing.CreatedBy(),
		wod.WODName(cmd.Name),
		wod.WODDescription(cmd.Description),
		stages,
		existing.Status(),
		existing.CreatedAt(),
		s.clock.Now(),
	)
	if err != nil {
		return WODDetailDTO{}, err
	}

	if err := s.repo.Save(ctx, updated); err != nil {
		return WODDetailDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toDetailDTO(updated), nil
}

func (s *Service) List(ctx context.Context) ([]WODSummaryDTO, error) {
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
		summaries = append(summaries, toSummaryDTO(aggregate))
	}
	return summaries, nil
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
		return wod.Stage{}, err
	}
	movements, err := buildMovements(s.idgen, input.Movements)
	if err != nil {
		return wod.Stage{}, err
	}
	return wod.NewStage(wod.StageID(s.idgen.NewID()), input.Kind, position, input.Instructions, cfg, movements)
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

func buildMovements(generator idgen.Generator, inputs []MovementInput) ([]wod.Movement, error) {
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
		if input.LoadValue != nil {
			value := wod.LoadValue(*input.LoadValue)
			loadValue = &value
		}
		var loadUnit *wod.LoadUnit
		if input.LoadUnit != nil {
			value := wod.LoadUnit(*input.LoadUnit)
			loadUnit = &value
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
			return nil, err
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func toCreateResultDTO(aggregate wod.WOD) CreateWODResultDTO {
	stages := aggregate.Stages()
	return CreateWODResultDTO{
		ID:         aggregate.ID().String(),
		Name:       string(aggregate.Name()),
		Status:     aggregate.Status(),
		StageCount: len(stages),
		Stages:     stagesToSummaryDTO(stages),
	}
}

func toSummaryDTO(aggregate wod.WOD) WODSummaryDTO {
	stages := aggregate.Stages()
	return WODSummaryDTO{
		ID:         aggregate.ID().String(),
		Name:       string(aggregate.Name()),
		Status:     aggregate.Status(),
		StageCount: len(stages),
		Stages:     stagesToSummaryDTO(stages),
		CreatedAt:  aggregate.CreatedAt(),
		UpdatedAt:  aggregate.UpdatedAt(),
	}
}

func toDetailDTO(aggregate wod.WOD) WODDetailDTO {
	stages := aggregate.Stages()
	stageDTOs := make([]StageDTO, 0, len(stages))
	for _, stage := range stages {
		stageDTOs = append(stageDTOs, stageToDTO(stage))
	}

	return WODDetailDTO{
		ID:          aggregate.ID().String(),
		Name:        string(aggregate.Name()),
		Description: string(aggregate.Description()),
		Status:      aggregate.Status(),
		Stages:      stageDTOs,
		CreatedAt:   aggregate.CreatedAt(),
		UpdatedAt:   aggregate.UpdatedAt(),
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
