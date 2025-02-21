FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN cd cmd && go build -o ../api-gateway
RUN ls -l 

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

ENTRYPOINT [ "./api-gateway" ]
