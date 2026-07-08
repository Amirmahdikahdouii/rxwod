# RXWOD — Product & Technical Context

> Generated for use as background context in an external PM chatbot. Summarizes the current state of the codebase so product decisions can be grounded in what is actually built vs. what is planned.

## 1. Project Overview

**RXWOD** is a full-stack SaaS application for CrossFit-style gyms to plan, schedule, and publish daily workout programs ("WODs" — Workouts of the Day).

Core value proposition: a gym owner or coach can build a structured, multi-stage training session (warm-up → strength → core → metcon → cooldown), assign a scoring/workout format to each stage (AMRAP, For Time, Tabata, EMOM, or free-form "Open" prescriptions), schedule it on a calendar, and publish it so athletes in that gym can view it.

The product is explicitly **multi-tenant**: it is organized around "gyms" (called "workspaces" in the frontend), each with its own members, roles, and WOD library. A single user account can belong to multiple gyms with different roles in each.

The README describes this as a **"Phase 1 WOD Creator"** — an early-stage product focused on program authoring and gym/team management, not yet on athlete-facing workout tracking, results logging, or leaderboards.

## 2. Tech Stack

**Backend:** Go 1.25, Echo v4 (HTTP framework), `pgx/v5` (PostgreSQL driver, no ORM), Viper (config), custom JWT access tokens + hashed refresh tokens, bcrypt password hashing, `google/uuid`. Architecture style: Clean Architecture / DDD.

**Frontend:** Vue 3 (Composition API, `<script setup>`), Vite 6, Vue Router 4, TypeScript (strict), Vitest + Vue Test Utils for testing. No state library (Vuex/Pinia) — plain module-level `ref`/`computed` singletons act as shared stores. No UI/CSS framework — hand-written global CSS with design tokens and a light/dark theme toggle.

**Database:** PostgreSQL 18 (Docker Compose for local dev). Schema managed via plain numbered SQL migrations (`db/migrations/*.up.sql`/`*.down.sql`) applied through `psql` via Makefile — no migration tool (e.g. golang-migrate) is actually wired in.

**Notable dependencies:**
| Library | Purpose |
|---|---|
| `labstack/echo` | HTTP server, routing, middleware |
| `jackc/pgx/v5` | Postgres driver/pool |
| `spf13/viper` | Config (YAML + env) |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `vue-router` | SPA routing + navigation guards |

No Axios/TanStack Query — a hand-rolled `fetch` wrapper is the sole API client.

## 3. Architecture & Directory Structure

```text
rxwod/
├── backend/                 Go API (Clean Architecture + DDD)
│   ├── cmd/api/main.go      Composition root
│   ├── internal/
│   │   ├── domain/          Pure business rules (wod/, gym/, user/, authz/) — no framework/DB deps
│   │   ├── application/     Use-case services, port interfaces, DTOs, commands
│   │   ├── delivery/http/   Echo handlers, middleware, request/response DTOs, error mapping
│   │   └── infrastructure/  Postgres repos, JWT issuer, bcrypt, Viper config
│   └── migrations/          Placeholder README (real SQL lives in top-level db/migrations)
├── db/migrations/           Schema of record (numbered up/down SQL)
├── frontend/src/
│   ├── app/                 App shell + router (auth/workspace guards)
│   ├── pages/               Thin route views
│   ├── features/            Vertical slices: auth/, wod/, workspace/ (api/, composables/, model/, components/)
│   ├── shared/               httpClient, Base* form inputs, AppHeader, Result<T> utility
│   └── styles/               Global design tokens
├── docker-compose.yml        Postgres only — no app containers
└── docs/                     Architecture/extension guides
```

Dependency rule: `delivery → application → domain`; `infrastructure` implements `application`-defined interfaces (ports & adapters). Domain code never imports Echo/pgx/JSON tags. Frontend uses a feature-folder pattern with thin page components.

## 4. Core Features (Implemented)

**Authentication & Accounts** — Email/password register & login; JWT access tokens + hashed, DB-stored refresh tokens; refresh endpoint; `GET /me` returns profile + gym memberships; session persisted in `localStorage`.

**Multi-Tenant Gym Workspaces** — Create a gym (creator becomes owner); belong to multiple gyms with different roles; three roles (owner/coach/athlete) with a defined permission matrix; owner-only member management (list, invite coach/athlete by email, change role, remove); workspace switcher auto-attaches `X-Gym-ID` to gym-scoped requests.

**WOD Authoring** — Multi-stage builder (WARMUP/STRENGTH/CORE/METCON/COOLDOWN stage kinds); five workout types per stage (AMRAP, FORTIME, TABATA, EMOM, OPEN free-text); per-movement label/name/prescription/sets/reps/load+unit/notes; derived scoring kind per stage; live score/format preview; lifecycle status DRAFT → PUBLISHED → ARCHIVED (archive not reachable via API/UI yet); explicit publish action; scheduling with a calendar view and date-based list filtering; create/list/get/update endpoints scoped to the active gym; role- and ownership-aware edit/view/publish UI gating.

