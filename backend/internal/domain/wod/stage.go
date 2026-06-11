package wod

type Stage struct {
	id        StageID
	kind      StageKind
	position  int
	config    Config
	scoring   ScoringConfig
	movements []Movement
}

func NewStage(
	id StageID,
	kind StageKind,
	position int,
	cfg Config,
	movements []Movement,
) (Stage, error) {
	if err := validateStageKind(kind); err != nil {
		return Stage{}, err
	}
	if position <= 0 {
		return Stage{}, ErrInvalidStagePosition
	}
	if cfg == nil {
		return Stage{}, ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return Stage{}, err
	}
	if len(movements) == 0 {
		return Stage{}, ErrMovementRequired
	}
	for _, movement := range movements {
		if err := movement.Validate(); err != nil {
			return Stage{}, err
		}
	}

	return Stage{
		id:        id,
		kind:      kind,
		position:  position,
		config:    cfg,
		scoring:   NewScoringConfig(cfg.ScoringKind()),
		movements: cloneMovements(movements),
	}, nil
}

func (s Stage) ID() StageID {
	return s.id
}

func (s Stage) Kind() StageKind {
	return s.kind
}

func (s Stage) Position() int {
	return s.position
}

func (s Stage) Config() Config {
	return s.config
}

func (s Stage) Type() WODType {
	return s.config.Type()
}

func (s Stage) Scoring() ScoringConfig {
	return s.scoring
}

func (s Stage) ScoringKind() ScoringKind {
	return s.scoring.Kind()
}

func (s Stage) Movements() []Movement {
	return cloneMovements(s.movements)
}

func cloneStages(stages []Stage) []Stage {
	cloned := make([]Stage, len(stages))
	copy(cloned, stages)
	return cloned
}
