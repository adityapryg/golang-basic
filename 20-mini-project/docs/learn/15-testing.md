# 15 - Testing

Implement integration tests with testify suite pattern.

---

## Overview

We'll create:

1. `tests/auth_test.go` - Authentication flow tests
2. `tests/todo_test.go` - Todo CRUD operation tests

Tests ensure code correctness and prevent regressions.

---

## Step 1: Understanding Testing in Go

### 1.1 Testing Basics

**Standard Go testing:**

```go
func TestSomething(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}
```

**With testify:**

```go
func TestSomething(t *testing.T) {
	result := Add(2, 3)
	assert.Equal(t, 5, result, "Should add correctly")
}
```

**Test suite pattern:**

```go
type MyTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

func (suite *MyTestSuite) SetupSuite() {
	// Run once before all tests
}

func (suite *MyTestSuite) TearDownSuite() {
	// Run once after all tests
}

func (suite *MyTestSuite) TestSomething() {
	// Individual test
}
```

### 1.2 Test File Structure

```
tests/
├── auth_test.go      # Authentication tests
├── todo_test.go      # Todo CRUD tests
└── helper_test.go    # Shared test utilities (optional)
```

---

## Step 2: Create Authentication Tests

### 2.1 Create auth_test.go

📝 **Create file:** `tests/auth_test.go`

```go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/route"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AuthTestSuite defines the test suite
type AuthTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

// SetupSuite runs once before all tests
func (suite *AuthTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to test database
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err, "Failed to connect to test database")

	suite.db = db

	// Auto-migrate models
	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	suite.Require().NoError(err, "Failed to migrate test database")

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := handler.NewUserHandler(authService)
	healthHandler := handler.NewHealthHandler(db)

	// Dummy handlers for routes that won't be tested
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	// Setup router
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	route.SetupRoutes(router, userHandler, todoHandler, healthHandler)

	suite.router = router
}

// TearDownSuite runs once after all tests
func (suite *AuthTestSuite) TearDownSuite() {
	// Clean up test data
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// SetupTest runs before each test
func (suite *AuthTestSuite) SetupTest() {
	// Clean tables before each test
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// TestRegisterSuccess tests successful user registration
func (suite *AuthTestSuite) TestRegisterSuccess() {
	// Prepare request
	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Message, "berhasil")

	// Verify user data in response
	userData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", userData["username"])
	assert.Equal(suite.T(), "test@example.com", userData["email"])
}

// TestRegisterDuplicateUsername tests registration with existing username
func (suite *AuthTestSuite) TestRegisterDuplicateUsername() {
	// Create first user
	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Try to create duplicate username
	reqBody2 := dto.RegisterRequest{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "password123",
		FullName: "Test User 2",
	}
	jsonBody2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	// Assert error response
	assert.Equal(suite.T(), http.StatusConflict, w2.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
}

// TestLoginSuccess tests successful login
func (suite *AuthTestSuite) TestLoginSuccess() {
	// Register user first
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	// Login
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify token in response
	loginData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), loginData["token"])
	assert.NotNil(suite.T(), loginData["user"])
}

// TestLoginInvalidCredentials tests login with wrong password
func (suite *AuthTestSuite) TestLoginInvalidCredentials() {
	// Register user first
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	// Login with wrong password
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert error response
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
}

// TestGetProfile tests getting user profile with authentication
func (suite *AuthTestSuite) TestGetProfile() {
	// Register and login to get token
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	suite.router.ServeHTTP(wLogin, reqLogin)

	var loginResponse dto.SuccessResponse
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	loginData := loginResponse.Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Get profile with token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", userData["username"])
}

// Run the test suite
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
```

- [ ] File created at `tests/auth_test.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Test Suite Pattern

```go
type AuthTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

func (suite *AuthTestSuite) SetupSuite() {
	// Runs ONCE before all tests
}

func (suite *AuthTestSuite) SetupTest() {
	// Runs BEFORE EACH test
}

func (suite *AuthTestSuite) TearDownTest() {
	// Runs AFTER EACH test
}

func (suite *AuthTestSuite) TearDownSuite() {
	// Runs ONCE after all tests
}
```

**Execution order:**

```
SetupSuite
  ↓
SetupTest → Test1 → TearDownTest
  ↓
SetupTest → Test2 → TearDownTest
  ↓
SetupTest → Test3 → TearDownTest
  ↓
TearDownSuite
```

### 2.3 Understanding httptest

```go
// Create test request
req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
req.Header.Set("Content-Type", "application/json")

// Record response
w := httptest.NewRecorder()

// Execute request
suite.router.ServeHTTP(w, req)

// Check response
assert.Equal(suite.T(), 201, w.Code)
assert.Contains(suite.T(), w.Body.String(), "success")
```

---

## Step 3: Create Todo Tests

### 3.1 Create todo_test.go

📝 **Create file:** `tests/todo_test.go`

```go
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/route"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TodoTestSuite defines the test suite
type TodoTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	token  string
	userID uint
}

