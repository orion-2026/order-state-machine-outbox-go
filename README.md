# Order State Machine + Outbox Demo (Go + Gin)

A Go + Gin implementation of the order status flow architecture lesson.

## What this demo shows
- Order service as source of truth
- Explicit state machine for allowed transitions
- Application/service layer centralizing business rules
- Outbox pattern for domain event recording
- Easy local run with SQLite
- Docker Compose option with PostgreSQL

## Order lifecycle
```text
PendingPayment -> Paid -> Fulfilling -> Shipped -> Completed
PendingPayment -> Cancelled
Paid -> Refunded
Fulfilling -> Cancelled
Shipped -> Refunded
```

## API surface
- `GET /health`
- `POST /api/orders`
- `GET /api/orders`
- `GET /api/orders/{id}`
- `GET /api/orders/{id}/actions`
- `POST /api/orders/{id}/transitions`
- `GET /api/outbox`
- `POST /api/outbox/publish`

## Environment variables
- `PORT=8080`
- `DATABASE_PROVIDER=sqlite | postgres`
- `DATABASE_CONNECTION_STRING=...`

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

Default URLs:
- API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

## Example flow
### Create order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"customerId":"cust-001","productSku":"sku-demo-001","quantity":2}'
```

### Check allowed actions
```bash
curl http://localhost:8080/api/orders/{orderId}/actions
```

### Transition state
```bash
curl -X POST http://localhost:8080/api/orders/{orderId}/transitions \
  -H "Content-Type: application/json" \
  -d '{"action":"pay","reason":"Payment callback received"}'
```

### Read outbox
```bash
curl http://localhost:8080/api/outbox
```

## Notes
- Uses GORM auto-migrate for demo simplicity
- Stores outbox payloads as JSON strings
- Keeps the project intentionally small and teaching-oriented

## Good next upgrades
- Replace auto-migrate with real migrations
- Add background publisher worker
- Add optimistic concurrency checks
- Add dedup / inbox handling
