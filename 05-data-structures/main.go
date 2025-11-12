package main

import "fmt"

func main() {
	// Header
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI STRUKTUR DATA")
	fmt.Println("===========================================")
	fmt.Println()

	// 1. ARRAY
	demonstrasiArray()

	// 2. SLICE
	demonstrasiSlice()

	// 3. MAP
	demonstrasiMap()

	// 4. OPERASI LANJUTAN
	demonstrasiOperasiLanjutan()
}

func demonstrasiArray() {
	fmt.Println("1. ARRAY")
	fmt.Println("--------")

	// Deklarasi array dengan ukuran tetap
	var angka [5]int = [5]int{1, 2, 3, 4, 5}
	buah := [3]string{"Apel", "Jeruk", "Mangga"}

	// Ukuran otomatis berdasarkan jumlah elemen
	warna := [...]string{"Merah", "Hijau", "Biru", "Kuning"}

	fmt.Printf("Array angka: %v\n", angka)
	fmt.Printf("Array buah: %v\n", buah)
	fmt.Printf("Array warna: %v\n", warna)
	fmt.Printf("Panjang array warna: %d\n", len(warna))

	// Akses elemen
	fmt.Printf("Buah pertama: %s\n", buah[0])
	fmt.Printf("Buah terakhir: %s\n", buah[len(buah)-1])

	fmt.Println()
}

func demonstrasiSlice() {
	fmt.Println("2. SLICE")
	fmt.Println("--------")

	// Membuat slice kosong
	var bahasa []string
	fmt.Printf("Slice awal: %v (length: %d, capacity: %d)\n", bahasa, len(bahasa), cap(bahasa))

	// Append elemen
	bahasa = append(bahasa, "Go")
	bahasa = append(bahasa, "Python", "Java")
	fmt.Printf("Setelah append: %v (length: %d, capacity: %d)\n", bahasa, len(bahasa), cap(bahasa))
	fmt.Println()

	// Membuat slice dengan make
	angka := make([]int, 5, 10) // length: 5, capacity: 10
	fmt.Printf("Slice dengan make: %v (length: %d, capacity: %d)\n", angka, len(angka), cap(angka))
	fmt.Println()

	// Slicing operation
	buah := []string{"Apel", "Jeruk", "Mangga", "Pisang", "Anggur"}
	fmt.Printf("Slice buah: %v\n", buah)
	fmt.Printf("buah[1:3]: %v\n", buah[1:3])
	fmt.Printf("buah[:3]: %v\n", buah[:3])
	fmt.Printf("buah[2:]: %v\n", buah[2:])

	fmt.Println()
}

func demonstrasiMap() {
	fmt.Println("3. MAP")
	fmt.Println("------")

	// Membuat map
	populasi := make(map[string]int)
	populasi["Jakarta"] = 10000000
	populasi["Surabaya"] = 3000000
	populasi["Bandung"] = 2500000
	populasi["Medan"] = 2200000

	// Iterasi map
	fmt.Println("Data Populasi Kota:")
	for kota, jumlah := range populasi {
		fmt.Printf("%s: %d jiwa\n", kota, jumlah)
	}
	fmt.Println()

	// Cek keberadaan key
	if pop, exists := populasi["Jakarta"]; exists {
		fmt.Printf("Cek key 'Jakarta': Ada dengan populasi %d\n", pop)
	}

	// Delete key
	delete(populasi, "Medan")
	fmt.Printf("Setelah delete 'Medan': %d kota tersisa\n", len(populasi))

	fmt.Println()
}

func demonstrasiOperasiLanjutan() {
	fmt.Println("4. OPERASI LANJUTAN")
	fmt.Println("-------------------")

	// Slice of maps
	mahasiswa := []map[string]interface{}{
		{"nim": "123", "nama": "Budi", "ipk": 3.5},
		{"nim": "124", "nama": "Siti", "ipk": 3.8},
	}

	fmt.Println("Data Mahasiswa (Slice of Maps):")
	for i, mhs := range mahasiswa {
		fmt.Printf("Mahasiswa %d: NIM=%s, Nama=%s, IPK=%.2f\n",
			i+1, mhs["nim"], mhs["nama"], mhs["ipk"])
	}
	fmt.Println()

	// Copy slice
	original := []int{1, 2, 3, 4, 5}
	copy1 := make([]int, len(original))
	copy(copy1, original)
	copy1[0] = 999

	fmt.Printf("Original slice: %v\n", original)
	fmt.Printf("Copy slice: %v\n", copy1)
	fmt.Println()
}