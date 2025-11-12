package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	// Mencetak header
	fmt.Println("===========================================")
	fmt.Println("   SELAMAT DATANG DI GOLANG!")
	fmt.Println("===========================================")
	fmt.Println()

	// Mengambil nama dari command-line arguments
	var nama string
	if len(os.Args) > 1 {
		nama = os.Args[1]
	} else {
		nama = "Dunia"
	}

	// Mencetak sapaan
	fmt.Printf("Halo, %s!\n\n", nama)

	// Informasi sistem
	fmt.Println("Informasi System:")
	fmt.Printf("- Waktu: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("- Go Version: %s\n", runtime.Version())
	fmt.Printf("- OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	// Pesan penutup
	fmt.Println("Program selesai dijalankan.")
}