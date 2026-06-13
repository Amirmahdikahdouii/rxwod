package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/infrastructure/postgres"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type handlerMemoryRepo struct {
	items map[domainwod.WODID]domainwod.WOD
}

func newHandlerMemoryRepo() *handlerMemoryRepo {
	return &handlerMemoryRepo{items: make(map[domainwod.WODID]domainwod.WOD)}
}

func (m *handlerMemoryRepo) Save(_ context.Context, aggregate domainwod.WOD) error {
	m.items[aggregate.ID()] = aggregate
	return nil
}

func (m *handlerMemoryRepo) FindByID(_ context.Context, id domainwod.WODID) (domainwod.WOD, error) {
	aggregate, ok := m.items[id]
	if !ok {
		return domainwod.WOD{}, postgres.ErrNotFound
	}
	return aggregate, nil
}

func (m *handlerMemoryRepo) List(_ context.Context) ([]domainwod.WOD, error) {
	items := make([]domainwod.WOD, 0, len(m.items))
	for _, aggregate := range m.items {
		items = append(items, aggregate)
	}
	return items, nil
}

func newTestRouter() *echo.Echo {
	repo := newHandlerMemoryRepo()
	service := appwod.NewService(repo, clock.System{}, idgen.UUIDGenerator{})
	return NewRouter(service)
}

func postWOD(t *testing.T, router *echo.Echo, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wods", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func putWOD(t *testing.T, router *echo.Echo, id string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/wods/"+id, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCreateMultiStageWODHandler(t *testing.T) {
	router := newTestRouter()

	body := `{
		"name": "Monday Session",
		"description": "Full class plan",
		"stages": [
			{ "kind": "WARMUP", "type": "FORTIME", "config": { "rounds": 2 }, "movements": [{ "position": 1, "name": "Jumping Jacks", "reps": 20 }] },
			{ "kind": "METCON", "type": "AMRAP", "config": { "timeCapSeconds": 900 }, "movements": [{ "position": 1, "name": "Burpee", "reps": 21 }] },
			{ "kind": "COOLDOWN", "type": "TABATA", "config": { "workSeconds": 20, "restSeconds": 10, "rounds": 8, "cycles": 1 }, "movements": [{ "position": 1, "name": "Plank", "reps": 1 }] }
		]
	}`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestUpdateWODHandler(t *testing.T) {
	router := newTestRouter()
	createBody := `{
		"name": "Original Program",
		"description": "Full class plan",
		"stages": [
			{ "kind": "WARMUP", "type": "OPEN", "config": {}, "movements": [{ "position": 1, "name": "Jumping Jacks", "sets": 3, "reps": 10 }] }
		]
	}`

	createRec := postWOD(t, router, createBody)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", createRec.Code, createRec.Body.String())
	}
	var created CreateWODResponse
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	updateBody := `{
		"name": "Updated Program",
		"description": "Updated plan",
		"stages": [
			{ "kind": "STRENGTH", "type": "OPEN", "instructions": "Complete in 20 minutes.", "config": {}, "movements": [{ "position": 1, "label": "A", "name": "Back Squat", "sets": 5, "reps": 3 }] }
		]
	}`

	updateRec := putWOD(t, router, created.ID, updateBody)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", updateRec.Code, updateRec.Body.String())
	}
	var updated WODDetailResponse
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.ID != created.ID || updated.Name != "Updated Program" {
		t.Fatalf("unexpected updated response: %+v", updated)
	}
	if updated.Stages[0].Movements[0].Sets == nil || *updated.Stages[0].Movements[0].Sets != 5 {
		t.Fatalf("expected sets in response, got %+v", updated.Stages[0].Movements[0])
	}
}

func TestUpdateWODNotFound(t *testing.T) {
	router := newTestRouter()
	body := `{
		"name": "Missing Program",
		"description": "",
		"stages": [
			{ "kind": "WARMUP", "type": "OPEN", "config": {}, "movements": [{ "position": 1, "name": "Jumping Jacks" }] }
		]
	}`

	rec := putWOD(t, router, "missing", body)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestUpdateWODInvalidSets(t *testing.T) {
	router := newTestRouter()
	createBody := `{
		"name": "Original Program",
		"description": "",
		"stages": [
			{ "kind": "WARMUP", "type": "OPEN", "config": {}, "movements": [{ "position": 1, "name": "Jumping Jacks" }] }
		]
	}`
	createRec := postWOD(t, router, createBody)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", createRec.Code, createRec.Body.String())
	}
	var created CreateWODResponse
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	updateBody := `{
		"name": "Invalid Program",
		"description": "",
		"stages": [
			{ "kind": "WARMUP", "type": "OPEN", "config": {}, "movements": [{ "position": 1, "name": "Jumping Jacks", "sets": 0 }] }
		]
	}`

	rec := putWOD(t, router, created.ID, updateBody)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateWODEmptyStages(t *testing.T) {
	router := newTestRouter()

	body := `{ "name": "Empty Program", "description": "", "stages": [] }`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateWODInvalidStageKind(t *testing.T) {
	router := newTestRouter()

	body := `{
		"name": "Bad Kind",
		"description": "",
		"stages": [
			{ "kind": "INVALID", "type": "AMRAP", "config": { "timeCapSeconds": 900 }, "movements": [{ "position": 1, "name": "Burpee", "reps": 21 }] }
		]
	}`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateWODInvalidName(t *testing.T) {
	router := newTestRouter()

	body := `{
		"name": "ab",
		"description": "",
		"stages": [
			{ "kind": "METCON", "type": "AMRAP", "config": { "timeCapSeconds": 900 }, "movements": [{ "position": 1, "name": "Burpee", "reps": 21 }] }
		]
	}`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
