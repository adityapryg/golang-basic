package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware adalah middleware custom untuk logging request
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Waktu sebelum request diproses
		startTime := time.Now()

		// Proses request
		c.Next()

		// Waktu setelah request diproses
		duration := time.Since(startTime)

		// Log informasi request
		fmt.Printf("[%s] %s %s - Status: %d - Duration: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}

// AuthMiddleware adalah middleware untuk autentikasi sederhana
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cek header Authorization
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak ditemukan",
			})
			c.Abort() // Hentikan eksekusi handler berikutnya
			return
		}

		// Validasi token sederhana (dalam real-world gunakan JWT)
		if token != "Bearer secret-token-123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid",
			})
			c.Abort()
			return
		}

		// Set user info ke context (bisa diakses di handler)
		c.Set("user_id", "12345")
		c.Set("username", "john_doe")

		c.Next() // Lanjutkan ke handler berikutnya
	}
}

// CORSMiddleware adalah middleware untuk menangani CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware adalah middleware sederhana untuk rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// Map untuk menyimpan jumlah request per IP
	requestCount := make(map[string]int)
	const maxRequests = 5 // Maksimal 5 request per IP

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Increment counter
		requestCount[clientIP]++

		// Cek apakah sudah melebihi limit
		if requestCount[clientIP] > maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak request, coba lagi nanti",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI MIDDLEWARE DI GIN")
	fmt.Println("===========================================\n")

	// Buat router Gin dengan mode release
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware global (berlaku untuk semua route)
	router.Use(LoggerMiddleware())
	router.Use(gin.Recovery()) // Built-in recovery middleware
	router.Use(CORSMiddleware())

	// Route tanpa autentikasi
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Selamat datang di API dengan Middleware",
			"endpoints": gin.H{
				"public":  "/public",
				"private": "/private (perlu Authorization header)",
			},
		})
	})

	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ini adalah endpoint public, tidak perlu autentikasi",
		})
	})

	// Group route dengan autentikasi
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/private", func(c *gin.Context) {
			// Ambil data user dari context
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")

			c.JSON(http.StatusOK, gin.H{
				"message":  "Ini adalah endpoint private",
				"user_id":  userID,
				"username": username,
			})
		})

		protected.GET("/profile", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"message": "Profil user",
				"data": gin.H{
					"username": username,
					"email":    "john@example.com",
					"role":     "admin",
				},
			})
		})
	}

	// Route dengan rate limiting
	limited := router.Group("/api")
	limited.Use(RateLimitMiddleware())
	{
		limited.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Data berhasil diambil",
				"data":    []string{"item1", "item2", "item3"},
			})
		})
	}

	// Informasi cara testing
	fmt.Println("Server berjalan di http://localhost:8080")
	fmt.Println("\nCara testing:")
	fmt.Println("1. Public endpoint (tidak perlu token):")
	fmt.Println("   curl http://localhost:8080/public")
	fmt.Println("\n2. Private endpoint (perlu token):")
	fmt.Println("   curl -H \"Authorization: Bearer secret-token-123\" http://localhost:8080/private")
	fmt.Println("\n3. Rate limited endpoint (maksimal 5 request):")
	fmt.Println("   curl http://localhost:8080/api/data")
	fmt.Println("\nTekan Ctrl+C untuk menghentikan server\n")

	// Jalankan server
	router.Run(":8080")
}
