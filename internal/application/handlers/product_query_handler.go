package handlers

import (
	"context"

	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// ProductQueryHandler handles product-related queries
type ProductQueryHandler struct {
	productRepo  interfaces.ProductRepository
	categoryRepo interfaces.CategoryRepository
	logger       logger.Logger
}

// NewProductQueryHandler creates a new ProductQueryHandler
func NewProductQueryHandler(
	productRepo interfaces.ProductRepository,
	categoryRepo interfaces.CategoryRepository,
	logger logger.Logger,
) *ProductQueryHandler {
	return &ProductQueryHandler{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// Handle handles queries
func (h *ProductQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	switch q := query.(type) {
	case *queries.GetProductByIDQuery:
		return h.handleGetProductByID(ctx, q)
	case *queries.GetProductBySKUQuery:
		return h.handleGetProductBySKU(ctx, q)
	case *queries.ListProductsQuery:
		return h.handleListProducts(ctx, q)
	case *queries.SearchProductsQuery:
		return h.handleSearchProducts(ctx, q)
	case *queries.GetProductsByCategoryQuery:
		return h.handleGetProductsByCategory(ctx, q)
	case *queries.GetLowStockProductsQuery:
		return h.handleGetLowStockProducts(ctx, q)
	case *queries.GetCategoryByIDQuery:
		return h.handleGetCategoryByID(ctx, q)
	case *queries.GetCategoryBySlugQuery:
		return h.handleGetCategoryBySlug(ctx, q)
	case *queries.ListCategoriesQuery:
		return h.handleListCategories(ctx, q)
	case *queries.GetCategoryChildrenQuery:
		return h.handleGetCategoryChildren(ctx, q)
	case *queries.GetRootCategoriesQuery:
		return h.handleGetRootCategories(ctx, q)
	default:
		return nil, errors.New("UNSUPPORTED_QUERY", "Unsupported query type", 400)
	}
}

// handleGetProductByID handles getting a product by ID
func (h *ProductQueryHandler) handleGetProductByID(ctx context.Context, query *queries.GetProductByIDQuery) (*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Getting product by ID: %s", query.ProductID)
	
	product, err := h.productRepo.GetByID(ctx, query.ProductID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved product: %s", product.ID)
	return product, nil
}

// handleGetProductBySKU handles getting a product by SKU
func (h *ProductQueryHandler) handleGetProductBySKU(ctx context.Context, query *queries.GetProductBySKUQuery) (*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Getting product by SKU: %s", query.SKU)
	
	product, err := h.productRepo.GetBySKU(ctx, query.SKU)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved product: %s", product.ID)
	return product, nil
}

// handleListProducts handles listing products with filtering
func (h *ProductQueryHandler) handleListProducts(ctx context.Context, query *queries.ListProductsQuery) ([]*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Listing products with filter")
	
	products, err := h.productRepo.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d products", len(products))
	return products, nil
}

// handleSearchProducts handles searching products
func (h *ProductQueryHandler) handleSearchProducts(ctx context.Context, query *queries.SearchProductsQuery) ([]*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Searching products with query: %s", query.Query)
	
	products, err := h.productRepo.Search(ctx, query.Query, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully found %d products", len(products))
	return products, nil
}

// handleGetProductsByCategory handles getting products by category
func (h *ProductQueryHandler) handleGetProductsByCategory(ctx context.Context, query *queries.GetProductsByCategoryQuery) ([]*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Getting products by category: %s", query.CategoryID)
	
	// Verify category exists
	_, err := h.categoryRepo.GetByID(ctx, query.CategoryID)
	if err != nil {
		return nil, err
	}
	
	products, err := h.productRepo.GetByCategory(ctx, query.CategoryID, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d products for category: %s", len(products), query.CategoryID)
	return products, nil
}

// handleGetLowStockProducts handles getting low stock products
func (h *ProductQueryHandler) handleGetLowStockProducts(ctx context.Context, query *queries.GetLowStockProductsQuery) ([]*entities.Product, error) {
	h.logger.WithContext(ctx).Debugf("Getting low stock products with threshold: %d", query.Threshold)
	
	products, err := h.productRepo.GetLowStockProducts(ctx, query.Threshold)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d low stock products", len(products))
	return products, nil
}

// handleGetCategoryByID handles getting a category by ID
func (h *ProductQueryHandler) handleGetCategoryByID(ctx context.Context, query *queries.GetCategoryByIDQuery) (*entities.Category, error) {
	h.logger.WithContext(ctx).Debugf("Getting category by ID: %s", query.CategoryID)
	
	category, err := h.categoryRepo.GetByID(ctx, query.CategoryID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved category: %s", category.ID)
	return category, nil
}

// handleGetCategoryBySlug handles getting a category by slug
func (h *ProductQueryHandler) handleGetCategoryBySlug(ctx context.Context, query *queries.GetCategoryBySlugQuery) (*entities.Category, error) {
	h.logger.WithContext(ctx).Debugf("Getting category by slug: %s", query.Slug)
	
	category, err := h.categoryRepo.GetBySlug(ctx, query.Slug)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved category: %s", category.ID)
	return category, nil
}

// handleListCategories handles listing categories with filtering
func (h *ProductQueryHandler) handleListCategories(ctx context.Context, query *queries.ListCategoriesQuery) ([]*entities.Category, error) {
	h.logger.WithContext(ctx).Debugf("Listing categories with filter")
	
	categories, err := h.categoryRepo.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d categories", len(categories))
	return categories, nil
}

// handleGetCategoryChildren handles getting category children
func (h *ProductQueryHandler) handleGetCategoryChildren(ctx context.Context, query *queries.GetCategoryChildrenQuery) ([]*entities.Category, error) {
	h.logger.WithContext(ctx).Debugf("Getting children for category: %s", query.ParentID)
	
	// Verify parent category exists
	_, err := h.categoryRepo.GetByID(ctx, query.ParentID)
	if err != nil {
		return nil, err
	}
	
	children, err := h.categoryRepo.GetChildren(ctx, query.ParentID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d children for category: %s", len(children), query.ParentID)
	return children, nil
}

// handleGetRootCategories handles getting root categories
func (h *ProductQueryHandler) handleGetRootCategories(ctx context.Context, query *queries.GetRootCategoriesQuery) ([]*entities.Category, error) {
	h.logger.WithContext(ctx).Debugf("Getting root categories")
	
	categories, err := h.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d root categories", len(categories))
	return categories, nil
}
