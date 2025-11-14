# 12 - Main Application

Create the main entry point with manual dependency injection.

---

## Overview

We'll create:

1. `cmd/api/main.go` - Application entry point with DI pattern

This ties everything together: config, database, repositories, services, handlers, and routes.

---

## Step 1: Create Main Application

### 1.1 Create main.go File

📝 **Create file:** `cmd/api/main.go`

```go
package main

import (
	"log"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/route"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("Configuration loaded")

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize database connection
	db := config.NewDatabase(cfg)
	log.Println("Database connected")

	// Manual Dependency Injection Pattern
	// Layer 1: Initialize Repositories (Data Access Layer)
	userRepository := repository.NewUserRepository(db)
	todoRepository := repository.NewTodoRepository(db)
	log.Println("Repositories initialized")

	// Layer 2: Initialize Services (Business Logic Layer)
	authService := service.NewAuthService(userRepository)
	todoService := service.NewTodoService(todoRepository)
	log.Println("Services initialized")

	// Layer 3: Initialize Handlers (HTTP Layer)
	userHandler := handler.NewUserHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealthHandler(db)
	log.Println("Handlers initialized")

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())
	log.Println("Middleware applied")

	// Setup routes
	route.SetupRoutes(router, userHandler, todoHandler, healthHandler)
	log.Println("Routes configured")

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	log.Printf("Environment: %s", cfg.GinMode)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

- [ ] File created at `cmd/api/main.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Main Function Flow

```go
func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Set Gin Mode
	gin.SetMode(cfg.GinMode)

	// 3. Initialize Database
	db := config.NewDatabase(cfg)

	// 4. Manual Dependency Injection (3 layers)
	// Layer 1: Repositories
	// Layer 2: Services
	// Layer 3: Handlers

	// 5. Initialize Router
	router := gin.Default()

	// 6. Apply Middleware
	router.Use(middleware.LoggerMiddleware())

	// 7. Setup Routes
	route.SetupRoutes(router, handlers...)

	// 8. Start Server
	router.Run(":8080")
}
```

**Execution order is important:**

- Config must be loaded first
- Database needs config
- Repositories need database
- Services need repositories
- Handlers need services
- Routes need handlers

---

## Step 2: Understanding Manual Dependency Injection

### 2.1 What is Dependency Injection?

**Without DI (bad):**

```go
type UserHandler struct {}

func (h *UserHandler) Register() {
	// Create dependencies inside handler
	db := connectDB()
	repo := NewUserRepository(db)
	service := NewAuthService(repo)

	// Use service
	service.Register(...)
}
```

**Problems:**

- Tight coupling
- Hard to test
- Duplicate code
- Can't replace dependencies

**With DI (good):**

```go
type UserHandler struct {
	authService *service.AuthService  // Injected
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}
```

**Benefits:**

- Loose coupling
- Easy testing (can mock services)
- Single initialization
- Clear dependencies

### 2.2 Three-Layer Dependency Pattern

```go
// Layer 1: Repositories (depend on database)
userRepo := repository.NewUserRepository(db)
todoRepo := repository.NewTodoRepository(db)

// Layer 2: Services (depend on repositories)
authService := service.NewAuthService(userRepo)
todoService := service.NewTodoService(todoRepo)

// Layer 3: Handlers (depend on services)
userHandler := handler.NewUserHandler(authService)
todoHandler := handler.NewTodoHandler(todoService)
```

**Dependency flow:**

```
DB
 ↓
Repositories
 ↓
Services
 ↓
Handlers
 ↓
Routes
```

**Why this order?**

- Each layer only knows about the layer below
- Changes to lower layers don't affect upper layers
- Easy to test each layer independently

### 2.3 Why Manual DI (No Framework)?

**Frameworks (like Wire, Fx):**

```go
// Automatic dependency injection
container.Provide(NewDB)
container.Provide(NewUserRepository)
container.Provide(NewAuthService)
container.Invoke(StartServer)
```

**Manual DI:**

```go
// Explicit dependency creation
db := NewDatabase()
userRepo := NewUserRepository(db)
authService := NewAuthService(userRepo)
```

**Benefits of manual DI:**

- **Explicit**: See all dependencies clearly
- **Simple**: No magic, no complex configuration
- **Debuggable**: Easy to trace execution
- **No learning curve**: Just function calls
- **No framework lock-in**: Standard Go code

**When to use frameworks:**

- Very large applications (100+ dependencies)
- Complex dependency graphs
- Need lifecycle management

---

## Step 3: Understanding gin.Default() vs gin.New()

### 3.1 Difference Between Default and New

```go
// gin.Default() includes built-in middleware
router := gin.Default()
// Equivalent to:
router := gin.New()
router.Use(gin.Logger())
router.Use(gin.Recovery())

// gin.New() is clean router
router := gin.New()
// No middleware, completely bare
```

