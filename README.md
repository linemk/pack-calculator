# Pack Calculator

HTTP API and UI for calculating the optimal number of packs to fulfil an order.

**Rules (in priority order):**
1. Only whole packs can be sent
2. Send the minimum number of items (≥ order amount)
3. Among solutions with equal item count, minimise the number of packs

## Algorithm

Dynamic programming over `dp[0..order+maxPack]` where `dp[i]` = minimum packs to ship exactly `i` items. The first reachable target `≥ order` gives the optimal solution.

This correctly handles edge cases where a greedy approach fails, e.g.:
- sizes `[23, 31, 53]`, order `500000` → `{53: 9429, 31: 7, 23: 2}` (exactly 500 000 items)

## Running locally

**Prerequisites:** Go 1.22+

```bash
make run
# or
go run ./cmd/server
```

Open http://localhost:8080

## Running with Docker

```bash
# Build
docker build -t pack-calculator .

# Run
docker run -p 8080:8080 pack-calculator
```

Open http://localhost:8080

## Running with docker-compose

```bash
docker-compose up --build
```

## Running tests

```bash
make test
# or
go test ./... -race -count=1 -v
```

## Linting

```bash
make lint
# or
golangci-lint run ./...
```

Full CI pipeline (lint + test + build):

```bash
make ci
```

## API

### GET /api/packs

Returns current pack sizes.

```bash
curl http://localhost:8080/api/packs
# {"sizes":[250,500,1000,2000,5000]}
```

### PUT /api/packs

Update pack sizes (persists until restart).

```bash
curl -X PUT http://localhost:8080/api/packs \
  -H "Content-Type: application/json" \
  -d '{"sizes":[23,31,53]}'
# {"sizes":[53,31,23]}
```

### POST /api/calculate

Calculate packs for an order.

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order":500000}'
# {"packs":{"23":2,"31":7,"53":9429},"total_items":500000}
```

### GET /health

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Project structure

```
cmd/server/         entry point, DI wiring (uber/dig), graceful shutdown
  web/              embedded UI (index.html)
internal/
  calculator/       DP algorithm + tests
  store/            thread-safe in-memory pack size store
  handler/          HTTP handlers (gorilla/mux) + tests
  middleware/       RequestID, Logger, Recoverer, CORS
```
