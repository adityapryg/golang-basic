package main

import "fmt"

func main() {
    // Contoh penggunaan if-else
    score := 80
    if score >= 75 {
        fmt.Println("Selamat! Anda lulus.")
    } else {
        fmt.Println("Maaf, Anda tidak lulus.")
    }

    // Contoh penggunaan switch
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

    // Contoh penggunaan for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    // Contoh break dan continue
    for i := 0; i < 5; i++ {
        if i == 2 {
            continue
        }
        fmt.Println(i)
    }
}