**gin.Default() includes:**

1. **Logger**: Logs all requests
2. **Recovery**: Catches panics and returns 500

**When to use each:**

| Use Case               | Choice          | Reason                   |
| ---------------------- | --------------- | ------------------------ |
| Production             | `gin.Default()` | Want panic recovery      |
| Development            | `gin.Default()` | Want request logging     |
| Testing                | `gin.New()`     | Don't want logs in tests |
| Custom middleware only | `gin.New()`     | Full control             |

**In our app:**

```go
router := gin.Default()  // Get recovery + built-in logger

// Add custom middleware
router.Use(middleware.LoggerMiddleware())  // Our logger
router.Use(middleware.CORSMiddleware())
router.Use(middleware.ErrorHandler())
```

### 3.2 Middleware Execution Order

```go
router := gin.Default()
// 1. gin.Logger()    (built-in)
// 2. gin.Recovery()  (built-in)

router.Use(middleware.LoggerMiddleware())
// 3. Our logger

router.Use(middleware.CORSMiddleware())
// 4. CORS

router.Use(middleware.ErrorHandler())
// 5. Error handler

// Then route-specific middleware
// 6. AuthMiddleware (if applied to route)

// Finally handler
// 7. Handler executes
```

---

## Step 4: Understanding Server Startup

### 4.1 Starting the Server

```go
serverAddr := ":" + cfg.ServerPort  // ":8080"
log.Printf("Server starting on %s", serverAddr)

if err := router.Run(serverAddr); err != nil {
	log.Fatalf("Failed to start server: %v", err)
}
```

**What `router.Run()` does:**

1. Creates HTTP server
2. Binds to port (e.g., 8080)
3. Starts listening for requests
4. **Blocks** until server stops

**Port formats:**

```go
router.Run(":8080")          // localhost:8080
router.Run("0.0.0.0:8080")   // All interfaces
router.Run("127.0.0.1:8080") // Localhost only
```

**Why bind to `0.0.0.0` in Docker?**

```go
// Docker needs this to accept external connections
router.Run("0.0.0.0:8080")

// This won't work in Docker (only internal)
router.Run("localhost:8080")
```

### 4.2 Graceful Shutdown

**Current implementation:**

```go
router.Run(":8080")  // Blocks forever
```

**Better implementation (with graceful shutdown):**

```go
srv := &http.Server{
	Addr:    ":8080",
	Handler: router,
}

// Start server in goroutine
go func() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("Shutting down server...")

// Graceful shutdown with 5 second timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
	log.Fatalf("Server forced to shutdown: %v", err)
}

log.Println("Server exited")
```

---

## Step 5: Verify Main Application

### 5.1 Check File Structure

```bash
$ ls -la cmd/api/    # Linux/Mac
$ dir cmd\api\       # Windows
```

**Expected output:**

```
main.go
```

- [ ] File present

### 5.2 Verify Compilation

```bash
$ go build cmd/api/main.go
```

**Expected:** Binary created

- [ ] Application compiles successfully

### 5.3 Check for Syntax Errors

```bash
$ go vet cmd/api/main.go
```

**Expected:** No issues found

- [ ] No syntax errors

---

## Step 6: Run the Application

### 6.1 First Run

```bash
$ go run cmd/api/main.go
```

**Expected output:**

```
Configuration loaded
Database connected
Repositories initialized
Services initialized
Handlers initialized
Middleware applied
Routes configured
Server starting on :8080
Environment: debug
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*HealthHandler).HealthCheck-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/register     --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*UserHandler).Register-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/login        --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*UserHandler).Login-fm (5 handlers)
[GIN-debug] GET    /api/v1/users/profile     --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*UserHandler).GetProfile-fm (6 handlers)
[GIN-debug] PUT    /api/v1/users/profile     --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*UserHandler).UpdateProfile-fm (6 handlers)
[GIN-debug] POST   /api/v1/todos             --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*TodoHandler).Create-fm (6 handlers)
[GIN-debug] GET    /api/v1/todos             --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*TodoHandler).GetAll-fm (6 handlers)
[GIN-debug] GET    /api/v1/todos/:id         --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*TodoHandler).GetByID-fm (6 handlers)
[GIN-debug] PUT    /api/v1/todos/:id         --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*TodoHandler).Update-fm (6 handlers)
[GIN-debug] DELETE /api/v1/todos/:id         --> github.com/adityapryg/golang-demo/20-mini-project/internal/handler.(*TodoHandler).Delete-fm (6 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

- [ ] Server starts without errors
- [ ] All routes listed
- [ ] Listening on port 8080

### 6.2 Test Health Check

**Open new terminal:**

```bash
$ curl http://localhost:8080/health
```

**Expected response:**

```json
{
  "status": "ok",
  "message": "Service is healthy",
  "database": "connected"
}
```

- [ ] Health check returns 200 OK
- [ ] Database shows connected

### 6.3 Stop Server

**In terminal running server:**

```
Press Ctrl+C
```

**Expected:**

```
^C
[GIN] Shutdown requested
```

- [ ] Server stops gracefully

---

## Step 7: Build Production Binary

### 7.1 Build Binary

```bash
# Linux/Mac
$ go build -o bin/api cmd/api/main.go

