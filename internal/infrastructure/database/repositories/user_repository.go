package repositories

import (
	"context"
	"errors" // For gorm.ErrRecordNotFound
	"fmt"    // For not implemented messages

	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	domainInterfaces "github.com/yourusername/electricity-shop-go/internal/domain/interfaces" // Alias for domain interfaces
	// "github.com/yourusername/electricity-shop-go/internal/infrastructure/database"           // For database.DB or GetDB() - Not directly needed if DB is passed in constructor
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

// NewGORMUserRepository creates a new instance of UserRepository.
// It expects a GORM DB connection to be passed.
func NewGORMUserRepository(db *gorm.DB) domainInterfaces.UserRepository {
	if db == nil {
		// This is a critical error, the repository cannot function without a DB.
		panic("gormUserRepository received a nil DB instance")
	}
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Using WithContext for good practice, though GORM might handle some context propagation.
	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Standard practice: return nil, nil if not found
		}
		return nil, err // Other DB errors
	}
	return &user, nil
}

func (r *gormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Standard practice: return nil, nil if not found
		}
		return nil, err // Other DB errors
	}
	return &user, nil
}

// --- Methods to be implemented later ---

func (r *gormUserRepository) Update(ctx context.Context, user *entities.User) error {
	return fmt.Errorf("Update not implemented")
}

func (r *gormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("Delete not implemented")
}

func (r *gormUserRepository) List(ctx context.Context, filter domainInterfaces.UserFilter) ([]*entities.User, error) {
	return nil, fmt.Errorf("List not implemented")
}

func (r *gormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, fmt.Errorf("ExistsByEmail not implemented")
}
