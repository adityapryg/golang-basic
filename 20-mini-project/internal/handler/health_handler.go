package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck checks if API and database are running
// @Summary Health check
// @Description Check if API and database connection are healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "error",
			"message":  "database connection failed",
			"database": "disconnected",
		})
		return
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "error",
			"message":  "database ping failed",
			"database": "disconnected",
		})
		return
	}

	// All checks passed
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "API is running",
		"database": "connected",
	})
}
