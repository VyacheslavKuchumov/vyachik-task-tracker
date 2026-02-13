# Architecture

## High-Level Components

- `web` (Nuxt 4): UI, auth state, API proxy routes
- `server` (Go + chi): auth, goals, tasks, business rules
- `postgres` (PostgreSQL 16): persistent data store

## Directory Overview

### Backend (`server/`)

- `cmd/main.go`: API app entrypoint
- `cmd/server/server.go`: HTTP router and service wiring
- `cmd/migrate/main.go`: migration runner
- `cmd/migrate/migrations/`: SQL migrations
- `service/user/`: register/login handlers and store
- `service/auth/`: JWT creation/validation and password hashing
- `service/tracker/`: goals/tasks handlers and store
- `types/`: API and domain structs
- `db/db.go`: PostgreSQL connection

### Frontend (`web/`)

- `app/pages/`: routes (`/`, `/login`, `/signup`)
- `app/components/`: UI cards, navbar, task board
- `app/stores/`: Pinia stores (`auth`, `tracker`)
- `app/middleware/auth.global.js`: route protection
- `server/api/`: Nuxt server route proxies
- `server/utils/backend.ts`: proxy helper used by API routes

## Request Flow

### Public auth flow

1. UI submits credentials to Nuxt route (`/api/auth/login` or `/api/auth/register`).
2. Nuxt route forwards to backend `/api/v1/login` or `/api/v1/register`.
3. Backend validates and returns JSON payload (JWT token on login).
4. Frontend stores token in Pinia persisted state.

### Protected flow

1. Frontend includes `Authorization: Bearer <token>`.
2. Nuxt server route enforces header presence (`requireAuth: true`).
3. Backend middleware validates JWT and loads user from DB.
4. Handler executes goal/task operation.

## Data Model

### `users`

- `id`, `first_name`, `last_name`, `email`, `password`, `created_at`

### `goals`

- `id`, `title`, `description`, `owner_id`, `created_at`

### `tasks`

- `id`, `goal_id`, `title`, `description`, `status`, `assignee_id`, `created_by`, `created_at`
- `status` allowed values: `todo`, `in_progress`, `done`

## Authorization Rules

- Goal owner can create tasks under that goal.
- Goal owner can assign tasks in that goal.
- User can list only their own goals (`GET /goals`).
- User can list tasks assigned to themselves (`GET /tasks/assigned`).
