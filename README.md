# Parish API

REST API backend for a parish website. Manages schedules, events, parish groups, materials, users, and roles with fine-grained permission control.

## Prerequisites

- Go 1.22+
- Docker & Docker Compose

## Quick start

```bash
# Start local infrastructure and the API server in one command
make run-local

# In another terminal, run the smoke tests
make smoke-test
```

The API will be available at `http://localhost:8080`. An admin user (`admin@parish.local`) is seeded automatically on startup.

## Makefile targets

| Target | Description |
|---|---|
| `make help` | Print all available targets |
| `make build` | Compile the binary to `bin/parish-api` |
| `make run` | Run the server directly (requires env vars to be set) |
| `make test` | Run all unit tests |
| `make clean` | Remove build artifacts |
| `make mocks` | Regenerate repository mocks via `go generate` |
| `make infra-start` | Start Firestore emulator + Redis via Docker Compose |
| `make infra-stop` | Tear down local infrastructure |
| `make run-local` | Start infrastructure, seed admin, and run the API |
| `make smoke-test` | Execute happy-path bash tests against the running API |

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `LOG_LEVEL` | `info` | Log verbosity |
| `GCP_PROJECT_ID` | — | Google Cloud project (required) |
| `FIRESTORE_EMULATOR_HOST` | — | Set to `host:port` to use the Firestore emulator (e.g. `localhost:8081` with `make run-local`) |
| `REDIS_URL` | `localhost:6379` | Redis address for token storage and rate limiting |
| `COOKIE_SECURE` | `true` | Set to `false` for local HTTP (no TLS) |
| `CORS_ORIGIN` | `http://localhost:3000` | Allowed origin for CORS credentials |
| `TLS_CERT` | — | Path to TLS certificate file (enables HTTPS when set with `TLS_KEY`) |
| `TLS_KEY` | — | Path to TLS private key file |
| `ADMIN_EMAIL` | — | Seeded admin user email (idempotent on startup) |
| `ADMIN_PASSWORD` | — | Seeded admin user password |

## Project structure

```
cmd/
  main.go              Application entry point
  config.go            Environment-based configuration
  container.go         Dependency injection setup (uber-go/dig)
  seed.go              Idempotent admin user/role seed on startup
  rest/
    rest.go            HTTP router (Go 1.22+ method-based routing)
    handler/           HTTP handlers and request/response DTOs
    middleware/        Auth, CORS, rate limiter, security headers, logging
internal/
  domain/              Entity definitions, validation, domain errors
  repository/          Repository interfaces + generated mocks (matryer/moq)
    firestore/         Google Cloud Firestore (native) implementations
  usecase/             Business logic interfaces + implementations
  cache/               Cache interface + Redis implementation
scripts/
  smoke_test.sh        Bash-based API integration tests
infra/
  firestore/           firestore.indexes.json + README for composite indexes (deploy to GCP / Firebase)
docker-compose.yml     Local Firestore emulator + Redis
```

## Design choices

### Clean Architecture

The codebase follows Clean Architecture with three distinct layers:

- **Domain** (`internal/domain`) — Pure entities, value objects, validation rules, and domain errors. No external dependencies.
- **Use cases** (`internal/usecase`) — Business logic defined as interfaces, each method accepting a typed input struct (e.g. `CreateEventInput`) to keep signatures clean and extensible.
- **Infrastructure** — Repository implementations (`internal/repository/firestore`), cache (`internal/cache`), HTTP handlers (`cmd/rest/handler`), and middleware.

Dependencies flow inward: handlers depend on use cases, use cases depend on repository and cache interfaces, and the domain depends on nothing.

### API responses (transport layer)

Success payloads do not return raw persistence entities. Each resource has a dedicated response type in [`internal/domain/response.go`](internal/domain/response.go) (e.g. `EventResponse`, `ScheduleResponse`) with only fields needed by clients. Domain types expose `ToResponse()` methods; list endpoints use helpers like `ToEventResponses` so empty results serialize as `[]` instead of `null`. Audit metadata (`createdAt`, `updatedAt`, `createdBy`, `updatedBy`) is omitted from public JSON.

