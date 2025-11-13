# 20. Mini Project - Todo List REST API

## Deskripsi

Mini project ini adalah implementasi lengkap REST API untuk aplikasi Todo List menggunakan Golang dengan **Clean Architecture**. Project ini mengintegrasikan semua konsep yang telah dipelajari dari lesson 1-19, termasuk Gin framework, GORM, JWT authentication, testing, dan dokumentasi Swagger.

Project ini mengadopsi struktur yang terinspirasi dari production-ready Go projects dengan pemisahan layer yang jelas: **Repository → Service → Controller**, serta penggunaan DTO (Data Transfer Objects) untuk memisahkan representasi data internal dan eksternal.

## Fitur

### ✅ Autentikasi & Autorisasi

- Register user baru dengan password hashing (bcrypt)
- Login dengan JWT token
- Middleware autentikasi untuk protected routes
- Token expiration (24 jam)

### 📝 User Management

- Get profil user
- Update profil user (email, full name)
- Validasi email unik

### ✔️ Todo Management (CRUD)

- Create todo dengan judul, deskripsi, status, priority, dan due date
- Read semua todos milik user (dengan filter status)
- Read detail todo by ID
- Update todo
- Delete todo (soft delete)
- Ownership validation (user hanya bisa akses todo miliknya)

### 🏥 Health Check

- Endpoint health check untuk monitoring
- Database connection check

### 📚 Dokumentasi API

- Swagger/OpenAPI documentation
- Interactive API testing via Swagger UI

### 🐳 Docker Deployment

- Dockerfile untuk containerization
- Docker Compose dengan PostgreSQL

### 🧪 Testing

- Unit tests untuk autentikasi
- Unit tests untuk todo management
- Test suite dengan testify

## Keuntungan Clean Architecture

### ✅ Separation of Concerns

- Setiap layer punya tanggung jawab yang jelas
- Mudah dipahami dan dimaintain
- Perubahan di satu layer tidak affect layer lain

### ✅ Testability

- Service layer bisa ditest tanpa HTTP
- Repository bisa ditest dengan mock database
- Handler bisa ditest dengan mock service

### ✅ Flexibility

- Mudah ganti framework (Gin → Fiber)
- Mudah ganti database (PostgreSQL → MySQL)
- Business logic tetap sama

### ✅ Reusability

- Service bisa dipanggil dari HTTP handler, CLI, atau job scheduler
- Repository bisa dipakai ulang di berbagai service

### ✅ Scalability

- Setiap layer bisa di-scale independent
- Mudah add fitur baru tanpa ubah kode existing

## Contoh Implementasi Layer

### Model (Entity)

```go
// internal/model/user.go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### DTO (Data Transfer Object)

```go
// internal/dto/user_dto.go
type UserRegisterRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Repository (Data Access)

```go
// internal/repository/user_repository.go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}
```

### Service (Business Logic)

```go
// internal/service/auth_service.go
type AuthService struct {
    userRepo *repository.UserRepository
}

func (s *AuthService) Register(req dto.UserRegisterRequest) (*dto.UserResponse, error) {
    // Business logic: check if username exists
    existing, _ := s.userRepo.FindByUsername(req.Username)
    if existing != nil {
        return nil, errcode.ErrUserExists
    }

    // Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    // Create user entity
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
    }

    // Save to database
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Convert to response DTO
    return &dto.UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

### Handler (HTTP Controller)

```go
// internal/handler/user_handler.go
type UserHandler struct {
    authService *service.AuthService
}

func (h *UserHandler) Register(c *gin.Context) {
    var req dto.UserRegisterRequest

    // Parse and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.ErrorResponse{Message: "Invalid input"})
        return
    }

    // Call service
    user, err := h.authService.Register(req)
    if err != nil {
        c.JSON(400, dto.ErrorResponse{Message: err.Error()})
        return
    }

    // Return success response
    c.JSON(201, dto.SuccessResponse{
        Message: "User registered successfully",
        Data:    user,
    })
}
```

## Tech Stack

- **Framework**: Gin (HTTP router)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Architecture**: Clean Architecture (Repository-Service-Handler pattern)
- **Authentication**: JWT (golang-jwt)
- **Password**: bcrypt (golang.org/x/crypto)
- **Documentation**: Swaggo
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose

## Arsitektur Project

### Clean Architecture Layers

```
┌──────────────────────────────────────┐
│         HTTP Request (Client)        │
└────────────────┬─────────────────────┘
                 │
