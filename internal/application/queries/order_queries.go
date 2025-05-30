package queries

import (
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
)

// GetOrderByIDQuery represents a query to get an order by ID
type GetOrderByIDQuery struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

func (q GetOrderByIDQuery) GetName() string {
	return "GetOrderByID"
}

// GetOrderByNumberQuery represents a query to get an order by order number
type GetOrderByNumberQuery struct {
	OrderNumber string `json:"order_number" validate:"required"`
}

func (q GetOrderByNumberQuery) GetName() string {
	return "GetOrderByNumber"
}

// GetOrdersByUserIDQuery represents a query to get orders for a user
type GetOrdersByUserIDQuery struct {
	UserID uuid.UUID              `json:"user_id" validate:"required"`
	Filter interfaces.OrderFilter `json:"filter"`
}

func (q GetOrdersByUserIDQuery) GetName() string {
	return "GetOrdersByUserID"
}

// ListOrdersQuery represents a query to list orders with filtering
type ListOrdersQuery struct {
	Filter interfaces.OrderFilter `json:"filter"`
}

func (q ListOrdersQuery) GetName() string {
	return "ListOrders"
}

// GetOrderItemsQuery represents a query to get order items
type GetOrderItemsQuery struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

func (q GetOrderItemsQuery) GetName() string {
	return "GetOrderItems"
}

// GetOrderPaymentsQuery represents a query to get order payments
type GetOrderPaymentsQuery struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

func (q GetOrderPaymentsQuery) GetName() string {
	return "GetOrderPayments"
}

// GetOrderShipmentsQuery represents a query to get order shipments
type GetOrderShipmentsQuery struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

func (q GetOrderShipmentsQuery) GetName() string {
	return "GetOrderShipments"
}

// GetPaymentByIDQuery represents a query to get a payment by ID
type GetPaymentByIDQuery struct {
	PaymentID uuid.UUID `json:"payment_id" validate:"required"`
}

func (q GetPaymentByIDQuery) GetName() string {
	return "GetPaymentByID"
}

// GetPaymentByTransactionIDQuery represents a query to get a payment by transaction ID
type GetPaymentByTransactionIDQuery struct {
	TransactionID string `json:"transaction_id" validate:"required"`
}

func (q GetPaymentByTransactionIDQuery) GetName() string {
	return "GetPaymentByTransactionID"
}

// ListPaymentsQuery represents a query to list payments with filtering
type ListPaymentsQuery struct {
	Filter interfaces.PaymentFilter `json:"filter"`
}

func (q ListPaymentsQuery) GetName() string {
	return "ListPayments"
}

// GetOrderSummaryQuery represents a query to get order summary/statistics
type GetOrderSummaryQuery struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	StartDate *string    `json:"start_date,omitempty"`
	EndDate   *string    `json:"end_date,omitempty"`
}

func (q GetOrderSummaryQuery) GetName() string {
	return "GetOrderSummary"
}

// GetOrdersToProcessQuery represents a query to get orders that need processing
type GetOrdersToProcessQuery struct{}

func (q GetOrdersToProcessQuery) GetName() string {
	return "GetOrdersToProcess"
}
