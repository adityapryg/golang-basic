package main

import (
	"errors"
	"fmt"
)

// Calculator adalah struct untuk operasi matematika
type Calculator struct{}

// Add menjumlahkan dua bilangan
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// Subtract mengurangkan dua bilangan
func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

// Multiply mengalikan dua bilangan
func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

// Divide membagi dua bilangan
func (c *Calculator) Divide(a, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("tidak bisa membagi dengan nol")
	}
	return float64(a) / float64(b), nil
}

// User adalah struct untuk data user
type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

// UserService adalah service untuk operasi user
type UserService struct {
	users map[int]*User
}

// NewUserService membuat instance UserService baru
func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]*User),
	}
}

// CreateUser membuat user baru
func (s *UserService) CreateUser(username, email string, age int) (*User, error) {
	// Validasi
	if username == "" {
		return nil, errors.New("username tidak boleh kosong")
	}
	if email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("umur tidak valid")
	}

	// Cek duplicate username
	for _, u := range s.users {
		if u.Username == username {
			return nil, errors.New("username sudah digunakan")
		}
	}

	// Buat user
	id := len(s.users) + 1
	user := &User{
		ID:       id,
		Username: username,
		Email:    email,
		Age:      age,
	}
	s.users[id] = user

	return user, nil
}

// GetUser mendapatkan user berdasarkan ID
func (s *UserService) GetUser(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

// GetAllUsers mendapatkan semua users
func (s *UserService) GetAllUsers() []*User {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// UpdateUser mengupdate data user
func (s *UserService) UpdateUser(id int, username, email string, age int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user tidak ditemukan")
	}

	// Validasi
	if username == "" {
		return nil, errors.New("username tidak boleh kosong")
	}
	if email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("umur tidak valid")
	}

	// Update
	user.Username = username
	user.Email = email
	user.Age = age

	return user, nil
}

// DeleteUser menghapus user
func (s *UserService) DeleteUser(id int) error {
	if _, exists := s.users[id]; !exists {
		return errors.New("user tidak ditemukan")
	}
	delete(s.users, id)
	return nil
}

// IsAdult mengecek apakah user sudah dewasa
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

// GetEmailDomain mendapatkan domain dari email
func (u *User) GetEmailDomain() string {
	for i, char := range u.Email {
		if char == '@' {
			return u.Email[i+1:]
		}
	}
	return ""
}

func main() {
	fmt.Println("===========================================")
	fmt.Println("   TESTING DI GO")
	fmt.Println("===========================================\n")

	fmt.Println("Program ini untuk mendemonstrasikan testing di Go.")
	fmt.Println("Kode utama ada di file ini (main.go)")
	fmt.Println("Test cases ada di file main_test.go")
	fmt.Println("\nUntuk menjalankan tests:")
	fmt.Println("  go test -v")
	fmt.Println("  go test -v -cover")
	fmt.Println("  go test -v -coverprofile=coverage.out")
	fmt.Println("\nUntuk melihat coverage report:")
	fmt.Println("  go tool cover -html=coverage.out")

	// Demo penggunaan
	fmt.Println("\n" + "===========================================")
	fmt.Println("   DEMO PENGGUNAAN")
	fmt.Println("===========================================\n")

	calc := &Calculator{}
	fmt.Printf("5 + 3 = %d\n", calc.Add(5, 3))
	fmt.Printf("5 - 3 = %d\n", calc.Subtract(5, 3))
	fmt.Printf("5 * 3 = %d\n", calc.Multiply(5, 3))
	result, _ := calc.Divide(10, 2)
	fmt.Printf("10 / 2 = %.2f\n", result)

	userService := NewUserService()
	user, _ := userService.CreateUser("john_doe", "john@example.com", 25)
	fmt.Printf("\nUser created: %s (%s), Age: %d\n", user.Username, user.Email, user.Age)
	fmt.Printf("Is adult: %v\n", user.IsAdult())
	fmt.Printf("Email domain: %s\n", user.GetEmailDomain())
}
