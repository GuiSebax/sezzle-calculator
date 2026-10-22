# Sezzle Calculator

A full-stack calculator: a **Go** microservice that performs arithmetic over a small REST API, and a **React + TypeScript** single-page app that calls it.

Supported operations: addition, subtraction, multiplication, division, exponentiation, square root and percentage.

## Contents

- [Tech stack](#tech-stack)
- [Project structure](#project-structure)
- [Running the project](#running-the-project)
- [Tests and coverage](#tests-and-coverage)
- [API](#api)
- [Design decisions and trade-offs](#design-decisions-and-trade-offs)
- [Assumptions](#assumptions)
- [Possible improvements](#possible-improvements)
- [Use of AI tooling](#use-of-ai-tooling)

## Tech stack

| Layer    | Technology                                                                                 |
| -------- | ------------------------------------------------------------------------------------------ |
| Backend  | Go 1.22+ (standard library only: `net/http`, `encoding/json`, `log/slog`)                  |
| Frontend | React 19, TypeScript (strict), Vite, CSS Modules — no UI library                           |
| Testing  | Go `testing` + `net/http/httptest`; Vitest + React Testing Library                          |
| Delivery | Docker multi-stage builds (distroless for the API, unprivileged nginx for the UI), Compose |

## Project structure

```
.
├── backend/
│   ├── cmd/server/main.go        # entry point: config from env, HTTP server, graceful shutdown
│   ├── internal/calculator/      # pure math functions + sentinel errors (no HTTP)
│   ├── internal/api/             # routing, JSON handling, error mapping, CORS middleware
│   ├── Dockerfile
│   └── Makefile                  # run, test, coverage
├── frontend/
│   ├── src/api/                  # typed API client and error types (only place that calls fetch)
│   ├── src/hooks/                # useCalculation: request lifecycle state
│   ├── src/lib/                  # input validation and display formatting
│   ├── src/components/           # Calculator, NumberField, ResultPanel (+ CSS modules)
│   ├── nginx.conf                # serves the SPA and proxies /api to the backend
│   └── Dockerfile
├── docker-compose.yml
├── Makefile                      # shortcuts for both layers
└── PROMPTS.md                    # AI prompts used during development
```

## Running the project

### With Docker Compose (recommended)

Requires Docker with Compose v2.

```bash
docker compose up --build
```

| Service  | URL                            |
| -------- | ------------------------------ |
| Frontend | http://localhost:3000          |
| Backend  | http://localhost:8080 (direct) |

Host ports can be changed if they are taken: `FRONTEND_PORT=3001 BACKEND_PORT=8081 docker compose up --build`.

### Locally, without Docker

Requires Go 1.22+ and Node.js 22+. Start the backend first, then the frontend in a second terminal.

**Backend**

```bash
cd backend
go run ./cmd/server          # or: make run
```

| Variable               | Default                 | Description                                 |
| ---------------------- | ----------------------- | ------------------------------------------- |
| `PORT`                 | `8080`                  | Port the API listens on                     |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Comma-separated list of allowed origins     |

**Frontend**

```bash
cd frontend
npm install
npm run dev                  # http://localhost:5173
```

| Variable       | Default                 | Description                                              |
| -------------- | ----------------------- | -------------------------------------------------------- |
| `VITE_API_URL` | `http://localhost:8080` | Backend base URL, read at build time (see `.env.example`) |

Other scripts: `npm run build` (type-check + production build) and `npm run preview`.

## Tests and coverage

| Command                         | What it does                                                     |
| ------------------------------- | ---------------------------------------------------------------- |
| `cd backend && make test`       | Go unit and handler tests (`go test ./...`)                      |
| `cd backend && make coverage`   | Writes `backend/coverage/coverage.out` and `coverage/index.html` |
| `cd frontend && npm test`       | Vitest suite                                                     |
| `cd frontend && npm run coverage` | Text summary + HTML report in `frontend/coverage/index.html`   |
| `make test` / `make coverage`   | Runs both layers from the repository root                        |

Without `make` (e.g. on Windows), run the commands directly: `go test ./... -cover`, or `go test -coverprofile=coverage.out ./...` followed by `go tool cover -html=coverage.out`.

Coverage reports are git-ignored.

### Test coverage

| Layer    | Scope                          | Statements |
| -------- | ------------------------------ | ---------- |
| Backend  | `internal/calculator`          | 100.0%     |
| Backend  | `internal/api`                 | 96.8%      |
| Frontend | All source files (45 tests)    | ~98%       |

`cmd/server` (process wiring only: env vars, signals, server start/stop) has no unit tests, which brings the backend total to 71.8%. Its graceful shutdown was checked manually by sending SIGTERM to the container.

## API

Base URL: `http://localhost:8080`. All responses are JSON.

### `GET /health`

```bash
curl http://localhost:8080/health
```

```json
{ "status": "ok" }
```

### `POST /api/v1/calculate/{operation}`

| Operation    | Computes     | Required fields |
| ------------ | ------------ | --------------- |
| `add`        | `a + b`      | `a`, `b`        |
| `subtract`   | `a - b`      | `a`, `b`        |
| `multiply`   | `a × b`      | `a`, `b`        |
| `divide`     | `a ÷ b`      | `a`, `b`        |
| `power`      | `a ^ b`      | `a`, `b`        |
| `percentage` | `a`% of `b`  | `a`, `b`        |
| `sqrt`       | `√a`         | `a`             |

**Success** — `200 OK`

```bash
curl -X POST http://localhost:8080/api/v1/calculate/add \
  -H "Content-Type: application/json" \
  -d '{"a": 5, "b": 3}'
```

```json
{ "operation": "add", "result": 8 }
```

```bash
curl -X POST http://localhost:8080/api/v1/calculate/sqrt -d '{"a": 16}'
# {"operation":"sqrt","result":4}

curl -X POST http://localhost:8080/api/v1/calculate/percentage -d '{"a": 10, "b": 50}'
# {"operation":"percentage","result":5}
```

**Errors** — every error uses the same shape:

```json
{ "error": "a human-readable message" }
```

```bash
curl -X POST http://localhost:8080/api/v1/calculate/divide -d '{"a": 5, "b": 0}'
# 422 {"error":"division by zero"}

curl -X POST http://localhost:8080/api/v1/calculate/add -d '{"a": 5}'
# 400 {"error":"field \"b\" is required"}
```

### Status codes

| Status | Meaning                  | Examples                                                                                  |
| ------ | ------------------------ | ----------------------------------------------------------------------------------------- |
| 200    | Success                  | —                                                                                         |
| 400    | The request is invalid   | malformed JSON, empty body, missing field, non-numeric value, unknown operation           |
| 413    | Body too large           | request body over 1 KiB                                                                   |
| 422    | Valid request, no result | `division by zero`, `square root of a negative number`, `result is not a finite number` (overflow, `0 ^ -1`, `(-8) ^ (1/3)`) |
| 405    | Wrong HTTP method        | `GET /api/v1/calculate/add` (generated by `ServeMux`, plain-text body)                    |

## Design decisions and trade-offs

**Pure calculator package, separate from HTTP.** `internal/calculator` only knows about numbers: each function returns `(float64, error)` and reports failures with sentinel errors (`ErrDivisionByZero`, `ErrNegativeSqrt`, `ErrNonFinite`). The HTTP layer maps them to status codes with `errors.Is`. The math can be tested exhaustively without HTTP, and the HTTP layer could be swapped (gRPC, CLI) without touching it. NaN and ±Inf are turned into errors because they are not valid JSON and are never a useful answer.

**One endpoint per operation.** `POST /api/v1/calculate/{operation}` uses Go 1.22 `ServeMux` patterns, so routing, method matching and path parameters need no framework. Internally, a map from operation name to function keeps adding an operation to a one-line change. The `/v1` prefix leaves room for breaking changes.

**400 vs 422.** 400 means the client sent something the server cannot interpret (bad JSON, missing fields, unknown operation). 422 means the request was well-formed but the math has no answer. The frontend uses this to label errors as "Invalid input" or "Math error". Any unmapped error falls through to 500.

**Defensive input handling.** Request fields are `*float64`, so a missing `b` is reported as missing instead of silently becoming `0` (which would turn into a misleading "division by zero"). Bodies are capped with `http.MaxBytesReader`, and trailing data after the JSON object is rejected.

**Standard library only.** The service is small enough that `net/http`, `encoding/json` and `log/slog` cover everything. That means no dependencies to audit or update, and the code is familiar to any Go developer. The trade-off is writing small pieces like CORS by hand (about 20 lines).

**Production-minded server.** The server sets explicit read, write and idle timeouts (Go's defaults have none) and shuts down gracefully on SIGINT/SIGTERM. All configuration comes from environment variables.

**Typed API client in the frontend.** `src/api` is the only code that calls `fetch`. It turns responses into either a typed result or one of two error classes: `NetworkError` (backend unreachable) and `ApiError` (with `status` and a `kind` of `invalid-input`, `math` or `server`). Components handle errors without knowing anything about HTTP. Response bodies are treated as untrusted and checked before use.

**Client-side validation, backend as source of truth.** The UI validates for fast feedback: required fields and a strict number format (it rejects `0x10`, `Infinity` and `1,5`, which `Number()` would partly accept). It shows inline messages and disables Calculate while the form is invalid. It does not repeat the math rules: division by zero or a negative square root is sent to the backend, and its message is shown. That way there is a single implementation of the rules.

**UI state as a discriminated union.** `idle | loading | success | error` makes impossible combinations, such as loading while showing an error, unrepresentable. It also keeps `Calculator` focused on the form, while `useCalculation` owns the request lifecycle.

**Docker: nginx proxies `/api` instead of baking in the backend URL.** The frontend image is built with `VITE_API_URL=/`, so the browser calls `/api/v1/...` on its own origin, and nginx forwards those requests to the `backend` service on the Compose network.
- The image doesn't depend on where the backend runs, so the same build works in any environment.
- No cross-origin requests are needed in production.
- The backend doesn't have to be publicly reachable. It is published on 8080 only for convenience.

The trade-off is that nginx is coupled to the backend's service name. Both images run as non-root users. The backend runs on distroless (about 16 MB image, no shell).

**Accessibility and responsiveness.** Every input has a label. Validation errors are linked to their inputs with `aria-describedby`. Results are announced through a live region and errors use `role="alert"`. Enter submits the form. The layout is mobile-first: fields stack on small screens and form a single `a [op] b` row on wider ones. Inputs use `inputMode="decimal"`, so phones show a numeric keyboard.

## Assumptions

- **Percentage** means "`a` percent of `b`": `percentage(10, 50) = 5`. It is computed as `a × b / 100`, so integer inputs stay exact.
- **Square root** uses only `a`. If `b` is sent, it is ignored.
- **Unknown operations** return `400` rather than `404`, as specified. The route exists, and the operation name is treated as invalid input.
- Numbers are IEEE-754 doubles. Results that overflow, or have no real value, return `422`. `-0` counts as zero for division.
- The API returns full precision. The UI rounds results to 12 significant digits for display, so `0.1 + 0.2` shows as `0.3`.
- The UI accepts `.` as the decimal separator only.

## Possible improvements

- **CI** (GitHub Actions) running tests, `go vet`, `-race`, type-checks and Docker builds on every push.
- **End-to-end tests** (e.g. Playwright) against the Compose stack.
- **OpenAPI spec** for the API, used to generate or validate the frontend client.
- **Observability:** request-logging middleware with request IDs, Prometheus metrics, and a container health check. The distroless image has no shell, so this would need a `-healthcheck` flag on the binary.
- **Rate limiting** and JSON bodies for `405` responses, for consistency.
- **Arbitrary-precision decimals** (`math/big`) if exact decimal arithmetic matters, e.g. for money.
- **UX:** calculation history, a keypad for touch devices, and locale-aware number input (`1,5`).

## Use of AI tooling

AI assistance (Claude Code) was used during development to scaffold code, write tests and draft documentation. All output was reviewed, run and verified. The prompts used are listed in [PROMPTS.md](PROMPTS.md).
