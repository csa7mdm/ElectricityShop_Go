package handlers

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// CartCommandHandler handles cart-related commands
type CartCommandHandler struct {
	cartRepo       interfaces.CartRepository
	productRepo    interfaces.ProductRepository
	userRepo       interfaces.UserRepository
	eventPublisher interfaces.EventPublisher
	logger         logger.Logger
}

// NewCartCommandHandler creates a new CartCommandHandler
func NewCartCommandHandler(
	cartRepo interfaces.CartRepository,
	productRepo interfaces.ProductRepository,
	userRepo interfaces.UserRepository,
	eventPublisher interfaces.EventPublisher,
	logger logger.Logger,
) *CartCommandHandler {
	return &CartCommandHandler{
		cartRepo:       cartRepo,
		productRepo:    productRepo,
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// Handle handles commands
func (h *CartCommandHandler) Handle(ctx context.Context, command mediator.Command) error {
	switch cmd := command.(type) {
	case *commands.AddToCartCommand:
		return h.handleAddToCart(ctx, cmd)
	case *commands.UpdateCartItemCommand:
		return h.handleUpdateCartItem(ctx, cmd)
	case *commands.RemoveFromCartCommand:
		return h.handleRemoveFromCart(ctx, cmd)
	case *commands.ClearCartCommand:
		return h.handleClearCart(ctx, cmd)
	default:
		return errors.New("UNSUPPORTED_COMMAND", "Unsupported command type", 400)
	}
}

// handleAddToCart handles adding an item to the cart
func (h *CartCommandHandler) handleAddToCart(ctx context.Context, cmd *commands.AddToCartCommand) error {
	h.logger.WithContext(ctx).Infof("Adding item to cart for user: %s", cmd.UserID)
	
	// Verify user exists
	_, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Verify product exists and is available
	product, err := h.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	
	if !product.IsActive {
		return errors.New("PRODUCT_UNAVAILABLE", "Product is not available", 400)
	}
	
	if !product.CanOrder(cmd.Quantity) {
		return errors.ErrInsufficientStock.WithDetails(fmt.Sprintf("Only %d items available", product.Stock))
	}
	
	// Get or create user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Create cart item
	cartItem := &entities.CartItem{
		CartID:    cart.ID,
		ProductID: cmd.ProductID,
		Quantity:  cmd.Quantity,
		UnitPrice: product.Price,
		Total:     product.Price.Mul(decimal.NewFromInt(int64(cmd.Quantity))),
	}
	
	// Add item to cart
	if err := h.cartRepo.AddItem(ctx, cartItem); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewCartItemAddedEvent(
		cart.ID,
		cmd.UserID,
		cmd.ProductID,
		cmd.Quantity,
		product.Price,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish CartItemAddedEvent: %v", err)
		// Don't fail the command for event publishing errors
	}
	
	h.logger.WithContext(ctx).Infof("Successfully added item to cart for user: %s", cmd.UserID)
	return nil
}

// handleUpdateCartItem handles updating cart item quantity
func (h *CartCommandHandler) handleUpdateCartItem(ctx context.Context, cmd *commands.UpdateCartItemCommand) error {
	h.logger.WithContext(ctx).Infof("Updating cart item for user: %s", cmd.UserID)
	
	// Get user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Get cart item
	cartItem, err := h.cartRepo.GetItemByProductID(ctx, cart.ID, cmd.ProductID)
	if err != nil {
		return err
	}
	
	// Verify product availability
	product, err := h.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	
	if !product.CanOrder(cmd.Quantity) {
		return errors.ErrInsufficientStock.WithDetails(fmt.Sprintf("Only %d items available", product.Stock))
	}
	
	// Update cart item
	cartItem.Quantity = cmd.Quantity
	cartItem.Total = cartItem.UnitPrice.Mul(decimal.NewFromInt(int64(cmd.Quantity)))
	
	if err := h.cartRepo.UpdateItem(ctx, cartItem); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated cart item for user: %s", cmd.UserID)
	return nil
}

// handleRemoveFromCart handles removing an item from the cart
func (h *CartCommandHandler) handleRemoveFromCart(ctx context.Context, cmd *commands.RemoveFromCartCommand) error {
	h.logger.WithContext(ctx).Infof("Removing item from cart for user: %s", cmd.UserID)
	
	// Get user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Remove item from cart
	if err := h.cartRepo.RemoveItem(ctx, cart.ID, cmd.ProductID); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully removed item from cart for user: %s", cmd.UserID)
	return nil
}

// handleClearCart handles clearing all items from the cart
func (h *CartCommandHandler) handleClearCart(ctx context.Context, cmd *commands.ClearCartCommand) error {
	h.logger.WithContext(ctx).Infof("Clearing cart for user: %s", cmd.UserID)
	
	// Get user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Clear all items
	if err := h.cartRepo.ClearItems(ctx, cart.ID); err != nil {
		return err
	}
	
	// Publish domain event
	reason := cmd.Reason
	if reason == "" {
		reason = "User action"
	}
	
	event := events.NewCartClearedEvent(cart.ID, cmd.UserID, reason)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish CartClearedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully cleared cart for user: %s", cmd.UserID)
	return nil
}
