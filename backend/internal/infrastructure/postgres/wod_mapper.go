package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type wodRecord struct {
	ID            string
	Name          string
	WODType       string
	Status        string
	Description   string
	Config        []byte
	ScoringKind   string
	ScoringConfig []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type movementRecord struct {
	ID        string
	WODID     string
	Position  int
	Name      string
	Reps      *int
	LoadValue *float64
	LoadUnit  *string
	Notes     string
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

func variantToRecord(variant domainwod.Variant) (wodRecord, []movementRecord, error) {
	switch v := variant.(type) {
	case domainwod.SavedAMRAP:
		return amrapToRecord(v)
	case domainwod.SavedForTime:
		return forTimeToRecord(v)
	case domainwod.SavedTabata:
		return tabataToRecord(v)
	case domainwod.SavedEMOM:
		return emomToRecord(v)
	default:
		return wodRecord{}, nil, domainwod.ErrUnknownWODType
	}
}

func amrapToRecord(saved domainwod.SavedAMRAP) (wodRecord, []movementRecord, error) {
	w := saved.WOD()
	payload, err := json.Marshal(amrapConfigPayload{TimeCapSeconds: int(w.Config().TimeCap())})
	if err != nil {
		return wodRecord{}, nil, fmt.Errorf("marshal amrap config: %w", err)
	}
	record := baseRecord(w, payload)
	return record, movementsToRecords(w), nil
}

func forTimeToRecord(saved domainwod.SavedForTime) (wodRecord, []movementRecord, error) {
	w := saved.WOD()
	payloadStruct := forTimeConfigPayload{Rounds: int(w.Config().Rounds())}
	if capValue := w.Config().TimeCap(); capValue != nil {
		value := int(*capValue)
		payloadStruct.TimeCapSeconds = &value
	}
	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		return wodRecord{}, nil, fmt.Errorf("marshal fortime config: %w", err)
	}
	record := baseRecord(w, payload)
	return record, movementsToRecords(w), nil
}

func tabataToRecord(saved domainwod.SavedTabata) (wodRecord, []movementRecord, error) {
	w := saved.WOD()
	payload, err := json.Marshal(tabataConfigPayload{
		WorkSeconds: int(w.Config().WorkSeconds()),
		RestSeconds: int(w.Config().RestSeconds()),
		Rounds:      int(w.Config().Rounds()),
		Cycles:      int(w.Config().Cycles()),
	})
	if err != nil {
		return wodRecord{}, nil, fmt.Errorf("marshal tabata config: %w", err)
	}
	record := baseRecord(w, payload)
	return record, movementsToRecords(w), nil
}

func emomToRecord(saved domainwod.SavedEMOM) (wodRecord, []movementRecord, error) {
	w := saved.WOD()
	payload, err := json.Marshal(emomConfigPayload{
		IntervalSeconds: int(w.Config().IntervalSeconds()),
		Rounds:          int(w.Config().Rounds()),
	})
	if err != nil {
		return wodRecord{}, nil, fmt.Errorf("marshal emom config: %w", err)
	}
	record := baseRecord(w, payload)
	return record, movementsToRecords(w), nil
}

func baseRecord[C domainwod.Config](w domainwod.WOD[C], config []byte) wodRecord {
	return wodRecord{
		ID:            w.ID().String(),
		Name:          string(w.Name()),
		WODType:       string(w.Config().Type()),
		Status:        string(w.Status()),
		Description:   string(w.Description()),
		Config:        config,
		ScoringKind:   string(w.Scoring().Kind()),
		ScoringConfig: []byte("{}"),
		CreatedAt:     w.CreatedAt(),
		UpdatedAt:     w.UpdatedAt(),
	}
}

func movementsToRecords[C domainwod.Config](w domainwod.WOD[C]) []movementRecord {
	records := make([]movementRecord, 0, len(w.Movements()))
	for _, movement := range w.Movements() {
		record := movementRecord{
			ID:       movement.ID().String(),
			WODID:    w.ID().String(),
			Position: movement.Position(),
			Name:     movement.Name(),
			Notes:    movement.Notes(),
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

func recordToVariant(record wodRecord, movements []movementRecord) (domainwod.Variant, error) {
	domainMovements, err := recordsToMovements(movements)
	if err != nil {
		return nil, err
	}

	switch domainwod.WODType(record.WODType) {
	case domainwod.WODTypeAMRAP:
		return amrapFromRecord(record, domainMovements)
	case domainwod.WODTypeForTime:
		return forTimeFromRecord(record, domainMovements)
	case domainwod.WODTypeTabata:
		return tabataFromRecord(record, domainMovements)
	case domainwod.WODTypeEMOM:
		return emomFromRecord(record, domainMovements)
	default:
		return nil, domainwod.ErrUnknownWODType
	}
}

func amrapFromRecord(record wodRecord, movements []domainwod.Movement) (domainwod.Variant, error) {
	var payload amrapConfigPayload
	if err := json.Unmarshal(record.Config, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal amrap config: %w", err)
	}
	cfg, err := domainwod.NewAMRAPConfig(domainwod.TimeCapSeconds(payload.TimeCapSeconds))
	if err != nil {
		return nil, err
	}
	w, err := domainwod.ReconstructWOD(
		domainwod.WODID(record.ID),
		domainwod.WODName(record.Name),
		domainwod.WODDescription(record.Description),
		cfg,
		movements,
		domainwod.WODStatus(record.Status),
		domainwod.NewScoringConfig(domainwod.ScoringKind(record.ScoringKind)),
		record.CreatedAt,
		record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return domainwod.NewSavedAMRAP(w), nil
}

func forTimeFromRecord(record wodRecord, movements []domainwod.Movement) (domainwod.Variant, error) {
	var payload forTimeConfigPayload
	if err := json.Unmarshal(record.Config, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal fortime config: %w", err)
	}
	var timeCap *domainwod.TimeCapSeconds
	if payload.TimeCapSeconds != nil {
		value := domainwod.TimeCapSeconds(*payload.TimeCapSeconds)
		timeCap = &value
	}
	cfg, err := domainwod.NewForTimeConfig(domainwod.RoundCount(payload.Rounds), timeCap)
	if err != nil {
		return nil, err
	}
	w, err := domainwod.ReconstructWOD(
		domainwod.WODID(record.ID),
		domainwod.WODName(record.Name),
		domainwod.WODDescription(record.Description),
		cfg,
		movements,
		domainwod.WODStatus(record.Status),
		domainwod.NewScoringConfig(domainwod.ScoringKind(record.ScoringKind)),
		record.CreatedAt,
		record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return domainwod.NewSavedForTime(w), nil
}

func tabataFromRecord(record wodRecord, movements []domainwod.Movement) (domainwod.Variant, error) {
	var payload tabataConfigPayload
	if err := json.Unmarshal(record.Config, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal tabata config: %w", err)
	}
	cfg, err := domainwod.NewTabataConfig(
		domainwod.WorkSeconds(payload.WorkSeconds),
		domainwod.RestSeconds(payload.RestSeconds),
		domainwod.RoundCount(payload.Rounds),
		domainwod.CycleCount(payload.Cycles),
	)
	if err != nil {
		return nil, err
	}
	w, err := domainwod.ReconstructWOD(
		domainwod.WODID(record.ID),
		domainwod.WODName(record.Name),
		domainwod.WODDescription(record.Description),
		cfg,
		movements,
		domainwod.WODStatus(record.Status),
		domainwod.NewScoringConfig(domainwod.ScoringKind(record.ScoringKind)),
		record.CreatedAt,
		record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return domainwod.NewSavedTabata(w), nil
}

func emomFromRecord(record wodRecord, movements []domainwod.Movement) (domainwod.Variant, error) {
	var payload emomConfigPayload
	if err := json.Unmarshal(record.Config, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal emom config: %w", err)
	}
	cfg, err := domainwod.NewEMOMConfig(domainwod.IntervalSeconds(payload.IntervalSeconds), domainwod.RoundCount(payload.Rounds))
	if err != nil {
		return nil, err
	}
	w, err := domainwod.ReconstructWOD(
		domainwod.WODID(record.ID),
		domainwod.WODName(record.Name),
		domainwod.WODDescription(record.Description),
		cfg,
		movements,
		domainwod.WODStatus(record.Status),
		domainwod.NewScoringConfig(domainwod.ScoringKind(record.ScoringKind)),
		record.CreatedAt,
		record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return domainwod.NewSavedEMOM(w), nil
}

func recordsToMovements(records []movementRecord) ([]domainwod.Movement, error) {
	movements := make([]domainwod.Movement, 0, len(records))
	for _, record := range records {
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
			record.Name,
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
