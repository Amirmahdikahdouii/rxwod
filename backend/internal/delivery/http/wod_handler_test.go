package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type handlerMemoryRepo struct {
	items map[domainwod.WODID]domainwod.Variant
}

func newHandlerMemoryRepo() *handlerMemoryRepo {
	return &handlerMemoryRepo{items: make(map[domainwod.WODID]domainwod.Variant)}
}

func (m *handlerMemoryRepo) Save(_ context.Context, variant domainwod.Variant) error {
	m.items[variant.ID()] = variant
	return nil
}

func (m *handlerMemoryRepo) FindByID(_ context.Context, id domainwod.WODID) (domainwod.Variant, error) {
	variant, ok := m.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return variant, nil
}

func (m *handlerMemoryRepo) List(_ context.Context) ([]domainwod.Variant, error) {
	items := make([]domainwod.Variant, 0, len(m.items))
	for _, variant := range m.items {
		items = append(items, variant)
	}
	return items, nil
}

func TestCreateWODHandler(t *testing.T) {
	repo := newHandlerMemoryRepo()
	service := appwod.NewService(repo, clock.System{}, idgen.UUIDGenerator{})
	router := NewRouter(service)

	body := `{
		"name": "Test AMRAP",
		"type": "AMRAP",
		"description": "desc",
		"config": { "timeCapSeconds": 900 },
		"movements": [{ "position": 1, "name": "Burpee", "reps": 21 }]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wods", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateWODValidation(t *testing.T) {
	repo := newHandlerMemoryRepo()
	service := appwod.NewService(repo, clock.System{}, idgen.UUIDGenerator{})
	router := NewRouter(service)

	body := `{
		"name": "ab",
		"type": "AMRAP",
		"description": "desc",
		"config": { "timeCapSeconds": 900 },
		"movements": [{ "position": 1, "name": "Burpee", "reps": 21 }]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wods", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
