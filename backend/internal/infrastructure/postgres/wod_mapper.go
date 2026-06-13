package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type wodRecord struct {
	ID          string
	Name        string
	Status      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type stageRecord struct {
	ID           string
	WODID        string
	Position     int
	StageKind    string
	WODType      string
	Instructions string
	Config       []byte
	ScoringKind  string
}

type movementRecord struct {
	ID           string
	StageID      string
	Position     int
	Label        string
	Name         string
	Prescription string
	Sets         *int
	Reps         *int
	LoadValue    *float64
	LoadUnit     *string
	Notes        string
}

type amrapConfigPayload struct {
	TimeCapSeconds int `json:"timeCapSeconds"`
}

type forTimeConfigPayload struct {
	Rounds         int  `json:"rounds"`
	TimeCapSeconds *int `json:"timeCapSeconds,omitempty"`
}

type tabataConfigPayload struct {
	WorkSeconds int `json:"workSeconds"`
	RestSeconds int `json:"restSeconds"`
	Rounds      int `json:"rounds"`
	Cycles      int `json:"cycles"`
}

type emomConfigPayload struct {
	IntervalSeconds int `json:"intervalSeconds"`
	Rounds          int `json:"rounds"`
}

func wodToRecords(w domainwod.WOD) (wodRecord, []stageRecord, []movementRecord, error) {
	record := wodRecord{
		ID:          w.ID().String(),
		Name:        string(w.Name()),
		Status:      string(w.Status()),
		Description: string(w.Description()),
		CreatedAt:   w.CreatedAt(),
		UpdatedAt:   w.UpdatedAt(),
	}

	var stages []stageRecord
	var movements []movementRecord
	for _, stage := range w.Stages() {
		config, err := configToJSON(stage.Config())
		if err != nil {
			return wodRecord{}, nil, nil, err
		}
		stages = append(stages, stageRecord{
			ID:           stage.ID().String(),
			WODID:        w.ID().String(),
			Position:     stage.Position(),
			StageKind:    string(stage.Kind()),
			WODType:      string(stage.Type()),
			Instructions: stage.Instructions(),
			Config:       config,
			ScoringKind:  string(stage.ScoringKind()),
		})
		movements = append(movements, movementsToRecords(stage)...)
	}

	return record, stages, movements, nil
}

func movementsToRecords(stage domainwod.Stage) []movementRecord {
	records := make([]movementRecord, 0, len(stage.Movements()))
	for _, movement := range stage.Movements() {
		record := movementRecord{
			ID:           movement.ID().String(),
			StageID:      stage.ID().String(),
			Position:     movement.Position(),
			Label:        movement.Label(),
			Name:         movement.Name(),
			Prescription: movement.Prescription(),
			Notes:        movement.Notes(),
		}
		if sets := movement.Sets(); sets != nil {
			value := int(*sets)
			record.Sets = &value
		}
		if reps := movement.Reps(); reps != nil {
			value := int(*reps)
			record.Reps = &value
		}
		if loadValue := movement.LoadValue(); loadValue != nil {
			value := float64(*loadValue)
			record.LoadValue = &value
		}
		if loadUnit := movement.LoadUnit(); loadUnit != nil {
			value := string(*loadUnit)
			record.LoadUnit = &value
		}
		records = append(records, record)
	}
	return records
}

func configToJSON(cfg domainwod.Config) ([]byte, error) {
	switch c := cfg.(type) {
	case domainwod.OpenConfig:
		return json.Marshal(map[string]any{})
	case domainwod.AMRAPConfig:
		return json.Marshal(amrapConfigPayload{TimeCapSeconds: int(c.TimeCap())})
	case domainwod.ForTimeConfig:
		payload := forTimeConfigPayload{Rounds: int(c.Rounds())}
		if capValue := c.TimeCap(); capValue != nil {
			value := int(*capValue)
			payload.TimeCapSeconds = &value
		}
		return json.Marshal(payload)
	case domainwod.TabataConfig:
		return json.Marshal(tabataConfigPayload{
			WorkSeconds: int(c.WorkSeconds()),
			RestSeconds: int(c.RestSeconds()),
			Rounds:      int(c.Rounds()),
			Cycles:      int(c.Cycles()),
		})
	case domainwod.EMOMConfig:
		return json.Marshal(emomConfigPayload{
			IntervalSeconds: int(c.IntervalSeconds()),
			Rounds:          int(c.Rounds()),
		})
	default:
		return nil, domainwod.ErrUnknownWODType
	}
}

func configFromJSON(wodType string, data []byte) (domainwod.Config, error) {
	switch domainwod.WODType(wodType) {
	case domainwod.WODTypeOpen:
		return domainwod.NewOpenConfig()
	case domainwod.WODTypeAMRAP:
		var payload amrapConfigPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal amrap config: %w", err)
		}
		return domainwod.NewAMRAPConfig(domainwod.TimeCapSeconds(payload.TimeCapSeconds))
	case domainwod.WODTypeForTime:
		var payload forTimeConfigPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal fortime config: %w", err)
		}
		var timeCap *domainwod.TimeCapSeconds
		if payload.TimeCapSeconds != nil {
			value := domainwod.TimeCapSeconds(*payload.TimeCapSeconds)
			timeCap = &value
		}
		return domainwod.NewForTimeConfig(domainwod.RoundCount(payload.Rounds), timeCap)
	case domainwod.WODTypeTabata:
		var payload tabataConfigPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal tabata config: %w", err)
		}
		return domainwod.NewTabataConfig(
			domainwod.WorkSeconds(payload.WorkSeconds),
			domainwod.RestSeconds(payload.RestSeconds),
			domainwod.RoundCount(payload.Rounds),
			domainwod.CycleCount(payload.Cycles),
		)
	case domainwod.WODTypeEMOM:
		var payload emomConfigPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal emom config: %w", err)
		}
		return domainwod.NewEMOMConfig(domainwod.IntervalSeconds(payload.IntervalSeconds), domainwod.RoundCount(payload.Rounds))
	default:
		return nil, domainwod.ErrUnknownWODType
	}
}

