package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// UserController handles user-related HTTP requests
type UserController struct {
	mediator  mediator.Mediator
	logger    logger.Logger
	validator *validator.Validate
}

// NewUserController creates a new UserController
func NewUserController(m mediator.Mediator, l logger.Logger) *UserController {
	return &UserController{
		mediator:  m,
		logger:    l,
		validator: validator.New(),
	}
}

// RegisterUser handles user registration
func (uc *UserController) RegisterUser(c *gin.Context) {
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

	// Create command
	cmd := &commands.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("User registration failed: %v", err)
		
		// Map domain errors to HTTP status codes
		switch {
		case errors.IsErrorType(err, "USER_ALREADY_EXISTS"):
			c.JSON(http.StatusConflict, responses.NewErrorResponse("User already exists", "USER_ALREADY_EXISTS"))
		case errors.IsErrorType(err, "PASSWORD_HASH_ERROR"):
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Internal server error", "INTERNAL_ERROR"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Registration failed", "REGISTRATION_FAILED"))
		}
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccessResponse(nil, "User registered successfully"))
}

// Login handles user authentication
func (uc *UserController) Login(c *gin.Context) {
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

	// Create login command (treated as query since it returns data)
	cmd := &commands.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	// Execute query
	result, err := uc.mediator.Query(c.Request.Context(), cmd)
	if err != nil {
		uc.logger.Errorf("User login failed: %v", err)
		
		// Map domain errors to HTTP status codes
		switch {
		case errors.IsErrorType(err, "INVALID_CREDENTIALS"):
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid credentials", "INVALID_CREDENTIALS"))
		case errors.IsErrorType(err, "USER_INACTIVE"):
			c.JSON(http.StatusForbidden, responses.NewErrorResponse("User account is inactive", "USER_INACTIVE"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Login failed", "LOGIN_FAILED"))
		}
		return
	}

	loginResponse, ok := result.(*dtos.LoginUserResponse)
	if !ok {
		uc.logger.Errorf("Login handler returned unexpected type: %T", result)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Internal server error", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(loginResponse, "Login successful"))
}

// GetUser handles getting user by ID
func (uc *UserController) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// Create query
	query := &queries.GetUserByIDQuery{
		UserID: userID,
	}

	// Execute query
	result, err := uc.mediator.Query(c.Request.Context(), query)
	if err != nil {
		uc.logger.Errorf("Failed to get user: %v", err)
		
		switch {
		case errors.IsErrorType(err, "USER_NOT_FOUND"):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("User not found", "USER_NOT_FOUND"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to get user", "GET_USER_FAILED"))
		}
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(result, "User retrieved successfully"))
}

// ListUsers handles getting list of users
func (uc *UserController) ListUsers(c *gin.Context) {
	// Create query with pagination
	query := &queries.ListUsersQuery{
		PageSize: 20, // Default page size
		Page:     1,  // Default page
	}

	// Parse query parameters if provided
	if pageStr := c.Query("page"); pageStr != "" {
		// Parse page number
	}
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		// Parse page size
	}

	// Execute query
	result, err := uc.mediator.Query(c.Request.Context(), query)
	if err != nil {
		uc.logger.Errorf("Failed to list users: %v", err)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to list users", "LIST_USERS_FAILED"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(result, "Users retrieved successfully"))
}

// UpdateUserProfile handles user profile updates
func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// TODO: Add UpdateUserProfileRequest DTO
	// var req dtos.UpdateUserProfileRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
	// 	return
	// }

	// Create command
	cmd := &commands.UpdateUserProfileCommand{
		UserID: userID,
		// Add other fields from request
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("Failed to update user profile: %v", err)
		
		switch {
		case errors.IsErrorType(err, "USER_NOT_FOUND"):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("User not found", "USER_NOT_FOUND"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to update profile", "UPDATE_PROFILE_FAILED"))
		}
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(nil, "Profile updated successfully"))
}

// DeleteUser handles user deletion
func (uc *UserController) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// Create command
	cmd := &commands.DeleteUserCommand{
		UserID: userID,
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("Failed to delete user: %v", err)
		
		switch {
		case errors.IsErrorType(err, "USER_NOT_FOUND"):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("User not found", "USER_NOT_FOUND"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to delete user", "DELETE_USER_FAILED"))
		}
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(nil, "User deleted successfully"))
}

// GetUserAddresses handles getting user addresses
func (uc *UserController) GetUserAddresses(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	// Create query
	query := &queries.GetUserAddressesQuery{
		UserID: userID,
	}

	// Execute query
	result, err := uc.mediator.Query(c.Request.Context(), query)
	if err != nil {
		uc.logger.Errorf("Failed to get user addresses: %v", err)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to get addresses", "GET_ADDRESSES_FAILED"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(result, "Addresses retrieved successfully"))
}

// AddAddress handles adding an address to a user
func (uc *UserController) AddAddress(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	var req dtos.AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for add address: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for add address: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	// Create command
	cmd := &commands.AddAddressCommand{
		UserID:     userID,
		Type:       entities.AddressType(req.Type),
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		IsDefault:  req.IsDefault,
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("Failed to add address: %v", err)
		
		switch {
		case errors.IsErrorType(err, "USER_NOT_FOUND"):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("User not found", "USER_NOT_FOUND"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to add address", "ADD_ADDRESS_FAILED"))
		}
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccessResponse(nil, "Address added successfully"))
}

// UpdateAddress handles updating an address
func (uc *UserController) UpdateAddress(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	addressIDStr := c.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid address ID", "INVALID_ADDRESS_ID"))
		return
	}

	var req dtos.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for update address: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// Validate request
	if err := uc.validator.Struct(&req); err != nil {
		uc.logger.Errorf("Validation failed for update address: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	// Create command
	cmd := &commands.UpdateAddressCommand{
		UserID:     userID,
		AddressID:  addressID,
		Type:       entities.AddressType(req.Type),
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		IsDefault:  req.IsDefault,
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("Failed to update address: %v", err)
		
		switch {
		case errors.IsErrorType(err, "ADDRESS_NOT_FOUND"):
			c.JSON(http.StatusNotFound, responses.NewErrorResponse("Address not found", "ADDRESS_NOT_FOUND"))
		case errors.IsErrorType(err, "UNAUTHORIZED"):
			c.JSON(http.StatusForbidden, responses.NewErrorResponse("Address does not belong to user", "UNAUTHORIZED"))
		default:
			c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to update address", "UPDATE_ADDRESS_FAILED"))
		}
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(nil, "Address updated successfully"))
}

// DeleteAddress handles deleting an address
func (uc *UserController) DeleteAddress(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	addressIDStr := c.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid address ID", "INVALID_ADDRESS_ID"))
		return
	}

	// Create command
	cmd := &commands.DeleteAddressCommand{
		UserID:    userID,
		AddressID: addressID,
	}

	// Execute command
	if err := uc.mediator.Send(c.Request.Context(), cmd); err != nil {
		uc.logger.Errorf("Failed to delete address: %v", err)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to delete address", "DELETE_ADDRESS_FAILED"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(nil, "Address deleted successfully"))
}
