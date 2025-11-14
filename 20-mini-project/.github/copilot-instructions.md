# Copilot Instructions - Todo REST API

## Architecture Overview

This is a **Clean Architecture** Go REST API with strict layer separation:

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Model (Entity)
```

**Dependency flow is unidirectional**: handlers depend on services, services depend on repositories. No upward dependencies.

## Project Structure

- `cmd/api/main.go` - Entry point with **manual dependency injection** pattern (no DI framework)
- `internal/handler/` - HTTP handlers (Gin controllers), extract context values, bind JSON
- `internal/service/` - Business logic, validation, and custom error definitions
- `internal/repository/` - GORM data access with query methods
- `internal/model/` - GORM entities with table name overrides
- `internal/dto/` - Request/response DTOs with binding tags (`binding:"required,oneof=..."`)
- `internal/middleware/` - Auth (JWT), CORS, error handling, logging
- `internal/config/` - Environment-based config with defaults
- `migrations/` - SQL migrations (manually applied, not auto-migrated)

## Key Patterns & Conventions

### 1. Layer Communication

**Handler → Service**:

```go
// Handler extracts userID from context (set by auth middleware)
userID, exists := c.Get("userID")
todo, err := h.todoService.CreateTodo(userID.(uint), req)
```

**Service → Repository**:

```go
// Service contains business logic and calls repository
todo := &model.Todo{Title: req.Title, UserID: userID}
if err := s.todoRepo.Create(todo); err != nil {
    return nil, err
}
```

### 2. Error Handling Convention

Services define **custom error variables** for domain errors:

```go
var (
    ErrTodoNotFound = errors.New("todo not found")
    ErrUnauthorizedAccess = errors.New("unauthorized access to todo")
)
```

Handlers check errors with `errors.Is()` to determine HTTP status codes:

```go
if errors.Is(err, service.ErrTodoNotFound) {
    statusCode = http.StatusNotFound
}
```

### 3. Authentication Flow

1. JWT middleware (`middleware.AuthMiddleware()`) validates token from `Authorization: Bearer <token>`
2. Extracts claims and sets `userID` in gin.Context: `c.Set("userID", claims.UserID)`
3. Handlers retrieve: `userID, exists := c.Get("userID")`
4. Services use userID for **ownership validation** before data operations

JWT utilities in `internal/utils/jwt.go`:

- `GenerateToken(userID, username)` - 24-hour expiration
- `ValidateToken(tokenString)` - Returns claims or error

### 4. DTO Binding & Validation

Use Gin struct tags for validation:

```go
type CreateTodoRequest struct {
    Title    string `json:"title" binding:"required,max=200"`
    Status   string `json:"status" binding:"required,oneof=pending in_progress completed"`
    Priority string `json:"priority" binding:"required,oneof=low medium high"`
}
```

Bind in handlers: `c.ShouldBindJSON(&req)`

### 5. Repository Query Patterns

Methods are **purpose-specific** (not generic):

- `FindByUserID(userID)` - User's todos ordered by `created_at DESC`
- `FindByUserIDWithFilters(userID, status, priority)` - Filtered queries
- `IsOwnedByUser(todoID, userID)` - Ownership check returning bool

### 6. Soft Deletes

Models use `gorm.DeletedAt` for soft deletes:

```go
DeletedAt gorm.DeletedAt `gorm:"index"`
```

Repository delete: `r.db.Delete(&model.Todo{}, id)` (automatically soft deletes)

## Configuration

Environment variables loaded via `config.LoadConfig()` with fallback defaults:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET` - **Critical**: Change in production (default is insecure)
- `SERVER_PORT` - Default 8080
- `GIN_MODE` - `debug` or `release`

## Development Workflow

### Build & Run

```bash
# Local development
go run cmd/api/main.go

# Build binary
go build -o bin/api.exe cmd/api/main.go

# Run binary
./bin/api.exe
```

### Docker Workflow

```bash
# Start with PostgreSQL
docker-compose up -d

# Rebuild after code changes
docker-compose up --build
```

The API depends on PostgreSQL health check (`condition: service_healthy` in docker-compose.yml).

### Testing

```bash
# Run all tests
go test ./tests/...

# Run specific test
go test ./tests/ -run TestAuthFlow
```

