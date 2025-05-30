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

// ProductController handles product-related HTTP requests
type ProductController struct {
	mediator mediator.Mediator
	logger   logger.Logger
}

// NewProductController creates a new ProductController
func NewProductController(mediator mediator.Mediator, logger logger.Logger) *ProductController {
	return &ProductController{
		mediator: mediator,
		logger:   logger,
	}
}

// CreateProduct handles product creation
// @Summary Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body commands.CreateProductCommand true "Product data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /api/v1/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var cmd commands.CreateProductCommand
	
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
		"message": "Product created successfully",
	})
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} responses.ProductResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	query := &queries.GetProductByIDQuery{ProductID: productID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	product := result.(*entities.Product)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// GetProductBySKU handles getting a product by SKU
// @Summary Get product by SKU
// @Tags Products
// @Produce json
// @Param sku path string true "Product SKU"
// @Success 200 {object} responses.ProductResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/products/sku/{sku} [get]
func (c *ProductController) GetProductBySKU(ctx *gin.Context) {
	sku := ctx.Param("sku")
	if sku == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "SKU is required",
		})
		return
	}
	
	query := &queries.GetProductBySKUQuery{SKU: sku}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	product := result.(*entities.Product)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// ListProducts handles listing products with filtering
// @Summary List products
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param category_id query string false "Category ID filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param in_stock query bool false "In stock filter"
// @Param is_active query bool false "Active filter"
// @Param is_featured query bool false "Featured filter"
// @Param brand query string false "Brand filter"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field"
// @Param sort_desc query bool false "Sort descending"
// @Success 200 {object} responses.ProductsListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/products [get]
func (c *ProductController) ListProducts(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	search := ctx.Query("search")
	brand := ctx.Query("brand")
	sortBy := ctx.Query("sort_by")
	sortDesc, _ := strconv.ParseBool(ctx.Query("sort_desc"))
	
	// Build filter
	filter := interfaces.ProductFilter{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Brand:    brand,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
	
	// Parse optional filters
	if categoryIDStr := ctx.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filter.CategoryID = &categoryID
		}
	}
	
	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}
	
	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}
	
	if inStockStr := ctx.Query("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			filter.InStock = &inStock
		}
	}
	
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}
	
	if isFeaturedStr := ctx.Query("is_featured"); isFeaturedStr != "" {
		if isFeatured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			filter.IsFeatured = &isFeatured
		}
	}
	
	query := &queries.ListProductsQuery{Filter: filter}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	products := result.([]*entities.Product)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(products),
		},
	})
}

// SearchProducts handles product search
// @Summary Search products
// @Tags Products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} responses.ProductsListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/products/search [get]
func (c *ProductController) SearchProducts(ctx *gin.Context) {
	searchQuery := ctx.Query("q")
	if searchQuery == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Search query is required",
		})
		return
	}
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	filter := interfaces.ProductFilter{
		Page:     page,
		PageSize: pageSize,
	}
	
	query := &queries.SearchProductsQuery{
		Query:  searchQuery,
		Filter: filter,
	}
	
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	products := result.([]*entities.Product)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(products),
		},
	})
}

// UpdateProduct handles product updates
// @Summary Update product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body commands.UpdateProductCommand true "Product data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	var cmd commands.UpdateProductCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.ProductID = productID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product updated successfully",
	})
}

// UpdateProductStock handles product stock updates
// @Summary Update product stock
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param stock body commands.UpdateProductStockCommand true "Stock data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/products/{id}/stock [put]
func (c *ProductController) UpdateProductStock(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	var cmd commands.UpdateProductStockCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.ProductID = productID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product stock updated successfully",
	})
}

// DeleteProduct handles product deletion
// @Summary Delete product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	cmd := &commands.DeleteProductCommand{ProductID: productID}
	
	if err := c.mediator.Send(ctx, cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product deleted successfully",
	})
}

// GetLowStockProducts handles getting low stock products
// @Summary Get low stock products
// @Tags Products
// @Produce json
// @Param threshold query int false "Stock threshold" default(5)
// @Success 200 {object} responses.ProductsListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/products/low-stock [get]
func (c *ProductController) GetLowStockProducts(ctx *gin.Context) {
	threshold, _ := strconv.Atoi(ctx.DefaultQuery("threshold", "5"))
	
	query := &queries.GetLowStockProductsQuery{Threshold: threshold}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	products := result.([]*entities.Product)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
		"total":   len(products),
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (c *ProductController) handleError(ctx *gin.Context, err error) {
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
