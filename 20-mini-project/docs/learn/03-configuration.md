# 03 - Configuration Files

Create configuration layer for environment variables and database connection.

---

## Overview

We'll create:

1. `internal/config/config.go` - Environment variable configuration
2. `internal/config/database.go` - Database connection setup

---

## Step 1: Create Configuration Structure

### 1.1 Create config.go File

📝 **Create file:** `internal/config/config.go`

Open your code editor and create the file with this **exact** content:

```go
package config

import "os"

// Config menyimpan konfigurasi aplikasi
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
	GinMode    string
}

// LoadConfig memuat konfigurasi dari environment variables
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "todolist_db"),
		JWTSecret:  getEnv("JWT_SECRET", "your-super-secret-key-change-this-in-production"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}
}

// getEnv mendapatkan environment variable dengan default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
```

- [ ] File created at `internal/config/config.go`
- [ ] Content copied exactly as shown
- [ ] File saved

### 1.2 Verify config.go Content

```bash
$ cat internal/config/config.go    # Linux/Mac
$ type internal\config\config.go   # Windows
```

- [ ] Content matches above exactly
- [ ] No typos or missing characters

💡 **Code Explanation:**

**Config struct:**

- Holds all configuration values
- Fields use PascalCase (Go convention)
- All fields are strings

**LoadConfig() function:**

- Creates new Config instance
- Calls `getEnv()` for each setting
- Returns pointer to Config

**getEnv() function:**

- Reads environment variable
- Returns default if not set
- Makes configuration flexible

---

## Step 2: Create Database Configuration

### 2.1 Create database.go File

📝 **Create file:** `internal/config/database.go`

```go
package config

import (
	"fmt"
	"log"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *Config) (*gorm.DB, error) {
	// Connection string PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Koneksi ke database dengan logging
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("✓ Successfully connected to database")

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✓ Database migration completed")

	return db, nil
}
```

- [ ] File created at `internal/config/database.go`
- [ ] Content copied exactly
- [ ] File saved

⚠️ **Important:** This file will show errors because `model.User` and `model.Todo` don't exist yet. This is **normal** and expected. We'll create those in the next section.

### 2.2 Verify database.go Content

```bash
$ cat internal/config/database.go    # Linux/Mac
$ type internal\config\database.go   # Windows
```

- [ ] Content matches above
- [ ] Package imports are correct

💡 **Code Explanation:**

**DSN (Data Source Name):**

- Connection string for PostgreSQL
- Contains host, user, password, database name, port
- `sslmode=disable` for local development
- `TimeZone=Asia/Jakarta` sets timezone

**GORM Configuration:**

- `logger.Info` mode shows SQL queries (helpful for debugging)
- Can change to `logger.Silent` in production

**AutoMigrate:**

- Creates/updates tables based on models
- Automatically handles schema changes
- Will execute when we run the app

---

## Step 3: Verify Configuration Layer

### 3.1 Check File Structure

```bash
$ ls -la internal/config/    # Linux/Mac
$ dir internal\config\       # Windows
```

**Expected output:**

```
config.go
database.go
```

- [ ] Both files present
- [ ] No extra files

### 3.2 Check Imports Compile (Will Have Errors)

```bash
$ go build internal/config/config.go
```

**Expected:** No errors (config.go is standalone)

- [ ] config.go builds successfully

```bash
$ go build internal/config/database.go
```

**Expected:** Errors about missing `model.User` and `model.Todo`

```
package command-line-arguments
	internal/config/database.go:7:2: no required module provides package github.com/adityapryg/golang-demo/20-mini-project/internal/model; to add it:
```

⚠️ **This is NORMAL** - we haven't created the models yet!

- [ ] Error mentions missing models (expected)

---

## Step 4: Understanding Environment Variables

### 4.1 How It Works

The configuration system works like this:

```
1. App starts
2. LoadConfig() called
3. Checks environment variables
4. Falls back to defaults if not set
5. Returns Config struct
```

### 4.2 Environment Variable Priority

```
Environment Variable > Default Value
```

**Example:**

- If `DB_HOST` is set in environment: uses that value
- If `DB_HOST` is not set: uses "localhost"

### 4.3 Testing Configuration (Conceptual)

You can test environment variables like this:

**Linux/Mac:**

```bash
$ export DB_HOST=192.168.1.100
$ export JWT_SECRET=my-secret-key
```

**Windows (Command Prompt):**

```bash
$ set DB_HOST=192.168.1.100
$ set JWT_SECRET=my-secret-key
```

**Windows (PowerShell):**

```powershell
$ $env:DB_HOST="192.168.1.100"
$ $env:JWT_SECRET="my-secret-key"
```

💡 These environment variables are temporary (current session only)

---

## Step 5: Create .env File for Development

### 5.1 Create .env File

📝 **Create file:** `.env`

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist_db

