package wod

type ScoringKind string

const (
	ScoringRoundsReps     ScoringKind = "ROUNDS_REPS"
	ScoringTimeToComplete ScoringKind = "TIME_TO_COMPLETE"
	ScoringTotalReps      ScoringKind = "TOTAL_REPS"
)

type ScoringConfig struct {
	kind ScoringKind
}

func NewScoringConfig(kind ScoringKind) ScoringConfig {
	return ScoringConfig{kind: kind}
}

func (s ScoringConfig) Kind() ScoringKind {
	return s.kind
}
