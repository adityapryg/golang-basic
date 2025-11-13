package dto

import (
	"time"
)

// ============================================
// TODO REQUEST DTOs
// ============================================

// CreateTodoRequest untuk membuat todo baru
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending in_progress completed"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     string `json:"due_date" binding:"omitempty"` // Format: YYYY-MM-DD
}

// UpdateTodoRequest untuk update todo
type UpdateTodoRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"due_date"` // Format: YYYY-MM-DD or empty string to clear
}

// TodoCreateRequest untuk backward compatibility (alias)
type TodoCreateRequest = CreateTodoRequest

// TodoUpdateRequest untuk backward compatibility (alias)
type TodoUpdateRequest = UpdateTodoRequest

// TodoQueryParams untuk filter dan pagination
type TodoQueryParams struct {
	Status   string `form:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority int    `form:"priority" binding:"omitempty,min=0,max=5"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// ============================================
// TODO RESPONSE DTOs
// ============================================

// TodoResponse untuk response todo
type TodoResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	UserID      uint       `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TodoListResponse untuk response list todos dengan pagination
type TodoListResponse struct {
	Todos      []TodoResponse `json:"todos"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}
