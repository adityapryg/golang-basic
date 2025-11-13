# 16. CRUD dengan Database + GORM

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Integrasi Gin dengan GORM untuk REST API
2. CRUD operations dengan database PostgreSQL
3. Relasi antar tabel (Category - Product)
4. Validasi foreign key
5. Preload untuk mengatasi N+1 query problem
6. Error handling yang proper

## Penjelasan

### Arsitektur Aplikasi

```
Client (curl/Postman)
    ↓
Gin Router & Handlers
    ↓
GORM Models & Queries
    ↓
PostgreSQL Database
```

### Model dengan Relasi

```go
type Category struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"unique;not null"`
}

type Product struct {
    ID         uint
    Name       string
    CategoryID uint     `gorm:"not null"`
    Category   Category `gorm:"foreignKey:CategoryID"`
}
```

Product **belongs to** Category (many-to-one relationship).

### Preload untuk Eager Loading

Tanpa Preload (N+1 problem):

```go
db.Find(&products)  // 1 query
for _, p := range products {
    db.First(&p.Category, p.CategoryID)  // N queries
}
```

Dengan Preload:

```go
db.Preload("Category").Find(&products)  // 2 queries total
```

### CRUD Operations

| Operation | Method | Endpoint             | GORM Method    |
| --------- | ------ | -------------------- | -------------- |
| Create    | POST   | /api/v1/products     | `db.Create()`  |
| Read All  | GET    | /api/v1/products     | `db.Find()`    |
| Read One  | GET    | /api/v1/products/:id | `db.First()`   |
| Update    | PUT    | /api/v1/products/:id | `db.Updates()` |
| Delete    | DELETE | /api/v1/products/:id | `db.Delete()`  |

## Prasyarat

### 1. Install Dependencies

```bash
cd 16-crud-with-db
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### 2. PostgreSQL dan Database

Pastikan sudah terinstall dan database sudah dibuat (lihat lesson 15):

```bash
# Cek PostgreSQL running
psql -U postgres -d golang_demo -c "SELECT 1"
```

### 2. Install Dependencies

```bash
cd 16-crud-with-db
go mod tidy
```

Dependencies yang akan diinstall:

- `github.com/gin-gonic/gin v1.9.1` - Gin web framework
- `gorm.io/driver/postgres v1.5.4` - PostgreSQL driver untuk GORM
- `gorm.io/gorm v1.25.5` - GORM ORM library

**Note:** Pastikan sudah mempelajari lesson 15 (database-gorm) terlebih dahulu.

## Cara Menjalankan

```bash
cd 16-crud-with-db
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Cara Testing

### 1. Lihat Semua Kategori

```bash
curl http://localhost:8080/api/v1/categories
```

**Response:**

```json
{
  "success": true,
  "message": "Data kategori berhasil diambil",
  "data": {
    "total": 3,
    "categories": [
      { "id": 1, "name": "Elektronik", "description": "Peralatan elektronik" },
      { "id": 2, "name": "Pakaian", "description": "Pakaian dan fashion" },
      { "id": 3, "name": "Makanan", "description": "Makanan dan minuman" }
    ]
  }
}
```

### 2. Buat Produk Baru

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Smartphone flagship dari Apple",
    "price": 20000000,
    "stock": 25,
    "category_id": 1
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Produk berhasil dibuat",
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "description": "Smartphone flagship dari Apple",
    "price": 20000000,
    "stock": 25,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Elektronik",
      "description": "Peralatan elektronik"
    },
    "created_at": "2025-11-13T10:30:00Z",
    "updated_at": "2025-11-13T10:30:00Z"
  }
}
```

### 3. Lihat Semua Produk

```bash
curl http://localhost:8080/api/v1/products
```

**Response:**

```json
{
  "success": true,
  "message": "Data produk berhasil diambil",
  "data": {
    "total": 1,
    "products": [
      {
        "id": 1,
        "name": "iPhone 15 Pro",
        "price": 20000000,
        "stock": 25,
        "category_id": 1,
        "category": {
          "id": 1,
          "name": "Elektronik"
        }
      }
    ]
  }
}
```

### 4. Lihat Produk by ID

```bash
curl http://localhost:8080/api/v1/products/1
```

### 5. Update Produk

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro Max",
    "description": "Versi terbesar dari iPhone 15",
    "price": 22000000,
    "stock": 20,
    "category_id": 1
  }'
```

### 6. Hapus Produk

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

### 7. Error - Kategori Tidak Ada

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "price": 1000,
    "stock": 10,
    "category_id": 999
  }'
```

**Response:**

```json
{
  "success": false,
  "message": "Kategori tidak ditemukan"
}
```

## Verifikasi di Database

```bash
psql -U postgres -d golang_demo
```

```sql
-- Lihat tabel
\dt

-- Query products dengan join
SELECT p.id, p.name, p.price, c.name as category
FROM products p
JOIN categories c ON p.category_id = c.id;

-- Count products per category
SELECT c.name, COUNT(p.id) as total_products
FROM categories c
LEFT JOIN products p ON c.id = p.category_id
GROUP BY c.name;
```

## Perbedaan dengan Lesson 14 (No DB)

| Aspek          | Lesson 14 (In-Memory) | Lesson 16 (Database)  |
| -------------- | --------------------- | --------------------- |
| Storage        | Map di memory         | PostgreSQL            |
| Persistence    | Hilang saat restart   | Permanen              |
| Concurrency    | sync.RWMutex          | Database lock         |
| Relasi         | Manual reference      | Foreign key & Preload |
| Query kompleks | Loop manual           | SQL joins             |
| Scalability    | Single instance       | Multiple instances OK |

## Best Practices

### 1. Validasi Foreign Key

```go
var category Category
if err := db.First(&category, product.CategoryID).Error; err != nil {
    return "Kategori tidak ditemukan"
}
```

### 2. Preload Relasi

```go
db.Preload("Category").Find(&products)  // Eager loading
```

### 3. Handle Not Found

```go
if err := db.First(&product, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return 404, "Produk tidak ditemukan"
    }
    return 500, "Database error"
}
```

### 4. Use Transactions

Untuk operasi yang melibatkan multiple tables:

```go
tx := db.Begin()
tx.Create(&category)
tx.Create(&product)
tx.Commit()
```

## Troubleshooting

### Error: "relation does not exist"

Auto migrate belum jalan. Restart aplikasi atau jalankan:

```go
db.AutoMigrate(&Category{}, &Product{})
```

### Error: "violates foreign key constraint"

Category yang direferensikan tidak ada. Cek dulu category exists.

### Data tidak muncul di response

Pastikan menggunakan `Preload()` untuk load relasi:

```go
db.Preload("Category").Find(&products)
```

## Referensi

- [GORM Associations](https://gorm.io/docs/belongs_to.html)
- [GORM Preloading](https://gorm.io/docs/preload.html)
- [Gin + GORM Tutorial](https://gorm.io/docs/connecting_to_the_database.html)
