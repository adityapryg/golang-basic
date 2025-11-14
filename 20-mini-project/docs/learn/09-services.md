# 09 - Service Layer

Create business logic layer with authentication and todo services.

---

## Overview

We'll create:

1. `internal/service/auth_service.go` - Authentication and user management
2. `internal/service/todo_service.go` - Todo business logic

Services contain business rules, validation, and coordinate repository operations.

---

## Step 1: Create Authentication Service

### 1.1 Create auth_service.go File

📝 **Create file:** `internal/service/auth_service.go`

```go
package service

import (
	"errors"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
)

// Custom errors for domain logic
var (
	ErrUserExists         = errors.New("username sudah digunakan")
	ErrEmailExists        = errors.New("email sudah digunakan")
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates new auth service instance
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register mendaftarkan user baru
func (s *AuthService) Register(req dto.RegisterRequest) (*model.User, error) {
	// Cek apakah username sudah ada
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// Cek apakah email sudah ada
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Buat user baru
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	// Simpan ke database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login melakukan autentikasi user
func (s *AuthService) Login(req dto.LoginRequest) (string, *model.User, error) {
	// Cari user berdasarkan username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Verifikasi password
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GetUserByID mengambil data user berdasarkan ID
func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser memperbarui data user
func (s *AuthService) UpdateUser(id uint, req dto.UpdateProfileRequest) (*model.User, error) {
	// Ambil user yang ada
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update email jika ada
	if req.Email != nil && *req.Email != user.Email {
		// Cek apakah email baru sudah digunakan
		exists, err := s.userRepo.ExistsByEmail(*req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	// Update full name jika ada
	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	// Update password jika ada
	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	// Simpan perubahan
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

- [ ] File created at `internal/service/auth_service.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Custom Errors

```go
var (
	ErrUserExists         = errors.New("username sudah digunakan")
	ErrEmailExists        = errors.New("email sudah digunakan")
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)
```

**Why package-level error variables?**

- Handlers can check specific errors with `errors.Is()`
- Consistent error messages
- Easy to maintain
- Enables proper HTTP status code mapping

**Usage in handler:**

```go
err := authService.Register(req)
if errors.Is(err, service.ErrUserExists) {
	c.JSON(409, gin.H{"error": "Username taken"})
} else if errors.Is(err, service.ErrEmailExists) {
	c.JSON(409, gin.H{"error": "Email taken"})
}
```

### 1.3 Understanding Register Method

```go
func (s *AuthService) Register(req dto.RegisterRequest) (*model.User, error) {
	// 1. Check username uniqueness
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if exists {
		return nil, ErrUserExists
	}

	// 2. Check email uniqueness
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if exists {
		return nil, ErrEmailExists
	}

	// 3. Hash password (NEVER store plain text!)
	hashedPassword, err := utils.HashPassword(req.Password)

	// 4. Create user model
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	// 5. Save to database
	err = s.userRepo.Create(user)

	return user, nil
}
```

**Business rules enforced:**

- Username must be unique
- Email must be unique
- Password must be hashed
- All fields required (validated by DTO)

### 1.4 Understanding Login Method

```go
func (s *AuthService) Login(req dto.LoginRequest) (string, *model.User, error) {
	// 1. Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, ErrInvalidCredentials  // Don't reveal user existence!
	}

	// 2. Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// 3. Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
```

**Security considerations:**

- Return same error for "user not found" and "wrong password"
- Don't reveal whether username exists
- Password comparison uses constant-time algorithm (bcrypt)

**⚠️ Security anti-pattern:**

```go
// ❌ Bad - reveals user existence
user, err := s.userRepo.FindByUsername(req.Username)
if err != nil {
	return "", nil, errors.New("user not found")
}
if !utils.CheckPassword(req.Password, user.Password) {
	return "", nil, errors.New("wrong password")
}
```

### 1.5 Understanding UpdateUser Method

```go
func (s *AuthService) UpdateUser(id uint, req dto.UpdateProfileRequest) (*model.User, error) {
	// 1. Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 2. Update email (if provided)
	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(*req.Email)
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	// 3. Update full name (if provided)
	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	// 4. Update password (if provided)
	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	// 5. Save changes
	err = s.userRepo.Update(user)

	return user, nil
}
```

**Why check `req.Email != nil`?**

- DTO uses pointer fields for optional updates
- `nil` = field not provided in request
- `*req.Email` = dereference pointer to get value

**Example requests:**

```json
// Update only email
{
  "email": "newemail@example.com"
}

// Update only full name
{
  "full_name": "New Name"
}

// Update multiple fields
{
  "email": "newemail@example.com",
  "full_name": "New Name",
  "password": "newpassword123"
}
```

---

## Step 2: Create Todo Service

### 2.1 Create todo_service.go File

