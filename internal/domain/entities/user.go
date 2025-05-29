package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"unique;not null;type:varchar(255)" json:"email"`
	Password  string    `gorm:"not null;type:varchar(255)" json:"-"`
	Role      UserRole  `gorm:"not null;type:varchar(50);default:'customer'" json:"role"`
	FirstName string    `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100)" json:"last_name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	Orders    []Order   `gorm:"foreignKey:UserID" json:"-"`
	Cart      Cart      `gorm:"foreignKey:UserID" json:"-"`
	Addresses []Address `gorm:"foreignKey:UserID" json:"-"`
}

// UserRole represents the role of a user
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
	RoleEmployee UserRole = "employee"
)

// Address represents a shipping/billing address
type Address struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Type         AddressType `gorm:"not null;type:varchar(20)" json:"type"`
	FirstName    string    `gorm:"not null;type:varchar(100)" json:"first_name"`
	LastName     string    `gorm:"not null;type:varchar(100)" json:"last_name"`
	Company      string    `gorm:"type:varchar(100)" json:"company"`
	AddressLine1 string    `gorm:"not null;type:varchar(255)" json:"address_line_1"`
	AddressLine2 string    `gorm:"type:varchar(255)" json:"address_line_2"`
	City         string    `gorm:"not null;type:varchar(100)" json:"city"`
	State        string    `gorm:"not null;type:varchar(100)" json:"state"`
	ZipCode      string    `gorm:"not null;type:varchar(20)" json:"zip_code"`
	Country      string    `gorm:"not null;type:varchar(100);default:'USA'" json:"country"`
	IsDefault    bool      `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeBilling  AddressType = "billing"
	AddressTypeShipping AddressType = "shipping"
)

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Address
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEmployee checks if user has employee role
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// FullAddress returns the full formatted address
func (a *Address) FullAddress() string {
	address := a.AddressLine1
	if a.AddressLine2 != "" {
		address += ", " + a.AddressLine2
	}
	address += ", " + a.City + ", " + a.State + " " + a.ZipCode
	if a.Country != "USA" {
		address += ", " + a.Country
	}
	return address
}
