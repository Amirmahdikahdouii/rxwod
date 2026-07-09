# RXWOD — Phase 1 WOD Creator

Full-stack CrossFit WOD creator built with Go (Echo, PostgreSQL, Viper) and Vue 3 (Composition API, Vite).

## Architecture

```text
frontend/   Vue 3 UI
backend/    Go API (Clean Architecture + DDD)
db/         PostgreSQL migrations
```

## Prerequisites

- Go 1.22+
- Node.js 20+
- Docker (for local PostgreSQL)

## Quick Start

1. Start PostgreSQL:

```bash
docker compose up -d postgres
```

1. Apply migrations:

```bash
make migrate-up
```

1. Run the API:

```bash
export AUTH_JWT_SECRET=your-local-dev-secret
cd backend && go run ./cmd/api
```

1. Run the frontend:

```bash
cd frontend
npm install
npm run dev
```

1. Open `http://localhost:5173`

Backend defaults live in [`backend/config.yaml`](backend/config.yaml). Environment
variables from [`.env.example`](.env.example) can override matching YAML values
when exported in the shell before starting the API. `AUTH_JWT_SECRET` is required
and has no default — the API will not start without it.

## API Endpoints

A WOD is a multi-stage program: one create request builds an ordered list of stages, each
with a stage kind (`WARMUP`, `STRENGTH`, `CORE`, `METCON`, `COOLDOWN`), its own workout type
(`AMRAP`, `FORTIME`, `TABATA`, `EMOM`) config, and its own movements.

Backend APIs now use email/password authentication and gym workspaces. Authenticated gym-scoped
requests must send both `Authorization: Bearer <accessToken>` and `X-Gym-ID: <gymId>`.

Auth:

- `POST /api/v1/auth/register` — create a user and return access/refresh tokens
- `POST /api/v1/auth/login` — authenticate and return access/refresh tokens
- `POST /api/v1/auth/refresh` — exchange a refresh token for a new access token
- `GET /api/v1/me` — current user plus active gym memberships

Gyms:

- `POST /api/v1/gyms` — create a gym; the caller becomes owner
- `GET /api/v1/gyms` — list gyms available to the caller
- `GET /api/v1/gyms/:gymId` — read gym details
- `GET /api/v1/gyms/:gymId/members` — owner-only member list
- `POST /api/v1/gyms/:gymId/coaches` — owner-only coach invite/add by email
- `POST /api/v1/gyms/:gymId/athletes` — owner-only athlete invite/add by email

WODs:

- `POST /api/v1/wods` — owner/coach create a multi-stage WOD program in the active gym
- `GET /api/v1/wods` — owner/coach/athlete list WOD programs in the active gym
- `GET /api/v1/wods/:id` — owner/coach/athlete fetch a program in the active gym
- `PUT /api/v1/wods/:id` — owner-only update in the active gym

Example create payload:

```json
{
  "name": "Monday Session",
  "description": "Full class plan",
  "stages": [
    { "kind": "WARMUP",   "type": "FORTIME", "config": { "rounds": 2 },           "movements": [{ "position": 1, "name": "Jumping Jacks", "reps": 20 }] },
    { "kind": "METCON",   "type": "AMRAP",   "config": { "timeCapSeconds": 900 }, "movements": [{ "position": 1, "name": "Burpee", "reps": 21 }] },
    { "kind": "COOLDOWN", "type": "TABATA",  "config": { "workSeconds": 20, "restSeconds": 10, "rounds": 8, "cycles": 1 }, "movements": [{ "position": 1, "name": "Plank", "reps": 1 }] }
  ]
}
```

## Quality Gates

```bash
make test
cd backend && go test ./...
cd frontend && npm run typecheck && npm run test:unit && npm run build
```

## Extending WOD Types

See [docs/EXTENDING_WOD_TYPES.md](docs/EXTENDING_WOD_TYPES.md).
