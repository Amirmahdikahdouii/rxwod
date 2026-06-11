.PHONY: migrate-up migrate-down migrate-up-docker migrate-down-docker backend-test frontend-test test

DATABASE_URL ?= postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable

migrate-up:
	psql "$(DATABASE_URL)" -f db/migrations/000001_create_wods.up.sql

migrate-down:
	psql "$(DATABASE_URL)" -f db/migrations/000001_create_wods.down.sql

migrate-up-docker:
	docker compose exec -T postgres psql -U rxwod -d rxwod -f /dev/stdin < db/migrations/000001_create_wods.up.sql

migrate-down-docker:
	docker compose exec -T postgres psql -U rxwod -d rxwod -f /dev/stdin < db/migrations/000001_create_wods.down.sql

backend-test:
	cd backend && go test ./...

frontend-test:
	cd frontend && npm run test:unit

test: backend-test frontend-test
