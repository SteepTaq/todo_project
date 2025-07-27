.PHONY: migrate-up migrate-down migrate-force migrate-create

migrate-up:
	@echo "Applying database migrations..."
	@go run ./cmd/migrate -action=up

migrate-up-steps:
	@go run ./cmd/migrate -action=up -steps=$(STEPS)

migrate-down:
	@echo "Reverting last migration..."
	@go run ./cmd/migrate -action=down

migrate-down-steps:
	@go run ./cmd/migrate -action=down -steps=$(STEPS)

migrate-force:
	@go run ./cmd/migrate -action=force -version=$(VERSION)

migrate-create:
	@echo "Creating new migration files..."
	@migrate create -ext sql -dir ./internal/dbservice/migrations -seq $(NAME)

migrate-version:
	@go run ./cmd/migrate -action=version