# 07 - Middleware

Create HTTP middleware for authentication, logging, CORS, and error handling.

---

## Overview

We'll create:

1. `internal/middleware/auth.go` - JWT authentication middleware
2. `internal/middleware/logger.go` - Request logging middleware
3. `internal/middleware/error.go` - Error handling and CORS middleware

Middleware intercepts HTTP requests/responses for cross-cutting concerns.

---

## Step 1: Create Authentication Middleware

### 1.1 Create auth.go File

📝 **Create file:** `internal/middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware memvalidasi JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Format token tidak valid",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Token tidak valid atau expired",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Simpan user ID ke context untuk digunakan di handler
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}
```

- [ ] File created at `internal/middleware/auth.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Auth Middleware Flow

```
1. Request arrives: GET /api/v1/todos
   Headers: Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

2. AuthMiddleware() executes
3. Extract "Authorization" header
4. Validate format: "Bearer <token>"
5. Extract token string
6. Validate token with utils.ValidateToken()
7. If valid:
   - Extract userID from claims
   - Store in context: c.Set("userID", claims.UserID)
   - Call c.Next() - proceed to handler
8. If invalid:
   - Return 401 Unauthorized
   - Call c.Abort() - stop execution
```

### 1.3 Context Storage Explained

```go
// In middleware - SET value
c.Set("userID", claims.UserID)

// In handler - GET value
userID, exists := c.Get("userID")
if exists {
	id := userID.(uint)  // Type assertion
	// Use id...
}
```

**Why use context?**

- Share data between middleware and handlers
- Type-safe within request scope
- Automatically cleaned up after request

### 1.4 c.Abort() vs c.Next()

```go
c.Abort()  // Stop execution, don't call next handlers
c.Next()   // Continue to next middleware/handler
```

**Example:**

```
Request → Middleware1 → Middleware2 → Handler

If Middleware1 calls c.Abort():
Request → Middleware1 [STOP]

If Middleware1 calls c.Next():
Request → Middleware1 → Middleware2 → Handler
```

---

## Step 2: Create Logger Middleware

### 2.1 Create logger.go File

📝 **Create file:** `internal/middleware/logger.go`

```go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware mencatat informasi request
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Hitung waktu eksekusi
		duration := time.Since(startTime)

		// Log informasi request
		log.Printf("[%s] %s %s - Status: %d - Duration: %v - IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}
```

- [ ] File created at `internal/middleware/logger.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Logger Middleware

**Logging Pattern:**

```
1. Record start time
2. Call c.Next() - let request process
3. After handler completes, calculate duration
4. Log request details
```

**Example log output:**

```
[GET] /api/v1/todos HTTP/1.1 - Status: 200 - Duration: 15ms - IP: 127.0.0.1
[POST] /api/v1/todos HTTP/1.1 - Status: 201 - Duration: 45ms - IP: 192.168.1.100
[PUT] /api/v1/todos/1 HTTP/1.1 - Status: 200 - Duration: 23ms - IP: 127.0.0.1
[DELETE] /api/v1/todos/1 HTTP/1.1 - Status: 204 - Duration: 12ms - IP: 127.0.0.1
```

### 2.3 What Gets Logged

| Field    | Description             | Example                  |
| -------- | ----------------------- | ------------------------ |
| Method   | HTTP method             | GET, POST, PUT, DELETE   |
| Path     | URL path                | /api/v1/todos            |
| Proto    | HTTP protocol           | HTTP/1.1, HTTP/2         |
| Status   | Response status code    | 200, 201, 404, 500       |
| Duration | Request processing time | 15ms, 250ms, 1.5s        |
| IP       | Client IP address       | 127.0.0.1, 192.168.1.100 |

---

## Step 3: Create Error and CORS Middleware

### 3.1 Create error.go File

📝 **Create file:** `internal/middleware/error.go`

```go
package middleware

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/gin-gonic/gin"
)

// ErrorHandler menangani error yang terjadi
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Cek jika ada error
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(500, dto.ErrorResponse{
				Success: false,
				Message: "Internal server error",
				Error:   err.Error(),
			})
		}
	}
}

// CORSMiddleware menangani CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

- [ ] File created at `internal/middleware/error.go`
- [ ] Content copied exactly
- [ ] File saved

### 3.2 Understanding Error Handler

```go
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()  // Process request first

		// After processing, check for errors
		if len(c.Errors) > 0 {
			// Return error response
		}
	}
}
```

**When to use:**

- Catch panic errors
- Handle unhandled errors from handlers
- Centralized error response format

**Example:**

```go
// In handler
if err := something(); err != nil {
	c.Error(err)  // Add error to context
	return
}

// ErrorHandler middleware catches it and returns JSON
```

### 3.3 Understanding CORS

**CORS (Cross-Origin Resource Sharing):**

- Allows frontend from different domain to access API
- Browser security feature
- Required for modern web apps

**Headers explained:**

```go
// Which domains can access (*)
"Access-Control-Allow-Origin": "*"

// Can send credentials (cookies)
"Access-Control-Allow-Credentials": "true"

// Allowed HTTP methods
"Access-Control-Allow-Methods": "POST, OPTIONS, GET, PUT, DELETE, PATCH"