Tests use `testify/suite` pattern with setup/teardown. Test database connection is initialized in `SetupSuite()`.

## Common Gotchas

1. **Import paths**: Module is `github.com/adityapryg/golang-demo/20-mini-project` - internal imports start here
2. **Middleware order matters**: In `main.go`, AuthMiddleware is applied per route group (not globally)
3. **Date format**: DueDate uses `"2006-01-02"` (YYYY-MM-DD) for parsing
4. **Password hashing**: Use `utils.HashPassword()` and `utils.CheckPasswordHash()` (bcrypt)
5. **GORM auto-migration not used**: Database schema managed via `migrations/*.sql` files
6. **Context values**: Only `userID` is set by middleware (type: `uint`)

## Route Structure

- Public: `/health`, `/api/v1/auth/register`, `/api/v1/auth/login`
- Protected: All `/api/v1/users/*` and `/api/v1/todos/*` require JWT

Routes defined in `internal/route/routes.go` with clear group separation.

## Adding New Features

When adding new entities, follow the layered approach:

1. Create model in `internal/model/` with GORM tags
2. Create DTOs in `internal/dto/` with binding validation
3. Implement repository in `internal/repository/` with specific query methods
4. Implement service in `internal/service/` with business logic and custom errors
5. Implement handler in `internal/handler/` with HTTP status code mapping
6. Register routes in `internal/route/routes.go`
7. Wire dependencies in `cmd/api/main.go` following existing DI pattern

## Testing Strategy

- **Service layer**: Mock repository, test business logic in isolation
- **Handler layer**: Use `httptest` with test routes, inject mock service
- **Auth flow**: Generate real JWT tokens for integration tests using `utils.GenerateToken()`
- Test database: Initialized per suite, use actual PostgreSQL (not in-memory)

---

## Project Recreation Checklist

Follow this step-by-step guide to recreate this project from scratch:

### Phase 1: Project Initialization

- [ ] **1.1** Create project directory: `mkdir 20-mini-project && cd 20-mini-project`
- [ ] **1.2** Initialize Go module: `go mod init github.com/adityapryg/golang-demo/20-mini-project`
- [ ] **1.3** Create directory structure:
  ```bash
  mkdir -p cmd/api
  mkdir -p internal/{config,dto,handler,middleware,model,repository,route,service,utils}
  mkdir -p migrations
  mkdir -p tests
  mkdir -p bin
  mkdir -p .github
  ```

### Phase 2: Install Dependencies

- [ ] **2.1** Install Gin framework: `go get github.com/gin-gonic/gin@v1.9.1`
- [ ] **2.2** Install GORM and PostgreSQL driver:
  ```bash
  go get gorm.io/gorm@v1.25.5
  go get gorm.io/driver/postgres@v1.5.4
  ```
- [ ] **2.3** Install JWT library: `go get github.com/golang-jwt/jwt/v5@v5.2.0`
- [ ] **2.4** Install bcrypt for password hashing: `go get golang.org/x/crypto@v0.17.0`
- [ ] **2.5** Install testing library: `go get github.com/stretchr/testify@v1.8.4`
- [ ] **2.6** (Optional) Install Swagger:
  ```bash
  go get github.com/swaggo/gin-swagger@v1.6.0
  go get github.com/swaggo/files@v1.0.1
  go get github.com/swaggo/swag@v1.16.2
  ```

### Phase 3: Configuration Layer

- [ ] **3.1** Create `internal/config/config.go`:

  - Define `Config` struct with DB credentials, JWT secret, server port, Gin mode
  - Implement `LoadConfig()` function with environment variable loading
  - Implement `getEnv(key, defaultValue)` helper function

- [ ] **3.2** Create `internal/config/database.go`:
  - Implement `NewDatabase(cfg *Config)` function
  - Build PostgreSQL DSN string with `fmt.Sprintf()`
  - Connect using `gorm.Open(postgres.Open(dsn), &gorm.Config{})`
  - Add auto-migration for models: `db.AutoMigrate(&model.User{}, &model.Todo{})`
  - Return `*gorm.DB` instance

### Phase 4: Model Layer (Entities)

