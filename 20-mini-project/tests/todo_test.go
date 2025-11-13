package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityapryg/golang-demo/20-mini-project/config"
	"github.com/adityapryg/golang-demo/20-mini-project/handlers"
	"github.com/adityapryg/golang-demo/20-mini-project/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/models"
	"github.com/adityapryg/golang-demo/20-mini-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TodoTestSuite struct {
	suite.Suite
	router *gin.Engine
	token  string
	userID uint
}

func (suite *TodoTestSuite) SetupSuite() {
	// Setup database untuk testing
	config.ConnectDatabase()

	// Setup router
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Setup routes
	v1 := suite.router.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware())
		{
			todos.POST("", handlers.CreateTodo)
			todos.GET("", handlers.GetAllTodos)
			todos.GET("/:id", handlers.GetTodoByID)
			todos.PUT("/:id", handlers.UpdateTodo)
			todos.DELETE("/:id", handlers.DeleteTodo)
		}
	}

	// Buat user dan token untuk testing
	hashedPassword, _ := utils.HashPassword("testpassword")
	user := models.User{
		Username: "todotest",
		Email:    "todotest@example.com",
		Password: hashedPassword,
		FullName: "Todo Test User",
	}
	config.DB.Create(&user)
	suite.userID = user.ID

	suite.token, _ = utils.GenerateToken(user.ID, user.Username)
}

func (suite *TodoTestSuite) TearDownSuite() {
	// Bersihkan data testing
	config.DB.Exec("DELETE FROM todos WHERE user_id = ?", suite.userID)
	config.DB.Exec("DELETE FROM users WHERE id = ?", suite.userID)
}

func (suite *TodoTestSuite) TestCreateTodoSuccess() {
	reqBody := models.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "This is a test todo",
		Status:      models.TodoStatusPending,
		Priority:    1,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Todo berhasil dibuat", response.Message)
}

func (suite *TodoTestSuite) TestCreateTodoWithoutAuth() {
	reqBody := models.TodoCreateRequest{
		Title:       "Test Todo Without Auth",
		Description: "This should fail",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// Tidak ada Authorization header

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 401, w.Code)
}

func (suite *TodoTestSuite) TestGetAllTodos() {
	// Buat beberapa todos dulu
	todos := []models.Todo{
		{Title: "Todo 1", Description: "Description 1", UserID: suite.userID, Status: models.TodoStatusPending},
		{Title: "Todo 2", Description: "Description 2", UserID: suite.userID, Status: models.TodoStatusCompleted},
	}
	for _, todo := range todos {
		config.DB.Create(&todo)
	}

	req, _ := http.NewRequest("GET", "/api/v1/todos", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)

	// Cek jumlah todos
	data := response.Data.([]interface{})
	assert.GreaterOrEqual(suite.T(), len(data), 2)
}

func (suite *TodoTestSuite) TestUpdateTodo() {
	// Buat todo dulu
	todo := models.Todo{
		Title:       "Original Title",
		Description: "Original Description",
		UserID:      suite.userID,
		Status:      models.TodoStatusPending,
	}
	config.DB.Create(&todo)

	// Update todo
	updateBody := models.TodoUpdateRequest{
		Title:  "Updated Title",
		Status: models.TodoStatusCompleted,
	}

	jsonBody, _ := json.Marshal(updateBody)
	url := fmt.Sprintf("/api/v1/todos/%d", todo.ID)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Todo berhasil diupdate", response.Message)
}

func (suite *TodoTestSuite) TestDeleteTodo() {
	// Buat todo dulu
	todo := models.Todo{
		Title:       "To Be Deleted",
		Description: "This will be deleted",
		UserID:      suite.userID,
		Status:      models.TodoStatusPending,
	}
	config.DB.Create(&todo)

	// Delete todo
	url := fmt.Sprintf("/api/v1/todos/%d", todo.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Todo berhasil dihapus", response.Message)

	// Verifikasi todo sudah dihapus (soft delete)
	var deletedTodo models.Todo
	err := config.DB.Unscoped().First(&deletedTodo, todo.ID).Error
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), deletedTodo.DeletedAt)
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
