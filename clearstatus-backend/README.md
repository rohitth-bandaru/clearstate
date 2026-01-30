# ClearStatus Backend

Go API for the ClearStatus status page application. Provides auth (Google + JWT), organizations, teams, services, incidents, and a public status endpoint.

## Requirements

- Go 1.21+
- MySQL 8+

## Environment

All configuration is loaded from environment variables; there are no hardcoded defaults. Copy `.env.example` to `.env` and set values. The server loads `.env` from the current working directory at startup (optional; ignored if the file is missing).

| Variable | Required | Description |
|----------|----------|-------------|
| PORT | Yes | Server port (e.g. `8080`) |
| DB_DRIVER | Yes | Database driver (e.g. `mysql`) |
| DB_DSN | Yes* | MySQL DSN (single URL). *Omit if using DB_HOST/DB_USER/DB_PASSWORD/DB_NAME. |
| DB_HOST | Yes* | Database host (e.g. `localhost`). *Required when DB_DSN is not set. |
| DB_PORT | No | Database port (default `3306` when using parts). |
| DB_USER | Yes* | Database user. *Required when DB_DSN is not set. |
| DB_PASSWORD | Yes* | Database password. *Required when DB_DSN is not set. |
| DB_NAME | Yes* | Database name (e.g. `clearstatus`). *Required when DB_DSN is not set. |
| JWT_SECRET | Yes | Secret for signing JWTs |
| GOOGLE_CLIENT_ID | No | Google OAuth client ID (for verifying ID tokens; required for sign-in) |
| PRODUCTION | No | Set to `true` in production |
| LOG_LEVEL | No | Log level: `debug`, `info`, `warn`, `error` (default `info`) |
| LOG_FORMAT | No | Log format: `json` (default) or `text` |

Use either **DB_DSN** (single URL) or **DB_HOST**, **DB_USER**, **DB_PASSWORD**, **DB_NAME** (and optionally **DB_PORT**). When using parts, the connection URL is built automatically and the DB connection is logged at startup (user@host:port/database, password never logged).

CORS is not used; all clients are expected to be in the same location. Startup failures and request errors are logged with structured logging (slog); the server exits with `os.Exit(1)` on startup errors (no panic).

### Generate JWT secret

Use a cryptographically secure random string for `JWT_SECRET`. Examples:

```bash
# Base64 (43 chars)
openssl rand -base64 32

# Hex (64 chars)
openssl rand -hex 32
```

Set the output in `.env` as `JWT_SECRET=<generated-value>`.

## Run locally

1. Create a MySQL database: `CREATE DATABASE clearstatus;`
2. Copy `.env.example` to `.env` and set all required variables.
3. Run: `go run ./cmd/server`

Migrations run automatically on startup from `migrations/001_init.up.sql`.

## API

- `POST /api/auth/google` — Body: `{ "id_token": "..." }`. Returns `{ "token": "<JWT>", "user": {...} }`.
- `GET/POST /api/orgs` — List, create orgs (JWT).
- `GET/PUT/DELETE /api/orgs/:orgID` — Get, update, delete org (JWT).
- `GET/POST /api/orgs/:orgID/members` — List, add members (JWT).
- `GET/POST /api/orgs/:orgID/teams`, `GET/PUT/DELETE /api/orgs/:orgID/teams/:teamID` — Teams (JWT).
- `GET/POST /api/orgs/:orgID/services`, `GET/PUT/DELETE /api/orgs/:orgID/services/:serviceID` — Services; updates use `version` for optimistic locking (JWT).
- `GET/POST /api/orgs/:orgID/incidents`, `GET/PUT/DELETE /api/orgs/:orgID/incidents/:incidentID`, `POST /api/orgs/:orgID/incidents/:incidentID/updates` — Incidents; updates use `version` (JWT).
- `GET /api/public/orgs/:slug/status` — Public status (no auth). For polling.

## Deploy

See [deploy.md](deploy.md) for AWS deployment.
