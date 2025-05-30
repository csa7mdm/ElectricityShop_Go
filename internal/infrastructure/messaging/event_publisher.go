package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

// InMemoryEventPublisher is a simple in-memory event publisher for development
// In production, you would replace this with a proper message broker like RabbitMQ, Kafka, etc.
type InMemoryEventPublisher struct {
	logger    logger.Logger
	handlers  map[string][]EventHandler
}

// EventHandler is a function that handles domain events
type EventHandler func(ctx context.Context, event events.DomainEvent) error

// NewInMemoryEventPublisher creates a new in-memory event publisher
func NewInMemoryEventPublisher(logger logger.Logger) interfaces.EventPublisher {
	return &InMemoryEventPublisher{
		logger:   logger,
		handlers: make(map[string][]EventHandler),
	}
}

// Publish publishes a single domain event
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event interface{}) error {
	domainEvent, ok := event.(events.DomainEvent)
	if !ok {
		return fmt.Errorf("event must implement DomainEvent interface")
	}
	
	eventType := domainEvent.GetEventType()
	
	// Log the event
	eventData, _ := json.Marshal(domainEvent.GetEventData())
	p.logger.WithContext(ctx).Infof("Publishing event: %s, AggregateID: %s, Data: %s", 
		eventType, domainEvent.GetAggregateID(), string(eventData))
	
	// Get handlers for this event type
	handlers, exists := p.handlers[eventType]
	if !exists || len(handlers) == 0 {
		p.logger.WithContext(ctx).Debugf("No handlers registered for event type: %s", eventType)
		return nil
	}
	
	// Execute all handlers
	for _, handler := range handlers {
		if err := handler(ctx, domainEvent); err != nil {
			p.logger.WithContext(ctx).Errorf("Error executing handler for event %s: %v", eventType, err)
			// Continue with other handlers even if one fails
			continue
		}
	}
	
	p.logger.WithContext(ctx).Debugf("Successfully published event: %s", eventType)
	return nil
}

// PublishBatch publishes multiple domain events
func (p *InMemoryEventPublisher) PublishBatch(ctx context.Context, events []interface{}) error {
	var errors []error
	
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			errors = append(errors, err)
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("failed to publish %d out of %d events", len(errors), len(events))
	}
	
	return nil
}

// Subscribe registers an event handler for a specific event type
func (p *InMemoryEventPublisher) Subscribe(eventType string, handler EventHandler) {
	if p.handlers[eventType] == nil {
		p.handlers[eventType] = make([]EventHandler, 0)
	}
	
	p.handlers[eventType] = append(p.handlers[eventType], handler)
	p.logger.Infof("Registered handler for event type: %s", eventType)
}

// GetHandlerCount returns the number of handlers registered for an event type
func (p *InMemoryEventPublisher) GetHandlerCount(eventType string) int {
	handlers, exists := p.handlers[eventType]
	if !exists {
		return 0
	}
	return len(handlers)
}

// Example event handlers that you can register

// LoggingEventHandler logs all events
func LoggingEventHandler(logger logger.Logger) EventHandler {
	return func(ctx context.Context, event events.DomainEvent) error {
		eventData, _ := json.Marshal(event.GetEventData())
		logger.WithContext(ctx).Infof("Event logged: Type=%s, AggregateID=%s, Data=%s", 
			event.GetEventType(), event.GetAggregateID(), string(eventData))
		return nil
	}
}

// EmailNotificationHandler sends email notifications for specific events
func EmailNotificationHandler(emailService interfaces.EmailService, logger logger.Logger) EventHandler {
	return func(ctx context.Context, event events.DomainEvent) error {
		switch e := event.(type) {
		case *events.UserRegisteredEvent:
			return emailService.SendWelcomeEmail(ctx, e.Email, e.FirstName+" "+e.LastName)
		case *events.OrderCreatedEvent:
			// You would need to get user email from the event or fetch it
			logger.WithContext(ctx).Infof("Order created notification for order: %s", e.OrderNumber)
			return nil
		case *events.ProductStockUpdatedEvent:
			if e.NewStock <= 5 && e.OldStock > 5 {
				// Send low stock alert
				logger.WithContext(ctx).Warnf("Low stock alert for product: %s", e.ProductID)
				return nil
			}
		}
		return nil
	}
}

// SetupDefaultHandlers sets up default event handlers
func (p *InMemoryEventPublisher) SetupDefaultHandlers() {
	// Register logging handler for all events
	loggingHandler := LoggingEventHandler(p.logger)
	
	// Register for all event types
	eventTypes := []string{
		"UserRegistered",
		"UserProfileUpdated",
		"ProductCreated",
		"ProductStockUpdated",
		"OrderCreated",
		"OrderStatusChanged",
		"OrderCancelled",
		"PaymentProcessed",
		"CartItemAdded",
		"CartCleared",
	}
	
	for _, eventType := range eventTypes {
		p.Subscribe(eventType, loggingHandler)
	}
	
	p.logger.Info("Default event handlers registered")
}
