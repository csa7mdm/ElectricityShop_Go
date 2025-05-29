package mediator

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

// Command represents a command in CQRS
type Command interface {
	GetName() string
}

// Query represents a query in CQRS
type Query interface {
	GetName() string
}

// CommandHandler handles commands
type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

// QueryHandler handles queries
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

// Mediator interface defines the contract for the mediator
type Mediator interface {
	RegisterCommandHandler(command Command, handler CommandHandler)
	RegisterQueryHandler(query Query, handler QueryHandler)
	Send(ctx context.Context, command Command) error
	Query(ctx context.Context, query Query) (interface{}, error)
}

// DefaultMediator is the default implementation of Mediator
type DefaultMediator struct {
	commandHandlers map[string]CommandHandler
	queryHandlers   map[string]QueryHandler
	logger          logger.Logger
	mutex           sync.RWMutex
}

// NewMediator creates a new instance of DefaultMediator
func NewMediator(logger logger.Logger) *DefaultMediator {
	return &DefaultMediator{
		commandHandlers: make(map[string]CommandHandler),
		queryHandlers:   make(map[string]QueryHandler),
		logger:          logger,
	}
}

// RegisterCommandHandler registers a command handler
func (m *DefaultMediator) RegisterCommandHandler(command Command, handler CommandHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	commandName := command.GetName()
	if commandName == "" {
		commandName = reflect.TypeOf(command).String()
	}
	
	m.commandHandlers[commandName] = handler
	m.logger.Infof("Registered command handler for: %s", commandName)
}

// RegisterQueryHandler registers a query handler
func (m *DefaultMediator) RegisterQueryHandler(query Query, handler QueryHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	queryName := query.GetName()
	if queryName == "" {
		queryName = reflect.TypeOf(query).String()
	}
	
	m.queryHandlers[queryName] = handler
	m.logger.Infof("Registered query handler for: %s", queryName)
}

// Send handles a command
func (m *DefaultMediator) Send(ctx context.Context, command Command) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	commandName := command.GetName()
	if commandName == "" {
		commandName = reflect.TypeOf(command).String()
	}
	
	handler, exists := m.commandHandlers[commandName]
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", commandName)
	}
	
	m.logger.Debugf("Handling command: %s", commandName)
	
	err := handler.Handle(ctx, command)
	if err != nil {
		m.logger.Errorf("Error handling command %s: %v", commandName, err)
		return fmt.Errorf("error handling command %s: %w", commandName, err)
	}
	
	m.logger.Debugf("Successfully handled command: %s", commandName)
	return nil
}

// Query handles a query
func (m *DefaultMediator) Query(ctx context.Context, query Query) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	queryName := query.GetName()
	if queryName == "" {
		queryName = reflect.TypeOf(query).String()
	}
	
	handler, exists := m.queryHandlers[queryName]
	if !exists {
		return nil, fmt.Errorf("no handler registered for query: %s", queryName)
	}
	
	m.logger.Debugf("Handling query: %s", queryName)
	
	result, err := handler.Handle(ctx, query)
	if err != nil {
		m.logger.Errorf("Error handling query %s: %v", queryName, err)
		return nil, fmt.Errorf("error handling query %s: %w", queryName, err)
	}
	
	m.logger.Debugf("Successfully handled query: %s", queryName)
	return result, nil
}

// GetRegisteredCommands returns the list of registered command names
func (m *DefaultMediator) GetRegisteredCommands() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	commands := make([]string, 0, len(m.commandHandlers))
	for commandName := range m.commandHandlers {
		commands = append(commands, commandName)
	}
	return commands
}

// GetRegisteredQueries returns the list of registered query names
func (m *DefaultMediator) GetRegisteredQueries() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	queries := make([]string, 0, len(m.queryHandlers))
	for queryName := range m.queryHandlers {
		queries = append(queries, queryName)
	}
	return queries
}

// Middleware for handling cross-cutting concerns
type Middleware func(ctx context.Context, request interface{}, next func(ctx context.Context, request interface{}) (interface{}, error)) (interface{}, error)

// EnhancedMediator extends DefaultMediator with middleware support
type EnhancedMediator struct {
	*DefaultMediator
	middlewares []Middleware
}

// NewEnhancedMediator creates a new enhanced mediator with middleware support
func NewEnhancedMediator(logger logger.Logger) *EnhancedMediator {
	return &EnhancedMediator{
		DefaultMediator: NewMediator(logger),
		middlewares:     make([]Middleware, 0),
	}
}

// Use adds a middleware to the mediator
func (m *EnhancedMediator) Use(middleware Middleware) {
	m.middlewares = append(m.middlewares, middleware)
}

// Send handles a command with middleware
func (m *EnhancedMediator) Send(ctx context.Context, command Command) error {
	if len(m.middlewares) == 0 {
		return m.DefaultMediator.Send(ctx, command)
	}
	
	handler := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, m.DefaultMediator.Send(ctx, request.(Command))
	}
	
	// Chain middlewares
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		middleware := m.middlewares[i]
		next := handler
		handler = func(ctx context.Context, request interface{}) (interface{}, error) {
			return middleware(ctx, request, next)
		}
	}
	
	_, err := handler(ctx, command)
	return err
}

// Query handles a query with middleware
func (m *EnhancedMediator) Query(ctx context.Context, query Query) (interface{}, error) {
	if len(m.middlewares) == 0 {
		return m.DefaultMediator.Query(ctx, query)
	}
	
	handler := func(ctx context.Context, request interface{}) (interface{}, error) {
		return m.DefaultMediator.Query(ctx, request.(Query))
	}
	
	// Chain middlewares
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		middleware := m.middlewares[i]
		next := handler
		handler = func(ctx context.Context, request interface{}) (interface{}, error) {
			return middleware(ctx, request, next)
		}
	}
	
	return handler(ctx, query)
}

// Common middleware implementations

// LoggingMiddleware logs all requests and responses
func LoggingMiddleware(logger logger.Logger) Middleware {
	return func(ctx context.Context, request interface{}, next func(ctx context.Context, request interface{}) (interface{}, error)) (interface{}, error) {
		requestType := reflect.TypeOf(request).String()
		logger.Infof("Processing request: %s", requestType)
		
		result, err := next(ctx, request)
		
		if err != nil {
			logger.Errorf("Request %s failed: %v", requestType, err)
		} else {
			logger.Infof("Request %s completed successfully", requestType)
		}
		
		return result, err
	}
}

// ValidationMiddleware validates requests before processing
func ValidationMiddleware(validator func(interface{}) error) Middleware {
	return func(ctx context.Context, request interface{}, next func(ctx context.Context, request interface{}) (interface{}, error)) (interface{}, error) {
		if err := validator(request); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		
		return next(ctx, request)
	}
}

// CachingMiddleware provides caching for queries
func CachingMiddleware(cache map[string]interface{}) Middleware {
	return func(ctx context.Context, request interface{}, next func(ctx context.Context, request interface{}) (interface{}, error)) (interface{}, error) {
		// Only cache queries, not commands
		if _, isQuery := request.(Query); !isQuery {
			return next(ctx, request)
		}
		
		requestType := reflect.TypeOf(request).String()
		cacheKey := fmt.Sprintf("%s:%+v", requestType, request)
		
		// Check cache
		if cachedResult, exists := cache[cacheKey]; exists {
			return cachedResult, nil
		}
		
		// Execute and cache result
		result, err := next(ctx, request)
		if err == nil {
			cache[cacheKey] = result
		}
		
		return result, err
	}
}
