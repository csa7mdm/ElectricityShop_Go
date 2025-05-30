package entities

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole defines the type for user roles.
type UserRole string

// Constants for UserRole.
const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      UserRole  `gorm:"not null;type:varchar(50)" json:"role"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TODO: Define UserProfile struct later.
// type UserProfile struct {
//    FirstName string
//    LastName  string
//    Address   string // This might be a more complex Address struct later
//    Phone     string
// }
