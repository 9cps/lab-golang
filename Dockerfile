# syntax=docker/dockerfile:1

# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs and compile binaries
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag fmt && \
    swag init --generalInfo cmd/api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api/ && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/migrate ./cmd/migrate/

# ── Runtime stage ────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server  ./server
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/docs    ./docs

EXPOSE 8080

CMD ["./server"]
