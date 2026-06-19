package wod

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
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

func (m *memoryRepo) FindByID(_ context.Context, gymID gym.GymID, id domainwod.WODID) (domainwod.WOD, error) {
	aggregate, ok := m.items[id]
	if !ok || aggregate.GymID() != gymID {
		return domainwod.WOD{}, errors.New("not found")
	}
	return aggregate, nil
}

func (m *memoryRepo) List(_ context.Context, gymID gym.GymID) ([]domainwod.WOD, error) {
	items := make([]domainwod.WOD, 0, len(m.items))
	for _, aggregate := range m.items {
		if aggregate.GymID() == gymID {
			items = append(items, aggregate)
		}
	}
	return items, nil
}

func (m *memoryRepo) ListCalendar(_ context.Context, gymID gym.GymID, from, to time.Time, includeDrafts bool) ([]CalendarEntry, error) {
	items := make([]CalendarEntry, 0)
	for _, aggregate := range m.items {
		if aggregate.GymID() != gymID {
			continue
		}
		scheduledDate := aggregate.ScheduledDate()
		if scheduledDate == nil {
			continue
		}
		if scheduledDate.Before(from) || scheduledDate.After(to) {
			continue
		}
		if !includeDrafts && aggregate.Status() != domainwod.WODStatusPublished {
			continue
		}
		items = append(items, CalendarEntry{
			ID:            aggregate.ID().String(),
			Name:          string(aggregate.Name()),
			Status:        aggregate.Status(),
			ScheduledDate: *scheduledDate,
		})
	}
	return items, nil
}

func testContext() context.Context {
	return appauthz.WithPrincipal(context.Background(), appauthz.Principal{
		UserID: user.UserID("user-1"),
		GymID:  gym.GymID("gym-1"),
		Role:   domainauthz.RoleOwner,
	})
}

func testContextForGym(gymID gym.GymID) context.Context {
	return appauthz.WithPrincipal(context.Background(), appauthz.Principal{
		UserID: user.UserID("user-1"),
		GymID:  gymID,
		Role:   domainauthz.RoleOwner,
	})
}

