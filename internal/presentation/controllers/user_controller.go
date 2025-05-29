package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// UserController handles user-related HTTP requests
type UserController struct {
	mediator mediator.Mediator
	logger   logger.Logger
}

// NewUserController creates a new UserController
func NewUserController(mediator mediator.Mediator, logger logger.Logger) *UserController {
	return &UserController{
		mediator: mediator,
		logger:   logger,
	}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body commands.RegisterUserCommand true "User registration data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /api/v1/users/register [post]
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var cmd commands.RegisterUserCommand
	
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		c.logger.WithContext(ctx).Errorf("Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
	})
}

// GetUser handles getting a user by ID
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	query := &queries.GetUserByIDQuery{UserID: userID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	user := result.(*entities.User)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// ListUsers handles listing users with filtering
// @Summary List users
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param role query string false "User role filter"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field"
// @Param sort_desc query bool false "Sort descending"
// @Success 200 {object} responses.UsersListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	role := ctx.Query("role")
	search := ctx.Query("search")
	sortBy := ctx.Query("sort_by")
	sortDesc, _ := strconv.ParseBool(ctx.Query("sort_desc"))
	
	// Build filter
	filter := interfaces.UserFilter{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
	
	if role != "" {
		filter.Role = entities.UserRole(role)
	}
	
	query := &queries.ListUsersQuery{Filter: filter}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	users := result.([]*entities.User)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(users),
		},
	})
}

// UpdateUserProfile handles updating user profile
// @Summary Update user profile
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body commands.UpdateUserProfileCommand true "User profile data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUserProfile(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	var cmd commands.UpdateUserProfileCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.UserID = userID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile updated successfully",
	})
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	cmd := &commands.DeleteUserCommand{UserID: userID}
	
	if err := c.mediator.Send(ctx, cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// GetUserAddresses handles getting user addresses
// @Summary Get user addresses
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.AddressesResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id}/addresses [get]
func (c *UserController) GetUserAddresses(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	query := &queries.GetUserAddressesQuery{UserID: userID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	addresses := result.([]*entities.Address)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    addresses,
	})
}

// AddAddress handles adding a user address
// @Summary Add user address
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param address body commands.AddAddressCommand true "Address data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id}/addresses [post]
func (c *UserController) AddAddress(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	var cmd commands.AddAddressCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.UserID = userID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Address added successfully",
	})
}

// UpdateAddress handles updating a user address
// @Summary Update user address
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param address_id path string true "Address ID"
// @Param address body commands.UpdateAddressCommand true "Address data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id}/addresses/{address_id} [put]
func (c *UserController) UpdateAddress(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	addressIDStr := ctx.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid address ID format",
		})
		return
	}
	
	var cmd commands.UpdateAddressCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.UserID = userID
	cmd.AddressID = addressID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address updated successfully",
	})
}

// DeleteAddress handles deleting a user address
// @Summary Delete user address
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Param address_id path string true "Address ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{id}/addresses/{address_id} [delete]
func (c *UserController) DeleteAddress(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	addressIDStr := ctx.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid address ID format",
		})
		return
	}
	
	cmd := &commands.DeleteAddressCommand{
		UserID:    userID,
		AddressID: addressID,
	}
	
	if err := c.mediator.Send(ctx, cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address deleted successfully",
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (c *UserController) handleError(ctx *gin.Context, err error) {
	if appErr, ok := errors.GetAppError(err); ok {
		ctx.JSON(appErr.HTTPStatus, gin.H{
			"success": false,
			"error":   appErr.Message,
			"code":    appErr.Code,
			"details": appErr.Details,
		})
		return
	}
	
	// Generic error
	c.logger.WithContext(ctx).Errorf("Unhandled error: %v", err)
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   "An internal server error occurred",
	})
}
