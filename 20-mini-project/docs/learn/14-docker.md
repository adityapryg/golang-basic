# 14 - Docker Configuration

Configure Docker for containerized deployment with PostgreSQL.

---

## Overview

We'll create:

1. `Dockerfile` - Container image for Go application
2. `docker-compose.yml` - Multi-container setup with PostgreSQL
3. `.dockerignore` - Exclude unnecessary files from image

Docker enables consistent deployments across environments.

---

## Step 1: Create Dockerfile

### 1.1 Create Dockerfile

📝 **Create file:** `Dockerfile`

```dockerfile
# Multi-stage build for smaller final image
# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/api/main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy migrations (optional, if needed at runtime)
COPY --from=builder /app/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
```

- [ ] File created at `Dockerfile`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Multi-Stage Builds

```dockerfile
# Stage 1: Builder
FROM golang:1.22-alpine AS builder
# ... build process ...

# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /app/main .
```

**Why multi-stage builds?**

| Aspect           | Single Stage                | Multi-Stage               |
| ---------------- | --------------------------- | ------------------------- |
| Final image size | ~800 MB (with Go toolchain) | ~20 MB (binary only)      |
| Security         | Includes build tools        | Only runtime dependencies |
| Attack surface   | Large                       | Minimal                   |
| Build time       | Slower                      | Faster (cached layers)    |

**Image size comparison:**

```
Single-stage:   golang:1.22        = 810 MB
Multi-stage:    alpine:latest      = 7 MB
                + binary           = ~15 MB
                + ca-certificates  = 1 MB
                -------------------------
                Total              = ~23 MB
```

### 1.3 Understanding Dockerfile Instructions

```dockerfile
# Base image
FROM golang:1.22-alpine

# Install packages
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Build flags explained:
CGO_ENABLED=0           # Disable CGO for static binary
GOOS=linux              # Target OS
-ldflags="-w -s"        # Strip debug info (-w) and symbol table (-s)
```

**Layer caching optimization:**

```dockerfile
# ❌ Bad - any code change invalidates dependency cache
COPY . .
RUN go mod download

# ✅ Good - dependencies cached separately
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

### 1.4 Understanding Security Best Practices

```dockerfile
# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Switch to non-root user
USER appuser
```

**Why run as non-root?**

- **Principle of least privilege**: App doesn't need root access
- **Container escape protection**: Limited damage if exploited
- **Best practice**: Kubernetes/production requires non-root

**Security checklist:**

- ✅ Run as non-root user
- ✅ Use minimal base image (alpine)
- ✅ Don't include build tools in final image
- ✅ Keep secrets out of image
- ✅ Scan for vulnerabilities

### 1.5 Understanding Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

**Health check parameters:**

| Parameter        | Value | Meaning                    |
| ---------------- | ----- | -------------------------- |
| `--interval`     | 30s   | Check every 30 seconds     |
| `--timeout`      | 3s    | Fail if check takes >3s    |
| `--start-period` | 5s    | Grace period for startup   |
| `--retries`      | 3     | Unhealthy after 3 failures |

**Health check states:**

```
starting → healthy → unhealthy → healthy
   ↓          ↓          ↓          ↓
  0s        35s        95s       125s
```

---

## Step 2: Create Docker Compose File

### 2.1 Create docker-compose.yml

📝 **Create file:** `docker-compose.yml`

```yaml
version: "3.8"

services:
  # PostgreSQL database service
  postgres:
    image: postgres:15-alpine
    container_name: todolist-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: todolist_db
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --locale=en_US.UTF-8"
    ports:
      - "5432:5432"
    volumes:
      # Persist data between restarts
      - postgres_data:/var/lib/postgresql/data
      # Optional: Auto-run migrations on startup
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - todolist-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

  # Go API service
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: todolist-api
    restart: unless-stopped
    environment:
      # Database configuration
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: todolist_db

      # Application configuration
      JWT_SECRET: your-super-secret-key-change-this-in-production
      SERVER_PORT: 8080
      GIN_MODE: release
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - todolist-network
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--no-verbose",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 40s

