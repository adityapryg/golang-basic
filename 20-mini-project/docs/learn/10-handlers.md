# 10 - Handler Layer

Create HTTP handlers (controllers) for user authentication, todos, and health checks.

---

## Overview

We'll create:

1. `internal/handler/user_handler.go` - User authentication endpoints
2. `internal/handler/todo_handler.go` - Todo CRUD endpoints
3. `internal/handler/health_handler.go` - Health check endpoint

Handlers handle HTTP requests, call services, and return HTTP responses.

---

## Step 1: Create User Handler

### 1.1 Create user_handler.go File

📝 **Create file:** `internal/handler/user_handler.go`

```go
package handler

import (
	"errors"
	"net/http"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	authService *service.AuthService
}

// NewUserHandler creates new user handler instance
func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Data tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	user, err := h.authService.Register(req)
	if err != nil {
		// Map service errors to HTTP status codes
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserExists) {
			statusCode = http.StatusConflict
		} else if errors.Is(err, service.ErrEmailExists) {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Registrasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses (jangan kirim password!)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Registrasi berhasil",
		Data: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Data tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	token, user, err := h.authService.Login(req)
	if err != nil {
		// Map service errors to HTTP status codes
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Login gagal",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses dengan token
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login berhasil",
		Data: dto.LoginResponse{
			Token: token,
			User: dto.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				FullName:  user.FullName,
				CreatedAt: user.CreatedAt,
			},
		},
	})
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Ambil user ID dari context (di-set oleh AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	// Panggil service
	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal mengambil profil",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profil berhasil diambil",
		Data: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		},
	})
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	var req dto.UpdateProfileRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Data tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	user, err := h.authService.UpdateUser(userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, service.ErrEmailExists) {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal memperbarui profil",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profil berhasil diperbarui",
		Data: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		},
	})
}
```

- [ ] File created at `internal/handler/user_handler.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Handler Pattern

```go
type UserHandler struct {
	authService *service.AuthService  // Dependency
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Register(c *gin.Context) {
	// Handler method
}
```

**Handler responsibilities:**

1. Parse HTTP request (JSON, query params, URL params)
2. Validate request format (Gin binding)
3. Extract context values (userID from JWT)
4. Call service methods
5. Map service errors to HTTP status codes
6. Return JSON responses

### 1.3 Understanding Request Binding

```go
var req dto.RegisterRequest

// ShouldBindJSON automatically:
// 1. Parses JSON body
// 2. Validates against struct tags
// 3. Returns error if validation fails
if err := c.ShouldBindJSON(&req); err != nil {
	c.JSON(400, dto.ErrorResponse{
		Success: false,
		Message: "Data tidak valid",
		Error:   err.Error(),
	})
	return
}
```

**What gets validated:**

```go
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}
```

**Validation errors:**

```json
// Missing required field
{
  "error": "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
}

