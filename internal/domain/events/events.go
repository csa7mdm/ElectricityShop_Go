package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DomainEvent represents a domain event interface
type DomainEvent interface {
	GetEventType() string
	GetAggregateID() uuid.UUID
	GetOccurredAt() time.Time
	GetEventData() interface{}
}

// BaseDomainEvent provides common event fields
type BaseDomainEvent struct {
	EventType   string    `json:"event_type"`
	AggregateID uuid.UUID `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (e BaseDomainEvent) GetEventType() string {
	return e.EventType
}

func (e BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

func (e BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// User Events
type UserRegisteredEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
}

func NewUserRegisteredEvent(userID uuid.UUID, email, firstName, lastName, role string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "UserRegistered",
			AggregateID: userID,
			OccurredAt:  time.Now(),
		},
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}
}

func (e UserRegisteredEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"user_id":    e.UserID,
		"email":      e.Email,
		"first_name": e.FirstName,
		"last_name":  e.LastName,
		"role":       e.Role,
	}
}

type UserProfileUpdatedEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
}

func NewUserProfileUpdatedEvent(userID uuid.UUID, email, firstName, lastName, phone string) *UserProfileUpdatedEvent {
	return &UserProfileUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "UserProfileUpdated",
			AggregateID: userID,
			OccurredAt:  time.Now(),
		},
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}
}

func (e UserProfileUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"user_id":    e.UserID,
		"email":      e.Email,
		"first_name": e.FirstName,
		"last_name":  e.LastName,
		"phone":      e.Phone,
	}
}

// Product Events
type ProductCreatedEvent struct {
	BaseDomainEvent
	ProductID   uuid.UUID       `json:"product_id"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Price       decimal.Decimal `json:"price"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Stock       int             `json:"stock"`
}

func NewProductCreatedEvent(productID uuid.UUID, name, sku string, price decimal.Decimal, categoryID uuid.UUID, stock int) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "ProductCreated",
			AggregateID: productID,
			OccurredAt:  time.Now(),
		},
		ProductID:  productID,
		Name:       name,
		SKU:        sku,
		Price:      price,
		CategoryID: categoryID,
		Stock:      stock,
	}
}

func (e ProductCreatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"product_id":  e.ProductID,
		"name":        e.Name,
		"sku":         e.SKU,
		"price":       e.Price,
		"category_id": e.CategoryID,
		"stock":       e.Stock,
	}
}

type ProductStockUpdatedEvent struct {
	BaseDomainEvent
	ProductID uuid.UUID `json:"product_id"`
	OldStock  int       `json:"old_stock"`
	NewStock  int       `json:"new_stock"`
	Reason    string    `json:"reason"`
}

func NewProductStockUpdatedEvent(productID uuid.UUID, oldStock, newStock int, reason string) *ProductStockUpdatedEvent {
	return &ProductStockUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "ProductStockUpdated",
			AggregateID: productID,
			OccurredAt:  time.Now(),
		},
		ProductID: productID,
		OldStock:  oldStock,
		NewStock:  newStock,
		Reason:    reason,
	}
}

func (e ProductStockUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"product_id": e.ProductID,
		"old_stock":  e.OldStock,
		"new_stock":  e.NewStock,
		"reason":     e.Reason,
	}
}

// Order Events
type OrderCreatedEvent struct {
	BaseDomainEvent
	OrderID     uuid.UUID       `json:"order_id"`
	UserID      uuid.UUID       `json:"user_id"`
	OrderNumber string          `json:"order_number"`
	Total       decimal.Decimal `json:"total"`
	ItemCount   int             `json:"item_count"`
}

func NewOrderCreatedEvent(orderID, userID uuid.UUID, orderNumber string, total decimal.Decimal, itemCount int) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "OrderCreated",
			AggregateID: orderID,
			OccurredAt:  time.Now(),
		},
		OrderID:     orderID,
		UserID:      userID,
		OrderNumber: orderNumber,
		Total:       total,
		ItemCount:   itemCount,
	}
}

func (e OrderCreatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"order_id":     e.OrderID,
		"user_id":      e.UserID,
		"order_number": e.OrderNumber,
		"total":        e.Total,
		"item_count":   e.ItemCount,
	}
}

type OrderStatusChangedEvent struct {
	BaseDomainEvent
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Reason    string    `json:"reason"`
}

func NewOrderStatusChangedEvent(orderID, userID uuid.UUID, oldStatus, newStatus, reason string) *OrderStatusChangedEvent {
	return &OrderStatusChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "OrderStatusChanged",
			AggregateID: orderID,
			OccurredAt:  time.Now(),
		},
		OrderID:   orderID,
		UserID:    userID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
	}
}

func (e OrderStatusChangedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"order_id":   e.OrderID,
		"user_id":    e.UserID,
		"old_status": e.OldStatus,
		"new_status": e.NewStatus,
		"reason":     e.Reason,
	}
}

type OrderCancelledEvent struct {
	BaseDomainEvent
	OrderID       uuid.UUID       `json:"order_id"`
	UserID        uuid.UUID       `json:"user_id"`
	OrderNumber   string          `json:"order_number"`
	RefundAmount  decimal.Decimal `json:"refund_amount"`
	CancelReason  string          `json:"cancel_reason"`
}

func NewOrderCancelledEvent(orderID, userID uuid.UUID, orderNumber string, refundAmount decimal.Decimal, cancelReason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "OrderCancelled",
			AggregateID: orderID,
			OccurredAt:  time.Now(),
		},
		OrderID:      orderID,
		UserID:       userID,
		OrderNumber:  orderNumber,
		RefundAmount: refundAmount,
		CancelReason: cancelReason,
	}
}

func (e OrderCancelledEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"order_id":      e.OrderID,
		"user_id":       e.UserID,
		"order_number":  e.OrderNumber,
		"refund_amount": e.RefundAmount,
		"cancel_reason": e.CancelReason,
	}
}

// Payment Events
type PaymentProcessedEvent struct {
	BaseDomainEvent
	PaymentID     uuid.UUID       `json:"payment_id"`
	OrderID       uuid.UUID       `json:"order_id"`
	UserID        uuid.UUID       `json:"user_id"`
	Amount        decimal.Decimal `json:"amount"`
	Method        string          `json:"method"`
	Status        string          `json:"status"`
	TransactionID string          `json:"transaction_id"`
}

func NewPaymentProcessedEvent(paymentID, orderID, userID uuid.UUID, amount decimal.Decimal, method, status, transactionID string) *PaymentProcessedEvent {
	return &PaymentProcessedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "PaymentProcessed",
			AggregateID: paymentID,
			OccurredAt:  time.Now(),
		},
		PaymentID:     paymentID,
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Method:        method,
		Status:        status,
		TransactionID: transactionID,
	}
}

func (e PaymentProcessedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"payment_id":     e.PaymentID,
		"order_id":       e.OrderID,
		"user_id":        e.UserID,
		"amount":         e.Amount,
		"method":         e.Method,
		"status":         e.Status,
		"transaction_id": e.TransactionID,
	}
}

// Cart Events
type CartItemAddedEvent struct {
	BaseDomainEvent
	CartID    uuid.UUID       `json:"cart_id"`
	UserID    uuid.UUID       `json:"user_id"`
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

func NewCartItemAddedEvent(cartID, userID, productID uuid.UUID, quantity int, unitPrice decimal.Decimal) *CartItemAddedEvent {
	return &CartItemAddedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "CartItemAdded",
			AggregateID: cartID,
			OccurredAt:  time.Now(),
		},
		CartID:    cartID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func (e CartItemAddedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"cart_id":    e.CartID,
		"user_id":    e.UserID,
		"product_id": e.ProductID,
		"quantity":   e.Quantity,
		"unit_price": e.UnitPrice,
	}
}

type CartClearedEvent struct {
	BaseDomainEvent
	CartID uuid.UUID `json:"cart_id"`
	UserID uuid.UUID `json:"user_id"`
	Reason string    `json:"reason"`
}

func NewCartClearedEvent(cartID, userID uuid.UUID, reason string) *CartClearedEvent {
	return &CartClearedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "CartCleared",
			AggregateID: cartID,
			OccurredAt:  time.Now(),
		},
		CartID: cartID,
		UserID: userID,
		Reason: reason,
	}
}

func (e CartClearedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"cart_id": e.CartID,
		"user_id": e.UserID,
		"reason":  e.Reason,
	}
}
