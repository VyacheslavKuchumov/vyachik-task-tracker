# Task Tracker

Task Tracker is a full-stack goals and tasks application:

- Go API (`server/`)
- Nuxt 4 frontend (`web/`)
- PostgreSQL database

## Quick Start

```bash
docker compose up -d --build --remove-orphans
```

Open:

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8000/api/v1`
- Swagger UI: `http://localhost:8000/swagger/index.html`

## Documentation

- `docs/README.md`: documentation index
- `docs/SETUP.md`: setup and local development
- `docs/ARCHITECTURE.md`: component and data flow
- `docs/API.md`: API contract
- `docs/OPERATIONS.md`: migrations, testing, troubleshooting

## Repository Structure

- `server/`: backend API, migrations, tests
- `web/`: frontend app and API proxy routes
- `docs/`: project docs
