package route

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	todoHandler *handler.TodoHandler,
) {
	// Health check endpoint (public)
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware()) // Apply JWT middleware
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
		}

		// Todo routes (protected)
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware())
		{
			todos.POST("", todoHandler.Create)
			todos.GET("", todoHandler.GetAll)
			todos.GET("/:id", todoHandler.GetByID)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}
	}
}
