package dtos

import (
	"time"
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
)

// Address DTOs

// AddAddressRequest represents the request to add a new address
type AddAddressRequest struct {
	Type       string `json:"type" validate:"required,oneof=home work billing shipping other"`
	Street     string `json:"street" validate:"required,min=1,max=255"`
	City       string `json:"city" validate:"required,min=1,max=100"`
	State      string `json:"state" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"max=20"`
	Country    string `json:"country" validate:"required,min=2,max=100"`
	IsDefault  bool   `json:"is_default"`
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	Type       string `json:"type" validate:"required,oneof=home work billing shipping other"`
	Street     string `json:"street" validate:"required,min=1,max=255"`
	City       string `json:"city" validate:"required,min=1,max=100"`
	State      string `json:"state" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"max=20"`
	Country    string `json:"country" validate:"required,min=2,max=100"`
	IsDefault  bool   `json:"is_default"`
}

// AddressResponse represents an address in responses
type AddressResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// User Profile DTOs

// UpdateUserProfileRequest represents the request to update user profile
type UpdateUserProfileRequest struct {
	// Add more fields as needed when UserProfile is implemented
	// FirstName string `json:"first_name" validate:"max=100"`
	// LastName  string `json:"last_name" validate:"max=100"`
	// Phone     string `json:"phone" validate:"max=20"`
}

// UserProfileResponse represents user profile in responses
type UserProfileResponse struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Role      string            `json:"role"`
	IsActive  bool              `json:"is_active"`
	Addresses []AddressResponse `json:"addresses,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Conversion methods

// ToAddressEntity converts AddAddressRequest to Address entity
func (req *AddAddressRequest) ToAddressEntity(userID uuid.UUID) *entities.Address {
	return &entities.Address{
		UserID:     userID,
		Type:       entities.AddressType(req.Type),
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		IsDefault:  req.IsDefault,
	}
}

// FromAddressEntity converts Address entity to AddressResponse
func FromAddressEntity(address *entities.Address) *AddressResponse {
	return &AddressResponse{
		ID:         address.ID.String(),
		Type:       string(address.Type),
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
		IsDefault:  address.IsDefault,
		CreatedAt:  address.CreatedAt,
		UpdatedAt:  address.UpdatedAt,
	}
}

// FromUserEntity converts User entity to UserProfileResponse
func FromUserEntity(user *entities.User) *UserProfileResponse {
	response := &UserProfileResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Convert addresses
	if len(user.Addresses) > 0 {
		response.Addresses = make([]AddressResponse, len(user.Addresses))
		for i, addr := range user.Addresses {
			response.Addresses[i] = *FromAddressEntity(&addr)
		}
	}

	return response
}

// Pagination DTOs

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination information in responses
type PaginationResponse struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users      []UserProfileResponse `json:"users"`
	Pagination PaginationResponse    `json:"pagination"`
}

// ListAddressesResponse represents the response for listing addresses
type ListAddressesResponse struct {
	Addresses []AddressResponse `json:"addresses"`
}
