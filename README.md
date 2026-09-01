# LocalOffers — API

Backend for LocalOffers, a hyperlocal, collaborative board of local offers: users see nearby offers, publish ones they find, and validate/comment on offers others post. The mobile client lives in the sibling repo `local-offers-app`.

## Status

Early WIP. A single vertical slice is implemented: creating an offer and finding offers near a point, backed by PostGIS geospatial queries. Users, auth, comments, and validate/invalidate voting are not modeled yet — see [CLAUDE.md](CLAUDE.md) for the current architecture and known gaps.

## Tech stack

- Go 1.23, [chi](https://github.com/go-chi/chi) router
- PostgreSQL + PostGIS, via GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`)
- SQL migrations in the [golang-migrate](https://github.com/golang-migrate/migrate) file format

## Getting started

Requires Docker and the [`migrate` CLI](https://github.com/golang-migrate/migrate#cli-usage).

```bash
# 1. Start Postgres+PostGIS
docker-compose up -d

# 2. Apply migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" up

# 3. (optional) seed sample offers
psql "postgres://postgres:postgres@localhost:5432/local_offers?sslmode=disable" -f seed_offers.sql

# 4. Run the API
go run ./cmd/api
```

Config is read from `.env` (see the committed `.env` for local defaults — `APP_PORT`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`).

The frontend (`local-offers-app`) expects the API reachable at your machine's LAN IP, not `localhost`, when testing on a device via Expo Go.

## API

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Healthcheck |
| POST | `/offers` | Create an offer |
| GET | `/offers/{id}` | Fetch an offer by id |
| GET | `/offers/nearby?lat=&lng=` | Offers within 5km of a point, nearest first |

No authentication, update/delete, or comment/vote endpoints exist yet.

## Project structure

```
cmd/api/            # entrypoint, wiring, routes
internal/config/     # env config loading
internal/database/   # gorm/Postgres connection
internal/models/      # gorm structs (DB row shape)
internal/dto/          # request/response JSON shapes
internal/repositories/ # DB access, incl. raw PostGIS queries
internal/services/     # business logic
internal/handlers/     # HTTP handlers
migrations/           # golang-migrate SQL migrations
```

## Docs

- [CLAUDE.md](CLAUDE.md) — architecture, data model, and conventions for working in this repo (also read by Claude Code).
