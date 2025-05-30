package commands

import (
	"github.com/google/uuid"
)

// AddToCartCommand represents adding an item to cart
type AddToCartCommand struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

func (c AddToCartCommand) GetName() string {
	return "AddToCart"
}

// UpdateCartItemCommand represents updating cart item quantity
type UpdateCartItemCommand struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

func (c UpdateCartItemCommand) GetName() string {
	return "UpdateCartItem"
}

// RemoveFromCartCommand represents removing an item from cart
type RemoveFromCartCommand struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (c RemoveFromCartCommand) GetName() string {
	return "RemoveFromCart"
}

// ClearCartCommand represents clearing all items from cart
type ClearCartCommand struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Reason string    `json:"reason,omitempty"`
}

func (c ClearCartCommand) GetName() string {
	return "ClearCart"
}

// CreateOrderFromCartCommand represents creating an order from cart items
type CreateOrderFromCartCommand struct {
	UserID             uuid.UUID `json:"user_id" validate:"required"`
	ShippingAddressID  uuid.UUID `json:"shipping_address_id" validate:"required"`
	BillingAddressID   uuid.UUID `json:"billing_address_id" validate:"required"`
	PaymentMethod      string    `json:"payment_method" validate:"required"`
	Notes              string    `json:"notes,omitempty"`
}

func (c CreateOrderFromCartCommand) GetName() string {
	return "CreateOrderFromCart"
}
