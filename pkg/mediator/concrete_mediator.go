package mediator

import (
	"context"
	"fmt"
	"sync"
)

// ConcreteMediator is a concrete implementation of the Mediator interface.
type ConcreteMediator struct {
	commandHandlers map[string]CommandHandler
	queryHandlers   map[string]QueryHandler
	mu              sync.RWMutex // To make concurrent access to handlers safe
}

// NewConcreteMediator creates a new instance of ConcreteMediator.
func NewConcreteMediator() *ConcreteMediator {
	return &ConcreteMediator{
		commandHandlers: make(map[string]CommandHandler),
		queryHandlers:   make(map[string]QueryHandler),
	}
}

// RegisterCommandHandler registers a command handler for a given command name.
func (m *ConcreteMediator) RegisterCommandHandler(commandName string, handler CommandHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.commandHandlers[commandName]; ok {
		return fmt.Errorf("command handler already registered for %s", commandName)
	}
	m.commandHandlers[commandName] = handler
	return nil
}

// RegisterQueryHandler registers a query handler for a given query name.
func (m *ConcreteMediator) RegisterQueryHandler(queryName string, handler QueryHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.queryHandlers[queryName]; ok {
		return fmt.Errorf("query handler already registered for %s", queryName)
	}
	m.queryHandlers[queryName] = handler
	return nil
}

// Send dispatches a command to its registered handler.
func (m *ConcreteMediator) Send(ctx context.Context, command Command) error {
	m.mu.RLock()
	handler, ok := m.commandHandlers[command.GetName()]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no command handler registered for %s", command.GetName())
	}
	return handler.Handle(ctx, command)
}

// Query dispatches a query to its registered handler.
func (m *ConcreteMediator) Query(ctx context.Context, query Query) (interface{}, error) {
	m.mu.RLock()
	handler, ok := m.queryHandlers[query.GetName()]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no query handler registered for %s", query.GetName())
	}
	return handler.Handle(ctx, query)
}
