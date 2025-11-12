# Aplikasi CLI Manajemen Mahasiswa

## Deskripsi
Aplikasi CLI **interaktif berbasis menu** untuk mengelola data mahasiswa dengan fitur CRUD lengkap.

## Cara Menjalankan

### Development (dengan Go)
```bash
cd 08-cli-app
go run main.go
```

### Production (Build Executable)
```bash
cd 08-cli-app
go build -o mahasiswa-app
./mahasiswa-app
```

## Tampilan Menu

Setelah menjalankan aplikasi, akan muncul menu interaktif:

```
=== Menu Mahasiswa ===
1. Tambah Mahasiswa
2. Lihat Semua Mahasiswa
3. Cari Mahasiswa
4. Update IPK
5. Hapus Mahasiswa
6. Tampilkan Statistik
7. Keluar
Silakan pilih opsi: _
```

## Contoh Penggunaan

### 1. Tambah Mahasiswa
```
Silakan pilih opsi: 1
Masukkan NIM: 123456
Masukkan Nama: Budi Santoso
Masukkan Jurusan: Teknik Informatika
Masukkan IPK: 3.75
✓ Mahasiswa berhasil ditambahkan
```

### 2. Lihat Semua Mahasiswa
```
Silakan pilih opsi: 2
NIM        Nama                 Jurusan              IPK       
------------------------------------------------------------
123456     Budi Santoso         Teknik Informatika   3.75
234567     Ani Wijaya           Sistem Informasi     3.85
```

### 3. Cari Mahasiswa
```
Silakan pilih opsi: 3
Masukkan NIM yang dicari: 123456
NIM: 123456
Nama: Budi Santoso
Jurusan: Teknik Informatika
IPK: 3.75
```

### 4. Update IPK
```
Silakan pilih opsi: 4
Masukkan NIM mahasiswa yang akan diupdate: 123456
Masukkan IPK baru: 3.80
✓ IPK mahasiswa berhasil diupdate
```

### 5. Hapus Mahasiswa
```
Silakan pilih opsi: 5
Masukkan NIM mahasiswa yang akan dihapus: 123456
Apakah Anda yakin ingin menghapus mahasiswa ini? (y/n): y
✓ Mahasiswa berhasil dihapus
```

### 6. Statistik
```
Silakan pilih opsi: 6
Total Mahasiswa: 5
Rata-rata IPK: 3.65
IPK Tertinggi: 3.95
IPK Terendah: 3.20
```

## Struktur Data

```go
type Mahasiswa struct {
    NIM     string  // Nomor Induk Mahasiswa
    Nama    string  // Nama Mahasiswa
    Jurusan string  // Jurusan Mahasiswa
    IPK     float64 // Indeks Prestasi Kumulatif
}
```

## Validasi
- **NIM**: Tidak boleh kosong
- **Nama**: Tidak boleh kosong
- **Jurusan**: Tidak boleh kosong
- **IPK**: Harus antara 0.0 - 4.0

## Konsep Go yang Diterapkan
- ✅ Structs untuk model data
- ✅ Slices untuk storage in-memory
- ✅ Functions modular
- ✅ Error handling (sentinel + custom errors)
- ✅ Control flow (switch, loops)
- ✅ bufio.Scanner untuk input
- ✅ String manipulation

## Repository
https://github.com/adityapryg/golang-demo
