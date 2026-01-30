# Running ClearStatus with Docker (local)

Use Docker Compose to run the full application (MySQL + backend + frontend) locally so others can test without installing Go, Node, or MySQL.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) (or Docker Desktop, which includes Compose)

## Quick start

1. **From the repo root**, copy the example env and set required values:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` and set at least:

   - `JWT_SECRET` – any secure string (e.g. `openssl rand -base64 32`)
   - `DB_PASSWORD` – MySQL root password (defaults to `clearstatus` if unset)
   - Optional: `GOOGLE_CLIENT_ID` and `NEXT_PUBLIC_GOOGLE_CLIENT_ID` for Google sign-in

2. **Start everything:**

   ```bash
   docker-compose up --build
   ```

3. **Open in the browser:**

   - **Frontend:** [http://localhost:3000](http://localhost:3000)
   - **Backend API:** [http://localhost:8080](http://localhost:8080)

4. **Stop:** `Ctrl+C`, then `docker-compose down` (add `-v` to remove the MySQL volume).

## Environment variables (root `.env`)

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_PASSWORD` | MySQL root password (used for backend and MySQL container) | `clearstatus` |
| `JWT_SECRET` | Backend JWT signing secret | `change-me-in-production` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID (optional) | - |
| `NEXT_PUBLIC_API_URL` | URL the **browser** uses to reach the API | `http://localhost:8080` |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | Same as `GOOGLE_CLIENT_ID` for frontend build (optional) | - |

**Important:** `NEXT_PUBLIC_API_URL` must be the URL the browser uses. With the default port mapping, that is `http://localhost:8080`. Do not use `http://backend:8080` (that only works inside the Docker network).

## Services

- **mysql** – MySQL 8, port 3306, database `clearstatus`. Backend runs migrations on startup.
- **backend** – Go API, port 8080, waits for MySQL to be healthy.
- **frontend** – Next.js app, port 3000, depends on backend for startup order.

## Building or running a single service

```bash
# Rebuild only the frontend
docker-compose up --build frontend

# Run backend and MySQL only (e.g. to test API alone)
docker-compose up backend
```

## Troubleshooting

- **Backend fails with "database ping failed"** – MySQL may still be starting. Wait a few seconds and try again, or run `docker-compose restart backend`.
- **Frontend shows API errors** – Ensure `NEXT_PUBLIC_API_URL` is `http://localhost:8080` (or the host/port where the backend is reachable from your browser). Rebuild the frontend after changing it: `docker-compose up --build frontend`.
- **Google sign-in doesn’t work** – Set `GOOGLE_CLIENT_ID` and `NEXT_PUBLIC_GOOGLE_CLIENT_ID` in `.env`, add `http://localhost:3000` (and your domain) to the OAuth client’s authorized origins, then rebuild the frontend.
