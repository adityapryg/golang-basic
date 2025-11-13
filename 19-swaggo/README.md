# 19. Dokumentasi API dengan Swaggo

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Cara generate dokumentasi API otomatis dengan Swaggo
2. Swagger/OpenAPI annotations di Go
3. Swagger UI untuk testing API interaktif
4. Best practices dokumentasi API
5. Integrasi Swagger dengan Gin framework

## Penjelasan

### Apa itu Swagger?

Swagger (OpenAPI) adalah standar untuk mendokumentasikan REST API. Swaggo adalah tool untuk generate Swagger documentation dari comment annotations di kode Go.

**Keuntungan:**

- ✅ Dokumentasi otomatis dari kode
- ✅ Interactive API testing (Swagger UI)
- ✅ Client SDK generation
- ✅ Selalu sync dengan kode
- ✅ Standar industry (OpenAPI spec)

### Swagger Annotations

```go
// @Summary      Ambil semua produk
// @Description  Mendapatkan list semua produk yang tersedia
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=[]Product}
// @Router       /products [get]
func getAllProducts(c *gin.Context) {
    // handler code
}
```

### Annotation Types

| Annotation     | Deskripsi                             |
| -------------- | ------------------------------------- |
| `@Summary`     | Ringkasan singkat endpoint            |
| `@Description` | Deskripsi detail                      |
| `@Tags`        | Grouping endpoints                    |
| `@Accept`      | Content-Type yang diterima            |
| `@Produce`     | Content-Type response                 |
| `@Param`       | Parameter (path, query, body, header) |
| `@Success`     | Success response                      |
| `@Failure`     | Error response                        |
| `@Router`      | Path dan HTTP method                  |
| `@Security`    | Authentication scheme                 |

### API Metadata (di main)

```go
// @title           Product API
// @version         1.0
// @description     REST API untuk manajemen produk
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

## Prasyarat

### 1. Install Dependencies

```bash
cd 19-swaggo
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 2. Install Swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify instalasi:

```bash
swag --version
```

**Note:** Pastikan `$GOPATH/bin` atau `%USERPROFILE%\go\bin` (Windows) sudah ada di PATH.

### 2. Install Dependencies

```bash
cd 19-swaggo
go mod tidy
```

Dependencies yang akan diinstall:

- `github.com/gin-gonic/gin v1.9.1` - Gin web framework
- `github.com/swaggo/swag v1.16.2` - Swagger generator
- `github.com/swaggo/gin-swagger v1.6.0` - Gin Swagger middleware
- `github.com/swaggo/files v1.0.1` - Swagger UI files

## Cara Menjalankan

### 1. Generate Swagger Docs

```bash
swag init
```

Output:

```
2025/11/13 10:30:00 Generate swagger docs....
2025/11/13 10:30:00 Generate general API Info
2025/11/13 10:30:00 Parsing go files...
2025/11/13 10:30:01 create docs.go at docs/docs.go
2025/11/13 10:30:01 create swagger.json at docs/swagger.json
2025/11/13 10:30:01 create swagger.yaml at docs/swagger.yaml
```

Ini akan generate folder `docs/` dengan:

- `docs.go`
- `swagger.json`
- `swagger.yaml`

### 2. Run Server

```bash
go run main.go
```

### 3. Buka Swagger UI

Browser: http://localhost:8080/swagger/index.html

## Cara Menggunakan Swagger UI

### 1. Browse Endpoints

Swagger UI menampilkan semua endpoints dengan:

- HTTP method dan path
- Summary dan description
- Parameters
- Request body schema
- Response schema
- Example values

### 2. Try It Out

1. Click endpoint (e.g., `GET /api/v1/products`)
2. Click **"Try it out"**
3. Fill parameters (jika ada)
4. Click **"Execute"**
5. Lihat response

### 3. Authentication

Untuk protected endpoints:

1. Click **"Authorize"** button (🔒 icon)
2. Masukkan token: `Bearer your-jwt-token-here`
3. Click **"Authorize"**
4. Click **"Close"**

Sekarang semua request akan include token di header.

### 4. Model Schemas

Scroll ke bawah untuk melihat:

- **Models** - Struktur data (Product, Response, dll)
- **Example values** - Contoh JSON

## Swagger Annotations Guide

### Basic Endpoint

```go
// @Summary      Endpoint summary
// @Description  Detailed description
// @Tags         tag-name
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseType
// @Router       /path [method]
```

### With Path Parameter

```go
// @Param        id   path      int  true  "Product ID"
// @Router       /products/{id} [get]
```

### With Query Parameter

