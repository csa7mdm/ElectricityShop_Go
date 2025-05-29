package mediator

import "context"

// Mediator defines the interface for sending commands and queries.
type Mediator interface {
	Send(ctx context.Context, command Command) error
	Query(ctx context.Context, query Query) (interface{}, error)
}

// Command represents an action that changes the state of the system.
type Command interface {
	GetName() string // GetName returns the name of the command.
}

// Query represents a request for information.
type Query interface {
	GetName() string // GetName returns the name of the query.
}

// CommandHandler defines the interface for handling commands.
type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

// QueryHandler defines the interface for handling queries.
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}
