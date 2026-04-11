# syntax=docker/dockerfile:1
FROM golang:1.23-alpine

RUN apk add --no-cache ca-certificates tzdata git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

EXPOSE 8080

CMD sh -c "go mod tidy && \
    go run migrate/migrateSchema.go && \
    swag fmt && \
    swag init && \
    go run main.go"
