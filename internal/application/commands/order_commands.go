package commands

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
)

// CreateOrderCommand represents creating an order
type CreateOrderCommand struct {
	UserID             uuid.UUID                  `json:"user_id" validate:"required"`
	Items              []CreateOrderItemCommand   `json:"items" validate:"required,min=1"`
	ShippingAddressID  uuid.UUID                  `json:"shipping_address_id" validate:"required"`
	BillingAddressID   uuid.UUID                  `json:"billing_address_id" validate:"required"`
	PaymentMethod      entities.PaymentMethod     `json:"payment_method" validate:"required"`
	Notes              string                     `json:"notes,omitempty"`
}

func (c CreateOrderCommand) GetName() string {
	return "CreateOrder"
}

// CreateOrderItemCommand represents an order item in create order command
type CreateOrderItemCommand struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderStatusCommand represents updating order status
type UpdateOrderStatusCommand struct {
	OrderID   uuid.UUID           `json:"order_id" validate:"required"`
	Status    entities.OrderStatus `json:"status" validate:"required"`
	Reason    string              `json:"reason,omitempty"`
	UpdatedBy uuid.UUID           `json:"updated_by" validate:"required"`
}

func (c UpdateOrderStatusCommand) GetName() string {
	return "UpdateOrderStatus"
}

// CancelOrderCommand represents cancelling an order
type CancelOrderCommand struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	UserID        uuid.UUID `json:"user_id" validate:"required"`
	CancelReason  string    `json:"cancel_reason" validate:"required"`
}

func (c CancelOrderCommand) GetName() string {
	return "CancelOrder"
}

// ProcessPaymentCommand represents processing a payment
type ProcessPaymentCommand struct {
	OrderID           uuid.UUID              `json:"order_id" validate:"required"`
	Amount            decimal.Decimal        `json:"amount" validate:"required"`
	PaymentMethod     entities.PaymentMethod `json:"payment_method" validate:"required"`
	TransactionID     string                 `json:"transaction_id,omitempty"`
	GatewayResponse   string                 `json:"gateway_response,omitempty"`
}

func (c ProcessPaymentCommand) GetName() string {
	return "ProcessPayment"
}

// UpdatePaymentStatusCommand represents updating payment status
type UpdatePaymentStatusCommand struct {
	PaymentID       uuid.UUID              `json:"payment_id" validate:"required"`
	Status          entities.PaymentStatus `json:"status" validate:"required"`
	TransactionID   string                 `json:"transaction_id,omitempty"`
	GatewayResponse string                 `json:"gateway_response,omitempty"`
	FailureReason   string                 `json:"failure_reason,omitempty"`
}

func (c UpdatePaymentStatusCommand) GetName() string {
	return "UpdatePaymentStatus"
}

// CreateShipmentCommand represents creating a shipment
type CreateShipmentCommand struct {
	OrderID        uuid.UUID              `json:"order_id" validate:"required"`
	TrackingNumber string                 `json:"tracking_number" validate:"required"`
	Carrier        string                 `json:"carrier" validate:"required"`
	EstimatedDelivery *string             `json:"estimated_delivery,omitempty"`
}

func (c CreateShipmentCommand) GetName() string {
	return "CreateShipment"
}

// UpdateShipmentStatusCommand represents updating shipment status
type UpdateShipmentStatusCommand struct {
	ShipmentID uuid.UUID                `json:"shipment_id" validate:"required"`
	Status     entities.ShippingStatus  `json:"status" validate:"required"`
}

func (c UpdateShipmentStatusCommand) GetName() string {
	return "UpdateShipmentStatus"
}
