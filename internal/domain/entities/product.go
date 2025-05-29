package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"not null;type:varchar(255)" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	SKU         string          `gorm:"unique;not null;type:varchar(100)" json:"sku"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Brand       string          `gorm:"type:varchar(100)" json:"brand"`
	Model       string          `gorm:"type:varchar(100)" json:"model"`
	Weight      decimal.Decimal `gorm:"type:decimal(8,3)" json:"weight"`
	Dimensions  string          `gorm:"type:varchar(100)" json:"dimensions"`
	Color       string          `gorm:"type:varchar(50)" json:"color"`
	Material    string          `gorm:"type:varchar(100)" json:"material"`
	Warranty    int             `gorm:"default:0" json:"warranty"` // in months
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	MinStock    int             `gorm:"default:5" json:"min_stock"`
	MaxStock    int             `gorm:"default:1000" json:"max_stock"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	IsFeatured  bool            `gorm:"default:false" json:"is_featured"`
	MetaTitle   string          `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDesc    string          `gorm:"type:text" json:"meta_description"`
	Tags        string          `gorm:"type:text" json:"tags"` // JSON array of tags
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
	
	// Relationships
	Category     Category        `gorm:"foreignKey:CategoryID" json:"category"`
	Images       []ProductImage  `gorm:"foreignKey:ProductID" json:"images"`
	Attributes   []ProductAttribute `gorm:"foreignKey:ProductID" json:"attributes"`
	Reviews      []ProductReview `gorm:"foreignKey:ProductID" json:"reviews"`
	CartItems    []CartItem      `gorm:"foreignKey:ProductID" json:"-"`
	OrderItems   []OrderItem     `gorm:"foreignKey:ProductID" json:"-"`
}

// Category represents a product category
type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null;type:varchar(255)" json:"name"`
	Slug        string    `gorm:"unique;not null;type:varchar(255)" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	MetaTitle   string    `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDesc    string    `gorm:"type:text" json:"meta_description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships  
	Parent     *Category  `gorm:"foreignKey:ParentID" json:"parent"`
	Children   []Category `gorm:"foreignKey:ParentID" json:"children"`
	Products   []Product  `gorm:"foreignKey:CategoryID" json:"products"`
}

// ProductImage represents product images
type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	URL       string    `gorm:"not null;type:varchar(500)" json:"url"`
	AltText   string    `gorm:"type:varchar(255)" json:"alt_text"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Product Product `gorm:"foreignKey:ProductID" json:"-"`
}

// ProductAttribute represents product attributes (like specifications)
type ProductAttribute struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Name      string    `gorm:"not null;type:varchar(100)" json:"name"`
	Value     string    `gorm:"not null;type:varchar(255)" json:"value"`
	Unit      string    `gorm:"type:varchar(50)" json:"unit"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Product Product `gorm:"foreignKey:ProductID" json:"-"`
}

// ProductReview represents customer reviews
type ProductReview struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Comment   string    `gorm:"type:text" json:"comment"`
	IsVerified bool     `gorm:"default:false" json:"is_verified"` // verified purchase
	IsApproved bool     `gorm:"default:false" json:"is_approved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	Product Product `gorm:"foreignKey:ProductID" json:"-"`
	User    User    `gorm:"foreignKey:UserID" json:"user"`
}

// BeforeCreate hooks
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (pi *ProductImage) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}
	return nil
}

func (pa *ProductAttribute) BeforeCreate(tx *gorm.DB) error {
	if pa.ID == uuid.Nil {
		pa.ID = uuid.New()
	}
	return nil
}

func (pr *ProductReview) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

// Business logic methods
func (p *Product) IsInStock() bool {
	return p.Stock > 0
}

func (p *Product) IsLowStock() bool {
	return p.Stock <= p.MinStock && p.Stock > 0
}

func (p *Product) CanOrder(quantity int) bool {
	return p.IsActive && p.Stock >= quantity
}

func (p *Product) GetPrimaryImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsPrimary {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

func (p *Product) GetAverageRating() float64 {
	if len(p.Reviews) == 0 {
		return 0.0
	}
	
	total := 0
	count := 0
	for _, review := range p.Reviews {
		if review.IsApproved {
			total += review.Rating
			count++
		}
	}
	
	if count == 0 {
		return 0.0
	}
	
	return float64(total) / float64(count)
}

func (p *Product) GetReviewCount() int {
	count := 0
	for _, review := range p.Reviews {
		if review.IsApproved {
			count++
		}
	}
	return count
}

func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}
