# Quick Reference - Command Cheatsheet

Quick reference for common commands used throughout this project.

---

## 📦 Go Module Commands

### Initialize Module

```bash
$ go mod init github.com/adityapryg/golang-demo/20-mini-project
```

### Install Package

```bash
$ go get package@version
$ go get github.com/gin-gonic/gin@v1.9.1
```

### Download Dependencies

```bash
$ go mod download
```

### Clean Up Dependencies

```bash
$ go mod tidy
```

### List All Dependencies

```bash
$ go list -m all
```

### Verify Dependencies

```bash
$ go mod verify
```

---

## 🏗️ Build Commands

### Run Application

```bash
$ go run cmd/api/main.go
```

### Build Binary

```bash
# Windows
$ go build -o bin/api.exe cmd/api/main.go

# Linux/Mac
$ go build -o bin/api cmd/api/main.go
```

### Build with Flags

```bash
$ go build -ldflags="-s -w" -o bin/api.exe cmd/api/main.go
# -s -w = Strip debug info (smaller binary)
```

### Cross-Platform Build

```bash
# For Linux from Windows
$ set GOOS=linux
$ set GOARCH=amd64
$ go build -o bin/api cmd/api/main.go

# For Windows from Linux/Mac
$ export GOOS=windows
$ export GOARCH=amd64
$ go build -o bin/api.exe cmd/api/main.go
```

---

## 🧪 Testing Commands

### Run All Tests

```bash
$ go test ./tests/...
```

### Run Tests with Verbose Output

```bash
$ go test -v ./tests/...
```

### Run Specific Test

```bash
$ go test ./tests/ -run TestAuthFlow
$ go test ./tests/ -run TestCreateTodo
```

### Run Tests with Coverage

```bash
$ go test -cover ./tests/...
```

### Generate Coverage Report

```bash
$ go test -coverprofile=coverage.out ./tests/...
$ go tool cover -html=coverage.out
```

### Run Tests in Specific Package

```bash
$ go test github.com/adityapryg/golang-demo/20-mini-project/internal/service
```

---

## 🐳 Docker Commands

### Build Image

```bash
$ docker build -t todolist-api .
```

### Run Container

```bash
$ docker run -p 8080:8080 todolist-api
```

### Docker Compose - Start

```bash
$ docker-compose up
```

### Docker Compose - Start in Background

```bash
$ docker-compose up -d
```

### Docker Compose - Rebuild and Start

```bash
$ docker-compose up --build
```

### Docker Compose - Stop

```bash
$ docker-compose down
```

### Docker Compose - Stop and Remove Volumes

```bash
$ docker-compose down -v
```

### Docker Compose - View Logs

```bash
$ docker-compose logs
$ docker-compose logs -f  # Follow logs
$ docker-compose logs api # Specific service
```

### Docker Compose - Start Only Postgres

```bash
$ docker-compose up -d postgres
```

---

## 🗄️ PostgreSQL Commands

### Connect to Database

```bash
$ psql -U postgres -d todolist_db
```

### Run SQL File

```bash
$ psql -U postgres -d todolist_db -f migrations/001_create_tables.sql
```

### Common psql Commands

```sql
\l              -- List databases
\c todolist_db  -- Connect to database
\dt             -- List tables
\d users        -- Describe users table
\d todos        -- Describe todos table
\q              -- Quit
```

### Export Database

```bash
$ pg_dump -U postgres todolist_db > backup.sql
```

### Import Database

```bash
$ psql -U postgres -d todolist_db < backup.sql
```

---

## 🔧 Environment Variables

### Set Environment Variable

**Windows (Command Prompt):**

```bash
$ set DB_HOST=localhost
$ set JWT_SECRET=my-secret
```

**Windows (PowerShell):**

```powershell
$ $env:DB_HOST="localhost"
$ $env:JWT_SECRET="my-secret"
```

**Linux/Mac:**

```bash
$ export DB_HOST=localhost
$ export JWT_SECRET=my-secret
```

### View Environment Variable

```bash
# Windows (CMD)
$ echo %DB_HOST%

# Windows (PowerShell)
$ echo $env:DB_HOST

# Linux/Mac
$ echo $DB_HOST
```

### Load from .env File

