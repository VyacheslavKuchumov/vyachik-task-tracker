# Task Tracker

## Structure

- `server/`: Go API server and database migrations
- `web/`: Nuxt frontend (Nuxt UI + Pinia)

## Run Full Stack With One Command

```bash
docker compose up -d --build --remove-orphans
```

Services:

- Frontend: `http://localhost:3000`
- API: `http://localhost:8000/api/v1`
- Postgres host port: `localhost:5433`

## Run Backend Locally

```bash
cd server
make docker-db-up
make migrate-up
make air
```
