# Learning Growth Platform

A minimal learning management MVP with a Go backend and a Vite frontend.

## What runs locally

- Backend: Go HTTP API on `http://localhost:8080`
- Frontend: Vite dev server on `http://localhost:5173`
- Database: MySQL 8.0.27 via Docker Compose

## Prerequisites

- Go 1.24+
- Node.js 18+
- Docker Desktop or another Docker Engine

## Start MySQL

The backend expects a local MySQL database named `learning_growth`.

```bash
docker compose up -d mysql
```

Default connection details:

- host: `127.0.0.1`
- port: `3306`
- database: `learning_growth`
- user: `root`
- password: `010511`

Example DSN:

```bash
MYSQL_DSN="root:010511@tcp(127.0.0.1:3306)/learning_growth?charset=utf8mb4&parseTime=True&loc=Local"
```

## Start the backend

From `backend/`:

```bash
copy .env.example .env
# edit .env if needed
# set MYSQL_DSN to your local MySQL instance
# set JWT_SECRET if you want a custom dev secret

go run ./cmd/server
```

If you use the Docker Compose MySQL service, the example DSN in `.env.example` already matches it.

## Start the frontend

From `frontend/`:

```bash
npm install
npm run dev
```

If you need the frontend to point at a non-default API origin, set `VITE_API_BASE_URL` before running Vite.

## Validation

Run the backend integration test first. It skips cleanly if MySQL is unavailable.

From `backend/`:

```bash
go test ./internal/integration -run TestMVPFlow -v -count=1 -timeout=180s
go test ./... -count=1 -timeout=180s
```

From `frontend/`:

```bash
npm test -- --run
npm run build
```

## MVP flow covered by integration test

- Register a new user
- Log in with the same user
- Create a subject
- Create a done task for today
- Create a study session for today
- Check in for today
- Fetch `/api/stats/overview` and verify the totals