# Named volumes for data persistence
volumes:
  postgres_data:
    driver: local

# Custom network for service communication
networks:
  todolist-network:
    driver: bridge
```

- [ ] File created at `docker-compose.yml`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding Docker Compose Services

```yaml
services:
  postgres:
    image: postgres:15-alpine
    # ...

  api:
    build:
      context: .
      dockerfile: Dockerfile
    # ...
```

**Service types:**

| Service    | Type        | Purpose                        |
| ---------- | ----------- | ------------------------------ |
| `postgres` | Image-based | Use pre-built PostgreSQL image |
| `api`      | Build-based | Build from Dockerfile          |

### 2.3 Understanding depends_on with Health Checks

```yaml
api:
  depends_on:
    postgres:
      condition: service_healthy
```

**Startup sequence:**

```
1. Start postgres container
2. Wait for postgres health check to pass
3. Only then start api container
```

**Without health check:**

```yaml
# ❌ Bad - API might start before DB is ready
depends_on:
  - postgres
```

**With health check:**

```yaml
# ✅ Good - API waits for DB to be ready
depends_on:
  postgres:
    condition: service_healthy
```

### 2.4 Understanding Volumes

```yaml
volumes:
  # Persist PostgreSQL data
  - postgres_data:/var/lib/postgresql/data

  # Mount migrations folder
  - ./migrations:/docker-entrypoint-initdb.d
```

**Volume types:**

| Type             | Syntax                | Use Case          | Data Persistence        |
| ---------------- | --------------------- | ----------------- | ----------------------- |
| Named volume     | `postgres_data:/path` | Database storage  | Yes (managed by Docker) |
| Bind mount       | `./migrations:/path`  | Development files | Yes (host filesystem)   |
| Anonymous volume | `/path`               | Temporary data    | No                      |

**PostgreSQL auto-initialization:**

```yaml
- ./migrations:/docker-entrypoint-initdb.d
```

- PostgreSQL runs any `.sql` files in this folder on first startup
- Automatic schema creation!

### 2.5 Understanding Networks

```yaml
networks:
  todolist-network:
    driver: bridge

services:
  postgres:
    networks:
      - todolist-network
  api:
    networks:
      - todolist-network
```

**Why custom networks?**

- **Service discovery**: Services can reach each other by name
- **Isolation**: Separate from other Docker containers
- **DNS**: Automatic DNS resolution

**Communication example:**

```go
// In api service, connect to postgres by service name
DB_HOST=postgres  // Not localhost!
```

---

## Step 3: Create .dockerignore

### 3.1 Create .dockerignore File

📝 **Create file:** `.dockerignore`

```
# Git
.git
.gitignore
.github

# Documentation
README.md
*.md
docs/

# IDE
.vscode
.idea
*.swp
*.swo

# Build artifacts
bin/
*.exe
*.exe~

# Dependencies (will be downloaded in container)
vendor/

# Environment files
.env
.env.local
.env.*.local

# Test files
*_test.go
tests/
coverage.out
*.test

# Logs
*.log

# OS files
.DS_Store
Thumbs.db

# Temporary files
tmp/
temp/
```

- [ ] File created at `.dockerignore`
- [ ] Content copied exactly
- [ ] File saved

### 3.2 Understanding .dockerignore

**Why use .dockerignore?**

- **Smaller build context**: Faster builds
- **Security**: Don't copy secrets (.env)
- **Cache efficiency**: Exclude changing files

**Impact on build time:**

```
Without .dockerignore:
  Sending build context to Docker daemon  1.2GB
  Build time: 45 seconds

With .dockerignore:
  Sending build context to Docker daemon  15MB
  Build time: 8 seconds
