package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// SimpleUserController handles user-related HTTP requests with direct dependencies
type SimpleUserController struct {
	authService *auth.AuthService
	logger      logger.Logger
	validator   *validator.Validate
}

// NewSimpleUserController creates a new SimpleUserController
func NewSimpleUserController(authService *auth.AuthService, logger logger.Logger) *SimpleUserController {
	return &SimpleUserController{
		authService: authService,
		logger:      logger,
		validator:   validator.New(),
	}
}

// RegisterUser handles user registration
func (uc *SimpleUserController) RegisterUser(c *gin.Context) {
	var req dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user registration: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for user registration: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	// For now, return success without actually creating user
	// TODO: Connect to actual user creation logic
	uc.logger.Infof("Registration request received for email: %s", req.Email)
	
	// Simulate user creation
	userID := uuid.New()
	
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(gin.H{
		"id": userID.String(),
		"email": req.Email,
		"message": "User registration successful (demo mode)",
	}, "User registered successfully"))
}

// Login handles user authentication
func (uc *SimpleUserController) Login(c *gin.Context) {
	var req dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user login: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for user login: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	uc.logger.Infof("Login request received for email: %s", req.Email)
	
	// For demo purposes, create a test user and token
	// TODO: Connect to actual user authentication logic
	
	// Demo: Accept any email/password combination for testing
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Email and password required", "INVALID_CREDENTIALS"))
		return
	}
	
	// Generate a real JWT token for testing
	userID := uuid.New()
	token, err := uc.authService.GenerateToken(userID, req.Email, entities.RoleCustomer)
	if err != nil {
		uc.logger.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to generate token", "TOKEN_ERROR"))
		return
	}
	
	loginResponse := &dtos.LoginUserResponse{
		ID:    userID.String(),
		Email: req.Email,
		Role:  "customer",
		Token: token,
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(loginResponse, "Login successful"))
}

// GetUser handles getting user by ID (protected endpoint)
func (uc *SimpleUserController) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// Get user info from JWT context
	contextUserID := c.GetString("user_id")
	contextUserEmail := c.GetString("user_email")
	contextUserRole := c.GetString("user_role")

	// For demo purposes, return the context user info
	userResponse := gin.H{
		"id":    userID.String(),
		"email": contextUserEmail,
		"role":  contextUserRole,
		"context_user_id": contextUserID,
		"message": "User retrieved successfully (demo mode)",
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(userResponse, "User retrieved successfully"))
}

// UpdateUserProfile handles user profile updates (protected endpoint)
func (uc *SimpleUserController) UpdateUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// Get user info from JWT context
	contextUserID := c.GetString("user_id")
	
	// Basic authorization check
	if contextUserID != userID.String() {
		c.JSON(http.StatusForbidden, responses.NewErrorResponse("Cannot update another user's profile", "FORBIDDEN"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(gin.H{
		"id": userID.String(),
		"message": "Profile update successful (demo mode)",
	}, "Profile updated successfully"))
}
