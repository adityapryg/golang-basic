# Alur Kendali di Go

## Introduction

Alur kendali adalah cara program menentukan urutan eksekusi instruksi. Dalam pemrograman Go, kita memiliki beberapa struktur kontrol seperti `if-else`, `switch`, dan `for` loop. Di sini, kita akan membahas masing-masing dengan contoh dan latihan.

## If-Else

Pernyataan `if` memungkinkan Anda untuk melakukan keputusan berdasarkan kondisi tertentu. Berikut adalah contoh penggunaannya:

```go
package main

import "fmt"

func main() {
    var score int = 80

    if score >= 75 {
        fmt.Println("Selamat! Anda lulus.")
    } else {
        fmt.Println("Maaf, Anda tidak lulus.")
    }
}
```

### Latihan

Tulis program yang memeriksa apakah angka itu genap atau ganjil menggunakan `if-else`.

## Switch

Pernyataan `switch` adalah alternatif untuk `if-else` dalam memeriksa banyak kondisi:

```go
package main

import "fmt"

func main() {
    day := 3

    switch day {
    case 1:
        fmt.Println("Hari Senin")
    case 2:
        fmt.Println("Hari Selasa")
    case 3:
        fmt.Println("Hari Rabu")
    default:
        fmt.Println("Hari tidak diketahui")
    }
}
```

### Latihan

Buat program menggunakan `switch` untuk mencetak nama bulan berdasarkan nomor bulan.

## For Loop

`for` loop digunakan untuk iterasi. Berikut adalah beberapa gaya:

### Gaya Dasar
```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### Gaya Range
```go
fruits := []string{"apel", "jeruk", "pisang"}
for i, fruit := range fruits {
    fmt.Println(i, fruit)
}
```

### Latihan

Buat program yang menggunakan `for` loop untuk menghitung jumlah dari 1 hingga 100.

## Break dan Continue

Dengan pernyataan `break`, Anda dapat keluar dari loop, sementara `continue` digunakan untuk melewatkan iterasi saat ini:

```go
for i := 0; i < 5; i++ {
    if i == 2 {
        continue
    }
    fmt.Println(i)
}
```

### Latihan

Buat program yang mencetak angka dari 0 sampai 10, tetapi lewati angka 5.

## Copilot Prompts

- Jelaskan apa itu alur kendali.
- Berikan contoh penggunaan `for` loop.
- Apa bedanya `switch` dan `if-else`?