package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
)

// CartRepository implements the CartRepository interface
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(db *gorm.DB) interfaces.CartRepository {
	return &CartRepository{db: db}
}

// Create creates a new cart
func (r *CartRepository) Create(ctx context.Context, cart *entities.Cart) error {
	if err := r.db.WithContext(ctx).Create(cart).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create cart", 500)
	}
	return nil
}

// GetByID retrieves a cart by ID
func (r *CartRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Cart, error) {
	var cart entities.Cart
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		First(&cart, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCartNotFound.WithDetails(fmt.Sprintf("Cart with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve cart", 500)
	}
	
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Cart, error) {
	var cart entities.Cart
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		First(&cart, "user_id = ?", userID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new cart if none exists
			newCart := &entities.Cart{
				UserID: userID,
				Items:  []entities.CartItem{},
			}
			if createErr := r.Create(ctx, newCart); createErr != nil {
				return nil, createErr
			}
			return newCart, nil
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve cart", 500)
	}
	
	return &cart, nil
}

// GetBySessionID retrieves a cart by session ID (for guest users)
func (r *CartRepository) GetBySessionID(ctx context.Context, sessionID string) (*entities.Cart, error) {
	var cart entities.Cart
	
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		First(&cart, "session_id = ?", sessionID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCartNotFound.WithDetails(fmt.Sprintf("Cart with session ID %s not found", sessionID))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve cart", 500)
	}
	
	return &cart, nil
}

// Update updates a cart
func (r *CartRepository) Update(ctx context.Context, cart *entities.Cart) error {
	if err := r.db.WithContext(ctx).Save(cart).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update cart", 500)
	}
	return nil
}

// Delete deletes a cart
func (r *CartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Cart{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete cart", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrCartNotFound.WithDetails(fmt.Sprintf("Cart with ID %s not found", id))
	}
	
	return nil
}

// AddItem adds an item to the cart
func (r *CartRepository) AddItem(ctx context.Context, cartItem *entities.CartItem) error {
	// Check if item already exists in cart
	existingItem, err := r.GetItemByProductID(ctx, cartItem.CartID, cartItem.ProductID)
	if err != nil && !errors.IsAppError(err) {
		return err
	}
	
	if existingItem != nil {
		// Update quantity if item already exists
		existingItem.Quantity += cartItem.Quantity
		existingItem.Total = existingItem.UnitPrice.Mul(decimal.NewFromInt(int64(existingItem.Quantity)))
		return r.UpdateItem(ctx, existingItem)
	}
	
	// Add new item
	if err := r.db.WithContext(ctx).Create(cartItem).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to add item to cart", 500)
	}
	
	return nil
}

// UpdateItem updates a cart item
func (r *CartRepository) UpdateItem(ctx context.Context, cartItem *entities.CartItem) error {
	if err := r.db.WithContext(ctx).Save(cartItem).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update cart item", 500)
	}
	return nil
}

// RemoveItem removes an item from the cart
func (r *CartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&entities.CartItem{})
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to remove item from cart", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.New("CART_ITEM_NOT_FOUND", "Cart item not found", 404)
	}
	
	return nil
}

// ClearItems removes all items from the cart
func (r *CartRepository) ClearItems(ctx context.Context, cartID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&entities.CartItem{}).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to clear cart items", 500)
	}
	
	return nil
}

// GetItems retrieves all items in a cart
func (r *CartRepository) GetItems(ctx context.Context, cartID uuid.UUID) ([]*entities.CartItem, error) {
	var items []*entities.CartItem
	
	if err := r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Find(&items).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve cart items", 500)
	}
	
	return items, nil
}

// GetItemByProductID retrieves a specific cart item by product ID
func (r *CartRepository) GetItemByProductID(ctx context.Context, cartID, productID uuid.UUID) (*entities.CartItem, error) {
	var item entities.CartItem
	
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Preload("Product").
		Preload("Product.Category").
		First(&item).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("CART_ITEM_NOT_FOUND", "Cart item not found", 404)
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve cart item", 500)
	}
	
	return &item, nil
}
