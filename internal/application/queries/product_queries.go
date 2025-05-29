package queries

import (
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
)

// GetProductByIDQuery represents a query to get a product by ID
type GetProductByIDQuery struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (q GetProductByIDQuery) GetName() string {
	return "GetProductByID"
}

// GetProductBySKUQuery represents a query to get a product by SKU
type GetProductBySKUQuery struct {
	SKU string `json:"sku" validate:"required"`
}

func (q GetProductBySKUQuery) GetName() string {
	return "GetProductBySKU"
}

// ListProductsQuery represents a query to list products with filtering
type ListProductsQuery struct {
	Filter interfaces.ProductFilter `json:"filter"`
}

func (q ListProductsQuery) GetName() string {
	return "ListProducts"
}

// SearchProductsQuery represents a query to search products
type SearchProductsQuery struct {
	Query  string                     `json:"query" validate:"required"`
	Filter interfaces.ProductFilter `json:"filter"`
}

func (q SearchProductsQuery) GetName() string {
	return "SearchProducts"
}

// GetProductsByCategoryQuery represents a query to get products by category
type GetProductsByCategoryQuery struct {
	CategoryID uuid.UUID               `json:"category_id" validate:"required"`
	Filter     interfaces.ProductFilter `json:"filter"`
}

func (q GetProductsByCategoryQuery) GetName() string {
	return "GetProductsByCategory"
}

// GetLowStockProductsQuery represents a query to get low stock products
type GetLowStockProductsQuery struct {
	Threshold int `json:"threshold" validate:"min=0"`
}

func (q GetLowStockProductsQuery) GetName() string {
	return "GetLowStockProducts"
}

// GetCategoryByIDQuery represents a query to get a category by ID
type GetCategoryByIDQuery struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

func (q GetCategoryByIDQuery) GetName() string {
	return "GetCategoryByID"
}

// GetCategoryBySlugQuery represents a query to get a category by slug
type GetCategoryBySlugQuery struct {
	Slug string `json:"slug" validate:"required"`
}

func (q GetCategoryBySlugQuery) GetName() string {
	return "GetCategoryBySlug"
}

// ListCategoriesQuery represents a query to list categories with filtering
type ListCategoriesQuery struct {
	Filter interfaces.CategoryFilter `json:"filter"`
}

func (q ListCategoriesQuery) GetName() string {
	return "ListCategories"
}

// GetCategoryChildrenQuery represents a query to get category children
type GetCategoryChildrenQuery struct {
	ParentID uuid.UUID `json:"parent_id" validate:"required"`
}

func (q GetCategoryChildrenQuery) GetName() string {
	return "GetCategoryChildren"
}

// GetRootCategoriesQuery represents a query to get root categories
type GetRootCategoriesQuery struct{}

func (q GetRootCategoriesQuery) GetName() string {
	return "GetRootCategories"
}
