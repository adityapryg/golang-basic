package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ========================================
// BASIC TESTING (Table-Driven Tests)
// ========================================

func TestCalculator_Add(t *testing.T) {
	calc := &Calculator{}

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed numbers", 5, -3, 2},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculator_Divide(t *testing.T) {
	calc := &Calculator{}

	t.Run("valid division", func(t *testing.T) {
		result, err := calc.Divide(10, 2)
		assert.NoError(t, err)
		assert.Equal(t, 5.0, result)
	})

	t.Run("divide by zero", func(t *testing.T) {
		result, err := calc.Divide(10, 0)
		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
		assert.EqualError(t, err, "tidak bisa membagi dengan nol")
	})
}

// ========================================
// TEST SUITE (testify/suite)
// ========================================

type UserServiceTestSuite struct {
	suite.Suite
	service *UserService
}

// SetupTest dijalankan sebelum setiap test
func (suite *UserServiceTestSuite) SetupTest() {
	suite.service = NewUserService()
}

// TearDownTest dijalankan setelah setiap test
func (suite *UserServiceTestSuite) TearDownTest() {
	// Cleanup jika diperlukan
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
	user, err := suite.service.CreateUser("john_doe", "john@example.com", 25)

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal("john_doe", user.Username)
	suite.Equal("john@example.com", user.Email)
	suite.Equal(25, user.Age)
	suite.Equal(1, user.ID)
}

func (suite *UserServiceTestSuite) TestCreateUser_EmptyUsername() {
	user, err := suite.service.CreateUser("", "john@example.com", 25)

	suite.Error(err)
	suite.Nil(user)
	suite.EqualError(err, "username tidak boleh kosong")
}

func (suite *UserServiceTestSuite) TestCreateUser_EmptyEmail() {
	user, err := suite.service.CreateUser("john_doe", "", 25)

	suite.Error(err)
	suite.Nil(user)
	suite.EqualError(err, "email tidak boleh kosong")
}

func (suite *UserServiceTestSuite) TestCreateUser_InvalidAge() {
	tests := []struct {
		name string
		age  int
	}{
		{"negative age", -1},
		{"too old", 151},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			user, err := suite.service.CreateUser("john_doe", "john@example.com", tt.age)
			suite.Error(err)
			suite.Nil(user)
			suite.EqualError(err, "umur tidak valid")
		})
	}
}

func (suite *UserServiceTestSuite) TestCreateUser_DuplicateUsername() {
	// Buat user pertama
	_, err := suite.service.CreateUser("john_doe", "john@example.com", 25)
	suite.NoError(err)

	// Coba buat user dengan username sama
	user, err := suite.service.CreateUser("john_doe", "jane@example.com", 30)
	suite.Error(err)
	suite.Nil(user)
	suite.EqualError(err, "username sudah digunakan")
}

func (suite *UserServiceTestSuite) TestGetUser_Success() {
	// Setup: buat user
	createdUser, _ := suite.service.CreateUser("john_doe", "john@example.com", 25)

	// Test
	user, err := suite.service.GetUser(createdUser.ID)

	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(createdUser.ID, user.ID)
	suite.Equal(createdUser.Username, user.Username)
}

func (suite *UserServiceTestSuite) TestGetUser_NotFound() {
	user, err := suite.service.GetUser(999)

	suite.Error(err)
	suite.Nil(user)
	suite.EqualError(err, "user tidak ditemukan")
}

func (suite *UserServiceTestSuite) TestGetAllUsers() {
	// Setup: buat beberapa users
	suite.service.CreateUser("user1", "user1@example.com", 20)
	suite.service.CreateUser("user2", "user2@example.com", 25)
	suite.service.CreateUser("user3", "user3@example.com", 30)

	// Test
	users := suite.service.GetAllUsers()

	suite.Len(users, 3)
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	// Setup
	createdUser, _ := suite.service.CreateUser("john_doe", "john@example.com", 25)

	// Test
	updatedUser, err := suite.service.UpdateUser(createdUser.ID, "john_updated", "new@example.com", 30)

	suite.NoError(err)
	suite.NotNil(updatedUser)
	suite.Equal("john_updated", updatedUser.Username)
	suite.Equal("new@example.com", updatedUser.Email)
	suite.Equal(30, updatedUser.Age)
}

func (suite *UserServiceTestSuite) TestUpdateUser_NotFound() {
	user, err := suite.service.UpdateUser(999, "john", "john@example.com", 25)

	suite.Error(err)
	suite.Nil(user)
	suite.EqualError(err, "user tidak ditemukan")
}

func (suite *UserServiceTestSuite) TestDeleteUser_Success() {
	// Setup
	createdUser, _ := suite.service.CreateUser("john_doe", "john@example.com", 25)

	// Test
	err := suite.service.DeleteUser(createdUser.ID)
	suite.NoError(err)

	// Verify deleted
	user, err := suite.service.GetUser(createdUser.ID)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UserServiceTestSuite) TestDeleteUser_NotFound() {
	err := suite.service.DeleteUser(999)
	suite.Error(err)
	suite.EqualError(err, "user tidak ditemukan")
}

// TestUserServiceTestSuite menjalankan test suite
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// ========================================
// UNIT TESTS untuk User methods
// ========================================

func TestUser_IsAdult(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"child", 10, false},
		{"teenager", 17, false},
		{"adult (18)", 18, true},
		{"adult (25)", 25, true},
		{"senior", 70, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Age: tt.age}
			result := user.IsAdult()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_GetEmailDomain(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"gmail", "user@gmail.com", "gmail.com"},
		{"company", "john@company.co.id", "company.co.id"},
		{"subdomain", "admin@mail.company.com", "mail.company.com"},
		{"no at sign", "invalidemail", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Email: tt.email}
			result := user.GetEmailDomain()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ========================================
// BENCHMARK TESTS
// ========================================

func BenchmarkCalculator_Add(b *testing.B) {
	calc := &Calculator{}
	for i := 0; i < b.N; i++ {
		calc.Add(5, 3)
	}
}

func BenchmarkCalculator_Multiply(b *testing.B) {
	calc := &Calculator{}
	for i := 0; i < b.N; i++ {
		calc.Multiply(5, 3)
	}
}

func BenchmarkUserService_CreateUser(b *testing.B) {
	service := NewUserService()
	b.ResetTimer() // Reset timer setelah setup

	for i := 0; i < b.N; i++ {
		service.CreateUser("user", "user@example.com", 25)
	}
}
