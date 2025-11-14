# 11 - Route Configuration

Configure HTTP routes and group endpoints with middleware.

---

## Overview

We'll create:

1. `internal/route/routes.go` - Route definitions and middleware application

Routes connect HTTP endpoints to handler methods and apply middleware.

---

## Step 1: Create Routes File

### 1.1 Create routes.go File

📝 **Create file:** `internal/route/routes.go`

```go
package route

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	todoHandler *handler.TodoHandler,
	healthHandler *handler.HealthHandler,
) {
	// Health check endpoint (no auth required)
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public - no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// User routes (protected - authentication required)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
		}

		// Todo routes (protected - authentication required)
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware())
		{
			todos.POST("", todoHandler.Create)
			todos.GET("", todoHandler.GetAll)
			todos.GET("/:id", todoHandler.GetByID)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}
	}
}
```

- [ ] File created at `internal/route/routes.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Route Groups

```go
// API v1 group - all routes start with /api/v1
v1 := router.Group("/api/v1")
{
	// Routes defined here
}
```

**What route grouping does:**

```
Before grouping:
router.POST("/api/v1/auth/register", handler)
router.POST("/api/v1/auth/login", handler)
router.GET("/api/v1/users/profile", handler)

After grouping:
v1 := router.Group("/api/v1")
auth := v1.Group("/auth")
auth.POST("/register", handler)  // Becomes /api/v1/auth/register
auth.POST("/login", handler)     // Becomes /api/v1/auth/login
```

**Benefits:**

- DRY (Don't Repeat Yourself)
- Easy versioning (v1, v2, v3)
- Apply middleware to specific groups
- Clear API structure

### 1.3 Understanding Public vs Protected Routes

```go
// Public routes (no authentication)
auth := v1.Group("/auth")
{
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)
}

// Protected routes (authentication required)
users := v1.Group("/users")
users.Use(middleware.AuthMiddleware())  // Apply middleware to group
{
	users.GET("/profile", userHandler.GetProfile)
	users.PUT("/profile", userHandler.UpdateProfile)
}
```

**Request flow:**

**Public route:**

```
POST /api/v1/auth/register
  ↓
UserHandler.Register
  ↓
Response
```

**Protected route:**

```
GET /api/v1/users/profile
Authorization: Bearer <token>
  ↓
AuthMiddleware (validates token)
  ↓
UserHandler.GetProfile (if token valid)
  ↓
Response
```

### 1.4 Understanding Route Parameters

```go
todos.GET("/:id", todoHandler.GetByID)
//         ^^^^
//         URL parameter

// Example requests:
// GET /api/v1/todos/1   → c.Param("id") = "1"
// GET /api/v1/todos/123 → c.Param("id") = "123"
```

**Parameter vs query string:**

```go
// URL parameter (part of path)
GET /api/v1/todos/:id
GET /api/v1/todos/123

// Query string (after ?)
GET /api/v1/todos?status=pending&priority=high
```

### 1.5 Understanding HTTP Methods

```go
todos.POST("", todoHandler.Create)      // Create new todo
todos.GET("", todoHandler.GetAll)       // Get all todos
todos.GET("/:id", todoHandler.GetByID)  // Get single todo
todos.PUT("/:id", todoHandler.Update)   // Update todo
todos.DELETE("/:id", todoHandler.Delete) // Delete todo
```

**RESTful conventions:**

| Method | Path         | Handler | Action         | Response                |
| ------ | ------------ | ------- | -------------- | ----------------------- |
| POST   | `/todos`     | Create  | Create new     | 201 Created             |
| GET    | `/todos`     | GetAll  | List all       | 200 OK                  |
| GET    | `/todos/:id` | GetByID | Get one        | 200 OK                  |
| PUT    | `/todos/:id` | Update  | Update         | 200 OK                  |
| PATCH  | `/todos/:id` | Update  | Partial update | 200 OK                  |
| DELETE | `/todos/:id` | Delete  | Delete         | 200 OK / 204 No Content |

**Why empty string `""`?**

```go
todos := v1.Group("/todos")
todos.POST("", handler)  // /api/v1/todos
todos.GET("", handler)   // /api/v1/todos

// Same as:
todos.POST("/", handler)
```

---

## Step 2: Complete API Endpoint Reference

### 2.1 All Available Endpoints

```
Health Check:
  GET /health

Authentication (Public):
  POST /api/v1/auth/register
  POST /api/v1/auth/login

User Profile (Protected):
  GET  /api/v1/users/profile
  PUT  /api/v1/users/profile

Todos (Protected):
  POST   /api/v1/todos
  GET    /api/v1/todos
  GET    /api/v1/todos/:id
  PUT    /api/v1/todos/:id
  DELETE /api/v1/todos/:id