// SetupSuite runs once before all tests
func (suite *TodoTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = db

	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	suite.Require().NoError(err)

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)
	userHandler := handler.NewUserHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealthHandler(db)

	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	route.SetupRoutes(router, userHandler, todoHandler, healthHandler)

	suite.router = router

	// Create test user and get token
	suite.createTestUserAndToken()
}

// createTestUserAndToken creates a test user and stores token
func (suite *TodoTestSuite) createTestUserAndToken() {
	// Register
	regBody := dto.RegisterRequest{
		Username: "todotest",
		Email:    "todotest@example.com",
		Password: "password123",
		FullName: "Todo Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	var regResponse dto.SuccessResponse
	json.Unmarshal(wReg.Body.Bytes(), &regResponse)
	userData := regResponse.Data.(map[string]interface{})
	suite.userID = uint(userData["id"].(float64))

	// Login
	loginBody := dto.LoginRequest{
		Username: "todotest",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	suite.router.ServeHTTP(wLogin, reqLogin)

	var loginResponse dto.SuccessResponse
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	loginData := loginResponse.Data.(map[string]interface{})
	suite.token = loginData["token"].(string)
}

// TearDownSuite runs once after all tests
func (suite *TodoTestSuite) TearDownSuite() {
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// SetupTest runs before each test
func (suite *TodoTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM todos WHERE user_id = ?", suite.userID)
}

// TestCreateTodo tests creating a new todo
func (suite *TodoTestSuite) TestCreateTodo() {
	reqBody := dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test description",
		Status:      "pending",
		Priority:    "high",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Test Todo", todoData["title"])
	assert.Equal(suite.T(), "pending", todoData["status"])
}

// TestGetAllTodos tests getting all user's todos
func (suite *TodoTestSuite) TestGetAllTodos() {
	// Create test todos
	suite.createTestTodo("Todo 1", "pending", "high")
	suite.createTestTodo("Todo 2", "in_progress", "medium")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todos, ok := response.Data.([]interface{})
	assert.True(suite.T(), ok)
	assert.GreaterOrEqual(suite.T(), len(todos), 2)
}

// TestGetTodoByID tests getting a specific todo
func (suite *TodoTestSuite) TestGetTodoByID() {
	todoID := suite.createTestTodo("Test Todo", "pending", "high")

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/todos/%d", todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Test Todo", todoData["title"])
}

// TestUpdateTodo tests updating a todo
func (suite *TodoTestSuite) TestUpdateTodo() {
	todoID := suite.createTestTodo("Original Title", "pending", "low")

	updateBody := map[string]interface{}{
		"title":  "Updated Title",
		"status": "completed",
	}
	jsonBody, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/todos/%d", todoID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Updated Title", todoData["title"])
	assert.Equal(suite.T(), "completed", todoData["status"])
}

// TestDeleteTodo tests deleting a todo
func (suite *TodoTestSuite) TestDeleteTodo() {
	todoID := suite.createTestTodo("To Delete", "pending", "low")

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/todos/%d", todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify todo is deleted (soft delete)
	var todo model.Todo
	err = suite.db.Where("id = ?", todoID).First(&todo).Error
	assert.Error(suite.T(), err) // Should not find (soft deleted)
}

// Helper function to create test todo
func (suite *TodoTestSuite) createTestTodo(title, status, priority string) uint {
	reqBody := dto.CreateTodoRequest{
		Title:    title,
		Status:   status,
		Priority: priority,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response dto.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	todoData := response.Data.(map[string]interface{})
	return uint(todoData["id"].(float64))
}

// Run the test suite
func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
```

- [ ] File created at `tests/todo_test.go`
- [ ] Content copied exactly
- [ ] File saved

---

## Step 4: Run Tests

### 4.1 Run All Tests

```bash
$ go test ./tests/... -v
```

**Expected output:**

```
=== RUN   TestAuthTestSuite
=== RUN   TestAuthTestSuite/TestRegisterSuccess
=== RUN   TestAuthTestSuite/TestRegisterDuplicateUsername
=== RUN   TestAuthTestSuite/TestLoginSuccess
=== RUN   TestAuthTestSuite/TestLoginInvalidCredentials
=== RUN   TestAuthTestSuite/TestGetProfile
--- PASS: TestAuthTestSuite (1.23s)
    --- PASS: TestAuthTestSuite/TestRegisterSuccess (0.15s)
    --- PASS: TestAuthTestSuite/TestRegisterDuplicateUsername (0.12s)
    --- PASS: TestAuthTestSuite/TestLoginSuccess (0.18s)
    --- PASS: TestAuthTestSuite/TestLoginInvalidCredentials (0.14s)
    --- PASS: TestAuthTestSuite/TestGetProfile (0.21s)

=== RUN   TestTodoTestSuite
=== RUN   TestTodoTestSuite/TestCreateTodo
=== RUN   TestTodoTestSuite/TestGetAllTodos
=== RUN   TestTodoTestSuite/TestGetTodoByID
=== RUN   TestTodoTestSuite/TestUpdateTodo
=== RUN   TestTodoTestSuite/TestDeleteTodo
--- PASS: TestTodoTestSuite (0.95s)
    --- PASS: TestTodoTestSuite/TestCreateTodo (0.12s)
    --- PASS: TestTodoTestSuite/TestGetAllTodos (0.18s)
    --- PASS: TestTodoTestSuite/TestGetTodoByID (0.11s)
    --- PASS: TestTodoTestSuite/TestUpdateTodo (0.14s)
    --- PASS: TestTodoTestSuite/TestDeleteTodo (0.13s)

PASS
ok      github.com/adityapryg/golang-demo/20-mini-project/tests    2.180s
```

- [ ] All tests pass
- [ ] No failures

### 4.2 Run Specific Test

```bash
# Run only auth tests
$ go test ./tests/ -run TestAuthTestSuite -v

# Run specific test method
$ go test ./tests/ -run TestAuthTestSuite/TestRegisterSuccess -v

# Run only todo tests
$ go test ./tests/ -run TestTodoTestSuite -v
```

### 4.3 Run with Coverage

```bash
$ go test ./tests/... -cover
```

**Expected output:**

```
ok      github.com/adityapryg/golang-demo/20-mini-project/tests    2.180s  coverage: 65.4% of statements
```

- [ ] Coverage report generated

### 4.4 Generate Coverage Report

```bash
# Generate coverage profile
$ go test ./tests/... -coverprofile=coverage.out

# View coverage in browser
$ go tool cover -html=coverage.out
```

- [ ] HTML coverage report opens in browser

---

## Step 5: Testing Best Practices

### 5.1 Test Naming Conventions

```go
// ✅ Good test names
func TestRegisterSuccess()
func TestLoginInvalidCredentials()
func TestUpdateTodoWithoutPermission()

// ❌ Bad test names
func TestRegister()
func TestLogin()
func TestUpdate()
```

**Pattern:** `Test<Function><Scenario>`

### 5.2 Arrange-Act-Assert Pattern

```go
func TestSomething(suite *Suite) {
	// Arrange - Setup test data
	user := createTestUser()

	// Act - Execute function being tested
	result := suite.service.DoSomething(user)

	// Assert - Verify results
	assert.Equal(suite.T(), expected, result)
}
```

### 5.3 Test Independence

```go
// ✅ Good - tests are independent
func (suite *Suite) SetupTest() {
	suite.db.Exec("DELETE FROM todos")  // Clean before each
}

// ❌ Bad - tests depend on each other
func TestCreate() { /* creates ID 1 */ }
func TestUpdate() { /* assumes ID 1 exists */ }
```

### 5.4 Don't Test External Dependencies

```go
// ✅ Good - test your code
func TestCreateTodo() {
	todo := service.CreateTodo(...)  // Your code
	assert.NotNil(todo)
}

// ❌ Bad - testing GORM
func TestDatabaseInsert() {
	db.Create(&todo)  // Testing GORM, not your code
}
```

---

## Step 6: Commit Changes

### 6.1 Check Status

```bash
$ git status
```

- [ ] Test files shown

### 6.2 Stage Tests

```bash
$ git add tests/
```

- [ ] Files staged

### 6.3 Commit

```bash
$ git commit -m "Add integration tests for auth and todo endpoints"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `tests/auth_test.go` created
- [ ] AuthTestSuite defined
- [ ] SetupSuite implemented
- [ ] TearDownSuite implemented
- [ ] SetupTest for cleaning data
- [ ] TestRegisterSuccess implemented
- [ ] TestRegisterDuplicateUsername implemented
- [ ] TestLoginSuccess implemented
- [ ] TestLoginInvalidCredentials implemented
- [ ] TestGetProfile with token implemented
- [ ] `tests/todo_test.go` created
- [ ] TodoTestSuite defined
- [ ] Test user and token creation
- [ ] TestCreateTodo implemented
- [ ] TestGetAllTodos implemented
- [ ] TestGetTodoByID implemented
- [ ] TestUpdateTodo implemented
- [ ] TestDeleteTodo implemented
- [ ] Helper functions for test data
- [ ] Understanding of test suites
- [ ] Understanding of httptest
- [ ] Understanding of testify assertions
- [ ] All tests run successfully
- [ ] Coverage report generated
- [ ] Changes committed to git

---

## 🐛 Common Issues

### Issue: "Database connection failed in tests"

**Solution:** Ensure PostgreSQL is running and test DB exists

### Issue: "Tests fail randomly"

**Solution:** Tests not independent. Clean data in `SetupTest()`

### Issue: "Token invalid in tests"

**Solution:** Generate fresh token for each test or in `SetupSuite()`

### Issue: "Type assertion fails"

**Solution:** JSON unmarshal returns `map[string]interface{}` and `[]interface{}`

---

## 📊 Current Progress

```
[████████████████████████████████████████] 93% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers, routes, main app, migrations, Docker, testing  
**Next:** Final verification and deployment

---

**Previous:** [14-docker.md](14-docker.md)  
**Next:** [16-verification.md](16-verification.md)
