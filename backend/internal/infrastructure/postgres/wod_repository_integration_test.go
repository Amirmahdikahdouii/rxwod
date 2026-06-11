package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

func TestWODRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	db, err := NewDB(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	repo := NewWODRepository(db)

	reps := domainwod.RepCount(15)
	warmupMovement, err := domainwod.NewMovement(domainwod.MovementID("00000000-0000-4000-8000-0000000000a1"), 1, "", "Jumping Jacks", "", &reps, nil, nil, "")
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}
	metconMovement, err := domainwod.NewMovement(domainwod.MovementID("00000000-0000-4000-8000-0000000000a2"), 1, "", "Air Squat", "", &reps, nil, nil, "")
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}

	forTime, err := domainwod.NewForTimeConfig(domainwod.RoundCount(2), nil)
	if err != nil {
		t.Fatalf("config error: %v", err)
	}
	tabata, err := domainwod.NewTabataConfig(domainwod.WorkSeconds(20), domainwod.RestSeconds(10), domainwod.RoundCount(8), domainwod.CycleCount(1))
	if err != nil {
		t.Fatalf("config error: %v", err)
	}

	warmup, err := domainwod.NewStage(domainwod.StageID("00000000-0000-4000-8000-0000000000b1"), domainwod.StageWarmup, 1, "", forTime, []domainwod.Movement{warmupMovement})
	if err != nil {
		t.Fatalf("stage error: %v", err)
	}
	metcon, err := domainwod.NewStage(domainwod.StageID("00000000-0000-4000-8000-0000000000b2"), domainwod.StageMetcon, 2, "", tabata, []domainwod.Movement{metconMovement})
	if err != nil {
		t.Fatalf("stage error: %v", err)
	}

	now := time.Now().UTC()
	wod, err := domainwod.NewWOD(
		domainwod.WODID("00000000-0000-4000-8000-000000000001"),
		domainwod.WODName("Integration Program"),
		domainwod.WODDescription("test"),
		[]domainwod.Stage{warmup, metcon},
		now,
	)
	if err != nil {
		t.Fatalf("wod error: %v", err)
	}

	if err := repo.Save(ctx, wod); err != nil {
		t.Fatalf("save error: %v", err)
	}

	found, err := repo.FindByID(ctx, wod.ID())
	if err != nil {
		t.Fatalf("find error: %v", err)
	}

	stages := found.Stages()
	if len(stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(stages))
	}
	if stages[0].Type() != domainwod.WODTypeForTime {
		t.Fatalf("expected first stage ForTime, got %s", stages[0].Type())
	}
	if stages[1].Type() != domainwod.WODTypeTabata {
		t.Fatalf("expected second stage Tabata, got %s", stages[1].Type())
	}
}
