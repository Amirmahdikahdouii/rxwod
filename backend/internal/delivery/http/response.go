package http

import (
	"time"

	appwod "github.com/rxwod/backend/internal/application/wod"
)

type CreateWODResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	ScoringKind string `json:"scoringKind"`
}

type WODSummaryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	ScoringKind string    `json:"scoringKind"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type WODDetailResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"type"`
	Status      string               `json:"status"`
	ScoringKind string               `json:"scoringKind"`
	Config      ConfigResponse       `json:"config"`
	Movements   []MovementResponse   `json:"movements"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
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
	ID        string   `json:"id"`
	Position  int      `json:"position"`
	Name      string   `json:"name"`
	Reps      *int     `json:"reps,omitempty"`
	LoadValue *float64 `json:"loadValue,omitempty"`
	LoadUnit  *string  `json:"loadUnit,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toCreateResponse(dto appwod.CreateWODResultDTO) CreateWODResponse {
	return CreateWODResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Type:        string(dto.Type),
		Status:      string(dto.Status),
		ScoringKind: string(dto.ScoringKind),
	}
}

func toSummaryResponse(dto appwod.WODSummaryDTO) WODSummaryResponse {
	return WODSummaryResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Type:        string(dto.Type),
		Status:      string(dto.Status),
		ScoringKind: string(dto.ScoringKind),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func toDetailResponse(dto appwod.WODDetailDTO) WODDetailResponse {
	movements := make([]MovementResponse, 0, len(dto.Movements))
	for _, movement := range dto.Movements {
		movements = append(movements, MovementResponse{
			ID:        movement.ID,
			Position:  movement.Position,
			Name:      movement.Name,
			Reps:      movement.Reps,
			LoadValue: movement.LoadValue,
			LoadUnit:  movement.LoadUnit,
			Notes:     movement.Notes,
		})
	}

	return WODDetailResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Type:        string(dto.Type),
		Status:      string(dto.Status),
		ScoringKind: string(dto.ScoringKind),
		Config: ConfigResponse{
			TimeCapSeconds:  dto.Config.TimeCapSeconds,
			Rounds:          dto.Config.Rounds,
			WorkSeconds:     dto.Config.WorkSeconds,
			RestSeconds:     dto.Config.RestSeconds,
			Cycles:          dto.Config.Cycles,
			IntervalSeconds: dto.Config.IntervalSeconds,
		},
		Movements: movements,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
