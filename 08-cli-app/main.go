package main

import (
	"bufio"
	"fmt"
	"os"
)

// Struct Mahasiswa represents a student.
type Mahasiswa struct {
	NIM     string
	Nama    string
	Jurusan string
	IPK     float64
}

// Custom error types
var (
	ErrNotFound       = fmt.Errorf("mahasiswa tidak ditemukan")
	ErrAlreadyExists   = fmt.Errorf("mahasiswa sudah ada")
	ErrInvalidInput   = fmt.Errorf("input tidak valid")
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Global slice to store mahasiswa data.
var mahasiswas []Mahasiswa

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		showMenu()
		scanner.Scan()
		option := scanner.Text()
		switch option {
		case "1":
			addMahasiswa(scanner)
		case "2":
			viewMahasiswa()
		case "3":
			searchMahasiswa(scanner)
		case "4":
			updateIPK(scanner)
		case "5":
			hapusMahasiswa(scanner)
		case "6":
			showStatistics()
		case "7":
			fmt.Println("Keluar...")
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}

func showMenu() {
	fmt.Println("Menu:")
	fmt.Println("1. Tambah")
	fmt.Println("2. Lihat Semua")
	fmt.Println("3. Cari")
	fmt.Println("4. Update IPK")
	fmt.Println("5. Hapus")
	fmt.Println("6. Statistik")
	fmt.Println("7. Keluar")
	fmt.Print("Pilih opsi: ")
}

func addMahasiswa(scanner *bufio.Scanner) {
	fmt.Print("Masukkan NIM: ")
	scanner.Scan()
	NIM := scanner.Text()

	if err := validateNIM(NIM); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Masukkan Nama: ")
	scanner.Scan()
	Nama := scanner.Text()

	if err := validateNama(Nama); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Masukkan Jurusan: ")
	scanner.Scan()
	Jurusan := scanner.Text()

	fmt.Print("Masukkan IPK: ")
	scanner.Scan()
	var IPK float64
	fmt.Sscan(scanner.Text(), &IPK)

	if err := validateIPK(IPK); err != nil {
		fmt.Println(err)
		return
	}

	mahasiswas = append(mahasiswas, Mahasiswa{NIM, Nama, Jurusan, IPK})
	fmt.Println("Mahasiswa ditambahkan.")
}

// Additional CRUD functions here... (viewMahasiswa, searchMahasiswa, etc.)

func validateNIM(NIM string) error {
	// Add logic to check for unique NIM
	return nil
}

func validateNama(Nama string) error {
	if len(Nama) < 3 {
		return &ValidationError{Field: "Nama", Message: "Nama minimal 3 karakter"}
	}
	return nil
}

func validateIPK(IPK float64) error {
	if IPK < 0.0 || IPK > 4.0 {
		return &ValidationError{Field: "IPK", Message: "IPK harus antara 0.0 dan 4.0"}
	}
	return nil
}

func showStatistics() {
	totalIPK := 0.0
	numMahasiswa := len(mahasiswas)
	if numMahasiswa == 0 {
		fmt.Println("Tidak ada mahasiswa.")
		return
	}

	for _, m := range mahasiswas {
		totalIPK += m.IPK
	}
	fmt.Printf("Total Mahasiswa: %d\n", numMahasiswa)
	fmt.Printf("Rata-rata IPK: %.2f\n", totalIPK/float64(numMahasiswa))
	// Additional statistics...
}

// Implement delete confirmation and remaining CRUD operations