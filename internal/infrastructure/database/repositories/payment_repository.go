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

// PaymentRepository implements the PaymentRepository interface
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(db *gorm.DB) interfaces.PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create payment", 500)
	}
	return nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.User").
		First(&payment, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPaymentNotFound.WithDetails(fmt.Sprintf("Payment with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve payment", 500)
	}
	
	return &payment, nil
}

// GetByOrderID retrieves all payments for an order
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve order payments", 500)
	}
	
	return payments, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error) {
	var payment entities.Payment
	
	err := r.db.WithContext(ctx).
		Preload("Order").
		First(&payment, "transaction_id = ?", transactionID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPaymentNotFound.WithDetails(fmt.Sprintf("Payment with transaction ID %s not found", transactionID))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve payment", 500)
	}
	
	return &payment, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	if err := r.db.WithContext(ctx).Save(payment).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update payment", 500)
	}
	return nil
}

// UpdateStatus updates payment status
func (r *PaymentRepository) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entities.PaymentStatus) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Payment{}).
		Where("id = ?", paymentID).
		Update("status", status)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to update payment status", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrPaymentNotFound.WithDetails(fmt.Sprintf("Payment with ID %s not found", paymentID))
	}
	
	return nil
}

// List retrieves payments with filtering
func (r *PaymentRepository) List(ctx context.Context, filter interfaces.PaymentFilter) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	
	query := r.db.WithContext(ctx).Model(&entities.Payment{})
	
	// Apply filters
	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}
	
	if filter.UserID != nil {
		query = query.Joins("JOIN orders ON payments.order_id = orders.id").
			Where("orders.user_id = ?", *filter.UserID)
	}
	
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	
	// Apply date filters
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	
	// Apply sorting
	if filter.SortBy != "" {
		orderClause := filter.SortBy
		if filter.SortDesc {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("created_at DESC")
	}
	
	// Apply pagination
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	
	if err := query.
		Preload("Order").
		Preload("Order.User").
		Find(&payments).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to list payments", 500)
	}
	
	return payments, nil
}
