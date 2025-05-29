package commands

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateProductCommand represents a product creation command
type CreateProductCommand struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	SKU         string          `json:"sku" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	CategoryID  uuid.UUID       `json:"category_id" validate:"required"`
	Brand       string          `json:"brand"`
	Model       string          `json:"model"`
	Weight      decimal.Decimal `json:"weight"`
	Dimensions  string          `json:"dimensions"`
	Color       string          `json:"color"`
	Material    string          `json:"material"`
	Warranty    int             `json:"warranty"`
	Stock       int             `json:"stock" validate:"min=0"`
	MinStock    int             `json:"min_stock" validate:"min=0"`
	MaxStock    int             `json:"max_stock" validate:"min=1"`
	IsFeatured  bool            `json:"is_featured"`
	MetaTitle   string          `json:"meta_title"`
	MetaDesc    string          `json:"meta_description"`
	Tags        string          `json:"tags"`
}

func (c CreateProductCommand) GetName() string {
	return "CreateProduct"
}

// UpdateProductCommand represents a product update command
type UpdateProductCommand struct {
	ProductID   uuid.UUID       `json:"product_id" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	CategoryID  uuid.UUID       `json:"category_id" validate:"required"`
	Brand       string          `json:"brand"`
	Model       string          `json:"model"`
	Weight      decimal.Decimal `json:"weight"`
	Dimensions  string          `json:"dimensions"`
	Color       string          `json:"color"`
	Material    string          `json:"material"`
	Warranty    int             `json:"warranty"`
	MinStock    int             `json:"min_stock" validate:"min=0"`
	MaxStock    int             `json:"max_stock" validate:"min=1"`
	IsFeatured  bool            `json:"is_featured"`
	MetaTitle   string          `json:"meta_title"`
	MetaDesc    string          `json:"meta_description"`
	Tags        string          `json:"tags"`
}

func (c UpdateProductCommand) GetName() string {
	return "UpdateProduct"
}

// UpdateProductStockCommand represents a product stock update command
type UpdateProductStockCommand struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"min=0"`
	Reason    string    `json:"reason"`
}

func (c UpdateProductStockCommand) GetName() string {
	return "UpdateProductStock"
}

// DeleteProductCommand represents a product deletion command
type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (c DeleteProductCommand) GetName() string {
	return "DeleteProduct"
}

// AddProductImageCommand represents adding a product image command
type AddProductImageCommand struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	URL       string    `json:"url" validate:"required,url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
}

func (c AddProductImageCommand) GetName() string {
	return "AddProductImage"
}

// AddProductAttributeCommand represents adding a product attribute command
type AddProductAttributeCommand struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Value     string    `json:"value" validate:"required"`
	Unit      string    `json:"unit"`
	SortOrder int       `json:"sort_order"`
}

func (c AddProductAttributeCommand) GetName() string {
	return "AddProductAttribute"
}

// CreateCategoryCommand represents a category creation command
type CreateCategoryCommand struct {
	Name        string     `json:"name" validate:"required"`
	Slug        string     `json:"slug" validate:"required"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	ImageURL    string     `json:"image_url"`
	SortOrder   int        `json:"sort_order"`
	MetaTitle   string     `json:"meta_title"`
	MetaDesc    string     `json:"meta_description"`
}

func (c CreateCategoryCommand) GetName() string {
	return "CreateCategory"
}

// UpdateCategoryCommand represents a category update command
type UpdateCategoryCommand struct {
	CategoryID  uuid.UUID  `json:"category_id" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	Slug        string     `json:"slug" validate:"required"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	ImageURL    string     `json:"image_url"`
	SortOrder   int        `json:"sort_order"`
	MetaTitle   string     `json:"meta_title"`
	MetaDesc    string     `json:"meta_description"`
}

func (c UpdateCategoryCommand) GetName() string {
	return "UpdateCategory"
}

// DeleteCategoryCommand represents a category deletion command
type DeleteCategoryCommand struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

func (c DeleteCategoryCommand) GetName() string {
	return "DeleteCategory"
}
