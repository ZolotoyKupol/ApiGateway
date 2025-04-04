FROM golang:1.23-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api-gateway ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

ENTRYPOINT [ "./api-gateway" ]
