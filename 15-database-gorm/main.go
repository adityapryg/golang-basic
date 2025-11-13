package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User adalah model untuk tabel users
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;unique;not null" json:"email"`
	Age       int       `gorm:"not null" json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category adalah model untuk tabel categories
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Product adalah model untuk tabel products dengan relasi ke Category
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// connectDB membuat koneksi ke PostgreSQL
func connectDB() (*gorm.DB, error) {
	// Connection string
	dsn := "host=localhost user=postgres password=postgres dbname=golang_demo port=5432 sslmode=disable"

	// Konfigurasi GORM dengan logger
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Buat koneksi
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sql.DB: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// migrateDatabase membuat/update schema database
func migrateDatabase(db *gorm.DB) error {
	fmt.Println("\n📦 Melakukan migrasi database...")

	// Auto migrate akan membuat/update tabel
	err := db.AutoMigrate(&User{}, &Category{}, &Product{})
	if err != nil {
		return fmt.Errorf("migrasi gagal: %w", err)
	}

	fmt.Println("✅ Migrasi database berhasil!")
	return nil
}

// seedData menambahkan data awal ke database
func seedData(db *gorm.DB) error {
	fmt.Println("\n🌱 Menambahkan data seed...")

	// Seed users
	users := []User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}

	for _, user := range users {
		// FirstOrCreate: buat jika belum ada
		result := db.Where("email = ?", user.Email).FirstOrCreate(&user)
		if result.Error != nil {
			return result.Error
		}
	}

	// Seed categories
	categories := []Category{
		{Name: "Elektronik", Description: "Peralatan elektronik"},
		{Name: "Pakaian", Description: "Pakaian dan fashion"},
		{Name: "Makanan", Description: "Makanan dan minuman"},
	}

	for _, category := range categories {
		result := db.Where("name = ?", category.Name).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}
	}

	// Seed products
	var elektronik Category
	db.Where("name = ?", "Elektronik").First(&elektronik)

	products := []Product{
		{
			Name:        "Laptop Gaming",
			Description: "Laptop gaming dengan spesifikasi tinggi",
			Price:       15000000,
			Stock:       10,
			CategoryID:  elektronik.ID,
		},
		{
			Name:        "Mouse Wireless",
			Description: "Mouse wireless ergonomis",
			Price:       150000,
			Stock:       50,
			CategoryID:  elektronik.ID,
		},
	}

	for _, product := range products {
		result := db.Where("name = ?", product.Name).FirstOrCreate(&product)
		if result.Error != nil {
			return result.Error
		}
	}

	fmt.Println("✅ Data seed berhasil ditambahkan!")
	return nil
}

// demonstrateQueries mendemonstrasikan berbagai query GORM
func demonstrateQueries(db *gorm.DB) {
	fmt.Println("\n" + "===========================================")
	fmt.Println("   DEMONSTRASI QUERY GORM")
	fmt.Println("===========================================")

	// 1. Find All Users
	fmt.Println("\n1️⃣  Semua Users:")
	var users []User
	db.Find(&users)
	for _, user := range users {
		fmt.Printf("   - %s (%s), Umur: %d\n", user.Name, user.Email, user.Age)
	}

	// 2. Find User by ID
	fmt.Println("\n2️⃣  User dengan ID 1:")
	var user User
	db.First(&user, 1)
	fmt.Printf("   - %s (%s)\n", user.Name, user.Email)

	// 3. Find User by condition
	fmt.Println("\n3️⃣  User dengan umur >= 30:")
	var oldUsers []User
	db.Where("age >= ?", 30).Find(&oldUsers)
	for _, u := range oldUsers {
		fmt.Printf("   - %s, Umur: %d\n", u.Name, u.Age)
	}

	// 4. Count records
	fmt.Println("\n4️⃣  Jumlah Users:")
	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Printf("   - Total: %d users\n", count)

	// 5. Find Products with Category (Join)
	fmt.Println("\n5️⃣  Products dengan Category:")
	var products []Product
	db.Preload("Category").Find(&products)
	for _, p := range products {
		fmt.Printf("   - %s (Rp %.0f) - Kategori: %s\n", p.Name, p.Price, p.Category.Name)
	}

	// 6. Update record
	fmt.Println("\n6️⃣  Update stock produk:")
	var laptop Product
	db.Where("name = ?", "Laptop Gaming").First(&laptop)
	oldStock := laptop.Stock
	db.Model(&laptop).Update("stock", laptop.Stock+5)
	fmt.Printf("   - Stock %s: %d → %d\n", laptop.Name, oldStock, laptop.Stock)

	// 7. Aggregate function
	fmt.Println("\n7️⃣  Total nilai inventory:")
	var totalValue float64
	db.Model(&Product{}).Select("SUM(price * stock)").Scan(&totalValue)
	fmt.Printf("   - Total nilai: Rp %.0f\n", totalValue)
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   SETUP DATABASE DENGAN GORM")
	fmt.Println("===========================================")

	// 1. Connect to database
	fmt.Println("\n🔌 Menghubungkan ke database PostgreSQL...")
	db, err := connectDB()
	if err != nil {
		log.Fatal("❌ Koneksi database gagal:", err)
	}
	fmt.Println("✅ Koneksi database berhasil!")

	// 2. Migrate database
	if err := migrateDatabase(db); err != nil {
		log.Fatal("❌ Migrasi database gagal:", err)
	}

	// 3. Seed data
	if err := seedData(db); err != nil {
		log.Fatal("❌ Seed data gagal:", err)
	}

	// 4. Demonstrate queries
	demonstrateQueries(db)

	fmt.Println("\n===========================================")
	fmt.Println("   PROGRAM SELESAI")
	fmt.Println("===========================================")
	fmt.Println("\nCek database Anda, tabel dan data sudah dibuat!")
	fmt.Println("Gunakan: psql -U postgres -d golang_demo")
	fmt.Println("Atau tool GUI seperti pgAdmin, DBeaver, TablePlus")
}
