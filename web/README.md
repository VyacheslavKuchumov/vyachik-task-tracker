# Task Tracker Frontend (Nuxt + Nuxt UI + Pinia)

## What this frontend does

- JWT auth against the Go server
- Register/login pages
- Goal management
- Task creation inside goals
- Task assignment by user ID
- "Assigned to me" task list

## Backend contract

This app expects the Go server to expose:

- `POST /api/v1/register`
- `POST /api/v1/login`
- `GET /api/v1/goals`
- `POST /api/v1/goals`
- `POST /api/v1/goals/{goalId}/tasks`
- `GET /api/v1/tasks/assigned`
- `PUT /api/v1/tasks/{taskId}/assign`

The frontend calls Nuxt server routes under `/api/...`, and those routes proxy requests to the Go backend URL.

## Environment

Create `.env` in `web/`:

```env
BACKEND_URL=http://127.0.0.1:8000
```

## Run

```bash
cd web
npm install
npm run dev
```

Frontend: `http://localhost:3000`

## Playwright E2E

These tests cover the main user flow in the UI:

- registration and authenticated navigation
- goal creation
- task creation/editing under a goal
- users/tasks board visibility
- profile update and password change

Prerequisites:

1. backend API available on `http://127.0.0.1:8000` (or set `PLAYWRIGHT_BACKEND_URL`)
2. database/backend started before test run
3. frontend app started on `http://127.0.0.1:3000` (or set `PLAYWRIGHT_BASE_URL`)

Run:

```bash
cd web
npm run test:e2e
```
