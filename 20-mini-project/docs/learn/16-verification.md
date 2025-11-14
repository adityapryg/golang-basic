# 16 - Final Verification & Deployment

Comprehensive testing and deployment verification.

---

## Overview

We'll verify:

1. Local development setup
2. API endpoints with cURL
3. Docker deployment
4. Production readiness

---

## Step 1: Local Development Verification

### 1.1 Start PostgreSQL

```bash
$ docker-compose up -d postgres
```

**Expected output:**

```
[+] Running 2/2
 ✔ Network 20-mini-project_default      Created
 ✔ Container 20-mini-project-postgres-1 Started
```

- [ ] PostgreSQL container running

### 1.2 Check Database Connection

```bash
$ docker exec -it 20-mini-project-postgres-1 psql -U postgres -d todolist_db
```

**Inside psql:**

```sql
\dt  -- List tables
SELECT * FROM users;
SELECT * FROM todos;
\q   -- Quit
```

- [ ] Tables exist
- [ ] Can query database

### 1.3 Start API Server

```bash
$ go run cmd/api/main.go
```

**Expected output:**

```
2024/01/15 10:30:00 [Config] Loaded configuration
2024/01/15 10:30:00 [Database] Connecting to PostgreSQL...
2024/01/15 10:30:01 [Database] Connected successfully
2024/01/15 10:30:01 [Database] Running auto-migration...
2024/01/15 10:30:02 [Routes] Registered routes
2024/01/15 10:30:02 [Server] Starting on :8080
```

- [ ] Server starts without errors
- [ ] Listening on port 8080

### 1.4 Test Health Endpoint

```bash
$ curl http://localhost:8080/health
```

**Expected response:**

```json
{
  "success": true,
  "message": "Server is running",
  "data": {
    "database": "connected",
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:15+07:00"
  }
}
```

- [ ] Health check returns 200 OK
- [ ] Database status is "connected"

---

## Step 2: API Testing with cURL

### 2.1 Test User Registration

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass123",
    "full_name": "John Doe"
  }'
```

**Expected response (201 Created):**

```json
{
  "success": true,
  "message": "User berhasil didaftarkan",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-15T10:35:00+07:00"
  }
}
```

- [ ] Status code 201
- [ ] User created successfully
- [ ] No password in response

### 2.2 Test Duplicate Registration

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "another@example.com",
    "password": "SecurePass123",
    "full_name": "John Duplicate"
  }'
```

**Expected response (409 Conflict):**

```json
{
  "success": false,
  "message": "Username sudah digunakan",
  "error": "user already exists"
}
```

- [ ] Status code 409
- [ ] Duplicate username rejected

### 2.3 Test Login

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123"
  }'
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe"
    }
  }
}
```

**💾 Save the token for next steps:**

```bash
$ export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

- [ ] Status code 200
- [ ] JWT token received
- [ ] Token saved to environment variable

### 2.4 Test Invalid Login

```bash
$ curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "WrongPassword"
  }'
```

**Expected response (401 Unauthorized):**

```json
{
  "success": false,
  "message": "Username atau password salah",
  "error": "invalid credentials"
}
```

- [ ] Status code 401
- [ ] Invalid credentials rejected

### 2.5 Test Get Profile (Protected)

```bash
$ curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Profil pengguna",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-15T10:35:00+07:00"
  }
}
```

- [ ] Status code 200
- [ ] Profile retrieved with token

### 2.6 Test Unauthorized Access

```bash
$ curl -X GET http://localhost:8080/api/v1/users/profile
```

**Expected response (401 Unauthorized):**

```json
{
  "success": false,
  "message": "Token tidak valid atau sudah kedaluwarsa",
  "error": "missing authorization header"
}
```

- [ ] Status code 401
- [ ] Unauthorized without token

### 2.7 Test Create Todo

```bash
$ curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive README and setup guides",
    "status": "pending",
    "priority": "high",
    "due_date": "2024-01-30"
  }'
```

**Expected response (201 Created):**

```json
{
  "success": true,
  "message": "Todo berhasil dibuat",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "description": "Write comprehensive README and setup guides",
    "status": "pending",
    "priority": "high",
    "due_date": "2024-01-30T00:00:00Z",
    "user_id": 1,
    "created_at": "2024-01-15T10:40:00+07:00"
  }
}
```

