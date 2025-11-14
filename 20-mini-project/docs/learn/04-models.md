# 04 - Database Models

Create GORM entity models for User and Todo.

---

## Overview

We'll create:

1. `internal/model/user.go` - User entity
2. `internal/model/todo.go` - Todo entity

These are the **database schema** represented as Go structs.

---

## Step 1: Create User Model

### 1.1 Create user.go File

📝 **Create file:** `internal/model/user.go`

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan entity pengguna di database
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null;size:50;index"`
	Email     string `gorm:"unique;not null;size:100;index"`
	Password  string `gorm:"not null"` // Hash password, jangan pernah expose ke luar
	FullName  string `gorm:"size:100"`
	Todos     []Todo `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName override nama tabel
func (User) TableName() string {
	return "users"
}
```

- [ ] File created at `internal/model/user.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding User Model

💡 **Field Breakdown:**

| Field       | Type           | GORM Tag                          | Description                                |
| ----------- | -------------- | --------------------------------- | ------------------------------------------ |
| `ID`        | uint           | primaryKey                        | Auto-increment primary key                 |
| `Username`  | string         | unique, not null, size:50, index  | Unique username, max 50 chars, indexed     |
| `Email`     | string         | unique, not null, size:100, index | Unique email, max 100 chars, indexed       |
| `Password`  | string         | not null                          | Bcrypt hashed password (never plain text!) |
| `FullName`  | string         | size:100                          | User's full name, max 100 chars, optional  |
| `Todos`     | []Todo         | foreignKey:UserID                 | One-to-many relationship with todos        |
| `CreatedAt` | time.Time      | (auto)                            | Automatically set when record created      |
| `UpdatedAt` | time.Time      | (auto)                            | Automatically updated when record modified |
| `DeletedAt` | gorm.DeletedAt | index                             | Soft delete support, indexed               |

💡 **GORM Tags Explained:**

- `primaryKey` - Makes field the primary key
- `unique` - Enforces uniqueness constraint
- `not null` - Field cannot be NULL
- `size:N` - Sets VARCHAR(N) size
- `index` - Creates database index for faster queries
- `foreignKey:UserID` - Defines foreign key relationship

💡 **Soft Delete:**

- `DeletedAt` field enables soft deletes
- Records are marked as deleted, not removed
- GORM automatically filters soft-deleted records
- Can be recovered if needed

### 1.3 TableName() Method

```go
func (User) TableName() string {
	return "users"
}
```

**Why?**

- By default, GORM pluralizes struct name (User → users)
- This explicitly sets table name
- Good practice for clarity
- Prevents naming issues

---

## Step 2: Create Todo Model

### 2.1 Create todo.go File

📝 **Create file:** `internal/model/todo.go`

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// Todo merepresentasikan entity todo di database
type Todo struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null;size:200"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);default:'pending';index"`
	Priority    string `gorm:"type:varchar(10);default:'medium'"`
	DueDate     *time.Time
	UserID      uint `gorm:"not null;index"`
	User        User `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName override nama tabel
func (Todo) TableName() string {
	return "todos"
}
```