func recordsToWOD(record wodRecord, stages []stageRecord, movementsByStage map[string][]movementRecord) (domainwod.WOD, error) {
	domainStages := make([]domainwod.Stage, 0, len(stages))
	for _, stage := range stages {
		cfg, err := configFromJSON(stage.WODType, stage.Config)
		if err != nil {
			return domainwod.WOD{}, err
		}
		movements, err := recordsToMovements(movementsByStage[stage.ID])
		if err != nil {
			return domainwod.WOD{}, err
		}
		domainStage, err := domainwod.NewStage(
			domainwod.StageID(stage.ID),
			domainwod.StageKind(stage.StageKind),
			stage.Position,
			stage.Instructions,
			cfg,
			movements,
		)
		if err != nil {
			return domainwod.WOD{}, err
		}
		domainStages = append(domainStages, domainStage)
	}

	return domainwod.ReconstructWOD(
		domainwod.WODID(record.ID),
		domainwod.WODName(record.Name),
		domainwod.WODDescription(record.Description),
		domainStages,
		domainwod.WODStatus(record.Status),
		record.CreatedAt,
		record.UpdatedAt,
	)
}

func recordsToMovements(records []movementRecord) ([]domainwod.Movement, error) {
	movements := make([]domainwod.Movement, 0, len(records))
	for _, record := range records {
		var sets *domainwod.SetCount
		if record.Sets != nil {
			value := domainwod.SetCount(*record.Sets)
			sets = &value
		}
		var reps *domainwod.RepCount
		if record.Reps != nil {
			value := domainwod.RepCount(*record.Reps)
			reps = &value
		}
		var loadValue *domainwod.LoadValue
		if record.LoadValue != nil {
			value := domainwod.LoadValue(*record.LoadValue)
			loadValue = &value
		}
		var loadUnit *domainwod.LoadUnit
		if record.LoadUnit != nil {
			value := domainwod.LoadUnit(*record.LoadUnit)
			loadUnit = &value
		}
		movement, err := domainwod.NewMovement(
			domainwod.MovementID(record.ID),
			record.Position,
			record.Label,
			record.Name,
			record.Prescription,
			sets,
			reps,
			loadValue,
			loadUnit,
			record.Notes,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}
	return movements, nil
}
