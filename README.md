# Task Tracker (Go + HTMX + PostgreSQL)

## Quick Start

```bash
docker compose up -d --build
```

Open:

- Login page: `http://localhost:8000/login`
- Register page: `http://localhost:8000/register`
- Goals page: `http://localhost:8000/goals`
- Tasks page: `http://localhost:8000/tasks`
- API base: `http://localhost:8000/api/v1`

For local hot-reload development instead of containers:

```bash
make docker-db-up
make migrate-up
make air
```

## Full Docs

Detailed documentation is available at:

- `docs/DOCUMENTATION.md`

It includes:

- Setup and environment variables
- Architecture and folder structure
- Complete API reference with payload examples
- HTMX frontend behavior and endpoints
- Troubleshooting and migration notes
