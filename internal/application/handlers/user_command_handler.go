package handlers

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// UserCommandHandler handles user-related commands
type UserCommandHandler struct {
	userRepo      interfaces.UserRepository
	addressRepo   interfaces.AddressRepository
	eventPublisher interfaces.EventPublisher
	logger        logger.Logger
}

// NewUserCommandHandler creates a new UserCommandHandler
func NewUserCommandHandler(
	userRepo interfaces.UserRepository,
	addressRepo interfaces.AddressRepository,
	eventPublisher interfaces.EventPublisher,
	logger logger.Logger,
) *UserCommandHandler {
	return &UserCommandHandler{
		userRepo:      userRepo,
		addressRepo:   addressRepo,
		eventPublisher: eventPublisher,
		logger:        logger,
	}
}

// Handle handles commands
func (h *UserCommandHandler) Handle(ctx context.Context, command mediator.Command) error {
	switch cmd := command.(type) {
	case *commands.RegisterUserCommand:
		return h.handleRegisterUser(ctx, cmd)
	case *commands.UpdateUserProfileCommand:
		return h.handleUpdateUserProfile(ctx, cmd)
	case *commands.DeleteUserCommand:
		return h.handleDeleteUser(ctx, cmd)
	case *commands.AddAddressCommand:
		return h.handleAddAddress(ctx, cmd)
	case *commands.UpdateAddressCommand:
		return h.handleUpdateAddress(ctx, cmd)
	case *commands.DeleteAddressCommand:
		return h.handleDeleteAddress(ctx, cmd)
	default:
		return errors.New("UNSUPPORTED_COMMAND", "Unsupported command type", 400)
	}
}

// handleRegisterUser handles user registration
func (h *UserCommandHandler) handleRegisterUser(ctx context.Context, cmd *commands.RegisterUserCommand) error {
	h.logger.WithContext(ctx).Infof("Registering user with email: %s", cmd.Email)
	
	// Check if user already exists
	exists, err := h.userRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUserAlreadyExists.WithDetails("User with this email already exists")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "BCRYPT_ERROR", "Failed to hash password", 500)
	}
	
	// Set default role if not provided
	role := cmd.Role
	if role == "" {
		role = entities.RoleCustomer
	}
	
	// Create user entity
	user := &entities.User{
		Email:     cmd.Email,
		Password:  string(hashedPassword),
		Role:      role,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Phone:     cmd.Phone,
		IsActive:  true,
	}
	
	// Save user
	if err := h.userRepo.Create(ctx, user); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewUserRegisteredEvent(
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		string(user.Role),
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish UserRegisteredEvent: %v", err)
		// Don't fail the command for event publishing errors
	}
	
	h.logger.WithContext(ctx).Infof("Successfully registered user: %s", user.ID)
	return nil
}

// handleUpdateUserProfile handles user profile updates
func (h *UserCommandHandler) handleUpdateUserProfile(ctx context.Context, cmd *commands.UpdateUserProfileCommand) error {
	h.logger.WithContext(ctx).Infof("Updating user profile: %s", cmd.UserID)
	
	// Get existing user
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Update fields
	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName
	user.Phone = cmd.Phone
	
	// Save user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewUserProfileUpdatedEvent(
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish UserProfileUpdatedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated user profile: %s", user.ID)
	return nil
}

// handleDeleteUser handles user deletion
func (h *UserCommandHandler) handleDeleteUser(ctx context.Context, cmd *commands.DeleteUserCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting user: %s", cmd.UserID)
	
	// Check if user exists
	_, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Delete user (soft delete)
	if err := h.userRepo.Delete(ctx, cmd.UserID); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully deleted user: %s", cmd.UserID)
	return nil
}

// handleAddAddress handles adding an address
func (h *UserCommandHandler) handleAddAddress(ctx context.Context, cmd *commands.AddAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Adding address for user: %s", cmd.UserID)
	
	// Verify user exists
	_, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Create address entity
	address := &entities.Address{
		UserID:       cmd.UserID,
		Type:         cmd.Type,
		FirstName:    cmd.FirstName,
		LastName:     cmd.LastName,
		Company:      cmd.Company,
		AddressLine1: cmd.AddressLine1,
		AddressLine2: cmd.AddressLine2,
		City:         cmd.City,
		State:        cmd.State,
		ZipCode:      cmd.ZipCode,
		Country:      cmd.Country,
		IsDefault:    cmd.IsDefault,
	}
	
	// If this is set as default, unset other defaults first
	if cmd.IsDefault {
		if err := h.addressRepo.SetAsDefault(ctx, uuid.Nil, cmd.Type); err != nil {
			return err
		}
	}
	
	// Save address
	if err := h.addressRepo.Create(ctx, address); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully added address for user: %s", cmd.UserID)
	return nil
}

// handleUpdateAddress handles updating an address
func (h *UserCommandHandler) handleUpdateAddress(ctx context.Context, cmd *commands.UpdateAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Updating address: %s", cmd.AddressID)
	
	// Get existing address and verify ownership
	address, err := h.addressRepo.GetByID(ctx, cmd.AddressID)
	if err != nil {
		return err
	}
	
	if address.UserID != cmd.UserID {
		return errors.ErrForbidden.WithDetails("You don't have permission to update this address")
	}
	
	// Update fields
	address.Type = cmd.Type
	address.FirstName = cmd.FirstName
	address.LastName = cmd.LastName
	address.Company = cmd.Company
	address.AddressLine1 = cmd.AddressLine1
	address.AddressLine2 = cmd.AddressLine2
	address.City = cmd.City
	address.State = cmd.State
	address.ZipCode = cmd.ZipCode
	address.Country = cmd.Country
	address.IsDefault = cmd.IsDefault
	
	// If this is set as default, unset other defaults first
	if cmd.IsDefault && !address.IsDefault {
		if err := h.addressRepo.SetAsDefault(ctx, cmd.AddressID, cmd.Type); err != nil {
			return err
		}
	}
	
	// Save address
	if err := h.addressRepo.Update(ctx, address); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated address: %s", cmd.AddressID)
	return nil
}

// handleDeleteAddress handles deleting an address
func (h *UserCommandHandler) handleDeleteAddress(ctx context.Context, cmd *commands.DeleteAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting address: %s", cmd.AddressID)
	
	// Get existing address and verify ownership
	address, err := h.addressRepo.GetByID(ctx, cmd.AddressID)
	if err != nil {
		return err
	}
	
	if address.UserID != cmd.UserID {
		return errors.ErrForbidden.WithDetails("You don't have permission to delete this address")
	}
	
	// Delete address
	if err := h.addressRepo.Delete(ctx, cmd.AddressID); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully deleted address: %s", cmd.AddressID)
	return nil
}
