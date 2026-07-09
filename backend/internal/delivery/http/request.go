package http

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type CreateGymRequest struct {
	Name string `json:"name"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role"`
}

type CreateWODRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ScheduledDate string         `json:"scheduledDate,omitempty"`
	Stages        []StageRequest `json:"stages"`
}

type StageRequest struct {
	Kind         string             `json:"kind"`
	Type         string             `json:"type"`
	Instructions string             `json:"instructions"`
	Config       StageConfigRequest `json:"config"`
	Movements    []MovementRequest  `json:"movements"`
}

type StageConfigRequest struct {
	TimeCapSeconds  *int `json:"timeCapSeconds,omitempty"`
	Rounds          *int `json:"rounds,omitempty"`
	WorkSeconds     *int `json:"workSeconds,omitempty"`
	RestSeconds     *int `json:"restSeconds,omitempty"`
	Cycles          *int `json:"cycles,omitempty"`
	IntervalSeconds *int `json:"intervalSeconds,omitempty"`
}

type MovementRequest struct {
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

type SubmitResultRequest struct {
	ScoreValue int    `json:"scoreValue"`
	IsRx       bool   `json:"isRx"`
	Notes      string `json:"notes"`
}

type CreateClassSessionRequest struct {
	WodID     *string `json:"wodId,omitempty"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Capacity  int     `json:"capacity"`
}

type OverbookRequest struct {
	AthleteUserID string `json:"athleteUserId"`
}

type CancelBookingRequest struct {
	AthleteUserID *string `json:"athleteUserId,omitempty"`
}

type SetDefaultSessionRequest struct {
	DayOfWeek int    `json:"dayOfWeek"`
	TimeSlot  string `json:"timeSlot"`
}