- [ ] **4.1** Create `internal/model/user.go`:

  - Define `User` struct with fields: `ID`, `Username`, `Email`, `Password`, `FullName`, `Todos`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
  - Add GORM tags: `gorm:"primaryKey"`, `gorm:"unique;not null"`, etc.
  - Add `DeletedAt gorm.DeletedAt` for soft delete
  - Override table name with `TableName()` method returning `"users"`

- [ ] **4.2** Create `internal/model/todo.go`:
  - Define `Todo` struct with fields: `ID`, `Title`, `Description`, `Status`, `Priority`, `DueDate`, `UserID`, `User`, timestamps
  - Add GORM tags with sizes, indexes, defaults
  - Use `*time.Time` for nullable `DueDate`
  - Add foreign key: `gorm:"foreignKey:UserID"`
  - Override table name to `"todos"`

### Phase 5: DTO Layer (Data Transfer Objects)

- [ ] **5.1** Create `internal/dto/user_dto.go`:

  - Define `RegisterRequest` with `Username`, `Email`, `Password`, `FullName` (add `binding:"required"` tags)
  - Define `LoginRequest` with `Username`, `Password`
  - Define `UpdateProfileRequest` with optional fields using pointers
  - Define `UserResponse` with safe user data (no password)
  - Define `LoginResponse` with `Token` and `User`

- [ ] **5.2** Create `internal/dto/todo_dto.go`:

  - Define `CreateTodoRequest` with validation: `binding:"required,max=200"`, `binding:"oneof=pending in_progress completed"`
  - Define `UpdateTodoRequest` with optional fields using pointers
  - Define `TodoResponse` for API responses
  - Define `TodoQueryParams` for filtering (optional)

- [ ] **5.3** Add common response DTOs in `user_dto.go` or separate file:
  - Define `SuccessResponse` struct with `Success bool`, `Message string`, `Data interface{}`
  - Define `ErrorResponse` struct with `Success bool`, `Message string`, `Error string`

### Phase 6: Utility Functions

- [ ] **6.1** Create `internal/utils/password.go`:

  - Implement `HashPassword(password string) (string, error)` using `bcrypt.GenerateFromPassword()`
  - Implement `CheckPassword(password, hash string) bool` using `bcrypt.CompareHashAndPassword()`

- [ ] **6.2** Create `internal/utils/jwt.go`:
  - Define `Claims` struct with `UserID uint`, `Username string`, and `jwt.RegisteredClaims`
  - Implement `GenerateToken(userID uint, username string) (string, error)`:
    - Set expiration to 24 hours
    - Create claims with `ExpiresAt` and `IssuedAt`
    - Sign token using `jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`
    - Return signed string
  - Implement `ValidateToken(tokenString string) (*Claims, error)`:
    - Parse token with `jwt.ParseWithClaims()`
    - Verify signature using JWT secret from config
    - Return claims if valid

### Phase 7: Middleware Layer

- [ ] **7.1** Create `internal/middleware/auth.go`:

  - Implement `AuthMiddleware() gin.HandlerFunc`:
    - Extract token from `Authorization: Bearer <token>` header
    - Validate format (must have "Bearer" prefix)
    - Call `utils.ValidateToken(tokenString)`
    - Set `userID` in context: `c.Set("userID", claims.UserID)`
    - Return 401 error if validation fails

- [ ] **7.2** Create `internal/middleware/logger.go`:

  - Implement `LoggerMiddleware() gin.HandlerFunc`:
    - Record start time with `time.Now()`
    - Call `c.Next()` to process request
    - Calculate duration with `time.Since(startTime)`
    - Log method, path, IP, status code, and duration

