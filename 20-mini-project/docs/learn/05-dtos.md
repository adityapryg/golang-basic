# 05 - Data Transfer Objects (DTOs)

Create request/response DTOs for API communication.

---

## Overview

We'll create:

1. `internal/dto/user_dto.go` - User-related DTOs (register, login, profile)
2. `internal/dto/todo_dto.go` - Todo-related DTOs (create, update, query)

DTOs separate external API representation from internal models.

---

## Step 1: Create User DTOs

### 1.1 Create user_dto.go File

📝 **Create file:** `internal/dto/user_dto.go`

```go
package dto

import "time"

// ============================================
// USER REQUEST DTOs
// ============================================

// RegisterRequest untuk registrasi user baru
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,max=100"`
}

// LoginRequest untuk login user
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest untuk update profil user
type UpdateProfileRequest struct {
	Email    *string `json:"email" binding:"omitempty,email,max=100"`
	FullName *string `json:"full_name" binding:"omitempty,max=100"`
}

// ============================================
// USER RESPONSE DTOs
// ============================================

// UserResponse untuk response data user (tanpa password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse untuk response setelah login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ============================================
// COMMON RESPONSE DTOs
// ============================================

// SuccessResponse untuk response sukses
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse untuk response error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
```

- [ ] File created at `internal/dto/user_dto.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Binding Tags

💡 **Gin Binding Tags Explained:**

| Tag         | Description           | Example                           |
| ----------- | --------------------- | --------------------------------- |
| `required`  | Field must be present | `binding:"required"`              |
| `omitempty` | Field is optional     | `binding:"omitempty"`             |
| `min=N`     | Minimum length/value  | `binding:"min=3"`                 |
| `max=N`     | Maximum length/value  | `binding:"max=50"`                |
| `email`     | Must be valid email   | `binding:"email"`                 |
| `oneof=A B` | Must be one of values | `binding:"oneof=low medium high"` |

**Example:**

```go
Username string `json:"username" binding:"required,min=3,max=50"`
```

- Must be present (required)
- Minimum 3 characters
- Maximum 50 characters
- JSON field name is "username"

### 1.3 Pointer Fields Explained

```go
type UpdateProfileRequest struct {
	Email    *string `json:"email" binding:"omitempty,email,max=100"`
	FullName *string `json:"full_name" binding:"omitempty,max=100"`
}
```

**Why pointers?**

- Distinguish between "not provided" (nil) and "empty string" ("")
- Allows partial updates
- `nil` = don't update this field
- `&""` = update to empty string

**Example usage:**

```go
email := "new@example.com"
req := UpdateProfileRequest{
	Email: &email,    // Update email
	FullName: nil,    // Don't update full name
}
```

### 1.4 Response DTOs Purpose

**Why separate response DTOs?**

```go
// Model (has password)
type User struct {
	Password string `gorm:"not null"`
	// ... other fields
}

// Response DTO (no password!)
type UserResponse struct {
	// NO PASSWORD FIELD
	Username string `json:"username"`
	Email    string `json:"email"`
}
```

⚠️ **Security:** Never expose password hashes in API responses!

---

## Step 2: Create Todo DTOs

### 2.1 Create todo_dto.go File

📝 **Create file:** `internal/dto/todo_dto.go`

```go
package dto

import (
	"time"
)

// ============================================
// TODO REQUEST DTOs
// ============================================

// CreateTodoRequest untuk membuat todo baru
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending in_progress completed"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     string `json:"due_date" binding:"omitempty"` // Format: YYYY-MM-DD
}

// UpdateTodoRequest untuk update todo
type UpdateTodoRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"due_date"` // Format: YYYY-MM-DD or empty string to clear
}

// TodoCreateRequest untuk backward compatibility (alias)
type TodoCreateRequest = CreateTodoRequest

// TodoUpdateRequest untuk backward compatibility (alias)
type TodoUpdateRequest = UpdateTodoRequest

