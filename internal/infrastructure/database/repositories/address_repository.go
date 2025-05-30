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

// AddressRepository implements the AddressRepository interface
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new AddressRepository
func NewAddressRepository(db *gorm.DB) interfaces.AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, address *entities.Address) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create address", 500)
	}
	return nil
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Address, error) {
	var address entities.Address
	
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&address, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrAddressNotFound.WithDetails(fmt.Sprintf("Address with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve address", 500)
	}
	
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error) {
	var addresses []*entities.Address
	
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve user addresses", 500)
	}
	
	return addresses, nil
}

// Update updates an address
func (r *AddressRepository) Update(ctx context.Context, address *entities.Address) error {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update address", 500)
	}
	return nil
}

// Delete deletes an address
func (r *AddressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Address{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete address", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrAddressNotFound.WithDetails(fmt.Sprintf("Address with ID %s not found", id))
	}
	
	return nil
}

// SetAsDefault sets an address as default for a specific type
func (r *AddressRepository) SetAsDefault(ctx context.Context, addressID uuid.UUID, addressType entities.AddressType) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "Failed to start transaction", 500)
	}
	
	// First, get the address to find the user ID
	var address entities.Address
	if err := tx.First(&address, "id = ?", addressID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return errors.ErrAddressNotFound.WithDetails(fmt.Sprintf("Address with ID %s not found", addressID))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve address", 500)
	}
	
	// Unset all other addresses of the same type for this user as default
	if err := tx.Model(&entities.Address{}).
		Where("user_id = ? AND type = ? AND id != ?", address.UserID, addressType, addressID).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to unset other default addresses", 500)
	}
	
	// Set the specified address as default
	if err := tx.Model(&entities.Address{}).
		Where("id = ?", addressID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to set address as default", 500)
	}
	
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to commit transaction", 500)
	}
	
	return nil
}

// UnsetDefaultForUser unsets all default addresses for a user
func (r *AddressRepository) UnsetDefaultForUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&entities.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error; err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to unset default addresses", 500)
	}
	return nil
}
