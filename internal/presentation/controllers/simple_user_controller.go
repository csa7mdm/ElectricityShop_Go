package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/application/commands" // Added import
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities" // Keep for Login response for now
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses"
	// "github.com/yourusername/electricity-shop-go/pkg/auth" // Removed authService dependency
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// SimpleUserController handles user-related HTTP requests with direct dependencies
type SimpleUserController struct {
	mediator  mediator.Mediator // Changed from authService
	logger    logger.Logger
	validator *validator.Validate
}

// NewSimpleUserController creates a new SimpleUserController
func NewSimpleUserController(mediator mediator.Mediator, logger logger.Logger) *SimpleUserController { // Changed signature
	return &SimpleUserController{
		mediator:  mediator, // Changed from authService
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterUser handles user registration
func (uc *SimpleUserController) RegisterUser(c *gin.Context) {
	var req dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user registration: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", errors.ErrorCodeInvalidInput))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for user registration: %v", err)
		// It's good practice to provide more specific validation error details if possible
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed: "+err.Error(), errors.ErrorCodeValidation))
		return
	}

	uc.logger.Infof("Registration request received for email: %s", req.Email)

	cmd := commands.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	_, err := uc.mediator.Send(c.Request.Context(), &cmd)
	if err != nil {
		uc.logger.Errorf("Failed to register user: %v", err)
		apiErr := errors.ExtractAPIError(err)
		if apiErr != nil {
			if apiErr.Code == errors.ErrorCodeUserAlreadyExists {
				c.JSON(http.StatusConflict, responses.NewErrorResponse(apiErr.Message, apiErr.Code))
				return
			}
			c.JSON(apiErr.StatusCode, responses.NewErrorResponse(apiErr.Message, apiErr.Code))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to register user", errors.ErrorCodeInternalServer))
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccessResponse(nil, "User registered successfully"))
}

// Login handles user authentication
func (uc *SimpleUserController) Login(c *gin.Context) {
	var req dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user login: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", errors.ErrorCodeInvalidInput))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for user login: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed: "+err.Error(), errors.ErrorCodeValidation))
		return
	}

	uc.logger.Infof("Login request received for email: %s", req.Email)

	cmd := commands.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := uc.mediator.Send(c.Request.Context(), &cmd)
	if err != nil {
		uc.logger.Errorf("Failed to login user: %v", err)
		apiErr := errors.ExtractAPIError(err)
		if apiErr != nil {
			if apiErr.Code == errors.ErrorCodeInvalidCredentials {
				c.JSON(http.StatusUnauthorized, responses.NewErrorResponse(apiErr.Message, apiErr.Code))
				return
			}
			c.JSON(apiErr.StatusCode, responses.NewErrorResponse(apiErr.Message, apiErr.Code))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to login", errors.ErrorCodeInternalServer))
		return
	}

	loginResponse, ok := result.(*dtos.LoginUserResponse)
	if !ok {
		uc.logger.Errorf("Login command returned unexpected type: %T", result)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to process login response", errors.ErrorCodeInternalServer))
		return
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
