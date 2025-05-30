package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/handlers"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// CartController handles cart-related HTTP requests
type CartController struct {
	mediator mediator.Mediator
	logger   logger.Logger
}

// NewCartController creates a new CartController
func NewCartController(mediator mediator.Mediator, logger logger.Logger) *CartController {
	return &CartController{
		mediator: mediator,
		logger:   logger,
	}
}

// GetCart handles getting user's cart
// @Summary Get user's cart
// @Tags Cart
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} responses.CartResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart [get]
func (c *CartController) GetCart(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	query := &queries.GetCartByUserIDQuery{UserID: userID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	cart := result.(*entities.Cart)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cart,
	})
}

// GetCartSummary handles getting cart summary
// @Summary Get cart summary
// @Tags Cart
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} responses.CartSummaryResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart/summary [get]
func (c *CartController) GetCartSummary(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	query := &queries.GetCartSummaryQuery{UserID: userID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	summary := result.(*handlers.CartSummary)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// AddToCart handles adding an item to cart
// @Summary Add item to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param item body commands.AddToCartCommand true "Cart item data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart/items [post]
func (c *CartController) AddToCart(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	var cmd commands.AddToCartCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		c.logger.WithContext(ctx).Errorf("Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.UserID = userID // Ensure user ID from URL is used
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Item added to cart successfully",
	})
}

// UpdateCartItem handles updating cart item quantity
// @Summary Update cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param product_id path string true "Product ID"
// @Param item body commands.UpdateCartItemCommand true "Updated cart item data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart/items/{product_id} [put]
func (c *CartController) UpdateCartItem(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	var cmd commands.UpdateCartItemCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.UserID = userID       // Ensure user ID from URL is used
	cmd.ProductID = productID // Ensure product ID from URL is used
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart item updated successfully",
	})
}

// RemoveFromCart handles removing an item from cart
// @Summary Remove item from cart
// @Tags Cart
// @Produce json
// @Param user_id path string true "User ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart/items/{product_id} [delete]
func (c *CartController) RemoveFromCart(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
		})
		return
	}
	
	cmd := &commands.RemoveFromCartCommand{
		UserID:    userID,
		ProductID: productID,
	}
	
	if err := c.mediator.Send(ctx, cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item removed from cart successfully",
	})
}

// ClearCart handles clearing all items from cart
// @Summary Clear all items from cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param clear body commands.ClearCartCommand true "Clear cart data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/cart/clear [post]
func (c *CartController) ClearCart(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	var cmd commands.ClearCartCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		// If no body provided, use default
		cmd = commands.ClearCartCommand{
			UserID: userID,
			Reason: "User action",
		}
	} else {
		cmd.UserID = userID // Ensure user ID from URL is used
	}
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart cleared successfully",
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (c *CartController) handleError(ctx *gin.Context, err error) {
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