# JWT Configuration
JWT_SECRET=dev-secret-key-do-not-use-in-production-12345

# Server Configuration
SERVER_PORT=8080

# Gin Mode (debug or release)
GIN_MODE=debug
```

- [ ] File `.env` created
- [ ] Content added
- [ ] File saved

⚠️ **Security Note:**

- `.env` is in `.gitignore`
- Never commit `.env` to git
- Use `.env.example` as template for others

### 5.2 Verify .env is Ignored by Git

```bash
$ git status
```

**Expected:** `.env` should **NOT** appear in untracked files

- [ ] `.env` is ignored by git

💡 If `.env` appears in git status:

1. Check `.gitignore` contains `.env`
2. Run `git rm --cached .env` if accidentally tracked

---

## Step 6: Configuration Constants Explained

### 6.1 Database Configuration

| Variable      | Default     | Description             |
| ------------- | ----------- | ----------------------- |
| `DB_HOST`     | localhost   | Database server address |
| `DB_PORT`     | 5432        | PostgreSQL default port |
| `DB_USER`     | postgres    | Database username       |
| `DB_PASSWORD` | postgres    | Database password       |
| `DB_NAME`     | todolist_db | Database name           |

### 6.2 JWT Configuration

| Variable     | Default       | Description                       |
| ------------ | ------------- | --------------------------------- |
| `JWT_SECRET` | (long string) | Secret key for signing JWT tokens |

⚠️ **Important:**

- **MUST** change in production
- Minimum 32 characters
- Use random, unguessable string
- Keep it secret!

### 6.3 Server Configuration

| Variable      | Default | Description        |
| ------------- | ------- | ------------------ |
| `SERVER_PORT` | 8080    | HTTP server port   |
| `GIN_MODE`    | debug   | Gin framework mode |

**GIN_MODE Options:**

- `debug` - Verbose logging, helpful errors
- `release` - Minimal logging, production use

---

## Step 7: Commit Changes

### 7.1 Check Status

```bash
$ git status
```

**Expected output:**

```
Untracked files:
  internal/config/config.go
  internal/config/database.go
```

- [ ] Only config files shown (not .env)

### 7.2 Stage Configuration Files

```bash
$ git add internal/config/
```

- [ ] Files staged

### 7.3 Commit

```bash
$ git commit -m "Add configuration layer with environment variables"
```

**Expected output:**

```
[main xxxxxxx] Add configuration layer with environment variables
 2 files changed, XX insertions(+)
 create mode 100644 internal/config/config.go
 create mode 100644 internal/config/database.go
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/config/config.go` created with Config struct
- [ ] `LoadConfig()` function reads environment variables
- [ ] `getEnv()` helper function provides defaults
- [ ] `internal/config/database.go` created with NewDatabase()
- [ ] PostgreSQL DSN string properly formatted
- [ ] GORM configuration with logging enabled
- [ ] AutoMigrate setup (will work once models exist)
- [ ] `.env` file created for development
- [ ] `.env` is ignored by git
- [ ] Environment variables understood
- [ ] Changes committed to git

---

## 🧪 Quick Verification

Create a test file to verify config works:

📝 **Create temporary file:** `test_config.go` (in project root)

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("Config loaded:\n")
	fmt.Printf("  DB Host: %s\n", cfg.DBHost)
	fmt.Printf("  DB Port: %s\n", cfg.DBPort)
	fmt.Printf("  DB Name: %s\n", cfg.DBName)
	fmt.Printf("  Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("  Gin Mode: %s\n", cfg.GinMode)
}
```

### Run Test

```bash
$ go run test_config.go
```

**Expected output:**

```
Config loaded:
  DB Host: localhost
  DB Port: 5432
  DB Name: todolist_db
  Server Port: 8080
  Gin Mode: debug
```

- [ ] Config prints correctly
- [ ] All values show defaults or from .env

### Clean Up

```bash
$ rm test_config.go    # Linux/Mac
$ del test_config.go   # Windows
```

- [ ] Test file removed

---

## 🐛 Common Issues

### Issue: "package config is not in GOROOT"

**Solution:**

- Make sure you're in the project root directory
- Verify `go.mod` exists
- Run `go mod tidy`

### Issue: database.go shows "model" errors

**Solution:**

- This is **expected** - models don't exist yet
- Ignore for now, will fix in next section

### Issue: .env appears in git status

**Solution:**

- Check `.gitignore` contains `.env`
- Run `git rm --cached .env`
- Add `.env` to `.gitignore`

### Issue: Environment variables not loading

**Solution:**

- Check `.env` file exists
- Verify environment variable names match exactly
- Restart terminal/IDE to reload environment

---

## 📊 Current Progress

```
[████████████░░░░] 18% Complete
```

**Completed:** Setup, dependencies, configuration  
**Next:** Database models

---

**Previous:** [02-dependencies.md](02-dependencies.md)  
**Next:** [04-models.md](04-models.md)
