# Demonstrasi Hello World

## Tujuan Pembelajaran  
Mengerti cara membuat aplikasi Go sederhana yang mencetak "Hello World" ke layar.

## Penjelasan Struktur Dasar  
1. **Package**: Bagian paling dasar dari program Go, setiap file memiliki package. Misalnya, `package main` berarti file ini adalah bagian dari package utama.
2. **Import**: Menyertakan pustaka yang diperlukan. Dalam demo ini kita menggunakan `fmt` untuk format output.
3. **Function main**: Titik awal eksekusi program, di mana kode kita akan mulai dieksekusi.

## Langkah-langkah Demonstrasi  
1. Buat file baru dengan nama `main.go` di dalam folder `01-hello-world`.
2. Tulis kode berikut:
   ```go
   package main
   
   import (
       "fmt"
   )
   
   func main() {
       fmt.Println("Hello World")
   }
   ```  
3. Simpan file tersebut.

## Cara Menjalankan  
Untuk menjalankan aplikasi, buka terminal dan navigasikan ke folder `01-hello-world`, lalu jalankan perintah:
```bash
go run main.go
```

## Output yang Diharapkan  
```
Hello World
```

## Konsep Penting  
- **Go**: Bahasa pemrograman yang sederhana dan efisien.
- **Fungsi**: Sekumpulan kode yang dijalankan saat dipanggil.
- **Output**: Hasil dari eksekusi program yang ditampilkan di layar.

## Latihan  
1. Modifikasi program untuk mencetak nama Anda.
2. Tambahkan lebih banyak fungsi untuk menghitung sum dua angka.

## Prompt GitHub Copilot  
"Buatlah program Go untuk mencetak 'Hello World' ke layar."

## Pemecahan Masalah  
- Jika terjadi kesalahan saat menjalankan program, periksa kembali kode Anda terutama pada bagian sintaks dan struktur. 
- Pastikan Anda telah menginstal Go dan menambahkannya ke PATH dengan benar.