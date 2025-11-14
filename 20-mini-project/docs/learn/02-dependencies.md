# 02 - Installing Dependencies

Install all required Go packages for the project.

---

## Overview

We'll install these packages:

1. **Gin** - HTTP web framework
2. **GORM** - ORM library
3. **PostgreSQL Driver** - Database driver
4. **JWT** - JSON Web Token library
5. **Crypto** - Password hashing (bcrypt)
6. **Testify** - Testing toolkit
7. **Swagger** - API documentation (optional)

---

## Step 1: Install Gin Framework

### 1.1 Install Gin

```bash
$ go get github.com/gin-gonic/gin@v1.9.1
```

**Expected output:**

```
go: downloading github.com/gin-gonic/gin v1.9.1
go: downloading github.com/gin-contrib/sse v0.1.0
...
go: added github.com/gin-gonic/gin v1.9.1
```

- [ ] Installation successful
- [ ] No error messages

### 1.2 Verify Installation

```bash
$ cat go.mod    # Linux/Mac
$ type go.mod   # Windows
```

**Should contain:**

```go
require github.com/gin-gonic/gin v1.9.1
```

- [ ] Gin appears in go.mod

💡 **What is Gin?**

- Fast HTTP web framework
- Handles routing, middleware, request/response
- Similar to Express.js in Node.js

---

## Step 2: Install GORM and PostgreSQL Driver

### 2.1 Install GORM

```bash
$ go get gorm.io/gorm@v1.25.5
```

**Expected output:**

```
go: downloading gorm.io/gorm v1.25.5
go: added gorm.io/gorm v1.25.5
```

- [ ] GORM installed

### 2.2 Install PostgreSQL Driver

```bash
$ go get gorm.io/driver/postgres@v1.5.4
```

**Expected output:**

```
go: downloading gorm.io/driver/postgres v1.5.4
go: downloading github.com/jackc/pgx/v5 v5.4.3
...
go: added gorm.io/driver/postgres v1.5.4
```

- [ ] PostgreSQL driver installed

### 2.3 Verify go.mod Updates

```bash
$ cat go.mod | grep gorm    # Linux/Mac
$ type go.mod | findstr gorm    # Windows
```

**Should show:**

```go
gorm.io/driver/postgres v1.5.4
gorm.io/gorm v1.25.5
```

- [ ] Both GORM packages in go.mod

💡 **What is GORM?**

- Object-Relational Mapping library
- Converts between Go structs and database tables
- Handles CRUD operations, relationships, migrations

---

## Step 3: Install JWT Library

### 3.1 Install JWT Package

```bash
$ go get github.com/golang-jwt/jwt/v5@v5.2.0
```

**Expected output:**

```
go: downloading github.com/golang-jwt/jwt/v5 v5.2.0
go: added github.com/golang-jwt/jwt/v5 v5.2.0
```

- [ ] JWT library installed

### 3.2 Verify Installation

```bash
$ cat go.mod | grep jwt    # Linux/Mac
$ type go.mod | findstr jwt    # Windows
```

**Should contain:**

```go
github.com/golang-jwt/jwt/v5 v5.2.0
```

- [ ] JWT package in go.mod

💡 **What is JWT?**

- JSON Web Token standard
- Used for authentication/authorization
- Creates and validates secure tokens

---

## Step 4: Install Crypto Package

### 4.1 Install golang.org/x/crypto

```bash
$ go get golang.org/x/crypto@v0.17.0
```

**Expected output:**

```
go: downloading golang.org/x/crypto v0.17.0
go: added golang.org/x/crypto v0.17.0
```

- [ ] Crypto package installed

### 4.2 Verify Installation

```bash
$ cat go.mod | grep crypto    # Linux/Mac
$ type go.mod | findstr crypto    # Windows
```

**Should contain:**

```go
golang.org/x/crypto v0.17.0
```

- [ ] Crypto package in go.mod

💡 **What is this for?**

- Provides bcrypt for password hashing
- Much more secure than plain text or simple hashing
- Industry standard for password storage

---

## Step 5: Install Testing Libraries

### 5.1 Install Testify

```bash
$ go get github.com/stretchr/testify@v1.8.4
```

**Expected output:**

```
go: downloading github.com/stretchr/testify v1.8.4
go: downloading github.com/davecgh/go-spew v1.1.1
...
go: added github.com/stretchr/testify v1.8.4
```

- [ ] Testify installed

### 5.2 Verify Installation

```bash
$ cat go.mod | grep testify    # Linux/Mac
$ type go.mod | findstr testify    # Windows
```

**Should contain:**

```go
github.com/stretchr/testify v1.8.4
```

- [ ] Testify package in go.mod

💡 **What is Testify?**

- Testing toolkit with assertions
- Test suite support
- Mock generation
- Makes tests easier to write and read

---

## Step 6: Install Swagger (Optional)

### 6.1 Install Swagger Packages

⚠️ **Note:** These are optional for API documentation. Skip if you want a minimal setup.

```bash
$ go get github.com/swaggo/gin-swagger@v1.6.0
$ go get github.com/swaggo/files@v1.0.1
$ go get github.com/swaggo/swag@v1.16.2
```

**Expected output for each:**

```
go: downloading github.com/swaggo/...
go: added github.com/swaggo/...
```

- [ ] Swagger packages installed (optional)

### 6.2 Verify Swagger Installation

