package queries

import (
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
)

// GetUserByIDQuery represents a query to get a user by ID
type GetUserByIDQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (q GetUserByIDQuery) GetName() string {
	return "GetUserByID"
}

// GetUserByEmailQuery represents a query to get a user by email
type GetUserByEmailQuery struct {
	Email string `json:"email" validate:"required,email"`
}

func (q GetUserByEmailQuery) GetName() string {
	return "GetUserByEmail"
}

// ListUsersQuery represents a query to list users with filtering
type ListUsersQuery struct {
	Filter interfaces.UserFilter `json:"filter"`
}

func (q ListUsersQuery) GetName() string {
	return "ListUsers"
}

// GetUserAddressesQuery represents a query to get user addresses
type GetUserAddressesQuery struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (q GetUserAddressesQuery) GetName() string {
	return "GetUserAddresses"
}
