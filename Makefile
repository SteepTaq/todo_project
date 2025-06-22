.PHONY: migrate-up migrate-down migrate-force migrate-create

# Основная команда для применения миграций
migrate-up:
	@echo "Applying database migrations..."
	@go run ./cmd/migrate -action=up

# Применить определённое количество миграций
migrate-up-steps:
	@go run ./cmd/migrate -action=up -steps=$(STEPS)

# Откатить последнюю миграцию
migrate-down:
	@echo "Reverting last migration..."
	@go run ./cmd/migrate -action=down

# Откатить несколько миграций
migrate-down-steps:
	@go run ./cmd/migrate -action=down -steps=$(STEPS)

# Принудительно установить версию миграции
migrate-force:
	@go run ./cmd/migrate -action=force -version=$(VERSION)

# Создать новую миграцию
migrate-create:
	@echo "Creating new migration files..."
	@migrate create -ext sql -dir ./internal/dbservice/migrations -seq $(NAME)

# Проверить текущее состояние миграций
migrate-version:
	@go run ./cmd/migrate -action=version