**💾 Save todo ID:**

```bash
$ export TODO_ID=1
```

- [ ] Status code 201
- [ ] Todo created successfully
- [ ] Todo ID saved

### 2.8 Test Get All Todos

```bash
$ curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN"
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Daftar todos",
  "data": [
    {
      "id": 1,
      "title": "Complete project documentation",
      "description": "Write comprehensive README and setup guides",
      "status": "pending",
      "priority": "high",
      "due_date": "2024-01-30T00:00:00Z",
      "created_at": "2024-01-15T10:40:00+07:00"
    }
  ]
}
```

- [ ] Status code 200
- [ ] Todo list retrieved

### 2.9 Test Get Todo by ID

```bash
$ curl -X GET http://localhost:8080/api/v1/todos/$TODO_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Detail todo",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "pending",
    "priority": "high"
  }
}
```

- [ ] Status code 200
- [ ] Specific todo retrieved

### 2.10 Test Update Todo

```bash
$ curl -X PUT http://localhost:8080/api/v1/todos/$TODO_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "in_progress",
    "priority": "medium"
  }'
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Todo berhasil diupdate",
  "data": {
    "id": 1,
    "title": "Complete project documentation",
    "status": "in_progress",
    "priority": "medium",
    "updated_at": "2024-01-15T11:00:00+07:00"
  }
}
```

- [ ] Status code 200
- [ ] Todo updated successfully

### 2.11 Test Filter Todos

```bash
# By status
$ curl -X GET "http://localhost:8080/api/v1/todos?status=in_progress" \
  -H "Authorization: Bearer $TOKEN"

# By priority
$ curl -X GET "http://localhost:8080/api/v1/todos?priority=high" \
  -H "Authorization: Bearer $TOKEN"

# By both
$ curl -X GET "http://localhost:8080/api/v1/todos?status=pending&priority=high" \
  -H "Authorization: Bearer $TOKEN"
```

- [ ] Filter by status works
- [ ] Filter by priority works
- [ ] Combined filters work

### 2.12 Test Delete Todo

```bash
$ curl -X DELETE http://localhost:8080/api/v1/todos/$TODO_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected response (200 OK):**

```json
{
  "success": true,
  "message": "Todo berhasil dihapus",
  "data": null
}
```

**Verify soft delete:**

```bash
$ curl -X GET http://localhost:8080/api/v1/todos/$TODO_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected: 404 Not Found**

- [ ] Status code 200 on delete
- [ ] Todo not accessible after deletion
- [ ] Soft delete confirmed

---

## Step 3: Docker Deployment Verification

### 3.1 Stop Local Server

```bash
# Press Ctrl+C to stop go run
```

- [ ] Local server stopped

### 3.2 Build and Start with Docker

```bash
$ docker-compose down
$ docker-compose up --build
```

**Expected output:**

```
[+] Building 45.2s (17/17) FINISHED
 => [internal] load build definition from Dockerfile
 => => transferring dockerfile: 523B
 => [internal] load .dockerignore
 => [stage-0 1/6] FROM golang:1.22-alpine
 => [stage-1 1/4] FROM alpine:latest
 ...
 => exporting to image
 => => naming to docker.io/library/20-mini-project-api

[+] Running 3/3
 ✔ Network 20-mini-project_default      Created
 ✔ Container 20-mini-project-postgres-1 Healthy
 ✔ Container 20-mini-project-api-1      Started
```

- [ ] Docker build successful
- [ ] PostgreSQL healthy
- [ ] API container started

### 3.3 Check Container Logs

```bash
$ docker-compose logs api
```

**Expected:**

```
api-1  | [Config] Loaded configuration
api-1  | [Database] Connecting to PostgreSQL...
api-1  | [Database] Connected successfully
api-1  | [Server] Starting on :8080
```

- [ ] No errors in logs
- [ ] Database connected

### 3.4 Check Container Health

```bash
$ docker-compose ps
```

**Expected:**

```
NAME                            STATUS                   PORTS
20-mini-project-postgres-1      Up (healthy)             5432/tcp
20-mini-project-api-1           Up (healthy)             0.0.0.0:8080->8080/tcp
```

- [ ] Both containers healthy
- [ ] Port 8080 exposed

### 3.5 Test Docker Deployment

```bash
$ curl http://localhost:8080/health
```

**Expected:**

