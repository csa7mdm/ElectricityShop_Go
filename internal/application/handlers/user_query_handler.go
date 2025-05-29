package handlers

import (
	"context"

	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// UserQueryHandler handles user-related queries
type UserQueryHandler struct {
	userRepo    interfaces.UserRepository
	addressRepo interfaces.AddressRepository
	logger      logger.Logger
}

// NewUserQueryHandler creates a new UserQueryHandler
func NewUserQueryHandler(
	userRepo interfaces.UserRepository,
	addressRepo interfaces.AddressRepository,
	logger logger.Logger,
) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		logger:      logger,
	}
}

// Handle handles queries
func (h *UserQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	switch q := query.(type) {
	case *queries.GetUserByIDQuery:
		return h.handleGetUserByID(ctx, q)
	case *queries.GetUserByEmailQuery:
		return h.handleGetUserByEmail(ctx, q)
	case *queries.ListUsersQuery:
		return h.handleListUsers(ctx, q)
	case *queries.GetUserAddressesQuery:
		return h.handleGetUserAddresses(ctx, q)
	default:
		return nil, errors.New("UNSUPPORTED_QUERY", "Unsupported query type", 400)
	}
}

// handleGetUserByID handles getting a user by ID
func (h *UserQueryHandler) handleGetUserByID(ctx context.Context, query *queries.GetUserByIDQuery) (*entities.User, error) {
	h.logger.WithContext(ctx).Debugf("Getting user by ID: %s", query.UserID)
	
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved user: %s", user.ID)
	return user, nil
}

// handleGetUserByEmail handles getting a user by email
func (h *UserQueryHandler) handleGetUserByEmail(ctx context.Context, query *queries.GetUserByEmailQuery) (*entities.User, error) {
	h.logger.WithContext(ctx).Debugf("Getting user by email: %s", query.Email)
	
	user, err := h.userRepo.GetByEmail(ctx, query.Email)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved user: %s", user.ID)
	return user, nil
}

// handleListUsers handles listing users with filtering
func (h *UserQueryHandler) handleListUsers(ctx context.Context, query *queries.ListUsersQuery) ([]*entities.User, error) {
	h.logger.WithContext(ctx).Debugf("Listing users with filter")
	
	users, err := h.userRepo.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d users", len(users))
	return users, nil
}

// handleGetUserAddresses handles getting user addresses
func (h *UserQueryHandler) handleGetUserAddresses(ctx context.Context, query *queries.GetUserAddressesQuery) ([]*entities.Address, error) {
	h.logger.WithContext(ctx).Debugf("Getting addresses for user: %s", query.UserID)
	
	// Verify user exists
	_, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	addresses, err := h.addressRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	h.logger.WithContext(ctx).Debugf("Successfully retrieved %d addresses for user: %s", len(addresses), query.UserID)
	return addresses, nil
}
