package main

import (
	"errors"
	"fmt"
)

func main() {
	// Header
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI FUNGSI DAN STRUCT")
	fmt.Println("===========================================")
	fmt.Println()

	// 1. FUNGSI SEDERHANA
	demonstrasiFungsiSederhana()

	// 2. MULTIPLE RETURN VALUES
	demonstrasiMultipleReturns()

	// 3. STRUCT DAN METHODS
	demonstrasiStruct()

	// 4. CONSTRUCTOR PATTERN
	demonstrasiConstructor()
}

// ============================================
// 1. FUNGSI SEDERHANA
// ============================================

func greet(name string) string {
	return "Halo, " + name + "!"
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func demonstrasiFungsiSederhana() {
	fmt.Println("1. FUNGSI SEDERHANA")
	fmt.Println("------------------")

	result := greet("Aditya")
	fmt.Println("Greeting:", result)

	sum := add(10, 5)
	fmt.Printf("Penjumlahan: 10 + 5 = %d\n", sum)

	product := multiply(4, 3)
	fmt.Printf("Perkalian: 4 * 3 = %d\n", product)

	fmt.Println()
}

// ============================================
// 2. MULTIPLE RETURN VALUES
// ============================================

// Fungsi dengan multiple return values untuk error handling
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Named return values
func minMax(numbers []int) (min, max int) {
	if len(numbers) == 0 {
		return 0, 0
	}

	min = numbers[0]
	max = numbers[0]

	for _, n := range numbers {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return // naked return
}

func demonstrasiMultipleReturns() {
	fmt.Println("2. MULTIPLE RETURN VALUES")
	fmt.Println("------------------------")

	// Success case
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Hasil pembagian: 10.00 / 2.00 = %.2f\n", result)
	}

	// Error case
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error pembagian:", err)
	}

	fmt.Println()

	// Named returns
	numbers := []int{5, 2, 9, 1, 7, 3}
	min, max := minMax(numbers)
	fmt.Printf("Dari %v: min=%d, max=%d\n", numbers, min, max)

	fmt.Println()
}

// ============================================
// 3. STRUCT DAN METHODS
// ============================================

// Definisi struct Mahasiswa
type Mahasiswa struct {
	Nama    string
	NIM     string
	Jurusan string
	IPK     float64
}

// Value receiver (read-only) - tidak mengubah original
func (m Mahasiswa) PrintInfo() {
	fmt.Println("Data Mahasiswa:")
	fmt.Println("Nama:", m.Nama)
	fmt.Println("NIM:", m.NIM)
	fmt.Println("Jurusan:", m.Jurusan)
	fmt.Printf("IPK: %.2f\n", m.IPK)
}

// Value receiver untuk mendapatkan info
func (m Mahasiswa) GetFullInfo() string {
	return fmt.Sprintf("%s (NIM: %s) - %s", m.Nama, m.NIM, m.Jurusan)
}

// Pointer receiver (can modify) - bisa mengubah original
func (m *Mahasiswa) UpdateIPK(ipkBaru float64) {
	m.IPK = ipkBaru
}

// Method untuk cek kelulusan
func (m Mahasiswa) IsLulus() bool {
	return m.IPK >= 2.75
}

func demonstrasiStruct() {
	fmt.Println("3. STRUCT DAN METHODS")
	fmt.Println("--------------------")

	// Membuat instance struct
	mahasiswa := Mahasiswa{
		Nama:    "Budi Santoso",
		NIM:     "123456",
		Jurusan: "Teknik Informatika",
		IPK:     3.75,
	}

	// Memanggil method dengan value receiver
	mahasiswa.PrintInfo()
	fmt.Println()

	// Method yang mengembalikan string
	fmt.Println("Info lengkap:", mahasiswa.GetFullInfo())

	// Memanggil method dengan pointer receiver
	mahasiswa.UpdateIPK(3.85)
	fmt.Printf("Setelah update IPK: %.2f\n", mahasiswa.IPK)

	// Cek kelulusan
	if mahasiswa.IsLulus() {
		fmt.Println("Status: LULUS")
	} else {
		fmt.Println("Status: TIDAK LULUS")
	}

	fmt.Println()
}

// ============================================
// 4. CONSTRUCTOR PATTERN
// ============================================

// Constructor function
func NewMahasiswa(nama, nim, jurusan string, ipk float64) *Mahasiswa {
	return &Mahasiswa{
		Nama:    nama,
		NIM:     nim,
		Jurusan: jurusan,
		IPK:     ipk,
	}
}

func demonstrasiConstructor() {
	fmt.Println("4. CONSTRUCTOR PATTERN")
	fmt.Println("---------------------")

	// Menggunakan constructor
	mhs1 := NewMahasiswa("Siti Rahayu", "654321", "Sistem Informasi", 3.60)
	mhs2 := NewMahasiswa("Ahmad Fauzi", "789012", "Teknik Elektro", 3.40)

	fmt.Println("Mahasiswa 1:", mhs1.GetFullInfo())
	fmt.Printf("IPK: %.2f\n", mhs1.IPK)
	fmt.Println()

	fmt.Println("Mahasiswa 2:", mhs2.GetFullInfo())
	fmt.Printf("IPK: %.2f\n", mhs2.IPK)

	fmt.Println()