**Frontend UX** — Client-side draft autosave per workspace/WOD in `localStorage` with "resume draft" and relative-time messaging; light/dark theme toggle; route guards for auth/workspace state; reusable base form components.

**Testing Infrastructure** — Unit tests across backend domain/application/delivery layers with hand-rolled fakes (no mocking framework); Postgres integration test gated on `DATABASE_URL`; frontend Vitest coverage for schemas, composables, and calendar utilities.

## 5. Data Flow & API Interaction

REST/JSON over HTTP, no GraphQL/WebSockets. Frontend `httpClient` (thin `fetch` wrapper) auto-attaches `Authorization: Bearer <accessToken>` and `X-Gym-ID: <activeWorkspaceId>` per-call as needed, normalizes all responses into a `Result<T>` type, and clears the session globally on any 401.

Backend middleware chain: logger → panic recovery → CORS → `AuthMiddleware` (validates JWT, attaches principal) → `GymContextMiddleware` (resolves role in the gym from `X-Gym-ID`) → `RequirePermission` (checks role against the `domain/authz` matrix).

**Key endpoints** (`/api/v1`): `POST /auth/register|login|refresh`, `GET /me`, `POST /gyms`, `GET /gyms`, `GET /gyms/:gymId`, `GET /gyms/:gymId/members`, `PATCH /gyms/:gymId/members/:userId`, `DELETE /gyms/:gymId/members/:userId`, `POST /gyms/:gymId/coaches|athletes`, `POST /wods`, `GET /wods`, `GET /wods/calendar`, `GET /wods/:id`, `PUT /wods/:id`, `POST /wods/:id/publish`.

> Note: `README.md` documents an older/smaller endpoint list (missing role update/removal, calendar, publish) — the router (`backend/internal/delivery/http/router.go`) is the source of truth.

Error contract: consistent `{"error": "<message>"}` with centralized status mapping (400 validation, 401 unauthenticated, 403 forbidden, 404 not found/gym mismatch, 500 fallback); internal/DB errors never leaked.

## 6. Missing Features & TODOs

No inline `TODO`/`FIXME`/`HACK` comments exist anywhere in the codebase. Gaps identified by comparing implementation against the domain model and typical product expectations:

- No athlete-facing results logging, PR tracking, leaderboards, or workout history — despite `ScoringKind` already being computed per stage.
- `ARCHIVED` WOD status exists in the domain but has no API route or UI action to reach it.
- No visible invite-acceptance flow for a newly invited user (invitations table has `pending/accepted/expired/revoked` status, but no claim/accept endpoint is wired) and no transactional email sending anywhere.
- No password reset / forgot-password flow.
- No email verification on registration.
- No pagination on `GET /wods` or `GET /gyms/:gymId/members`.
- No search/filter on WOD list beyond calendar date.
- `wod:delete` permission is defined in the authz matrix (owner-only) but there is no delete service method, repository method, or route behind it.
- Two migrations directories exist (`backend/migrations/` placeholder vs. `db/migrations/` actual) — potential onboarding confusion.
- `README.md`'s endpoint list is stale relative to the real router.

## 7. Production Readiness Checklist

**Security**
- Hardcoded default JWT secret (`dev-secret-change-me` in `backend/config.yaml`) — needs a mandatory, validated override per environment.
- `middleware.CORS()` used with no config — defaults to allowing all origins; must be locked to real frontend origin(s) per environment.
- No rate limiting anywhere, especially `auth/login|register|refresh` — brute-force/credential-stuffing exposure.
- No security headers middleware (HSTS/CSP/X-Frame-Options).
- Logout only clears client-side storage — verify server-side refresh token revocation is actually invoked (the `revoked_at` column exists but usage should be confirmed).
- `.env.example` is present but empty in this checkout — should list all required variable names.

**Reliability & Observability**
- No structured logging (only Echo's default request logger), no request IDs, no log aggregation integration.
- No health check endpoint for load balancer/orchestrator probes.
- No metrics/tracing (Prometheus/OpenTelemetry).
- No error tracking (e.g. Sentry) on backend or frontend.

**Data & Migrations**
- Migrations applied via a raw `psql` shell loop with no version-tracking tool — no protection against partial/out-of-order runs in CI/CD.
- No documented backup/restore strategy.
- Production Postgres pool tuning (size/timeouts) not yet addressed beyond local defaults.

**Deployment & Infra**
- No Dockerfiles for the backend or frontend — only Postgres is containerized.
- No CI/CD pipeline (no `.github/workflows`); `make test` isn't enforced automatically.
- No documented staging/production config story beyond local `config.yaml` + env overrides.

**Product Completeness** (see Section 6) — pagination/filtering, WOD delete/archive, invitation acceptance + transactional email, and a decision on athlete-facing engagement features before broader launch.

**Documentation** — reconcile `README.md`'s endpoint list with the real router; resolve the duplicate migrations-directory naming.