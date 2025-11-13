package main

import (
	"log"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/handler"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/route"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("===========================================")
	log.Println("   STARTING TODO REST API SERVER")
	log.Println("===========================================")

	// Load configuration
	cfg := config.LoadConfig()
	gin.SetMode(cfg.GinMode)

	// Initialize database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("✓ Configuration loaded")

	// ============================================
	// DEPENDENCY INJECTION PATTERN
	// ============================================

	// Layer 1: Initialize Repositories (Data Access Layer)
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	log.Println("✓ Repositories initialized")

	// Layer 2: Initialize Services (Business Logic Layer)
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)
	log.Println("✓ Services initialized")

	// Layer 3: Initialize Handlers (HTTP Layer)
	userHandler := handler.NewUserHandler(authService)
	healthHandler := handler.NewHealthHandler(db)
	todoHandler := handler.NewTodoHandler(todoService)
	log.Println("✓ Handlers initialized")

	// ============================================
	// GIN ROUTER SETUP
	// ============================================

	router := gin.Default()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// Setup routes
	route.SetupRoutes(router, userHandler, healthHandler, todoHandler)
	log.Println("✓ Routes configured")

	// ============================================
	// START SERVER
	// ============================================

	log.Println("===========================================")
	log.Printf("   Server running on :%s", cfg.ServerPort)
	log.Println("   API Endpoint: http://localhost:" + cfg.ServerPort)
	log.Println("===========================================")

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
