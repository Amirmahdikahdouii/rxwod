package wod

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time { return f.now }

type sequentialIDGen struct {
	counter int
}

func (s *sequentialIDGen) NewID() string {
	s.counter++
	return fmt.Sprintf("generated-id-%d", s.counter)
}

type memoryRepo struct {
	items map[domainwod.WODID]domainwod.WOD
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{items: make(map[domainwod.WODID]domainwod.WOD)}
}

func (m *memoryRepo) Save(_ context.Context, aggregate domainwod.WOD) error {
	m.items[aggregate.ID()] = aggregate
	return nil
}

func (m *memoryRepo) FindByID(_ context.Context, id domainwod.WODID) (domainwod.WOD, error) {
	aggregate, ok := m.items[id]
	if !ok {
		return domainwod.WOD{}, errors.New("not found")
	}
	return aggregate, nil
}

func (m *memoryRepo) List(_ context.Context) ([]domainwod.WOD, error) {
	items := make([]domainwod.WOD, 0, len(m.items))
	for _, aggregate := range m.items {
		items = append(items, aggregate)
	}
	return items, nil
}

func TestServiceCreateMultiStageProgram(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	result, err := service.Create(context.Background(), CreateWODCommand{
		Name:        "Monday Session",
		Description: "Full class plan",
		Stages: []StageInput{
			{
				Kind:      domainwod.StageWarmup,
				Config:    StageConfigInput{Type: domainwod.WODTypeForTime, Rounds: intPtr(2)},
				Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks", Reps: intPtr(20)}},
			},
			{
				Kind:      domainwod.StageMetcon,
				Config:    StageConfigInput{Type: domainwod.WODTypeAMRAP, TimeCapSeconds: intPtr(900)},
				Movements: []MovementInput{{Position: 1, Name: "Burpee", Reps: intPtr(21)}},
			},
			{
				Kind: domainwod.StageCooldown,
				Config: StageConfigInput{
					Type:        domainwod.WODTypeTabata,
					WorkSeconds: intPtr(20),
					RestSeconds: intPtr(10),
					Rounds:      intPtr(8),
					Cycles:      intPtr(1),
				},
				Movements: []MovementInput{{Position: 1, Name: "Plank", Reps: intPtr(1)}},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StageCount != 3 {
		t.Fatalf("expected 3 stages, got %d", result.StageCount)
	}
	if result.Stages[0].Kind != domainwod.StageWarmup || result.Stages[0].ScoringKind != domainwod.ScoringTimeToComplete {
		t.Fatalf("unexpected warmup stage summary: %+v", result.Stages[0])
	}
	if result.Stages[1].ScoringKind != domainwod.ScoringRoundsReps {
		t.Fatalf("expected metcon rounds-reps scoring, got %s", result.Stages[1].ScoringKind)
	}
	if result.Stages[2].ScoringKind != domainwod.ScoringTotalReps {
		t.Fatalf("expected cooldown total-reps scoring, got %s", result.Stages[2].ScoringKind)
	}
	if len(repo.items) != 1 {
		t.Fatalf("expected one saved wod")
	}
}

func TestServiceCreateRejectsInvalidStageKind(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, clock.System{}, idgen.UUIDGenerator{})

	_, err := service.Create(context.Background(), CreateWODCommand{
		Name: "Bad Kind",
		Stages: []StageInput{
			{
				Kind:      domainwod.StageKind("INVALID"),
				Config:    StageConfigInput{Type: domainwod.WODTypeAMRAP, TimeCapSeconds: intPtr(900)},
				Movements: []MovementInput{{Position: 1, Name: "Burpee", Reps: intPtr(21)}},
			},
		},
	})
	if !errors.Is(err, domainwod.ErrInvalidStageKind) {
		t.Fatalf("expected ErrInvalidStageKind, got %v", err)
	}
}

func TestServiceCreateRejectsMissingConfigField(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, clock.System{}, idgen.UUIDGenerator{})

	_, err := service.Create(context.Background(), CreateWODCommand{
		Name: "Missing Field",
		Stages: []StageInput{
			{
				Kind:      domainwod.StageMetcon,
				Config:    StageConfigInput{Type: domainwod.WODTypeAMRAP},
				Movements: []MovementInput{{Position: 1, Name: "Burpee", Reps: intPtr(21)}},
			},
		},
	})
	if !errors.Is(err, ErrMissingConfigField) {
		t.Fatalf("expected ErrMissingConfigField, got %v", err)
	}
}

func TestServiceCreateRejectsTabataValidation(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, clock.System{}, idgen.UUIDGenerator{})

	_, err := service.Create(context.Background(), CreateWODCommand{
		Name: "Tabata Test",
		Stages: []StageInput{
			{
				Kind: domainwod.StageMetcon,
				Config: StageConfigInput{
					Type:        domainwod.WODTypeTabata,
					WorkSeconds: intPtr(20),
					RestSeconds: intPtr(10),
					Rounds:      intPtr(8),
					Cycles:      intPtr(0),
				},
				Movements: []MovementInput{{Position: 1, Name: "Air Squat", Reps: intPtr(10)}},
			},
		},
	})
	if !errors.Is(err, domainwod.ErrInvalidCycles) {
		t.Fatalf("expected ErrInvalidCycles, got %v", err)
	}
}

func TestServiceCreateRequiresStage(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, clock.System{}, idgen.UUIDGenerator{})

	_, err := service.Create(context.Background(), CreateWODCommand{
		Name:   "Empty Program",
		Stages: nil,
	})
	if !errors.Is(err, domainwod.ErrStageRequired) {
		t.Fatalf("expected ErrStageRequired, got %v", err)
	}
}

func intPtr(v int) *int { return &v }
