# LocalOffers API

Go REST API for LocalOffers. It provides authentication, nearby-offer discovery, offer publishing, comments, voting, and image uploads. Product-level documentation lives in the sibling coordination repository's [`docs/`](../docs/) directory.

## Stack

- Go 1.23 and chi
- PostgreSQL with PostGIS
- GORM
- JWT bearer authentication and bcrypt
- SQL migrations compatible with `golang-migrate`

## Local setup

```bash
docker-compose up -d

migrate -path migrations \
  -database "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" up

# Optional sample data
psql "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" \
  -f seed_offers.sql

go run ./cmd/api
```

The API reads `.env`. `JWT_SECRET` is required. The committed configuration is for local development only and must not be reused in a deployed environment.

## Checks

```bash
go test ./...
go vet ./...
go build ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The race detector requires CGO. Repository integration tests require an isolated PostgreSQL/PostGIS database whose name contains `test`:

```bash
docker compose up -d
docker exec local-offers-postgres createdb -U postgres local_offers_test

LOCAL_OFFERS_TEST_DATABASE_URL='postgres://postgres:postgres@localhost:5432/local_offers_test?sslmode=disable' \
  go test -tags=integration ./internal/repositories
```

The integration suite drops and rebuilds its tables from the real migrations. Its database-name guard prevents it from targeting the normal development database. GitHub Actions runs formatting, vet, race-enabled unit tests, core-package coverage, and PostGIS integration tests for pull requests and pushes to `main`.

## Error responses

API failures use a stable JSON envelope:

```json
{
  "error": {
    "code": "offer_not_found",
    "message": "The offer was not found."
  }
}
```

Request JSON is size-bounded and rejects malformed input, unknown fields, and trailing values. Successful response payloads remain unchanged.

## Structure

```text
cmd/api/               application entrypoint and route wiring
internal/auth/         JWT generation, parsing, and middleware
internal/config/       environment configuration
internal/database/     database connection
internal/handlers/     HTTP boundary
internal/httpx/        strict JSON and response helpers
internal/services/     business logic
internal/repositories/ persistence and PostGIS queries
internal/models/       database models
internal/dto/          public request and response shapes
internal/server/       router construction
migrations/            schema migrations
```

See [`../docs/architecture.md`](../docs/architecture.md) for the current API surface and cross-stack design, and [`../docs/development.md`](../docs/development.md) for the complete local workflow.