┌────────────────▼─────────────────────┐
│      Controller/Handler Layer        │  ◄── Terima HTTP request, validasi input
│   (handlers/user_handler.go)         │      Return HTTP response
└────────────────┬─────────────────────┘
                 │
┌────────────────▼─────────────────────┐
│         Service Layer                │  ◄── Business logic, orchestration
│   (services/user_service.go)         │      Transaction management
└────────────────┬─────────────────────┘
                 │
┌────────────────▼─────────────────────┐
│       Repository Layer               │  ◄── Data access, CRUD operations
│  (repositories/user_repository.go)   │      Database queries
└────────────────┬─────────────────────┘
                 │
┌────────────────▼─────────────────────┐
│         Database (PostgreSQL)        │
└──────────────────────────────────────┘
```

### Struktur Folder

```
20-mini-project/
├── cmd/
│   └── api/
│       └── main.go             # Entry point aplikasi
├── internal/
│   ├── config/
│   │   ├── config.go           # Konfigurasi aplikasi
│   │   └── database.go         # Setup database & migrasi
│   ├── dto/
│   │   ├── user_dto.go         # Data Transfer Objects untuk User
│   │   ├── todo_dto.go         # Data Transfer Objects untuk Todo
│   │   └── response.go         # Standard API response format
│   ├── handler/
│   │   ├── user_handler.go     # HTTP handlers untuk User
│   │   ├── todo_handler.go     # HTTP handlers untuk Todo
│   │   └── health_handler.go   # Health check handler
│   ├── middleware/
│   │   ├── auth.go             # JWT authentication middleware
│   │   ├── logger.go           # Logging middleware
│   │   └── error.go            # Error handling & CORS
│   ├── model/
│   │   ├── user.go             # Entity User (database model)
│   │   └── todo.go             # Entity Todo (database model)
│   ├── repository/
│   │   ├── user_repository.go  # User data access layer
│   │   └── todo_repository.go  # Todo data access layer
│   ├── service/
│   │   ├── auth_service.go     # Authentication business logic
│   │   ├── user_service.go     # User business logic
│   │   └── todo_service.go     # Todo business logic
│   ├── utils/
│   │   ├── errcode/
│   │   │   └── errcode.go      # Custom error codes
│   │   ├── jwt.go              # JWT utilities
│   │   └── password.go         # Password hashing utilities
│   └── route/
│       └── route.go            # Route configuration
├── tests/
│   ├── integration/
│   │   ├── auth_test.go        # Integration test autentikasi
│   │   └── todo_test.go        # Integration test todo
│   └── unit/
│       ├── service_test.go     # Unit test service layer
│       └── repository_test.go  # Unit test repository layer
├── go.mod                      # Dependencies
├── .env.example                # Template environment variables
├── Dockerfile                  # Docker build instructions
├── docker-compose.yml          # Docker Compose configuration
└── README.md                   # Dokumentasi (file ini)
```

### Penjelasan Layer

#### 1. **Model Layer** (`internal/model/`)

- Entity database murni dengan GORM tags
- Representasi tabel database
- Tidak ada business logic
- **Contoh**: `User`, `Todo`

#### 2. **DTO Layer** (`internal/dto/`)

- Data Transfer Objects untuk komunikasi dengan client
- Memisahkan representasi internal dan eksternal
- Validasi input dengan struct tags
- **Contoh**: `UserRegisterRequest`, `TodoResponse`

#### 3. **Repository Layer** (`internal/repository/`)

- Data access layer (DAL)
- CRUD operations dengan database
- Hanya interaksi dengan database, tanpa business logic
- Return model entities
- **Contoh**: `CreateUser()`, `FindTodosByUserID()`

#### 4. **Service Layer** (`internal/service/`)

- Business logic aplikasi
- Orchestration antar repositories
- Transaction management
- Error handling business rules
- **Contoh**: `RegisterUser()`, `CreateTodoForUser()`

#### 5. **Handler Layer** (`internal/handler/`)

- HTTP request/response handling
- Parsing request body
- Call service layer
- Return HTTP response dengan status code
- **Contoh**: `HandleRegister()`, `HandleCreateTodo()`

### Alur Data Request

```
1. Client → HTTP Request (POST /api/users/register)
             ↓