```json
{
  "success": true,
  "message": "Server is running",
  "data": {
    "database": "connected",
    "status": "healthy"
  }
}
```

- [ ] API responds
- [ ] Health check passes

### 3.6 Repeat API Tests

Run all cURL commands from Step 2 again:

- [ ] Registration works
- [ ] Login works
- [ ] Protected routes work
- [ ] Todo CRUD works
- [ ] Filters work

---

## Step 4: Test Suite Verification

### 4.1 Run Full Test Suite

```bash
$ go test ./tests/... -v
```

**Expected:**

```
=== RUN   TestAuthTestSuite
--- PASS: TestAuthTestSuite (1.23s)
=== RUN   TestTodoTestSuite
--- PASS: TestTodoTestSuite (0.95s)
PASS
ok      ...     2.180s
```

- [ ] All tests pass
- [ ] No test failures

### 4.2 Generate Coverage Report

```bash
$ go test ./tests/... -coverprofile=coverage.out
$ go tool cover -html=coverage.out
```

- [ ] Coverage report generated
- [ ] Coverage ≥ 60%

---

## Step 5: Production Readiness Checklist

### 5.1 Security

- [ ] **JWT Secret**: Change from default in production
  ```bash
  # Generate secure secret
  $ openssl rand -base64 32
  ```
- [ ] **Database Password**: Use strong password
- [ ] **CORS**: Restrict allowed origins in production
  ```go
  // In middleware/error.go
  c.Writer.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
  ```
- [ ] **HTTPS**: Use TLS/SSL certificates
- [ ] **Rate Limiting**: Implement rate limiting middleware
- [ ] **Input Validation**: Review all DTO validation tags

### 5.2 Environment Configuration

**Create `.env` file for production:**

```bash
# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=todolist_production

# JWT
JWT_SECRET=your-very-secure-jwt-secret-generated-with-openssl

# Server
SERVER_PORT=8080
GIN_MODE=release

# Optional: Logging
LOG_LEVEL=info
LOG_FILE=/var/log/api.log
```

- [ ] `.env` created
- [ ] Secure values set
- [ ] `.env` in `.gitignore`

### 5.3 Database

- [ ] **Migrations Applied**: Run `001_create_tables.sql`
- [ ] **Indexes Created**: Verify all indexes exist
  ```sql
  SELECT * FROM pg_indexes WHERE tablename IN ('users', 'todos');
  ```
- [ ] **Backup Strategy**: Set up automated backups
  ```bash
  # Manual backup
  $ pg_dump -U postgres todolist_db > backup.sql
  ```
- [ ] **Connection Pooling**: Configure GORM connection pool
  ```go
  sqlDB, _ := db.DB()
  sqlDB.SetMaxIdleConns(10)
  sqlDB.SetMaxOpenConns(100)
  sqlDB.SetConnMaxLifetime(time.Hour)
  ```

### 5.4 Monitoring & Logging

- [ ] **Logging**: Implement structured logging
  ```go
  // Use zerolog or logrus
  log.Info().
      Str("method", c.Request.Method).
      Str("path", c.Request.URL.Path).
      Int("status", c.Writer.Status()).
      Msg("Request processed")
  ```
- [ ] **Metrics**: Add Prometheus metrics
- [ ] **Health Checks**: Verify `/health` endpoint

- [ ] **Error Tracking**: Integrate Sentry or similar

### 5.5 Performance

- [ ] **Database Indexes**: Review query performance
  ```sql
  EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 1 AND status = 'pending';
  ```
- [ ] **Caching**: Consider Redis for sessions/cache
- [ ] **Compression**: Enable gzip middleware
  ```go
  router.Use(gzip.Gzip(gzip.DefaultCompression))
  ```
- [ ] **Connection Limits**: Configure max connections

### 5.6 Docker Production

**Update `docker-compose.yml` for production:**

```yaml
version: "3.8"
services:
  postgres:
    restart: always
    deploy:
      resources:
        limits:
          cpus: "2"
          memory: 2G
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups # For backups

  api:
    restart: always
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 512M
    environment:
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET} # From .env file
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
      timeout: 10s
      retries: 3
      start_period: 40s
```

- [ ] Restart policies set
- [ ] Resource limits configured
- [ ] Health checks enabled
- [ ] Secrets from environment

### 5.7 Build Production Binary

