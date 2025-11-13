package service

import (
	"errors"
	"time"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"gorm.io/gorm"
)

var (
	// ErrTodoNotFound is returned when todo is not found
	ErrTodoNotFound = errors.New("todo not found")
	// ErrUnauthorizedAccess is returned when user tries to access todo they don't own
	ErrUnauthorizedAccess = errors.New("unauthorized access to todo")
	// ErrInvalidStatus is returned when status value is invalid
	ErrInvalidStatus = errors.New("invalid status value")
	// ErrInvalidPriority is returned when priority value is invalid
	ErrInvalidPriority = errors.New("invalid priority value")
)

// TodoService handles todo business logic
type TodoService struct {
	todoRepo *repository.TodoRepository
}

// NewTodoService creates a new todo service instance
func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

// CreateTodo creates a new todo for a user
func (s *TodoService) CreateTodo(userID uint, req dto.CreateTodoRequest) (*model.Todo, error) {
	// Validate status
	if !isValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	// Validate priority
	if !isValidPriority(req.Priority) {
		return nil, ErrInvalidPriority
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		dueDate = &parsedDate
	}

	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     dueDate,
		UserID:      userID,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetTodoByID retrieves a todo by ID with authorization check
func (s *TodoService) GetTodoByID(todoID, userID uint) (*model.Todo, error) {
	todo, err := s.todoRepo.FindByID(todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	// Check if user owns this todo
	if todo.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	return todo, nil
}

// GetUserTodos retrieves all todos for a user with optional filters
func (s *TodoService) GetUserTodos(userID uint, status, priority string) ([]model.Todo, error) {
	// Validate filters if provided
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	if priority != "" && !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}

	return s.todoRepo.FindByUserIDWithFilters(userID, status, priority)
}

// UpdateTodo updates a todo with authorization check
func (s *TodoService) UpdateTodo(todoID, userID uint, req dto.UpdateTodoRequest) (*model.Todo, error) {
	// Check if todo exists and user owns it
	todo, err := s.GetTodoByID(todoID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}

	if req.Description != nil {
		todo.Description = *req.Description
	}

	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		todo.Status = *req.Status
	}

	if req.Priority != nil {
		if !isValidPriority(*req.Priority) {
			return nil, ErrInvalidPriority
		}
		todo.Priority = *req.Priority
	}

	if req.DueDate != nil {
		if *req.DueDate == "" {
			todo.DueDate = nil
		} else {
			parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				return nil, errors.New("invalid date format, use YYYY-MM-DD")
			}
			todo.DueDate = &parsedDate
		}
	}

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo with authorization check
func (s *TodoService) DeleteTodo(todoID, userID uint) error {
	// Check if todo exists and user owns it
	_, err := s.GetTodoByID(todoID, userID)
	if err != nil {
		return err
	}

	return s.todoRepo.Delete(todoID)
}

// Helper functions for validation

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}
	return validStatuses[status]
}

func isValidPriority(priority string) bool {
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return validPriorities[priority]
}
