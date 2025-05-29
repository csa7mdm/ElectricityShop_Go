package entities

import (
	"time"
	"github.com/google/uuid"
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
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      UserRole  `gorm:"not null;type:varchar(50)"` // Added type for DB
	// Profile   UserProfile // UserProfile is not defined yet, commenting out for now.
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TODO: Define UserProfile struct later.
// type UserProfile struct {
//    FirstName string
//    LastName  string
//    Address   string // This might be a more complex Address struct later
//    Phone     string
// }
