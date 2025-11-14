# 13 - Database Migrations

Create SQL migration files and understand database schema management.

---

## Overview

We'll create:

1. `migrations/001_create_tables.sql` - Initial schema with users and todos tables

Migrations define database schema and can be applied/rolled back systematically.

---

## Step 1: Understanding Database Migrations

### 1.1 What are Migrations?

**Definition:**
Database migrations are version-controlled changes to your database schema.

**Why use migrations?**

- **Version control**: Track database changes like code
- **Reproducible**: Same schema across environments
- **Team collaboration**: Everyone has same database structure
- **Rollback**: Undo changes if needed
- **Documentation**: Schema changes are documented

**Migration lifecycle:**

```
Development → Create migration
           → Apply to dev DB
           → Test
           → Commit to git
           → Apply to staging DB
           → Apply to production DB
```

### 1.2 Migration Strategies

**Option 1: Manual SQL files (our choice)**

```sql
-- migrations/001_create_tables.sql
CREATE TABLE users (...);
CREATE TABLE todos (...);

-- Apply manually:
psql -U postgres -d todolist_db -f migrations/001_create_tables.sql
```

**Benefits:**

- Full control over SQL
- No dependencies
- Clear and explicit
- Easy to review

**Option 2: GORM Auto-Migrate**

```go
db.AutoMigrate(&model.User{}, &model.Todo{})
```

**Benefits:**

- Automatic schema generation
- No SQL needed
- Works across databases

**Drawbacks:**

- Less control
- Can't handle complex changes
- No rollback

**Option 3: Migration tools (migrate, goose, etc.)**

```bash
migrate -path ./migrations -database postgres://... up
```

**Benefits:**

- Up/down migrations
- Version tracking
- Rollback support

---

## Step 2: Create Migration File

### 2.1 Create 001_create_tables.sql

📝 **Create file:** `migrations/001_create_tables.sql`

```sql
-- Migration: Create users and todos tables
-- Version: 001
-- Description: Initial schema for Todo REST API

-- ============================================
-- Create users table
-- ============================================

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Add comments for documentation
COMMENT ON TABLE users IS 'User accounts for authentication';
COMMENT ON COLUMN users.id IS 'Primary key, auto-increment';
COMMENT ON COLUMN users.username IS 'Unique username for login';
COMMENT ON COLUMN users.email IS 'Unique email address';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.full_name IS 'User''s full name';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp (NULL = not deleted)';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ============================================
-- Create todos table
-- ============================================

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    due_date DATE,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign key constraint
    CONSTRAINT fk_todos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Add comments for documentation
COMMENT ON TABLE todos IS 'Todo items belonging to users';
COMMENT ON COLUMN todos.id IS 'Primary key, auto-increment';
COMMENT ON COLUMN todos.title IS 'Todo title (max 200 chars)';
COMMENT ON COLUMN todos.description IS 'Optional detailed description';
COMMENT ON COLUMN todos.status IS 'Status: pending, in_progress, completed';
COMMENT ON COLUMN todos.priority IS 'Priority: low, medium, high';
COMMENT ON COLUMN todos.due_date IS 'Optional due date';
COMMENT ON COLUMN todos.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN todos.deleted_at IS 'Soft delete timestamp (NULL = not deleted)';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_deleted_at ON todos(deleted_at);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);

-- Composite index for common query (user's todos by status)
CREATE INDEX IF NOT EXISTS idx_todos_user_status ON todos(user_id, status);

-- ============================================
-- Add constraints for data validation
-- ============================================

-- Check constraint for status values
ALTER TABLE todos
    ADD CONSTRAINT check_todos_status
    CHECK (status IN ('pending', 'in_progress', 'completed'));

-- Check constraint for priority values
ALTER TABLE todos
    ADD CONSTRAINT check_todos_priority
    CHECK (priority IN ('low', 'medium', 'high'));

-- ============================================
-- Create updated_at trigger function
-- ============================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for todos table
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
CREATE TRIGGER update_todos_updated_at
    BEFORE UPDATE ON todos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Migration complete
-- ============================================

-- Display success message
DO $$
BEGIN
    RAISE NOTICE 'Migration 001 completed successfully';
    RAISE NOTICE 'Created tables: users, todos';
    RAISE NOTICE 'Created indexes for performance';
    RAISE NOTICE 'Created triggers for updated_at';
END $$;
```

- [ ] File created at `migrations/001_create_tables.sql`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding the Schema

