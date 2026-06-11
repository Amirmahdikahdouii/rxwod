# Database Migrations

SQL migrations live in the repository root at `db/migrations/`.

Apply migrations locally (requires `psql` on the host):

```bash
docker compose up -d postgres
make migrate-up
```

Or via Docker (no local `psql` needed):

```bash
docker compose up -d postgres
make migrate-up-docker
```

Equivalent manual command:

```bash
docker compose exec -T postgres psql -U rxwod -d rxwod -f /dev/stdin < db/migrations/000001_create_wods.up.sql
```
