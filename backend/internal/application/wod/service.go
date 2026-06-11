package wod

import (
	"context"
	"fmt"
	"time"

	"github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo   Repository
	clock  clock.Clock
	idgen  idgen.Generator
}

func NewService(repo Repository, clock clock.Clock, idgen idgen.Generator) *Service {
	return &Service{
		repo:  repo,
		clock: clock,
		idgen: idgen,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (CreateWODResultDTO, error) {
	variant, err := s.buildVariant(cmd)
	if err != nil {
		return CreateWODResultDTO{}, err
	}

	if err := s.repo.Save(ctx, variant); err != nil {
		return CreateWODResultDTO{}, fmt.Errorf("save wod: %w", err)
	}

	return toCreateResultDTO(variant), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (WODDetailDTO, error) {
	variant, err := s.repo.FindByID(ctx, wod.WODID(id))
	if err != nil {
		return WODDetailDTO{}, err
	}
	return toDetailDTO(variant)
}

func (s *Service) List(ctx context.Context) ([]WODSummaryDTO, error) {
	variants, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	summaries := make([]WODSummaryDTO, 0, len(variants))
	for _, variant := range variants {
		summaries = append(summaries, toSummaryDTO(variant))
	}
	return summaries, nil
}

func (s *Service) buildVariant(cmd CreateCommand) (wod.Variant, error) {
	now := s.clock.Now()
	id := wod.WODID(s.idgen.NewID())

	switch cmd.Type {
	case wod.WODTypeAMRAP:
		if cmd.AMRAP == nil {
			return nil, fmt.Errorf("amrap command is required")
		}
		return s.buildAMRAP(id, *cmd.AMRAP, now)
	case wod.WODTypeForTime:
		if cmd.ForTime == nil {
			return nil, fmt.Errorf("fortime command is required")
		}
		return s.buildForTime(id, *cmd.ForTime, now)
	case wod.WODTypeTabata:
		if cmd.Tabata == nil {
			return nil, fmt.Errorf("tabata command is required")
		}
		return s.buildTabata(id, *cmd.Tabata, now)
	case wod.WODTypeEMOM:
		if cmd.EMOM == nil {
			return nil, fmt.Errorf("emom command is required")
		}
		return s.buildEMOM(id, *cmd.EMOM, now)
	default:
		return nil, wod.ErrUnknownWODType
	}
}

func (s *Service) buildAMRAP(id wod.WODID, cmd CreateAMRAPCommand, now time.Time) (wod.Variant, error) {
	cfg, err := wod.NewAMRAPConfig(wod.TimeCapSeconds(cmd.TimeCap))
	if err != nil {
		return nil, err
	}
	movements, err := buildMovements(s.idgen, cmd.Movements)
	if err != nil {
		return nil, err
	}
	aggregate, err := wod.NewWOD(id, wod.WODName(cmd.Name), wod.WODDescription(cmd.Description), cfg, movements, now)
	if err != nil {
		return nil, err
	}
	return wod.NewSavedAMRAP(aggregate), nil
}

func (s *Service) buildForTime(id wod.WODID, cmd CreateForTimeCommand, now time.Time) (wod.Variant, error) {
	var timeCap *wod.TimeCapSeconds
	if cmd.TimeCap != nil {
		value := wod.TimeCapSeconds(*cmd.TimeCap)
		timeCap = &value
	}
	cfg, err := wod.NewForTimeConfig(wod.RoundCount(cmd.Rounds), timeCap)
	if err != nil {
		return nil, err
	}
	movements, err := buildMovements(s.idgen, cmd.Movements)
	if err != nil {
		return nil, err
	}
	aggregate, err := wod.NewWOD(id, wod.WODName(cmd.Name), wod.WODDescription(cmd.Description), cfg, movements, now)
	if err != nil {
		return nil, err
	}
	return wod.NewSavedForTime(aggregate), nil
}

func (s *Service) buildTabata(id wod.WODID, cmd CreateTabataCommand, now time.Time) (wod.Variant, error) {
	cfg, err := wod.NewTabataConfig(
		wod.WorkSeconds(cmd.WorkSeconds),
		wod.RestSeconds(cmd.RestSeconds),
		wod.RoundCount(cmd.Rounds),
		wod.CycleCount(cmd.Cycles),
	)
	if err != nil {
		return nil, err
	}
	movements, err := buildMovements(s.idgen, cmd.Movements)
	if err != nil {
		return nil, err
	}
	aggregate, err := wod.NewWOD(id, wod.WODName(cmd.Name), wod.WODDescription(cmd.Description), cfg, movements, now)
	if err != nil {
		return nil, err
	}
	return wod.NewSavedTabata(aggregate), nil
}

func (s *Service) buildEMOM(id wod.WODID, cmd CreateEMOMCommand, now time.Time) (wod.Variant, error) {
	cfg, err := wod.NewEMOMConfig(wod.IntervalSeconds(cmd.IntervalSeconds), wod.RoundCount(cmd.Rounds))
	if err != nil {
		return nil, err
	}
	movements, err := buildMovements(s.idgen, cmd.Movements)
	if err != nil {
		return nil, err
	}
	aggregate, err := wod.NewWOD(id, wod.WODName(cmd.Name), wod.WODDescription(cmd.Description), cfg, movements, now)
	if err != nil {
		return nil, err
	}
	return wod.NewSavedEMOM(aggregate), nil
}

func buildMovements(generator idgen.Generator, inputs []MovementInput) ([]wod.Movement, error) {
	movements := make([]wod.Movement, 0, len(inputs))
	for _, input := range inputs {
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
		movement, err := wod.NewMovement(
			wod.MovementID(generator.NewID()),
			input.Position,
			input.Name,
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

func toCreateResultDTO(variant wod.Variant) CreateWODResultDTO {
	return CreateWODResultDTO{
		ID:          variant.ID().String(),
		Name:        string(variant.Name()),
		Type:        variant.Type(),
		Status:      variant.Status(),
		ScoringKind: variant.ScoringKind(),
	}
}

func toSummaryDTO(variant wod.Variant) WODSummaryDTO {
	return WODSummaryDTO{
		ID:          variant.ID().String(),
		Name:        string(variant.Name()),
		Type:        variant.Type(),
		Status:      variant.Status(),
		ScoringKind: variant.ScoringKind(),
		CreatedAt:   variant.CreatedAt(),
		UpdatedAt:   variant.UpdatedAt(),
	}
}

func toDetailDTO(variant wod.Variant) (WODDetailDTO, error) {
	config, err := configFromVariant(variant)
	if err != nil {
		return WODDetailDTO{}, err
	}

	movements := make([]MovementDTO, 0, len(variant.Movements()))
	for _, movement := range variant.Movements() {
		movements = append(movements, movementToDTO(movement))
	}

	return WODDetailDTO{
		ID:          variant.ID().String(),
		Name:        string(variant.Name()),
		Description: string(variant.Description()),
		Type:        variant.Type(),
		Status:      variant.Status(),
		ScoringKind: variant.ScoringKind(),
		Config:      config,
		Movements:   movements,
		CreatedAt:   variant.CreatedAt(),
		UpdatedAt:   variant.UpdatedAt(),
	}, nil
}

func configFromVariant(variant wod.Variant) (ConfigDTO, error) {
	switch v := variant.(type) {
	case wod.SavedAMRAP:
		timeCap := int(v.WOD().Config().TimeCap())
		return ConfigDTO{TimeCapSeconds: &timeCap}, nil
	case wod.SavedForTime:
		rounds := int(v.WOD().Config().Rounds())
		dto := ConfigDTO{Rounds: &rounds}
		if capValue := v.WOD().Config().TimeCap(); capValue != nil {
			value := int(*capValue)
			dto.TimeCapSeconds = &value
		}
		return dto, nil
	case wod.SavedTabata:
		work := int(v.WOD().Config().WorkSeconds())
		rest := int(v.WOD().Config().RestSeconds())
		rounds := int(v.WOD().Config().Rounds())
		cycles := int(v.WOD().Config().Cycles())
		return ConfigDTO{
			WorkSeconds: &work,
			RestSeconds: &rest,
			Rounds:      &rounds,
			Cycles:      &cycles,
		}, nil
	case wod.SavedEMOM:
		interval := int(v.WOD().Config().IntervalSeconds())
		rounds := int(v.WOD().Config().Rounds())
		return ConfigDTO{
			IntervalSeconds: &interval,
			Rounds:          &rounds,
		}, nil
	default:
		return ConfigDTO{}, wod.ErrUnknownWODType
	}
}

func movementToDTO(movement wod.Movement) MovementDTO {
	dto := MovementDTO{
		ID:       movement.ID().String(),
		Position: movement.Position(),
		Name:     movement.Name(),
		Notes:    movement.Notes(),
	}
	if reps := movement.Reps(); reps != nil {
		value := int(*reps)
		dto.Reps = &value
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
