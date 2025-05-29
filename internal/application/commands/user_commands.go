package commands

import (
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
)

// RegisterUserCommand represents the command to register a new user.
type RegisterUserCommand struct {
	Email    string
	Password string
	// Role can be added here if clients can specify it, or set by default in handler
}

func (c *RegisterUserCommand) GetName() string {
	return "RegisterUserCommand"
}

// UpdateUserProfileCommand represents a user profile update command
type UpdateUserProfileCommand struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Phone     string    `json:"phone,omitempty"`
}

func (c UpdateUserProfileCommand) GetName() string {
	return "UpdateUserProfile"
}

// DeleteUserCommand represents a user deletion command
type DeleteUserCommand struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (c DeleteUserCommand) GetName() string {
	return "DeleteUser"
}

// AddAddressCommand represents adding an address command
type AddAddressCommand struct {
	UserID       uuid.UUID            `json:"user_id" validate:"required"`
	Type         entities.AddressType `json:"type" validate:"required"`
	FirstName    string               `json:"first_name" validate:"required"`
	LastName     string               `json:"last_name" validate:"required"`
	Company      string               `json:"company,omitempty"`
	AddressLine1 string               `json:"address_line_1" validate:"required"`
	AddressLine2 string               `json:"address_line_2,omitempty"`
	City         string               `json:"city" validate:"required"`
	State        string               `json:"state" validate:"required"`
	ZipCode      string               `json:"zip_code" validate:"required"`
	Country      string               `json:"country" validate:"required"`
	IsDefault    bool                 `json:"is_default"`
}

func (c AddAddressCommand) GetName() string {
	return "AddAddress"
}

// UpdateAddressCommand represents updating an address command
type UpdateAddressCommand struct {
	AddressID    uuid.UUID            `json:"address_id" validate:"required"`
	UserID       uuid.UUID            `json:"user_id" validate:"required"`
	Type         entities.AddressType `json:"type" validate:"required"`
	FirstName    string               `json:"first_name" validate:"required"`
	LastName     string               `json:"last_name" validate:"required"`
	Company      string               `json:"company,omitempty"`
	AddressLine1 string               `json:"address_line_1" validate:"required"`
	AddressLine2 string               `json:"address_line_2,omitempty"`
	City         string               `json:"city" validate:"required"`
	State        string               `json:"state" validate:"required"`
	ZipCode      string               `json:"zip_code" validate:"required"`
	Country      string               `json:"country" validate:"required"`
	IsDefault    bool                 `json:"is_default"`
}

func (c UpdateAddressCommand) GetName() string {
	return "UpdateAddress"
}

// DeleteAddressCommand represents deleting an address command
type DeleteAddressCommand struct {
	AddressID uuid.UUID `json:"address_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
}

func (c DeleteAddressCommand) GetName() string {
	return "DeleteAddress"
}

// LoginUserCommand represents the command to log in a user.
type LoginUserCommand struct {
	Email    string
	Password string
}

func (c *LoginUserCommand) GetName() string {
	return "LoginUserCommand"
}
