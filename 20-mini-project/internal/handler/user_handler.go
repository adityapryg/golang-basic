package handler

import (
	"errors"
	"net/http"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/dto"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/middleware"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	authService *service.AuthService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Call service
	user, err := h.authService.Register(req)
	if err != nil {
		// Map service errors to HTTP status codes
		statusCode := http.StatusInternalServerError
		message := "Failed to register user"

		if errors.Is(err, service.ErrUserExists) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login
// @Summary User login
// @Description Login with username and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequest true "Login credentials"
// @Success 200 {object} dto.SuccessResponse{data=dto.AuthResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Call service
	authResp, err := h.authService.Login(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to login"

		if errors.Is(err, service.ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    authResp,
	})
}

// GetProfile handles get user profile (requires auth)
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from JWT (set by auth middleware)
	userID := middleware.GetUserID(c)

	// Call service
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to get profile"

		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// UpdateProfile handles update user profile (requires auth)
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.UserUpdateRequest true "Profile update data"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT
	userID := middleware.GetUserID(c)

	var req dto.UserUpdateRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Call service
	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to update profile"

		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
			message = err.Error()
		} else if errors.Is(err, service.ErrUserExists) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    user,
	})
}
