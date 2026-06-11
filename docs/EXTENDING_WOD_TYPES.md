# Extending WOD Types

A WOD is a multi-stage program: an aggregate root holding ordered stages, each with a
stage kind (`WARMUP`, `STRENGTH`, `CORE`, `METCON`, `COOLDOWN`) and its own workout-type
config. This guide explains how to add a new workout type (for example `CHIPPER`) that
stages can use, without rewriting the core domain.

## 1. Domain Layer

Add to [`backend/internal/domain/wod/types.go`](../backend/internal/domain/wod/types.go):

- A new `WODType` constant
- Any new value object types (avoid raw primitives)

Add a config struct in [`backend/internal/domain/wod/config.go`](../backend/internal/domain/wod/config.go):

- Implement the sealed `Config` interface (`wodConfig()`, `Type()`, `Validate()`, `ScoringKind()`)
- Validate invariants in `Validate()`
- Derive `ScoringKind()` from domain rules

Stages hold a `Config` value directly via [`backend/internal/domain/wod/stage.go`](../backend/internal/domain/wod/stage.go),
so there is no per-type aggregate wrapper to add — the sealed `Config` interface is the only
polymorphism mechanism.

Add domain unit tests in [`backend/internal/domain/wod/aggregate_test.go`](../backend/internal/domain/wod/aggregate_test.go).

## 2. Application Layer

The create contract is a single `CreateWODCommand` with a `Stages []StageInput` slice; each
stage carries a `StageConfigInput` discriminated by `Type`. To support a new type, extend the
`buildConfig()` switch in [`backend/internal/application/wod/service.go`](../backend/internal/application/wod/service.go)
to map the new `StageConfigInput` fields to your config constructor, and extend `configToDTO()`
for read responses. Add any new optional fields to `StageConfigInput` in
[`backend/internal/application/wod/commands.go`](../backend/internal/application/wod/commands.go)
and `ConfigDTO` in [`backend/internal/application/wod/dto.go`](../backend/internal/application/wod/dto.go).

Add application tests in [`backend/internal/application/wod/service_test.go`](../backend/internal/application/wod/service_test.go).

## 3. Infrastructure Layer

Update [`backend/internal/infrastructure/postgres/wod_mapper.go`](../backend/internal/infrastructure/postgres/wod_mapper.go):

- Add a JSON payload struct
- Extend the `configToJSON()` and `configFromJSON()` switches over the sealed `Config` type

If the database enum must include the new type, add a migration under [`db/migrations/`](../db/migrations/).

## 4. Delivery Layer

The HTTP handler passes stage `type`/`kind` strings through as-is for domain validation, so no
per-type branching is needed in [`backend/internal/delivery/http/wod_handler.go`](../backend/internal/delivery/http/wod_handler.go).
Add any new config fields to `StageConfigRequest`/`ConfigResponse` in
[`request.go`](../backend/internal/delivery/http/request.go) and
[`response.go`](../backend/internal/delivery/http/response.go).

Add handler tests in [`backend/internal/delivery/http/wod_handler_test.go`](../backend/internal/delivery/http/wod_handler_test.go).

## 5. Frontend

Update:

- [`frontend/src/features/wod/model/wodTypes.ts`](../frontend/src/features/wod/model/wodTypes.ts)
- [`frontend/src/features/wod/model/wodSchemas.ts`](../frontend/src/features/wod/model/wodSchemas.ts)
- [`frontend/src/features/wod/components/ScorePreview.vue`](../frontend/src/features/wod/components/ScorePreview.vue)

Add schema tests in [`frontend/src/features/wod/model/wodSchemas.test.ts`](../frontend/src/features/wod/model/wodSchemas.test.ts).

## Design Rules

- Domain models must never import Echo, PostgreSQL, JSON tags, or Viper.
- Use typed commands instead of `map[string]any`.
- Keep type-specific persistence in JSONB mappers, not in the domain.
- Prefer sealed interfaces and exhaustive switches over open-ended dynamic typing.
