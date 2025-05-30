package handlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// ProductCommandHandler handles product-related commands
type ProductCommandHandler struct {
	productRepo     interfaces.ProductRepository
	categoryRepo    interfaces.CategoryRepository
	eventPublisher  interfaces.EventPublisher
	logger          logger.Logger
}

// NewProductCommandHandler creates a new ProductCommandHandler
func NewProductCommandHandler(
	productRepo interfaces.ProductRepository,
	categoryRepo interfaces.CategoryRepository,
	eventPublisher interfaces.EventPublisher,
	logger logger.Logger,
) *ProductCommandHandler {
	return &ProductCommandHandler{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// Handle handles commands
func (h *ProductCommandHandler) Handle(ctx context.Context, command mediator.Command) error {
	switch cmd := command.(type) {
	case *commands.CreateProductCommand:
		return h.handleCreateProduct(ctx, cmd)
	case *commands.UpdateProductCommand:
		return h.handleUpdateProduct(ctx, cmd)
	case *commands.UpdateProductStockCommand:
		return h.handleUpdateProductStock(ctx, cmd)
	case *commands.DeleteProductCommand:
		return h.handleDeleteProduct(ctx, cmd)
	case *commands.CreateCategoryCommand:
		return h.handleCreateCategory(ctx, cmd)
	case *commands.UpdateCategoryCommand:
		return h.handleUpdateCategory(ctx, cmd)
	case *commands.DeleteCategoryCommand:
		return h.handleDeleteCategory(ctx, cmd)
	default:
		return errors.New("UNSUPPORTED_COMMAND", "Unsupported command type", 400)
	}
}

// handleCreateProduct handles product creation
func (h *ProductCommandHandler) handleCreateProduct(ctx context.Context, cmd *commands.CreateProductCommand) error {
	h.logger.WithContext(ctx).Infof("Creating product with SKU: %s", cmd.SKU)
	
	// Check if product with SKU already exists
	exists, err := h.productRepo.ExistsBySKU(ctx, cmd.SKU)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrProductAlreadyExists.WithDetails("Product with this SKU already exists")
	}
	
	// Verify category exists
	_, err = h.categoryRepo.GetByID(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}
	
	// Create product entity
	product := &entities.Product{
		Name:        cmd.Name,
		Description: cmd.Description,
		SKU:         cmd.SKU,
		Price:       cmd.Price,
		CategoryID:  cmd.CategoryID,
		Brand:       cmd.Brand,
		Model:       cmd.Model,
		Weight:      cmd.Weight,
		Dimensions:  cmd.Dimensions,
		Color:       cmd.Color,
		Material:    cmd.Material,
		Warranty:    cmd.Warranty,
		Stock:       cmd.Stock,
		MinStock:    cmd.MinStock,
		MaxStock:    cmd.MaxStock,
		IsActive:    true,
		IsFeatured:  cmd.IsFeatured,
		MetaTitle:   cmd.MetaTitle,
		MetaDesc:    cmd.MetaDesc,
		Tags:        cmd.Tags,
	}
	
	// Save product
	if err := h.productRepo.Create(ctx, product); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewProductCreatedEvent(
		product.ID,
		product.Name,
		product.SKU,
		product.Price,
		product.CategoryID,
		product.Stock,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish ProductCreatedEvent: %v", err)
		// Don't fail the command for event publishing errors
	}
	
	h.logger.WithContext(ctx).Infof("Successfully created product: %s", product.ID)
	return nil
}

// handleUpdateProduct handles product updates
func (h *ProductCommandHandler) handleUpdateProduct(ctx context.Context, cmd *commands.UpdateProductCommand) error {
	h.logger.WithContext(ctx).Infof("Updating product: %s", cmd.ProductID)
	
	// Get existing product
	product, err := h.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	
	// Verify category exists
	_, err = h.categoryRepo.GetByID(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}
	
	// Update fields
	product.Name = cmd.Name
	product.Description = cmd.Description
	product.Price = cmd.Price
	product.CategoryID = cmd.CategoryID
	product.Brand = cmd.Brand
	product.Model = cmd.Model
	product.Weight = cmd.Weight
	product.Dimensions = cmd.Dimensions
	product.Color = cmd.Color
	product.Material = cmd.Material
	product.Warranty = cmd.Warranty
	product.MinStock = cmd.MinStock
	product.MaxStock = cmd.MaxStock
	product.IsFeatured = cmd.IsFeatured
	product.MetaTitle = cmd.MetaTitle
	product.MetaDesc = cmd.MetaDesc
	product.Tags = cmd.Tags
	
	// Save product
	if err := h.productRepo.Update(ctx, product); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated product: %s", product.ID)
	return nil
}

// handleUpdateProductStock handles product stock updates
func (h *ProductCommandHandler) handleUpdateProductStock(ctx context.Context, cmd *commands.UpdateProductStockCommand) error {
	h.logger.WithContext(ctx).Infof("Updating stock for product: %s", cmd.ProductID)
	
	// Get existing product to get old stock value
	product, err := h.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	
	oldStock := product.Stock
	
	// Update stock
	if err := h.productRepo.UpdateStock(ctx, cmd.ProductID, cmd.Quantity); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewProductStockUpdatedEvent(
		cmd.ProductID,
		oldStock,
		cmd.Quantity,
		cmd.Reason,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish ProductStockUpdatedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated stock for product: %s", cmd.ProductID)
	return nil
}

// handleDeleteProduct handles product deletion
func (h *ProductCommandHandler) handleDeleteProduct(ctx context.Context, cmd *commands.DeleteProductCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting product: %s", cmd.ProductID)
	
	// Check if product exists
	_, err := h.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	
	// Delete product (soft delete)
	if err := h.productRepo.Delete(ctx, cmd.ProductID); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully deleted product: %s", cmd.ProductID)
	return nil
}

// handleCreateCategory handles category creation
func (h *ProductCommandHandler) handleCreateCategory(ctx context.Context, cmd *commands.CreateCategoryCommand) error {
	h.logger.WithContext(ctx).Infof("Creating category with slug: %s", cmd.Slug)
	
	// Check if category with slug already exists
	exists, err := h.categoryRepo.ExistsBySlug(ctx, cmd.Slug)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrCategoryAlreadyExists.WithDetails("Category with this slug already exists")
	}
	
	// Verify parent category exists if provided
	if cmd.ParentID != nil {
		_, err = h.categoryRepo.GetByID(ctx, *cmd.ParentID)
		if err != nil {
			return err
		}
	}
	
	// Create category entity
	category := &entities.Category{
		Name:        cmd.Name,
		Slug:        cmd.Slug,
		Description: cmd.Description,
		ParentID:    cmd.ParentID,
		ImageURL:    cmd.ImageURL,
		SortOrder:   cmd.SortOrder,
		IsActive:    true,
		MetaTitle:   cmd.MetaTitle,
		MetaDesc:    cmd.MetaDesc,
	}
	
	// Save category
	if err := h.categoryRepo.Create(ctx, category); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully created category: %s", category.ID)
	return nil
}

// handleUpdateCategory handles category updates
func (h *ProductCommandHandler) handleUpdateCategory(ctx context.Context, cmd *commands.UpdateCategoryCommand) error {
	h.logger.WithContext(ctx).Infof("Updating category: %s", cmd.CategoryID)
	
	// Get existing category
	category, err := h.categoryRepo.GetByID(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}
	
	// Check if slug is unique (excluding current category)
	if category.Slug != cmd.Slug {
		exists, err := h.categoryRepo.ExistsBySlug(ctx, cmd.Slug)
		if err != nil {
			return err
		}
		if exists {
			return errors.ErrCategoryAlreadyExists.WithDetails("Category with this slug already exists")
		}
	}
	
	// Verify parent category exists if provided
	if cmd.ParentID != nil {
		if *cmd.ParentID != cmd.CategoryID { // Prevent self-referencing
			_, err = h.categoryRepo.GetByID(ctx, *cmd.ParentID)
			if err != nil {
				return err
			}
		} else {
			return errors.NewBusinessLogicError("INVALID_PARENT", "Category cannot be its own parent")
		}
	}
	
	// Update fields
	category.Name = cmd.Name
	category.Slug = cmd.Slug
	category.Description = cmd.Description
	category.ParentID = cmd.ParentID
	category.ImageURL = cmd.ImageURL
	category.SortOrder = cmd.SortOrder
	category.MetaTitle = cmd.MetaTitle
	category.MetaDesc = cmd.MetaDesc
	
	// Save category
	if err := h.categoryRepo.Update(ctx, category); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated category: %s", category.ID)
	return nil
}

// handleDeleteCategory handles category deletion
func (h *ProductCommandHandler) handleDeleteCategory(ctx context.Context, cmd *commands.DeleteCategoryCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting category: %s", cmd.CategoryID)
	
	// Check if category exists
	category, err := h.categoryRepo.GetByID(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}
	
	// Check if category has children
	children, err := h.categoryRepo.GetChildren(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.ErrResourceInUse.WithDetails("Category has child categories and cannot be deleted")
	}
	
	// Check if category has products (you might want to implement this check)
	filter := interfaces.ProductFilter{CategoryID: &cmd.CategoryID, PageSize: 1}
	products, err := h.productRepo.List(ctx, filter)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return errors.ErrResourceInUse.WithDetails("Category has products and cannot be deleted")
	}
	
	// Delete category (soft delete)
	if err := h.categoryRepo.Delete(ctx, cmd.CategoryID); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully deleted category: %s", cmd.CategoryID)
	return nil
}