### Authentication & Authorization

- **Token-based auth** — Login returns a random token stored in Redis with a 24-hour TTL. The token is delivered as an `HttpOnly` secure cookie (`auth_token`), with a `Bearer` header fallback.
- **Role-Based Access Control (RBAC)** — Users are assigned roles, each carrying a list of resource-level permissions (`read`/`write`). The `RequirePermission` middleware checks these before every protected endpoint.
- **Admin seeding** — On startup, an admin role and user are created idempotently from the `ADMIN_EMAIL`/`ADMIN_PASSWORD` env vars, avoiding bootstrap chicken-and-egg problems.

### Caching & Rate Limiting

Both the token store and the per-IP rate limiter are backed by Redis through a thin `cache.Cache` interface (`internal/cache`). This allows the API to scale horizontally without sharing in-memory state, and makes TTL-based expiry automatic. A `CacheMock` generated by `moq` (same as the repository mocks) is used in unit tests so they run without a Redis dependency.

### Error Handling

Domain errors use prefixed sentinel messages (e.g. `ERR_NOT_FOUND:`, `ERR_UNAUTHORIZED:`) that are mapped to HTTP status codes by a single `HandleDomainError` function. This centralizes the translation between business errors and REST semantics.

### Mock Generation

Repository mocks are generated with [matryer/moq](https://github.com/matryer/moq) via `//go:generate` directives. Run `make mocks` after changing any repository interface.

### TLS

When `TLS_CERT` and `TLS_KEY` are set the server starts with HTTPS. For local development these are left empty so the server runs plain HTTP.

### Structured Logging

All error paths in use cases emit structured JSON logs via the standard library `log/slog` package.

## API endpoints

### Public

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/login` | Login (sets auth cookie) |
| `POST` | `/api/auth/logout` | Logout (clears auth cookie) |
| `GET` | `/api/schedule` | Get parish schedule |
| `GET` | `/api/events` | List events |
| `GET` | `/api/events/{id}` | Get event by ID |
| `GET` | `/api/parish-groups` | List parish groups |
| `GET` | `/api/parish-groups/{id}` | Get parish group by ID |
| `GET` | `/api/materials` | List materials (supports `?type=` and `?label=` filters) |
| `GET` | `/api/materials/{id}` | Get material by ID |

### Protected (requires authentication + permission)

| Method | Path | Permission | Description |
|---|---|---|---|
| `PUT` | `/api/schedule` | `schedule:write` | Update schedule |
| `POST` | `/api/events` | `events:write` | Create event |
| `PUT` | `/api/events/{id}` | `events:write` | Update event |
| `DELETE` | `/api/events/{id}` | `events:write` | Delete event |
| `POST` | `/api/parish-groups` | `parish_groups:write` | Create parish group |
| `PUT` | `/api/parish-groups/{id}` | `parish_groups:write` | Update parish group |
| `DELETE` | `/api/parish-groups/{id}` | `parish_groups:write` | Delete parish group |
| `POST` | `/api/materials` | `materials:write` | Create material |
| `PUT` | `/api/materials/{id}` | `materials:write` | Update material |
| `DELETE` | `/api/materials/{id}` | `materials:write` | Delete material |
| `GET` | `/api/roles` | `roles:read` | List roles |
| `GET` | `/api/roles/{id}` | `roles:read` | Get role by ID |
| `POST` | `/api/roles` | `roles:write` | Create role |
| `PUT` | `/api/roles/{id}` | `roles:write` | Update role |
| `DELETE` | `/api/roles/{id}` | `roles:write` | Delete role |
| `PUT` | `/api/users/{id}/roles` | `roles:write` | Assign roles to a user |