#### Users Table Structure

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,              -- Auto-increment primary key
    username VARCHAR(50) UNIQUE NOT NULL,  -- Max 50 chars, must be unique
    email VARCHAR(100) UNIQUE NOT NULL,    -- Max 100 chars, must be unique
    password VARCHAR(255) NOT NULL,        -- Bcrypt hash (~60 chars)
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP                   -- Soft delete (NULL = active)
);
```

**Column types explained:**

| Type         | Usage                  | Size     | Example                    |
| ------------ | ---------------------- | -------- | -------------------------- |
| `SERIAL`     | Auto-increment integer | 4 bytes  | 1, 2, 3...                 |
| `VARCHAR(n)` | Variable-length string | n chars  | "john", "test@example.com" |
| `TEXT`       | Unlimited text         | Variable | Long descriptions          |
| `TIMESTAMP`  | Date + time            | 8 bytes  | 2024-03-15 10:30:00        |
| `DATE`       | Date only              | 4 bytes  | 2024-03-15                 |
| `INTEGER`    | Whole number           | 4 bytes  | 1, 100, 9999               |

#### Todos Table Structure

```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,                      -- Optional, can be NULL
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    due_date DATE,                         -- Optional
    user_id INTEGER NOT NULL,              -- Foreign key
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign key relationship
    CONSTRAINT fk_todos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE    -- Delete todos when user is deleted
        ON UPDATE CASCADE    -- Update user_id if users.id changes
);
```

### 2.3 Understanding Indexes

```sql
-- Single column index
CREATE INDEX idx_users_username ON users(username);

-- Composite index (multiple columns)
CREATE INDEX idx_todos_user_status ON todos(user_id, status);
```

**Why indexes?**

- Speed up queries (like book index)
- Trade-off: Slower writes, faster reads
- Essential for foreign keys and WHERE clauses

**When to index:**

```sql
-- ✅ Index these
WHERE user_id = 1           → Index user_id
WHERE status = 'pending'    → Index status
ORDER BY created_at DESC    → Index created_at

-- ❌ Don't index these
SELECT *                    → No filter
WHERE description LIKE '%x%' → Full text search needed
Small tables (<1000 rows)   → Index overhead not worth it
```

**Index types:**

| Index Type       | Command                         | Use Case                |
| ---------------- | ------------------------------- | ----------------------- |
| B-Tree (default) | `CREATE INDEX`                  | Equality, range queries |
| Unique           | `CREATE UNIQUE INDEX`           | Enforce uniqueness      |
| Partial          | `CREATE INDEX ... WHERE`        | Index subset of rows    |
| Composite        | `CREATE INDEX ... (col1, col2)` | Multi-column queries    |

### 2.4 Understanding Constraints

```sql
-- Check constraint (validation at DB level)
ALTER TABLE todos
    ADD CONSTRAINT check_todos_status
    CHECK (status IN ('pending', 'in_progress', 'completed'));
```

**Why use CHECK constraints?**

- Database-level validation
- Can't be bypassed by code
- Protects data integrity
- Self-documenting

**Constraint types:**

| Type          | Example                                | Purpose               |
| ------------- | -------------------------------------- | --------------------- |
| `NOT NULL`    | `email VARCHAR(100) NOT NULL`          | Prevent NULL values   |
| `UNIQUE`      | `username VARCHAR(50) UNIQUE`          | No duplicates         |
| `PRIMARY KEY` | `id SERIAL PRIMARY KEY`                | Unique identifier     |
| `FOREIGN KEY` | `REFERENCES users(id)`                 | Referential integrity |
| `CHECK`       | `CHECK (status IN ('a','b'))`          | Custom validation     |
| `DEFAULT`     | `status VARCHAR(20) DEFAULT 'pending'` | Default value         |

### 2.5 Understanding Triggers

```sql
-- Function to update updated_at
CREATE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger that calls function
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**How triggers work:**

```
1. UPDATE users SET email = 'new@example.com' WHERE id = 1
2. Trigger fires BEFORE UPDATE
3. Function sets NEW.updated_at = CURRENT_TIMESTAMP
4. Row updated with new timestamp
```

**Why use triggers?**

- Automatic timestamp updates
- Audit logging
- Data validation
- Consistency enforcement

---

## Step 3: Apply Migration

### 3.1 Using psql Command

```bash
# Connect to database and apply migration
$ psql -U postgres -d todolist_db -f migrations/001_create_tables.sql
```