📝 **Create file:** `internal/service/todo_service.go`

```go
package service

import (
	"errors"
	"time"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
)

// Custom errors for todo operations
var (
	ErrTodoNotFound       = errors.New("todo tidak ditemukan")
	ErrUnauthorizedAccess = errors.New("tidak memiliki akses ke todo ini")
	ErrInvalidStatus      = errors.New("status tidak valid")
	ErrInvalidPriority    = errors.New("priority tidak valid")
)

// TodoService handles todo business logic
type TodoService struct {
	todoRepo *repository.TodoRepository
}

// NewTodoService creates new todo service instance
func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

// CreateTodo membuat todo baru
func (s *TodoService) CreateTodo(userID uint, req dto.CreateTodoRequest) (*model.Todo, error) {
	// Validasi status
	if !isValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	// Validasi priority
	if !isValidPriority(req.Priority) {
		return nil, ErrInvalidPriority
	}

	// Parse due date jika ada
	var dueDate *time.Time
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, errors.New("format due date tidak valid (gunakan YYYY-MM-DD)")
		}
		dueDate = &parsedDate
	}

	// Buat todo baru
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     dueDate,
		UserID:      userID,
	}

	// Simpan ke database
	err := s.todoRepo.Create(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// GetTodoByID mengambil todo berdasarkan ID
func (s *TodoService) GetTodoByID(todoID, userID uint) (*model.Todo, error) {
	// Cek ownership
	owned, err := s.todoRepo.IsOwnedByUser(todoID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrUnauthorizedAccess
	}

	// Ambil todo
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

// GetUserTodos mengambil semua todo user dengan filter opsional
func (s *TodoService) GetUserTodos(userID uint, status, priority string) ([]model.Todo, error) {
	// Validasi filter jika ada
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}
	if priority != "" && !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}

	// Ambil todos
	if status != "" || priority != "" {
		return s.todoRepo.FindByUserIDWithFilters(userID, status, priority)
	}

	return s.todoRepo.FindByUserID(userID)
}

// UpdateTodo memperbarui todo
func (s *TodoService) UpdateTodo(todoID, userID uint, req dto.UpdateTodoRequest) (*model.Todo, error) {
	// Cek ownership
	owned, err := s.todoRepo.IsOwnedByUser(todoID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrUnauthorizedAccess
	}

	// Ambil todo yang ada
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	// Update fields jika ada
	if req.Title != nil {
		todo.Title = *req.Title
	}

	if req.Description != nil {
		todo.Description = *req.Description
	}

	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		todo.Status = *req.Status
	}

	if req.Priority != nil {
		if !isValidPriority(*req.Priority) {
			return nil, ErrInvalidPriority
		}
		todo.Priority = *req.Priority
	}

	if req.DueDate != nil {
		if *req.DueDate == "" {
			todo.DueDate = nil
		} else {
			parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				return nil, errors.New("format due date tidak valid (gunakan YYYY-MM-DD)")
			}
			todo.DueDate = &parsedDate
		}
	}

	// Simpan perubahan
	err = s.todoRepo.Update(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo menghapus todo
func (s *TodoService) DeleteTodo(todoID, userID uint) error {
	// Cek ownership
	owned, err := s.todoRepo.IsOwnedByUser(todoID, userID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrUnauthorizedAccess
	}

	// Hapus todo
	return s.todoRepo.Delete(todoID)
}

// Helper functions

func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func isValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high"}
	for _, p := range validPriorities {
		if p == priority {
			return true
		}
	}
	return false
}
```

- [ ] File created at `internal/service/todo_service.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding CreateTodo Method

```go
func (s *TodoService) CreateTodo(userID uint, req dto.CreateTodoRequest) (*model.Todo, error) {
	// 1. Validate status (business rule)
	if !isValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	// 2. Validate priority (business rule)
	if !isValidPriority(req.Priority) {
		return nil, ErrInvalidPriority
	}

	// 3. Parse optional due date
	var dueDate *time.Time
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, errors.New("format due date tidak valid")
		}
		dueDate = &parsedDate
	}

	// 4. Create todo model
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     dueDate,
		UserID:      userID,  // From JWT token
	}

	// 5. Save to database
	err := s.todoRepo.Create(todo)

	return todo, nil
}
```

**Date parsing explained:**

```go
// Format: "2006-01-02" is Go's reference time
// 2006 = year, 01 = month, 02 = day

time.Parse("2006-01-02", "2024-03-15")  // ✅ Valid
time.Parse("2006-01-02", "15-03-2024")  // ❌ Invalid format
time.Parse("2006-01-02", "2024/03/15")  // ❌ Invalid format
```

**Why use pointer for DueDate?**

```go
var dueDate *time.Time  // nil = no due date

