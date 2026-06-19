package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	"github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
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

func (m *handlerMemoryRepo) FindByID(_ context.Context, gymID gym.GymID, id domainwod.WODID) (domainwod.WOD, error) {
	aggregate, ok := m.items[id]
	if !ok || aggregate.GymID() != gymID {
		return domainwod.WOD{}, postgres.ErrNotFound
	}
	return aggregate, nil
}

func (m *handlerMemoryRepo) List(_ context.Context, gymID gym.GymID) ([]domainwod.WOD, error) {
	items := make([]domainwod.WOD, 0, len(m.items))
	for _, aggregate := range m.items {
		if aggregate.GymID() == gymID {
			items = append(items, aggregate)
		}
	}
	return items, nil
}

func (m *handlerMemoryRepo) ListCalendar(_ context.Context, gymID gym.GymID, from, to time.Time, includeDrafts bool) ([]appwod.CalendarEntry, error) {
	items := make([]appwod.CalendarEntry, 0)
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
		items = append(items, appwod.CalendarEntry{
			ID:            aggregate.ID().String(),
			Name:          string(aggregate.Name()),
			Status:        aggregate.Status(),
			ScheduledDate: *scheduledDate,
		})
	}
	return items, nil
}

func newTestRouter() *echo.Echo {
	repo := newHandlerMemoryRepo()
	service := appwod.NewService(repo, clock.System{}, idgen.UUIDGenerator{})
	handler := NewWODHandler(service)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := appauthz.WithPrincipal(c.Request().Context(), appauthz.Principal{
				UserID: user.UserID("user-1"),
				GymID:  gym.GymID("gym-1"),
				Role:   domainauthz.RoleOwner,
			})
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
	e.POST("/api/v1/wods", handler.Create)
	e.PUT("/api/v1/wods/:id", handler.Update)
	return e
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

func TestCreateWODLoadValueWithoutUnit(t *testing.T) {
	router := newTestRouter()

	body := `{
		"name": "Load Error Program",
		"description": "",
		"stages": [
			{
				"kind": "METCON",
				"type": "AMRAP",
				"config": { "timeCapSeconds": 900 },
				"movements": [{ "position": 1, "label": "A", "name": "Power Snatch", "reps": 12, "loadValue": 28 }]
			}
		]
	}`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}

	var response ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if !strings.Contains(response.Error, "METCON stage, item A") {
		t.Fatalf("expected contextual error, got %q", response.Error)
	}
	if !strings.Contains(response.Error, "load value requires a load unit") {
		t.Fatalf("expected load unit error, got %q", response.Error)
	}
}

func TestCreateWODZeroLoadValue(t *testing.T) {
	router := newTestRouter()

	body := `{
		"name": "June 18",
		"description": "",
		"stages": [
			{
				"kind": "WARMUP",
				"type": "OPEN",
				"config": {},
				"movements": [{
					"position": 1,
					"label": "A",
					"name": "Wall facing handstand plate stepup",
					"reps": 20,
					"sets": 1,
					"loadValue": 0
				}]
			}
		]
	}`

	rec := postWOD(t, router, body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}
}
