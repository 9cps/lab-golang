# Backend Golang Best Practices Skill

## โครงสร้างโปรเจกต์ตามมาตรฐาน [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

### 🏗️ 1. Standard Go Project Layout (แนะนำสำหรับโปรเจกต์ใหญ่)
```
myapp/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── handler/                 # HTTP handlers (controller)
│   │   ├── user_handler.go
│   │   └── product_handler.go
│   ├── service/                 # Business logic
│   │   ├── interfaces/
│   │   │   └── user_iservice.go # Interface definition
│   │   └── user_service.go      # Implementation
│   ├── repository/              # Database layer
│   │   ├── interfaces/
│   │   │   └── user_irepo.go    # Interface definition
│   │   └── user_repo.go         # Implementation
│   ├── model/                   # Structs / Domain models
│   │   ├── user.go
│   │   └── product.go
│   ├── dtos/                    # Data Transfer Objects
│   │   ├── request/
│   │   │   └── user_request.go
│   │   └── response/
│   │       └── user_response.go
│   ├── middleware/              # Auth, logging, CORS ฯลฯ
│   │   └── auth.go
│   └── router/                  # Route definitions
│       └── router.go
├── pkg/                         # Reusable packages (ใช้ได้จากภายนอก)
│   ├── logger/
│   └── validator/
├── config/                      # Config structs & loaders
│   └── config.go
├── migrations/                  # SQL migration files
├── docs/                        # Swagger / API docs
├── test/                        # Unit / Integration tests
├── .env
├── .env.example
├── go.mod
├── go.sum
└── Makefile
```

### 🧩 2. Domain-Driven Structure (แนะนำสำหรับ Microservice / DDD)
```
myapp/
├── cmd/
│   └── main.go
├── internal/
│   ├── user/                    # Domain: User
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── dto.go               # Request/Response DTOs
│   ├── product/                 # Domain: Product
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── dto.go
│   └── middleware/
├── config/
├── pkg/
├── docs/
├── go.mod
└── go.sum
```

---

## 📐 Layer Architecture (หลักการสำคัญ)

```
Request → Router → Middleware → Handler → Service → Repository → Database
                                   ↑           ↑           ↑
                                  DTO        Interface   Interface
```

| Layer | หน้าที่ |
|---|---|
| **Router** | กำหนด route และ group ของ API |
| **Middleware** | Auth, logging, CORS, rate limiting |
| **Handler** | รับ request, bind & validate DTO, ส่ง response |
| **Service** | Business logic ทั้งหมด, ห้าม query DB ตรง |
| **Repository** | Query database เท่านั้น, ห้ามมี business logic |
| **Model** | Data structures / entities ที่ map กับ DB |
| **DTO** | Request/Response objects ที่ใช้รับ-ส่งข้อมูลกับ client |

---

## Best Practices สำหรับ Backend Golang

### 1. Interface-Based Design (สำคัญมากสำหรับ Testability)

ทุก Service และ Repository ต้องมี interface คู่กันเสมอ เพื่อให้ mock ได้ในการทดสอบ

```go
// interfaces/user_iservice.go
type UserService interface {
    GetUserByID(ctx context.Context, id uint) (*model.User, error)
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
}

// service/user_service.go
type userService struct {
    userRepo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
    return &userService{userRepo: repo}
}
```

### 2. DTOs — แยก Model ออกจาก Request/Response

ห้าม expose `model.User` ตรงกับ client ใช้ DTO แทนเสมอ

```go
// dtos/request/user_request.go
type CreateUserRequest struct {
    Name  string `json:"name"  binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// dtos/response/user_response.go
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 3. Error Handling — ตรวจสอบและส่ง error พร้อม context เสมอ

```go
// ❌ Bad
user, _ := repo.FindByID(ctx, id)

// ✅ Good
user, err := repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("user %d not found: %w", id, err)
    }
    return nil, fmt.Errorf("get user by id: %w", err)
}
```

### 4. Context Propagation — ส่ง `context.Context` เป็น argument แรกเสมอ

```go
// ✅ ทุก function ที่เรียก DB หรือ external service ต้องรับ ctx
func (r *userRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

### 5. Dependency Injection ใน main.go

```go
// cmd/server/main.go
func main() {
    db := config.InitDB()

    // Wire up dependencies
    userRepo := repository.NewUserRepository(db)
    userSvc  := service.NewUserService(userRepo)
    userHdl  := handler.NewUserHandler(userSvc)

    r := router.NewRouter(userHdl)
    r.Run(":8080")
}
```

### 6. Naming Conventions

| สิ่งที่ตั้งชื่อ | รูปแบบ | ตัวอย่าง |
|---|---|---|
| Interface | `<Name>Service` / `<Name>Repository` | `UserService`, `UserRepository` |
| Implementation | `<name>Service` / `<name>Repository` | `userService`, `userRepository` |
| Constructor | `New<Name>` | `NewUserService` |
| Handler file | `<name>_handler.go` | `user_handler.go` |
| Test file | `<name>_test.go` | `user_service_test.go` |

### 7. Configuration — แยก config ออกจากโค้ดด้วย `.env`

```go
// config/config.go
type Config struct {
    DBHost string `mapstructure:"DB_HOST"`
    DBPort string `mapstructure:"DB_PORT"`
    DBName string `mapstructure:"DB_NAME"`
}

func Load() (*Config, error) {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    var cfg Config
    return &cfg, viper.Unmarshal(&cfg)
}
```

### 8. แนวทางปฏิบัติอื่น ๆ

| หัวข้อ | แนวทาง |
|---|---|
| **Logging** | ใช้ structured logging (zerolog/zap), ห้าม log password หรือ token |
| **API Design** | RESTful, ใช้ HTTP status code ให้ถูกต้อง (200/201/400/404/500) |
| **Testing** | Table-driven tests, mock ด้วย interface, แยก unit / integration test |
| **Security** | Validate/sanitize input ทุก request, จัดการ secrets ผ่าน env |
| **Performance** | ใช้ goroutine อย่างระมัดระวัง, profile ก่อน optimize |
| **Documentation** | ใช้ Swagger/OpenAPI, comment public functions |

---

> อ้างอิง: https://github.com/golang-standards/project-layout