```

---

## Step 4: Build and Run with Docker

### 4.1 Build Docker Image

```bash
# Build image
$ docker build -t todolist-api:latest .
```

**Expected output:**

```
[+] Building 45.2s (18/18) FINISHED
 => [internal] load build definition from Dockerfile
 => => transferring dockerfile: 1.2kB
 => [internal] load .dockerignore
 => [builder 1/7] FROM docker.io/library/golang:1.22-alpine
 => [builder 2/7] WORKDIR /app
 => [builder 3/7] COPY go.mod go.sum ./
 => [builder 4/7] RUN go mod download
 => [builder 5/7] COPY . .
 => [builder 6/7] RUN CGO_ENABLED=0 GOOS=linux go build ...
 => [stage-1 1/5] FROM docker.io/library/alpine:latest
 => [stage-1 2/5] RUN apk --no-cache add ca-certificates
 => [stage-1 3/5] COPY --from=builder /app/main .
 => exporting to image
 => => writing image sha256:abc123...
 => => naming to docker.io/library/todolist-api:latest
```

- [ ] Image built successfully
- [ ] No errors

### 4.2 Run with Docker Compose

```bash
# Start all services (detached mode)
$ docker-compose up -d

# Check status
$ docker-compose ps

# Expected output:
NAME                    STATUS              PORTS
todolist-api            Up (healthy)        0.0.0.0:8080->8080/tcp
todolist-postgres       Up (healthy)        0.0.0.0:5432->5432/tcp
```

- [ ] Services started
- [ ] Both services healthy

### 4.3 View Logs

```bash
# View all logs
$ docker-compose logs

# Follow logs (live updates)
$ docker-compose logs -f

# Logs for specific service
$ docker-compose logs api
$ docker-compose logs postgres
```

**Expected API logs:**

```
todolist-api | Configuration loaded
todolist-api | Database connected
todolist-api | Repositories initialized
todolist-api | Services initialized
todolist-api | Handlers initialized
todolist-api | Server starting on :8080
```

- [ ] No error logs
- [ ] Server started successfully

### 4.4 Test Application

```bash
# Test health check
$ curl http://localhost:8080/health

# Expected response:
{
  "status": "ok",
  "message": "Service is healthy",
  "database": "connected"
}

# Test registration
$ curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

- [ ] Health check returns 200
- [ ] Registration works
- [ ] Database connected

---

## Step 5: Docker Commands Cheatsheet

### 5.1 Docker Compose Commands

```bash
# Start services
docker-compose up              # Foreground
docker-compose up -d           # Background (detached)
docker-compose up --build      # Rebuild before starting

# Stop services
docker-compose stop            # Stop containers
docker-compose down            # Stop and remove containers
docker-compose down -v         # Also remove volumes (data loss!)

# View status
docker-compose ps              # List containers
docker-compose top             # Show running processes

# View logs
docker-compose logs            # All logs
docker-compose logs -f         # Follow logs
docker-compose logs api        # Specific service

# Execute commands
docker-compose exec api sh     # Shell in api container
docker-compose exec postgres psql -U postgres  # PostgreSQL CLI

# Rebuild
docker-compose build           # Rebuild images
docker-compose build --no-cache  # Fresh build
```

### 5.2 Docker Commands

```bash
# List images
docker images

# Remove image
docker rmi todolist-api:latest

# List containers
docker ps                      # Running containers
docker ps -a                   # All containers

# Stop container
docker stop todolist-api

# Remove container
docker rm todolist-api

# View logs
docker logs todolist-api
docker logs -f todolist-api    # Follow

# Execute command in container
docker exec -it todolist-api sh

# Inspect container
docker inspect todolist-api
```

### 5.3 Volume Commands

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect 20-mini-project_postgres_data

# Remove volume
docker volume rm 20-mini-project_postgres_data

# Remove unused volumes
docker volume prune
```

---

## Step 6: Production Considerations

### 6.1 Environment Variables for Production

```yaml
# docker-compose.prod.yml
services:
  api:
    environment:
      # Use secrets or external config
      JWT_SECRET: ${JWT_SECRET}
      DB_PASSWORD: ${DB_PASSWORD}
      GIN_MODE: release

      # Disable debug features
      DEBUG: false
```

**Load from .env file:**

```bash
# Create .env file (don't commit!)
JWT_SECRET=very-strong-production-secret-minimum-32-chars
DB_PASSWORD=strong-database-password

