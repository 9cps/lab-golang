# Lab-Golang

A REST API built with **Gin + GORM + PostgreSQL** following **Layered Architecture + Microservice** design, featuring Swagger, JWT Auth, CORS middleware, unit tests, and Docker.

---

## ✨ Features

- CRUD Expenses (Create / Read / Update / Delete + ExpensesDetail)
- Health Check (API + Database)
- Swagger UI (`/swagger/`)
- JWT Authentication middleware (api-gateway)
- CORS middleware (api-gateway)
- Reverse Proxy (api-gateway → expenses-service)
- Graceful shutdown
- Unit tests covering handler / service / repository layers

---

## 🧱 Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) v1.9 |
| ORM | [GORM](https://gorm.io/) v1.25 |
| Database | PostgreSQL 16 |
| API Docs | [Swaggo](https://github.com/swaggo/swag) v1.16 |
| Auth | JWT HS256 (`golang-jwt/jwt/v5`) |
| Testing | `testify`, `go-sqlmock` |
| Container | Docker + docker-compose |

---

## 📁 Project Structure

```
lab-golang/
├── services/
│   ├── api-gateway/                  # Reverse proxy + Auth/CORS middleware
│   │   ├── cmd/
│   │   │   └── main.go
│   │   └── internal/
│   │       ├── config/
│   │       │   └── env.go
│   │       ├── middleware/
│   │       │   └── middleware.go     # CorsMiddleware, AuthMiddleware
│   │       ├── proxy/
│   │       │   └── proxy.go          # httputil.ReverseProxy wrapper
│   │       └── router/
│   │           └── router.go
│   └── expenses-service/             # Business logic microservice
│       ├── cmd/
│       │   ├── api/
│       │   │   └── main.go           # DI wiring + graceful shutdown
│       │   └── migrate/
│       │       └── main.go           # AutoMigrate runner
│       ├── internal/
│       │   ├── config/
│       │   │   ├── database.go       # PostgreSQL connection + pool
│       │   │   └── env.go
│       │   ├── dtos/
│       │   │   ├── request/
│       │   │   │   └── expenses_request.go
│       │   │   └── response/
│       │   │       ├── default_response.go
│       │   │       └── expenses_response.go
│       │   ├── handler/
│       │   │   ├── expenses_handler.go
│       │   │   └── health_check_handler.go
│       │   ├── model/
│       │   │   └── expenses.go
│       │   ├── repository/
│       │   │   ├── interfaces/
│       │   │   │   ├── expenses_irepo.go
│       │   │   │   └── health_check_irepo.go
│       │   │   ├── expenses_repository.go
│       │   │   └── health_check_repository.go
│       │   ├── router/
│       │   │   └── router.go
│       │   └── service/
│       │       ├── interfaces/
│       │       │   ├── expenses_iservice.go
│       │       │   └── health_check_iservice.go
│       │       ├── expenses_service.go
│       │       └── health_check_service.go
│       └── test/
│           ├── handler/
│           │   └── expenses_handler_test.go
│           ├── repository/
│           │   └── expenses_repository_test.go
│           └── service/
│               └── expenses_service_test.go
├── docs/                             # Swagger (pre-generated)
├── docker-compose.yml
├── Dockerfile.gateway
├── Dockerfile.expenses
├── go.mod
└── .env.example
```

### Layer Flow

```
Client → api-gateway (CORS + JWT) → expenses-service
                                         ↓
                               Handler (bind DTO)
                                         ↓
                               Service (business logic)
                                         ↓
                               Repository (GORM queries)
                                         ↓
                                    PostgreSQL
```

### Naming Conventions

| Symbol | Pattern | Example |
|---|---|---|
| Interface | `<Name>Service` / `<Name>Repository` / `<Name>Handler` | `ExpensesService`, `ExpensesRepository`, `ExpensesHandler` |
| Implementation struct | `<name>Service` / `<name>Repository` / `<name>Handler` | `expensesService`, `expensesRepository`, `expensesHandler` |
| Constructor | `New<Name>` | `NewExpensesService`, `NewExpensesRepository`, `NewExpensesHandler` |
| File | `<name>_<layer>.go` (snake_case) | `expenses_handler.go`, `expenses_service.go`, `expenses_repository.go` |

---

## ⚙️ Environment Variables

Create a `.env` file at the project root based on `.env.example`.

```env
# Database
SERVER_NAME=localhost
SERVER_PORT=5432
DATABASE_NAME=lab-golang
USER_DB=root
PASSWORD_DB=root
DB_SSLMODE=disable

# Services
SERVER_ADDR=:8080
EXPENSES_SERVICE_URL=http://localhost:8081

# Auth
JWT_SECRET=dev-only-insecure-secret-change-me

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

---

## 🚀 Getting Started (Local)

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Run database migration

```bash
go run ./services/expenses-service/cmd/migrate/
```

### 3. Generate Swagger docs

> Skip this step if `docs/` is already up to date.

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init \
  --generalInfo services/expenses-service/cmd/api/main.go \
  --dir services/expenses-service \
  --output docs \
  --parseInternal
```

### 4. Run services in separate terminals

```bash
# Terminal 1 — expenses-service
go run ./services/expenses-service/cmd/api/

# Terminal 2 — api-gateway
go run ./services/api-gateway/cmd/
```

### 5. Open Swagger UI

```
http://localhost:8080/swagger/
```

---

## 🐳 Run with Docker

```bash
docker compose up -d
```

| Service | URL |
|---------|-----|
| API Gateway | http://localhost:8080 |
| Expenses Service | http://localhost:8081 |
| Swagger UI | http://localhost:8080/swagger/ |
| PostgreSQL | localhost:5432 |

Stop all services:

```bash
docker compose down
```

---

## 🧪 Testing

Run all unit tests:

```bash
go test ./services/expenses-service/test/... -v
```

Run by layer:

```bash
go test ./services/expenses-service/test/handler/...    -v
go test ./services/expenses-service/test/service/...    -v
go test ./services/expenses-service/test/repository/... -v
```

---

## 📌 API Endpoints

Base URL: `http://localhost:8080/api/v1`

> Expenses endpoints require the `Authorization: Bearer <token>` header.

### Health Check

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | ❌ | Check if the API is running |
| GET | `/health/database` | ❌ | Check database connectivity |

### Expenses

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/expenses` | ✅ | Create an expense record |
| GET | `/expenses` | ✅ | List all expense cards |
| POST | `/expenses/details` | ✅ | Create an expense detail |
| GET | `/expenses/details?id=<id>` | ✅ | Get details by expense ID |
| PUT | `/expenses/details/:id` | ✅ | Update an expense detail |
| DELETE | `/expenses/details/:id` | ✅ | Delete an expense detail |

For full request/response schemas, see the Swagger UI: `http://localhost:8080/swagger/`

---