**Expected output:**

```
CREATE TABLE
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE TABLE
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
ALTER TABLE
ALTER TABLE
CREATE FUNCTION
DROP TRIGGER
CREATE TRIGGER
DROP TRIGGER
CREATE TRIGGER
DO
NOTICE:  Migration 001 completed successfully
NOTICE:  Created tables: users, todos
NOTICE:  Created indexes for performance
NOTICE:  Created triggers for updated_at
```

- [ ] Migration applied successfully
- [ ] No errors

### 3.2 Using Docker Exec

If using Docker:

```bash
# Copy migration file to container
$ docker cp migrations/001_create_tables.sql 20-mini-project-postgres-1:/tmp/

# Execute migration inside container
$ docker exec -it 20-mini-project-postgres-1 psql -U postgres -d todolist_db -f /tmp/001_create_tables.sql
```

- [ ] Migration copied to container
- [ ] Migration applied in container

### 3.3 Verify Schema

```bash
# Connect to database
$ psql -U postgres -d todolist_db

# List tables
\dt

# Expected output:
         List of relations
 Schema |  Name  | Type  |  Owner
--------+--------+-------+----------
 public | todos  | table | postgres
 public | users  | table | postgres
```

**Check table structure:**

```sql
\d users

-- Expected:
                                      Table "public.users"
   Column   |            Type             | Collation | Nullable |              Default
------------+-----------------------------+-----------+----------+-----------------------------------
 id         | integer                     |           | not null | nextval('users_id_seq'::regclass)
 username   | character varying(50)       |           | not null |
 email      | character varying(100)      |           | not null |
 password   | character varying(255)      |           | not null |
 full_name  | character varying(100)      |           | not null |
 created_at | timestamp without time zone |           | not null | CURRENT_TIMESTAMP
 updated_at | timestamp without time zone |           | not null | CURRENT_TIMESTAMP
 deleted_at | timestamp without time zone |           |          |
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)
    "idx_users_deleted_at" btree (deleted_at)
    "idx_users_email" btree (email)
    "idx_users_username" btree (username)
```

- [ ] Tables created
- [ ] Indexes present
- [ ] Constraints active

---

## Step 4: Optional Seed Data

### 4.1 Create Seed File

📝 **Create file:** `migrations/002_seed_data.sql`

```sql
-- Seed data for testing
-- NOTE: This is for development only!

-- Insert test user
-- Password: "password123" (bcrypt hash)
INSERT INTO users (username, email, password, full_name)
VALUES (
    'testuser',
    'test@example.com',
    '$2a$10$YourBcryptHashHere',  -- Replace with actual hash
    'Test User'
) ON CONFLICT (username) DO NOTHING;

-- Get user ID for reference
DO $$
DECLARE
    test_user_id INTEGER;
BEGIN
    SELECT id INTO test_user_id FROM users WHERE username = 'testuser';

    -- Insert sample todos
    INSERT INTO todos (title, description, status, priority, user_id)
    VALUES
        ('Complete project documentation', 'Write comprehensive README and API docs', 'pending', 'high', test_user_id),
        ('Review pull requests', 'Review and merge pending PRs', 'in_progress', 'medium', test_user_id),
        ('Update dependencies', 'Update Go modules to latest versions', 'pending', 'low', test_user_id);

    RAISE NOTICE 'Seed data inserted successfully';
END $$;
```

**Generate bcrypt hash for seed data:**

```go
// Quick script to generate hash
package main
import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

- [ ] Seed file created (optional)
- [ ] Password hash generated (optional)

---

## Step 5: Migration Best Practices

### 5.1 Naming Conventions

```
migrations/
├── 001_create_tables.sql
├── 002_add_user_roles.sql
├── 003_add_todos_tags.sql
└── 004_alter_todos_add_priority.sql
```

**Format:**

- `NNN_description.sql` (3-digit number + description)
- Sequential numbering
- Descriptive names
- Use underscores

### 5.2 Migration Rules

**✅ Do:**

- Keep migrations small and focused
- Test migrations before committing
- Add comments explaining WHY
- Make migrations idempotent (can run multiple times)
- Use `IF NOT EXISTS` for safety

**❌ Don't:**

- Modify existing migrations (create new ones)
- Mix DDL and DML (schema + data)
- Delete old migrations
- Assume migration order
- Hardcode values

### 5.3 Rollback Strategy

**Create down migration:**

📝 **Create file:** `migrations/001_create_tables_down.sql`

```sql
-- Rollback migration 001

