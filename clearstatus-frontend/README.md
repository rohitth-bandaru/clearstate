# ClearStatus Frontend

Next.js frontend for the ClearStatus status page application. Sign in with Google, manage organizations, teams, services, and incidents, and view the public status page with 2s polling (UI updates only when data changes).

## Requirements

- Node.js 18+
- Backend API running (see clearstatus-backend)
- Google OAuth Client ID (Web application)

## Environment

All configuration is loaded from environment variables. Copy `.env.example` to `.env.local` and set values. No hardcoded URLs or secrets in code.

| Variable | Required | Description |
|----------|----------|-------------|
| NEXT_PUBLIC_API_URL | Yes | Backend API base URL (e.g. `http://localhost:8080`) |
| NEXT_PUBLIC_GOOGLE_CLIENT_ID | Yes | Google OAuth 2.0 Client ID (Web application). Same as backend. |
| NEXT_PUBLIC_POLL_INTERVAL_MS | No | Status page poll interval in ms (default 2000) |

## Run locally

1. Copy `.env.example` to `.env.local` and set all required variables.
2. Run: `npm run dev`
3. Open [http://localhost:3000](http://localhost:3000).
4. Sign in with Google (ensure backend is running; clients are expected to be in the same location as the API).

## Features

- **Sign in with Google** — Single auth flow; account created on first sign-in.
- **Organizations** — Create and manage orgs; switch via sidebar.
- **Teams** — CRUD teams per org.
- **Services** — CRUD services; set status (operational, degraded, partial outage, major outage). Optimistic locking with `version`.
- **Incidents** — Create incidents/maintenance, link services, add updates. Optimistic locking with `version`.
- **Public status page** — `/status/[orgSlug]`. Polls every 2s; re-renders only when response changes (deep equality).

## Deploy

See [deploy.md](deploy.md) for Vercel deployment.
