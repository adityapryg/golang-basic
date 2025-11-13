package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityapryg/golang-demo/20-mini-project/config"
	"github.com/adityapryg/golang-demo/20-mini-project/handlers"
	"github.com/adityapryg/golang-demo/20-mini-project/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *AuthTestSuite) SetupSuite() {
	// Setup database untuk testing
	config.ConnectDatabase()

	// Setup router
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	auth := suite.router.Group("/api/v1/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}
}

func (suite *AuthTestSuite) TearDownTest() {
	// Bersihkan data testing setelah setiap test
	config.DB.Exec("DELETE FROM users WHERE username LIKE 'testuser%'")
}

func (suite *AuthTestSuite) TestRegisterSuccess() {
	reqBody := models.UserRegisterRequest{
		Username: "testuser1",
		Email:    "testuser1@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Registrasi berhasil", response.Message)
}

func (suite *AuthTestSuite) TestRegisterDuplicateUsername() {
	// Register pertama kali
	reqBody := models.UserRegisterRequest{
		Username: "testuser2",
		Email:    "testuser2@example.com",
		Password: "password123",
		FullName: "Test User 2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Register dengan username yang sama
	req2, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), 400, w2.Code)

	var response models.APIResponse
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Username sudah digunakan", response.Message)
}

func (suite *AuthTestSuite) TestLoginSuccess() {
	// Register user dulu
	registerBody := models.UserRegisterRequest{
		Username: "testuser3",
		Email:    "testuser3@example.com",
		Password: "password123",
		FullName: "Test User 3",
	}

	jsonBody, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Login
	loginBody := models.UserLoginRequest{
		Username: "testuser3",
		Password: "password123",
	}

	jsonBody, _ = json.Marshal(loginBody)
	req2, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), 200, w2.Code)

	var response models.APIResponse
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Login berhasil", response.Message)

	// Cek token ada di response
	data := response.Data.(map[string]interface{})
	assert.NotEmpty(suite.T(), data["token"])
}

func (suite *AuthTestSuite) TestLoginWrongPassword() {
	// Register user dulu
	registerBody := models.UserRegisterRequest{
		Username: "testuser4",
		Email:    "testuser4@example.com",
		Password: "password123",
		FullName: "Test User 4",
	}

	jsonBody, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Login dengan password salah
	loginBody := models.UserLoginRequest{
		Username: "testuser4",
		Password: "wrongpassword",
	}

	jsonBody, _ = json.Marshal(loginBody)
	req2, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), 401, w2.Code)

	var response models.APIResponse
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Username atau password salah", response.Message)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