```bash
# Build for Linux
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/api-linux cmd/api/main.go

# Build for Windows
$ GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/api.exe cmd/api/main.go

# Check binary size
$ ls -lh bin/
```

- [ ] Binary builds successfully
- [ ] Binary size optimized (~15-20MB)

---

## Step 6: Deployment Steps

### 6.1 Deploy to VPS/Cloud

**Example deployment to Ubuntu VPS:**

```bash
# 1. Copy files to server
$ scp -r . user@your-server:/opt/todolist-api

# 2. SSH to server
$ ssh user@your-server

# 3. Install Docker & Docker Compose
$ sudo apt update
$ sudo apt install docker.io docker-compose -y

# 4. Navigate to project
$ cd /opt/todolist-api

# 5. Set environment variables
$ nano .env  # Edit with production values

# 6. Start with Docker Compose
$ sudo docker-compose up -d

# 7. Check logs
$ sudo docker-compose logs -f

# 8. Test deployment
$ curl http://localhost:8080/health
```

- [ ] Files copied to server
- [ ] Docker installed
- [ ] Environment configured
- [ ] Containers running
- [ ] API accessible

### 6.2 Setup Nginx Reverse Proxy (Optional)

**Install Nginx:**

```bash
$ sudo apt install nginx -y
```

**Configure Nginx:**

```bash
$ sudo nano /etc/nginx/sites-available/todolist-api
```

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Enable site:**

```bash
$ sudo ln -s /etc/nginx/sites-available/todolist-api /etc/nginx/sites-enabled/
$ sudo nginx -t
$ sudo systemctl restart nginx
```

- [ ] Nginx installed
- [ ] Reverse proxy configured
- [ ] API accessible via domain

### 6.3 Setup SSL with Let's Encrypt

```bash
$ sudo apt install certbot python3-certbot-nginx -y
$ sudo certbot --nginx -d api.yourdomain.com
```

- [ ] SSL certificate installed
- [ ] HTTPS enabled
- [ ] Auto-renewal configured

---

## Step 7: Post-Deployment Verification

### 7.1 Test Production API

```bash
# Replace with your domain
$ curl https://api.yourdomain.com/health

# Test registration
$ curl -X POST https://api.yourdomain.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","full_name":"Test User"}'

# Test login
$ curl -X POST https://api.yourdomain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

- [ ] Health check responds
- [ ] Registration works
- [ ] Login works
- [ ] HTTPS certificate valid

### 7.2 Monitor Logs

```bash
# Container logs
$ sudo docker-compose logs -f api

# Nginx logs
$ sudo tail -f /var/log/nginx/access.log
$ sudo tail -f /var/log/nginx/error.log
```

- [ ] No errors in logs
- [ ] Requests being processed

### 7.3 Database Backup

```bash
# Create backup script
$ sudo nano /opt/backup-db.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/todolist-api/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
docker exec 20-mini-project-postgres-1 pg_dump -U postgres todolist_db > "$BACKUP_DIR/backup_$TIMESTAMP.sql"
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete  # Keep 7 days
```

**Schedule with cron:**

```bash
$ sudo chmod +x /opt/backup-db.sh
$ sudo crontab -e
```

Add line:

```
0 2 * * * /opt/backup-db.sh  # Daily at 2 AM
```

- [ ] Backup script created
- [ ] Cron job scheduled
- [ ] Test backup manually

---

## ✅ Final Completion Checklist

### Development

- [ ] Project structure matches specification
- [ ] All dependencies installed
- [ ] Configuration layer implemented
- [ ] Models with GORM tags created
- [ ] DTOs with validation created
- [ ] Utilities (password, JWT) implemented
- [ ] Middleware (auth, logger, CORS, error) implemented
- [ ] Repositories with specific queries created
- [ ] Services with business logic implemented
- [ ] Handlers with error mapping created
- [ ] Routes properly organized
- [ ] Main application with DI pattern
- [ ] Database migrations created
- [ ] Dockerfile multi-stage build configured
- [ ] docker-compose.yml with health checks
- [ ] Tests implemented with testify suite
- [ ] All tests passing

### Local Verification

- [ ] PostgreSQL running
- [ ] API server starts successfully
- [ ] Health endpoint responds
- [ ] User registration works
- [ ] Login returns JWT token
- [ ] Protected routes require authentication
- [ ] Todo CRUD operations work
- [ ] Filters work correctly
- [ ] Soft delete confirmed
- [ ] Error handling correct

### Docker Verification

- [ ] Docker build successful
- [ ] Containers start correctly
- [ ] Health checks passing
- [ ] API accessible on port 8080
- [ ] All API endpoints work in Docker
- [ ] Container logs show no errors

### Testing

- [ ] Test suite runs successfully
- [ ] Coverage ≥ 60%
- [ ] Auth tests pass
- [ ] Todo tests pass
- [ ] Test independence verified

### Production Readiness

- [ ] JWT secret changed from default
- [ ] Strong database password set
- [ ] CORS configured for production
- [ ] HTTPS enabled
- [ ] Rate limiting considered
- [ ] Input validation reviewed
- [ ] Environment variables configured
- [ ] Database migrations applied
- [ ] Indexes verified
- [ ] Backup strategy implemented
- [ ] Logging configured
- [ ] Monitoring setup
- [ ] Health checks working
- [ ] Resource limits set
- [ ] Restart policies configured

### Deployment

- [ ] Files copied to server
- [ ] Docker installed on server
- [ ] Environment configured for production
- [ ] Containers running on server
- [ ] Nginx reverse proxy setup (optional)
- [ ] SSL certificate installed (optional)
- [ ] Production API responding
- [ ] All endpoints tested in production
- [ ] Logs monitored
- [ ] Backup cron job scheduled

---

## 🐛 Common Production Issues

### Issue: "Connection refused"

**Solution:**

```bash
# Check if containers are running
$ docker-compose ps

