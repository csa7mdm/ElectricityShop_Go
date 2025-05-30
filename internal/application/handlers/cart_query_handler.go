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

// CartQueryHandler handles cart-related queries
type CartQueryHandler struct {
	cartRepo interfaces.CartRepository
	logger   logger.Logger
}

// NewCartQueryHandler creates a new CartQueryHandler
func NewCartQueryHandler(
	cartRepo interfaces.CartRepository,
	logger logger.Logger,
) *CartQueryHandler {
	return &CartQueryHandler{
		cartRepo: cartRepo,
		logger:   logger,
	}
}

// Handle handles queries
func (h *CartQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	switch q := query.(type) {
	case *queries.GetCartByUserIDQuery:
		return h.handleGetCartByUserID(ctx, q)
	case *queries.GetCartByIDQuery:
		return h.handleGetCartByID(ctx, q)
	case *queries.GetCartItemsQuery:
		return h.handleGetCartItems(ctx, q)
	case *queries.GetCartSummaryQuery:
		return h.handleGetCartSummary(ctx, q)
	default:
		return nil, errors.New("UNSUPPORTED_QUERY", "Unsupported query type", 400)
	}
}

// handleGetCartByUserID handles getting a cart by user ID
func (h *CartQueryHandler) handleGetCartByUserID(ctx context.Context, query *queries.GetCartByUserIDQuery) (*entities.Cart, error) {
	h.logger.WithContext(ctx).Debugf("Getting cart for user: %s", query.UserID)
	
	cart, err := h.cartRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved cart: %s", cart.ID)
	return cart, nil
}

// handleGetCartByID handles getting a cart by ID
func (h *CartQueryHandler) handleGetCartByID(ctx context.Context, query *queries.GetCartByIDQuery) (*entities.Cart, error) {
	h.logger.WithContext(ctx).Debugf("Getting cart by ID: %s", query.CartID)
	
	cart, err := h.cartRepo.GetByID(ctx, query.CartID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved cart: %s", cart.ID)
	return cart, nil
}

// handleGetCartItems handles getting cart items
func (h *CartQueryHandler) handleGetCartItems(ctx context.Context, query *queries.GetCartItemsQuery) ([]*entities.CartItem, error) {
	h.logger.WithContext(ctx).Debugf("Getting items for cart: %s", query.CartID)
	
	items, err := h.cartRepo.GetItems(ctx, query.CartID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d cart items", len(items))
	return items, nil
}

// CartSummary represents a cart summary with totals and counts
type CartSummary struct {
	CartID       string          `json:"cart_id"`
	UserID       string          `json:"user_id"`
	ItemCount    int             `json:"item_count"`
	TotalItems   int             `json:"total_items"` // sum of all quantities
	Subtotal     decimal.Decimal `json:"subtotal"`
	TaxAmount    decimal.Decimal `json:"tax_amount"`
	Total        decimal.Decimal `json:"total"`
	IsEmpty      bool            `json:"is_empty"`
	LastUpdated  string          `json:"last_updated"`
}

// handleGetCartSummary handles getting cart summary
func (h *CartQueryHandler) handleGetCartSummary(ctx context.Context, query *queries.GetCartSummaryQuery) (*CartSummary, error) {
	h.logger.WithContext(ctx).Debugf("Getting cart summary for user: %s", query.UserID)
	
	// Get user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	// Get cart items
	items, err := h.cartRepo.GetItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	
	// Calculate summary
	summary := &CartSummary{
		CartID:      cart.ID.String(),
		UserID:      cart.UserID.String(),
		ItemCount:   len(items),
		TotalItems:  0,
		Subtotal:    decimal.Zero,
		TaxAmount:   decimal.Zero, // TODO: Implement tax calculation
		LastUpdated: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	
	for _, item := range items {
		summary.TotalItems += item.Quantity
		summary.Subtotal = summary.Subtotal.Add(item.Total)
	}
	
	// TODO: Calculate tax based on business rules
	// For now, assume 8% tax rate
	taxRate := decimal.NewFromFloat(0.08)
	summary.TaxAmount = summary.Subtotal.Mul(taxRate)
	summary.Total = summary.Subtotal.Add(summary.TaxAmount)
	summary.IsEmpty = len(items) == 0
	
	h.logger.WithContext(ctx).Debugf("Successfully calculated cart summary for user: %s", query.UserID)
	return summary, nil
}
