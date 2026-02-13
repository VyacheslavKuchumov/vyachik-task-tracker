# AGENTS.md

## Purpose

This document gives coding agents project-specific context for the `vyachik-task-tracker` repository.

## Project Summary

Task Tracker is a full-stack app with:

- Go backend API (`server/`)
- Nuxt 4 frontend (`web/`)
- PostgreSQL database
- SQL migrations managed with `golang-migrate`
- Docker Compose orchestration at repository root

## Repository Layout

- `server/`: Go API, business logic, migrations, tests
- `web/`: Nuxt UI app, Pinia stores, server API proxy routes
- `docker-compose.yml`: full stack runtime (postgres + server + web)
- `docs/`: project documentation

## Core Runtime Flow

1. Frontend calls local Nuxt server routes under `web/server/api/...`.
2. Nuxt route handlers proxy to Go API `/api/v1/...` using `callBackend`.
3. Go API validates JWT and executes user/tracker operations against PostgreSQL.

## Commands Agents Should Use

### Full stack (recommended)

```bash
docker compose up -d --build --remove-orphans
```

### Backend local development

```bash
cd server
cp example.env .env
make docker-db-up
make migrate-up
make air
```

### Frontend local development

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

## Testing and Verification

Current automated test suites in this repository:

- Backend Go tests: `cd server && make test`
- Frontend Playwright e2e tests: `cd web && npm run test:e2e`

Before finishing backend-impacting changes:

```bash
cd server
make test
```

For frontend-impacting changes:

```bash
cd web
npm run build
```

For frontend user-flow/auth/API integration changes, run end-to-end tests (requires running backend + frontend):

```bash
cd web
npm run test:e2e
```

Playwright helper commands:

```bash
cd web
npm run test:e2e:headed
npm run test:e2e:ui
npm run test:e2e:report
```

If API handlers or payloads changed, regenerate OpenAPI docs:

```bash
cd server
make swagger
```

## Environment Notes

- Backend expects env vars from `server/.env` in local development.
- Frontend expects `BACKEND_URL` in `web/.env`.
- Default local ports:
  - Web: `3000`
  - API: `8000`
  - Postgres host: `5433`

## API and Auth Notes

- API base path is `/api/v1`.
- Protected backend routes accept either:
  - `Authorization: Bearer <token>`
  - `task_tracker_token` cookie
- Frontend currently uses `Authorization` header from Pinia token state for protected calls.

## Migration Notes

- Migration files are in `server/cmd/migrate/migrations`.
- Migration file names include older labels (`products`, `orders`) but actual schema is `users`, `goals`, `tasks`.
- Run migrations with:

```bash
cd server
make migrate-up
```

## Change Discipline

- Keep backend contracts and frontend proxy routes in sync.
- If endpoint shape changes, update:
  - `web/server/api/...`
  - `web/app/stores/...`
  - `docs/API.md`
- Prefer minimal, focused edits and keep tests green.

## Git Workflow Reminder

- Do not work directly on `main` (it is protected).
- Before starting development and creating a feature branch, sync with latest `origin/main`:
  - `git fetch origin`
  - `git switch main`
  - `git pull --ff-only origin main`
- Create and use a feature branch for every task.
- Before committing, run relevant verification commands:
  - backend-impacting changes: `cd server && make test`
  - frontend-impacting changes: `cd web && npm run build`
  - frontend user-flow/auth/API integration changes: `cd web && npm run test:e2e`
- Commit all agent-made changes to that feature branch.
- Open a pull request from the feature branch to `main` after pushing changes.
- Do not merge directly to `main`; merging is done via the pull request workflow.
- After opening a PR, reset local DB state so a human can re-verify from clean data:
  - from repository root: `docker compose down -v --remove-orphans`

## Commit Message Guidelines

- Use clear, imperative commit subjects (example: `fix: validate empty task title`).
- Prefer `type: short summary` format (for example: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`).
- Keep the subject concise (around 50-72 chars) and focused on one change.
- Avoid vague messages like `update`, `changes`, or `fix stuff`.

## Pull Request Guidelines

- Keep each PR scoped to one logical change.
- PR title should summarize the change and match the commit style where possible.
- In PR description, include:
  - what changed
  - why it changed
  - how it was tested (commands run)
- If backend/frontend contracts changed, update relevant docs and proxy/store code in the same PR.
- If API shapes changed, regenerate Swagger docs before opening PR (`cd server && make swagger`).