# Docker Compose automatically loads .env
docker-compose up -d
```

### 6.2 Health Check Configuration

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
  interval: 30s # Check every 30s
  timeout: 3s # Fail if check takes >3s
  retries: 3 # Mark unhealthy after 3 failures
  start_period: 40s # Grace period during startup
```

### 6.3 Resource Limits

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: "0.5" # Max 50% of one CPU
          memory: 512M # Max 512MB RAM
        reservations:
          cpus: "0.25" # Guaranteed 25% CPU
          memory: 256M # Guaranteed 256MB RAM
```

### 6.4 Restart Policies

```yaml
restart: unless-stopped # Recommended for production

# Options:
# no           - Never restart
# always       - Always restart
# on-failure   - Only on error
# unless-stopped - Always unless manually stopped
```

---

## Step 7: Troubleshooting

### 7.1 Common Issues

**Issue: "port already in use"**

```bash
# Check what's using port
lsof -i :8080              # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Change port in docker-compose.yml
ports:
  - "8081:8080"  # Host:Container
```

**Issue: "database connection refused"**

```bash
# Check postgres is healthy
docker-compose ps

# View postgres logs
docker-compose logs postgres

# Test connection
docker-compose exec postgres psql -U postgres -d todolist_db
```

**Issue: "API not starting"**

```bash
# View API logs
docker-compose logs api

# Check health
docker-compose exec api wget -O- http://localhost:8080/health

# Access container shell
docker-compose exec api sh
```

### 7.2 Reset Everything

```bash
# Stop and remove everything
docker-compose down -v

# Remove images
docker rmi $(docker images -q todolist*)

# Fresh start
docker-compose up --build
```

---

## Step 8: Commit Changes

### 8.1 Check Status

```bash
$ git status
```

- [ ] Docker files shown

### 8.2 Stage Files

```bash
$ git add Dockerfile docker-compose.yml .dockerignore
```

- [ ] Files staged

### 8.3 Commit

```bash
$ git commit -m "Add Docker configuration with multi-stage build and compose setup"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `Dockerfile` created
- [ ] Multi-stage build configured
- [ ] Builder stage with Go toolchain
- [ ] Runtime stage with Alpine
- [ ] Static binary compilation
- [ ] Non-root user configured
- [ ] Health check configured
- [ ] `docker-compose.yml` created
- [ ] PostgreSQL service configured
- [ ] API service configured
- [ ] Service dependencies configured
- [ ] Health check conditions set
- [ ] Named volumes for persistence
- [ ] Custom network configured
- [ ] Auto-migration on startup
- [ ] Environment variables configured
- [ ] `.dockerignore` created
- [ ] Unnecessary files excluded
- [ ] Understanding of multi-stage builds
- [ ] Understanding of Docker Compose
- [ ] Understanding of volumes
- [ ] Understanding of networks
- [ ] Understanding of health checks
- [ ] Docker image builds successfully
- [ ] Services start with docker-compose
- [ ] Application accessible on port 8080
- [ ] Health check passes
- [ ] Changes committed to git

---

## 🧪 Quick Test

```bash
# 1. Build and start
docker-compose up --build -d

# 2. Wait for healthy status
docker-compose ps

# 3. Test health
curl http://localhost:8080/health

# 4. Test registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"docker","email":"docker@example.com","password":"pass123","full_name":"Docker User"}'

# 5. Check logs
docker-compose logs api

# 6. Clean up
docker-compose down
```

- [ ] All tests pass
- [ ] No errors in logs

---

## 🐛 Common Issues

### Issue: "Cannot connect to Docker daemon"

**Solution:** Start Docker Desktop or Docker service

### Issue: "Build fails: go.mod not found"

**Solution:** Run `go mod tidy` before building

### Issue: "API can't connect to postgres"

**Solution:** Use service name `postgres`, not `localhost`

### Issue: "Volume permission denied"

**Solution:** Check volume ownership and permissions

---

## 📊 Current Progress

```
[██████████████████████████████████████] 87% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities, middleware, repositories, services, handlers, routes, main app, migrations, Docker  
**Next:** Testing implementation

---

**Previous:** [13-database.md](13-database.md)  
**Next:** [15-testing.md](15-testing.md)
