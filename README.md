# Lab-Golang

REST API สำหรับฝึก Golang ด้วย **Gin + GORM + PostgreSQL** ออกแบบตามหลัก **Layered Architecture (MVC + Repository Pattern)** พร้อม Swagger, JWT Auth, CORS middleware, Unit test และ Docker

---

## ✨ Features

- CRUD Expenses (Create / Read / Update / Delete + ExpensesDetail)
- Health Check (API + Database)
- Swagger UI (auto-generated via swaggo)
- JWT Authentication middleware
- CORS middleware
- Graceful shutdown
- Unit test ครอบคลุม controllers / services / repositories

---

## 🧱 Tech Stack

| ด้าน | เทคโนโลยี |
|------|-----------|
| Language | Go 1.19 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| Database | PostgreSQL 16 |
| API Docs | [Swaggo](https://github.com/swaggo/swag) |
| Auth | JWT (`golang-jwt/jwt/v4`) |
| Testing | `testify`, `go-sqlmock` |
| Container | Docker + docker-compose |

---

## 📁 Project Structure

```
lab-golang/
├── main.go                 # Entry point + DI wiring + graceful shutdown
├── initializers/           # โหลด .env และเชื่อมต่อ DB
│   ├── loadEnv.go
│   └── connectDatabase.go
├── routers/                # Route registration
│   └── router.go
├── middleware/             # CORS, Auth, Logger
│   └── middleware.go
├── controllers/            # HTTP handlers (parse request → call service)
│   ├── expensesController.go
│   └── healthCheckController.go
├── services/               # Business logic
│   ├── interfaces/
│   └── service/
├── repositories/           # Data access (GORM)
│   ├── interfaces/
│   └── repository/
├── models/                 # GORM entities
│   └── expensesModel.go
├── dtos/                   # Request/Response objects
│   ├── request/
│   └── response/
├── helper/                 # Utility functions
├── migrate/                # DB schema migration script
├── docs/                   # Swagger auto-generated
├── test/                   # Unit tests
│   ├── controllers/
│   ├── services/
│   └── repositories/
├── Dockerfile
└── docker-compose.yml
```

### Layered Flow

```
Client → Router → Middleware → Controller → Service → Repository → Database
                                    ↓            ↑
                                   DTO        Model
```

---

## ⚙️ Environment Variables

สร้างไฟล์ `.env` ที่ root ของโปรเจค

```env
SERVER_ADDR=:8080
SERVER_NAME=localhost
SERVER_PORT=5432
DATABASE_NAME=lab-golang
USER_DB=root
PASSWORD_DB=root
DB_SSLMODE=disable
JWT_SECRET=dev-only-insecure-secret-change-me
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

---

## 🚀 Getting Started

### 1. ติดตั้ง dependencies

```bash
go mod tidy
```

### 2. สร้าง Schema ในฐานข้อมูล

```bash
go run migrate/migrateSchema.go
```

### 3. Generate Swagger docs

```bash
swag fmt
swag init
```

> ติดตั้ง swag CLI: `go install github.com/swaggo/swag/cmd/swag@latest`

### 4. รัน API

```bash
go run main.go
```

### 5. เปิด Swagger UI

```
http://localhost:8080/swagger/index.html
```

---

## 🐳 Run with Docker

รันทั้ง API + PostgreSQL พร้อมกันด้วย docker-compose

```bash
docker-compose up -d --build
```

- API: http://localhost:8080
- PostgreSQL: localhost:5432

หยุดการทำงาน

```bash
docker-compose down
```

---

## 🧪 Testing

รัน unit test ทั้งหมด

```bash
go test ./test/... -v
```

รันเฉพาะ layer

```bash
go test ./test/services/... -v
go test ./test/controllers/... -v
go test ./test/repositories/... -v
```

---

## 📌 API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Health Check

| Method | Path | Description |
|--------|------|-------------|
| GET | `/HealthCheck/Api` | ตรวจสอบว่า API ทำงานอยู่ |
| GET | `/HealthCheck/Database` | ตรวจสอบการเชื่อมต่อฐานข้อมูล |

### Expenses

| Method | Path | Description |
|--------|------|-------------|
| PUT | `/Expenses/CreateExpenses` | สร้างรายการค่าใช้จ่าย |
| PUT | `/Expenses/CreateExpensesDetail` | สร้างรายละเอียดค่าใช้จ่าย |
| GET | `/Expenses/GetListMoneyCard` | ดึงรายการค่าใช้จ่ายทั้งหมด |
| POST | `/Expenses/GetListMoneyCardDetail` | ดึงรายละเอียดของรายการ |
| PUT | `/Expenses/UpdateExpensesDetail` | แก้ไขรายละเอียดค่าใช้จ่าย |
| DELETE | `/Expenses/DeleteExpensesDetail` | ลบรายละเอียดค่าใช้จ่าย |

ดู request/response schema เพิ่มเติมได้ที่ Swagger UI

---
