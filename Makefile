# Переменные
DB_DRIVER=postgres
DB_SOURCE=./migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable


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
	
# Собрать приложение
build:
	go build -o api-gateway .

# Запустить приложение
run:
	go run main.go