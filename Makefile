# Define the migration directory once to avoid repetition
MIGRATION_DIR := cmd/migrate/migrations

# Ensure 'migration' is marked as a phony target
.PHONY: migration

clean:
	@rm -rf bin

build:
	@go build -o bin/server cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/server

migration:
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-force:
	@go run cmd/migrate/main.go force

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down