package postgres

import (
	"testing"
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

const (
	testGymID  gym.GymID   = "gym-1"
	testUserID user.UserID = "user-1"
)

func buildStage(t *testing.T, id string, kind domainwod.StageKind, position int, cfg domainwod.Config) domainwod.Stage {
	t.Helper()
	reps := domainwod.RepCount(21)
	movement, err := domainwod.NewMovement(domainwod.MovementID(id+"-mov-1"), 1, "", "Burpee", "", nil, &reps, nil, nil, "")
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}
	stage, err := domainwod.NewStage(domainwod.StageID(id), kind, position, "", cfg, []domainwod.Movement{movement})
	if err != nil {
		t.Fatalf("stage error: %v", err)
	}
	return stage
}

func TestWODRecordsRoundTripMultiStage(t *testing.T) {
	amrap, err := domainwod.NewAMRAPConfig(domainwod.TimeCapSeconds(900))
	if err != nil {
		t.Fatalf("config error: %v", err)
	}
	forTime, err := domainwod.NewForTimeConfig(domainwod.RoundCount(2), nil)
	if err != nil {
		t.Fatalf("config error: %v", err)
	}
	tabata, err := domainwod.NewTabataConfig(domainwod.WorkSeconds(20), domainwod.RestSeconds(10), domainwod.RoundCount(8), domainwod.CycleCount(1))
	if err != nil {
		t.Fatalf("config error: %v", err)
	}

	stages := []domainwod.Stage{
		buildStage(t, "stage-1", domainwod.StageWarmup, 1, forTime),
		buildStage(t, "stage-2", domainwod.StageMetcon, 2, amrap),
		buildStage(t, "stage-3", domainwod.StageCooldown, 3, tabata),
	}

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	wod, err := domainwod.NewWOD(domainwod.WODID("wod-1"), testGymID, testUserID, domainwod.WODName("Monday Session"), domainwod.WODDescription("desc"), stages, now)
	if err != nil {
		t.Fatalf("wod error: %v", err)
	}

	record, stageRecords, movementRecords, err := wodToRecords(wod)
	if err != nil {
		t.Fatalf("to records error: %v", err)
	}
	if len(stageRecords) != 3 {
		t.Fatalf("expected 3 stage records, got %d", len(stageRecords))
	}
	if record.GymID != testGymID.String() || record.CreatedBy != testUserID.String() {
		t.Fatalf("expected gym and creator IDs to round trip, got %+v", record)
	}
	if len(movementRecords) != 3 {
		t.Fatalf("expected 3 movement records, got %d", len(movementRecords))
	}

	movementsByStage := make(map[string][]movementRecord)
	for _, m := range movementRecords {
		movementsByStage[m.StageID] = append(movementsByStage[m.StageID], m)
	}

	restored, err := recordsToWOD(record, stageRecords, movementsByStage)
	if err != nil {
		t.Fatalf("from records error: %v", err)
	}

	restoredStages := restored.Stages()
	if len(restoredStages) != 3 {
		t.Fatalf("expected 3 restored stages, got %d", len(restoredStages))
	}
	if restoredStages[0].Kind() != domainwod.StageWarmup || restoredStages[0].Type() != domainwod.WODTypeForTime {
		t.Fatalf("unexpected warmup stage: %+v", restoredStages[0])
	}
	if restoredStages[1].Type() != domainwod.WODTypeAMRAP {
		t.Fatalf("expected metcon AMRAP, got %s", restoredStages[1].Type())
	}
	if restoredStages[2].Type() != domainwod.WODTypeTabata {
		t.Fatalf("expected cooldown Tabata, got %s", restoredStages[2].Type())
	}
}

func TestWODRecordsRoundTripOpenPrescriptions(t *testing.T) {
	openCfg, err := domainwod.NewOpenConfig()
	if err != nil {
		t.Fatalf("config error: %v", err)
	}
	sets := domainwod.SetCount(3)
	reps := domainwod.RepCount(10)

	prescriptionMovement, err := domainwod.NewMovement(
		domainwod.MovementID("mov-rx"),
		1,
		"D",
		"Crow + Knee Lift Off",
		"1-2X4 lift offs per side, rest as needed.",
		&sets,
		&reps,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}

	stage, err := domainwod.NewStage(
		domainwod.StageID("stage-open"),
		domainwod.StageWarmup,
		1,
		"",
		openCfg,
		[]domainwod.Movement{prescriptionMovement},
	)
	if err != nil {
		t.Fatalf("stage error: %v", err)
	}

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	wod, err := domainwod.NewWOD(
		domainwod.WODID("wod-open"),
		testGymID,
		testUserID,
		domainwod.WODName("Prescription Program"),
		domainwod.WODDescription("desc"),
		[]domainwod.Stage{stage},
		now,
	)
	if err != nil {
		t.Fatalf("wod error: %v", err)
	}

	record, stageRecords, movementRecords, err := wodToRecords(wod)
	if err != nil {
		t.Fatalf("to records error: %v", err)
	}

	movementsByStage := map[string][]movementRecord{
		stageRecords[0].ID: movementRecords,
	}

	restored, err := recordsToWOD(record, stageRecords, movementsByStage)
	if err != nil {
		t.Fatalf("from records error: %v", err)
	}

	got := restored.Stages()[0].Movements()[0]
	if got.Label() != "D" || got.Prescription() == "" {
		t.Fatalf("expected prescription movement, got label=%q prescription=%q", got.Label(), got.Prescription())
	}
	if got.Sets() == nil || *got.Sets() != sets {
		t.Fatalf("expected sets %d, got %v", sets, got.Sets())
	}
	if restored.Stages()[0].ScoringKind() != domainwod.ScoringNone {
		t.Fatalf("expected NONE scoring, got %s", restored.Stages()[0].ScoringKind())
	}
}
