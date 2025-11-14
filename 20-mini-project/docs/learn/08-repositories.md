# 08 - Repository Layer

Create data access layer with GORM repositories for users and todos.

---

## Overview

We'll create:

1. `internal/repository/user_repository.go` - User data access
2. `internal/repository/todo_repository.go` - Todo data access

Repositories encapsulate all database operations with specific query methods.

---

## Step 1: Create User Repository

### 1.1 Create user_repository.go File

📝 **Create file:** `internal/repository/user_repository.go`

```go
package repository

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create menyimpan user baru
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID mencari user berdasarkan ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername mencari user berdasarkan username
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update memperbarui data user
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// ExistsByUsername mengecek apakah username sudah ada
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail mengecek apakah email sudah ada
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
```

- [ ] File created at `internal/repository/user_repository.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Repository Pattern

**What is a repository?**

- Abstraction layer between business logic and database
- Encapsulates all SQL queries
- Makes testing easier (can mock repository)

**Structure:**

```go
type UserRepository struct {
	db *gorm.DB  // Database connection
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}  // Constructor
}
```

### 1.3 Understanding Repository Methods

#### Create Method

```go
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
```

**What it does:**

```sql
INSERT INTO users (username, email, password, full_name, created_at, updated_at)
VALUES ('john', 'john@example.com', '$2a$10$...', 'John Doe', NOW(), NOW())
```

**Usage in service:**

```go
user := &model.User{
	Username: "john",
	Email:    "john@example.com",
	Password: hashedPassword,
	FullName: "John Doe",
}
err := userRepo.Create(user)
```

#### FindByID Method

```go
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

**What it does:**

```sql
SELECT * FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1
```

**Returns:**

- `*model.User` if found
- `nil, error` if not found or database error

#### FindByUsername Method

```go
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

**Why use `Where()` instead of `First()`?**

- `First(&user, id)` - for primary key lookup
- `Where("field = ?", value)` - for other fields

**⚠️ SQL injection prevention:**

```go
// ✅ Safe (parameterized query)
Where("username = ?", username)

// ❌ Dangerous (SQL injection risk)
Where("username = '" + username + "'")
```

#### ExistsByUsername Method

```go
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
```

**What it does:**

```sql
SELECT COUNT(*) FROM users WHERE username = ? AND deleted_at IS NULL
```

**Why not use `Find()`?**

- More efficient - only counts, doesn't fetch data
- Faster for existence checks
- Returns `true/false` instead of full object

---

## Step 2: Create Todo Repository

### 2.1 Create todo_repository.go File

📝 **Create file:** `internal/repository/todo_repository.go`

```go
package repository

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/gorm"
)

