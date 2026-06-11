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

- `POST /api/v1/wods` — create a WOD
- `GET /api/v1/wods` — list WODs
- `GET /api/v1/wods/:id` — fetch WOD detail

## Quality Gates

```bash
make test
cd backend && go test ./...
cd frontend && npm run typecheck && npm run test:unit && npm run build
```

## Extending WOD Types

See [docs/EXTENDING_WOD_TYPES.md](docs/EXTENDING_WOD_TYPES.md).