// Allowed headers in request
"Access-Control-Allow-Headers": "Content-Type, Authorization, ..."
```

**⚠️ Production tip:**
Replace `"*"` with specific domain:

```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
```

### 3.4 OPTIONS Preflight Request

```go
if c.Request.Method == "OPTIONS" {
	c.AbortWithStatus(204)
	return
}
```

**What is preflight?**

- Browser sends OPTIONS request before actual request
- Checks if CORS is allowed
- Must respond with 204 No Content
- Then browser sends actual request

**Flow:**

```
1. Browser: OPTIONS /api/v1/todos (preflight)
2. Server: 204 with CORS headers
3. Browser: POST /api/v1/todos (actual request)
4. Server: 201 with data
```

---

## Step 4: Middleware Order Matters

### 4.1 Typical Middleware Chain

```go
router := gin.Default()

// Global middleware (applied to all routes)
router.Use(LoggerMiddleware())    // 1. Log all requests
router.Use(CORSMiddleware())      // 2. Handle CORS
router.Use(ErrorHandler())        // 3. Catch errors

// Route-specific middleware
protected := router.Group("/api")
protected.Use(AuthMiddleware())   // 4. Require auth
protected.GET("/todos", handler)  // 5. Handler
```

**Execution order:**

```
Request
  ↓
Logger (logs start)
  ↓
CORS (sets headers)
  ↓
ErrorHandler (sets up catch)
  ↓
Auth (validates token)
  ↓
Handler (processes request)
  ↓
Auth (after c.Next())
  ↓
ErrorHandler (checks errors)
  ↓
CORS (after c.Next())
  ↓
Logger (logs result)
  ↓
Response
```

### 4.2 Why Order Matters

**❌ Wrong order:**

```go
router.Use(AuthMiddleware())  // Requires token
router.Use(CORSMiddleware())  // Sets CORS headers
```

Problem: Browser OPTIONS request blocked by auth!

**✅ Correct order:**

```go
router.Use(CORSMiddleware())  // Handle CORS first
router.Use(AuthMiddleware())  // Then check auth
```

---

## Step 5: Verify Middleware

### 5.1 Check File Structure

```bash
$ ls -la internal/middleware/    # Linux/Mac
$ dir internal\middleware\       # Windows
```

**Expected output:**

```
auth.go
error.go
logger.go
```

- [ ] All three files present

### 5.2 Verify Middleware Compiles

```bash
$ go build internal/middleware/*.go
```

**Expected:** No errors

- [ ] Middleware compiles successfully

---

## Step 6: Understanding Middleware Patterns

### 6.1 Middleware Function Signature

```go
func MiddlewareName() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before handler
		c.Next()
		// After handler
	}
}
```

**Pattern explained:**

```go
func MiddlewareName() gin.HandlerFunc {
	// Setup (runs once when middleware registered)
	config := loadConfig()

	return func(c *gin.Context) {
		// Runs on every request
		// Do something with config
		c.Next()
	}
}
```

### 6.2 Stopping Execution

```go
// Early return stops execution
if unauthorized {
	c.JSON(401, ErrorResponse{...})
	c.Abort()
	return  // Important: return after Abort
}
```

### 6.3 Passing Data

```go
// Set in middleware
c.Set("key", value)

// Get in handler
value, exists := c.Get("key")
if exists {
	// Use value
}
```

---

## Step 7: Commit Changes

### 7.1 Check Status

```bash
$ git status
```

- [ ] Middleware files shown

### 7.2 Stage Middleware

```bash
$ git add internal/middleware/
```

- [ ] Files staged

### 7.3 Commit

```bash
$ git commit -m "Add authentication, logging, CORS, and error handling middleware"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/middleware/auth.go` created
- [ ] AuthMiddleware validates JWT tokens
- [ ] Authorization header parsing
- [ ] Bearer token format validation
- [ ] Token validation with utils
- [ ] UserID stored in context
- [ ] Proper error responses with 401 status
- [ ] `internal/middleware/logger.go` created
- [ ] LoggerMiddleware logs requests
- [ ] Request duration calculation
- [ ] All request details logged
- [ ] `internal/middleware/error.go` created
- [ ] ErrorHandler catches errors
- [ ] CORSMiddleware handles CORS
- [ ] OPTIONS preflight handling
- [ ] Understanding of middleware chain
- [ ] Understanding of execution order
- [ ] Middleware compiles without errors
- [ ] Changes committed to git

---

## 🧪 Quick Test

We can't fully test middleware without handlers, but we can verify the structure:

📝 **Create temporary file:** `test_middleware.go`

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// Test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Middleware loaded successfully"})
	})

	fmt.Println("✅ Middleware functions loaded successfully")
	fmt.Println("✅ Router configured with middleware chain")
	fmt.Println("✅ Ready for testing with actual handlers")
}
```

### Run Test

```bash
$ go run test_middleware.go
```

**Expected output:**

```
✅ Middleware functions loaded successfully
✅ Router configured with middleware chain
✅ Ready for testing with actual handlers
```

- [ ] Test runs without errors
- [ ] Middleware loads successfully

### Clean Up

```bash
$ rm test_middleware.go    # Linux/Mac
$ del test_middleware.go   # Windows
```

---

## 🐛 Common Issues

### Issue: "c.Set not working"

**Solution:** Make sure you call `c.Next()` in middleware

### Issue: "CORS not working"

**Solution:** Ensure CORSMiddleware is applied BEFORE route handlers

### Issue: "Auth blocking OPTIONS requests"

**Solution:** Apply CORS middleware before Auth middleware

### Issue: "userID not found in context"

**Solution:** Ensure AuthMiddleware is applied to the route group

---

## 📊 Current Progress

```
[████████████████████████] 43% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware  
**Next:** Repository layer

---

**Previous:** [06-utilities.md](06-utilities.md)  
**Next:** [08-repositories.md](08-repositories.md)
