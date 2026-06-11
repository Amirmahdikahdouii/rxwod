package http

type CreateWODRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Stages      []StageRequest `json:"stages"`
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
	Reps         *int     `json:"reps,omitempty"`
	LoadValue    *float64 `json:"loadValue,omitempty"`
	LoadUnit     *string  `json:"loadUnit,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}