# Windows
$ go build -o bin/api.exe cmd/api/main.go
```

- [ ] Binary created in bin/

### 7.2 Run Binary

```bash
# Linux/Mac
$ ./bin/api

# Windows
$ bin\api.exe
```

**Expected:** Server starts

- [ ] Binary runs successfully

### 7.3 Check Binary Size

```bash
$ ls -lh bin/    # Linux/Mac
$ dir bin\       # Windows
```

**Expected:** ~20-30 MB

- [ ] Binary size reasonable

---

## Step 8: Environment Configuration

### 8.1 Create .env File (Optional)

📝 **Create file:** `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist_db

JWT_SECRET=your-super-secret-key-change-this-in-production

SERVER_PORT=8080
GIN_MODE=debug
```

- [ ] `.env` file created
- [ ] Variables match config.go

### 8.2 Using Environment Variables

```bash
# Linux/Mac
$ export DB_HOST=localhost
$ export DB_PORT=5432
$ export JWT_SECRET=my-secret-key
$ go run cmd/api/main.go

# Windows
$ set DB_HOST=localhost
$ set DB_PORT=5432
$ set JWT_SECRET=my-secret-key
$ go run cmd/api/main.go
```

### 8.3 Production Environment

```bash
# Production settings
export GIN_MODE=release
export SERVER_PORT=8080
export JWT_SECRET=very-strong-production-secret-minimum-32-chars
export DB_PASSWORD=strong-database-password

./bin/api
```

---

## Step 9: Commit Changes

### 9.1 Check Status

```bash
$ git status
```

- [ ] main.go shown

### 9.2 Stage Main File

```bash
$ git add cmd/api/main.go
```

- [ ] File staged

### 9.3 Commit

```bash
$ git commit -m "Add main application entry point with dependency injection"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `cmd/api/main.go` created
- [ ] Configuration loading implemented
- [ ] Gin mode setting
- [ ] Database initialization
- [ ] Manual dependency injection pattern
- [ ] Layer 1: Repositories initialized
- [ ] UserRepository created
- [ ] TodoRepository created
- [ ] Layer 2: Services initialized
- [ ] AuthService created
- [ ] TodoService created
- [ ] Layer 3: Handlers initialized
- [ ] UserHandler created
- [ ] TodoHandler created
- [ ] HealthHandler created
- [ ] Gin router initialization
- [ ] Global middleware applied
- [ ] Routes setup
- [ ] Server startup logic
- [ ] Error handling on startup
- [ ] Logging at each step
- [ ] Understanding of DI pattern
- [ ] Understanding of layer dependencies
- [ ] Understanding of gin.Default()
- [ ] Understanding of router.Run()
- [ ] Application compiles successfully
- [ ] Server starts and runs
- [ ] Health check endpoint works
- [ ] Binary builds successfully
- [ ] Environment variables understood
- [ ] Changes committed to git

---

## 🧪 Quick Test

### Test Complete Startup

```bash
# 1. Ensure PostgreSQL is running
$ docker-compose up -d postgres

# 2. Start application
$ go run cmd/api/main.go

# 3. In another terminal, test endpoints
$ curl http://localhost:8080/health

# 4. Test registration
$ curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# 5. Test login
$ curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

- [ ] Health check works
- [ ] Registration works
- [ ] Login works and returns token
- [ ] All endpoints accessible

---

## 🐛 Common Issues

### Issue: "Address already in use"

**Solution:**

```bash
# Find process using port 8080
lsof -i :8080        # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Kill process or change SERVER_PORT
```

### Issue: "Database connection failed"

**Solution:**

- Ensure PostgreSQL is running
- Check DB credentials in config
- Verify database exists

### Issue: "Module not found"

**Solution:**

```bash
$ go mod tidy
$ go mod download
```

### Issue: "Middleware not executing"

**Solution:** Ensure `.Use()` is called before `SetupRoutes()`

---

## 📊 Current Progress

```
[██████████████████████████████████] 75% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers, routes, main app  
**Next:** Database migrations

---

**Previous:** [11-routes.md](11-routes.md)  
**Next:** [13-database.md](13-database.md)
