FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY migrations/ ./migrations/
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/api-gateway ./cmd
RUN ls -l 

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

ENTRYPOINT [ "./api-gateway" ]
