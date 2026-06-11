package wod

import (
	"context"
	"errors"
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
	return "generated-id-" + string(rune('0'+s.counter))
}

type memoryRepo struct {
	items map[domainwod.WODID]domainwod.Variant
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{items: make(map[domainwod.WODID]domainwod.Variant)}
}

func (m *memoryRepo) Save(_ context.Context, variant domainwod.Variant) error {
	m.items[variant.ID()] = variant
	return nil
}

func (m *memoryRepo) FindByID(_ context.Context, id domainwod.WODID) (domainwod.Variant, error) {
	variant, ok := m.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return variant, nil
}

func (m *memoryRepo) List(_ context.Context) ([]domainwod.Variant, error) {
	items := make([]domainwod.Variant, 0, len(m.items))
	for _, variant := range m.items {
		items = append(items, variant)
	}
	return items, nil
}

func TestServiceCreateAMRAP(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	result, err := service.Create(context.Background(), CreateCommand{
		Type: domainwod.WODTypeAMRAP,
		AMRAP: &CreateAMRAPCommand{
			Name:        "Test AMRAP",
			Description: "desc",
			TimeCap:     900,
			Movements: []MovementInput{
				{Position: 1, Name: "Burpee", Reps: intPtr(21)},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ScoringKind != domainwod.ScoringRoundsReps {
		t.Fatalf("expected rounds reps scoring")
	}
	if len(repo.items) != 1 {
		t.Fatalf("expected one saved wod")
	}
}

func TestServiceCreateTabataValidation(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, clock.System{}, idgen.UUIDGenerator{})

	_, err := service.Create(context.Background(), CreateCommand{
		Type: domainwod.WODTypeTabata,
		Tabata: &CreateTabataCommand{
			Name:        "Tabata Test",
			WorkSeconds: 20,
			RestSeconds: 10,
			Rounds:      8,
			Cycles:      0,
			Movements: []MovementInput{
				{Position: 1, Name: "Air Squat"},
			},
		},
	})
	if !errors.Is(err, domainwod.ErrInvalidCycles) {
		t.Fatalf("expected ErrInvalidCycles, got %v", err)
	}
}

func intPtr(v int) *int { return &v }
