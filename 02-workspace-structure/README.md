# Demonstrasi Go Workspace dan Package

## Tujuan Pembelajaran

- Memahami struktur workspace Go
- Mengerti konsep Go Modules
- Dapat membuat dan menggunakan package sendiri
- Memahami cara import package lokal dan eksternal

## Penjelasan

### Go Modules

Go Modules adalah sistem manajemen dependensi resmi di Go. Setiap project Go modern menggunakan Go Modules untuk mengelola package dan versinya.

### Struktur Folder

```
02-workspace-structure/
├── go.mod              # File konfigurasi module
├── main.go             # Entry point aplikasi
├── utils/              # Package lokal
│   └── greetings.go    # File dalam package utils
└── README.md           # Dokumentasi
```

## Konsep Penting

### 1. go.mod

File `go.mod` adalah jantung dari Go Module. Berisi:
- Nama module
- Versi Go yang digunakan
- Daftar dependensi eksternal

```go
module github.com/adityapryg/golang-demo/02-workspace-structure

go 1.22
```

### 2. Package

- Package adalah cara Go mengorganisir kode
- Setiap folder adalah satu package
- File dalam satu folder harus memiliki nama package yang sama
- Package `main` adalah package khusus untuk executable

### 3. Import

```go
// Import package standar
import "fmt"

// Import package lokal
import "github.com/adityapryg/golang-demo/02-workspace-structure/utils"

// Import multiple packages
import (
    "fmt"
    "time"
)
```

## Cara Menjalankan

```bash
# 1. Masuk ke folder
cd 02-workspace-structure

# 2. Jalankan program
go run main.go

# 3. Build executable
go build -o app main.go
./app
```

## Output yang Diharapkan

```
===========================================
   DEMONSTRASI GO WORKSPACE
===========================================

Menggunakan Package Lokal:
Selamat Pagi, Aditya!

Informasi Module:
- Module Path: github.com/adityapryg/golang-demo/02-workspace-structure
- Go Version: go1.22

Contoh Fungsi Exported vs Unexported:
✓ SayHello() - Exported (huruf kapital)
✗ sayGoodbye() - Unexported (huruf kecil)
```

## Exported vs Unexported

### Exported (Public)
- Dimulai dengan huruf KAPITAL
- Bisa diakses dari package lain
- Contoh: `SayHello()`, `UserName`, `MaxValue`

### Unexported (Private)
- Dimulai dengan huruf kecil
- Hanya bisa diakses dalam package yang sama
- Contoh: `sayGoodbye()`, `userName`, `maxValue`

## Perintah Go Module

```bash
# Inisialisasi module baru
go mod init <module-name>

# Download dependensi
go mod download

# Bersihkan dependensi yang tidak terpakai
go mod tidy

# Lihat dependensi
go list -m all
```

## Latihan

1. Buat package baru bernama `calculator` dengan fungsi `Add()`, `Subtract()`
2. Buat package `models` dengan struct `User` dan method-nya
3. Tambahkan dependensi eksternal dan gunakan
4. Ekspor beberapa fungsi dan buat beberapa unexported

## GitHub Copilot Prompt

```
Create a Go workspace demonstration that:
- Initializes a Go module
- Creates a local package named 'utils' with exported and unexported functions
- Demonstrates importing and using the local package in main
- Shows the difference between exported and unexported identifiers
- Includes proper package documentation
```