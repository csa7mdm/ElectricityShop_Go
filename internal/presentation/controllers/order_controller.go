package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/handlers"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// OrderController handles order-related HTTP requests
type OrderController struct {
	mediator mediator.Mediator
	logger   logger.Logger
}

// NewOrderController creates a new OrderController
func NewOrderController(mediator mediator.Mediator, logger logger.Logger) *OrderController {
	return &OrderController{
		mediator: mediator,
		logger:   logger,
	}
}

// CreateOrder handles order creation
// @Summary Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body commands.CreateOrderCommand true "Order data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var cmd commands.CreateOrderCommand
	
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
		"message": "Order created successfully",
	})
}

// CreateOrderFromCart handles creating order from cart
// @Summary Create order from cart
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body commands.CreateOrderFromCartCommand true "Order from cart data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/orders/from-cart [post]
func (c *OrderController) CreateOrderFromCart(ctx *gin.Context) {
	var cmd commands.CreateOrderFromCartCommand
	
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
		"message": "Order created from cart successfully",
	})
}

// GetOrder handles getting an order by ID
// @Summary Get order by ID
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} responses.OrderResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id} [get]
func (c *OrderController) GetOrder(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid order ID format",
		})
		return
	}
	
	query := &queries.GetOrderByIDQuery{OrderID: orderID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	order := result.(*entities.Order)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetOrderByNumber handles getting an order by order number
// @Summary Get order by order number
// @Tags Orders
// @Produce json
// @Param number path string true "Order number"
// @Success 200 {object} responses.OrderResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/number/{number} [get]
func (c *OrderController) GetOrderByNumber(ctx *gin.Context) {
	orderNumber := ctx.Param("number")
	if orderNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Order number is required",
		})
		return
	}
	
	query := &queries.GetOrderByNumberQuery{OrderNumber: orderNumber}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	order := result.(*entities.Order)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetUserOrders handles getting orders for a user
// @Summary Get user orders
// @Tags Orders
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Order status filter"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} responses.OrdersListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/users/{user_id}/orders [get]
func (c *OrderController) GetUserOrders(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}
	
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status := ctx.Query("status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	
	// Build filter
	filter := interfaces.OrderFilter{
		Page:     page,
		PageSize: pageSize,
	}
	
	if status != "" {
		filter.Status = entities.OrderStatus(status)
	}
	
	if startDate != "" {
		filter.StartDate = &startDate
	}
	
	if endDate != "" {
		filter.EndDate = &endDate
	}
	
	query := &queries.GetOrdersByUserIDQuery{
		UserID: userID,
		Filter: filter,
	}
	
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	orders := result.([]*entities.Order)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(orders),
		},
	})
}

// ListOrders handles listing all orders with filtering
// @Summary List orders
// @Tags Orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Order status filter"
// @Param payment_status query string false "Payment status filter"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} responses.OrdersListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/orders [get]
func (c *OrderController) ListOrders(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status := ctx.Query("status")
	paymentStatus := ctx.Query("payment_status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	
	// Build filter
	filter := interfaces.OrderFilter{
		Page:     page,
		PageSize: pageSize,
	}
	
	if status != "" {
		filter.Status = entities.OrderStatus(status)
	}
	
	if paymentStatus != "" {
		filter.PaymentStatus = entities.PaymentStatus(paymentStatus)
	}
	
	if startDate != "" {
		filter.StartDate = &startDate
	}
	
	if endDate != "" {
		filter.EndDate = &endDate
	}
	
	query := &queries.ListOrdersQuery{Filter: filter}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	orders := result.([]*entities.Order)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     len(orders),
		},
	})
}

// UpdateOrderStatus handles updating order status
// @Summary Update order status
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body commands.UpdateOrderStatusCommand true "Status update data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id}/status [put]
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid order ID format",
		})
		return
	}
	
	var cmd commands.UpdateOrderStatusCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.OrderID = orderID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order status updated successfully",
	})
}

// CancelOrder handles order cancellation
// @Summary Cancel order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param cancel body commands.CancelOrderCommand true "Cancel order data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id}/cancel [post]
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid order ID format",
		})
		return
	}
	
	var cmd commands.CancelOrderCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.OrderID = orderID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order cancelled successfully",
	})
}

// ProcessPayment handles payment processing
// @Summary Process payment for order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param payment body commands.ProcessPaymentCommand true "Payment data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id}/payment [post]
func (c *OrderController) ProcessPayment(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid order ID format",
		})
		return
	}
	
	var cmd commands.ProcessPaymentCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	cmd.OrderID = orderID
	
	if err := c.mediator.Send(ctx, &cmd); err != nil {
		c.handleError(ctx, err)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment processed successfully",
	})
}

// GetOrderPayments handles getting order payments
// @Summary Get order payments
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} responses.PaymentsResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /api/v1/orders/{id}/payments [get]
func (c *OrderController) GetOrderPayments(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid order ID format",
		})
		return
	}
	
	query := &queries.GetOrderPaymentsQuery{OrderID: orderID}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	payments := result.([]*entities.Payment)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payments,
	})
}

// GetOrderSummary handles getting order summary/statistics
// @Summary Get order summary
// @Tags Orders
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} responses.OrderSummaryResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/orders/summary [get]
func (c *OrderController) GetOrderSummary(ctx *gin.Context) {
	query := &queries.GetOrderSummaryQuery{}
	
	// Parse optional filters
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			query.UserID = &userID
		}
	}
	
	if startDate := ctx.Query("start_date"); startDate != "" {
		query.StartDate = &startDate
	}
	
	if endDate := ctx.Query("end_date"); endDate != "" {
		query.EndDate = &endDate
	}
	
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	summary := result.(*handlers.OrderSummary)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetOrdersToProcess handles getting orders that need processing
// @Summary Get orders to process
// @Tags Orders
// @Produce json
// @Success 200 {object} responses.OrdersListResponse
// @Failure 400 {object} responses.ErrorResponse
// @Router /api/v1/orders/to-process [get]
func (c *OrderController) GetOrdersToProcess(ctx *gin.Context) {
	query := &queries.GetOrdersToProcessQuery{}
	result, err := c.mediator.Query(ctx, query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	
	orders := result.([]*entities.Order)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"total":   len(orders),
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (c *OrderController) handleError(ctx *gin.Context, err error) {
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
