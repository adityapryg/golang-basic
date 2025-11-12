package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Struct for Mahasiswa
type Mahasiswa struct {
	NIM    string
	Nama   string
	Jurusan string
	IPK    float64
}

// Sentinel errors
var (
	ErrNotFound      = fmt.Errorf("mahasiswa tidak ditemukan")
	ErrAlreadyExists  = fmt.Errorf("mahasiswa sudah ada")
	ErrInvalidInput   = fmt.Errorf("input tidak valid")
)

// Custom ValidationError type
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Global slice to store Mahasiswa
var daftarMahasiswa []Mahasiswa

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
	
tampilkanMenu()
		s := scanner.Scan()
		if !s {
			fmt.Println("Terjadi kesalahan input")
			continue
		}
		option := strings.TrimSpace(scanner.Text())
		
		switch option {
		case "1":
			tambahMahasiswa(scanner)
		case "2":
			lihatSemuaMahasiswa()
		case "3":
			 cariMahasiswa(scanner)
		case "4":
			 updateIPK(scanner)
		case "5":
			 hapusMahasiswa(scanner)
		case "6":
			 tampilkanStatistik()
		case "7":
			 fmt.Println("Keluar dari program...")
			 return
		default:
			 fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func tampilkanMenu() {
	fmt.Println("=== Menu Mahasiswa ===")
	fmt.Println("1. Tambah Mahasiswa")
	fmt.Println("2. Lihat Semua Mahasiswa")
	fmt.Println("3. Cari Mahasiswa")
	fmt.Println("4. Update IPK")
	fmt.Println("5. Hapus Mahasiswa")
	fmt.Println("6. Tampilkan Statistik")
	fmt.Println("7. Keluar")
	fmt.Print("Silakan pilih opsi: ")
}

func tambahMahasiswa(scanner *bufio.Scanner) {
	var m Mahasiswa
	fmt.Print("Masukkan NIM: ")
	scanner.Scan()
	m.NIM = strings.TrimSpace(scanner.Text())

	fmt.Print("Masukkan Nama: ")
	scanner.Scan()
	m.Nama = strings.TrimSpace(scanner.Text())

	fmt.Print("Masukkan Jurusan: ")
	scanner.Scan()
	m.Jurusan = strings.TrimSpace(scanner.Text())

	fmt.Print("Masukkan IPK: ")
	scanner.Scan()
	ipkStr := strings.TrimSpace(scanner.Text())
	var err error
	if m.IPK, err = validasiIPK(ipkStr); err != nil {
		fmt.Println(err)
		return
	}
	
	if err := validasiMahasiswa(m.NIM, m.Nama, m.Jurusan, m.IPK); err != nil {
		fmt.Println(err)
		return
	}

	daftarMahasiswa = append(daftarMahasiswa, m)
	fmt.Println("✓ Mahasiswa berhasil ditambahkan")
}

func lihatSemuaMahasiswa() {
	fmt.Printf("%-10s %-20s %-20s %-10s\n", "NIM", "Nama", "Jurusan", "IPK")
	fmt.Println(strings.Repeat("-", 60))
	for _, m := range daftarMahasiswa {
		fmt.Printf("%-10s %-20s %-20s %-10.2f\n", m.NIM, m.Nama, m.Jurusan, m.IPK)
	}
}

func cariMahasiswa(scanner *bufio.Scanner) {
	fmt.Print("Masukkan NIM yang dicari: ")
	scanner.Scan()
	nim := strings.TrimSpace(scanner.Text())
	for _, m := range daftarMahasiswa {
		if m.NIM == nim {
			fmt.Printf("NIM: %s\nNama: %s\nJurusan: %s\nIPK: %.2f\n", m.NIM, m.Nama, m.Jurusan, m.IPK)
			return
		}
	}
	fmt.Println("❌ Mahasiswa tidak ditemukan")
}

func updateIPK(scanner *bufio.Scanner) {
	fmt.Print("Masukkan NIM mahasiswa yang akan diupdate: ")
	scanner.Scan()
	nim := strings.TrimSpace(scanner.Text())

	for i, m := range daftarMahasiswa {
		if m.NIM == nim {
			fmt.Print("Masukkan IPK baru: ")
			scanner.Scan()
			ipkStr := strings.TrimSpace(scanner.Text())
			var err error
			if m.IPK, err = validasiIPK(ipkStr); err != nil {
				fmt.Println(err)
				return
			}
			daftarMahasiswa[i].IPK = m.IPK
			fmt.Println("✓ IPK mahasiswa berhasil diupdate")
			return
		}
	}
	fmt.Println("❌ Mahasiswa tidak ditemukan")
}

func hapusMahasiswa(scanner *bufio.Scanner) {
	fmt.Print("Masukkan NIM mahasiswa yang akan dihapus: ")
	scanner.Scan()
	nim := strings.TrimSpace(scanner.Text())

	for i, m := range daftarMahasiswa {
		if m.NIM == nim {
			var konfirmasi string
			fmt.Print("Apakah Anda yakin ingin menghapus mahasiswa ini? (y/n): ")
			scanner.Scan()
			konfirmasi = strings.TrimSpace(scanner.Text())
			if konfirmasi == "y" {
				daftarMahasiswa = append(daftarMahasiswa[:i], daftarMahasiswa[i+1:]...)
				fmt.Println("✓ Mahasiswa berhasil dihapus")
				return
			}
			return
		}
	}
	fmt.Println("❌ Mahasiswa tidak ditemukan")
}

func tampilkanStatistik() {
	var totalIPK, highestIPK, lowestIPK float64
	var totalMahasiswa = len(daftarMahasiswa)

	if totalMahasiswa == 0 {
		fmt.Println("Tidak ada mahasiswa untuk ditampilkan")
		return
	}

	highestIPK = daftarMahasiswa[0].IPK
	lowestIPK = daftarMahasiswa[0].IPK

	for _, m := range daftarMahasiswa {
		totalIPK += m.IPK
		if m.IPK > highestIPK {
			highestIPK = m.IPK
		}
		if m.IPK < lowestIPK {
			lowestIPK = m.IPK
		}
	}

	avgIPK := totalIPK / float64(totalMahasiswa)
	fmt.Printf("Total Mahasiswa: %d\nRata-rata IPK: %.2f\nIPK Tertinggi: %.2f\nIPK Terendah: %.2f\n", totalMahasiswa, avgIPK, highestIPK, lowestIPK)
}

func cariByNIM(nim string) (*Mahasiswa, error) {
	for _, m := range daftarMahasiswa {
		if m.NIM == nim {
			return &m, nil
		}
	}
	return nil, ErrNotFound
}

func validasiMahasiswa(nim, nama, jurusan string, ipk float64) error {
	if nim == "" || nama == "" || jurusan == "" {
		return &ValidationError{Field: "Input", Message: "semua field harus diisi"}
	}
	return nil
}

func validasiIPK(ipkStr string) (float64, error) {
	ipk, err := strconv.ParseFloat(ipkStr, 64)
	if err != nil || ipk < 0 || ipk > 4 {
		return 0, ErrInvalidInput
	}
	return ipk, nil
}