# Go Learning Repository - AI Agent Instructions

## Repository Purpose

This is a **progressive learning repository** for teaching Golang fundamentals through standalone, runnable examples. Each numbered folder (`01-hello-world` through `19-swaggo`) is a self-contained lesson, not a unified application.

## Project Structure Pattern

### Lesson Organization

- Each `XX-topic-name/` folder contains:
  - `main.go` - Runnable demonstration code
  - `README.md` - Indonesian language documentation with learning objectives, code explanations, and expected output
  - `go.mod` - Module definition with path pattern `github.com/adityapryg/golang-demo/XX-topic-name`
  - Optional subdirectories (e.g., `utils/`) for package demonstrations

### Module Path Convention

All go.mod files use the pattern: `module github.com/adityapryg/golang-demo/XX-folder-name`

Example from `02-workspace-structure/go.mod`:

```go
module github.com/adityapryg/golang-demo/02-workspace-structure
go 1.22
```

## Running Examples

### Standard Workflow

```bash
cd XX-topic-name    # Navigate to specific lesson
go run main.go      # Run the example directly
```

### With External Dependencies

For lessons using frameworks (e.g., `12-gin-framework`):

```bash
cd 12-gin-framework
go mod download     # Download dependencies if needed
go run main.go      # Dependencies auto-fetched on first run
```

### CLI Applications with Flags

Some examples (`09-defer-panic-recover`, `12-gin-framework`) include flag parsing:

```bash
go run main.go -port 8080 -host localhost -debug
```

## Code Patterns & Conventions

### Error Handling Strategy

This codebase demonstrates **three error patterns** explicitly:

1. **Basic errors** - `errors.New()` for simple cases
2. **Sentinel errors** - Package-level variables for specific cases:
   ```go
   var (
       ErrNotFound = errors.New("data tidak ditemukan")
       ErrAlreadyExists = errors.New("data sudah ada")
   )
   ```
3. **Custom error types** - Structs implementing `Error()` for contextual information:
   ```go
   type ValidationError struct {
       Field   string
       Message string
   }
   ```

### Exported vs Unexported (Critical Go Concept)

The codebase **explicitly demonstrates** Go's capitalization-based visibility:

- Capital first letter = Exported (public): `SayHello()`, `User{}`
- Lowercase first letter = Unexported (private): `greet()`, `validateAge()`

See `02-workspace-structure/utils/greetings.go` for teaching examples.

### Struct Tag Validation (Gin Framework)

`12-gin-framework` uses Gin's binding validation:

```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,min=1,max=150"`
}
```

### Concurrency Safety Pattern

Examples with shared state (`11-rest-api-intro`, `12-gin-framework`) use `sync.RWMutex`:

```go
var (
    users   = make(map[int]User)
    usersMu sync.RWMutex  // Always lock before accessing users
)
```

## Documentation Language

All README files and comments are in **Indonesian (Bahasa Indonesia)**. Code should maintain this pattern for consistency with the learning material.

## When Adding New Examples

1. **Create numbered folder** following existing pattern (`13-new-topic/`)
2. **Initialize module**: `go mod init github.com/adityapryg/golang-demo/13-new-topic`
3. **Write self-contained code** - No cross-folder dependencies except for package demo purposes
4. **Include README.md** with:
   - Tujuan Pembelajaran (Learning Objectives)
   - Penjelasan (Explanation)
   - Cara Menjalankan (How to Run)
   - Output yang Diharapkan (Expected Output)
5. **Add visual separators** in terminal output:
   ```go
   fmt.Println("===========================================")
   fmt.Println("   JUDUL DEMONSTRASI")
   fmt.Println("===========================================")
   ```

## Testing & Debugging

- **No formal test suites by design** - each `main.go` demonstrates concepts through runnable output for pedagogical purposes
- Verify changes by running `go run main.go` in affected folder
- Check for compilation errors across all modules:
  ```bash
  for d in */; do (cd "$d" && [ -f go.mod ] && go build); done
  ```

## Lessons Overview

The repository contains 19 lessons covering:

**Fundamentals (01-10):**

- Basic syntax, control flow, data structures
- Functions, structs, error handling
- Command-line applications and defer/panic/recover

**REST API Development (11-19):**

- `11-rest-api-intro` - REST API basics with net/http
- `12-gin-framework` - Gin framework and routing
- `13-middleware` - Middleware (logging, auth, CORS, rate limiting)
- `14-crud-no-db` - CRUD API with in-memory storage
- `15-database-gorm` - PostgreSQL connection and GORM basics
- `16-crud-with-db` - CRUD API with database integration
- `17-jwt-auth` - JWT authentication and authorization
- `18-testing-docs` - Testing best practices and coverage
- `19-swaggo` - API documentation with Swagger/OpenAPI

**Planned (20+):**

- `20-mini-project` - Full-stack Todo REST API with Vue.js integration

When maintaining these lessons:

- Each lesson is self-contained with go.mod, main.go, and README.md
- All documentation is in Indonesian (Bahasa Indonesia)
- READMEs include: Tujuan Pembelajaran, Penjelasan, Cara Menjalankan, and expected output

## Common Dependencies

- **Standard library only** for lessons 01-10
- **Gin framework** (`github.com/gin-gonic/gin v1.9.1`) for lessons 11-14, 16-17, 19
- **PostgreSQL + GORM** (`gorm.io/driver/postgres`, `gorm.io/gorm`) for lessons 15-16
- **JWT** (`github.com/golang-jwt/jwt/v5`) for lesson 17
- **Testing** (`github.com/stretchr/testify`) for lesson 18
- **Swagger** (`github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`) for lesson 19
- Go version: **1.22** consistently across all modules
