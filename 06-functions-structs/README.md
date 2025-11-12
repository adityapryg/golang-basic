# Fungsi dan Struct di Go

Go adalah bahasa pemrograman yang mendukung pemrograman fungsional. Berikut adalah beberapa konsep terkait fungsi dan structs.

## Fungsi Sederhana
Fungsi sederhana di Go didefinisikan dengan kata kunci `func`:
```go
func HelloWorld() {
    fmt.Println("Hello, World!")
}
```

## Beberapa Nilai Kembali
Go memungkinkan fungsi mengembalikan lebih dari satu nilai:
```go
func GetValues() (int, string) {
    return 1, "Hello"
}
```

## Nilai yang Diberi Nama
Kita bisa memberi nama pada nilai kembali untuk lebih mudah dipahami:
```go
func GetDetails() (id int, name string) {
    return 1, "Aditya"
}
```

## Definisi Struct
Struct digunakan untuk mengelompokkan data:
```go
type Person struct {
    Name string
    Age  int
}
```

## Receiver Value dan Pointer
Kita dapat menggunakan receiver value dan pointer saat mendefinisikan metode:
```go
func (p Person) Greet() {
    fmt.Println("Hello, my name is", p.Name)
}

func (p *Person) HaveBirthday() {
    p.Age++
}
```

## Pola Konstruktor
Kita dapat menggunakan fungsi untuk membuat struct baru (konstruktor):
```go
func NewPerson(name string, age int) *Person {
    return &Person{Name: name, Age: age}
}
```

## Fungsi Variadic
Fungsi yang menerima argumen variable:
```go
func Sum(numbers ...int) int {
    total := 0
    for _, number := range numbers {
        total += number
    }
    return total
}
```

## Praktik Terbaik
- Selalu gunakan nama fungsi yang deskriptif.
- Jaga agar fungsi tetap kecil dan fokus pada satu tugas.

## Pemecahan Masalah
Beberapa kesalahan umum yang bisa terjadi dan cara mengatasinya:
- Kesalahan dalam mendefinisikan fungsi: cek jumlah argumen dan tipe return.

## Latihan Praktis
Buat fungsi yang menghitung luas segitiga berdasarkan alas dan tinggi.