// TodoQueryParams untuk filter dan pagination
type TodoQueryParams struct {
	Status   string `form:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority string `form:"priority" binding:"omitempty,oneof=low medium high"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// ============================================
// TODO RESPONSE DTOs
// ============================================

// TodoResponse untuk response todo
type TodoResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	UserID      uint       `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
```

- [ ] File created at `internal/dto/todo_dto.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding oneof Validation

```go
Status string `json:"status" binding:"required,oneof=pending in_progress completed"`
```

**What it does:**

- Value must be exactly one of: `pending`, `in_progress`, or `completed`
- Any other value will be rejected
- Spaces separate allowed values

**Examples:**

```json
// ✅ Valid
{"status": "pending"}
{"status": "in_progress"}
{"status": "completed"}

// ❌ Invalid
{"status": "done"}        // Not in allowed list
{"status": "PENDING"}     // Case sensitive
{"status": "in progress"} // Space vs underscore
```

### 2.3 Query Parameters vs JSON Body

**JSON Body (POST/PUT):**

```go
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}
```

Used in: `c.ShouldBindJSON(&req)`

**Query Parameters (GET):**

```go
type TodoQueryParams struct {
	Status string `form:"status" binding:"omitempty"`
}
```

Used in: `c.ShouldBindQuery(&params)`

**Example URLs:**

```
GET /api/v1/todos?status=pending&priority=high
GET /api/v1/todos?page=2&limit=20
```

### 2.4 Date Format Convention

```go
DueDate string `json:"due_date" binding:"omitempty"` // Format: YYYY-MM-DD
```

**Expected format:** `"2006-01-02"` (Go's reference date)

**Valid examples:**

```json
{"due_date": "2025-11-14"}
{"due_date": "2025-12-31"}
{"due_date": ""}  // Empty = no due date
```

**Invalid examples:**

```json
{"due_date": "14-11-2025"}    // Wrong format
{"due_date": "2025/11/14"}    // Wrong separator
{"due_date": "Nov 14, 2025"}  // Wrong format
```

---

## Step 3: DTO Validation Flow

### 3.1 How Gin Validates DTOs

```
1. Client sends JSON request
2. Handler calls c.ShouldBindJSON(&dto)
3. Gin validates:
   - JSON structure matches DTO
   - Required fields present
   - Field types correct
   - Validation tags pass
4. If valid: proceeds to service
5. If invalid: returns 400 Bad Request
```

### 3.2 Example Validation Process

**Request:**

```json
POST /api/v1/auth/register
{
  "username": "ab",
  "email": "invalid-email",
  "password": "123",
  "full_name": ""
}
```

**Validation checks:**

```go
Username string `json:"username" binding:"required,min=3,max=50"`
// ❌ Fails: "ab" is < 3 characters

Email string `json:"email" binding:"required,email,max=100"`
// ❌ Fails: "invalid-email" is not valid email format

Password string `json:"password" binding:"required,min=6"`
// ❌ Fails: "123" is < 6 characters

FullName string `json:"full_name" binding:"required,max=100"`
// ❌ Fails: empty string but required
```

**Response:**

```json
{
  "success": false,
  "message": "Invalid request data",
  "error": "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag"
}
```

---

## Step 4: Understanding Type Aliases

```go
type TodoCreateRequest = CreateTodoRequest
type TodoUpdateRequest = UpdateTodoRequest
```

**Why use aliases?**

- Backward compatibility
- Multiple names for same type
- No runtime cost

**Both are valid:**

```go
var req1 CreateTodoRequest
var req2 TodoCreateRequest  // Same type!
```

---

## Step 5: Verify DTOs

### 5.1 Check File Structure

```bash
$ ls -la internal/dto/    # Linux/Mac
$ dir internal\dto\       # Windows
```

**Expected output:**

```
todo_dto.go
user_dto.go
```

- [ ] Both files present

### 5.2 Verify DTOs Compile

```bash
$ go build internal/dto/user_dto.go internal/dto/todo_dto.go
```

**Expected:** No errors

- [ ] DTOs compile successfully

---

## Step 6: DTO Best Practices

### 6.1 Naming Conventions

```go
// ✅ Good
type CreateTodoRequest struct { }
type TodoResponse struct { }
type UpdateProfileRequest struct { }

