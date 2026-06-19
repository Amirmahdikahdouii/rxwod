package http

import (
	"time"

	appauth "github.com/rxwod/backend/internal/application/auth"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type MeResponse struct {
	User UserResponse        `json:"user"`
	Gyms []WorkspaceResponse `json:"gyms"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type GymResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"ownerId"`
}

type WorkspaceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type MemberResponse struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}

type InvitationResponse struct {
	ID    string `json:"id"`
	GymID string `json:"gymId"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type StageSummaryResponse struct {
	Kind        string `json:"kind"`
	Position    int    `json:"position"`
	Type        string `json:"type"`
	ScoringKind string `json:"scoringKind"`
}

type CreateWODResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Status        string                 `json:"status"`
	StageCount    int                    `json:"stageCount"`
	Stages        []StageSummaryResponse `json:"stages"`
	ScheduledDate string                 `json:"scheduledDate,omitempty"`
}

type WODSummaryResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Status        string                 `json:"status"`
	StageCount    int                    `json:"stageCount"`
	Stages        []StageSummaryResponse `json:"stages"`
	CreatedBy     string                 `json:"createdBy"`
	ScheduledDate string                 `json:"scheduledDate,omitempty"`
	PublishedAt   *time.Time             `json:"publishedAt,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

type WODDetailResponse struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Status        string          `json:"status"`
	Stages        []StageResponse `json:"stages"`
	CreatedBy     string          `json:"createdBy"`
	ScheduledDate string          `json:"scheduledDate,omitempty"`
	PublishedAt   *time.Time      `json:"publishedAt,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type CalendarPlanResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CalendarDayResponse struct {
	Date           string                 `json:"date"`
	PublishedCount int                    `json:"publishedCount"`
	DraftCount     int                    `json:"draftCount"`
	Plans          []CalendarPlanResponse `json:"plans"`
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
	Sets         *int     `json:"sets,omitempty"`
	Reps         *int     `json:"reps,omitempty"`
	LoadValue    *float64 `json:"loadValue,omitempty"`
	LoadUnit     *string  `json:"loadUnit,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toTokenResponse(dto appauth.TokenDTO) TokenResponse {
	return TokenResponse{
		AccessToken:  dto.AccessToken,
		RefreshToken: dto.RefreshToken,
		ExpiresIn:    dto.ExpiresIn,
	}
}

func toAccessTokenResponse(dto appauth.AccessTokenDTO) TokenResponse {
	return TokenResponse{
		AccessToken: dto.AccessToken,
		ExpiresIn:   dto.ExpiresIn,
	}
}

func toUserResponse(dto appauth.UserDTO) UserResponse {
	return UserResponse{
		ID:          dto.ID,
		Email:       dto.Email,
		DisplayName: dto.DisplayName,
	}
}

func toGymResponse(dto appgym.GymDTO) GymResponse {
	return GymResponse{
		ID:      dto.ID,
		Name:    dto.Name,
		OwnerID: dto.OwnerID,
	}
}

func toWorkspaceResponses(dtos []appgym.WorkspaceDTO) []WorkspaceResponse {
	responses := make([]WorkspaceResponse, 0, len(dtos))
	for _, dto := range dtos {
		responses = append(responses, WorkspaceResponse{
			ID:   dto.ID,
			Name: dto.Name,
			Role: string(dto.Role),
		})
	}
	return responses
}

func toMemberResponses(dtos []appgym.MemberDTO) []MemberResponse {
	responses := make([]MemberResponse, 0, len(dtos))
	for _, dto := range dtos {
		responses = append(responses, toMemberResponse(dto))
	}
	return responses
}

func toMemberResponse(dto appgym.MemberDTO) MemberResponse {
	return MemberResponse{
		UserID:      dto.UserID,
		Email:       dto.Email,
		DisplayName: dto.DisplayName,
		Role:        string(dto.Role),
		Status:      string(dto.Status),
	}
}

func toInvitationResponse(dto appgym.InvitationDTO) InvitationResponse {
	return InvitationResponse{
		ID:    dto.ID,
		GymID: dto.GymID,
		Email: dto.Email,
		Role:  string(dto.Role),
	}
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

func formatDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func toCreateResponse(dto appwod.CreateWODResultDTO) CreateWODResponse {
	return CreateWODResponse{
		ID:            dto.ID,
		Name:          dto.Name,
		Status:        string(dto.Status),
		StageCount:    dto.StageCount,
		Stages:        toStageSummaryResponses(dto.Stages),
		ScheduledDate: formatDate(dto.ScheduledDate),
	}
}

func toSummaryResponse(dto appwod.WODSummaryDTO) WODSummaryResponse {
	return WODSummaryResponse{
		ID:            dto.ID,
		Name:          dto.Name,
		Status:        string(dto.Status),
		StageCount:    dto.StageCount,
		Stages:        toStageSummaryResponses(dto.Stages),
		CreatedBy:     dto.CreatedBy,
		ScheduledDate: formatDate(dto.ScheduledDate),
		PublishedAt:   dto.PublishedAt,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
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
				Sets:         movement.Sets,
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
		ID:            dto.ID,
		Name:          dto.Name,
		Description:   dto.Description,
		Status:        string(dto.Status),
		Stages:        stages,
		CreatedBy:     dto.CreatedBy,
		ScheduledDate: formatDate(dto.ScheduledDate),
		PublishedAt:   dto.PublishedAt,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func toCalendarResponses(days []appwod.CalendarDayDTO) []CalendarDayResponse {
	responses := make([]CalendarDayResponse, 0, len(days))
	for _, day := range days {
		plans := make([]CalendarPlanResponse, 0, len(day.Plans))
		for _, plan := range day.Plans {
			plans = append(plans, CalendarPlanResponse{
				ID:     plan.ID,
				Name:   plan.Name,
				Status: string(plan.Status),
			})
		}
		responses = append(responses, CalendarDayResponse{
			Date:           day.Date,
			PublishedCount: day.PublishedCount,
			DraftCount:     day.DraftCount,
			Plans:          plans,
		})
	}
	return responses
}
