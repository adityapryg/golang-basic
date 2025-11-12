# Aplikasi CLI untuk Sistem Manajemen Mahasiswa

## Tujuan
Aplikasi ini dirancang untuk memudahkan pengelolaan data mahasiswa secara efisien melalui antarmuka baris perintah.

## Fitur
- **CRUD Operasi untuk Mahasiswa**: Membuat, Membaca, Memperbarui, dan Menghapus record mahasiswa.

## Struktur Data
```go
type Mahasiswa struct {
    NIM   string // Nomor Induk Mahasiswa
    Nama  string // Nama Mahasiswa
    Jurusan string // Jurusan Mahasiswa
    IPK   float64 // Indeks Prestasi Kumulatif
}
```

## Pola Penanganan Kesalahan
Aplikasi ini menggunakan pola penanganan kesalahan yang jelas agar pengguna dapat memahami dan menangani kesalahan yang terjadi.

## Contoh Penggunaan
```bash
# Menambahkan mahasiswa baru
cli-app add --nim 123456 --nama "Ali" --jurusan "Teknik Informatika" --ipk 3.5

# Melihat semua mahasiswa
cli-app list

# Memperbarui data mahasiswa
cli-app update --nim 123456 --ipk 3.7

# Menghapus mahasiswa
cli-app delete --nim 123456
```

## Aturan Validasi
- NIM harus 6 digit
- Nama tidak boleh kosong
- Jurusan tidak boleh kosong
- IPK harus antara 0.0 dan 4.0

## Praktik Terbaik
- Gunakan komentar dan dokumentasi yang jelas dalam kode.
- Pastikan semua input dari pengguna divalidasi.
- Tangani kesalahan dengan baik dan beri tahu pengguna.
- Lakukan pengujian untuk memastikan semua fungsi berjalan dengan baik.
