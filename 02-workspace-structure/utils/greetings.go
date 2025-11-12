package utils

import "fmt"

// SayHello adalah fungsi exported (huruf kapital)
// Fungsi ini bisa diakses dari package lain
func SayHello(name string) string {
	return fmt.Sprintf("Selamat Pagi, %s!", name)
}

// SayGoodMorning - fungsi exported dengan dokumentasi
// Dokumentasi ini akan muncul di IDE dan godoc
func SayGoodMorning(name string) string {
	return greet("Pagi", name)
}

// SayGoodAfternoon - exported function
func SayGoodAfternoon(name string) string {
	return greet("Siang", name)
}

// greet adalah fungsi unexported (huruf kecil)
// Fungsi ini hanya bisa diakses dalam package utils
func greet(timeOfDay, name string) string {
	return fmt.Sprintf("Selamat %s, %s!", timeOfDay, name)
}

// GetWelcomeMessage - contoh fungsi exported lainnya
func GetWelcomeMessage() string {
	return "Selamat datang di demonstrasi Go Workspace!"
}
