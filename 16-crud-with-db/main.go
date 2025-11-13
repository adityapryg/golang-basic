package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Product model
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,gt=0"`
	Stock       int       `gorm:"not null;default:0" json:"stock" binding:"required,gte=0"`
	CategoryID  uint      `gorm:"not null" json:"category_id" binding:"required"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category model
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name" binding:"required"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Response format
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var db *gorm.DB

// initDB initializes database connection
func initDB() error {
	dsn := "host=localhost user=postgres password=postgres dbname=golang_demo port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("gagal koneksi database: %w", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&Category{}, &Product{})
	if err != nil {
		return fmt.Errorf("migrasi gagal: %w", err)
	}

	// Seed categories if empty
	var count int64
	db.Model(&Category{}).Count(&count)
	if count == 0 {
		categories := []Category{
			{Name: "Elektronik", Description: "Peralatan elektronik"},
			{Name: "Pakaian", Description: "Pakaian dan fashion"},
			{Name: "Makanan", Description: "Makanan dan minuman"},
		}
		db.Create(&categories)
	}

	return nil
}

// Category handlers
func createCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal membuat kategori",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Kategori berhasil dibuat",
		Data:    category,
	})
}

func getAllCategories(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal mengambil kategori",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Data kategori berhasil diambil",
		Data: gin.H{
			"total":      len(categories),
			"categories": categories,
		},
	})
}

// Product handlers
func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Validasi category exists
	var category Category
	if err := db.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Kategori tidak ditemukan",
		})
		return
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal membuat produk",
			Error:   err.Error(),
		})
		return
	}

	// Load category untuk response
	db.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Produk berhasil dibuat",
		Data:    product,
	})
}

func getAllProducts(c *gin.Context) {
	var products []Product

	// Preload category
	if err := db.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal mengambil produk",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Data produk berhasil diambil",
		Data: gin.H{
			"total":    len(products),
			"products": products,
		},
	})
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")

	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Produk tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Produk ditemukan",
		Data:    product,
	})
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Produk tidak ditemukan",
		})
		return
	}

	var updates Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Validasi category exists jika diubah
	if updates.CategoryID != 0 {
		var category Category
		if err := db.First(&category, updates.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Kategori tidak ditemukan",
			})
			return
		}
	}

	if err := db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal mengupdate produk",
			Error:   err.Error(),
		})
		return
	}

	// Load category untuk response
	db.Preload("Category").First(&product, id)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Produk berhasil diupdate",
		Data:    product,
	})
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Produk tidak ditemukan",
		})
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Gagal menghapus produk",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Produk berhasil dihapus",
	})
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   CRUD API DENGAN DATABASE")
	fmt.Println("===========================================\n")

	// Initialize database
	fmt.Println("🔌 Menghubungkan ke database...")
	if err := initDB(); err != nil {
		log.Fatal("❌ Gagal inisialisasi database:", err)
	}
	fmt.Println("✅ Database berhasil terkoneksi!")

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		// Categories
		api.POST("/categories", createCategory)
		api.GET("/categories", getAllCategories)

		// Products
		api.POST("/products", createProduct)
		api.GET("/products", getAllProducts)
		api.GET("/products/:id", getProductByID)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Informasi
	fmt.Println("\n📡 Server berjalan di http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  Categories:")
	fmt.Println("    POST   /api/v1/categories")
	fmt.Println("    GET    /api/v1/categories")
	fmt.Println("  Products:")
	fmt.Println("    POST   /api/v1/products")
	fmt.Println("    GET    /api/v1/products")
	fmt.Println("    GET    /api/v1/products/:id")
	fmt.Println("    PUT    /api/v1/products/:id")
	fmt.Println("    DELETE /api/v1/products/:id")
	fmt.Println("\nTekan Ctrl+C untuk menghentikan server\n")

	router.Run(":8080")
}