// Invalid email format
{
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

### 1.4 Understanding Error Mapping

```go
user, err := h.authService.Register(req)
if err != nil {
	// Default to 500 Internal Server Error
	statusCode := http.StatusInternalServerError

	// Map specific errors to appropriate status codes
	if errors.Is(err, service.ErrUserExists) {
		statusCode = http.StatusConflict  // 409
	} else if errors.Is(err, service.ErrEmailExists) {
		statusCode = http.StatusConflict  // 409
	}

	c.JSON(statusCode, dto.ErrorResponse{
		Success: false,
		Message: "Registrasi gagal",
		Error:   err.Error(),
	})
	return
}
```

**HTTP status code mapping:**

| Service Error           | HTTP Status           | Code | Meaning                 |
| ----------------------- | --------------------- | ---- | ----------------------- |
| `ErrUserExists`         | Conflict              | 409  | Resource already exists |
| `ErrEmailExists`        | Conflict              | 409  | Resource already exists |
| `ErrInvalidCredentials` | Unauthorized          | 401  | Authentication failed   |
| `ErrUserNotFound`       | Not Found             | 404  | Resource not found      |
| Any other error         | Internal Server Error | 500  | Server error            |

### 1.5 Understanding Context Extraction

```go
// Get userID from context (set by AuthMiddleware)
userID, exists := c.Get("userID")
if !exists {
	c.JSON(401, dto.ErrorResponse{
		Success: false,
		Message: "User ID tidak ditemukan",
	})
	return
}

// Type assertion (interface{} to uint)
id := userID.(uint)
user, err := h.authService.GetUserByID(id)
```

**How it works:**

```
1. Request with JWT token
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

2. AuthMiddleware extracts and validates token
   claims := ValidateToken(token)
   c.Set("userID", claims.UserID)

3. Handler retrieves userID
   userID, exists := c.Get("userID")
   id := userID.(uint)
```

---

## Step 2: Create Todo Handler

### 2.1 Create todo_handler.go File

📝 **Create file:** `internal/handler/todo_handler.go`

```go
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	todoService *service.TodoService
}

// NewTodoHandler creates new todo handler instance
func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// Create handles creating new todo
func (h *TodoHandler) Create(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	var req dto.CreateTodoRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Data tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	todo, err := h.todoService.CreateTodo(userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidStatus) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal membuat todo",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Todo berhasil dibuat",
		Data:    todo,
	})
}

// GetAll handles getting all user's todos
func (h *TodoHandler) GetAll(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	// Ambil query parameters untuk filter
	status := c.Query("status")
	priority := c.Query("priority")

	// Panggil service
	todos, err := h.todoService.GetUserTodos(userID.(uint), status, priority)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal mengambil todos",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todos berhasil diambil",
		Data:    todos,
	})
}

// GetByID handles getting single todo
func (h *TodoHandler) GetByID(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	// Parse todo ID dari URL parameter
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "ID todo tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	todo, err := h.todoService.GetTodoByID(uint(todoID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal mengambil todo",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo berhasil diambil",
		Data:    todo,
	})
}

// Update handles updating todo
func (h *TodoHandler) Update(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	// Parse todo ID dari URL parameter
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "ID todo tidak valid",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateTodoRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Data tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	todo, err := h.todoService.UpdateTodo(uint(todoID), userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
		} else if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal memperbarui todo",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo berhasil diperbarui",
		Data:    todo,
	})
}

// Delete handles deleting todo
func (h *TodoHandler) Delete(c *gin.Context) {
	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "User ID tidak ditemukan",
		})
		return
	}

	// Parse todo ID dari URL parameter
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "ID todo tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Panggil service
	err = h.todoService.DeleteTodo(uint(todoID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: "Gagal menghapus todo",
			Error:   err.Error(),
		})
		return
	}

	// Response sukses (204 No Content atau 200 OK)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo berhasil dihapus",
		Data:    nil,
	})
}
```

- [ ] File created at `internal/handler/todo_handler.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Query Parameters

```go
// Get query parameters
status := c.Query("status")       // ?status=pending
priority := c.Query("priority")   // ?priority=high
```

**Example requests:**

```bash
# No filters
GET /api/v1/todos

# Filter by status
GET /api/v1/todos?status=pending

# Filter by priority
GET /api/v1/todos?priority=high

# Multiple filters
GET /api/v1/todos?status=pending&priority=high
```

**What `c.Query()` returns:**

```go
// If parameter exists: returns value
status := c.Query("status")  // "pending"

// If parameter doesn't exist: returns empty string
status := c.Query("status")  // ""

// With default value
status := c.DefaultQuery("status", "all")  // "all" if not provided
```

### 2.3 Understanding URL Parameters

```go
// Parse ID from URL parameter
todoIDStr := c.Param("id")  // Get as string
todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
if err != nil {
	c.JSON(400, dto.ErrorResponse{
		Success: false,
		Message: "ID todo tidak valid",
	})
	return
}
```

**Route definition:**

```go
router.GET("/todos/:id", handler.GetByID)
//                  ^^^
//                  URL parameter
```

**Example requests:**

```bash
GET /api/v1/todos/1      # c.Param("id") = "1"
GET /api/v1/todos/123    # c.Param("id") = "123"
GET /api/v1/todos/abc    # c.Param("id") = "abc" → ParseUint error
```

**Why parse to uint?**

```go
// String from URL
todoIDStr := c.Param("id")  // "123"

// Parse to uint64
todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
// Base 10, fit in 32 bits

// Convert to uint for service call
todo, err := h.todoService.GetTodoByID(uint(todoID), userID)
```

### 2.4 Understanding HTTP Status Codes for Todos

```go
statusCode := http.StatusInternalServerError  // 500 (default)

if errors.Is(err, service.ErrTodoNotFound) {
	statusCode = http.StatusNotFound  // 404
} else if errors.Is(err, service.ErrUnauthorizedAccess) {
	statusCode = http.StatusForbidden  // 403
} else if errors.Is(err, service.ErrInvalidStatus) {
	statusCode = http.StatusBadRequest  // 400
}
```

**Status code meanings:**

| Code | Constant                    | Meaning      | When to use                      |
| ---- | --------------------------- | ------------ | -------------------------------- |
| 200  | `StatusOK`                  | Success      | GET, PUT successful              |
| 201  | `StatusCreated`             | Created      | POST successful                  |
| 204  | `StatusNoContent`           | No content   | DELETE successful (no body)      |
| 400  | `StatusBadRequest`          | Bad request  | Invalid input, validation failed |
| 401  | `StatusUnauthorized`        | Unauthorized | Missing/invalid token            |
| 403  | `StatusForbidden`           | Forbidden    | Valid token but no permission    |
| 404  | `StatusNotFound`            | Not found    | Resource doesn't exist           |
| 409  | `StatusConflict`            | Conflict     | Duplicate username/email         |
| 500  | `StatusInternalServerError` | Server error | Unexpected errors                |

**401 vs 403:**

```
401 Unauthorized:
- No token provided
- Invalid/expired token
- "Who are you?"

403 Forbidden:
- Valid token provided
- But user doesn't own the resource
- "I know who you are, but you can't do that"
```

---

## Step 3: Create Health Handler

### 3.1 Create health_handler.go File

📝 **Create file:** `internal/handler/health_handler.go`

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates new health handler instance
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles health check endpoint
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Cek koneksi database
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   "error",
			"message":  "Database connection failed",
			"database": "disconnected",
		})
		return
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   "error",
			"message":  "Database ping failed",
			"database": "disconnected",
		})
		return
	}

	// Semua OK
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "Service is healthy",
		"database": "connected",
	})
}
```

- [ ] File created at `internal/handler/health_handler.go`
- [ ] Content copied exactly
- [ ] File saved

### 3.2 Understanding Health Check

**Why health check endpoint?**

- Monitor service availability
- Docker health checks
- Kubernetes liveness/readiness probes
- Load balancer health checks

**What to check:**

```go
// 1. Get underlying SQL DB
sqlDB, err := h.db.DB()

// 2. Ping to verify connection
err := sqlDB.Ping()

// 3. Return status
{
  "status": "ok",
  "database": "connected"
}
```

**Docker health check usage:**

```yaml
# docker-compose.yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

---

## Step 4: Verify Handlers

### 4.1 Check File Structure

```bash
$ ls -la internal/handler/    # Linux/Mac
$ dir internal\handler\       # Windows
```

**Expected output:**

```
health_handler.go
todo_handler.go
user_handler.go
```

- [ ] All three files present

### 4.2 Verify Compilation

```bash
$ go build internal/handler/*.go
```

**Expected:** No errors

- [ ] Handlers compile successfully

### 4.3 Check Handler Methods

```bash
$ grep -r "^func (h" internal/handler/    # Linux/Mac
$ findstr /B "func (h" internal\handler\* # Windows
```

**Expected methods:**

- UserHandler: Register, Login, GetProfile, UpdateProfile
- TodoHandler: Create, GetAll, GetByID, Update, Delete
- HealthHandler: HealthCheck

- [ ] All methods defined

---

## Step 5: Understanding Handler Testing

### 5.1 Handler Test Pattern

```go
func TestUserHandler_Register(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock service
	mockService := &MockAuthService{}
	handler := NewUserHandler(mockService)

	// Register route
	router.POST("/register", handler.Register)

	// 2. Create request
	reqBody := `{
		"username": "testuser",
		"email": "test@example.com",
		"password": "password123",
		"full_name": "Test User"
	}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 3. Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assert
	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}
```

### 5.2 Testing with httptest

```go
// Create test request
req := httptest.NewRequest("GET", "/todos", nil)

// Add headers
req.Header.Set("Authorization", "Bearer token123")

// Record response
w := httptest.NewRecorder()

// Execute request
router.ServeHTTP(w, req)

// Check response
assert.Equal(t, 200, w.Code)
body := w.Body.String()
assert.Contains(t, body, "success")
```

---

## Step 6: Commit Changes

### 6.1 Check Status

```bash
$ git status
```

- [ ] Handler files shown

### 6.2 Stage Handlers

```bash
$ git add internal/handler/
```

- [ ] Files staged

### 6.3 Commit

```bash
$ git commit -m "Add user, todo, and health HTTP handlers"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/handler/user_handler.go` created
- [ ] UserHandler struct defined
- [ ] NewUserHandler constructor implemented
- [ ] Register handler with request binding
- [ ] Error mapping to HTTP status codes
- [ ] Password excluded from responses
- [ ] Login handler with token response
- [ ] GetProfile handler with context extraction
- [ ] UpdateProfile handler with partial updates
- [ ] `internal/handler/todo_handler.go` created
- [ ] TodoHandler struct defined
- [ ] NewTodoHandler constructor implemented
- [ ] Create handler with validation
- [ ] GetAll handler with query parameters
- [ ] Status and priority filters
- [ ] GetByID handler with URL parameter parsing
- [ ] Update handler with ownership check
- [ ] Delete handler with proper response
- [ ] URL parameter parsing with strconv
- [ ] 401 vs 403 status code distinction
- [ ] `internal/handler/health_handler.go` created
- [ ] HealthHandler struct defined
- [ ] NewHealthHandler constructor implemented
- [ ] HealthCheck handler with database ping
- [ ] Understanding of request binding
- [ ] Understanding of context extraction
- [ ] Understanding of error mapping
- [ ] Understanding of query vs URL parameters
- [ ] Handlers compile without errors
- [ ] Changes committed to git

---

## 🐛 Common Issues

### Issue: "userID not found in context"

**Solution:** Ensure AuthMiddleware is applied to the route

### Issue: "ShouldBindJSON always fails"

**Solution:** Check JSON tags in DTO structs match request field names

### Issue: "strconv.ParseUint fails"

**Solution:** Validate ID is numeric before parsing

### Issue: "Password appears in response"

**Solution:** Use UserResponse DTO, not model.User directly

---

## 📊 Current Progress

```
[█████████████████████████████] 62% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers  
**Next:** Route configuration

---

**Previous:** [09-services.md](09-services.md)  
**Next:** [11-routes.md](11-routes.md)
