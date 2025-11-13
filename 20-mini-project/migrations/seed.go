package main

import (
	"log"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
)

// Seeder populates database with test data
func main() {
	log.Println("===========================================")
	log.Println("   RUNNING DATABASE SEEDER")
	log.Println("===========================================")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✓ Connected to database")

	// ============================================
	// SEED USERS
	// ============================================

	users := []model.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: mustHashPassword("admin123"),
			FullName: "Administrator",
		},
		{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: mustHashPassword("password123"),
			FullName: "John Doe",
		},
		{
			Username: "jane_smith",
			Email:    "jane@example.com",
			Password: mustHashPassword("password123"),
			FullName: "Jane Smith",
		},
	}

	for _, user := range users {
		// Check if user already exists
		var existing model.User
		result := db.Where("username = ?", user.Username).First(&existing)

		if result.Error == nil {
			log.Printf("⚠ User '%s' already exists, skipping...", user.Username)
			continue
		}

		// Create user
		if err := db.Create(&user).Error; err != nil {
			log.Printf("✗ Failed to create user '%s': %v", user.Username, err)
			continue
		}

		log.Printf("✓ Created user: %s (%s)", user.Username, user.Email)
	}

	// ============================================
	// SEED TODOS
	// ============================================

	// Get user IDs for todo assignment
	var adminUser, johnUser model.User
	db.Where("username = ?", "admin").First(&adminUser)
	db.Where("username = ?", "john_doe").First(&johnUser)

	todos := []model.Todo{
		{
			Title:       "Complete project documentation",
			Description: "Write comprehensive README and API documentation",
			Status:      model.TodoStatusInProgress,
			Priority:    5,
			UserID:      adminUser.ID,
		},
		{
			Title:       "Review pull requests",
			Description: "Review and merge pending pull requests",
			Status:      model.TodoStatusPending,
			Priority:    4,
			UserID:      adminUser.ID,
		},
		{
			Title:       "Fix authentication bug",
			Description: "Resolve token expiration issue",
			Status:      model.TodoStatusCompleted,
			Priority:    5,
			UserID:      adminUser.ID,
		},
		{
			Title:       "Learn Clean Architecture",
			Description: "Study Clean Architecture patterns in Go",
			Status:      model.TodoStatusInProgress,
			Priority:    3,
			UserID:      johnUser.ID,
		},
		{
			Title:       "Build REST API",
			Description: "Create a REST API using Gin framework",
			Status:      model.TodoStatusPending,
			Priority:    4,
			UserID:      johnUser.ID,
		},
		{
			Title:       "Write unit tests",
			Description: "Add unit tests for all services",
			Status:      model.TodoStatusPending,
			Priority:    3,
			UserID:      johnUser.ID,
		},
	}

	for _, todo := range todos {
		// Check if todo already exists
		var existing model.Todo
		result := db.Where("title = ? AND user_id = ?", todo.Title, todo.UserID).First(&existing)

		if result.Error == nil {
			log.Printf("⚠ Todo '%s' already exists, skipping...", todo.Title)
			continue
		}

		// Create todo
		if err := db.Create(&todo).Error; err != nil {
			log.Printf("✗ Failed to create todo '%s': %v", todo.Title, err)
			continue
		}

		log.Printf("✓ Created todo: %s (User: %d, Status: %s)", todo.Title, todo.UserID, todo.Status)
	}

	log.Println("===========================================")
	log.Println("   SEEDING COMPLETED")
	log.Println("===========================================")
	log.Println("\nTest Credentials:")
	log.Println("  Username: admin      | Password: admin123")
	log.Println("  Username: john_doe   | Password: password123")
	log.Println("  Username: jane_smith | Password: password123")
	log.Println("===========================================")
}

// mustHashPassword hashes password or panics on error
func mustHashPassword(password string) string {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return hashed
}
