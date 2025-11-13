package middleware

import (
	"net/http"
	"strings"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware memvalidasi JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Format token tidak valid",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Token tidak valid atau expired",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Simpan user info di context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID mengambil user ID dari context
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername mengambil username dari context
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
