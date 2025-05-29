package handlers

import (
	"context"
	"fmt" // For errors

	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	// "github.com/yourusername/electricity-shop-go/internal/domain/entities" // Not directly needed if dtos.UserResponse takes basic types. Actually it is for user.Role etc.
	domainInterfaces "github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/mediator" // For mediator.Query
)

// GetUserByIdQueryHandler handles GetUserByIdQuery.
type GetUserByIdQueryHandler struct {
	userRepository domainInterfaces.UserRepository
}

// NewGetUserByIdQueryHandler creates a new GetUserByIdQueryHandler.
func NewGetUserByIdQueryHandler(userRepo domainInterfaces.UserRepository) *GetUserByIdQueryHandler {
	return &GetUserByIdQueryHandler{userRepository: userRepo}
}

func (h *GetUserByIdQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	cmd, ok := query.(*queries.GetUserByIdQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetUserByIdQueryHandler")
	}

	user, err := h.userRepository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by ID: %w", err) // Propagate error
	}
	if user == nil {
		// Consider returning a domain-specific error like pkg/errors.ErrUserNotFound
		return nil, fmt.Errorf("user not found")
	}

	response := &dtos.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return response, nil
}

// GetUserByEmailQueryHandler handles GetUserByEmailQuery.
type GetUserByEmailQueryHandler struct {
	userRepository domainInterfaces.UserRepository
}

// NewGetUserByEmailQueryHandler creates a new GetUserByEmailQueryHandler.
func NewGetUserByEmailQueryHandler(userRepo domainInterfaces.UserRepository) *GetUserByEmailQueryHandler {
	return &GetUserByEmailQueryHandler{userRepository: userRepo}
}

func (h *GetUserByEmailQueryHandler) Handle(ctx context.Context, query mediator.Query) (interface{}, error) {
	cmd, ok := query.(*queries.GetUserByEmailQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetUserByEmailQueryHandler")
	}

	user, err := h.userRepository.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %w", err) // Propagate error
	}
	if user == nil {
		// Consider returning a domain-specific error like pkg/errors.ErrUserNotFound
		return nil, fmt.Errorf("user not found")
	}

	response := &dtos.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return response, nil
}
