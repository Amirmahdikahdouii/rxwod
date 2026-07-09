package wod

import (
	"testing"
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

const (
	testGymID  gym.GymID   = "gym-1"
	testUserID user.UserID = "user-1"
)

func sampleMovement(t *testing.T) Movement {
	t.Helper()
	reps := RepCount(21)
	movement, err := NewMovement(
		MovementID("mov-1"),
		1,
		"",
		"Burpee",
		"",
		nil,
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

func prescriptionMovement(t *testing.T) Movement {
	t.Helper()
	movement, err := NewMovement(
		MovementID("mov-2"),
		1,
		"D",
		"Crow + Knee Lift Off",
		"1-2X4 lift offs per side, rest as needed.",
		nil,
		nil,
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
		"",
		cfg,
		[]Movement{sampleMovement(t)},
	)
	if err != nil {
		t.Fatalf("unexpected stage error: %v", err)
	}
	return stage
}

func openStage(t *testing.T, kind StageKind, position int, instructions string) Stage {
	t.Helper()
	cfg, err := NewOpenConfig()
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}
	stage, err := NewStage(
		StageID("stage-open"),
		kind,
		position,
		instructions,
		cfg,
		[]Movement{prescriptionMovement(t)},
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
		testGymID,
		testUserID,
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

func TestOpenStageWithPrescriptionItems(t *testing.T) {
	stage := openStage(t, StageWarmup, 1, "")
	if stage.ScoringKind() != ScoringNone {
		t.Fatalf("expected NONE scoring for OPEN stage, got %s", stage.ScoringKind())
	}
	if stage.Type() != WODTypeOpen {
		t.Fatalf("expected OPEN type, got %s", stage.Type())
	}

	movements := stage.Movements()
	if movements[0].Label() != "D" {
		t.Fatalf("expected label D, got %s", movements[0].Label())
	}
	if movements[0].Prescription() == "" {
		t.Fatalf("expected prescription text")
	}
}

func TestMovementRequiresNameOrPrescription(t *testing.T) {
	_, err := NewMovement(MovementID("mov-1"), 1, "", "", "", nil, nil, nil, nil, "")
	if err != ErrMovementTextRequired {
		t.Fatalf("expected ErrMovementTextRequired, got %v", err)
	}

	movement, err := NewMovement(MovementID("mov-1"), 1, "", "", "Accumulate 20 reps.", nil, nil, nil, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.Prescription() != "Accumulate 20 reps." {
		t.Fatalf("expected prescription to be set")
	}
}

func TestMovementRejectsLoadValueWithoutUnit(t *testing.T) {
	load := LoadValue(28)
	_, err := NewMovement(MovementID("mov-1"), 1, "A", "Power Snatch", "", nil, nil, &load, nil, "")
	if err != ErrLoadValueRequiresUnit {
		t.Fatalf("expected ErrLoadValueRequiresUnit, got %v", err)
	}
}

func TestMovementRejectsInvalidSets(t *testing.T) {
	sets := SetCount(0)
	_, err := NewMovement(MovementID("mov-1"), 1, "", "Burpee", "", &sets, nil, nil, nil, "")
	if err != ErrInvalidSets {
		t.Fatalf("expected ErrInvalidSets, got %v", err)
	}
}

func TestNewWODRequiresStage(t *testing.T) {
	_, err := NewWOD(
		WODID("wod-1"),
		testGymID,
		testUserID,
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
		testGymID,
		testUserID,
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
		"",
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
		"",
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
		testGymID,
		testUserID,
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

func TestOpenConfigValidation(t *testing.T) {
	cfg, err := NewOpenConfig()
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}
	if cfg.ScoringKind() != ScoringNone {
		t.Fatalf("expected ScoringNone")
	}
}

func publishedWOD(t *testing.T) WOD {
	t.Helper()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	stages := []Stage{sampleStage(t, StageMetcon, 1, amrapConfig(t))}
	wod, err := NewWOD(
		WODID("wod-1"),
		testGymID,
		testUserID,
		WODName("Published Program"),
		WODDescription("test"),
		stages,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected wod error: %v", err)
	}
	scheduledDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	wod.SetScheduledDate(&scheduledDate, now)
	if err := wod.Publish(now); err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}
	return wod
}

func TestArchivePublishedWOD(t *testing.T) {
	wod := publishedWOD(t)
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	if err := wod.Archive(now); err != nil {
		t.Fatalf("unexpected archive error: %v", err)
	}
	if wod.Status() != WODStatusArchived {
		t.Fatalf("expected ARCHIVED status, got %s", wod.Status())
	}
	if !wod.UpdatedAt().Equal(now) {
		t.Fatalf("expected updatedAt %v, got %v", now, wod.UpdatedAt())
	}
}

func TestArchiveRejectsDraft(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	stages := []Stage{sampleStage(t, StageMetcon, 1, amrapConfig(t))}
	wod, err := NewWOD(
		WODID("wod-1"),
		testGymID,
		testUserID,
		WODName("Draft Program"),
		WODDescription("test"),
		stages,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected wod error: %v", err)
	}

	if err := wod.Archive(now); err != ErrInvalidStatusTransition {
		t.Fatalf("expected ErrInvalidStatusTransition, got %v", err)
	}
}

func TestArchiveRejectsAlreadyArchived(t *testing.T) {
	wod := publishedWOD(t)
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if err := wod.Archive(now); err != nil {
		t.Fatalf("unexpected archive error: %v", err)
	}

	if err := wod.Archive(now); err != ErrInvalidStatusTransition {
		t.Fatalf("expected ErrInvalidStatusTransition, got %v", err)
	}
}

func TestCanDelete(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	stages := []Stage{sampleStage(t, StageMetcon, 1, amrapConfig(t))}
	draft, err := NewWOD(
		WODID("wod-draft"),
		testGymID,
		testUserID,
		WODName("Draft Program"),
		WODDescription("test"),
		stages,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected wod error: %v", err)
	}
	if !draft.CanDelete() {
		t.Fatalf("expected draft to be deletable")
	}

	published := publishedWOD(t)
	if published.CanDelete() {
		t.Fatalf("expected published wod to not be deletable")
	}

	archived := publishedWOD(t)
	if err := archived.Archive(now); err != nil {
		t.Fatalf("unexpected archive error: %v", err)
	}
	if !archived.CanDelete() {
		t.Fatalf("expected archived wod to be deletable")
	}
}