func TestServiceCreateMultiStageProgram(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	result, err := service.Create(testContext(), CreateWODCommand{
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

	_, err := service.Create(testContext(), CreateWODCommand{
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

	_, err := service.Create(testContext(), CreateWODCommand{
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

	_, err := service.Create(testContext(), CreateWODCommand{
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

	_, err := service.Create(testContext(), CreateWODCommand{
		Name:   "Empty Program",
		Stages: nil,
	})
	if !errors.Is(err, domainwod.ErrStageRequired) {
		t.Fatalf("expected ErrStageRequired, got %v", err)
	}
}

func TestServiceCreatePrescriptionProgram(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	result, err := service.Create(testContext(), CreateWODCommand{
		Name:        "Prescription Program",
		Description: "Coach plan",
		Stages: []StageInput{
			{
				Kind: domainwod.StageWarmup,
				Config: StageConfigInput{
					Type: domainwod.WODTypeOpen,
				},
				Movements: []MovementInput{{
					Position:     1,
					Label:        "D",
					Name:         "Crow + Knee Lift Off",
					Prescription: "1-2X4 lift offs per side, rest as needed.",
				}},
			},
			{
				Kind:         domainwod.StageStrength,
				Instructions: "Complete in 20 minutes.",
				Config: StageConfigInput{
					Type: domainwod.WODTypeOpen,
				},
				Movements: []MovementInput{{
					Position:     1,
					Label:        "B",
					Name:         "Close Grip Bench Press",
					Prescription: "3RM",
				}},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Stages[0].ScoringKind != domainwod.ScoringNone {
		t.Fatalf("expected NONE scoring for OPEN warmup, got %s", result.Stages[0].ScoringKind)
	}
	if result.StageCount != 2 {
		t.Fatalf("expected 2 stages, got %d", result.StageCount)
	}
}

func TestServiceUpdatePreservesIdentityAndTimestamps(t *testing.T) {
	repo := newMemoryRepo()
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	generator := &sequentialIDGen{}
	createService := NewService(repo, fixedClock{now: createdAt}, generator)

	created, err := createService.Create(testContext(), CreateWODCommand{
		Name:        "Original Program",
		Description: "Original notes",
		Stages: []StageInput{
			{
				Kind:      domainwod.StageWarmup,
				Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
				Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks", Sets: intPtr(3), Reps: intPtr(10)}},
			},
		},
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	updatedAt := createdAt.Add(2 * time.Hour)
	updateService := NewService(repo, fixedClock{now: updatedAt}, generator)
	updated, err := updateService.Update(testContext(), created.ID, CreateWODCommand{
		Name:        "Updated Program",
		Description: "Updated notes",
		Stages: []StageInput{
			{
				Kind:         domainwod.StageStrength,
				Instructions: "Complete in 20 minutes.",
				Config:       StageConfigInput{Type: domainwod.WODTypeOpen},
				Movements: []MovementInput{{
					Position:     1,
					Label:        "A",
					Name:         "Back Squat",
					Prescription: "Build to a heavy triple.",
					Sets:         intPtr(5),
					Reps:         intPtr(3),
				}},
			},
		},
	})
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	if updated.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, updated.ID)
	}
	if updated.Name != "Updated Program" || updated.Description != "Updated notes" {
		t.Fatalf("unexpected updated detail: %+v", updated)
	}
	if !updated.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected createdAt %s, got %s", createdAt, updated.CreatedAt)
	}
	if !updated.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("expected updatedAt %s, got %s", updatedAt, updated.UpdatedAt)
	}
	if updated.Status != domainwod.WODStatusDraft {
		t.Fatalf("expected draft status, got %s", updated.Status)
	}
	if len(updated.Stages) != 1 || updated.Stages[0].Kind != domainwod.StageStrength {
		t.Fatalf("unexpected stages: %+v", updated.Stages)
	}
	if updated.Stages[0].Movements[0].Sets == nil || *updated.Stages[0].Movements[0].Sets != 5 {
		t.Fatalf("expected sets to round trip, got %+v", updated.Stages[0].Movements[0])
	}
}

func TestServiceRejectsLoadValueWithoutUnit(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	load := 28.0
	_, err := service.Create(testContext(), CreateWODCommand{
		Name: "Load Error Program",
		Stages: []StageInput{
			{
				Kind:   domainwod.StageMetcon,
				Config: StageConfigInput{Type: domainwod.WODTypeAMRAP, TimeCapSeconds: intPtr(900)},
				Movements: []MovementInput{{
					Position:  1,
					Label:     "A",
					Name:      "Power Snatch",
					Reps:      intPtr(12),
					LoadValue: &load,
				}},
			},
		},
	})
	if err == nil {
		t.Fatal("expected load value without unit to fail")
	}
	if !errors.Is(err, domainwod.ErrLoadValueRequiresUnit) {
		t.Fatalf("expected ErrLoadValueRequiresUnit, got %v", err)
	}
	if !strings.Contains(err.Error(), "METCON stage, item A") {
		t.Fatalf("expected contextual error, got %v", err)
	}
}

func TestServiceCreateAcceptsZeroLoadValue(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	zeroLoad := 0.0
	_, err := service.Create(testContext(), CreateWODCommand{
		Name: "June 18",
		Stages: []StageInput{
			{
				Kind: domainwod.StageWarmup,
				Config: StageConfigInput{Type: domainwod.WODTypeOpen},
				Movements: []MovementInput{{
					Position: 1,
					Label:    "A",
					Name:     "Wall facing handstand plate stepup",
					Reps:     intPtr(20),
					Sets:     intPtr(1),
					LoadValue: &zeroLoad,
				}},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected zero load value to be accepted, got %v", err)
	}
}

func TestServiceRejectsCrossGymRead(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	created, err := service.Create(testContextForGym(gym.GymID("gym-1")), CreateWODCommand{
		Name: "Gym One Program",
		Stages: []StageInput{
			{
				Kind:      domainwod.StageWarmup,
				Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
				Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks"}},
			},
		},
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	_, err = service.GetByID(testContextForGym(gym.GymID("gym-2")), created.ID)
	if err == nil {
		t.Fatal("expected cross-gym read to fail")
	}
}

func intPtr(v int) *int { return &v }

func testContextForRole(role domainauthz.Role, userID user.UserID) context.Context {
	return appauthz.WithPrincipal(context.Background(), appauthz.Principal{
		UserID: userID,
		GymID:  gym.GymID("gym-1"),
		Role:   role,
	})
}

func scheduledDate(year int, month time.Month, day int) *time.Time {
	value := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &value
}

func TestServicePublishRequiresScheduledDate(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	created, err := service.Create(testContext(), CreateWODCommand{
		Name: "Draft Program",
		Stages: []StageInput{{
			Kind:      domainwod.StageWarmup,
			Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
			Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks"}},
		}},
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	_, err = service.Publish(testContext(), created.ID)
	if !errors.Is(err, domainwod.ErrScheduledDateRequired) {
		t.Fatalf("expected ErrScheduledDateRequired, got %v", err)
	}
}

func TestServiceAthleteListExcludesDrafts(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	created, err := service.Create(testContext(), CreateWODCommand{
		Name:          "Draft Program",
		ScheduledDate: scheduledDate(2026, time.June, 20),
		Stages: []StageInput{{
			Kind:      domainwod.StageWarmup,
			Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
			Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks"}},
		}},
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if _, err := service.Publish(testContext(), created.ID); err != nil {
		t.Fatalf("publish error: %v", err)
	}

	_, err = service.Create(testContext(), CreateWODCommand{
		Name:          "Another Draft",
		ScheduledDate: scheduledDate(2026, time.June, 21),
		Stages: []StageInput{{
			Kind:      domainwod.StageWarmup,
			Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
			Movements: []MovementInput{{Position: 1, Name: "Air Squat"}},
		}},
	})
	if err != nil {
		t.Fatalf("create draft error: %v", err)
	}

	items, err := service.List(testContextForRole(domainauthz.RoleAthlete, user.UserID("athlete-1")), nil)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 published item for athlete, got %d", len(items))
	}
}

func TestServiceCoachCannotUpdateAnotherCoachDraft(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	created, err := service.Create(testContextForRole(domainauthz.RoleCoach, user.UserID("coach-1")), CreateWODCommand{
		Name: "Coach One Draft",
		Stages: []StageInput{{
			Kind:      domainwod.StageWarmup,
			Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
			Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks"}},
		}},
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	_, err = service.Update(testContextForRole(domainauthz.RoleCoach, user.UserID("coach-2")), created.ID, CreateWODCommand{
		Name: "Stolen Draft",
		Stages: []StageInput{{
			Kind:      domainwod.StageWarmup,
			Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
			Movements: []MovementInput{{Position: 1, Name: "Air Squat"}},
		}},
	})
	if !errors.Is(err, appauthz.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestServiceCalendarGroupsMultiplePlansPerDate(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo, fixedClock{now: time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC)}, &sequentialIDGen{})

	for _, name := range []string{"Morning", "Evening"} {
		_, err := service.Create(testContext(), CreateWODCommand{
			Name:          name,
			ScheduledDate: scheduledDate(2026, time.June, 20),
			Stages: []StageInput{{
				Kind:      domainwod.StageWarmup,
				Config:    StageConfigInput{Type: domainwod.WODTypeOpen},
				Movements: []MovementInput{{Position: 1, Name: "Jumping Jacks"}},
			}},
		})
		if err != nil {
			t.Fatalf("create error: %v", err)
		}
	}

	days, err := service.Calendar(testContext(), time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("calendar error: %v", err)
	}
	if len(days) != 1 || days[0].DraftCount != 2 {
		t.Fatalf("expected one day with two drafts, got %+v", days)
	}
}
