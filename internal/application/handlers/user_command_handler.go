package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// UserCommandHandler handles user-related commands
type UserCommandHandler struct {
	userRepo       interfaces.UserRepository
	addressRepo    interfaces.AddressRepository
	eventPublisher interfaces.EventPublisher
	authService    *auth.AuthService
	logger         logger.Logger
}

// NewUserCommandHandler creates a new UserCommandHandler
func NewUserCommandHandler(
	userRepo interfaces.UserRepository,
	addressRepo interfaces.AddressRepository,
	eventPublisher interfaces.EventPublisher,
	authService *auth.AuthService,
	logger logger.Logger,
) *UserCommandHandler {
	return &UserCommandHandler{
		userRepo:       userRepo,
		addressRepo:    addressRepo,
		eventPublisher: eventPublisher,
		authService:    authService,
		logger:         logger,
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

// HandleQuery handles queries that return data (like login)
func (h *UserCommandHandler) HandleQuery(ctx context.Context, query mediator.Query) (interface{}, error) {
	switch q := query.(type) {
	case *commands.LoginUserCommand:
		return h.handleLoginUser(ctx, q)
	default:
		return nil, errors.New("UNSUPPORTED_QUERY", "Unsupported query type", 400)
	}
}

// handleRegisterUser handles user registration
func (h *UserCommandHandler) handleRegisterUser(ctx context.Context, cmd *commands.RegisterUserCommand) error {
	h.logger.WithContext(ctx).Infof("Registering user with email: %s", cmd.Email)

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.ErrUserAlreadyExists.WithDetails("User with this email already exists")
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(cmd.Password)
	if err != nil {
		return errors.Wrap(err, "PASSWORD_HASH_ERROR", "Failed to hash password", 500)
	}

	// Create user entity
	user := &entities.User{
		ID:        uuid.New(),
		Email:     cmd.Email,
		Password:  hashedPassword,
		Role:      entities.RoleCustomer, // Default role
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	if err := h.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Publish domain event
	event := events.NewUserRegisteredEvent(
		user.ID,
		user.Email,
		string(user.Role),
	)

	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish UserRegisteredEvent: %v", err)
		// Don't fail the command for event publishing errors
	}

	h.logger.WithContext(ctx).Infof("Successfully registered user: %s", user.ID)
	return nil
}

// handleLoginUser handles user login
func (h *UserCommandHandler) handleLoginUser(ctx context.Context, cmd *commands.LoginUserCommand) (*dtos.LoginUserResponse, error) {
	h.logger.WithContext(ctx).Infof("User login attempt for email: %s", cmd.Email)

	// Get user by email
	user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.ErrUserInactive
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.Password, cmd.Password); err != nil {
		h.logger.WithContext(ctx).Warnf("Invalid password for user: %s", user.Email)
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.Wrap(err, "TOKEN_GENERATION_ERROR", "Failed to generate token", 500)
	}

	// Create response
	response := &dtos.LoginUserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  string(user.Role),
		Token: token,
	}

	h.logger.WithContext(ctx).Infof("Successfully logged in user: %s", user.ID)
	return response, nil
}

// handleUpdateUserProfile handles user profile updates
func (h *UserCommandHandler) handleUpdateUserProfile(ctx context.Context, cmd *commands.UpdateUserProfileCommand) error {
	h.logger.WithContext(ctx).Infof("Updating user profile: %s", cmd.UserID)

	// Get existing user
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Update fields (add more fields as needed based on your UpdateUserProfileCommand)
	// user.FirstName = cmd.FirstName
	// user.LastName = cmd.LastName
	// user.Phone = cmd.Phone
	user.UpdatedAt = time.Now()

	// Save user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return err
	}

	h.logger.WithContext(ctx).Infof("Successfully updated user profile: %s", user.ID)
	return nil
}

// handleDeleteUser handles user deletion
func (h *UserCommandHandler) handleDeleteUser(ctx context.Context, cmd *commands.DeleteUserCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting user: %s", cmd.UserID)

	// Check if user exists
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Soft delete user
	if err := h.userRepo.Delete(ctx, cmd.UserID); err != nil {
		return err
	}

	h.logger.WithContext(ctx).Infof("Successfully deleted user: %s", cmd.UserID)
	return nil
}

// handleAddAddress handles adding an address to a user
func (h *UserCommandHandler) handleAddAddress(ctx context.Context, cmd *commands.AddAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Adding address for user: %s", cmd.UserID)

	// Verify user exists
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Create address entity
	address := &entities.Address{
		ID:         uuid.New(),
		UserID:     cmd.UserID,
		Type:       cmd.Type,
		Street:     cmd.Street,
		City:       cmd.City,
		State:      cmd.State,
		PostalCode: cmd.PostalCode,
		Country:    cmd.Country,
		IsDefault:  cmd.IsDefault,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// If this is being set as default, unset other default addresses
	if cmd.IsDefault {
		if err := h.addressRepo.UnsetDefaultForUser(ctx, cmd.UserID); err != nil {
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

// handleUpdateAddress handles address updates
func (h *UserCommandHandler) handleUpdateAddress(ctx context.Context, cmd *commands.UpdateAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Updating address: %s for user: %s", cmd.AddressID, cmd.UserID)

	// Get existing address
	address, err := h.addressRepo.GetByID(ctx, cmd.AddressID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.ErrAddressNotFound
	}

	// Verify address belongs to user
	if address.UserID != cmd.UserID {
		return errors.ErrUnauthorized.WithDetails("Address does not belong to user")
	}

	// Update fields
	address.Type = cmd.Type
	address.Street = cmd.Street
	address.City = cmd.City
	address.State = cmd.State
	address.PostalCode = cmd.PostalCode
	address.Country = cmd.Country
	address.IsDefault = cmd.IsDefault
	address.UpdatedAt = time.Now()

	// If this is being set as default, unset other default addresses
	if cmd.IsDefault {
		if err := h.addressRepo.UnsetDefaultForUser(ctx, cmd.UserID); err != nil {
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

// handleDeleteAddress handles address deletion
func (h *UserCommandHandler) handleDeleteAddress(ctx context.Context, cmd *commands.DeleteAddressCommand) error {
	h.logger.WithContext(ctx).Infof("Deleting address: %s for user: %s", cmd.AddressID, cmd.UserID)

	// Get existing address
	address, err := h.addressRepo.GetByID(ctx, cmd.AddressID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.ErrAddressNotFound
	}

	// Verify address belongs to user
	if address.UserID != cmd.UserID {
		return errors.ErrUnauthorized.WithDetails("Address does not belong to user")
	}

	// Delete address
	if err := h.addressRepo.Delete(ctx, cmd.AddressID); err != nil {
		return err
	}

	h.logger.WithContext(ctx).Infof("Successfully deleted address: %s", cmd.AddressID)
	return nil
}