2. Handler → Parse & validate DTO (UserRegisterRequest)
             ↓
3. Handler → Call Service (authService.Register())
             ↓
4. Service → Business logic (check duplicate, hash password)
             ↓
5. Service → Call Repository (userRepo.Create())
             ↓
6. Repository → Execute SQL query (INSERT INTO users...)
             ↓
7. Repository → Return model entity (User)
             ↓
8. Service → Convert model to DTO (UserResponse)
             ↓
9. Service → Return DTO to Handler
             ↓
10. Handler → Return HTTP response (201 Created + JSON)
```

## Instalasi & Setup

### Prasyarat

- Go 1.22 atau lebih baru
- PostgreSQL 15+
- Docker & Docker Compose (opsional, untuk deployment)

### 1. Install Dependencies

```bash
cd 20-mini-project
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto
go get -u github.com/stretchr/testify
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 2. Install Swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 3. Setup Database PostgreSQL

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE todolist_db;

# Keluar
\q
```

### 4. Setup Environment Variables

```bash
# Copy file .env.example
cp .env.example .env

# Edit .env sesuai konfigurasi database Anda
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=todolist_db
# JWT_SECRET=your-super-secret-key-change-this-in-production
# SERVER_PORT=8080
# GIN_MODE=debug
```

### 5. Generate Swagger Documentation

```bash
swag init
```

Ini akan generate folder `docs/` dengan file swagger.

### 6. Jalankan Aplikasi

```bash
go run main.go
```

Aplikasi akan berjalan di `http://localhost:8080`

Swagger UI tersedia di `http://localhost:8080/swagger/index.html`

## API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Health Check

| Method | Endpoint  | Deskripsi            |
| ------ | --------- | -------------------- |
| GET    | `/health` | Cek kesehatan sistem |

### Authentication (Public)

| Method | Endpoint         | Deskripsi                    |
| ------ | ---------------- | ---------------------------- |
| POST   | `/auth/register` | Register user baru           |
| POST   | `/auth/login`    | Login dan dapatkan JWT token |

### Users (Protected)

| Method | Endpoint         | Deskripsi          | Auth |
| ------ | ---------------- | ------------------ | ---- |
| GET    | `/users/profile` | Get profil user    | ✅   |
| PUT    | `/users/profile` | Update profil user | ✅   |

### Todos (Protected)

| Method | Endpoint     | Deskripsi                                 | Auth |
| ------ | ------------ | ----------------------------------------- | ---- |
| POST   | `/todos`     | Buat todo baru                            | ✅   |
| GET    | `/todos`     | Get semua todos (filter: ?status=pending) | ✅   |
| GET    | `/todos/:id` | Get detail todo                           | ✅   |
| PUT    | `/todos/:id` | Update todo                               | ✅   |
| DELETE | `/todos/:id` | Hapus todo                                | ✅   |

