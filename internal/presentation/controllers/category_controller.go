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

// CategoryController handles category-related HTTP requests
type CategoryController struct {
	mediator mediator.Mediator
	logger   logger.Logger
}

// NewCategoryController creates a new CategoryController
func NewCategoryController(mediator mediator.Mediator, logger logger.Logger) *CategoryController {
	return &CategoryController{
		mediator: mediator,
		logger:   logger,
	}
}

// CreateCategory handles category creation
// @Summary Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body commands.CreateCategoryCommand true "Category data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /api/v1/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var cmd commands.CreateCategoryCommand
	
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
		"message": "Category created successfully",
	})
}

// GetCategory handles getting a category by ID
// @Summary Get category by ID
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} responses.CategoryResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/categories/{id} [get]
func (c *CategoryController) GetCategory(ctx *gin.Context) {
	categoryIDStr := ctx.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid category ID format",
		})
		return
	}
	
	query := &queries.GetCategoryByIDQuery{CategoryID: categoryID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	category := result.(*entities.Category)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})
}

// GetCategoryBySlug handles getting a category by slug
// @Summary Get category by slug
// @Tags Categories
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} responses.CategoryResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/categories/slug/{slug} [get]
func (c *CategoryController) GetCategoryBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Slug is required",
		})
		return
	}
	
	query := &queries.GetCategoryBySlugQuery{Slug: slug}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	category := result.(*entities.Category)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})
}

// ListCategories handles listing categories with filtering
// @Summary List categories
// @Tags Categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param parent_id query string false "Parent category ID filter"
// @Param is_active query bool false "Active filter"
// @Param sort_by query string false "Sort field"
// @Param sort_desc query bool false "Sort descending"
// @Success 200 {object} responses.CategoriesListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/categories [get]
func (c *CategoryController) ListCategories(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	sortBy := ctx.Query("sort_by")
	sortDesc, _ := strconv.ParseBool(ctx.Query("sort_desc"))
	
	// Build filter
	filter := interfaces.CategoryFilter{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
	
	// Parse optional filters
	if parentIDStr := ctx.Query("parent_id"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			filter.ParentID = &parentID
		}
	}
	
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}
	
	query := &queries.ListCategoriesQuery{Filter: filter}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	categories := result.([]*entities.Category)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(categories),
		},
	})
}

// GetRootCategories handles getting root categories
// @Summary Get root categories
// @Tags Categories
// @Produce json
// @Success 200 {object} responses.CategoriesListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/categories/root [get]
func (c *CategoryController) GetRootCategories(ctx *gin.Context) {
	query := &queries.GetRootCategoriesQuery{}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	categories := result.([]*entities.Category)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"total":   len(categories),
	})
}

// GetCategoryChildren handles getting category children
// @Summary Get category children
// @Tags Categories
// @Produce json
// @Param id path string true "Parent category ID"
// @Success 200 {object} responses.CategoriesListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/categories/{id}/children [get]
func (c *CategoryController) GetCategoryChildren(ctx *gin.Context) {
	parentIDStr := ctx.Param("id")
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid parent category ID format",
		})
		return
	}
	
	query := &queries.GetCategoryChildrenQuery{ParentID: parentID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	categories := result.([]*entities.Category)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"total":   len(categories),
	})
}

// UpdateCategory handles category updates
// @Summary Update category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body commands.UpdateCategoryCommand true "Category data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	categoryIDStr := ctx.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid category ID format",
		})
		return
	}
	
	var cmd commands.UpdateCategoryCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.CategoryID = categoryID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
	})
}

// DeleteCategory handles category deletion
// @Summary Delete category
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	categoryIDStr := ctx.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid category ID format",
		})
		return
	}
	
	cmd := &commands.DeleteCategoryCommand{CategoryID: categoryID}
	
	if err := c.mediator.Send(ctx, cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (c *CategoryController) handleError(ctx *gin.Context, err error) {
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