-- Drop triggers
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (CASCADE removes dependent objects)
DROP TABLE IF EXISTS todos CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Migration 001 rolled back successfully';
END $$;
```

**Apply rollback:**

```bash
$ psql -U postgres -d todolist_db -f migrations/001_create_tables_down.sql
```

---

## Step 6: GORM Auto-Migrate Alternative

If you prefer GORM's auto-migration instead:

### 6.1 Update database.go

```go
func NewDatabase(cfg *Config) *gorm.DB {
    dsn := fmt.Sprintf(...)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto-migrate models
    err = db.AutoMigrate(&model.User{}, &model.Todo{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    return db
}
```

**Pros:**

- Automatic schema generation
- No manual SQL
- Cross-database compatible

**Cons:**

- Less control over indexes
- Can't add custom constraints easily
- No rollback mechanism
- Harder to review changes

---

## Step 7: Verify Migration

### 7.1 Check Tables Exist

```bash
$ psql -U postgres -d todolist_db -c "\dt"
```

**Expected:**

```
         List of relations
 Schema |  Name  | Type  |  Owner
--------+--------+-------+----------
 public | todos  | table | postgres
 public | users  | table | postgres
```

- [ ] Both tables exist

### 7.2 Test Insert

```sql
-- Test user insert
INSERT INTO users (username, email, password, full_name)
VALUES ('testuser', 'test@example.com', 'hashedpass', 'Test User');

-- Test todo insert
INSERT INTO todos (title, status, priority, user_id)
VALUES ('Test Todo', 'pending', 'high', 1);

-- Verify
SELECT * FROM users;
SELECT * FROM todos;
```

- [ ] Inserts work
- [ ] Data retrieved successfully

### 7.3 Test Constraints

```sql
-- Test unique constraint
INSERT INTO users (username, email, password, full_name)
VALUES ('testuser', 'duplicate@example.com', 'pass', 'User');
-- Expected: ERROR:  duplicate key value violates unique constraint "users_username_key"

-- Test check constraint
INSERT INTO todos (title, status, priority, user_id)
VALUES ('Test', 'invalid_status', 'high', 1);
-- Expected: ERROR:  new row for relation "todos" violates check constraint "check_todos_status"

-- Test foreign key
INSERT INTO todos (title, status, priority, user_id)
VALUES ('Test', 'pending', 'high', 9999);
-- Expected: ERROR:  insert or update on table "todos" violates foreign key constraint "fk_todos_user"
```

- [ ] Constraints enforced

---

## Step 8: Commit Changes

### 8.1 Check Status

```bash
$ git status
```

- [ ] Migration files shown

### 8.2 Stage Migrations

```bash
$ git add migrations/
```

- [ ] Files staged

### 8.3 Commit

```bash
$ git commit -m "Add database migration with users and todos schema"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `migrations/001_create_tables.sql` created
- [ ] Users table schema defined
- [ ] Todos table schema defined
- [ ] Primary keys configured
- [ ] Foreign key relationship established
- [ ] Unique constraints on username and email
- [ ] Check constraints for status and priority
- [ ] Indexes created for performance
- [ ] Composite index for common queries
- [ ] Table comments added
- [ ] Column comments added
- [ ] Trigger function created
- [ ] Triggers for updated_at configured
- [ ] ON DELETE CASCADE configured
- [ ] Soft delete support (deleted_at)
- [ ] Understanding of migrations
- [ ] Understanding of indexes
- [ ] Understanding of constraints
- [ ] Understanding of triggers
- [ ] Migration applied successfully
- [ ] Schema verified in database
- [ ] Constraints tested
- [ ] Changes committed to git

---

## 🐛 Common Issues

### Issue: "relation already exists"

**Solution:** Tables already exist. Drop them first or use `IF NOT EXISTS`

### Issue: "permission denied"

**Solution:** Check PostgreSQL user has CREATE privilege

### Issue: "foreign key constraint violated"

**Solution:** Insert parent record (user) before child (todo)

### Issue: "constraint violation"

**Solution:** Data doesn't match CHECK constraint (e.g., invalid status)

---

## 📊 Current Progress

```
[████████████████████████████████████] 81% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers, routes, main app, migrations  
**Next:** Docker configuration

---

**Previous:** [12-main-app.md](12-main-app.md)  
**Next:** [14-docker.md](14-docker.md)
