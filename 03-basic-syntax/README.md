# Demonstrasi Sintaks Dasar Go

## Tujuan Pembelajaran

- Memahami cara deklarasi variabel di Go
- Mengerti berbagai tipe data built-in
- Dapat menggunakan konstanta dengan benar
- Memahami type inference dan zero values
- Mengetahui konversi tipe data

## Penjelasan

### Variabel di Go

Go menyediakan beberapa cara untuk mendeklarasikan variabel:

#### 1. Deklarasi dengan var
```go
var nama string = "Aditya"
var umur int = 25
var tinggi float64 = 175.5
```

#### 2. Type Inference (Go menebak tipe)
```go
var nama = "Aditya"  // Go tahu ini string
var umur = 25        // Go tahu ini int
```

#### 3. Short Declaration (hanya di dalam fungsi)
```go
nama := "Aditya"
umur := 25
```

### Tipe Data Built-in

#### Numerik
- **int**, int8, int16, int32, int64
- **uint**, uint8, uint16, uint32, uint64
- **float32**, float64
- **complex64**, complex128

#### String dan Boolean
- **string** - teks
- **bool** - true atau false

#### Lainnya
- **byte** - alias untuk uint8
- **rune** - alias untuk int32 (untuk karakter Unicode)

### Konstanta

Konstanta adalah nilai yang tidak bisa diubah:

```go
const Pi = 3.14159
const AppName = "MyApp"
const MaxUsers = 100
```

### Zero Values

Jika variabel dideklarasikan tanpa nilai awal, Go memberikan "zero value":
- Numerik: `0`
- Boolean: `false`
- String: `""` (string kosong)
- Pointer: `nil`

## Cara Menjalankan

```bash
cd 03-basic-syntax
go run main.go
```

## Output yang Diharapkan

```
===========================================
   DEMONSTRASI SINTAKS DASAR GO
===========================================

1. DEKLARASI VARIABEL
---------------------
Menggunakan var:
- Nama: Aditya Prayoga
- Umur: 25 tahun
- Tinggi: 175.50 cm

Menggunakan := (short declaration):
- Kota: Jakarta
- Negara: Indonesia

2. TIPE DATA NUMERIK
-------------------
Integer:
- int8: -128 hingga 127
- int16: -32768 hingga 32767
- int32: -2147483648 hingga 2147483647

Float:
- Gaji: Rp 15000000.50
- Nilai: 89.75

3. KONSTANTA
-----------
- Nama Aplikasi: Golang Demo
- Versi: 1.0.0
- Maksimal User: 1000

4. ZERO VALUES
-------------
var x int = 0
var y string = ""
var z bool = false

5. KONVERSI TIPE DATA
--------------------
int ke float64: 42 → 42.00
float64 ke int: 3.14 → 3
string ke []byte: Hello → [72 101 108 108 111]
```

## Konsep Penting

### 1. Type Safety
Go adalah bahasa strongly typed - tipe data harus eksplisit dan tidak bisa dicampur:
```go
var x int = 10
var y float64 = 20.5
// z := x + y  // ERROR! Tidak bisa tambah int dengan float64
z := float64(x) + y  // Harus konversi dulu
```

### 2. Short Declaration Rules
- Hanya bisa digunakan di dalam fungsi
- Minimal satu variabel baru harus ada
- Tipe ditentukan otomatis dari nilai

### 3. Multiple Declaration
```go
var a, b, c int = 1, 2, 3
x, y := "hello", 42
```

## Latihan

1. Buat program yang menghitung luas lingkaran (gunakan const untuk Pi)
2. Deklarasikan variabel untuk data mahasiswa (nim, nama, ipk)
3. Praktikkan konversi antara berbagai tipe numerik
4. Buat konstanta untuk konfigurasi aplikasi
5. Coba zero values untuk berbagai tipe data

## GitHub Copilot Prompt

```
Create a Go program demonstrating variables and data types that:
- Shows all three ways to declare variables (var, var with inference, :=)
- Demonstrates all numeric types (int, float, uint)
- Shows string and boolean usage
- Includes constants for configuration
- Demonstrates zero values for each type
- Shows type conversion between numeric types
- Has clear Indonesian comments explaining each concept
```

## Best Practices

1. **Gunakan := untuk variabel lokal** - lebih ringkas dan mudah dibaca
2. **Gunakan var untuk package-level** - lebih jelas dan eksplisit
3. **Nama variabel deskriptif** - `userName` lebih baik dari `u`
4. **Konstanta untuk nilai tetap** - lebih aman dan jelas
5. **Zero values** - manfaatkan untuk inisialisasi sederhana

## Troubleshooting

**Problem**: `undefined: variableName`
**Solution**: Pastikan variabel sudah dideklarasikan sebelum digunakan

**Problem**: `cannot use x (type int) as type float64`
**Solution**: Lakukan konversi tipe: `float64(x)`

**Problem**: `no new variables on left side of :=`
**Solution**: Gunakan `=` untuk assignment, atau deklarasikan variabel baru

**Problem**: `syntax error: unexpected :=`
**Solution**: `:=` hanya bisa digunakan di dalam fungsi, gunakan `var` di package level
