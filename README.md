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
```

There are currently no automated test files, so `go test` primarily verifies compilation.

## Structure

```text
cmd/api/               application entrypoint and route wiring
internal/auth/         JWT generation, parsing, and middleware
internal/config/       environment configuration
internal/database/     database connection
internal/handlers/     HTTP boundary
internal/services/     business logic
internal/repositories/ persistence and PostGIS queries
internal/models/       database models
internal/dto/          public request and response shapes
migrations/            schema migrations
```

See [`../docs/architecture.md`](../docs/architecture.md) for the current API surface and cross-stack design, and [`../docs/development.md`](../docs/development.md) for the complete local workflow.
