package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("   DEMONSTRASI ERROR HANDLING")
	fmt.Println("===========================================")
	fmt.Println()

	demonstrasiBasicError()
	demonstrasiCustomError()
	demonstrasiSentinelError()
	demonstrasiErrorWrapping()
	demonstrasiMultipleErrors()
}

// 1. BASIC ERROR HANDLING
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("pembagian dengan nol tidak diperbolehkan")
	}
	return a / b, nil
}

func demonstrasiBasicError() {
	fmt.Println("1. BASIC ERROR HANDLING")
	fmt.Println("----------------------")
	
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Hasil pembagian: %.2f\n", result)
	}
	
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println()
}

// 2. CUSTOM ERROR TYPES
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validasi gagal pada field '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "umur tidak boleh negatif"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "umur tidak realistis"}
	}
	return nil
}

func demonstrasiCustomError() {
	fmt.Println("2. CUSTOM ERROR TYPES")
	fmt.Println("--------------------")
	
	if err := validateAge(25); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Validasi berhasil untuk umur 25")
	}
	
	if err := validateAge(-5); err != nil {
		fmt.Println("Error:", err)
	}
	
	if err := validateAge(200); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println()
}

// 3. SENTINEL ERRORS
var (
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrAlreadyExists = errors.New("data sudah ada")
	ErrInvalidInput  = errors.New("input tidak valid")
)

type User struct {
	ID   string
	Name string
}

var users = map[string]User{
	"1": {ID: "1", Name: "John Doe"},
	"2": {ID: "2", Name: "Jane Smith"},
}

func findUser(id string) (*User, error) {
	user, exists := users[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &user, nil
}

func demonstrasiSentinelError() {
	fmt.Println("3. SENTINEL ERRORS")
	fmt.Println("-----------------")
	
	user, err := findUser("1")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("User ditemukan: %s\n", user.Name)
	}
	
	_, err = findUser("999")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Error: data tidak ditemukan (sentinel error)")
		} else {
			fmt.Println("Error:", err)
		}
	}
	fmt.Println()
}

// 4. ERROR WRAPPING
func processUser(id string) error {
	user, err := findUser(id)
	if err != nil {
		return fmt.Errorf("gagal memproses user %s: %w", id, err)
	}
	fmt.Printf("Memproses user: %s\n", user.Name)
	return nil
}

func demonstrasiErrorWrapping() {
	fmt.Println("4. ERROR WRAPPING")
	fmt.Println("----------------")
	
	if err := processUser("1"); err != nil {
		fmt.Println("Error:", err)
	}
	
	err := processUser("999")
	if err != nil {
		fmt.Println("Error:", err)
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Original error adalah ErrNotFound: true")
		}
	}
	fmt.Println()
}

// 5. MULTIPLE ERRORS
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred", len(m.Errors))
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

func validateInput(email, password string) error {
	var errs MultiError
	
	if email == "" {
		errs.Add(&ValidationError{Field: "email", Message: "email tidak boleh kosong"})
	} else if len(email) < 5 {
		errs.Add(&ValidationError{Field: "email", Message: "email terlalu pendek"})
	}
	
	if password == "" {
		errs.Add(&ValidationError{Field: "password", Message: "password tidak boleh kosong"})
	} else if len(password) < 8 {
		errs.Add(&ValidationError{Field: "password", Message: "password minimal 8 karakter"})
	}
	
	if errs.HasErrors() {
		return &errs
	}
	return nil
}

func demonstrasiMultipleErrors() {
	fmt.Println("5. MULTIPLE ERRORS")
	fmt.Println("-----------------")
	
	if err := validateInput("user@email.com", "password123"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Validasi berhasil!")
	}
	
	err := validateInput("abc", "pass")
	if err != nil {
		if multiErr, ok := err.(*MultiError); ok {
			fmt.Printf("Validasi memiliki %d error:\n", len(multiErr.Errors))
			for _, e := range multiErr.Errors {
				fmt.Printf("- %s\n", e.Error())
			}
		}
	}
	fmt.Println()
}