# Делаем api-gateway для гостиницы

<!-- 1. Напиши hello world и запусти
2. 
    Добавить ручки: (пример https://go.dev/doc/tutorial/web-service-gin)
    (будем исопльзовать gin)
    В main.go без БД
    1. Создать гостя `POST` `/guest` body{"name", ...} reponse 203, {"id"}
    2. Обновить `GET` `/guest/:id` body{"name", ...} reponse 200, {"id"} 
    3. Обновить `PUT` `/guest` body{"name", ...} reponse 200, {"id"} 
    4. Удалить `DELETE` `/guest/:id` reponse 200

3. В постмане отправить запросы и проверить -->

4. отличие слайса от массива, структура слайса, что происходит при append
5. маппа, бакеты, миграции, что проиходит при коллизиях

<!-- 6. переход с переменной на базу -->
<!-- 7. добавить таблицу guest -->

8. Добавить docker compose
    1. Поставить docker hub, brew install ...
    2. Добавить docker-compose.yaml
    ```yaml
    services:

    api-gateway:
    build: ./
    ports:
    - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

    postgres:
    image: postgres:14.10-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - pgdata:/var/lib/postgresql/data  
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 10s
      retries: 5
    ```

    Добавить ./Dockerfile
    ```Dockerfile
    FROM golang:1.23-alpine AS builder

    WORKDIR /

    COPY go.mod go.sum ./
    RUN go mod download

    COPY . .

    RUN go build -o /app/api-gateway ./cmd/
    RUN ls -l /app

    FROM alpine:latest 

    WORKDIR /app
    COPY --from=builder /app/api-gateway .

    ENTRYPOINT ["./api-gateway"]
    ```

    1. docker compose build
    2. docker compose up

    docker compose up --build