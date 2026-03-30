# Order State Machine + Outbox Demo (Go)

A Go + Gin implementation of the order status flow architecture lesson.

## What this demo shows
- Order service as source of truth
- Explicit state machine for allowed transitions
- Service layer centralizing business rules
- Outbox pattern for domain event recording
- SQLite for easy local run
- PostgreSQL via Docker Compose for a more realistic setup

## Endpoints
- `GET /health`
- `POST /api/orders`
- `GET /api/orders`
- `GET /api/orders/{id}`
- `GET /api/orders/{id}/actions`
- `POST /api/orders/{id}/transitions`
- `GET /api/outbox`
- `POST /api/outbox/publish`

## Run locally
```bash
cp .env.example .env 2>/dev/null || true
go mod tidy
go run ./cmd/api
```

Default local DB:
- `order-demo.db`

## Run with Docker + PostgreSQL
```bash
docker compose up --build
```

## Notes
- Uses GORM auto-migrate for demo simplicity
- Keeps payloads as JSON strings in the outbox table
- Intentionally small and teaching-oriented, not production-complete

## Good next upgrades
- Replace auto-migrate with real migrations
- Add background publisher worker
- Add optimistic concurrency checks
- Add dedup / inbox handling for callbacks
