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

## UI Tests (Selenium)

Selenium UI tests are in `web/selenium_ui_test.go`.

They verify:

- `/login` and `/register` are accessible without authentication
- unauthenticated users are redirected away from protected pages
- authenticated users can open protected pages

Run them with a Selenium server:

```bash
# Example: start selenium standalone chrome
docker run --rm -d -p 4444:4444 --name selenium selenium/standalone-chrome:latest

# Run tests (set endpoint explicitly)
SELENIUM_URL=http://localhost:4444/wd/hub go test ./web -run TestSelenium_PageAccessControl -v
```
