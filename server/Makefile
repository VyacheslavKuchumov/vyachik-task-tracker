# Define the migration directory once to avoid repetition
MIGRATION_DIR := cmd/migrate/migrations

# Ensure 'migration' is marked as a phony target
.PHONY: migration swagger

clean:
	@rm -rf bin

build:
	@go build -o bin/server cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/server

air:
	@air

docker-up:
	@docker compose -f ../docker-compose.yml up -d --build

docker-db-up:
	@docker compose -f ../docker-compose.yml up -d postgres

docker-down:
	@docker compose -f ../docker-compose.yml down

migration:
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-force:
	@go run cmd/migrate/main.go force

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

swagger:
	@go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d cmd,service/user,service/tracker,types,utils -o docs --parseInternal