- [ ] **7.3** Create `internal/middleware/error.go`:
  - Implement `ErrorHandler() gin.HandlerFunc` to catch errors from `c.Errors`
  - Implement `CORSMiddleware() gin.HandlerFunc`:
    - Set `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, etc.
    - Handle OPTIONS preflight requests with `c.AbortWithStatus(204)`

### Phase 8: Repository Layer (Data Access)

- [ ] **8.1** Create `internal/repository/user_repository.go`:

  - Define `UserRepository` struct with `db *gorm.DB`
  - Implement constructor: `NewUserRepository(db *gorm.DB) *UserRepository`
  - Implement methods:
    - `Create(user *model.User) error`
    - `FindByID(id uint) (*model.User, error)`
    - `FindByUsername(username string) (*model.User, error)`
    - `FindByEmail(email string) (*model.User, error)`
    - `Update(user *model.User) error`
    - `ExistsByUsername(username string) (bool, error)`
    - `ExistsByEmail(email string) (bool, error)`

- [ ] **8.2** Create `internal/repository/todo_repository.go`:
  - Define `TodoRepository` struct with `db *gorm.DB`
  - Implement constructor: `NewTodoRepository(db *gorm.DB) *TodoRepository`
  - Implement methods:
    - `Create(todo *model.Todo) error`
    - `FindByID(id uint) (*model.Todo, error)`
    - `FindByUserID(userID uint) ([]model.Todo, error)` with `Order("created_at DESC")`
    - `FindByUserIDWithFilters(userID uint, status, priority string) ([]model.Todo, error)`
    - `Update(todo *model.Todo) error`
    - `Delete(id uint) error` (soft delete with `db.Delete()`)
    - `IsOwnedByUser(todoID, userID uint) (bool, error)`

### Phase 9: Service Layer (Business Logic)

- [ ] **9.1** Create `internal/service/auth_service.go`:

  - Define custom errors: `ErrUserExists`, `ErrInvalidCredentials`, `ErrEmailExists`
  - Define `AuthService` struct with `userRepo *repository.UserRepository`
  - Implement constructor: `NewAuthService(userRepo *repository.UserRepository) *AuthService`
  - Implement `Register(req dto.RegisterRequest) (*model.User, error)`:
    - Check if username/email exists
    - Hash password with `utils.HashPassword()`
    - Create user via repository
  - Implement `Login(req dto.LoginRequest) (string, *model.User, error)`:
    - Find user by username
    - Verify password with `utils.CheckPassword()`
    - Generate JWT token with `utils.GenerateToken()`
  - Implement `GetUserByID(id uint) (*model.User, error)`
  - Implement `UpdateUser(id uint, req dto.UpdateProfileRequest) (*model.User, error)`:
    - Check email uniqueness if changed
    - Update only provided fields

- [ ] **9.2** Create `internal/service/todo_service.go`:
  - Define custom errors: `ErrTodoNotFound`, `ErrUnauthorizedAccess`, `ErrInvalidStatus`, `ErrInvalidPriority`
  - Define `TodoService` struct with `todoRepo *repository.TodoRepository`
  - Implement constructor: `NewTodoService(todoRepo *repository.TodoRepository) *TodoService`
  - Implement `CreateTodo(userID uint, req dto.CreateTodoRequest) (*model.Todo, error)`:
    - Validate status and priority
    - Parse due date if provided (format: `"2006-01-02"`)
    - Create todo via repository
  - Implement `GetTodoByID(todoID, userID uint) (*model.Todo, error)`:
    - Check ownership with `IsOwnedByUser()`
    - Return `ErrUnauthorizedAccess` if not owner
  - Implement `GetUserTodos(userID uint, status, priority string) ([]model.Todo, error)`
  - Implement `UpdateTodo(todoID, userID uint, req dto.UpdateTodoRequest) (*model.Todo, error)`:
    - Verify ownership first
    - Update only non-nil fields
  - Implement `DeleteTodo(todoID, userID uint) error`:
    - Verify ownership before deletion
  - Add helper functions: `isValidStatus(status string) bool`, `isValidPriority(priority string) bool`

### Phase 10: Handler Layer (HTTP Controllers)

- [ ] **10.1** Create `internal/handler/user_handler.go`:

  - Define `UserHandler` struct with `authService *service.AuthService`
  - Implement constructor: `NewUserHandler(authService *service.AuthService) *UserHandler`
  - Implement `Register(c *gin.Context)`:
    - Bind JSON with `c.ShouldBindJSON(&req)`
    - Call `authService.Register()`
    - Return 201 with user data (exclude password)
  - Implement `Login(c *gin.Context)`:
    - Bind credentials
    - Call `authService.Login()`
    - Return token and user data
  - Implement `GetProfile(c *gin.Context)`:
    - Extract `userID` from context: `c.Get("userID")`
    - Call `authService.GetUserByID()`
  - Implement `UpdateProfile(c *gin.Context)`:
    - Extract userID from context
    - Bind update request
    - Call `authService.UpdateUser()`

- [ ] **10.2** Create `internal/handler/todo_handler.go`:

  - Define `TodoHandler` struct with `todoService *service.TodoService`
  - Implement constructor: `NewTodoHandler(todoService *service.TodoService) *TodoHandler`
  - Implement `Create(c *gin.Context)`:
    - Extract userID from context
    - Bind request and call `todoService.CreateTodo()`
    - Map service errors to HTTP status codes with `errors.Is()`
  - Implement `GetAll(c *gin.Context)`:
    - Extract userID and query params (status, priority filters)
    - Call `todoService.GetUserTodos()`
  - Implement `GetByID(c *gin.Context)`:
    - Parse `:id` param with `strconv.ParseUint()`
    - Verify ownership in service
  - Implement `Update(c *gin.Context)`:
    - Extract ID, userID, and bind update request
    - Call `todoService.UpdateTodo()`
  - Implement `Delete(c *gin.Context)`:
    - Extract ID and userID
    - Call `todoService.DeleteTodo()`

- [ ] **10.3** Create `internal/handler/health_handler.go`:
  - Define `HealthHandler` struct with `db *gorm.DB`
  - Implement `HealthCheck(c *gin.Context)`:
    - Ping database with `db.Raw("SELECT 1")`
    - Return health status and DB connection status

### Phase 11: Route Configuration

- [ ] **11.1** Create `internal/route/routes.go`:
  - Implement `SetupRoutes(router *gin.Engine, userHandler, healthHandler, todoHandler)`:
    - Add health check route: `router.GET("/health", healthHandler.HealthCheck)`
    - Create API v1 group: `v1 := router.Group("/api/v1")`
    - Add auth routes (public): `/auth/register`, `/auth/login`
    - Add user routes (protected with `AuthMiddleware()`): `/users/profile` (GET, PUT)
    - Add todo routes (protected): `/todos` (POST, GET), `/todos/:id` (GET, PUT, DELETE)

### Phase 12: Main Application

- [ ] **12.1** Create `cmd/api/main.go`:
  - Import all necessary packages
  - Implement `main()` function:
    - Load config with `config.LoadConfig()`
    - Set Gin mode with `gin.SetMode(cfg.GinMode)`
    - Initialize database with `config.NewDatabase(cfg)`
    - **Manual Dependency Injection Pattern**:
      - Layer 1: Initialize repositories (UserRepository, TodoRepository)
      - Layer 2: Initialize services (AuthService, TodoService)
      - Layer 3: Initialize handlers (UserHandler, HealthHandler, TodoHandler)
    - Create Gin router: `router := gin.Default()`
    - Apply global middleware: `LoggerMiddleware()`, `CORSMiddleware()`, `ErrorHandler()`
    - Setup routes with `route.SetupRoutes()`
    - Start server: `router.Run(":" + cfg.ServerPort)`

### Phase 13: Database Migrations

- [ ] **13.1** Create `migrations/001_create_tables.sql`:

  - Create `users` table with columns: id, username, email, password, full_name, timestamps
  - Add unique constraints on username and email
  - Add indexes: idx_users_username, idx_users_email, idx_users_deleted_at
  - Create `todos` table with columns: id, title, description, status, priority, due_date, user_id, timestamps
  - Add foreign key constraint to users table with CASCADE delete
  - Add indexes: idx_todos_user_id, idx_todos_status, idx_todos_deleted_at
  - Add check constraints for status (`IN ('pending', 'in_progress', 'completed')`)
  - Add check constraints for priority (0-5 range)
  - Add table and column comments

- [ ] **13.2** Apply migrations manually:
  ```bash
  psql -U postgres -d todolist_db -f migrations/001_create_tables.sql
  ```

### Phase 14: Docker Configuration

- [ ] **14.1** Create `Dockerfile`:

  - Multi-stage build with `golang:1.22-alpine` as builder
  - Copy go.mod and go.sum, run `go mod download`
  - Copy source code and build: `go build -o main .`
  - Use `alpine:latest` for runtime
  - Copy binary from builder stage
  - Expose port 8080
  - Set CMD to run binary

- [ ] **14.2** Create `docker-compose.yml`:

  - Define `postgres` service:
    - Image: `postgres:15-alpine`
    - Environment: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB
    - Port: 5432
    - Volume: postgres_data
    - Healthcheck: `pg_isready -U postgres`
  - Define `api` service:
    - Build from Dockerfile
    - Environment: DB_HOST=postgres, DB credentials, JWT_SECRET, GIN_MODE
    - Port: 8080
    - Depends on postgres with `condition: service_healthy`
    - Restart policy: unless-stopped

- [ ] **14.3** Create `.dockerignore`:
  - Add: bin/, .git/, .env, \*.md, tests/, .gitignore

### Phase 15: Testing Setup

- [ ] **15.1** Create `tests/auth_test.go`:

  - Define `AuthTestSuite` struct embedding `suite.Suite`
  - Implement `SetupSuite()` to initialize database and router
  - Implement test methods:
    - `TestRegister()` - Test user registration
    - `TestLoginSuccess()` - Test successful login
    - `TestLoginInvalidCredentials()` - Test failed login
  - Use `httptest.NewRecorder()` and `httptest.NewRequest()`
  - Run suite with `suite.Run(t, new(AuthTestSuite))`

- [ ] **15.2** Create `tests/todo_test.go`:
  - Define `TodoTestSuite` with router, token, and userID
  - Create test user and JWT token in `SetupSuite()`
  - Implement test methods:
    - `TestCreateTodo()` - Test todo creation
    - `TestGetAllTodos()` - Test listing todos
    - `TestGetTodoByID()` - Test getting single todo
    - `TestUpdateTodo()` - Test updating todo
    - `TestDeleteTodo()` - Test soft delete
  - Add `Authorization: Bearer <token>` header to requests
  - Clean up in `TearDownTest()`

### Phase 16: Additional Files

- [ ] **16.1** Create `.gitignore`:

  - Add: bin/, .env, _.exe, _.log, vendor/, coverage.out

- [ ] **16.2** Create `.env.example`:

  ```
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=postgres
  DB_NAME=todolist_db
  JWT_SECRET=your-super-secret-key-change-this-in-production
  SERVER_PORT=8080
  GIN_MODE=debug
  ```

- [ ] **16.3** Create `README.md` with:

  - Project description and features
  - Architecture diagram
  - Setup instructions
  - API endpoints documentation
  - Docker commands
  - Testing commands

- [ ] **16.4** Create `.github/copilot-instructions.md` (this file!)

### Phase 17: Verification & Testing

- [ ] **17.1** Verify project structure matches expected layout
- [ ] **17.2** Run `go mod tidy` to clean up dependencies
- [ ] **17.3** Start PostgreSQL: `docker-compose up -d postgres`
- [ ] **17.4** Apply migrations manually (if not using auto-migrate)
- [ ] **17.5** Run application: `go run cmd/api/main.go`
- [ ] **17.6** Test health endpoint: `curl http://localhost:8080/health`
- [ ] **17.7** Test registration:
  ```bash
  curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","email":"test@example.com","password":"password123","full_name":"Test User"}'
  ```
