# Copilot Instructions for incd-mgnt-system

Use these concise, repository-specific notes to help AI coding agents be productive immediately.

1. Big picture
   - This is a monorepo with a Go backend (cmd/server) and a Vue 3 + TypeScript frontend (web). See `Low-Level-Architecture.md` for detailed flows.
   - Main runtime: backend handles Alertmanager webhooks, groups alerts into incidents, persists (in-memory or Postgres), and sends notifications (Slack/Email/Telegram).
   - Key folders: `cmd/` (entrypoints), `internal/handlers`, `internal/services`, `internal/storage`, `web/` (frontend), `migrations/` (DB schema).

2. How to run locally (developer-first)
   - Frontend dev server: `cd web && npm install && npm run dev` (port 5173). Uses Vite with hot reload and API proxy to backend.
   - Backend dev: `go run cmd/server/main.go` (port from PORT env) or `go install github.com/cosmtrek/air@latest && air` for hot reload.
   - Infrastructure: `docker-compose up -d` (Postgres, Prometheus, Alertmanager, Node Exporter).
   - Full local stack: `docker-compose --profile development up -d` (includes frontend dev server via `docker-compose.app.yaml`).

3. Configuration patterns
   - All runtime config from environment variables; defaults and validation in `internal/config/config.go`.
   - Critical envs: `DATABASE_URL`, `PORT`, `LOG_LEVEL`, `ALERTMANAGER_URL`, `SLACK_TOKEN`+`SLACK_CHANNEL`, `EMAIL_SMTP_*`, `TELEGRAM_*`.
   - `LoadAndValidateConfig()` enforces required combinations (e.g., Slack token requires channel) — extend `Validate()` when adding new config.
   - Hot-reloadable: log level, timeouts, CORS, debug mode; non-reloadable: port, DB, TLS, notification credentials.

4. Storage & migrations
   - Pluggable storage: memory (`internal/storage/memory.go`) or Postgres (`internal/storage/postgres.go`).
   - When adding persistence, implement `Store` interface and follow repository patterns in `Low-Level-Architecture.md`.
   - Migrations: `migrations/*.up.sql` and `*.down.sql`; run with `go run cmd/migrate/main.go up/down`.
   - Connection pooling: configure `DB_MAX_OPEN_CONNS` (default 25), `DB_MAX_IDLE_CONNS` (5), `DB_CONN_MAX_LIFETIME` (5m).

5. Services & patterns to follow
   - Business logic in `internal/services/*`; inject dependencies (store, metrics, logger) via constructors.
   - Handlers in `internal/handlers` are thin HTTP adapters; validate input, call services, return JSON.
   - Notifications: template-first with `NotificationTemplateService`, batching via `NotificationBatchProcessor`, retries via `internal/retry`.
   - Metrics: use `metricsService.RecordDBQuery()`, `RecordIncidentCreated()`; structured logging with `logger.Info()/Error()`.
   - Authentication: JWT-based with `AuthService`; user management via `UserService`.

6. Tests & CI
   - Unit tests: `go test ./...` (CI uses `-short`); integration tests require `TEST_DATABASE_URL`.
   - Config tests: validate combinations in `internal/config/config_test.go`.
   - Frontend: `npm run type-check`, `npm run lint`; E2E with Playwright in `web/e2e/`.
   - CI workflow: `.github/workflows/ci.yml` — mirror environment for local integration tests.

7. Common pitfalls
   - Memory store loses data on restart; use `DATABASE_URL` in tests for Postgres behavior.
   - Config validation is strict — update `Validate()` and tests when adding env vars.
   - Frontend proxy: set `VITE_API_TARGET=http://host.docker.internal:8080` when running frontend outside Docker.
   - API routes must be prefixed with `/api/` for proxy; frontend routes are SPA catch-all.
   - Notification channels require complete config (e.g., Slack needs both token and channel).

8. Useful file references (examples)
   - App bootstrap: `cmd/server/main.go` (service wiring, middleware stack)
   - Config schema: `internal/config/config.go` (validation rules, env parsing)
   - Alert grouping: `internal/services/alert.go` (webhook processing, fingerprint deduplication)
   - Notification flow: `internal/services/notification.go` (templating, batching, retry logic)
   - Storage interface: `internal/storage/memory.go` or `postgres.go` (CRUD patterns)
   - Frontend proxy: `web/vite.config.ts` and `web/FRONTEND_INTEGRATION.md`
   - Docker setup: `docker-compose.app.yaml` (full app with profiles); `docker-compose.yml` (infrastructure only)

9. When creating changes
   - Update `Low-Level-Architecture.md` and `README.md` for architecture changes.
   - Add unit tests under `internal/..._test.go`; integration tests mirroring CI.
   - Follow service injection pattern; use structured logging and metrics.
   - Frontend: build with `npm run build`, deploy dist/ contents; ensure API compatibility.
   - Database changes: add migrations, update storage interfaces, test with both memory and Postgres stores.

If anything above is unclear or you want a different level of detail (more walkthroughs, examples, or editable templates for tests), tell me which sections to expand.