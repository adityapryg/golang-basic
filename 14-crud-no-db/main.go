package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Product adalah model untuk data produk
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Stock       int       `json:"stock" binding:"required,gte=0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Response adalah format standar untuk response API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// In-memory storage
var (
	products   = make(map[int]Product)
	productsMu sync.RWMutex
	nextID     = 1
)

// createProduct adalah handler untuk membuat produk baru
func createProduct(c *gin.Context) {
	var product Product

	// Bind dan validasi JSON request
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	// Lock untuk thread safety
	productsMu.Lock()
	defer productsMu.Unlock()

	// Set ID dan timestamp
	product.ID = nextID
	nextID++
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Simpan ke map
	products[product.ID] = product

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Produk berhasil dibuat",
		Data:    product,
	})
}

// getAllProducts adalah handler untuk mendapatkan semua produk
func getAllProducts(c *gin.Context) {
	productsMu.RLock()
	defer productsMu.RUnlock()

	// Convert map ke slice
	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Data produk berhasil diambil",
		Data: gin.H{
			"total":    len(productList),
			"products": productList,
		},
	})
}

// getProductByID adalah handler untuk mendapatkan produk berdasarkan ID
func getProductByID(c *gin.Context) {
	// Parse ID dari URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
		return
	}

	productsMu.RLock()
	product, exists := products[id]
	productsMu.RUnlock()

	if !exists {
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

// updateProduct adalah handler untuk mengupdate produk
func updateProduct(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Bind JSON request
	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	productsMu.Lock()
	defer productsMu.Unlock()

	// Cek apakah produk ada
	product, exists := products[id]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Produk tidak ditemukan",
		})
		return
	}

	// Update fields (preserve ID dan CreatedAt)
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.Stock = updatedProduct.Stock
	product.UpdatedAt = time.Now()

	products[id] = product

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Produk berhasil diupdate",
		Data:    product,
	})
}

// deleteProduct adalah handler untuk menghapus produk
func deleteProduct(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
		return
	}

	productsMu.Lock()
	defer productsMu.Unlock()

	// Cek apakah produk ada
	if _, exists := products[id]; !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Produk tidak ditemukan",
		})
		return
	}

	// Hapus dari map
	delete(products, id)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Produk berhasil dihapus",
	})
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   CRUD API TANPA DATABASE")
	fmt.Println("===========================================\n")

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		// Products endpoints
		api.POST("/products", createProduct)
		api.GET("/products", getAllProducts)
		api.GET("/products/:id", getProductByID)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Informasi
	fmt.Println("Server berjalan di http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  POST   /api/v1/products      - Buat produk baru")
	fmt.Println("  GET    /api/v1/products      - Lihat semua produk")
	fmt.Println("  GET    /api/v1/products/:id  - Lihat produk by ID")
	fmt.Println("  PUT    /api/v1/products/:id  - Update produk")
	fmt.Println("  DELETE /api/v1/products/:id  - Hapus produk")
	fmt.Println("\nContoh testing dengan curl ada di README.md")
	fmt.Println("\nTekan Ctrl+C untuk menghentikan server\n")

	// Jalankan server
	router.Run(":8080")
}
