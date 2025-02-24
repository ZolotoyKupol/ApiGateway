# Переменные
DB_DRIVER=postgres
DB_SOURCE=./migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

# Цели (targets)
.PHONY: help migrate up down create clean build run

# Помощь
help:
	@echo "Available commands:"
	@echo "  make migrate   - Apply all migrations"
	@echo "  make up        - Apply pending migrations"
	@echo "  make down      - Rollback the last migration"
	@echo "  make create    - Create a new migration file"
	@echo "  make clean     - Remove all migration files"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"

# Применить все миграции
migrate:
	goose -dir $(DB_SOURCE) $(DB_DRIVER) "$(DATABASE_URL)" up

# Применить только новые миграции
up:
	goose -dir $(DB_SOURCE) $(DB_DRIVER) "$(DATABASE_URL)" up

# Откатить последнюю миграцию
down:
	goose -dir $(DB_SOURCE) $(DB_DRIVER) "$(DATABASE_URL)" down

# Создать новую миграцию
create:
ifndef name
	$(error "Variable 'name' is not set. Usage: make create name=<migration_name>")
endif
	goose -dir $(DB_SOURCE) create $(name) sql

# Очистить все миграции
clean:
	rm -rf $(DB_SOURCE)/*.sql

# Собрать приложение
build:
	go build -o api-gateway .

# Запустить приложение
run:
	go run main.go