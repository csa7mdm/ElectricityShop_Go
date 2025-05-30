package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddressType defines the type of address
type AddressType string

const (
	AddressTypeHome     AddressType = "home"
	AddressTypeWork     AddressType = "work"
	AddressTypeBilling  AddressType = "billing"
	AddressTypeShipping AddressType = "shipping"
	AddressTypeOther    AddressType = "other"
)

// Address represents a user's address
type Address struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	Type       AddressType `gorm:"type:varchar(50);not null" json:"type"`
	Street     string      `gorm:"type:varchar(255);not null" json:"street"`
	City       string      `gorm:"type:varchar(100);not null" json:"city"`
	State      string      `gorm:"type:varchar(100)" json:"state"`
	PostalCode string      `gorm:"type:varchar(20)" json:"postal_code"`
	Country    string      `gorm:"type:varchar(100);not null;default:'US'" json:"country"`
	IsDefault  bool        `gorm:"default:false" json:"is_default"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EmbeddableAddress for embedding in other entities (like Order)
type EmbeddableAddress struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// FormatOneLine returns the address formatted as a single line
func (a *Address) FormatOneLine() string {
	address := a.Street
	if a.City != "" {
		address += ", " + a.City
	}
	if a.State != "" {
		address += ", " + a.State
	}
	if a.PostalCode != "" {
		address += " " + a.PostalCode
	}
	if a.Country != "" {
		address += ", " + a.Country
	}
	return address
}

// ToEmbeddable converts Address to EmbeddableAddress
func (a *Address) ToEmbeddable() EmbeddableAddress {
	return EmbeddableAddress{
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		Country:    a.Country,
	}
}

// IsComplete checks if all required fields are filled
func (a *Address) IsComplete() bool {
	return a.Street != "" && a.City != "" && a.Country != ""
}
