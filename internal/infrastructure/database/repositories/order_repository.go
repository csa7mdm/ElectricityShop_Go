package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
)

// OrderRepository implements the OrderRepository interface
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrDuplicateOrderNumber.WithDetails(fmt.Sprintf("Order with number %s already exists", order.OrderNumber))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create order", 500)
	}
	return nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Preload("Payments").
		Preload("Shipments").
		First(&order, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrOrderNotFound.WithDetails(fmt.Sprintf("Order with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve order", 500)
	}
	
	return &order, nil
}

// GetByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error) {
	var order entities.Order
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payments").
		Preload("Shipments").
		First(&order, "order_number = ?", orderNumber).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrOrderNotFound.WithDetails(fmt.Sprintf("Order with number %s not found", orderNumber))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve order", 500)
	}
	
	return &order, nil
}

// Update updates an order
func (r *OrderRepository) Update(ctx context.Context, order *entities.Order) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update order", 500)
	}
	return nil
}

// Delete soft deletes an order
func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Order{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete order", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrOrderNotFound.WithDetails(fmt.Sprintf("Order with ID %s not found", id))
	}
	
	return nil
}

// GetByUserID retrieves orders for a specific user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filter interfaces.OrderFilter) ([]*entities.Order, error) {
	var orders []*entities.Order
	
	query := r.db.WithContext(ctx).Model(&entities.Order{}).Where("user_id = ?", userID)
	
	// Apply filters
	query = r.applyOrderFilters(query, filter)
	
	if err := query.
		Preload("Items").
		Preload("Items.Product").
		Preload("Payments").
		Preload("Shipments").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve user orders", 500)
	}
	
	return orders, nil
}

// List retrieves orders with filtering
func (r *OrderRepository) List(ctx context.Context, filter interfaces.OrderFilter) ([]*entities.Order, error) {
	var orders []*entities.Order
	
	query := r.db.WithContext(ctx).Model(&entities.Order{})
	
	// Apply filters
	query = r.applyOrderFilters(query, filter)
	
	if err := query.
		Preload("User").
		Preload("Items").
		Preload("Payments").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to list orders", 500)
	}
	
	return orders, nil
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("id = ?", orderID).
		Update("status", status)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to update order status", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrOrderNotFound.WithDetails(fmt.Sprintf("Order with ID %s not found", orderID))
	}
	
	return nil
}

// GetOrdersToProcess retrieves orders that need processing
func (r *OrderRepository) GetOrdersToProcess(ctx context.Context) ([]*entities.Order, error) {
	var orders []*entities.Order
	
	if err := r.db.WithContext(ctx).
		Where("status IN ? AND payment_status = ?", 
			[]entities.OrderStatus{entities.OrderStatusConfirmed, entities.OrderStatusProcessing},
			entities.PaymentStatusCompleted).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve orders to process", 500)
	}
	
	return orders, nil
}

// GetOrdersByDateRange retrieves orders within a date range
func (r *OrderRepository) GetOrdersByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Order, error) {
	var orders []*entities.Order
	
	query := r.db.WithContext(ctx).Model(&entities.Order{})
	
	if startDate != "" {
		query = query.Where("ordered_at >= ?", startDate)
	}
	
	if endDate != "" {
		query = query.Where("ordered_at <= ?", endDate)
	}
	
	if err := query.
		Preload("User").
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve orders by date range", 500)
	}
	
	return orders, nil
}

// applyOrderFilters applies filtering to order queries
func (r *OrderRepository) applyOrderFilters(query *gorm.DB, filter interfaces.OrderFilter) *gorm.DB {
	// Apply user filter
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	
	// Apply status filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	
	// Apply date filters
	if filter.StartDate != nil {
		query = query.Where("ordered_at >= ?", *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("ordered_at <= ?", *filter.EndDate)
	}
	
	// Apply amount filters
	if filter.MinTotal != nil {
		query = query.Where("total >= ?", *filter.MinTotal)
	}
	
	if filter.MaxTotal != nil {
		query = query.Where("total <= ?", *filter.MaxTotal)
	}
	
	// Apply sorting
	if filter.SortBy != "" {
		orderClause := filter.SortBy
		if filter.SortDesc {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("ordered_at DESC")
	}
	
	// Apply pagination
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	
	return query
}