```

### 2.2 Request/Response Examples

#### Register User

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}

# Response 201
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-03-15T10:30:00Z"
  }
}
```

#### Login

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john",
  "password": "password123"
}

# Response 200
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "full_name": "John Doe",
      "created_at": "2024-03-15T10:30:00Z"
    }
  }
}
```

#### Get Profile (Protected)

```bash
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Response 200
{
  "success": true,
  "message": "Profil berhasil diambil",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-03-15T10:30:00Z"
  }
}
```

#### Create Todo (Protected)

```bash
POST /api/v1/todos
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending",
  "priority": "high",
  "due_date": "2024-03-20"
}

# Response 201
{
  "success": true,
  "message": "Todo berhasil dibuat",
  "data": {
    "id": 1,
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "status": "pending",
    "priority": "high",
    "due_date": "2024-03-20T00:00:00Z",
    "user_id": 1,
    "created_at": "2024-03-15T10:35:00Z",
    "updated_at": "2024-03-15T10:35:00Z"
  }
}
```

#### Get All Todos with Filters (Protected)

```bash
# All todos
GET /api/v1/todos
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Filter by status
GET /api/v1/todos?status=pending
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Filter by priority
GET /api/v1/todos?priority=high
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Multiple filters
GET /api/v1/todos?status=pending&priority=high
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Response 200
{
  "success": true,
  "message": "Todos berhasil diambil",
  "data": [
    {
      "id": 1,
      "title": "Buy groceries",
      "status": "pending",
      "priority": "high",
      "created_at": "2024-03-15T10:35:00Z"
    }
  ]
}
```

#### Update Todo (Protected)

```bash
PUT /api/v1/todos/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "status": "completed"
}

# Response 200
{
  "success": true,
  "message": "Todo berhasil diperbarui",
  "data": {
    "id": 1,
    "title": "Buy groceries",
    "status": "completed",
    "priority": "high",
    "updated_at": "2024-03-15T11:00:00Z"
  }
}
```

#### Delete Todo (Protected)

```bash
DELETE /api/v1/todos/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Response 200
{
  "success": true,
  "message": "Todo berhasil dihapus",
  "data": null
}
```

---

## Step 3: Understanding Middleware Application

### 3.1 Global Middleware

```go
// Applied to ALL routes
router.Use(middleware.LoggerMiddleware())
router.Use(middleware.CORSMiddleware())
```

**When to use:**

- Logging all requests
- CORS for all endpoints
- Error handling

### 3.2 Group Middleware

```go
// Applied to specific group only
users := v1.Group("/users")
users.Use(middleware.AuthMiddleware())
```

**When to use:**

- Authentication for protected routes
- Rate limiting for API routes
- Specific validation

### 3.3 Route-Specific Middleware

```go
// Applied to single route only
router.GET("/admin", middleware.AdminMiddleware(), handler.Admin)
```

**When to use:**

- Special permissions (admin only)
- Route-specific validation
- Custom rate limits

### 3.4 Middleware Order Visualization

```
Request
  ↓
Global Middleware 1 (Logger)
  ↓
Global Middleware 2 (CORS)
  ↓
Group Middleware (Auth)
  ↓
Handler
  ↓
Response
```

---

## Step 4: Advanced Route Patterns

### 4.1 Route Versioning

```go
// API v1
v1 := router.Group("/api/v1")
{
	v1.GET("/todos", todoHandler.GetAll)
}

// API v2 (when you need breaking changes)
v2 := router.Group("/api/v2")
{
	v2.GET("/todos", todoHandlerV2.GetAll)
}
```

**Both versions can coexist:**

```
GET /api/v1/todos  → Old version
GET /api/v2/todos  → New version
```

### 4.2 Nested Resources

```go
// Get user's todos
router.GET("/users/:userId/todos", handler.GetUserTodos)

// In handler:
userID := c.Param("userId")
todos := service.GetTodosByUser(userID)
```

### 4.3 Multiple Parameters

```go
// Get specific comment on todo
router.GET("/todos/:todoId/comments/:commentId", handler.GetComment)

