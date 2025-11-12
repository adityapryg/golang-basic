package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/02-workspace-structure/utils"
)

func main() {
	// Mencetak header
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI GO WORKSPACE")
	fmt.Println("===========================================")
	fmt.Println()

	// Demonstrasi menggunakan package lokal
	fmt.Println("Menggunakan Package Lokal:")
	greeting := utils.SayHello("Aditya")
	fmt.Println(greeting)
	fmt.Println()

	// Demonstrasi berbagai fungsi dari utils
	fmt.Println("Berbagai Fungsi Greeting:")
	fmt.Println(utils.SayGoodMorning("Budi"))
	fmt.Println(utils.SayGoodAfternoon("Siti"))
	fmt.Println()

	// Menampilkan welcome message
	fmt.Println(utils.GetWelcomeMessage())
	fmt.Println()

	// Penjelasan Exported vs Unexported
	fmt.Println("Contoh Fungsi Exported vs Unexported:")
	fmt.Println("✓ SayHello() - Exported (huruf kapital)")
	fmt.Println("✗ greet() - Unexported (huruf kecil, tidak bisa diakses dari luar package)")
	fmt.Println()

	// Demonstrasi struktur workspace
	fmt.Println("Struktur Workspace:")
	fmt.Println("02-workspace-structure/")
	fmt.Println("├── go.mod")
	fmt.Println("├── main.go")
	fmt.Println("├── utils/")
	fmt.Println("│   └── greetings.go")
	fmt.Println("└── README.md")
}