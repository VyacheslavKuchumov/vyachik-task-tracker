# Project Documentation

This folder contains documentation for the full Task Tracker project.

## Contents

- `docs/SETUP.md`: local and Docker setup
- `docs/ARCHITECTURE.md`: system design and data flow
- `docs/API.md`: backend API contract
- `docs/OPERATIONS.md`: migrations, tests, and troubleshooting

## Quick Start

### Run full stack with Docker

```bash
docker compose up -d --build --remove-orphans
```

Open:

- Frontend: `http://localhost:3000`
- Backend API base: `http://localhost:8000/api/v1`

### Local development split

Backend:

```bash
cd server
cp example.env .env
make docker-db-up
make migrate-up
make air
```

Frontend:

```bash
cd web
cp .env.example .env
npm install
npm run dev
```