- [ ] **17.8** Test login and save token
- [ ] **17.9** Test protected endpoints with token
- [ ] **17.10** Run test suite: `go test ./tests/... -v`
- [ ] **17.11** Build Docker image: `docker-compose build`
- [ ] **17.12** Run full stack: `docker-compose up`

### Phase 18: Final Checklist

- [ ] **18.1** All imports use correct module path: `github.com/adityapryg/golang-demo/20-mini-project`
- [ ] **18.2** All handlers properly extract `userID` from context
- [ ] **18.3** All services define custom errors as package-level variables
- [ ] **18.4** All repositories return specific errors (not generic ones)
- [ ] **18.5** All DTOs have proper validation tags
- [ ] **18.6** All routes are properly grouped (public vs protected)
- [ ] **18.7** Middleware is applied in correct order
- [ ] **18.8** Database connection includes timezone setting
- [ ] **18.9** JWT secret is loaded from environment (not hardcoded)
- [ ] **18.10** Password is hashed before storage, never exposed in responses

---

## Quick Start After Recreation

```bash
# 1. Start infrastructure
docker-compose up -d postgres

# 2. Install dependencies
go mod download

# 3. Run application
go run cmd/api/main.go

# 4. Run tests
go test ./tests/... -v

# 5. Build and run with Docker
docker-compose up --build
```
