package service

import (
	"errors"
	"fmt"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/repository"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
)

var (
	// ErrUserExists ketika username atau email sudah terdaftar
	ErrUserExists = errors.New("username or email already exists")
	// ErrInvalidCredentials ketika username atau password salah
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUserNotFound ketika user tidak ditemukan
	ErrUserNotFound = errors.New("user not found")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register mendaftarkan user baru
func (s *AuthService) Register(req dto.UserRegisterRequest) (*dto.UserResponse, error) {
	// Business Rule 1: Check if username already exists
	existsUsername, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existsUsername {
		return nil, ErrUserExists
	}

	// Business Rule 2: Check if email already exists
	existsEmail, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existsEmail {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	// Save to database via repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert model to DTO response (without password)
	return s.toUserResponse(user), nil
}

// Login melakukan autentikasi user dan mengembalikan JWT token
func (s *AuthService) Login(req dto.UserLoginRequest) (*dto.AuthResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return auth response with token and user data
	return &dto.AuthResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

// GetProfile mendapatkan profile user berdasarkan ID
func (s *AuthService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.toUserResponse(user), nil
}

// UpdateProfile mengupdate profile user
func (s *AuthService) UpdateProfile(userID uint, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	// Find existing user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Business Rule: Check if email already used by another user
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrUserExists
		}
		user.Email = req.Email
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// Helper: Convert model.User to dto.UserResponse
func (s *AuthService) toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
