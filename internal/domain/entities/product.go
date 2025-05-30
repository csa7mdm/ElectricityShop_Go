package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Category represents a product category.
type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Slug        string    `gorm:"unique;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	MetaTitle   string    `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDesc    string    `gorm:"type:varchar(500)" json:"meta_desc"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// Product represents an electrical product in the shop.
type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	SKU         string          `gorm:"unique;not null" json:"sku"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Brand       string          `gorm:"type:varchar(100)" json:"brand"`
	Model       string          `gorm:"type:varchar(100)" json:"model"`
	Weight      *decimal.Decimal `gorm:"type:decimal(8,2)" json:"weight"`
	Dimensions  string          `gorm:"type:varchar(100)" json:"dimensions"`
	Color       string          `gorm:"type:varchar(50)" json:"color"`
	Material    string          `gorm:"type:varchar(100)" json:"material"`
	Warranty    string          `gorm:"type:varchar(100)" json:"warranty"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	MinStock    int             `gorm:"default:0" json:"min_stock"`
	MaxStock    int             `gorm:"default:1000" json:"max_stock"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	IsFeatured  bool            `gorm:"default:false" json:"is_featured"`
	MetaTitle   string          `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDesc    string          `gorm:"type:varchar(500)" json:"meta_desc"`
	Tags        string          `gorm:"type:text" json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
	
	// Relationships
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hooks
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Business logic methods

// CanOrder checks if the product can be ordered
func (p *Product) CanOrder(quantity int) bool {
	return p.IsActive && p.Stock >= quantity && quantity > 0
}

// IsLowStock checks if the product is low on stock
func (p *Product) IsLowStock() bool {
	return p.Stock <= p.MinStock
}

// IsOutOfStock checks if the product is out of stock
func (p *Product) IsOutOfStock() bool {
	return p.Stock <= 0
}

// GetAvailableStock returns the available stock quantity
func (p *Product) GetAvailableStock() int {
	if !p.IsActive {
		return 0
	}
	return p.Stock
}

// TODO: Define ProductImage struct later.
// type ProductImage struct {
//    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
//    ProductID uuid.UUID
//    URL       string `gorm:"not null"`
//    AltText   string
//    IsDefault bool
//    CreatedAt time.Time
// }
