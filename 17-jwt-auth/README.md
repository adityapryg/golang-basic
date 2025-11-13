# 17. JWT Authentication

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Konsep JWT (JSON Web Token) dan cara kerjanya
2. Implementasi register dan login dengan JWT
3. Password hashing dengan bcrypt
4. Middleware untuk validasi JWT token
5. Role-based access control (RBAC)
6. Best practices keamanan autentikasi

## Penjelasan

### Apa itu JWT?

JWT (JSON Web Token) adalah standar untuk membuat token akses yang berisi informasi user dalam format JSON yang di-encode dan di-sign.

**Struktur JWT:**

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

Terdiri dari 3 bagian (dipisahkan `.`):

1. **Header** - Algoritma dan tipe token
2. **Payload** - Data/claims (user_id, username, role, exp, dll)
3. **Signature** - Verifikasi integritas token

### Alur Autentikasi

```
1. Register/Login → Server validasi kredensial
2. Server generate JWT token
3. Client simpan token (localStorage/cookie)
4. Client kirim token di header setiap request
5. Server validasi token di middleware
6. Jika valid, lanjutkan ke handler
```

### Password Hashing dengan bcrypt

```go
// Hash password
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verify password
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
```

**JANGAN** simpan password plain text di database!

### JWT Claims

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims  // exp, iat, iss, dll
}
```

### Middleware Chain

```
Request → AuthMiddleware → RoleMiddleware → Handler
           (validasi JWT)    (cek role)      (bisnis logic)
```

## Prasyarat

### Install Dependencies

```bash
cd 17-jwt-auth
go get -u github.com/gin-gonic/gin
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto
```

Dependencies yang akan diinstall:

- `github.com/gin-gonic/gin v1.9.1` - Gin web framework
- `github.com/golang-jwt/jwt/v5 v5.2.0` - JWT library
- `golang.org/x/crypto v0.14.0` - Password hashing (bcrypt)

## Cara Menjalankan

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Cara Testing

### 1. Register User Baru

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": {
    "id": 3,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

**Simpan token** dari response untuk request selanjutnya!

### 3. Akses Protected Endpoint (dengan token)

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**

```json
{
  "success": true,
  "message": "Profil user",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

### 4. Akses Protected Endpoint (tanpa token)

```bash
curl http://localhost:8080/api/profile
```

**Response:**

```json
{
  "success": false,
  "message": "Token tidak ditemukan"
}
```

### 5. Akses Admin Endpoint (dengan role admin)

```bash
curl http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Response:**

```json
{
  "success": true,
  "message": "Welcome Admin!",
  "data": {
    "users": [...]
  }
}
```

### 6. Akses Admin Endpoint (dengan role user)

```bash
curl http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Response:**

```json
{
  "success": false,
  "message": "Akses ditolak: role tidak memiliki permission"
}
```

## Testing dengan Postman

### Setup

1. Buat collection "JWT Auth"
2. Buat environment variable `token`

### Workflow

1. **Login** → Copy token dari response
2. **Set Authorization**:
   - Type: Bearer Token
   - Token: `{{token}}`
3. Test protected endpoints

### Postman Script (Auto-save token)

Di tab **Tests** pada request login:

```javascript
pm.test("Status code is 200", function () {
  pm.response.to.have.status(200);
});

var jsonData = pm.response.json();
pm.environment.set("token", jsonData.data.token);
```

## Security Best Practices

### 1. Secret Key

❌ Jangan hardcode:

```go
var jwtSecret = []byte("my-secret-key")
```

✅ Gunakan environment variable:

```go
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

### 2. Token Expiration

Atur expiration time yang reasonable:

```go
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
```

### 3. Password Hashing

Selalu gunakan bcrypt atau argon2:

```go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### 4. HTTPS Only

Dalam production, SELALU gunakan HTTPS untuk mencegah token dicuri.

### 5. Token Refresh

Implementasikan refresh token untuk UX yang lebih baik:

- Access token: short-lived (15 menit)
- Refresh token: long-lived (7 hari)

### 6. Token Revocation

Untuk logout yang proper, implementasikan token blacklist atau token versioning.

## Struktur JWT Token

Decode token di [jwt.io](https://jwt.io):

**Header:**

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

**Payload:**

```json
{
  "user_id": 1,
  "username": "admin",
  "role": "admin",
  "exp": 1699900800,
  "iat": 1699814400,
  "iss": "golang-demo-api"
}
```

## Role-Based Access Control (RBAC)

```go
// Middleware untuk validasi role
admin := api.Group("/admin")
admin.Use(RoleMiddleware("admin"))
{
    admin.GET("/users", adminOnly)
}
```

Role hierarchy bisa diperluas:

- `admin` - Full access
- `moderator` - Moderate content
- `user` - Basic access
- `guest` - Read only

## Common Errors & Solutions

### Error: "Token tidak valid atau sudah expired"

Token sudah expired. Login ulang untuk mendapat token baru.

### Error: "Format token tidak valid"

Header Authorization harus format: `Bearer <token>`

### Error: "Akses ditolak: role tidak memiliki permission"

User tidak memiliki role yang diperlukan untuk akses endpoint tersebut.

## Referensi

- [JWT Introduction](https://jwt.io/introduction)
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [bcrypt Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP Authentication](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