```bash
$ cat go.mod | grep swaggo    # Linux/Mac
$ type go.mod | findstr swaggo    # Windows
```

**Should contain (if installed):**

```go
github.com/swaggo/files v1.0.1
github.com/swaggo/gin-swagger v1.6.0
github.com/swaggo/swag v1.16.2
```

- [ ] Swagger packages in go.mod (if installed)

💡 **What is Swagger?**

- API documentation generator
- Interactive API testing UI
- Auto-generates docs from code comments

---

## Step 7: Download All Dependencies

### 7.1 Download Dependencies

```bash
$ go mod download
```

**Expected output:**

```
go: downloading ...
(lists all packages being downloaded)
```

- [ ] All dependencies downloaded
- [ ] No errors

### 7.2 Tidy Dependencies

```bash
$ go mod tidy
```

**Expected output:**

```
(may show some additions or removals of indirect dependencies)
```

- [ ] Command completed successfully

💡 **What does `go mod tidy` do?**

- Adds missing dependencies
- Removes unused dependencies
- Updates go.sum file
- Organizes go.mod

---

## Step 8: Verify go.mod Content

### 8.1 View Complete go.mod

```bash
$ cat go.mod    # Linux/Mac
$ type go.mod   # Windows
```

**Expected content:**

```go
module github.com/adityapryg/golang-demo/20-mini-project

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/stretchr/testify v1.8.4
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
	golang.org/x/crypto v0.17.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	... (more indirect dependencies)
)
```

- [ ] Module path correct
- [ ] All required packages present
- [ ] Go version is 1.22

### 8.2 Check go.sum Exists

```bash
$ ls go.sum    # Linux/Mac
$ dir go.sum   # Windows
```

- [ ] `go.sum` file exists

💡 **What is go.sum?**

- Checksums for all dependencies
- Ensures package integrity
- Used for verification
- Should be committed to git

---

## Step 9: Verify Dependencies

### 9.1 List All Dependencies

```bash
$ go list -m all
```

**Expected output:**

```
github.com/adityapryg/golang-demo/20-mini-project
github.com/gin-gonic/gin v1.9.1
github.com/golang-jwt/jwt/v5 v5.2.0
...
(lists all direct and indirect dependencies)
```

- [ ] Command runs without errors
- [ ] All main packages listed

### 9.2 Check for Vulnerabilities (Optional)

```bash
$ go list -m -u all
```

This shows if any dependencies have updates available.

- [ ] Command runs (optional check)

---

## Step 10: Commit Changes

### 10.1 Check Git Status

```bash
$ git status
```

**Expected output:**

```
Changes not staged for commit:
  modified:   go.mod

Untracked files:
  go.sum
```

- [ ] Changes detected

### 10.2 Stage Changes

```bash
$ git add go.mod go.sum
```

- [ ] Files staged

### 10.3 Commit

```bash
$ git commit -m "Add project dependencies"
```

**Expected output:**

```
[main xxxxxxx] Add project dependencies
 2 files changed, XXX insertions(+)
 create mode 100644 go.sum
```

- [ ] Commit successful

---

## ✅ Completion Checklist

Verify all packages are installed:

- [ ] **Gin Framework** v1.9.1 - Web framework
- [ ] **GORM** v1.25.5 - ORM library
- [ ] **PostgreSQL Driver** v1.5.4 - Database driver
- [ ] **JWT** v5.2.0 - Authentication tokens
- [ ] **Crypto** v0.17.0 - Password hashing
- [ ] **Testify** v1.8.4 - Testing toolkit
- [ ] **Swagger** (optional) - API documentation
- [ ] `go.mod` contains all packages
- [ ] `go.sum` file created
- [ ] `go mod download` completed successfully
- [ ] `go mod tidy` ran without errors
- [ ] Changes committed to git

---

## 🧪 Testing Dependencies

Create a quick test file to verify imports work:

### Test File

📝 **Create temporary file:** `test_imports.go` (in project root)

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Just testing imports compile
	_ = gin.New()
	_ = jwt.New(jwt.SigningMethodHS256)
	_ = assert.Equal
	_ = bcrypt.GenerateFromPassword
	_ = postgres.Open
	_ = &gorm.Config{}
}
```

### Run Test

```bash
$ go run test_imports.go
```

**Expected:** No errors (program may not output anything, which is fine)

- [ ] File compiles without errors

### Clean Up

```bash
$ rm test_imports.go    # Linux/Mac
$ del test_imports.go   # Windows
```

- [ ] Test file removed

---

## 🐛 Common Issues

### Issue: "go get: timeout"

**Solution:**

- Check internet connection
- Try with Go proxy: `export GOPROXY=https://proxy.golang.org,direct`
- Or: `go env -w GOPROXY=https://proxy.golang.org,direct`

### Issue: "invalid version"

**Solution:**

- Check exact version numbers
- Use `@latest` to get newest: `go get package@latest`

### Issue: "ambiguous import"

**Solution:**

- Run `go mod tidy` to clean up
- Delete go.sum and run `go mod download` again

### Issue: "package not found"

**Solution:**

- Ensure you're in project directory
- Verify go.mod exists
- Run `go mod download` again

---

## 📊 Current Progress

```
[████████░░░░░░░░] 12% Complete
```

**Completed:** Project setup, dependencies installed  
**Next:** Configuration files

---

**Previous:** [01-initial-setup.md](01-initial-setup.md)  
**Next:** [03-configuration.md](03-configuration.md)
