package http

import (
	"time"

	appathletedefaultsession "github.com/rxwod/backend/internal/application/athletedefaultsession"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appclassbooking "github.com/rxwod/backend/internal/application/classbooking"
	appclasssession "github.com/rxwod/backend/internal/application/classsession"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	appwodresult "github.com/rxwod/backend/internal/application/wodresult"
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
	ID            string `json:"id"`
	Email         string `json:"email"`
	DisplayName   string `json:"displayName"`
	EmailVerified bool   `json:"emailVerified"`
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
	MembershipID string `json:"membershipId"`
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	Role         string `json:"role"`
	Status       string `json:"status"`
}

type PaginatedMemberResponse struct {
	Data []MemberResponse       `json:"data"`
	Meta PaginationMetaResponse `json:"meta"`
}

type InvitationResponse struct {
	ID    string `json:"id"`
	GymID string `json:"gymId"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token,omitempty"`
}

type InvitationPreviewResponse struct {
	GymID   string `json:"gymId"`
	GymName string `json:"gymName"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Status  string `json:"status"`
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

type PaginationMetaResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type PaginatedWODSummaryResponse struct {
	Data []WODSummaryResponse   `json:"data"`
	Meta PaginationMetaResponse `json:"meta"`
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

type HealthResponse struct {
	Status string `json:"status"`
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
		ID:            dto.ID,
		Email:         dto.Email,
		DisplayName:   dto.DisplayName,
		EmailVerified: dto.EmailVerified,
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
		MembershipID: dto.MembershipID,
		UserID:       dto.UserID,
		Email:        dto.Email,
		DisplayName:  dto.DisplayName,
		Role:         string(dto.Role),
		Status:       string(dto.Status),
	}
}

func toPaginatedMemberResponse(dto appgym.PaginatedMembersDTO) PaginatedMemberResponse {
	return PaginatedMemberResponse{
		Data: toMemberResponses(dto.Data),
		Meta: PaginationMetaResponse{
			Page:       dto.Meta.Page,
			Limit:      dto.Meta.Limit,
			Total:      dto.Meta.Total,
			TotalPages: dto.Meta.TotalPages,
		},
	}
}

func toInvitationResponse(dto appgym.InvitationDTO) InvitationResponse {
	return InvitationResponse{
		ID:    dto.ID,
		GymID: dto.GymID,
		Email: dto.Email,
		Role:  string(dto.Role),
		Token: dto.Token,
	}
}

func toInvitationPreviewResponse(dto appgym.InvitationPreviewDTO) InvitationPreviewResponse {
	return InvitationPreviewResponse{
		GymID:   dto.GymID,
		GymName: dto.GymName,
		Email:   dto.Email,
		Role:    string(dto.Role),
		Status:  string(dto.Status),
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

func toPaginatedSummaryResponse(dto appwod.PaginatedWODSummariesDTO) PaginatedWODSummaryResponse {
	items := make([]WODSummaryResponse, 0, len(dto.Data))
	for _, item := range dto.Data {
		items = append(items, toSummaryResponse(item))
	}
	return PaginatedWODSummaryResponse{
		Data: items,
		Meta: PaginationMetaResponse{
			Page:       dto.Meta.Page,
			Limit:      dto.Meta.Limit,
			Total:      dto.Meta.Total,
			TotalPages: dto.Meta.TotalPages,
		},
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

type ResultResponse struct {
	ID              string    `json:"id"`
	WODID           string    `json:"wodId"`
	GymMembershipID string    `json:"gymMembershipId"`
	ScoreValue      int       `json:"scoreValue"`
	IsRx            bool      `json:"isRx"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type LeaderboardEntryResponse struct {
	Rank            int       `json:"rank"`
	GymMembershipID string    `json:"gymMembershipId"`
	DisplayName     string    `json:"displayName"`
	ScoreValue      int       `json:"scoreValue"`
	IsRx            bool      `json:"isRx"`
	Notes           string    `json:"notes"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type LeaderboardResponse struct {
	WODID   string                     `json:"wodId"`
	Entries []LeaderboardEntryResponse `json:"entries"`
}

func toResultResponse(dto appwodresult.ResultDTO) ResultResponse {
	return ResultResponse{
		ID:              dto.ID,
		WODID:           dto.WODID,
		GymMembershipID: dto.GymMembershipID,
		ScoreValue:      dto.ScoreValue,
		IsRx:            dto.IsRx,
		Notes:           dto.Notes,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

func toLeaderboardResponse(dto appwodresult.LeaderboardDTO) LeaderboardResponse {
	entries := make([]LeaderboardEntryResponse, 0, len(dto.Entries))
	for _, entry := range dto.Entries {
		entries = append(entries, LeaderboardEntryResponse{
			Rank:            entry.Rank,
			GymMembershipID: entry.GymMembershipID,
			DisplayName:     entry.DisplayName,
			ScoreValue:      entry.ScoreValue,
			IsRx:            entry.IsRx,
			Notes:           entry.Notes,
			UpdatedAt:       entry.UpdatedAt,
		})
	}
	return LeaderboardResponse{
		WODID:   dto.WODID,
		Entries: entries,
	}
}

type ClassSessionResponse struct {
	ID              string    `json:"id"`
	GymID           string    `json:"gymId"`
	WodID           *string   `json:"wodId,omitempty"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Capacity        int       `json:"capacity"`
	CoachID         string    `json:"coachId"`
	BookedCount     int       `json:"bookedCount"`
	MyBookingStatus *string   `json:"myBookingStatus,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateClassSessionResponse struct {
	Session                 ClassSessionResponse `json:"session"`
	AutoBookedCount         int                  `json:"autoBookedCount"`
	AutoBookedMembershipIDs []string             `json:"autoBookedMembershipIds"`
}

type ClassBookingResponse struct {
	ID              string    `json:"id"`
	SessionID       string    `json:"sessionId"`
	GymMembershipID string    `json:"gymMembershipId"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type DefaultSessionResponse struct {
	ID              string    `json:"id"`
	GymMembershipID string    `json:"gymMembershipId"`
	DayOfWeek       int       `json:"dayOfWeek"`
	TimeSlot        string    `json:"timeSlot"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func toClassSessionResponse(dto appclasssession.ClassSessionDTO) ClassSessionResponse {
	return ClassSessionResponse{
		ID:              dto.ID,
		GymID:           dto.GymID,
		WodID:           dto.WodID,
		StartTime:       dto.StartTime,
		EndTime:         dto.EndTime,
		Capacity:        dto.Capacity,
		CoachID:         dto.CoachID,
		BookedCount:     dto.BookedCount,
		MyBookingStatus: dto.MyBookingStatus,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

func toCreateClassSessionResponse(dto appclasssession.CreateClassSessionResultDTO) CreateClassSessionResponse {
	return CreateClassSessionResponse{
		Session:                 toClassSessionResponse(dto.Session),
		AutoBookedCount:         dto.AutoBookedCount,
		AutoBookedMembershipIDs: dto.AutoBookedMembershipIDs,
	}
}

func toClassBookingResponse(dto appclassbooking.BookingDTO) ClassBookingResponse {
	return ClassBookingResponse{
		ID:              dto.ID,
		SessionID:       dto.SessionID,
		GymMembershipID: dto.GymMembershipID,
		Status:          dto.Status,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

func toDefaultSessionResponse(dto appathletedefaultsession.DefaultSessionDTO) DefaultSessionResponse {
	return DefaultSessionResponse{
		ID:              dto.ID,
		GymMembershipID: dto.GymMembershipID,
		DayOfWeek:       dto.DayOfWeek,
		TimeSlot:        dto.TimeSlot,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
