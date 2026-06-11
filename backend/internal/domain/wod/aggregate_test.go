package wod

import (
	"testing"
	"time"
)

func sampleMovement(t *testing.T) Movement {
	t.Helper()
	reps := RepCount(21)
	movement, err := NewMovement(
		MovementID("mov-1"),
		1,
		"Burpee",
		&reps,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("unexpected movement error: %v", err)
	}
	return movement
}

func sampleStage(t *testing.T, kind StageKind, position int, cfg Config) Stage {
	t.Helper()
	stage, err := NewStage(
		StageID("stage-1"),
		kind,
		position,
		cfg,
		[]Movement{sampleMovement(t)},
	)
	if err != nil {
		t.Fatalf("unexpected stage error: %v", err)
	}
	return stage
}

func amrapConfig(t *testing.T) AMRAPConfig {
	t.Helper()
	cfg, err := NewAMRAPConfig(TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}
	return cfg
}

func TestNewMultiStageWOD(t *testing.T) {
	forTime, err := NewForTimeConfig(RoundCount(2), nil)
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}
	tabata, err := NewTabataConfig(WorkSeconds(20), RestSeconds(10), RoundCount(8), CycleCount(1))
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}

	stages := []Stage{
		sampleStage(t, StageWarmup, 1, forTime),
		sampleStage(t, StageMetcon, 2, amrapConfig(t)),
		sampleStage(t, StageCooldown, 3, tabata),
	}

	wod, err := NewWOD(
		WODID("wod-1"),
		WODName("Monday Session"),
		WODDescription("test"),
		stages,
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("unexpected wod error: %v", err)
	}

	got := wod.Stages()
	if len(got) != 3 {
		t.Fatalf("expected 3 stages, got %d", len(got))
	}
	if got[0].ScoringKind() != ScoringTimeToComplete {
		t.Fatalf("expected warmup fortime scoring TIME_TO_COMPLETE, got %s", got[0].ScoringKind())
	}
	if got[1].ScoringKind() != ScoringRoundsReps {
		t.Fatalf("expected metcon amrap scoring ROUNDS_REPS, got %s", got[1].ScoringKind())
	}
	if got[2].ScoringKind() != ScoringTotalReps {
		t.Fatalf("expected cooldown tabata scoring TOTAL_REPS, got %s", got[2].ScoringKind())
	}
}

func TestNewWODRequiresStage(t *testing.T) {
	_, err := NewWOD(
		WODID("wod-1"),
		WODName("Empty Program"),
		WODDescription("test"),
		nil,
		time.Now().UTC(),
	)
	if err != ErrStageRequired {
		t.Fatalf("expected ErrStageRequired, got %v", err)
	}
}

func TestNewWODRejectsNonContiguousPositions(t *testing.T) {
	stages := []Stage{
		sampleStage(t, StageWarmup, 1, amrapConfig(t)),
		sampleStage(t, StageMetcon, 3, amrapConfig(t)),
	}

	_, err := NewWOD(
		WODID("wod-1"),
		WODName("Bad Positions"),
		WODDescription("test"),
		stages,
		time.Now().UTC(),
	)
	if err != ErrInvalidStagePosition {
		t.Fatalf("expected ErrInvalidStagePosition, got %v", err)
	}
}

func TestNewStageRejectsInvalidKind(t *testing.T) {
	_, err := NewStage(
		StageID("stage-1"),
		StageKind("INVALID"),
		1,
		amrapConfig(t),
		[]Movement{sampleMovement(t)},
	)
	if err != ErrInvalidStageKind {
		t.Fatalf("expected ErrInvalidStageKind, got %v", err)
	}
}

func TestNewStageRequiresMovement(t *testing.T) {
	_, err := NewStage(
		StageID("stage-1"),
		StageWarmup,
		1,
		amrapConfig(t),
		nil,
	)
	if err != ErrMovementRequired {
		t.Fatalf("expected ErrMovementRequired, got %v", err)
	}
}

func TestInvalidTimeCap(t *testing.T) {
	_, err := NewAMRAPConfig(TimeCapSeconds(0))
	if err != ErrInvalidTimeCap {
		t.Fatalf("expected ErrInvalidTimeCap, got %v", err)
	}
}

func TestTabataConfigValidation(t *testing.T) {
	_, err := NewTabataConfig(WorkSeconds(20), RestSeconds(10), RoundCount(8), CycleCount(0))
	if err != ErrInvalidCycles {
		t.Fatalf("expected ErrInvalidCycles, got %v", err)
	}
}

func TestInvalidName(t *testing.T) {
	stages := []Stage{sampleStage(t, StageWarmup, 1, amrapConfig(t))}

	_, err := NewWOD(
		WODID("wod-1"),
		WODName("ab"),
		WODDescription("test"),
		stages,
		time.Now().UTC(),
	)
	if err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

func TestForTimeConfigOptionalTimeCap(t *testing.T) {
	cfg, err := NewForTimeConfig(RoundCount(5), nil)
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}
	if cfg.Rounds() != 5 {
		t.Fatalf("expected 5 rounds")
	}
}

func TestEMOMConfigValidation(t *testing.T) {
	_, err := NewEMOMConfig(IntervalSeconds(0), RoundCount(10))
	if err != ErrInvalidInterval {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}
