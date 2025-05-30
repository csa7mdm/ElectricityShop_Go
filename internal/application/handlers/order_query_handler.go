package handlers

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// OrderQueryHandler handles order-related queries
type OrderQueryHandler struct {
	orderRepo   interfaces.OrderRepository
	paymentRepo interfaces.PaymentRepository
	logger      logger.Logger
}

// NewOrderQueryHandler creates a new OrderQueryHandler
func NewOrderQueryHandler(
	orderRepo interfaces.OrderRepository,
	paymentRepo interfaces.PaymentRepository,
	logger logger.Logger,
) *OrderQueryHandler {
	return &OrderQueryHandler{
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// Handle handles queries
func (h *OrderQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	switch q := query.(type) {
	case *queries.GetOrderByIDQuery:
		return h.handleGetOrderByID(ctx, q)
	case *queries.GetOrderByNumberQuery:
		return h.handleGetOrderByNumber(ctx, q)
	case *queries.GetOrdersByUserIDQuery:
		return h.handleGetOrdersByUserID(ctx, q)
	case *queries.ListOrdersQuery:
		return h.handleListOrders(ctx, q)
	case *queries.GetOrderItemsQuery:
		return h.handleGetOrderItems(ctx, q)
	case *queries.GetOrderPaymentsQuery:
		return h.handleGetOrderPayments(ctx, q)
	case *queries.GetOrderShipmentsQuery:
		return h.handleGetOrderShipments(ctx, q)
	case *queries.GetPaymentByIDQuery:
		return h.handleGetPaymentByID(ctx, q)
	case *queries.GetPaymentByTransactionIDQuery:
		return h.handleGetPaymentByTransactionID(ctx, q)
	case *queries.ListPaymentsQuery:
		return h.handleListPayments(ctx, q)
	case *queries.GetOrderSummaryQuery:
		return h.handleGetOrderSummary(ctx, q)
	case *queries.GetOrdersToProcessQuery:
		return h.handleGetOrdersToProcess(ctx, q)
	default:
		return nil, errors.New("UNSUPPORTED_QUERY", "Unsupported query type", 400)
	}
}

// handleGetOrderByID handles getting an order by ID
func (h *OrderQueryHandler) handleGetOrderByID(ctx context.Context, query *queries.GetOrderByIDQuery) (*entities.Order, error) {
	h.logger.WithContext(ctx).Debugf("Getting order by ID: %s", query.OrderID)
	
	order, err := h.orderRepo.GetByID(ctx, query.OrderID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved order: %s", order.ID)
	return order, nil
}

// handleGetOrderByNumber handles getting an order by order number
func (h *OrderQueryHandler) handleGetOrderByNumber(ctx context.Context, query *queries.GetOrderByNumberQuery) (*entities.Order, error) {
	h.logger.WithContext(ctx).Debugf("Getting order by number: %s", query.OrderNumber)
	
	order, err := h.orderRepo.GetByOrderNumber(ctx, query.OrderNumber)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved order: %s", order.ID)
	return order, nil
}

// handleGetOrdersByUserID handles getting orders for a user
func (h *OrderQueryHandler) handleGetOrdersByUserID(ctx context.Context, query *queries.GetOrdersByUserIDQuery) ([]*entities.Order, error) {
	h.logger.WithContext(ctx).Debugf("Getting orders for user: %s", query.UserID)
	
	orders, err := h.orderRepo.GetByUserID(ctx, query.UserID, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d orders for user: %s", len(orders), query.UserID)
	return orders, nil
}

// handleListOrders handles listing orders with filtering
func (h *OrderQueryHandler) handleListOrders(ctx context.Context, query *queries.ListOrdersQuery) ([]*entities.Order, error) {
	h.logger.WithContext(ctx).Debugf("Listing orders with filter")
	
	orders, err := h.orderRepo.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d orders", len(orders))
	return orders, nil
}

// handleGetOrderItems handles getting order items
func (h *OrderQueryHandler) handleGetOrderItems(ctx context.Context, query *queries.GetOrderItemsQuery) ([]*entities.OrderItem, error) {
	h.logger.WithContext(ctx).Debugf("Getting items for order: %s", query.OrderID)
	
	// Get order first to ensure it exists and get items
	order, err := h.orderRepo.GetByID(ctx, query.OrderID)
	if err != nil {
		return nil, err
	}
	
	// Convert slice of OrderItem to slice of *OrderItem
	items := make([]*entities.OrderItem, len(order.Items))
	for i := range order.Items {
		items[i] = &order.Items[i]
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d items for order: %s", len(items), query.OrderID)
	return items, nil
}

// handleGetOrderPayments handles getting order payments
func (h *OrderQueryHandler) handleGetOrderPayments(ctx context.Context, query *queries.GetOrderPaymentsQuery) ([]*entities.Payment, error) {
	h.logger.WithContext(ctx).Debugf("Getting payments for order: %s", query.OrderID)
	
	payments, err := h.paymentRepo.GetByOrderID(ctx, query.OrderID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d payments for order: %s", len(payments), query.OrderID)
	return payments, nil
}

// handleGetOrderShipments handles getting order shipments
func (h *OrderQueryHandler) handleGetOrderShipments(ctx context.Context, query *queries.GetOrderShipmentsQuery) ([]*entities.Shipment, error) {
	h.logger.WithContext(ctx).Debugf("Getting shipments for order: %s", query.OrderID)
	
	// Get order first to ensure it exists and get shipments
	order, err := h.orderRepo.GetByID(ctx, query.OrderID)
	if err != nil {
		return nil, err
	}
	
	// Convert slice of Shipment to slice of *Shipment
	shipments := make([]*entities.Shipment, len(order.Shipments))
	for i := range order.Shipments {
		shipments[i] = &order.Shipments[i]
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d shipments for order: %s", len(shipments), query.OrderID)
	return shipments, nil
}

// handleGetPaymentByID handles getting a payment by ID
func (h *OrderQueryHandler) handleGetPaymentByID(ctx context.Context, query *queries.GetPaymentByIDQuery) (*entities.Payment, error) {
	h.logger.WithContext(ctx).Debugf("Getting payment by ID: %s", query.PaymentID)
	
	payment, err := h.paymentRepo.GetByID(ctx, query.PaymentID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved payment: %s", payment.ID)
	return payment, nil
}

// handleGetPaymentByTransactionID handles getting a payment by transaction ID
func (h *OrderQueryHandler) handleGetPaymentByTransactionID(ctx context.Context, query *queries.GetPaymentByTransactionIDQuery) (*entities.Payment, error) {
	h.logger.WithContext(ctx).Debugf("Getting payment by transaction ID: %s", query.TransactionID)
	
	payment, err := h.paymentRepo.GetByTransactionID(ctx, query.TransactionID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved payment: %s", payment.ID)
	return payment, nil
}

// handleListPayments handles listing payments with filtering
func (h *OrderQueryHandler) handleListPayments(ctx context.Context, query *queries.ListPaymentsQuery) ([]*entities.Payment, error) {
	h.logger.WithContext(ctx).Debugf("Listing payments with filter")
	
	payments, err := h.paymentRepo.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d payments", len(payments))
	return payments, nil
}

// OrderSummary represents order summary statistics
type OrderSummary struct {
	TotalOrders       int             `json:"total_orders"`
	TotalRevenue      decimal.Decimal `json:"total_revenue"`
	PendingOrders     int             `json:"pending_orders"`
	ProcessingOrders  int             `json:"processing_orders"`
	CompletedOrders   int             `json:"completed_orders"`
	CancelledOrders   int             `json:"cancelled_orders"`
	AverageOrderValue decimal.Decimal `json:"average_order_value"`
	TopSellingProducts []ProductSales `json:"top_selling_products,omitempty"`
}

// ProductSales represents product sales data
type ProductSales struct {
	ProductID   string          `json:"product_id"`
	ProductName string          `json:"product_name"`
	QuantitySold int            `json:"quantity_sold"`
	Revenue     decimal.Decimal `json:"revenue"`
}

// handleGetOrderSummary handles getting order summary/statistics
func (h *OrderQueryHandler) handleGetOrderSummary(ctx context.Context, query *queries.GetOrderSummaryQuery) (*OrderSummary, error) {
	h.logger.WithContext(ctx).Debugf("Getting order summary")
	
	// Build filter for orders
	filter := interfaces.OrderFilter{
		Page:     1,
		PageSize: 1000, // Get a large number for summary calculation
	}
	
	if query.UserID != nil {
		filter.UserID = query.UserID
	}
	
	if query.StartDate != nil {
		filter.StartDate = query.StartDate
	}
	
	if query.EndDate != nil {
		filter.EndDate = query.EndDate
	}
	
	// Get orders
	orders, err := h.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	// Calculate summary statistics
	summary := &OrderSummary{
		TotalOrders:      len(orders),
		TotalRevenue:     decimal.Zero,
		AverageOrderValue: decimal.Zero,
	}
	
	// Count orders by status and calculate revenue
	for _, order := range orders {
		summary.TotalRevenue = summary.TotalRevenue.Add(order.Total)
		
		switch order.Status {
		case entities.OrderStatusPending:
			summary.PendingOrders++
		case entities.OrderStatusProcessing:
			summary.ProcessingOrders++
		case entities.OrderStatusDelivered:
			summary.CompletedOrders++
		case entities.OrderStatusCancelled:
			summary.CancelledOrders++
		}
	}
	
	// Calculate average order value
	if summary.TotalOrders > 0 {
		summary.AverageOrderValue = summary.TotalRevenue.Div(decimal.NewFromInt(int64(summary.TotalOrders)))
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully calculated order summary")
	return summary, nil
}

// handleGetOrdersToProcess handles getting orders that need processing
func (h *OrderQueryHandler) handleGetOrdersToProcess(ctx context.Context, query *queries.GetOrdersToProcessQuery) ([]*entities.Order, error) {
	h.logger.WithContext(ctx).Debugf("Getting orders to process")
	
	orders, err := h.orderRepo.GetOrdersToProcess(ctx)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d orders to process", len(orders))
	return orders, nil
}