// In handler:
todoID := c.Param("todoId")
commentID := c.Param("commentId")
```

---

## Step 5: Verify Routes

### 5.1 Check File Structure

```bash
$ ls -la internal/route/    # Linux/Mac
$ dir internal\route\       # Windows
```

**Expected output:**

```
routes.go
```

- [ ] File present

### 5.2 Verify Compilation

```bash
$ go build internal/route/*.go
```

**Expected:** No errors

- [ ] Routes compile successfully

### 5.3 List Route Groups

```bash
$ grep -r "Group(" internal/route/    # Linux/Mac
$ findstr "Group(" internal\route\*   # Windows
```

**Expected groups:**

- `/api/v1`
- `/api/v1/auth`
- `/api/v1/users`
- `/api/v1/todos`

- [ ] All groups defined

---

## Step 6: Testing Routes

### 6.1 Route Visualization

Once the app is running, you can see all routes:

```go
// Add to main.go temporarily for debugging
routes := router.Routes()
for _, route := range routes {
	log.Printf("%s %s", route.Method, route.Path)
}
```

**Expected output:**

```
GET /health
POST /api/v1/auth/register
POST /api/v1/auth/login
GET /api/v1/users/profile
PUT /api/v1/users/profile
POST /api/v1/todos
GET /api/v1/todos
GET /api/v1/todos/:id
PUT /api/v1/todos/:id
DELETE /api/v1/todos/:id
```

### 6.2 Testing Public Routes

```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"pass123","full_name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"pass123"}'
```

### 6.3 Testing Protected Routes

```bash
# Save token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get profile
curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN"

# Create todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Todo","status":"pending","priority":"high"}'

# Get todos
curl http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN"
```

---

## Step 7: Commit Changes

### 7.1 Check Status

```bash
$ git status
```

- [ ] Route file shown

### 7.2 Stage Routes

```bash
$ git add internal/route/
```

- [ ] File staged

### 7.3 Commit

```bash
$ git commit -m "Add route configuration with API v1 endpoints"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/route/routes.go` created
- [ ] SetupRoutes function defined
- [ ] Health check route configured
- [ ] API v1 group created
- [ ] Auth routes group (public)
- [ ] Register endpoint
- [ ] Login endpoint
- [ ] Users routes group (protected)
- [ ] AuthMiddleware applied to users group
- [ ] Get profile endpoint
- [ ] Update profile endpoint
- [ ] Todos routes group (protected)
- [ ] AuthMiddleware applied to todos group
- [ ] Create todo endpoint
- [ ] Get all todos endpoint
- [ ] Get todo by ID endpoint
- [ ] Update todo endpoint
- [ ] Delete todo endpoint
- [ ] Understanding of route groups
- [ ] Understanding of public vs protected routes
- [ ] Understanding of URL parameters
- [ ] Understanding of HTTP methods
- [ ] Understanding of RESTful conventions
- [ ] Understanding of middleware application
- [ ] Routes compile without errors
- [ ] Changes committed to git

---

## 🧪 Quick Test

We can verify route structure:

📝 **Create temporary file:** `test_routes.go`

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/route"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Create dummy handlers (won't work without real services)
	userHandler := &handler.UserHandler{}
	todoHandler := &handler.TodoHandler{}
	healthHandler := &handler.HealthHandler{}

	// Setup routes
	route.SetupRoutes(router, userHandler, todoHandler, healthHandler)

	// List all routes
	fmt.Println("📍 Registered Routes:")
	fmt.Println("====================")
	routes := router.Routes()
	for _, r := range routes {
		fmt.Printf("%-7s %s\n", r.Method, r.Path)
	}

	fmt.Println("\n✅ Route configuration loaded successfully!")
	fmt.Printf("📊 Total routes: %d\n", len(routes))
}
```

### Run Test

```bash
$ go run test_routes.go
```

**Expected output:**

```
📍 Registered Routes:
====================
GET     /health
POST    /api/v1/auth/register
POST    /api/v1/auth/login
GET     /api/v1/users/profile
PUT     /api/v1/users/profile
POST    /api/v1/todos
GET     /api/v1/todos
GET     /api/v1/todos/:id
PUT     /api/v1/todos/:id
DELETE  /api/v1/todos/:id

✅ Route configuration loaded successfully!
📊 Total routes: 10
```

- [ ] Test runs without errors
- [ ] All routes listed
- [ ] Total routes = 10

### Clean Up

```bash
$ rm test_routes.go    # Linux/Mac
$ del test_routes.go   # Windows
```

---

## 🐛 Common Issues

### Issue: "Route not found 404"

**Solution:** Check route path matches exactly (case-sensitive, slashes)

### Issue: "Middleware not executing"

**Solution:** Ensure `.Use()` is called before route definitions

### Issue: "Auth required for public routes"

**Solution:** Apply AuthMiddleware to specific groups, not globally

### Issue: "URL parameter always empty"

**Solution:** Route definition must have `:paramName`, use `c.Param("paramName")`

---

## 📊 Current Progress

```
[████████████████████████████████] 68% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers, routes  
**Next:** Main application with dependency injection

---

**Previous:** [10-handlers.md](10-handlers.md)  
**Next:** [12-main-app.md](12-main-app.md)