```go
// @Param        page      query     int     false  "Page number"
// @Param        limit     query     int     false  "Items per page"
// @Param        search    query     string  false  "Search keyword"
// @Router       /products [get]
```

### With Request Body

```go
// @Param        product  body      CreateProductRequest  true  "Product data"
// @Router       /products [post]
```

### With Authentication

```go
// @Security     BearerAuth
// @Router       /products [post]
```

### Multiple Response Codes

```go
// @Success      200  {object}  Response{data=Product}
// @Failure      400  {object}  Response
// @Failure      404  {object}  Response
// @Failure      500  {object}  Response
```

## Struct Documentation

```go
// Product model
// @Description Product information
type Product struct {
    ID          int       `json:"id" example:"1"`
    Name        string    `json:"name" example:"Laptop Gaming"`
    Description string    `json:"description" example:"High-end laptop"`
    Price       float64   `json:"price" example:"15000000"`
    Stock       int       `json:"stock" example:"10"`
}
```

**Struct tags:**

- `json:"field_name"` - JSON field name
- `example:"value"` - Example value di Swagger

## Testing API via Swagger UI

### GET Request

1. Expand `GET /api/v1/products`
2. Click "Try it out"
3. Click "Execute"
4. Lihat Response

### POST Request

1. Expand `POST /api/v1/products`
2. Click "Try it out"
3. Edit Request body:
   ```json
   {
     "name": "New Product",
     "description": "Product description",
     "price": 100000,
     "stock": 20
   }
   ```
4. Click "Execute"

### PUT/DELETE Request

Similar dengan POST, tapi perlu parameter ID di path.

## Update Dokumentasi

Setiap kali ubah annotations:

```bash
swag init
```

Kemudian restart server.

## Best Practices

### 1. Always Document

Dokumentasikan **semua** public endpoints:

```go
// ✅ Good
// @Summary Ambil produk
// @Router /products [get]

// ❌ Bad - no documentation
func getAllProducts(c *gin.Context) {}
```

### 2. Use Descriptive Summaries

```go
// ✅ Good
// @Summary Ambil semua produk dengan pagination dan filter

// ❌ Bad
// @Summary Get products
```

### 3. Group with Tags

```go
// @Tags products
// @Tags users
// @Tags auth
```

### 4. Include Examples

```go
type Product struct {
    Price float64 `json:"price" example:"15000000"`
    Stock int     `json:"stock" example:"10"`
}
```

### 5. Document Error Responses

```go
// @Failure 400 {object} Response "Validation error"
// @Failure 404 {object} Response "Not found"
// @Failure 500 {object} Response "Internal server error"
```

### 6. Security Definitions

Untuk JWT:

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
```

### 7. Versioning

Include version di BasePath:

```go
// @BasePath /api/v1
```

## Swagger vs Postman

| Feature        | Swagger              | Postman         |
| -------------- | -------------------- | --------------- |
| Auto-generate  | ✅ Dari kode         | ❌ Manual       |
| Interactive UI | ✅ Ya                | ✅ Ya           |
| Sharing        | ✅ URL (self-hosted) | ✅ Cloud/export |
| Client gen     | ✅ Ya                | ❌ No           |
| Testing        | ⚠️ Basic             | ✅ Advanced     |
| Collaboration  | ⚠️ Limited           | ✅ Teams        |

**Recommendation:** Gunakan keduanya!

- Swagger - Untuk dokumentasi dan quick testing
- Postman - Untuk advanced testing dan automation

## Troubleshooting

### Error: "swag: command not found"

Install swag CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Pastikan `$GOPATH/bin` ada di PATH.

### Error: "docs not found"

Jalankan `swag init` untuk generate docs.

### Swagger UI tidak muncul

Cek:

1. Import `_ "package/docs"`
2. Route `/swagger/*any` registered
3. Docs folder exists

### Changes tidak muncul

Regenerate docs:

```bash
swag init
# Restart server
go run main.go
```

## Advanced Features

### Custom Response Types

```go
// @Success 200 {object} Response{data=[]Product} "List of products"
```

### Enum Values

```go
// @Param status query string false "Status" Enums(active, inactive, deleted)
```

### File Upload

```go
// @Param file formData file true "File to upload"
```

### Multiple Tags

```go
// @Tags products,inventory
```

## Export Swagger Spec

Swagger spec tersimpan di `docs/swagger.json` dan `docs/swagger.yaml`.

Gunakan untuk:

- Client SDK generation (openapi-generator)
- Import ke Postman
- API Gateway configuration
- Contract testing

## Referensi

- [Swaggo GitHub](https://github.com/swaggo/swag)
- [Swagger Documentation](https://swagger.io/docs/)
- [OpenAPI Specification](https://spec.openapis.org/oas/latest.html)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
