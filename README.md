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

2. Apply migrations:

```bash
make migrate-up
```

3. Run the API:

```bash
cd backend && go run ./cmd/api
```

4. Run the frontend:

```bash
cd frontend
npm install
npm run dev
```

5. Open `http://localhost:5173`

Copy [`.env.example`](.env.example) and adjust values as needed.

## API Endpoints

A WOD is a multi-stage program: one create request builds an ordered list of stages, each
with a stage kind (`WARMUP`, `STRENGTH`, `CORE`, `METCON`, `COOLDOWN`), its own workout type
(`AMRAP`, `FORTIME`, `TABATA`, `EMOM`) config, and its own movements.

- `POST /api/v1/wods` — create a multi-stage WOD program
- `GET /api/v1/wods` — list WOD programs (with stage summaries)
- `GET /api/v1/wods/:id` — fetch a program with full nested stages

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