## Contoh Penggunaan API

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### 3. Create Todo (dengan JWT token)

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Belajar Golang",
    "description": "Menyelesaikan lesson 1-20",
    "status": "in_progress",
    "priority": 3,
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Todo berhasil dibuat",
  "data": {
    "id": 1,
    "title": "Belajar Golang",
    "description": "Menyelesaikan lesson 1-20",
    "status": "in_progress",
    "priority": 3,
    "due_date": "2024-12-31T23:59:59Z",
    "user_id": 1,
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

### 4. Get All Todos

```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Dengan filter status
curl -X GET "http://localhost:8080/api/v1/todos?status=pending" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Update Todo

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Belajar Golang - Updated",
    "status": "completed"
  }'
```

### 6. Delete Todo

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Format Response API

Semua endpoint menggunakan format response yang konsisten:

### Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": {
    /* data object atau array */
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Pesan error",
  "errors": "Detail error"
}
```

## Testing

### Menjalankan Tests

```bash
# Run all tests
go test ./tests/... -v

# Run specific test
go test ./tests/auth_test.go -v
go test ./tests/todo_test.go -v

# Run dengan coverage
go test ./tests/... -v -cover
```

### Test Coverage

Project ini mencakup:

- ✅ Test autentikasi (register, login, duplicate user, wrong password)
- ✅ Test todo CRUD operations
- ✅ Test authorization (ownership validation)
- ✅ Test middleware (auth middleware)

## Deployment dengan Docker

### Build & Run dengan Docker Compose

```bash
# Build dan run semua services (API + PostgreSQL)
docker-compose up -d

# Lihat logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop dan hapus volumes (database akan dihapus)
docker-compose down -v
```

Aplikasi akan tersedia di `http://localhost:8080`

### Build Docker Image Manual

```bash
# Build image
docker build -t todolist-api:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=todolist_db \
  -e JWT_SECRET=your-secret \
  --name todolist-api \
  todolist-api:latest
```

## Security Best Practices

### ⚠️ PENTING untuk Production:

1. **JWT Secret Key**: Ganti `JWT_SECRET` di `.env` dengan key yang kuat dan random

   ```bash
   # Generate random secret
   openssl rand -base64 32
   ```

2. **Database Password**: Gunakan password yang kuat untuk PostgreSQL

3. **Environment Variables**: Jangan commit file `.env` ke git (sudah ada di `.gitignore`)

4. **HTTPS**: Gunakan HTTPS di production (bisa dengan reverse proxy seperti Nginx)

5. **Rate Limiting**: Implementasi rate limiting untuk API endpoints

6. **Input Validation**: Semua input sudah divalidasi dengan Gin binding

7. **SQL Injection**: GORM secara otomatis mencegah SQL injection

8. **CORS**: Configure CORS sesuai kebutuhan di `middleware/error.go`

## Troubleshooting

### Database Connection Error

```
Gagal koneksi ke database: dial tcp [::1]:5432: connect: connection refused
```

**Solusi:**

- Pastikan PostgreSQL berjalan: `sudo systemctl status postgresql` (Linux) atau cek Services (Windows)
- Cek kredensial database di `.env`
- Cek firewall tidak memblok port 5432

### Port Already in Use

```
listen tcp :8080: bind: address already in use
```

**Solusi:**

- Ganti `SERVER_PORT` di `.env` ke port lain (contoh: 8081)
- Atau stop aplikasi lain yang menggunakan port 8080

### Swagger Docs Not Found

```
404 page not found
```

**Solusi:**

- Pastikan sudah run `swag init` untuk generate docs
- Cek folder `docs/` ada dan berisi file swagger

### JWT Token Invalid

```
Token tidak valid atau expired
```

**Solusi:**

- Token berlaku 24 jam, login ulang untuk dapatkan token baru
- Pastikan format Authorization header benar: `Bearer <token>`

## Perbandingan dengan Lesson Sebelumnya

### Lesson 16 (CRUD dengan Database) vs Lesson 20 (Clean Architecture)

| Aspek           | Lesson 16                                    | Lesson 20                                    |
| --------------- | -------------------------------------------- | -------------------------------------------- |
| **Struktur**    | Flat (handlers langsung akses DB)            | Layered (Handler → Service → Repository)     |
| **Handler**     | Berisi business logic + DB query             | Hanya HTTP handling                          |
| **Database**    | Langsung dari handler (`config.DB.Create()`) | Melalui repository pattern                   |
| **Testability** | Sulit (perlu mock DB di handler test)        | Mudah (tiap layer independent)               |
| **Reusability** | Terbatas                                     | Business logic bisa dipanggil dari mana saja |
| **Scalability** | Sulit ditambah fitur baru                    | Mudah extend dengan pattern yang konsisten   |

