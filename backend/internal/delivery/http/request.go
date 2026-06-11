package http

type CreateWODRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Config      CreateWODConfigRequest `json:"config"`
	Movements   []MovementRequest      `json:"movements"`
}

type CreateWODConfigRequest struct {
	TimeCapSeconds  *int `json:"timeCapSeconds,omitempty"`
	Rounds          *int `json:"rounds,omitempty"`
	WorkSeconds     *int `json:"workSeconds,omitempty"`
	RestSeconds     *int `json:"restSeconds,omitempty"`
	Cycles          *int `json:"cycles,omitempty"`
	IntervalSeconds *int `json:"intervalSeconds,omitempty"`
}

type MovementRequest struct {
	Position  int      `json:"position"`
	Name      string   `json:"name"`
	Reps      *int     `json:"reps,omitempty"`
	LoadValue *float64 `json:"loadValue,omitempty"`
	LoadUnit  *string  `json:"loadUnit,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}
