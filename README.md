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

Generate root `.env` (contains ACME email, hosts, DB password, JWT secret):

```bash
python3 scripts/generate_compose_env.py \
  --acme-email your-email@example.com \
  --web-host home.your-domain.tld \
  --api-host home-server.your-domain.tld
```

```bash
docker compose up -d --build --remove-orphans
```

Open:

- Frontend: `https://<TRAEFIK_WEB_HOST>`
- Backend API: `https://<TRAEFIK_API_HOST>/api/v1`
- Swagger UI (requires auth): `https://<TRAEFIK_API_HOST>/swagger/index.html`

## Standalone Docker Dev Script

Run Postgres + backend + frontend as separate containers (without changing `docker-compose.yml`):

```bash
python3 scripts/dev_docker_stack.py up
```

Stop containers:

```bash
python3 scripts/dev_docker_stack.py down --remove-network
```

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
