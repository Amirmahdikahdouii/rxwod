.PHONY: migrate-up migrate-down migrate-up-docker migrate-down-docker backend-test frontend-test test

DATABASE_URL ?= postgres://rxwod:rxwod@localhost:5432/rxwod?sslmode=disable

migrate-up:
	@for f in $$(ls db/migrations/*.up.sql | sort); do \
		echo "applying $$f"; \
		psql "$(DATABASE_URL)" -f $$f || exit 1; \
	done

migrate-down:
	@for f in $$(ls db/migrations/*.down.sql | sort -r); do \
		echo "reverting $$f"; \
		psql "$(DATABASE_URL)" -f $$f || exit 1; \
	done

migrate-up-docker:
	@for f in $$(ls db/migrations/*.up.sql | sort); do \
		echo "applying $$f"; \
		docker compose exec -T postgres psql -U rxwod -d rxwod -f /dev/stdin < $$f || exit 1; \
	done

migrate-down-docker:
	@for f in $$(ls db/migrations/*.down.sql | sort -r); do \
		echo "reverting $$f"; \
		docker compose exec -T postgres psql -U rxwod -d rxwod -f /dev/stdin < $$f || exit 1; \
	done

backend-test:
	cd backend && go test ./...

frontend-test:
	cd frontend && npm run test:unit

test: backend-test frontend-test

run:
	docker compose up -d postgres
	cd backend && go run ./cmd/api
	cd frontend && npm run dev