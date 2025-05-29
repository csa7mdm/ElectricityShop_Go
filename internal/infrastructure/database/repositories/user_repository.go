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

// UserRepository implements the UserRepository interface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrUserAlreadyExists.WithDetails(fmt.Sprintf("User with email %s already exists", user.Email))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create user", 500)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		First(&user, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound.WithDetails(fmt.Sprintf("User with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve user", 500)
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		First(&user, "email = ?", email).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound.WithDetails(fmt.Sprintf("User with email %s not found", email))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve user", 500)
	}
	
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrUserAlreadyExists.WithDetails(fmt.Sprintf("User with email %s already exists", user.Email))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update user", 500)
	}
	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete user", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound.WithDetails(fmt.Sprintf("User with ID %s not found", id))
	}
	
	return nil
}

// List retrieves users with filtering
func (r *UserRepository) List(ctx context.Context, filter interfaces.UserFilter) ([]*entities.User, error) {
	var users []*entities.User
	
	query := r.db.WithContext(ctx).Model(&entities.User{})
	
	// Apply filters
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", 
			searchTerm, searchTerm, searchTerm)
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
	
	if err := query.Preload("Addresses").Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to list users", 500)
	}
	
	return users, nil
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	
	if err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "DATABASE_ERROR", "Failed to check user existence", 500)
	}
	
	return count > 0, nil
}

// Helper function to check if error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	// PostgreSQL unique constraint error typically contains "duplicate key value violates unique constraint"
	errStr := err.Error()
	return containsAny(errStr, []string{
		"duplicate key value violates unique constraint",
		"UNIQUE constraint failed",
		"violates unique constraint",
	})
}

// Helper function to check if string contains any of the substrings
func containsAny(str string, substrings []string) bool {
	for _, substring := range substrings {
		if len(str) >= len(substring) {
			for i := 0; i <= len(str)-len(substring); i++ {
				if str[i:i+len(substring)] == substring {
					return true
				}
			}
		}
	}
	return false
}