# Check logs
$ docker-compose logs api

# Verify network
$ docker network ls
$ docker network inspect 20-mini-project_default
```

### Issue: "Database connection failed"

**Solution:**

```bash
# Check PostgreSQL health
$ docker exec -it 20-mini-project-postgres-1 pg_isready

# Verify environment variables
$ docker-compose exec api env | grep DB_

# Check connection string
$ docker-compose logs postgres
```

### Issue: "Token invalid/expired"

**Solution:**

- Tokens expire after 24 hours (default)
- User must login again to get new token
- Check JWT_SECRET matches between containers
- Verify token format: "Bearer <token>"

### Issue: "502 Bad Gateway (Nginx)"

**Solution:**

```bash
# Check if API container is running
$ sudo docker-compose ps

# Test direct connection
$ curl http://localhost:8080/health

# Check Nginx config
$ sudo nginx -t

# Check Nginx logs
$ sudo tail -f /var/log/nginx/error.log
```

### Issue: "Out of memory"

**Solution:**

- Increase Docker memory limits in docker-compose.yml
- Optimize queries (add indexes)
- Implement connection pooling
- Consider pagination for large datasets

---

## 📊 Final Progress

```
[████████████████████████████████████████] 100% Complete
```

**✅ Project Recreation Complete!**

All 16 steps have been completed:

1. ✅ Initial setup
2. ✅ Dependencies installation
3. ✅ Configuration layer
4. ✅ Database models
5. ✅ DTOs and validation
6. ✅ Utilities (password/JWT)
7. ✅ Middleware
8. ✅ Repositories
9. ✅ Services
10. ✅ Handlers
11. ✅ Routes
12. ✅ Main application
13. ✅ Database migrations
14. ✅ Docker configuration
15. ✅ Testing
16. ✅ **Final verification** ← You are here

---

## 🎉 Congratulations!

You have successfully:

- ✅ Built a complete Go REST API with Clean Architecture
- ✅ Implemented JWT authentication
- ✅ Created CRUD operations with validation
- ✅ Dockerized the application
- ✅ Written comprehensive tests
- ✅ Deployed to production (optional)

### Next Steps

**To deepen your knowledge:**

1. Add more features (comments, tags, attachments)
2. Implement WebSocket for real-time updates
3. Add pagination and search
4. Integrate Redis for caching
5. Implement GraphQL API
6. Add API documentation with Swagger
7. Implement rate limiting
8. Add email notifications
9. Create a frontend (React/Vue)
10. Deploy to Kubernetes

**Resources:**

- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Go Testing](https://golang.org/pkg/testing/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Previous:** [15-testing.md](15-testing.md)  
**Return to:** [00-overview.md](00-overview.md)  
**Index:** [README.md](README.md)