// If provided
dueDate = &parsedDate   // pointer to parsed date

// In model
DueDate *time.Time      // nullable field
```

### 2.3 Understanding GetTodoByID Method

```go
func (s *TodoService) GetTodoByID(todoID, userID uint) (*model.Todo, error) {
	// 1. Security check - verify ownership FIRST
	owned, err := s.todoRepo.IsOwnedByUser(todoID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrUnauthorizedAccess
	}

	// 2. After security check passed, fetch todo
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}
```

**Why check ownership first?**

- **Security**: Prevent unauthorized access
- **IDOR prevention**: Insecure Direct Object Reference
- **Principle of least privilege**: Check permission before data access

**Attack scenario prevented:**

```
User A (ID: 1) tries to access User B's todo (ID: 100)

Without ownership check:
  GET /todos/100 → Returns todo → ❌ Security breach

With ownership check:
  GET /todos/100 → 403 Unauthorized → ✅ Prevented
```

### 2.4 Understanding UpdateTodo Method

```go
func (s *TodoService) UpdateTodo(todoID, userID uint, req dto.UpdateTodoRequest) (*model.Todo, error) {
	// 1. Security check
	owned, err := s.todoRepo.IsOwnedByUser(todoID, userID)
	if !owned {
		return nil, ErrUnauthorizedAccess
	}

	// 2. Get existing todo
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	// 3. Update only provided fields
	if req.Title != nil {
		todo.Title = *req.Title
	}

	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		todo.Status = *req.Status
	}

	// 4. Save changes
	err = s.todoRepo.Update(todo)

	return todo, nil
}
```

**Partial update pattern:**

```json
// Only update title
{
  "title": "New Title"
}

// Only update status
{
  "status": "completed"
}

// Update multiple fields
{
  "title": "New Title",
  "status": "completed",
  "priority": "high"
}
```

### 2.5 Understanding Helper Functions

```go
func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
```

**Why not use database enum?**

- More flexible (easy to change)
- Centralized validation
- Clear business rules in code
- No database migration for status changes

**Alternative implementations:**

```go
// Using map
var validStatuses = map[string]bool{
	"pending":     true,
	"in_progress": true,
	"completed":   true,
}

func isValidStatus(status string) bool {
	return validStatuses[status]
}

// Using constants
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)
```

---

## Step 3: Service Layer Best Practices

### 3.1 Service Layer Responsibilities

✅ **Should do:**

- Business logic and validation
- Coordinate multiple repository calls
- Define custom domain errors
- Transform DTOs to models
- Enforce business rules

❌ **Should NOT do:**

- HTTP handling (that's handler layer)
- Direct database queries (that's repository layer)
- JSON parsing (that's DTO layer)

### 3.2 Error Handling Pattern

```go
// Define errors at package level
var (
	ErrUserExists = errors.New("username already exists")
	ErrNotFound   = errors.New("user not found")
)

// Return specific errors from service
func (s *Service) Create(user *User) error {
	exists, _ := s.repo.Exists(user.Username)
	if exists {
		return ErrUserExists  // Specific error
	}
	return s.repo.Create(user)
}

// Check errors in handler
err := service.Create(user)
if errors.Is(err, service.ErrUserExists) {
	c.JSON(409, gin.H{"error": "Username taken"})
} else if errors.Is(err, service.ErrNotFound) {
	c.JSON(404, gin.H{"error": "Not found"})
}
```

### 3.3 Dependency Injection Pattern

```go
// Service depends on repository interface (not implementation)
type AuthService struct {
	userRepo *repository.UserRepository  // Dependency
}

