# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

This is the backend (Go) half of LocalOffers, a hyperlocal collaborative offers board. The frontend lives in the sibling repo `local-offers-app` (Expo/React Native).

## Commands

```bash
# Start Postgres+PostGIS locally
docker-compose up -d

# Run the API (reads .env for config)
go run ./cmd/api

# Apply/rollback migrations (golang-migrate CLI, not vendored — install separately)
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" up
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" down

# Seed sample data (after migrations)
psql "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" -f seed_offers.sql

go build ./...
go vet ./...
```

`.env` must set `JWT_SECRET` or `config.Load()` exits immediately.

No tests exist yet (`go test ./...` has nothing to run). Run `go mod tidy` after adding new imports to keep direct/indirect markings correct. `gorm.io/driver/mysql`/`go-sql-driver/mysql` will keep showing up as indirect deps — that's `gorm.io/datatypes` pulling them in transitively (`go mod why gorm.io/driver/mysql` to confirm), not something removable while `datatypes` (used for the `offer_type` JSONB field) stays in use.

## Architecture

Layered, manually-wired (no DI framework, no interfaces):

```
cmd/api/main.go       → config.Load() → database.Connect() → wires repo→service→handler → chi router
internal/config       → env vars via godotenv (.env)
internal/database     → gorm.Open(postgres) connection
internal/models        → gorm structs (DB row shape)
internal/dto           → request/response JSON shapes, decoupled from models
internal/repositories  → DB access; PostGIS queries use db.Raw(...).Scan(...), not GORM's query builder
internal/services      → business logic; DTO↔model translation
internal/handlers      → HTTP layer: decode request, call service, write JSON response
```

Three domains exist: **Offer**, **User/Auth**, and **Comment** (`internal/{models,dto,repositories,services,handlers}/{offer,auth,user,comment}*.go`, plus the standalone `internal/auth` package for JWT). When adding a new domain (e.g. votes), mirror this same repo→service→handler split and register routes under the `/api/v1` group in `main.go`.

`main.go` wires chi middleware in this order: `middleware.Logger` (request logging), `middleware.Recoverer` (panic recovery → 500 instead of a crash), then `cors.Handler` (wide open `AllowedOrigins: []string{"*"}` — fine for now since auth is Bearer-token, not cookie-based). `/health` stays unversioned at the router root; everything else lives under `/api/v1`. `POST /api/v1/offers` is wrapped in its own `r.Group` with `auth.RequireAuth(cfg.JWTSecret)` — that's the only protected route so far; browsing offers stays public per the product's intent.

### Auth

`internal/auth` holds JWT generation/parsing (`GenerateToken`/`ParseToken`, HS256, 30-day expiry, subject = user UUID) and the `RequireAuth` chi middleware, which puts the authenticated user id in the request context (read it back with `auth.UserIDFromContext`). `internal/services/auth_service.go` handles register/login: bcrypt-hashes passwords (min 8 chars, enforced in `Register`), looks up by lowercased email, returns `{token, user}`. Duplicate-email registration relies on the `users.email` unique constraint plus `gorm.Config{TranslateError: true}` (set in `internal/database/database.go`) so `errors.Is(err, gorm.ErrDuplicatedKey)` works. `JWT_SECRET` is required — `config.Load()` calls `log.Fatal` if it's unset; the committed `.env` has an obviously-dev value, replace it for any non-local environment.

`POST /api/v1/offers` now requires auth and reads the creator from the JWT (not the request body, to prevent spoofing) — see `internal/handlers/offer_handler.go`'s `Create`. `internal/repositories/offer_repository.go` LEFT JOINs `users` to populate the real `PostedBy` on every offer response, falling back to `{name: "Usuario LocalOffers"}` for the seeded rows that predate the `offers.user_id` column (nullable, added in migration `000004`).

### Data model

`offers` table (see `migrations/`): `id` (UUID), `headline`, `description`, `business_name`, `category`, `image_url`, `offer_type` (JSONB — polymorphic shape, see below), `location` (`GEOGRAPHY(POINT,4326)`, GIST-indexed), `user_id` (UUID, nullable FK → `users.id`), `expires_at`, `confirmations_count`, `invalidations_count`, `created_at`.

`users` table: `id` (UUID), `name`, `email` (unique), `password_hash`, `created_at`. `comments` table (migration `000005`): `id` (UUID), `offer_id` (FK → `offers.id`), `user_id` (FK → `users.id`), `body`, `created_at`. `confirmations_count`/`invalidations_count` are still plain counters with no per-user vote tracking (nothing can move them yet).

`offer_type` JSONB shapes (no DB-level schema/CHECK constraint — validate in Go): `{"type":"percentage","percentage":N}`, `{"type":"discount","percentage":N}`, `{"type":"bundle","label":"..."}`, `{"type":"text","label":"..."}`.

### Geospatial queries (PostGIS)

- Offer creation builds a WKT point (`SRID=4326;POINT(lon lat)`) in `internal/services/offer_service.go` and lets GORM insert it as the `location` string field.
- `GET /api/v1/offers/nearby?lat=&lng=&radius=` runs a raw SQL query in `internal/repositories/offer_repository.go` using `ST_DWithin(...)` filtered and `ST_Distance(...)` for `distance_meters`/ordering. `radius` (meters) is optional, parsed/clamped in `internal/handlers/offer_handler.go` (default 5000, max 20000, must be a positive integer or the request 400s).
- `GET /api/v1/offers/{id}` extracts lat/lng back out via `ST_Y`/`ST_X` on `location::geometry`.

### Current API surface

`GET /health` (unversioned), `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/offers` (auth required), `GET /api/v1/offers/{id}`, `GET /api/v1/offers/nearby?lat=&lng=&radius=`, `GET /api/v1/offers/{id}/comments`, `POST /api/v1/offers/{id}/comments` (auth required) — no update/delete, no listing-all, no vote endpoints yet. `PostedBy` and `CommentsCount` are now real (see Auth above and `internal/repositories/offer_repository.go`'s `comments_count` subquery); `IsVerifiedBusiness` is still always `false` — a placeholder until verified-business is modeled.
