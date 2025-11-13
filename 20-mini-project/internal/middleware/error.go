package middleware

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/gin-gonic/gin"
)

// ErrorHandler menangani error yang terjadi
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Cek jika ada error
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(500, dto.ErrorResponse{
				Success: false,
				Message: "Internal server error",
				Error:   err.Error(),
			})
		}
	}
}

// CORSMiddleware menangani CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
