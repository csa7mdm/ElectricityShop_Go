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

// CategoryRepository implements the CategoryRepository interface
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrCategoryAlreadyExists.WithDetails(fmt.Sprintf("Category with slug %s already exists", category.Slug))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create category", 500)
	}
	return nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var category entities.Category
	
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		First(&category, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound.WithDetails(fmt.Sprintf("Category with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve category", 500)
	}
	
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	var category entities.Category
	
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		First(&category, "slug = ?", slug).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound.WithDetails(fmt.Sprintf("Category with slug %s not found", slug))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve category", 500)
	}
	
	return &category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrCategoryAlreadyExists.WithDetails(fmt.Sprintf("Category with slug %s already exists", category.Slug))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update category", 500)
	}
	return nil
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Category{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete category", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrCategoryNotFound.WithDetails(fmt.Sprintf("Category with ID %s not found", id))
	}
	
	return nil
}

// List retrieves categories with filtering
func (r *CategoryRepository) List(ctx context.Context, filter interfaces.CategoryFilter) ([]*entities.Category, error) {
	var categories []*entities.Category
	
	query := r.db.WithContext(ctx).Model(&entities.Category{})
	
	// Apply filters
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	// Apply sorting
	if filter.SortBy != "" {
		orderClause := filter.SortBy
		if filter.SortDesc {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("sort_order ASC, name ASC")
	}
	
	// Apply pagination
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	
	if err := query.
		Preload("Parent").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Find(&categories).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to list categories", 500)
	}
	
	return categories, nil
}

// GetChildren retrieves child categories of a parent
func (r *CategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category
	
	if err := r.db.WithContext(ctx).
		Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve child categories", 500)
	}
	
	return categories, nil
}

// GetRootCategories retrieves root categories (categories without parent)
func (r *CategoryRepository) GetRootCategories(ctx context.Context) ([]*entities.Category, error) {
	var categories []*entities.Category
	
	if err := r.db.WithContext(ctx).
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Find(&categories).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve root categories", 500)
	}
	
	return categories, nil
}

// ExistsBySlug checks if a category exists by slug
func (r *CategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	
	if err := r.db.WithContext(ctx).Model(&entities.Category{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "DATABASE_ERROR", "Failed to check category existence", 500)
	}
	
	return count > 0, nil
}