```bash
# Using export (Linux/Mac)
$ export $(cat .env | xargs)

# Using set (Windows - in PowerShell)
$ Get-Content .env | ForEach-Object { $var = $_.Split('='); [Environment]::SetEnvironmentVariable($var[0], $var[1]) }
```

---

## 📁 File Operations

### Create Directory

```bash
# Single directory
$ mkdir foldername

# Nested directories (Windows)
$ mkdir parent\child\grandchild

# Nested directories (Linux/Mac)
$ mkdir -p parent/child/grandchild
```

### Create File

```bash
# Windows
$ type nul > filename.txt
$ echo. > filename.txt

# Linux/Mac
$ touch filename.txt
```

### View File Content

```bash
# Windows
$ type filename.txt

# Linux/Mac
$ cat filename.txt
```

### Copy File

```bash
# Windows
$ copy source.txt destination.txt

# Linux/Mac
$ cp source.txt destination.txt
```

### Delete File

```bash
# Windows
$ del filename.txt

# Linux/Mac
$ rm filename.txt
```

---

## 📊 Git Commands

### Initialize Repository

```bash
$ git init
```

### Check Status

```bash
$ git status
```

### Stage Files

```bash
$ git add .
$ git add filename.txt
$ git add internal/
```

### Commit Changes

```bash
$ git commit -m "Commit message"
```

### View Commit History

```bash
$ git log
$ git log --oneline
$ git log --graph --oneline
```

### View Differences

```bash
$ git diff
$ git diff filename.txt
```

### Create Branch

```bash
$ git branch branch-name
$ git checkout -b branch-name  # Create and switch
```

### Switch Branch

```bash
$ git checkout branch-name
```

### Merge Branch

```bash
$ git merge branch-name
```

### Push to Remote

```bash
$ git push origin main
```

### Pull from Remote

```bash
$ git pull origin main
```

---

## 🌐 API Testing with cURL

### Health Check

```bash
$ curl http://localhost:8080/health
```

### Register User

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","full_name":"Test User"}'
```

### Login

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Get Profile (with token)

```bash
$ curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create Todo

```bash
$ curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"My Todo","description":"Description","status":"pending","priority":"high"}'
```

### Get All Todos

```bash
$ curl http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get Todo by ID

```bash
$ curl http://localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update Todo

```bash
$ curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"status":"completed"}'
```

### Delete Todo

```bash
$ curl -X DELETE http://localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🔍 Debugging Commands

### Check Go Environment

```bash
$ go env
```

### Format Code

```bash
$ go fmt ./...
```

### Vet Code (Find Issues)

```bash
$ go vet ./...
```

### Static Analysis

```bash
$ go install honnef.co/go/tools/cmd/staticcheck@latest
$ staticcheck ./...
```

### Check for Vulnerabilities

```bash
$ go install golang.org/x/vuln/cmd/govulncheck@latest
$ govulncheck ./...
```

### View Documentation

```bash
$ go doc package
$ go doc package.Function
```

---

## 📝 Quick Notes

### Project Structure Check

```bash
# Count files
$ find . -type f | wc -l  # Linux/Mac
$ (Get-ChildItem -Recurse -File).Count  # Windows PowerShell

# Show tree structure
$ tree -L 3  # Linux/Mac
$ tree /F /A  # Windows
```

### Port Already in Use

```bash
# Find process using port 8080 (Windows)
$ netstat -ano | findstr :8080
$ taskkill /PID <PID> /F

# Find process using port 8080 (Linux/Mac)
$ lsof -i :8080
$ kill -9 <PID>
```

### Clear Go Cache

```bash
$ go clean -cache
$ go clean -modcache
$ go clean -testcache
```

---

## 🎯 Common Command Sequences

### Fresh Start

```bash
$ rm -rf bin/  # or: del /s /q bin\
$ go mod tidy
$ go build -o bin/api.exe cmd/api/main.go
$ ./bin/api.exe
```

### Test and Run

```bash
$ go test ./tests/... -v
$ go run cmd/api/main.go
```

### Docker Fresh Start

```bash
$ docker-compose down -v
$ docker-compose build --no-cache
$ docker-compose up
```

### Full Rebuild

```bash
$ go clean -cache
$ go mod tidy
$ go test ./tests/...
$ go build -o bin/api.exe cmd/api/main.go
$ ./bin/api.exe
```

---

**Back to:** [README.md](README.md) | [Documentation Index](00-overview.md)
