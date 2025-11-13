package config

import (
	"fmt"
	"log"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *Config) (*gorm.DB, error) {
	// Connection string PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Koneksi ke database dengan logging
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("✓ Successfully connected to database")

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✓ Database migration completed")

	return db, nil
}
