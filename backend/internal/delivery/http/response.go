package http

import (
	"time"

	appwod "github.com/rxwod/backend/internal/application/wod"
)

type StageSummaryResponse struct {
	Kind        string `json:"kind"`
	Position    int    `json:"position"`
	Type        string `json:"type"`
	ScoringKind string `json:"scoringKind"`
}

type CreateWODResponse struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	StageCount int                    `json:"stageCount"`
	Stages     []StageSummaryResponse `json:"stages"`
}

type WODSummaryResponse struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	StageCount int                    `json:"stageCount"`
	Stages     []StageSummaryResponse `json:"stages"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

type WODDetailResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Stages      []StageResponse `json:"stages"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type StageResponse struct {
	ID           string             `json:"id"`
	Kind         string             `json:"kind"`
	Position     int                `json:"position"`
	Instructions string             `json:"instructions"`
	Type         string             `json:"type"`
	ScoringKind  string             `json:"scoringKind"`
	Config       ConfigResponse     `json:"config"`
	Movements    []MovementResponse `json:"movements"`
}

type ConfigResponse struct {
	TimeCapSeconds  *int `json:"timeCapSeconds,omitempty"`
	Rounds          *int `json:"rounds,omitempty"`
	WorkSeconds     *int `json:"workSeconds,omitempty"`
	RestSeconds     *int `json:"restSeconds,omitempty"`
	Cycles          *int `json:"cycles,omitempty"`
	IntervalSeconds *int `json:"intervalSeconds,omitempty"`
}

type MovementResponse struct {
	ID           string   `json:"id"`
	Position     int      `json:"position"`
	Label        string   `json:"label,omitempty"`
	Name         string   `json:"name"`
	Prescription string   `json:"prescription,omitempty"`
	Reps         *int     `json:"reps,omitempty"`
	LoadValue    *float64 `json:"loadValue,omitempty"`
	LoadUnit     *string  `json:"loadUnit,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toStageSummaryResponses(summaries []appwod.StageSummaryDTO) []StageSummaryResponse {
	responses := make([]StageSummaryResponse, 0, len(summaries))
	for _, summary := range summaries {
		responses = append(responses, StageSummaryResponse{
			Kind:        string(summary.Kind),
			Position:    summary.Position,
			Type:        string(summary.Type),
			ScoringKind: string(summary.ScoringKind),
		})
	}
	return responses
}

func toCreateResponse(dto appwod.CreateWODResultDTO) CreateWODResponse {
	return CreateWODResponse{
		ID:         dto.ID,
		Name:       dto.Name,
		Status:     string(dto.Status),
		StageCount: dto.StageCount,
		Stages:     toStageSummaryResponses(dto.Stages),
	}
}

func toSummaryResponse(dto appwod.WODSummaryDTO) WODSummaryResponse {
	return WODSummaryResponse{
		ID:         dto.ID,
		Name:       dto.Name,
		Status:     string(dto.Status),
		StageCount: dto.StageCount,
		Stages:     toStageSummaryResponses(dto.Stages),
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

func toConfigResponse(dto appwod.ConfigDTO) ConfigResponse {
	return ConfigResponse{
		TimeCapSeconds:  dto.TimeCapSeconds,
		Rounds:          dto.Rounds,
		WorkSeconds:     dto.WorkSeconds,
		RestSeconds:     dto.RestSeconds,
		Cycles:          dto.Cycles,
		IntervalSeconds: dto.IntervalSeconds,
	}
}

func toDetailResponse(dto appwod.WODDetailDTO) WODDetailResponse {
	stages := make([]StageResponse, 0, len(dto.Stages))
	for _, stage := range dto.Stages {
		movements := make([]MovementResponse, 0, len(stage.Movements))
		for _, movement := range stage.Movements {
			movements = append(movements, MovementResponse{
				ID:           movement.ID,
				Position:     movement.Position,
				Label:        movement.Label,
				Name:         movement.Name,
				Prescription: movement.Prescription,
				Reps:         movement.Reps,
				LoadValue:    movement.LoadValue,
				LoadUnit:     movement.LoadUnit,
				Notes:        movement.Notes,
			})
		}

		stages = append(stages, StageResponse{
			ID:           stage.ID,
			Kind:         string(stage.Kind),
			Position:     stage.Position,
			Instructions: stage.Instructions,
			Type:         string(stage.Type),
			ScoringKind:  string(stage.ScoringKind),
			Config:       toConfigResponse(stage.Config),
			Movements:    movements,
		})
	}

	return WODDetailResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      string(dto.Status),
		Stages:      stages,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}