// ❌ Avoid
type CreateTodo struct { }        // Ambiguous
type Todo struct { }              // Conflicts with model
type UpdateProfile struct { }     // Not clear if request/response
```

### 6.2 Request vs Response

**Request DTOs:**

- Contain validation tags
- May have required fields
- Used for input

**Response DTOs:**

- No validation tags needed
- Never contain sensitive data (passwords!)
- Used for output

### 6.3 Optional Fields Pattern

```go
// For creation - all required
type CreateTodoRequest struct {
	Title  string `json:"title" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// For updates - all optional (pointers)
type UpdateTodoRequest struct {
	Title  *string `json:"title" binding:"omitempty"`
	Status *string `json:"status" binding:"omitempty,oneof=..."`
}
```

---

## Step 7: Common Response Patterns

### 7.1 Success Response

```go
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
```

**Usage:**

```go
c.JSON(200, dto.SuccessResponse{
	Success: true,
	Message: "Todo created successfully",
	Data:    todoResponse,
})
```

**Output:**

```json
{
  "success": true,
  "message": "Todo created successfully",
  "data": {
    "id": 1,
    "title": "My Todo",
    ...
  }
}
```

### 7.2 Error Response

```go
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
```

**Usage:**

```go
c.JSON(400, dto.ErrorResponse{
	Success: false,
	Message: "Validation failed",
	Error:   err.Error(),
})
```

**Output:**

```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'RegisterRequest.Email' Error:Field validation..."
}
```

---

## Step 8: Commit Changes

### 8.1 Check Status

```bash
$ git status
```

- [ ] DTO files shown as untracked

### 8.2 Stage DTOs

```bash
$ git add internal/dto/
```

- [ ] Files staged

### 8.3 Commit

```bash
$ git commit -m "Add DTOs for requests and responses with validation"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/dto/user_dto.go` created
- [ ] RegisterRequest with validation tags
- [ ] LoginRequest defined
- [ ] UpdateProfileRequest with pointer fields
- [ ] UserResponse without password
- [ ] LoginResponse with token and user
- [ ] SuccessResponse and ErrorResponse structs
- [ ] `internal/dto/todo_dto.go` created
- [ ] CreateTodoRequest with oneof validation
- [ ] UpdateTodoRequest with optional fields
- [ ] TodoQueryParams for filtering
- [ ] TodoResponse defined
- [ ] Type aliases for backward compatibility
- [ ] Understanding of binding tags
- [ ] Understanding of pointer fields
- [ ] Understanding of validation flow
- [ ] DTOs compile without errors
- [ ] Changes committed to git

---

## 🧪 Quick Test

📝 **Create temporary file:** `test_dto.go`

```go
package main

import (
	"fmt"
	"encoding/json"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
)

func main() {
	// Test RegisterRequest
	reg := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	// Test CreateTodoRequest
	todo := dto.CreateTodoRequest{
		Title:    "Test Todo",
		Status:   "pending",
		Priority: "high",
	}

	// Convert to JSON
	regJSON, _ := json.MarshalIndent(reg, "", "  ")
	todoJSON, _ := json.MarshalIndent(todo, "", "  ")

	fmt.Println("RegisterRequest:")
	fmt.Println(string(regJSON))
	fmt.Println("\nCreateTodoRequest:")
	fmt.Println(string(todoJSON))
}
```

### Run Test

```bash
$ go run test_dto.go
```

**Expected output:**

```json
RegisterRequest:
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}

CreateTodoRequest:
{
  "title": "Test Todo",
  "description": "",
  "status": "pending",
  "priority": "high",
  "due_date": ""
}
```

- [ ] Test runs successfully
- [ ] JSON structure correct

### Clean Up

```bash
$ rm test_dto.go    # Linux/Mac
$ del test_dto.go   # Windows
```

---

## 🐛 Common Issues

### Issue: "binding tags not working"

**Solution:** Make sure you're using Gin's binding, not just JSON unmarshaling

### Issue: "required field shows as optional"

**Solution:** Check spelling: `binding:"required"` not `binding:"require"`

### Issue: "oneof validation not working"

**Solution:** Values are space-separated and case-sensitive

---

## 📊 Current Progress

```
[████████████████████░░] 31% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs  
**Next:** Utility functions

---

**Previous:** [04-models.md](04-models.md)  
**Next:** [06-utilities.md](06-utilities.md)
