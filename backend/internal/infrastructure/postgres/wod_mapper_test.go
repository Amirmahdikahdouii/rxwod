package postgres

import (
	"testing"
	"time"

	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

func TestVariantToRecordRoundTripAMRAP(t *testing.T) {
	reps := domainwod.RepCount(21)
	movement, err := domainwod.NewMovement(domainwod.MovementID("mov-1"), 1, "Burpee", &reps, nil, nil, "")
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}

	cfg, err := domainwod.NewAMRAPConfig(domainwod.TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("config error: %v", err)
	}

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	wod, err := domainwod.NewWOD(domainwod.WODID("wod-1"), domainwod.WODName("Test"), domainwod.WODDescription("desc"), cfg, []domainwod.Movement{movement}, now)
	if err != nil {
		t.Fatalf("wod error: %v", err)
	}

	variant := domainwod.NewSavedAMRAP(wod)
	record, movements, err := variantToRecord(variant)
	if err != nil {
		t.Fatalf("to record error: %v", err)
	}

	restored, err := recordToVariant(record, movements)
	if err != nil {
		t.Fatalf("from record error: %v", err)
	}

	if restored.Type() != domainwod.WODTypeAMRAP {
		t.Fatalf("expected AMRAP type")
	}
}