// TodoRepository handles todo data access
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates new todo repository instance
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create menyimpan todo baru
func (r *TodoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

// FindByID mencari todo berdasarkan ID
func (r *TodoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// FindByUserID mengambil semua todo milik user
func (r *TodoRepository) FindByUserID(userID uint) ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// FindByUserIDWithFilters mengambil todo dengan filter
func (r *TodoRepository) FindByUserIDWithFilters(userID uint, status, priority string) ([]model.Todo, error) {
	var todos []model.Todo
	query := r.db.Where("user_id = ?", userID)

	// Filter status jika ada
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter priority jika ada
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	err := query.Order("created_at DESC").Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// Update memperbarui data todo
func (r *TodoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

// Delete menghapus todo (soft delete)
func (r *TodoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Todo{}, id).Error
}

// IsOwnedByUser mengecek apakah todo dimiliki oleh user
func (r *TodoRepository) IsOwnedByUser(todoID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Todo{}).
		Where("id = ? AND user_id = ?", todoID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
```

- [ ] File created at `internal/repository/todo_repository.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Todo Repository Methods

#### FindByUserID Method

```go
func (r *TodoRepository) FindByUserID(userID uint) ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}
```

**What it does:**

```sql
SELECT * FROM todos
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
```

**Why `Order("created_at DESC")`?**

- Shows newest todos first
- Better user experience
- Consistent ordering

#### FindByUserIDWithFilters Method

```go
func (r *TodoRepository) FindByUserIDWithFilters(userID uint, status, priority string) ([]model.Todo, error) {
	var todos []model.Todo
	query := r.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	err := query.Order("created_at DESC").Find(&todos).Error
	// ...
}
```

**What it does:**

```sql
-- No filters
SELECT * FROM todos WHERE user_id = ?

-- With status filter
SELECT * FROM todos WHERE user_id = ? AND status = 'pending'

-- With both filters
SELECT * FROM todos WHERE user_id = ? AND status = 'pending' AND priority = 'high'
```

**Query building pattern:**

```go
query := r.db.Where("user_id = ?", userID)  // Base query

if status != "" {
	query = query.Where("status = ?", status)  // Add condition
}

if priority != "" {
	query = query.Where("priority = ?", priority)  // Add another
}

err := query.Find(&todos).Error  // Execute final query
```

#### Delete Method (Soft Delete)

```go
func (r *TodoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Todo{}, id).Error
}
```

**What it does:**

```sql
-- Not actual DELETE!
UPDATE todos SET deleted_at = NOW() WHERE id = ?
```

**Why soft delete?**

- Data recovery possible
- Audit trail maintained
- Relationships preserved

**How it works:**

```go
// Model has DeletedAt field
type Todo struct {
	// ...
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// GORM automatically:
// - Sets deleted_at on Delete()
// - Filters deleted_at IS NULL on queries
```

#### IsOwnedByUser Method

```go
func (r *TodoRepository) IsOwnedByUser(todoID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Todo{}).
		Where("id = ? AND user_id = ?", todoID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
```

**Why this is important:**

- **Security**: Prevents users from accessing others' todos
- **Authorization**: Checks ownership before update/delete
- **Efficient**: Only counts, doesn't fetch data

**Usage in service:**

```go
// Before updating
owned, err := r.todoRepo.IsOwnedByUser(todoID, userID)
if !owned {
	return ErrUnauthorizedAccess
}

// Now safe to update
todo.Title = newTitle
r.todoRepo.Update(todo)
```

---

## Step 3: Understanding GORM Query Methods

### 3.1 GORM Method Comparison

| Method     | Use Case               | Returns | Example                       |
| ---------- | ---------------------- | ------- | ----------------------------- |
| `Create()` | Insert new record      | Error   | `db.Create(&user)`            |
| `First()`  | Find first matching    | Error   | `db.First(&user, id)`         |
| `Find()`   | Find all matching      | Error   | `db.Find(&users)`             |
| `Save()`   | Insert or update       | Error   | `db.Save(&user)`              |
| `Update()` | Update specific fields | Error   | `db.Update("name", "John")`   |
| `Delete()` | Soft delete            | Error   | `db.Delete(&user, id)`        |
| `Where()`  | Add WHERE condition    | \*DB    | `db.Where("age > ?", 18)`     |
| `Order()`  | Add ORDER BY           | \*DB    | `db.Order("created_at DESC")` |
| `Count()`  | Count records          | Error   | `db.Count(&count)`            |

### 3.2 Query Building Examples

#### Single Record Queries

```go
// By primary key
db.First(&user, 1)
// SELECT * FROM users WHERE id = 1 LIMIT 1

// By field
db.Where("username = ?", "john").First(&user)
// SELECT * FROM users WHERE username = 'john' LIMIT 1

// Multiple conditions
db.Where("age > ? AND city = ?", 18, "NYC").First(&user)
// SELECT * FROM users WHERE age > 18 AND city = 'NYC' LIMIT 1
```

#### Multiple Records Queries

```go
// All records
db.Find(&users)
// SELECT * FROM users WHERE deleted_at IS NULL

// With condition
db.Where("age > ?", 18).Find(&users)
// SELECT * FROM users WHERE age > 18 AND deleted_at IS NULL

// With ordering
db.Order("created_at DESC").Find(&users)
// SELECT * FROM users ORDER BY created_at DESC
```

#### Update Queries

```go
// Save entire struct
db.Save(&user)
// UPDATE users SET username=?, email=?, ... WHERE id = ?

// Update specific fields
db.Model(&user).Update("email", "new@example.com")
// UPDATE users SET email='new@example.com' WHERE id = ?

// Update multiple fields
db.Model(&user).Updates(map[string]interface{}{
	"email": "new@example.com",
	"name":  "New Name",
})
```

---

## Step 4: Repository Testing Concepts

### 4.1 Why Test Repositories?

- Verify database queries work correctly
- Catch SQL errors early
- Ensure data integrity
- Test edge cases (not found, duplicates)

### 4.2 Repository Test Structure

```go
func TestUserRepository(t *testing.T) {
	// Setup test database
	db := setupTestDB()
	defer db.Close()

	repo := repository.NewUserRepository(db)

	// Test Create
	user := &model.User{Username: "test"}
	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)  // ID should be set

	// Test FindByID
	found, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test", found.Username)

	// Test not found
	_, err = repo.FindByID(9999)
	assert.Error(t, err)
}
```

---

## Step 5: Verify Repositories

### 5.1 Check File Structure

```bash
$ ls -la internal/repository/    # Linux/Mac
$ dir internal\repository\       # Windows
```

**Expected output:**

```
todo_repository.go
user_repository.go
```

- [ ] Both files present

### 5.2 Verify Compilation

```bash
$ go build internal/repository/*.go
```

**Expected:** No errors

- [ ] Repositories compile successfully

### 5.3 Check Imports

```bash
$ go list -f '{{.Imports}}' github.com/adityapryg/golang-demo/20-mini-project/internal/repository
```

**Expected imports:**

- `github.com/adityapryg/golang-demo/20-mini-project/internal/model`
- `gorm.io/gorm`

- [ ] Imports verified

---

## Step 6: Commit Changes

### 6.1 Check Status

```bash
$ git status
```

- [ ] Repository files shown

### 6.2 Stage Repositories

```bash
$ git add internal/repository/
```

- [ ] Files staged

### 6.3 Commit

```bash
$ git commit -m "Add user and todo repositories with GORM queries"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/repository/user_repository.go` created
- [ ] UserRepository struct defined
- [ ] NewUserRepository constructor implemented
- [ ] Create method for inserting users
- [ ] FindByID method for primary key lookup
- [ ] FindByUsername method for username lookup
- [ ] FindByEmail method for email lookup
- [ ] Update method for updating users
- [ ] ExistsByUsername method for checking duplicates
- [ ] ExistsByEmail method for checking duplicates
- [ ] `internal/repository/todo_repository.go` created
- [ ] TodoRepository struct defined
- [ ] NewTodoRepository constructor implemented
- [ ] Create method for inserting todos
- [ ] FindByID method for getting single todo
- [ ] FindByUserID method for user's todos
- [ ] FindByUserIDWithFilters method for filtered queries
- [ ] Update method for updating todos
- [ ] Delete method for soft delete
- [ ] IsOwnedByUser method for ownership check
- [ ] Understanding of GORM query methods
- [ ] Understanding of soft delete
- [ ] Understanding of parameterized queries
- [ ] Repositories compile without errors
- [ ] Changes committed to git

---

## 🧪 Quick Test

We can verify repositories structure with a simple test:

📝 **Create temporary file:** `test_repository.go`

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// In-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto-migrate
	db.AutoMigrate(&model.User{}, &model.Todo{})

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Test user repository
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		FullName: "Test User",
	}
	err = userRepo.Create(user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User created with ID: %d\n", user.ID)

	// Test todo repository
	todo := &model.Todo{
		Title:    "Test Todo",
		Status:   "pending",
		Priority: "high",
		UserID:   user.ID,
	}
	err = todoRepo.Create(todo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Todo created with ID: %d\n", todo.ID)

	// Test queries
	found, err := userRepo.FindByUsername("testuser")
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User found: %s (%s)\n", found.Username, found.Email)

	todos, err := todoRepo.FindByUserID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Found %d todos for user\n", len(todos))

	// Test ownership
	owned, err := todoRepo.IsOwnedByUser(todo.ID, user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Ownership check: %v\n", owned)

	fmt.Println("\n✅ All repository methods working correctly!")
}
```

### Install SQLite Driver

```bash
$ go get gorm.io/driver/sqlite
```

### Run Test

```bash
$ go run test_repository.go
```

**Expected output:**

```
✅ User created with ID: 1
✅ Todo created with ID: 1
✅ User found: testuser (test@example.com)
✅ Found 1 todos for user
✅ Ownership check: true

✅ All repository methods working correctly!
```

- [ ] Test runs without errors
- [ ] All methods work correctly

### Clean Up

```bash
$ rm test_repository.go    # Linux/Mac
$ del test_repository.go   # Windows
```

---

## 🐛 Common Issues

### Issue: "record not found"

**Solution:** This is expected for `First()` when no matching record. Handle in service layer.

### Issue: "Soft delete not working"

**Solution:** Ensure model has `DeletedAt gorm.DeletedAt` field

### Issue: "Query returns deleted records"

**Solution:** GORM automatically filters deleted_at IS NULL. Check if you're using raw SQL.

### Issue: "Count returns wrong number"

**Solution:** Ensure you're using `Model(&Model{})` or passing struct type

---

## 📊 Current Progress

```
[█████████████████████████] 50% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories  
**Next:** Service layer with business logic

---

**Previous:** [07-middleware.md](07-middleware.md)  
**Next:** [09-services.md](09-services.md)
