package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Secret key untuk signing JWT (dalam production, simpan di environment variable)
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// User model (simplified, biasanya dari database)
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" agar tidak muncul di JSON response
	Role     string `json:"role"`
}

// In-memory user storage (dalam production gunakan database)
var users = map[string]User{
	"admin": {
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Password: "$2a$10$YourHashedPasswordHere", // "admin123"
		Role:     "admin",
	},
	"user": {
		ID:       2,
		Username: "user",
		Email:    "user@example.com",
		Password: "$2a$10$YourHashedPasswordHere", // "user123"
		Role:     "user",
	},
}

// JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Response format
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// generateToken membuat JWT token
func generateToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid 24 jam

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "golang-demo-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// hashPassword meng-hash password dengan bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash memverifikasi password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthMiddleware adalah middleware untuk validasi JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "Format token tidak valid",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse dan validasi token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "Token tidak valid atau sudah expired",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Set user info ke context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware adalah middleware untuk validasi role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, Response{
				Success: false,
				Message: "Role tidak ditemukan",
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		allowed := false
		for _, r := range allowedRoles {
			if r == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, Response{
				Success: false,
				Message: "Akses ditolak: role tidak memiliki permission",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Handler: Register
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Cek apakah username sudah ada
	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Username sudah digunakan",
		})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal meng-hash password",
			Error:   err.Error(),
		})
		return
	}

	// Simpan user (dalam production, simpan ke database)
	newUser := User{
		ID:       uint(len(users) + 1),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}
	users[req.Username] = newUser

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Registrasi berhasil",
		Data: gin.H{
			"id":       newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
			"role":     newUser.Role,
		},
	})
}

// Handler: Login
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Cek user exists
	user, exists := users[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Username atau password salah",
		})
		return
	}

	// Validasi password
	if !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Username atau password salah",
		})
		return
	}

	// Generate token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login berhasil",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

// Handler: Get Profile (protected)
func getProfile(c *gin.Context) {
	username, _ := c.Get("username")
	user := users[username.(string)]

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Profil user",
		Data: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Handler: Admin Only
func adminOnly(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Welcome Admin!",
		Data: gin.H{
			"users": users,
		},
	})
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   JWT AUTHENTICATION")
	fmt.Println("===========================================\n")

	// Hash password untuk demo users
	adminHash, _ := hashPassword("admin123")
	userHash, _ := hashPassword("user123")
	users["admin"] = User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Password: adminHash,
		Role:     "admin",
	}
	users["user"] = User{
		ID:       2,
		Username: "user",
		Email:    "user@example.com",
		Password: userHash,
		Role:     "user",
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
	}

	// Protected routes (perlu JWT)
	api := router.Group("/api")
	api.Use(AuthMiddleware())
	{
		api.GET("/profile", getProfile)

		// Admin only routes
		admin := api.Group("/admin")
		admin.Use(RoleMiddleware("admin"))
		{
			admin.GET("/users", adminOnly)
		}
	}

	fmt.Println("📡 Server berjalan di http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  Public:")
	fmt.Println("    POST /api/auth/register")
	fmt.Println("    POST /api/auth/login")
	fmt.Println("  Protected (perlu token):")
	fmt.Println("    GET /api/profile")
	fmt.Println("  Admin only:")
	fmt.Println("    GET /api/admin/users")
	fmt.Println("\nDefault users:")
	fmt.Println("  admin / admin123 (role: admin)")
	fmt.Println("  user / user123 (role: user)")
	fmt.Println("\nTekan Ctrl+C untuk menghentikan server\n")

	router.Run(":8080")
}
