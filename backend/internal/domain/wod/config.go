package wod

type Config interface {
	wodConfig()
	Type() WODType
	Validate() error
	ScoringKind() ScoringKind
}

type AMRAPConfig struct {
	timeCap TimeCapSeconds
}

func NewAMRAPConfig(timeCap TimeCapSeconds) (AMRAPConfig, error) {
	cfg := AMRAPConfig{timeCap: timeCap}
	return cfg, cfg.Validate()
}

func (AMRAPConfig) wodConfig()                 {}
func (c AMRAPConfig) Type() WODType            { return WODTypeAMRAP }
func (c AMRAPConfig) ScoringKind() ScoringKind { return ScoringRoundsReps }
func (c AMRAPConfig) TimeCap() TimeCapSeconds  { return c.timeCap }

func (c AMRAPConfig) Validate() error {
	if c.timeCap <= 0 {
		return ErrInvalidTimeCap
	}
	return nil
}

type ForTimeConfig struct {
	rounds  RoundCount
	timeCap *TimeCapSeconds
}

func NewForTimeConfig(rounds RoundCount, timeCap *TimeCapSeconds) (ForTimeConfig, error) {
	cfg := ForTimeConfig{rounds: rounds, timeCap: timeCap}
	return cfg, cfg.Validate()
}

func (ForTimeConfig) wodConfig()                 {}
func (c ForTimeConfig) Type() WODType            { return WODTypeForTime }
func (c ForTimeConfig) ScoringKind() ScoringKind { return ScoringTimeToComplete }
func (c ForTimeConfig) Rounds() RoundCount       { return c.rounds }
func (c ForTimeConfig) TimeCap() *TimeCapSeconds {
	if c.timeCap == nil {
		return nil
	}
	value := *c.timeCap
	return &value
}

func (c ForTimeConfig) Validate() error {
	if c.rounds <= 0 {
		return ErrInvalidRounds
	}
	if c.timeCap != nil && *c.timeCap <= 0 {
		return ErrInvalidTimeCap
	}
	return nil
}

type TabataConfig struct {
	workSeconds WorkSeconds
	restSeconds RestSeconds
	rounds      RoundCount
	cycles      CycleCount
}

func NewTabataConfig(workSeconds WorkSeconds, restSeconds RestSeconds, rounds RoundCount, cycles CycleCount) (TabataConfig, error) {
	cfg := TabataConfig{
		workSeconds: workSeconds,
		restSeconds: restSeconds,
		rounds:      rounds,
		cycles:      cycles,
	}
	return cfg, cfg.Validate()
}

func (TabataConfig) wodConfig()                 {}
func (c TabataConfig) Type() WODType            { return WODTypeTabata }
func (c TabataConfig) ScoringKind() ScoringKind { return ScoringTotalReps }
func (c TabataConfig) WorkSeconds() WorkSeconds { return c.workSeconds }
func (c TabataConfig) RestSeconds() RestSeconds { return c.restSeconds }
func (c TabataConfig) Rounds() RoundCount       { return c.rounds }
func (c TabataConfig) Cycles() CycleCount       { return c.cycles }

func (c TabataConfig) Validate() error {
	if c.workSeconds <= 0 {
		return ErrInvalidWorkSeconds
	}
	if c.restSeconds < 0 {
		return ErrInvalidRestSeconds
	}
	if c.rounds <= 0 {
		return ErrInvalidRounds
	}
	if c.cycles <= 0 {
		return ErrInvalidCycles
	}
	return nil
}

type EMOMConfig struct {
	intervalSeconds IntervalSeconds
	rounds          RoundCount
}

func NewEMOMConfig(intervalSeconds IntervalSeconds, rounds RoundCount) (EMOMConfig, error) {
	cfg := EMOMConfig{intervalSeconds: intervalSeconds, rounds: rounds}
	return cfg, cfg.Validate()
}

func (EMOMConfig) wodConfig()                         {}
func (c EMOMConfig) Type() WODType                    { return WODTypeEMOM }
func (c EMOMConfig) ScoringKind() ScoringKind         { return ScoringRoundsReps }
func (c EMOMConfig) IntervalSeconds() IntervalSeconds { return c.intervalSeconds }
func (c EMOMConfig) Rounds() RoundCount               { return c.rounds }

func (c EMOMConfig) Validate() error {
	if c.intervalSeconds <= 0 {
		return ErrInvalidInterval
	}
	if c.rounds <= 0 {
		return ErrInvalidRounds
	}
	return nil
}

type OpenConfig struct{}

func NewOpenConfig() (OpenConfig, error) {
	cfg := OpenConfig{}
	return cfg, cfg.Validate()
}

func (OpenConfig) wodConfig()                 {}
func (OpenConfig) Type() WODType            { return WODTypeOpen }
func (OpenConfig) ScoringKind() ScoringKind { return ScoringNone }

func (OpenConfig) Validate() error {
	return nil
}
