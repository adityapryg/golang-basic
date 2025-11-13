# 13. Middleware dan Error Handling

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Apa itu middleware dan bagaimana cara kerjanya
2. Cara membuat custom middleware di Gin
3. Penggunaan middleware untuk logging, autentikasi, dan CORS
4. Cara menerapkan middleware secara global atau per-group
5. Konsep chaining middleware dan `c.Next()` / `c.Abort()`

## Penjelasan

### Apa itu Middleware?

Middleware adalah fungsi yang dieksekusi **sebelum** atau **setelah** handler utama. Middleware berguna untuk:

- **Logging** - Mencatat setiap request
- **Autentikasi** - Memvalidasi token/kredensial
- **CORS** - Menangani Cross-Origin Resource Sharing
- **Rate Limiting** - Membatasi jumlah request
- **Error Recovery** - Menangani panic agar server tidak crash

### Anatomi Middleware di Gin

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Kode sebelum request diproses

        c.Next() // Lanjutkan ke handler/middleware berikutnya

        // Kode setelah request diproses
    }
}
```

### Perbedaan `c.Next()` dan `c.Abort()`

- **`c.Next()`** - Melanjutkan eksekusi ke handler berikutnya
- **`c.Abort()`** - Menghentikan eksekusi, handler berikutnya tidak dipanggil

### Middleware yang Dicontohkan

1. **LoggerMiddleware** - Mencatat method, path, status, dan durasi request
2. **AuthMiddleware** - Validasi token di header Authorization
3. **CORSMiddleware** - Mengatur header CORS untuk cross-origin requests
4. **RateLimitMiddleware** - Membatasi jumlah request per IP (simplified)

### Penggunaan Middleware

```go
// Global middleware (semua route)
router.Use(LoggerMiddleware())

// Group middleware (hanya route tertentu)
protected := router.Group("/")
protected.Use(AuthMiddleware())
{
    protected.GET("/private", handler)
}
```

## Prasyarat

### Install Dependencies

```bash
cd 13-middleware
go get -u github.com/gin-gonic/gin
```

Dependencies yang akan diinstall:

- `github.com/gin-gonic/gin` - Gin web framework

## Cara Menjalankan

```bash
cd 13-middleware
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Cara Testing

### 1. Test Public Endpoint (tanpa autentikasi)

```bash
curl http://localhost:8080/public
```

**Response:**

```json
{
  "message": "Ini adalah endpoint public, tidak perlu autentikasi"
}
```

### 2. Test Private Endpoint (dengan token yang benar)

```bash
curl -H "Authorization: Bearer secret-token-123" http://localhost:8080/private
```

**Response:**

```json
{
  "message": "Ini adalah endpoint private",
  "user_id": "12345",
  "username": "john_doe"
}
```

### 3. Test Private Endpoint (tanpa token)

```bash
curl http://localhost:8080/private
```

**Response:**

```json
{
  "error": "Token tidak ditemukan"
}
```

### 4. Test Rate Limiting

Jalankan perintah ini lebih dari 5 kali:

```bash
curl http://localhost:8080/api/data
```

Setelah request ke-6, Anda akan mendapat:

```json
{
  "error": "Terlalu banyak request, coba lagi nanti"
}
```

## Output di Terminal Server

Setiap request akan dicatat dengan format:

```
[2025-11-13 10:30:15] GET /public - Status: 200 - Duration: 245.3µs
[2025-11-13 10:30:20] GET /private - Status: 200 - Duration: 180.7µs
```

## Konsep Penting

### 1. Middleware Chaining

Middleware dieksekusi secara berurutan:

```
Request → Middleware1 → Middleware2 → Handler → Middleware2 (after) → Middleware1 (after) → Response
```

### 2. Context Data

Middleware dapat menyimpan data ke context:

```go
c.Set("user_id", "12345")      // Di middleware
userID, _ := c.Get("user_id")  // Di handler
```

### 3. Early Return

Gunakan `c.Abort()` untuk menghentikan eksekusi:

```go
if !isValid {
    c.JSON(401, gin.H{"error": "Unauthorized"})
    c.Abort()
    return
}
```

## Catatan Tambahan

- Middleware ini adalah contoh sederhana untuk pembelajaran
- Untuk production, gunakan:
  - JWT untuk autentikasi (bukan token string sederhana)
  - Redis untuk rate limiting yang lebih robust
  - Structured logging library (seperti Zap atau Logrus)
- Middleware di Gin sangat powerful dan bisa dikombinasikan sesuai kebutuhan

## Referensi

- [Gin Middleware Documentation](https://gin-gonic.com/docs/examples/custom-middleware/)
- [Gin Built-in Middlewares](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)
