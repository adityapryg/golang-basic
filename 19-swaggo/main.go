package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/adityapryg/golang-demo/19-swaggo/docs" // Import generated docs
)

// Product model
// @Description Product information
type Product struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Laptop Gaming"`
	Description string    `json:"description" example:"Laptop gaming dengan spesifikasi tinggi"`
	Price       float64   `json:"price" example:"15000000"`
	Stock       int       `json:"stock" example:"10"`
	CreatedAt   time.Time `json:"created_at" example:"2025-11-13T10:30:00Z"`
}

// Response model
// @Description API response format
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error message"`
}

// CreateProductRequest model
// @Description Request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"Laptop Gaming"`
	Description string  `json:"description" example:"Laptop gaming dengan spesifikasi tinggi"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"15000000"`
	Stock       int     `json:"stock" binding:"required,gte=0" example:"10"`
}

// In-memory storage
var products = []Product{
	{
		ID:          1,
		Name:        "Laptop Gaming",
		Description: "Laptop gaming dengan spesifikasi tinggi",
		Price:       15000000,
		Stock:       10,
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Name:        "Mouse Wireless",
		Description: "Mouse wireless ergonomis",
		Price:       150000,
		Stock:       50,
		CreatedAt:   time.Now(),
	},
}

var nextID = 3

// @title           Product API
// @version         1.0
// @description     REST API untuk manajemen produk dengan dokumentasi Swagger
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// GetAllProducts godoc
// @Summary      Ambil semua produk
// @Description  Mendapatkan list semua produk yang tersedia
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=[]Product}
// @Router       /products [get]
func getAllProducts(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Data produk berhasil diambil",
		Data:    products,
	})
}

// GetProductByID godoc
// @Summary      Ambil produk by ID
// @Description  Mendapatkan detail produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Response{data=Product}
// @Failure      404  {object}  Response
// @Router       /products/{id} [get]
func getProductByID(c *gin.Context) {
	id := c.Param("id")

	for _, product := range products {
		if product.ID == parseID(id) {
			c.JSON(http.StatusOK, Response{
				Success: true,
				Message: "Produk ditemukan",
				Data:    product,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: "Produk tidak ditemukan",
	})
}

// CreateProduct godoc
// @Summary      Buat produk baru
// @Description  Membuat produk baru dengan data yang diberikan
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      CreateProductRequest  true  "Product data"
// @Success      201  {object}  Response{data=Product}
// @Failure      400  {object}  Response
// @Security     BearerAuth
// @Router       /products [post]
func createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	product := Product{
		ID:          nextID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
	}
	nextID++

	products = append(products, product)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Produk berhasil dibuat",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary      Update produk
// @Description  Mengupdate data produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "Product ID"
// @Param        product  body      CreateProductRequest  true  "Updated product data"
// @Success      200  {object}  Response{data=Product}
// @Failure      400  {object}  Response
// @Failure      404  {object}  Response
// @Security     BearerAuth
// @Router       /products/{id} [put]
func updateProduct(c *gin.Context) {
	id := parseID(c.Param("id"))

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validasi gagal",
			Error:   err.Error(),
		})
		return
	}

	for i, product := range products {
		if product.ID == id {
			products[i].Name = req.Name
			products[i].Description = req.Description
			products[i].Price = req.Price
			products[i].Stock = req.Stock

			c.JSON(http.StatusOK, Response{
				Success: true,
				Message: "Produk berhasil diupdate",
				Data:    products[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: "Produk tidak ditemukan",
	})
}

// DeleteProduct godoc
// @Summary      Hapus produk
// @Description  Menghapus produk berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Response
// @Failure      404  {object}  Response
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func deleteProduct(c *gin.Context) {
	id := parseID(c.Param("id"))

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)

			c.JSON(http.StatusOK, Response{
				Success: true,
				Message: "Produk berhasil dihapus",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: "Produk tidak ditemukan",
	})
}

// Helper function
func parseID(idStr string) int {
	var id int
	for _, char := range idStr {
		id = id*10 + int(char-'0')
	}
	return id
}

// @title Product API
// @version 1.0
// @description REST API untuk manajemen produk dengan dokumentasi Swagger
func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/products", getAllProducts)
		api.GET("/products/:id", getProductByID)
		api.POST("/products", createProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	println("===========================================")
	println("   SWAGGER API DOCUMENTATION")
	println("===========================================")
	println()
	println("📡 Server: http://localhost:8080")
	println("📚 Swagger UI: http://localhost:8080/swagger/index.html")
	println()
	println("Endpoints:")
	println("  GET    /api/v1/products")
	println("  GET    /api/v1/products/:id")
	println("  POST   /api/v1/products")
	println("  PUT    /api/v1/products/:id")
	println("  DELETE /api/v1/products/:id")
	println()
	println("Tekan Ctrl+C untuk menghentikan server")
	println()

	router.Run(":8080")
}