### Contoh Code Comparison

**Lesson 16 (Handlers langsung ke DB):**

```go
func CreateProduct(c *gin.Context) {
    var product Product
    c.BindJSON(&product)

    // Handler berisi business logic DAN database query
    if product.CategoryID == 0 {
        c.JSON(400, gin.H{"error": "category required"})
        return
    }

    config.DB.Create(&product)  // ❌ Direct DB access
    c.JSON(201, product)
}
```

**Lesson 20 (Clean Architecture):**

```go
// Handler
func (h *TodoHandler) CreateTodo(c *gin.Context) {
    var req dto.TodoCreateRequest
    c.ShouldBindJSON(&req)

    userID := middleware.GetUserID(c)
    todo, err := h.todoService.Create(userID, req)  // ✅ Call service
    if err != nil {
        c.JSON(400, dto.ErrorResponse{Message: err.Error()})
        return
    }
    c.JSON(201, todo)
}

// Service (business logic)
func (s *TodoService) Create(userID uint, req dto.TodoCreateRequest) (*dto.TodoResponse, error) {
    // Validasi business rules
    if req.Priority > 5 {
        return nil, errcode.ErrInvalidPriority
    }

    todo := &model.Todo{
        Title:  req.Title,
        UserID: userID,
    }

    if err := s.todoRepo.Create(todo); err != nil {  // ✅ Call repository
        return nil, err
    }

    return s.toDTO(todo), nil
}

// Repository (data access)
func (r *TodoRepository) Create(todo *model.Todo) error {
    return r.db.Create(todo).Error  // ✅ Only database operation
}
```

## Pembelajaran yang Diterapkan

Mini project ini mengintegrasikan konsep dari lesson sebelumnya dengan tambahan **Clean Architecture**:

### Dari Lesson 01-19:

- **01-10**: Fundamentals (variabel, functions, structs, error handling)
- **11**: REST API dengan net/http
- **12**: Gin framework dan routing
- **13**: Middleware (auth, logging, error handling, CORS)
- **14**: CRUD API tanpa database
- **15**: Database GORM dan PostgreSQL
- **16**: CRUD dengan database
- **17**: JWT authentication dan bcrypt
- **18**: Testing dengan testify
- **19**: Swagger documentation

### Konsep Baru di Lesson 20:

- ✨ **Clean Architecture** (layered architecture)
- ✨ **Repository Pattern** (data access abstraction)
- ✨ **Service Layer** (business logic separation)
- ✨ **DTO Pattern** (data transfer objects)
- ✨ **Dependency Injection** (loose coupling)
- ✨ **Custom Error Codes** (centralized error handling)
- ✨ **Unit Testing per Layer** (better testability)

### Konsep Software Engineering:

- **Separation of Concerns** - Tiap layer fokus pada tanggung jawabnya
- **Dependency Inversion** - High-level modules tidak depend on low-level
- **Interface Segregation** - Repository menggunakan interfaces
- **Single Responsibility** - Satu struct satu purpose
- **Open/Closed Principle** - Open for extension, closed for modification

## Pengembangan Lebih Lanjut

Ide untuk pengembangan project ini:

- [ ] Implementasi refresh token
- [ ] Email verification saat register
- [ ] Forgot password functionality
- [ ] Role-based access control (admin/user)
- [ ] Todo categories/tags
- [ ] Todo sharing antar user
- [ ] Real-time notifications (WebSocket)
- [ ] File attachment untuk todo
- [ ] Activity logs
- [ ] Dashboard analytics

## Lisensi

Materi ini dibuat untuk keperluan pembelajaran internal.

## Kontributor

- Tutorial dibuat sebagai bagian dari Golang Backend Training

---

**Happy Coding! 🚀**
