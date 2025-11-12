# README untuk Topik 05: Data Structures

## Tujuan Pembelajaran
Pembelajaran ini bertujuan untuk memahami konsep dasar data structures dalam bahasa pemrograman Go, termasuk cara menggunakan Array, Slice, dan Map.

## Penjelasan
- **Array**: Struktur data dengan ukuran tetap, yang menyimpan elemen-elemen dengan tipe yang sama dan dapat diakses melalui indeks.
- **Slice**: Struktur data dinamis yang dapat berubah ukuran, merupakan bagian dari array dan memudahkan pengelolaan koleksi data.
- **Map**: Struktur data yang menyimpan pasangan key-value, memungkinkan pencarian data yang cepat berdasarkan kunci unik.

## Karakteristik Masing-Masing
- **Array**: Ukuran tidak dapat diubah setelah deklarasi, semua elemen harus memiliki tipe data yang sama.
- **Slice**: Ukuran dapat diubah, memiliki kemampuan untuk menambah dan mengurangi elemen, lebih fleksibel dibandingkan array.
- **Map**: Tidak berurutan, memungkinkan operasi pencarian yang efisien, dapat menyimpan data dengan tipe key dan value yang berbeda.

## Cara Menjalankan Demo
1. Clone repository ini.
2. Jalankan perintah `go run main.go` di terminal.

## Output yang Diharapkan
Output dari demo ini akan menunjukkan bagaimana masing-masing struktur data bekerja dan cara implementasinya dalam kode Go.

## Operasi Penting untuk Slice dan Map
- **Slice**: `append()`, `copy()`, `len()`, `cap()`.
- **Map**: Menambah item baru, menghapus item, dan melakukan pencarian menggunakan kunci.

## Latihan Praktis
Buatlah program sederhana yang menggunakan Array, Slice, dan Map untuk mengelola data mahasiswa, termasuk informasi seperti nama, umur, dan jurusan.

## GitHub Copilot Prompts
Gunakan prompt berikut untuk meminta bantuan dari GitHub Copilot:
- "Tulis fungsi untuk menambah elemen ke dalam slice."
- "Buat program yang mencari nilai dalam map."

## Best Practices
- Gunakan slice saat ukuran data berfariasi, dan map saat data perlu diakses dengan kunci.
- Hindari penggunaan array kecuali untuk ukuran tetap.

## Troubleshooting Common Issues
- Jika program tidak berjalan, periksa tipe data yang digunakan di array dan map.
- Pastikan untuk menginisialisasi slice sebelum menambah elemen.

## Tabel Perbandingan
| Kriteria  | Array | Slice | Map  |
|-----------|-------|-------|------|
| Ukuran    | Tetap | Dinamis | Dinamis |
| Akses     | Indeks | Indeks | Key   |
| Penyimpanan | Homogen | Homogen | Heterogen |
| Kinerja   | Tinggi  | Menengah | Tinggi  |