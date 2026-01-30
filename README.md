# ClearStatus

Status page application: sign in with Google, manage organizations, teams, services, and incidents, and publish a public status page.

- **Backend:** Go API (MySQL, JWT)
- **Frontend:** Next.js (ShadcnUI)

---

## Run with Docker (easiest)

From the repo root:

```bash
cp .env.example .env
# Edit .env: set JWT_SECRET and (for Google sign-in) NEXT_PUBLIC_GOOGLE_CLIENT_ID + GOOGLE_CLIENT_ID

docker-compose up --build
```

- **Frontend:** http://localhost:3000  
- **Backend:** http://localhost:8080  

See [docker.md](docker.md) for details.

---

## Run from console (local)

You need **MySQL** running and **Node.js** / **Go** installed.

### 1. Backend

```bash
cd clearstatus-backend
cp .env.example .env
# Edit .env: PORT, DB_*, JWT_SECRET, GOOGLE_CLIENT_ID

go run cmd/server/main.go
```

Runs on **http://localhost:8080** (or the port in `.env`).

### 2. Frontend

In a **second terminal**:

```bash
cd clearstatus-frontend
cp .env.example .env.local
# Edit .env.local: NEXT_PUBLIC_API_URL=http://localhost:8080, NEXT_PUBLIC_GOOGLE_CLIENT_ID

npm install
npm run dev
```

Runs on **http://localhost:3000**.

---

## More

- Backend API and env: [clearstatus-backend/README.md](clearstatus-backend/README.md)
- Frontend and env: [clearstatus-frontend/README.md](clearstatus-frontend/README.md)
