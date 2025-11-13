# 14. CRUD API dengan Gin Tanpa Database

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Cara membuat REST API dengan operasi CRUD lengkap (Create, Read, Update, Delete)
2. Penggunaan in-memory storage dengan map dan sync.RWMutex
3. Binding dan validasi JSON request
4. Format response API yang konsisten
5. HTTP status codes yang tepat untuk setiap operasi

## Penjelasan

### Operasi CRUD

API ini mengimplementasikan 5 endpoint untuk manajemen produk:

| Method | Endpoint             | Fungsi             | Status Code |
| ------ | -------------------- | ------------------ | ----------- |
| POST   | /api/v1/products     | Buat produk baru   | 201         |
| GET    | /api/v1/products     | Ambil semua produk | 200         |
| GET    | /api/v1/products/:id | Ambil produk by ID | 200/404     |
| PUT    | /api/v1/products/:id | Update produk      | 200/404     |
| DELETE | /api/v1/products/:id | Hapus produk       | 200/404     |

### Model Product

```go
type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name" binding:"required"`
    Description string    `json:"description"`
    Price       float64   `json:"price" binding:"required,gt=0"`
    Stock       int       `json:"stock" binding:"required,gte=0"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Validasi:**

- `Name` - Wajib diisi
- `Price` - Wajib diisi dan harus > 0
- `Stock` - Wajib diisi dan harus >= 0

### Format Response

Semua response menggunakan format standar:

```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

### Thread Safety dengan sync.RWMutex

Karena menggunakan in-memory map, kita perlu `sync.RWMutex` untuk mencegah race condition:

```go
var (
    products   = make(map[int]Product)
    productsMu sync.RWMutex
)

// Read operation
productsMu.RLock()
defer productsMu.RUnlock()

// Write operation
productsMu.Lock()
defer productsMu.Unlock()
```

## Prasyarat

### Install Dependencies

```bash
cd 14-crud-no-db
go get -u github.com/gin-gonic/gin
```

Dependencies yang akan diinstall:

- `github.com/gin-gonic/gin` - Gin web framework

## Cara Menjalankan

```bash
cd 14-crud-no-db
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Cara Testing dengan curl

### 1. Create - Buat Produk Baru

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Gaming",
    "description": "Laptop gaming dengan spesifikasi tinggi",
    "price": 15000000,
    "stock": 10
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Produk berhasil dibuat",
  "data": {
    "id": 1,
    "name": "Laptop Gaming",
    "description": "Laptop gaming dengan spesifikasi tinggi",
    "price": 15000000,
    "stock": 10,
    "created_at": "2025-11-13T10:30:00Z",
    "updated_at": "2025-11-13T10:30:00Z"
  }
}
```

### 2. Read All - Ambil Semua Produk

```bash
curl http://localhost:8080/api/v1/products
```

**Response:**

```json
{
  "success": true,
  "message": "Data produk berhasil diambil",
  "data": {
    "total": 2,
    "products": [
      {
        "id": 1,
        "name": "Laptop Gaming",
        "price": 15000000,
        "stock": 10
      },
      {
        "id": 2,
        "name": "Mouse Wireless",
        "price": 150000,
        "stock": 50
      }
    ]
  }
}
```

### 3. Read One - Ambil Produk by ID

```bash
curl http://localhost:8080/api/v1/products/1
```

**Response:**

```json
{
  "success": true,
  "message": "Produk ditemukan",
  "data": {
    "id": 1,
    "name": "Laptop Gaming",
    "description": "Laptop gaming dengan spesifikasi tinggi",
    "price": 15000000,
    "stock": 10,
    "created_at": "2025-11-13T10:30:00Z",
    "updated_at": "2025-11-13T10:30:00Z"
  }
}
```

### 4. Update - Update Produk

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Gaming ROG",
    "description": "ASUS ROG dengan RTX 4090",
    "price": 18000000,
    "stock": 8
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Produk berhasil diupdate",
  "data": {
    "id": 1,
    "name": "Laptop Gaming ROG",
    "description": "ASUS ROG dengan RTX 4090",
    "price": 18000000,
    "stock": 8,
    "created_at": "2025-11-13T10:30:00Z",
    "updated_at": "2025-11-13T10:35:00Z"
  }
}
```

### 5. Delete - Hapus Produk

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

**Response:**

```json
{
  "success": true,
  "message": "Produk berhasil dihapus"
}
```

### 6. Error Handling - Validasi Gagal

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "price": -1000,
    "stock": -5
  }'
```

**Response:**

```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Product.Price' Error:Field validation for 'Price' failed on the 'gt' tag"
}
```

## Testing dengan Postman

1. Buat collection baru: "CRUD Products"
2. Tambahkan 5 requests sesuai endpoints di atas
3. Set `Content-Type: application/json` di Headers
4. Untuk POST dan PUT, masukkan JSON body di Body → raw → JSON

## Konsep Penting

### 1. HTTP Status Codes

- **200 OK** - Operasi berhasil (GET, PUT, DELETE)
- **201 Created** - Resource berhasil dibuat (POST)
- **400 Bad Request** - Request tidak valid (validasi gagal)
- **404 Not Found** - Resource tidak ditemukan

### 2. RESTful Principles

- Gunakan HTTP methods sesuai fungsi (GET, POST, PUT, DELETE)
- URL merepresentasikan resource (`/products`, `/products/:id`)
- Status code yang konsisten

### 3. Data Binding & Validation

Gin otomatis bind dan validasi:

```go
if err := c.ShouldBindJSON(&product); err != nil {
    // Tangani error validasi
}
```

### 4. Thread Safety

Map tidak thread-safe, gunakan mutex:

```go
productsMu.Lock()   // Write lock
productsMu.RLock()  // Read lock (multiple readers OK)
```

## Keterbatasan In-Memory Storage

- Data hilang saat server restart
- Tidak scalable untuk production
- Race condition jika tidak hati-hati dengan mutex
- Tidak bisa distributed/multi-instance

**Solusi:** Gunakan database (akan dipelajari di lesson berikutnya)

## Referensi

- [Gin Binding](https://gin-gonic.com/docs/examples/binding-and-validation/)
- [Go sync.RWMutex](https://pkg.go.dev/sync#RWMutex)
- [REST API Best Practices](https://restfulapi.net/)
