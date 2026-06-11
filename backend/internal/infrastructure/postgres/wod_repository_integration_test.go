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
	movement, err := domainwod.NewMovement(domainwod.MovementID("mov-int-1"), 1, "Air Squat", &reps, nil, nil, "")
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}

	cfg, err := domainwod.NewTabataConfig(domainwod.WorkSeconds(20), domainwod.RestSeconds(10), domainwod.RoundCount(8), domainwod.CycleCount(1))
	if err != nil {
		t.Fatalf("config error: %v", err)
	}

	now := time.Now().UTC()
	wod, err := domainwod.NewWOD(domainwod.WODID("00000000-0000-4000-8000-000000000001"), domainwod.WODName("Integration Tabata"), domainwod.WODDescription("test"), cfg, []domainwod.Movement{movement}, now)
	if err != nil {
		t.Fatalf("wod error: %v", err)
	}

	variant := domainwod.NewSavedTabata(wod)
	if err := repo.Save(ctx, variant); err != nil {
		t.Fatalf("save error: %v", err)
	}

	found, err := repo.FindByID(ctx, wod.ID())
	if err != nil {
		t.Fatalf("find error: %v", err)
	}

	if found.Type() != domainwod.WODTypeTabata {
		t.Fatalf("expected tabata type")
	}
}
