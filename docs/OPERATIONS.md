# Operations Guide

## Backend Make Targets

Run from `server/`.

- `make build`: compile API binary into `bin/server`
- `make run`: build and run compiled binary
- `make air`: run backend with hot reload
- `make test`: run all Go tests
- `make swagger`: regenerate Swagger/OpenAPI files
- `make docker-up`: build/start compose stack from repository root file
- `make docker-db-up`: start only postgres service
- `make docker-down`: stop compose services

## Migrations

### Apply all migrations

```bash
cd server
make migrate-up
```

### Rollback all migrations

```bash
cd server
make migrate-down
```

### Force migration version

```bash
cd server
make migrate-force
```

Current `force` target pins version to `1` in `cmd/migrate/main.go`.

### Create a new migration template

```bash
cd server
make migration add_some_change
```

This writes SQL files into `server/cmd/migrate/migrations`.

### Regenerate Swagger docs

```bash
cd server
make swagger
```

## Test Commands

### Backend

```bash
cd server
make test
```

### Frontend

```bash
cd web
npm run build
```

## Troubleshooting

### Migrations fail with path error

Run migration commands from `server/`, because the runner expects:

- `file://cmd/migrate/migrations`

### `permission denied` on protected route

Check one of:

- missing `Authorization` header
- expired JWT (`expiredAt` claim)
- token signed with different `JWT_SECRET`

### Login succeeds but protected Nuxt API calls fail

The UI must send `Authorization` header from stored token. If token is absent/expired, re-login.
