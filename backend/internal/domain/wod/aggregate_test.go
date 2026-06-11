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

func TestNewAMRAPWOD(t *testing.T) {
	cfg, err := NewAMRAPConfig(TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}

	wod, err := NewWOD(
		WODID("wod-1"),
		WODName("Open AMRAP"),
		WODDescription("test"),
		cfg,
		[]Movement{sampleMovement(t)},
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("unexpected wod error: %v", err)
	}

	if wod.Scoring().Kind() != ScoringRoundsReps {
		t.Fatalf("expected scoring rounds reps, got %s", wod.Scoring().Kind())
	}
}

func TestNewWODRequiresMovement(t *testing.T) {
	cfg, err := NewAMRAPConfig(TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}

	_, err = NewWOD(
		WODID("wod-1"),
		WODName("Open AMRAP"),
		WODDescription("test"),
		cfg,
		nil,
		time.Now().UTC(),
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
	cfg, err := NewAMRAPConfig(TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}

	_, err = NewWOD(
		WODID("wod-1"),
		WODName("ab"),
		WODDescription("test"),
		cfg,
		[]Movement{sampleMovement(t)},
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
