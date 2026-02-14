# Task Tracker

Task Tracker is a full-stack goals and tasks application:

- Go API (`server/`)
- Nuxt 4 frontend (`web/`)
- PostgreSQL database

Frontend views:

- Home: your assigned tasks
- Goals: goals list and navigation to goal tasks
- Goal Tasks (`/tasks/:goalId`): task CRUD with user lookup assignment
- Users (`/users`): all users and their current tasks
- Profile: update first name/last name and password

## Quick Start

Make sure these DNS records point to your VPS public IP before starting:

- `home.vyachik-dev.ru`
- `home-server.vyachik-dev.ru`

Set your ACME email (used by Traefik Let's Encrypt):

```bash
export TRAEFIK_ACME_EMAIL=slavakuchumov@gmail.com
```

```bash
docker compose up -d --build --remove-orphans
```

Open:

- Frontend: `https://home.vyachik-dev.ru`
- Backend API: `https://home-server.vyachik-dev.ru/api/v1`
- Swagger UI (requires auth): `https://home-server.vyachik-dev.ru/swagger/index.html`

## Documentation

- `docs/README.md`: documentation index
- `docs/SETUP.md`: setup and local development
- `docs/ARCHITECTURE.md`: component and data flow
- `docs/API.md`: API contract
- `docs/OPERATIONS.md`: migrations, testing, troubleshooting
- `docs/CONTRIBUTING.md`: branch, commit, and pull request guidelines

## Repository Structure

- `server/`: backend API, migrations, tests
- `web/`: frontend app and API proxy routes
- `docs/`: project docs
