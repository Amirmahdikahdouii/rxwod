package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type handlerGymMemoryRepo struct {
	members []appgym.MemberDTO
}

func (m *handlerGymMemoryRepo) CreateGymWithOwner(context.Context, domaingym.Gym, domaingym.Membership) error {
	return nil
}

func (m *handlerGymMemoryRepo) ListForUser(context.Context, user.UserID) ([]appgym.WorkspaceDTO, error) {
	return nil, nil
}

func (m *handlerGymMemoryRepo) FindByID(context.Context, domaingym.GymID) (domaingym.Gym, error) {
	return domaingym.Gym{}, nil
}

func (m *handlerGymMemoryRepo) FindActiveMembership(context.Context, domaingym.GymID, user.UserID) (domaingym.Membership, error) {
	return domaingym.Membership{}, nil
}

func (m *handlerGymMemoryRepo) FindMembership(context.Context, domaingym.GymID, user.UserID) (domaingym.Membership, error) {
	return domaingym.Membership{}, nil
}

func (m *handlerGymMemoryRepo) FindMember(context.Context, domaingym.GymID, user.UserID) (appgym.MemberDTO, error) {
	return appgym.MemberDTO{}, nil
}

func (m *handlerGymMemoryRepo) ListMembers(_ context.Context, _ domaingym.GymID, filter appgym.ListMembersFilter) (appgym.ListMembersResult, error) {
	total := len(m.members)
	start := (filter.Page - 1) * filter.Limit
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}
	return appgym.ListMembersResult{Items: m.members[start:end], Total: total}, nil
}

func (m *handlerGymMemoryRepo) DeleteMembership(context.Context, domaingym.GymID, user.UserID) error {
	return nil
}

func (m *handlerGymMemoryRepo) FindUserByEmail(context.Context, user.Email) (user.User, error) {
	return user.User{}, nil
}

func (m *handlerGymMemoryRepo) FindUserByID(context.Context, user.UserID) (user.User, error) {
	return user.User{}, nil
}

func (m *handlerGymMemoryRepo) UpsertMembership(context.Context, domaingym.Membership) error {
	return nil
}

func (m *handlerGymMemoryRepo) SaveInvitation(context.Context, domaingym.Invitation, string) error {
	return nil
}

func (m *handlerGymMemoryRepo) FindPendingInvitationsByEmail(context.Context, user.Email) ([]domaingym.Invitation, error) {
	return nil, nil
}

func (m *handlerGymMemoryRepo) FindPendingInvitationByTokenHash(context.Context, domaingym.GymID, string) (domaingym.Invitation, error) {
	return domaingym.Invitation{}, nil
}

func (m *handlerGymMemoryRepo) FindInvitationPreviewByTokenHash(context.Context, string) (appgym.InvitationPreviewDTO, error) {
	return appgym.InvitationPreviewDTO{}, nil
}

func (m *handlerGymMemoryRepo) AcceptInvitationWithMembership(context.Context, domaingym.Invitation, domaingym.Membership) error {
	return nil
}

func newGymTestRouter(repo *handlerGymMemoryRepo) *echo.Echo {
	service := appgym.NewService(repo, nil, clock.System{}, idgen.UUIDGenerator{}, time.Hour)
	handler := NewGymHandler(service)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := appauthz.WithPrincipal(c.Request().Context(), appauthz.Principal{
				UserID: user.UserID("user-1"),
				GymID:  domaingym.GymID("gym-1"),
				Role:   domainauthz.RoleOwner,
			})
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
	e.GET("/api/v1/gyms/:gymId/members", handler.Members)
	return e
}

func getMembers(t *testing.T, router *echo.Echo, query string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gyms/gym-1/members"+query, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func seedMembers(count int) []appgym.MemberDTO {
	members := make([]appgym.MemberDTO, 0, count)
	for i := 0; i < count; i++ {
		members = append(members, appgym.MemberDTO{
			UserID:      "user-" + string(rune('a'+i)),
			Email:       "member@example.com",
			DisplayName: "Member",
			Role:        domainauthz.RoleAthlete,
			Status:      domaingym.MembershipStatusActive,
		})
	}
	return members
}

func TestListMembersPagination(t *testing.T) {
	repo := &handlerGymMemoryRepo{members: seedMembers(3)}
	router := newGymTestRouter(repo)

	rec := getMembers(t, router, "?page=1&limit=2")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var page PaginatedMemberResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(page.Data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(page.Data))
	}
	if page.Meta.Page != 1 || page.Meta.Limit != 2 || page.Meta.Total != 3 || page.Meta.TotalPages != 2 {
		t.Fatalf("unexpected meta: %+v", page.Meta)
	}
}

func TestListMembersDefaultsWithoutParams(t *testing.T) {
	repo := &handlerGymMemoryRepo{members: seedMembers(3)}
	router := newGymTestRouter(repo)

	rec := getMembers(t, router, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var page PaginatedMemberResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(page.Data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(page.Data))
	}
	if page.Meta.Page != defaultPage || page.Meta.Limit != defaultLimit || page.Meta.Total != 3 || page.Meta.TotalPages != 1 {
		t.Fatalf("unexpected meta: %+v", page.Meta)
	}
}

func TestListMembersInvalidPageParam(t *testing.T) {
	repo := &handlerGymMemoryRepo{members: seedMembers(1)}
	router := newGymTestRouter(repo)

	rec := getMembers(t, router, "?page=0")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestListMembersInvalidLimitParam(t *testing.T) {
	repo := &handlerGymMemoryRepo{members: seedMembers(1)}
	router := newGymTestRouter(repo)

	rec := getMembers(t, router, "?limit=101")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}
