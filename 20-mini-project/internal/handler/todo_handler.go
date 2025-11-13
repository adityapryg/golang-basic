package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TodoHandler handles todo HTTP requests
type TodoHandler struct {
	todoService *service.TodoService
}

// NewTodoHandler creates a new todo handler instance
func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// Create handles POST /api/v1/todos
// @Summary Create a new todo
// @Description Create a new todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body dto.CreateTodoRequest true "Todo data"
// @Success 201 {object} dto.SuccessResponse{data=dto.TodoResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/todos [post]
// @Security BearerAuth
func (h *TodoHandler) Create(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	todo, err := h.todoService.CreateTodo(userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to create todo"

		if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	response := dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		Priority:    todo.Priority,
		DueDate:     todo.DueDate,
		UserID:      todo.UserID,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Todo created successfully",
		Data:    response,
	})
}

// GetAll handles GET /api/v1/todos
// @Summary Get all todos for authenticated user
// @Description Retrieve all todos for the authenticated user with optional filters
// @Tags todos
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (pending, in_progress, completed)"
// @Param priority query string false "Filter by priority (low, medium, high)"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.TodoResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/todos [get]
// @Security BearerAuth
func (h *TodoHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	// Get query parameters
	status := c.Query("status")
	priority := c.Query("priority")

	todos, err := h.todoService.GetUserTodos(userID.(uint), status, priority)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to retrieve todos"

		if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	// Convert to response DTOs
	responses := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = dto.TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			UserID:      todo.UserID,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todos retrieved successfully",
		Data:    responses,
	})
}

// GetByID handles GET /api/v1/todos/:id
// @Summary Get a specific todo
// @Description Retrieve a specific todo by ID for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.TodoResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/todos/{id} [get]
// @Security BearerAuth
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid todo ID",
			Error:   err.Error(),
		})
		return
	}

	todo, err := h.todoService.GetTodoByID(uint(todoID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to retrieve todo"

		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
			message = "Todo not found"
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
			message = "You don't have permission to access this todo"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	response := dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		Priority:    todo.Priority,
		DueDate:     todo.DueDate,
		UserID:      todo.UserID,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo retrieved successfully",
		Data:    response,
	})
}

// Update handles PUT /api/v1/todos/:id
// @Summary Update a todo
// @Description Update a specific todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body dto.UpdateTodoRequest true "Todo data to update"
// @Success 200 {object} dto.SuccessResponse{data=dto.TodoResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/todos/{id} [put]
// @Security BearerAuth
func (h *TodoHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid todo ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	todo, err := h.todoService.UpdateTodo(uint(todoID), userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to update todo"

		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
			message = "Todo not found"
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
			message = "You don't have permission to update this todo"
		} else if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrInvalidPriority) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusNotFound
			message = "Todo not found"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	response := dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		Priority:    todo.Priority,
		DueDate:     todo.DueDate,
		UserID:      todo.UserID,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo updated successfully",
		Data:    response,
	})
}

// Delete handles DELETE /api/v1/todos/:id
// @Summary Delete a todo
// @Description Delete a specific todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/todos/{id} [delete]
// @Security BearerAuth
func (h *TodoHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid todo ID",
			Error:   err.Error(),
		})
		return
	}

	err = h.todoService.DeleteTodo(uint(todoID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to delete todo"

		if errors.Is(err, service.ErrTodoNotFound) {
			statusCode = http.StatusNotFound
			message = "Todo not found"
		} else if errors.Is(err, service.ErrUnauthorizedAccess) {
			statusCode = http.StatusForbidden
			message = "You don't have permission to delete this todo"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Todo deleted successfully",
		Data:    nil,
	})
}
