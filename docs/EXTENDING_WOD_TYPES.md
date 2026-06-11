# Extending WOD Types

This guide explains how to add a new WOD type (for example `CHIPPER`) without rewriting the core domain.

## 1. Domain Layer

Add to [`backend/internal/domain/wod/types.go`](../backend/internal/domain/wod/types.go):

- A new `WODType` constant
- Any new value object types (avoid raw primitives)

Add a config struct in [`backend/internal/domain/wod/config.go`](../backend/internal/domain/wod/config.go):

- Implement the sealed `Config` interface
- Validate invariants in `Validate()`
- Derive `ScoringKind()` from domain rules

Add a saved variant wrapper in [`backend/internal/domain/wod/variant.go`](../backend/internal/domain/wod/variant.go):

- `type ChipperWOD = WOD[ChipperConfig]`
- `type SavedChipper struct{ wod ChipperWOD }`
- Implement `Variant` accessors

Add domain unit tests in [`backend/internal/domain/wod/aggregate_test.go`](../backend/internal/domain/wod/aggregate_test.go).

## 2. Application Layer

Add a typed command in [`backend/internal/application/wod/commands.go`](../backend/internal/application/wod/commands.go).

Extend:

- `CreateCommand`
- `Service.buildVariant()`
- `configFromVariant()`

Add application tests in [`backend/internal/application/wod/service_test.go`](../backend/internal/application/wod/service_test.go).

## 3. Infrastructure Layer

Update [`backend/internal/infrastructure/postgres/wod_mapper.go`](../backend/internal/infrastructure/postgres/wod_mapper.go):

- Add JSON payload struct
- Add `chipperToRecord()` and `chipperFromRecord()`
- Extend `variantToRecord()` and `recordToVariant()` switches

If the database enum must include the new type, add a migration under [`db/migrations/`](../db/migrations/).

## 4. Delivery Layer

Update [`backend/internal/delivery/http/wod_handler.go`](../backend/internal/delivery/http/wod_handler.go):

- Extend `toCreateCommand()` with validation for the new type

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