// Constructor injection
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Usage in main.go
userRepo := repository.NewUserRepository(db)
authService := service.NewAuthService(userRepo)  // Inject dependency
```

**Benefits:**

- Easy testing (can mock repository)
- Loose coupling
- Clear dependencies

---

## Step 4: Verify Services

### 4.1 Check File Structure

```bash
$ ls -la internal/service/    # Linux/Mac
$ dir internal\service\       # Windows
```

**Expected output:**

```
auth_service.go
todo_service.go
```

- [ ] Both files present

### 4.2 Verify Compilation

```bash
$ go build internal/service/*.go
```

**Expected:** No errors

- [ ] Services compile successfully

### 4.3 Check Error Definitions

```bash
$ grep -r "^var (" internal/service/    # Linux/Mac
$ findstr /B "var" internal\service\*   # Windows
```

**Expected:**

- Custom errors in auth_service.go
- Custom errors in todo_service.go

- [ ] Error variables defined

---

## Step 5: Commit Changes

### 5.1 Check Status

```bash
$ git status
```

- [ ] Service files shown

### 5.2 Stage Services

```bash
$ git add internal/service/
```

- [ ] Files staged

### 5.3 Commit

```bash
$ git commit -m "Add authentication and todo services with business logic"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/service/auth_service.go` created
- [ ] AuthService struct defined
- [ ] Custom error variables defined
- [ ] NewAuthService constructor implemented
- [ ] Register method with validation
- [ ] Username uniqueness check
- [ ] Email uniqueness check
- [ ] Password hashing
- [ ] Login method with authentication
- [ ] Password verification
- [ ] JWT token generation
- [ ] GetUserByID method
- [ ] UpdateUser method with partial updates
- [ ] Email uniqueness check on update
- [ ] Password hashing on update
- [ ] `internal/service/todo_service.go` created
- [ ] TodoService struct defined
- [ ] Todo custom errors defined
- [ ] NewTodoService constructor implemented
- [ ] CreateTodo method with validation
- [ ] Status validation with isValidStatus
- [ ] Priority validation with isValidPriority
- [ ] Due date parsing
- [ ] GetTodoByID method with ownership check
- [ ] GetUserTodos method with filters
- [ ] UpdateTodo method with security check
- [ ] Partial update pattern
- [ ] DeleteTodo method with ownership check
- [ ] Helper functions for validation
- [ ] Understanding of business logic separation
- [ ] Understanding of error handling pattern
- [ ] Services compile without errors
- [ ] Changes committed to git

---

## 🧪 Quick Test

We can verify service logic with a simple test:

📝 **Create temporary file:** `test_service.go`

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}, &model.Todo{})

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Create services
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Test registration
	registerReq := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	user, err := authService.Register(registerReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User registered: %s (ID: %d)\n", user.Username, user.ID)

	// Test duplicate username
	_, err = authService.Register(registerReq)
	if err == service.ErrUserExists {
		fmt.Println("✅ Duplicate username detected correctly")
	}

	// Test login
	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	token, loginUser, err := authService.Login(loginReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Login successful, token generated: %s...\n", token[:20])

	// Test wrong password
	wrongReq := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	_, _, err = authService.Login(wrongReq)
	if err == service.ErrInvalidCredentials {
		fmt.Println("✅ Wrong password detected correctly")
	}

	// Test create todo
	createTodoReq := dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test description",
		Status:      "pending",
		Priority:    "high",
	}
	todo, err := todoService.CreateTodo(loginUser.ID, createTodoReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Todo created: %s (ID: %d)\n", todo.Title, todo.ID)

	// Test invalid status
	invalidReq := dto.CreateTodoRequest{
		Title:    "Invalid Todo",
		Status:   "invalid_status",
		Priority: "high",
	}
	_, err = todoService.CreateTodo(loginUser.ID, invalidReq)
	if err == service.ErrInvalidStatus {
		fmt.Println("✅ Invalid status detected correctly")
	}

	// Test get user todos
	todos, err := todoService.GetUserTodos(loginUser.ID, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Retrieved %d todos for user\n", len(todos))

	// Test ownership check
	_, err = todoService.GetTodoByID(todo.ID, 9999)  // Wrong user ID
	if err == service.ErrUnauthorizedAccess {
		fmt.Println("✅ Unauthorized access prevented correctly")
	}

	fmt.Println("\n✅ All service methods working correctly!")
}
```

### Install SQLite Driver (if not already)

```bash
$ go get gorm.io/driver/sqlite
```

### Run Test

```bash
$ go run test_service.go
```

**Expected output:**

```
✅ User registered: testuser (ID: 1)
✅ Duplicate username detected correctly
✅ Login successful, token generated: eyJhbGciOiJIUzI1NiIsInR...
✅ Wrong password detected correctly
✅ Todo created: Test Todo (ID: 1)
✅ Invalid status detected correctly
✅ Retrieved 1 todos for user
✅ Unauthorized access prevented correctly

✅ All service methods working correctly!
```

- [ ] Test runs without errors
- [ ] All validations work
- [ ] Security checks pass

### Clean Up

```bash
$ rm test_service.go    # Linux/Mac
$ del test_service.go   # Windows
```

---

## 🐛 Common Issues

### Issue: "ErrUserExists not working"

**Solution:** Ensure error variable is exported (starts with capital letter)

### Issue: "Ownership check always fails"

**Solution:** Verify userID from JWT token is correctly passed to service

### Issue: "Date parsing fails"

**Solution:** Use exact format "2006-01-02" (YYYY-MM-DD)

### Issue: "Password comparison fails"

**Solution:** Ensure password is hashed before comparing

---

## 📊 Current Progress

```
[███████████████████████████] 56% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services  
**Next:** Handler layer for HTTP endpoints

---

**Previous:** [08-repositories.md](08-repositories.md)  
**Next:** [10-handlers.md](10-handlers.md)
