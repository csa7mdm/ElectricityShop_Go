package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Category represents a product category.
type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"not null"`
	Slug      string    `gorm:"unique;not null"`
	Products  []Product `gorm:"foreignKey:CategoryID"` // Define relationship if applicable
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Product represents an electrical product in the shop.
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"not null"`
	Description string
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	SKU         string          `gorm:"unique;not null"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null"`
	Category    Category        // GORM will handle this relationship via CategoryID
	Stock       int             `gorm:"not null"`
	// Images      []ProductImage // ProductImage is not defined yet, commenting out.
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
