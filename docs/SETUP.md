# Setup Guide

## Prerequisites

- Docker + Docker Compose
- Go `1.25.x` (for backend local development)
- Node.js `>=20.9.0` (for frontend local development)
- `air` for backend hot reload:

```bash
go install github.com/air-verse/air@latest
```

## Option 1: Run Entire Stack with Docker

From repository root:

1. Point these DNS records to your server public IP:

- `home.vyachik-dev.ru`
- `home-server.vyachik-dev.ru`

2. Set an ACME email for Let's Encrypt certificate registration:

```bash
export TRAEFIK_ACME_EMAIL=slavakuchumov@gmail.com
```

3. Start the stack:

```bash
docker compose up -d --build --remove-orphans
```

Services:

- `web` at `https://home.vyachik-dev.ru`
- `server` at `https://home-server.vyachik-dev.ru`
- `postgres` at host port `5433`

Stop:

```bash
docker compose down
```

## Option 2: Local Backend + Local Frontend

### 1. Start Postgres

```bash
cd server
make docker-db-up
```

### 2. Configure backend env

```bash
cd server
cp example.env .env
```

Defaults in `server/.env`:

- `PORT=:8000`
- `DB_HOST=127.0.0.1`
- `DB_PORT=5433`
- `DB_NAME=task_tracker`
- `JWT_SECRET=CHANGE_ME`

### 3. Apply migrations

```bash
cd server
make migrate-up
```

### 4. Run backend in watch mode

```bash
cd server
make air
```

### 5. Configure frontend env

```bash
cd web
cp .env.example .env
```

Set:

```env
BACKEND_URL=http://127.0.0.1:8000
```

### 6. Run frontend

```bash
cd web
npm install
npm run dev
```

Open `http://localhost:3000`.

## Common Problems

- `connect: connection refused` from backend: Postgres is not up on `5433`.
- `permission denied` from protected APIs: missing/expired JWT header.
- Frontend cannot reach backend in Docker: ensure `BACKEND_URL` points to `http://server:8000` inside container environment.
- Let's Encrypt certificate not issued: ensure DNS A records are fully propagated for both domains and ports `80`/`443` are open on the host firewall.
