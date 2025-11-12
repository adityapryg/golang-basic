package main

import "fmt"

// Konstanta di package level
const (
	AppName    = "Golang Demo"
	AppVersion = "1.0.0"
	MaxUsers   = 1000
)

func main() {
	// Header
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI SINTAKS DASAR GO")
	fmt.Println("===========================================")
	fmt.Println()

	// 1. DEKLARASI VARIABEL
	demonstrasiVariabel()

	// 2. TIPE DATA NUMERIK
	demonstrasiTipeData()

	// 3. KONSTANTA
	demonstrasiKonstanta()

	// 4. ZERO VALUES
	demonstrasiZeroValues()

	// 5. KONVERSI TIPE DATA
	demonstrasiKonversi()
}

func demonstrasiVariabel() {
	fmt.Println("1. DEKLARASI VARIABEL")
	fmt.Println("---------------------")

	// Cara 1: Menggunakan var dengan tipe eksplisit
	var nama string = "Aditya Prayoga"
	var umur int = 25
	var tinggi float64 = 175.5

	fmt.Println("Menggunakan var:")
	fmt.Printf("- Nama: %s\n", nama)
	fmt.Printf("- Umur: %d tahun\n", umur)
	fmt.Printf("- Tinggi: %.2f cm\n", tinggi)
	fmt.Println()

	// Cara 2: Type inference (Go menebak tipe)
	var negara = "Indonesia"

	// Cara 3: Short declaration (hanya di dalam fungsi)
	kota := "Jakarta"
	isActive := true

	fmt.Println("Menggunakan := (short declaration):")
	fmt.Printf("- Kota: %s\n", kota)
	fmt.Printf("- Negara: %s\n", negara)
	fmt.Printf("- Status Aktif: %t\n", isActive)
	fmt.Println()
}

func demonstrasiTipeData() {
	fmt.Println("2. TIPE DATA NUMERIK")
	fmt.Println("-------------------")

	// Integer types
	var i8 int8 = 127
	var i16 int16 = 32767
	var i32 int32 = 2147483647

	fmt.Println("Integer:")
	fmt.Printf("- int8: -128 hingga %d\n", i8)
	fmt.Printf("- int16: -32768 hingga %d\n", i16)
	fmt.Printf("- int32: -2147483648 hingga %d\n", i32)
	fmt.Println()

	// Float types
	var gaji float64 = 15000000.50
	var nilai float32 = 89.75

	fmt.Println("Float:")
	fmt.Printf("- Gaji: Rp %.2f\n", gaji)
	fmt.Printf("- Nilai: %.2f\n", nilai)
	fmt.Println()

	// Boolean
	var lulus bool = true
	fmt.Printf("Status Kelulusan: %t\n", lulus)
	fmt.Println()
}

func demonstrasiKonstanta() {
	fmt.Println("3. KONSTANTA")
	fmt.Println("-----------")

	// Konstanta lokal
	const Pi = 3.14159
	const JumlahHari = 7

	fmt.Printf("- Nama Aplikasi: %s\n", AppName)
	fmt.Printf("- Versi: %s\n", AppVersion)
	fmt.Printf("- Maksimal User: %d\n", MaxUsers)
	fmt.Printf("- Pi: %.5f\n", Pi)
	fmt.Printf("- Jumlah Hari dalam Seminggu: %d\n", JumlahHari)
	fmt.Println()

	// Contoh penggunaan konstanta
	radius := 7.0
	luas := Pi * radius * radius
	fmt.Printf("Luas lingkaran dengan radius %.1f cm: %.2f cm²\n", radius, luas)
	fmt.Println()
}

func demonstrasiZeroValues() {
	fmt.Println("4. ZERO VALUES")
	fmt.Println("-------------")
	fmt.Println("Variabel tanpa nilai awal mendapat 'zero value':")

	var x int
	var y string
	var z bool
	var f float64

	fmt.Printf("var x int = %d\n", x)
	fmt.Printf("var y string = \"%s\"\n", y)
	fmt.Printf("var z bool = %t\n", z)
	fmt.Printf("var f float64 = %.1f\n", f)
	fmt.Println()
}

func demonstrasiKonversi() {
	fmt.Println("5. KONVERSI TIPE DATA")
	fmt.Println("--------------------")

	// Integer ke Float
	angka := 42
	desimal := float64(angka)
	fmt.Printf("int ke float64: %d → %.2f\n", angka, desimal)

	// Float ke Integer (kehilangan nilai desimal)
	pi := 3.14159
	piInt := int(pi)
	fmt.Printf("float64 ke int: %.5f → %d\n", pi, piInt)

	// String ke byte slice
	teks := "Hello"
	bytes := []byte(teks)
	fmt.Printf("string ke []byte: %s → %v\n", teks, bytes)
	fmt.Println()

	// Demonstrasi type safety
	var a int = 10
	var b float64 = 20.5
	// var c = a + b  // ERROR: tidak bisa langsung
	var c = float64(a) + b // Harus konversi dulu
	fmt.Printf("Type Safety: int(%d) + float64(%.1f) = %.1f (setelah konversi)\n", a, b, c)
	fmt.Println()
}