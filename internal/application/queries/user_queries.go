package queries

import "github.com/google/uuid"

// GetUserByIdQuery represents the query to get a user by their ID.
type GetUserByIdQuery struct {
	ID uuid.UUID
}

func (q *GetUserByIdQuery) GetName() string {
	return "GetUserByIdQuery"
}

// GetUserByEmailQuery represents the query to get a user by their email.
type GetUserByEmailQuery struct {
	Email string
}

func (q *GetUserByEmailQuery) GetName() string {
	return "GetUserByEmailQuery"
}
