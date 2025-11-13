# 15. Setup Database dengan GORM

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Cara koneksi ke PostgreSQL menggunakan GORM
2. Definisi model dan mapping ke tabel database
3. Auto Migration untuk membuat/update schema
4. Relasi antar tabel (Foreign Key)
5. Query dasar GORM (CRUD, Where, Preload, Aggregate)
6. Connection pooling dan konfigurasi database

## Penjelasan

### Apa itu GORM?

GORM adalah ORM (Object-Relational Mapping) library untuk Go yang:

- Mapping struct Go ke tabel database
- Auto migration schema
- Mendukung PostgreSQL, MySQL, SQLite, SQL Server
- Query builder yang powerful
- Hooks, transactions, dan banyak fitur lainnya

### Model Definition

```go
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    Email     string    `gorm:"size:100;unique;not null" json:"email"`
    Age       int       `gorm:"not null" json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**GORM Tags:**

- `primaryKey` - Primary key
- `size:100` - VARCHAR(100)
- `not null` - NOT NULL constraint
- `unique` - UNIQUE constraint
- `default:0` - Default value
- `type:text` - Custom column type

### Relasi Antar Tabel

```go
type Product struct {
    CategoryID uint     `gorm:"not null" json:"category_id"`
    Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
}
```

GORM mendukung:

- **Has One** - User has one Profile
- **Has Many** - Category has many Products
- **Belongs To** - Product belongs to Category
- **Many to Many** - Students and Courses

### Connection String PostgreSQL

```go
dsn := "host=localhost user=postgres password=postgres dbname=golang_demo port=5432 sslmode=disable"
```

Format: `host={host} user={user} password={password} dbname={dbname} port={port} sslmode={mode}`

## Prasyarat

### 1. Install PostgreSQL

**Windows:**

1. Download dari https://www.postgresql.org/download/windows/
2. Install dengan default settings
3. Password default: `postgres`

**macOS:**

```bash
brew install postgresql
brew services start postgresql
```

**Linux (Ubuntu):**

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 2. Buat Database

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE golang_demo;

# Keluar
\q
```

Atau gunakan GUI tool seperti pgAdmin, DBeaver, atau TablePlus.

### 3. Install Dependencies

```bash
cd 15-database-gorm
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

Dependencies yang akan diinstall:

- `gorm.io/gorm` - GORM ORM library
- `gorm.io/driver/postgres` - PostgreSQL driver untuk GORM

## Cara Menjalankan

```bash
cd 15-database-gorm
go run main.go
```

## Output yang Diharapkan

```
===========================================
   SETUP DATABASE DENGAN GORM
===========================================

🔌 Menghubungkan ke database PostgreSQL...
✅ Koneksi database berhasil!

📦 Melakukan migrasi database...
✅ Migrasi database berhasil!

🌱 Menambahkan data seed...
✅ Data seed berhasil ditambahkan!

===========================================
   DEMONSTRASI QUERY GORM
===========================================

1️⃣  Semua Users:
   - John Doe (john@example.com), Umur: 30
   - Jane Smith (jane@example.com), Umur: 25
   - Bob Johnson (bob@example.com), Umur: 35

2️⃣  User dengan ID 1:
   - John Doe (john@example.com)

3️⃣  User dengan umur >= 30:
   - John Doe, Umur: 30
   - Bob Johnson, Umur: 35

4️⃣  Jumlah Users:
   - Total: 3 users

5️⃣  Products dengan Category:
   - Laptop Gaming (Rp 15000000) - Kategori: Elektronik
   - Mouse Wireless (Rp 150000) - Kategori: Elektronik

6️⃣  Update stock produk:
   - Stock Laptop Gaming: 10 → 15

7️⃣  Total nilai inventory:
   - Total nilai: Rp 225000000

===========================================
   PROGRAM SELESAI
===========================================
```

## Query GORM yang Didemonstrasikan

### 1. Find All

```go
var users []User
db.Find(&users)
```

### 2. Find by ID

```go
var user User
db.First(&user, 1)  // WHERE id = 1
```

### 3. Find with Condition

```go
var users []User
db.Where("age >= ?", 30).Find(&users)
```

### 4. Count

```go
var count int64
db.Model(&User{}).Count(&count)
```

### 5. Preload (Eager Loading)

```go
var products []Product
db.Preload("Category").Find(&products)
```

### 6. Update

```go
db.Model(&product).Update("stock", newStock)
```

### 7. Aggregate

```go
var total float64
db.Model(&Product{}).Select("SUM(price * stock)").Scan(&total)
```

## Verifikasi Database

### Menggunakan psql (CLI)

```bash
psql -U postgres -d golang_demo
```

**Commands:**

```sql
-- Lihat semua tabel
\dt

-- Query data
SELECT * FROM users;
SELECT * FROM categories;
SELECT * FROM products;

-- Cek struktur tabel
\d users
```

### Menggunakan GUI Tools

Recommended tools:

- **pgAdmin** - Official PostgreSQL GUI
- **DBeaver** - Universal database tool
- **TablePlus** - Modern database GUI (Mac/Windows)
- **DataGrip** - JetBrains IDE

## Connection Pooling

```go
sqlDB, _ := db.DB()

// Maximum idle connections
sqlDB.SetMaxIdleConns(10)

// Maximum open connections
sqlDB.SetMaxOpenConns(100)

// Maximum connection lifetime
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Penjelasan:**

- **MaxIdleConns** - Koneksi idle yang disimpan di pool
- **MaxOpenConns** - Maksimal koneksi yang bisa dibuka
- **ConnMaxLifetime** - Durasi maksimal koneksi sebelum di-recycle

## Best Practices

### 1. Gunakan Environment Variables

Jangan hardcode credentials:

```go
dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
)
```

### 2. Handle Errors

```go
result := db.Create(&user)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

### 3. Use Transactions

```go
tx := db.Begin()
// ... operations
if err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}
```

### 4. Index untuk Performance

```go
type User struct {
    Email string `gorm:"size:100;unique;index" json:"email"`
}
```

## Troubleshooting

### Error: "dial tcp: connect: connection refused"

PostgreSQL tidak running. Start service:

```bash
# Windows (Services.msc)
# macOS
brew services start postgresql
# Linux
sudo systemctl start postgresql
```

### Error: "password authentication failed"

Password salah atau user tidak ada. Reset password:

```bash
sudo -u postgres psql
ALTER USER postgres PASSWORD 'new_password';
```

### Error: "database does not exist"

Buat database terlebih dahulu:

```sql
CREATE DATABASE golang_demo;
```

## Referensi

- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Tutorial](https://www.postgresql.org/docs/current/tutorial.html)
- [GORM Guides](https://gorm.io/docs/connecting_to_the_database.html)
