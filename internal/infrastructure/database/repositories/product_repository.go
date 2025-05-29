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

// ProductRepository implements the ProductRepository interface
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrProductAlreadyExists.WithDetails(fmt.Sprintf("Product with SKU %s already exists", product.SKU))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to create product", 500)
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Attributes", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_approved = ?", true).Order("created_at DESC")
		}).
		Preload("Reviews.User").
		First(&product, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrProductNotFound.WithDetails(fmt.Sprintf("Product with ID %s not found", id))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve product", 500)
	}
	
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Attributes", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&product, "sku = ?", sku).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrProductNotFound.WithDetails(fmt.Sprintf("Product with SKU %s not found", sku))
		}
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve product", 500)
	}
	
	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.ErrProductAlreadyExists.WithDetails(fmt.Sprintf("Product with SKU %s already exists", product.SKU))
		}
		return errors.Wrap(err, "DATABASE_ERROR", "Failed to update product", 500)
	}
	return nil
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Product{}, "id = ?", id)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to delete product", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrProductNotFound.WithDetails(fmt.Sprintf("Product with ID %s not found", id))
	}
	
	return nil
}

// List retrieves products with filtering
func (r *ProductRepository) List(ctx context.Context, filter interfaces.ProductFilter) ([]*entities.Product, error) {
	var products []*entities.Product
	
	query := r.db.WithContext(ctx).Model(&entities.Product{})
	
	// Apply filters
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	
	if filter.InStock != nil {
		if *filter.InStock {
			query = query.Where("stock > 0")
		} else {
			query = query.Where("stock = 0")
		}
	}
	
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}
	
	if filter.Brand != "" {
		query = query.Where("brand ILIKE ?", "%"+filter.Brand+"%")
	}
	
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ? OR brand ILIKE ?", 
			searchTerm, searchTerm, searchTerm, searchTerm)
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
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to list products", 500)
	}
	
	return products, nil
}

// Search searches products by query with filters
func (r *ProductRepository) Search(ctx context.Context, query string, filter interfaces.ProductFilter) ([]*entities.Product, error) {
	// Set search term in filter and use List method
	filter.Search = query
	return r.List(ctx, filter)
}

// GetByCategory retrieves products by category
func (r *ProductRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, filter interfaces.ProductFilter) ([]*entities.Product, error) {
	// Set category ID in filter and use List method
	filter.CategoryID = &categoryID
	return r.List(ctx, filter)
}

// UpdateStock updates product stock
func (r *ProductRepository) UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Update("stock", quantity)
	
	if result.Error != nil {
		return errors.Wrap(result.Error, "DATABASE_ERROR", "Failed to update product stock", 500)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrProductNotFound.WithDetails(fmt.Sprintf("Product with ID %s not found", productID))
	}
	
	return nil
}

// GetLowStockProducts retrieves products with stock below threshold
func (r *ProductRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]*entities.Product, error) {
	var products []*entities.Product
	
	if err := r.db.WithContext(ctx).
		Where("stock <= ? AND stock > 0 AND is_active = ?", threshold, true).
		Preload("Category").
		Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Failed to retrieve low stock products", 500)
	}
	
	return products, nil
}

// ExistsBySKU checks if a product exists by SKU
func (r *ProductRepository) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var count int64
	
	if err := r.db.WithContext(ctx).Model(&entities.Product{}).Where("sku = ?", sku).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "DATABASE_ERROR", "Failed to check product existence", 500)
	}
	
	return count > 0, nil
}
