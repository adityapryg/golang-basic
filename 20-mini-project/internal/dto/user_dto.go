package dto

import "time"

// ============================================
// USER REQUEST DTOs
// ============================================

// UserRegisterRequest untuk registrasi user baru
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
}

// UserLoginRequest untuk login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest untuk update profile
type UserUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	FullName string `json:"full_name" binding:"max=100"`
}

// ============================================
// USER RESPONSE DTOs
// ============================================

// UserResponse untuk response user tanpa password
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse untuk response setelah login/register
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ============================================
// GENERIC RESPONSE DTOs
// ============================================

// SuccessResponse untuk response sukses
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse untuk response error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