- [ ] File created at `internal/model/todo.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Todo Model

💡 **Field Breakdown:**

| Field         | Type           | GORM Tag                              | Description                            |
| ------------- | -------------- | ------------------------------------- | -------------------------------------- |
| `ID`          | uint           | primaryKey                            | Auto-increment primary key             |
| `Title`       | string         | not null, size:200                    | Todo title, max 200 chars, required    |
| `Description` | string         | type:text                             | Long text description, optional        |
| `Status`      | string         | varchar(20), default:'pending', index | Todo status, default 'pending'         |
| `Priority`    | string         | varchar(10), default:'medium'         | Priority level, default 'medium'       |
| `DueDate`     | \*time.Time    | (none)                                | Optional due date (pointer = nullable) |
| `UserID`      | uint           | not null, index                       | Foreign key to users table             |
| `User`        | User           | foreignKey:UserID                     | Belongs-to relationship                |
| `CreatedAt`   | time.Time      | (auto)                                | Creation timestamp                     |
| `UpdatedAt`   | time.Time      | (auto)                                | Last update timestamp                  |
| `DeletedAt`   | gorm.DeletedAt | index                                 | Soft delete timestamp                  |

💡 **Special Field Types:**

**Text vs VARCHAR:**

```go
Description string `gorm:"type:text"`     // Unlimited length
Title       string `gorm:"size:200"`     // VARCHAR(200)
Status      string `gorm:"type:varchar(20)"` // Explicit VARCHAR(20)
```

**Nullable Fields:**

```go
DueDate *time.Time  // Pointer = can be NULL
UserID  uint        // Not a pointer = NOT NULL
```

💡 **Status Values:**

- `pending` - Not started
- `in_progress` - Currently working on it
- `completed` - Finished

💡 **Priority Values:**

- `low` - Low priority
- `medium` - Medium priority (default)
- `high` - High priority

### 2.3 Relationships Explained

**User has many Todos:**

```go
// In User model
Todos []Todo `gorm:"foreignKey:UserID"`
```

**Todo belongs to User:**

```go
// In Todo model
UserID uint `gorm:"not null;index"`
User   User `gorm:"foreignKey:UserID"`
```

This creates a **one-to-many** relationship:

```
User (1) ←→ (many) Todos
```

---

## Step 3: Verify Models

### 3.1 Check File Structure

```bash
$ ls -la internal/model/    # Linux/Mac
$ dir internal\model\       # Windows
```

**Expected output:**

```
todo.go
user.go
```

- [ ] Both files present

### 3.2 Verify Models Compile

```bash
$ go build internal/model/user.go internal/model/todo.go
```

⚠️ **Expected:** Error about circular import or missing main

Try this instead:

```bash
$ cd internal/model
$ go build .
$ cd ../..
```

**Expected:** No output = success!

- [ ] Models compile without errors

---

## Step 4: Test Database Configuration

Now that models exist, let's test if database.go works:

### 4.1 Check database.go Compiles

```bash
$ go build internal/config/database.go
```

**Expected:** No errors (models now exist)

- [ ] database.go compiles successfully

---

## Step 5: Understanding Database Schema

### 5.1 Resulting Database Tables

When `AutoMigrate` runs, it creates:

**users table:**

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

**todos table:**

```sql
CREATE TABLE todos (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    priority VARCHAR(10) DEFAULT 'medium',
    due_date TIMESTAMP NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_deleted_at ON todos(deleted_at);
```

### 5.2 Indexes Explained

**Why indexes?**

- Speed up queries on indexed columns
- Username/Email: for login lookups
- Status: for filtering todos
- UserID: for joining with users
- DeletedAt: for filtering soft-deleted records

**Cost:**

- Slightly slower writes (must update index)
- More disk space
- Worth it for read-heavy operations

---

## Step 6: Model Validation Rules

### 6.1 Database-Level Constraints

GORM tags create database constraints:

```go
`gorm:"unique"`          // UNIQUE constraint
`gorm:"not null"`        // NOT NULL constraint
`gorm:"size:50"`         // VARCHAR(50) size limit
`gorm:"default:'value'"` // DEFAULT value
```

### 6.2 Application-Level Validation

We'll add validation later in:

- **DTOs** - Request validation (binding tags)
- **Services** - Business logic validation

---

## Step 7: Commit Changes

### 7.1 Check Status

```bash
$ git status
```

**Expected:**

```
Untracked files:
  internal/model/user.go
  internal/model/todo.go
```

- [ ] Model files shown

### 7.2 Stage Models

```bash
$ git add internal/model/
```

- [ ] Files staged

### 7.3 Commit

```bash
$ git commit -m "Add User and Todo database models"
```

**Expected output:**

```
[main xxxxxxx] Add User and Todo database models
 2 files changed, XX insertions(+)
 create mode 100644 internal/model/todo.go
 create mode 100644 internal/model/user.go
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/model/user.go` created with User struct
- [ ] User model has all required fields
- [ ] User model has proper GORM tags
- [ ] User model has soft delete support
- [ ] User TableName() returns "users"
- [ ] `internal/model/todo.go` created with Todo struct
- [ ] Todo model has all required fields
- [ ] Todo model has proper GORM tags
- [ ] Todo model has soft delete support
- [ ] Todo TableName() returns "todos"
- [ ] One-to-many relationship configured
- [ ] Models compile without errors
- [ ] database.go now compiles (models exist)
- [ ] Understanding of GORM tags
- [ ] Understanding of relationships
- [ ] Changes committed to git

---

## 🧪 Quick Test

Create a test to verify models work:

📝 **Create temporary file:** `test_models.go` (in project root)

```go
package main

import (
	"fmt"
	"time"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
)

func main() {
	// Create a user
	user := model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		FullName: "Test User",
	}

	// Create a todo
	todo := model.Todo{
		Title:       "Test Todo",
		Description: "This is a test",
		Status:      "pending",
		Priority:    "medium",
		UserID:      1,
	}

	fmt.Printf("User struct created: %s (%s)\n", user.Username, user.Email)
	fmt.Printf("Todo struct created: %s\n", todo.Title)
	fmt.Printf("Models work correctly!\n")
}
```

### Run Test

```bash
$ go run test_models.go
```

**Expected output:**

```
User struct created: testuser (test@example.com)
Todo struct created: Test Todo
Models work correctly!
```

- [ ] Test runs successfully
- [ ] No compilation errors

### Clean Up

```bash
$ rm test_models.go    # Linux/Mac
$ del test_models.go   # Windows
```

- [ ] Test file removed

---

## 🐛 Common Issues

### Issue: "undefined: gorm.DeletedAt"

**Solution:**

- Import missing: `"gorm.io/gorm"`
- Run `go mod tidy`

### Issue: Circular import error

**Solution:**

- Don't import model package into itself
- Build from parent directory

### Issue: "cannot use User literal as type"

**Solution:**

- Check struct field names match exactly
- Ensure all imports are correct

---

## 📊 Current Progress

```
[████████████████░░] 25% Complete
```

**Completed:** Setup, dependencies, config, models  
**Next:** Data Transfer Objects (DTOs)

---

**Previous:** [03-configuration.md](03-configuration.md)  
**Next:** [05-dtos.md](05-dtos.md)
