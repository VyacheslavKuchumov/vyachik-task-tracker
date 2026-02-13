# Task Tracker

## Structure

- `backend/`: Go API server and database migrations
- `frontend/`: reserved for Nuxt app (currently empty)

## Run Backend With Docker

```bash
docker compose up -d --build
```

API base: `http://localhost:8000/api/v1`

## Run Backend Locally

```bash
cd backend
make docker-db-up
make migrate-up
make air
```
