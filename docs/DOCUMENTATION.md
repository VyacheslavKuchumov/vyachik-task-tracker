# Task Tracker Documentation

## Overview

This project is a Go backend + HTMX frontend task tracker where users can:

- Register and log in with JWT auth
- Create big goals
- Create smaller tasks under goals
- Assign tasks to users
- View tasks assigned to themselves

Main stack:

- Go + `chi` router
- PostgreSQL
- `golang-migrate` for SQL migrations
- HTMX for a minimal frontend dashboard
- `air` for local hot reload

## Project Structure

- `cmd/main.go`: API server entrypoint
- `cmd/api/api.go`: route wiring
- `cmd/migrate/main.go`: migration runner
- `cmd/migrate/migrations/`: SQL migrations
- `service/user/`: register/login logic
- `service/auth/`: JWT + password helpers
- `service/tracker/`: goals/tasks business logic + HTMX handlers
- `db/db.go`: PostgreSQL connection

## Prerequisites

- Go (1.25+)
- Docker + Docker Compose
- `air` installed (`go install github.com/air-verse/air@latest`)

## Environment Variables

Copy from `example.env` or edit `.env` directly:

| Variable | Default | Description |
|---|---|---|
| `PUBLIC_HOST` | `http://localhost` | Public base URL |
| `PORT` | `:8000` | API listen address |
| `DB_USER` | `postgres` | PostgreSQL username |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_HOST` | `127.0.0.1` | PostgreSQL host |
| `DB_PORT` | `5433` | PostgreSQL host port |
| `DB_NAME` | `task_tracker` | PostgreSQL database |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_EXP` | `604800` | JWT expiration in seconds |
| `JWT_SECRET` | `CHANGE_ME` | JWT signing secret |

## Run Locally

1. Start PostgreSQL:

```bash
docker compose up -d
```

2. Apply migrations:

```bash
make migrate-up
```

3. Run server with hot reload:

```bash
make air
```

Server URL:

- API + dashboard: `http://localhost:8000`
- HTMX dashboard page: `http://localhost:8000/`

## Common Make Targets

- `make test`: run tests
- `make migrate-up`: apply migrations
- `make migrate-down`: rollback all migrations
- `make migrate-force`: force migration version
- `make docker-up`: start Postgres container
- `make docker-down`: stop Postgres container
- `make air`: run app with hot reload

## Authentication

Login endpoint returns a JWT token:

- `POST /api/v1/login`

Use token in `Authorization` header:

- `Authorization: Bearer <token>`

Raw token without `Bearer` is also supported.

## Data Model

### users

- `id`, `first_name`, `last_name`, `email`, `password`, `created_at`

### goals

- `id`, `title`, `description`, `owner_id`, `created_at`

### tasks

- `id`, `goal_id`, `title`, `description`, `status`, `assignee_id`, `created_by`, `created_at`
- Status constraint: `todo`, `in_progress`, `done`

## API Reference

Base path: `/api/v1`

### 1) Register

- `POST /register`
- Auth: no

Request:

```json
{
  "firstName": "Alice",
  "lastName": "Smith",
  "email": "alice@example.com",
  "password": "secret123"
}
```

Response:

- `201 Created`

### 2) Login

- `POST /login`
- Auth: no

Request:

```json
{
  "email": "alice@example.com",
  "password": "secret123"
}
```

Response:

```json
{
  "token": "..."
}
```

### 3) Create Goal

- `POST /goals`
- Auth: required

Request:

```json
{
  "title": "Launch mobile app",
  "description": "Ship MVP with onboarding and first task flow"
}
```

Response:

```json
{
  "id": 1,
  "title": "Launch mobile app",
  "description": "Ship MVP with onboarding and first task flow",
  "ownerId": 1,
  "createdAt": "2026-02-12T18:00:00Z"
}
```

### 4) Get My Goals (with tasks)

- `GET /goals`
- Auth: required

Response:

```json
[
  {
    "id": 1,
    "title": "Launch mobile app",
    "description": "Ship MVP with onboarding and first task flow",
    "ownerId": 1,
    "createdAt": "2026-02-12T18:00:00Z",
    "tasks": [
      {
        "id": 3,
        "goalId": 1,
        "title": "Design onboarding",
        "description": "Create wireframes and copy",
        "status": "todo",
        "assigneeId": 2,
        "createdBy": 1,
        "createdAt": "2026-02-12T18:05:00Z"
      }
    ]
  }
]
```

### 5) Create Task Under Goal

- `POST /goals/{goalID}/tasks`
- Auth: required
- Only goal owner can create tasks under that goal

Request:

```json
{
  "title": "Set up CI",
  "description": "Add lint + test pipeline",
  "assigneeId": 2
}
```

Response:

```json
{
  "id": 5,
  "goalId": 1,
  "title": "Set up CI",
  "description": "Add lint + test pipeline",
  "status": "todo",
  "assigneeId": 2,
  "createdBy": 1,
  "createdAt": "2026-02-12T18:10:00Z"
}
```

### 6) Assign or Unassign Task

- `PUT /tasks/{taskID}/assign`
- Auth: required
- Only goal owner can assign/unassign tasks

Assign request:

```json
{
  "assigneeId": 2
}
```

Unassign request:

```json
{
  "assigneeId": null
}
```

### 7) Get Tasks Assigned To Me

- `GET /tasks/assigned`
- Auth: required

Response:

```json
[
  {
    "id": 5,
    "goalId": 1,
    "title": "Set up CI",
    "description": "Add lint + test pipeline",
    "status": "todo",
    "assigneeId": 2,
    "createdBy": 1,
    "createdAt": "2026-02-12T18:10:00Z"
  }
]
```

## Error Format

JSON API handlers return:

```json
{
  "error": "message"
}
```

Common status codes:

- `400`: invalid payload / invalid route param
- `403`: forbidden (auth failed or insufficient ownership)
- `500`: internal server/database error

## HTMX Frontend

Dashboard route:

- `GET /`

HTMX endpoints (JWT required):

- `GET /htmx/goals`
- `POST /htmx/goals/create`
- `GET /htmx/tasks/assigned`
- `POST /htmx/tasks/create`
- `POST /htmx/tasks/assign`

Flow:

1. Register and login on dashboard form
2. Token is stored in page input
3. HTMX requests include token in `Authorization` header
4. Create goals/tasks and refresh lists from UI buttons

## Migration Notes (MySQL -> PostgreSQL)

This codebase was migrated from MySQL to PostgreSQL:

- SQL placeholders changed from `?` to `$1...$n`
- Schema updated to PostgreSQL naming/types (`BIGSERIAL`, `TIMESTAMPTZ`, constraints)
- DB driver changed to `pgx` stdlib
- Migration driver changed to `migrate/database/postgres`
- Old product service removed and replaced by tracker domain

## Troubleshooting

### Port conflict on 5432

Project maps Postgres to host port `5433` by default to avoid conflicts:

- Compose mapping: `5433:5432`
- `.env`: `DB_PORT=5433`

### Auth always returns 403

- Confirm `Authorization` header has a valid JWT from `/api/v1/login`
- Confirm token uses same `JWT_SECRET` as backend config

### Migration errors

- Check database is healthy: `docker compose ps`
- Rerun: `make migrate-up`
- If migration state is dirty, use: `make migrate-force` and then `make migrate-up`
