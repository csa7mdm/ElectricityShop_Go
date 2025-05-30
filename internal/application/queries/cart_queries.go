package queries

import (
	"github.com/google/uuid"
)

// GetCartByUserIDQuery represents a query to get user's cart
type GetCartByUserIDQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (q GetCartByUserIDQuery) GetName() string {
	return "GetCartByUserID"
}

// GetCartByIDQuery represents a query to get cart by ID
type GetCartByIDQuery struct {
	CartID uuid.UUID `json:"cart_id" validate:"required"`
}

func (q GetCartByIDQuery) GetName() string {
	return "GetCartByID"
}

// GetCartItemsQuery represents a query to get cart items
type GetCartItemsQuery struct {
	CartID uuid.UUID `json:"cart_id" validate:"required"`
}

func (q GetCartItemsQuery) GetName() string {
	return "GetCartItems"
}

// GetCartSummaryQuery represents a query to get cart summary (totals, item count, etc.)
type GetCartSummaryQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (q GetCartSummaryQuery) GetName() string {
	return "GetCartSummary"
}
