package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	SessionID string    `gorm:"type:varchar(255)" json:"session_id"` // for guest users
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	User  User       `gorm:"foreignKey:UserID" json:"user"`
	Items []CartItem `gorm:"foreignKey:CartID" json:"items"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CartID    uuid.UUID       `gorm:"type:uuid;not null" json:"cart_id"`
	ProductID uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int             `gorm:"not null;check:quantity > 0" json:"quantity"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Total     decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	
	// Relationships
	Cart    Cart    `gorm:"foreignKey:CartID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Order represents a customer order
type Order struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	OrderNumber     string          `gorm:"unique;not null;type:varchar(50)" json:"order_number"`
	Status          OrderStatus     `gorm:"not null;type:varchar(50);default:'pending'" json:"status"`
	PaymentStatus   PaymentStatus   `gorm:"not null;type:varchar(50);default:'pending'" json:"payment_status"`
	ShippingStatus  ShippingStatus  `gorm:"not null;type:varchar(50);default:'pending'" json:"shipping_status"`
	Subtotal        decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	ShippingAmount  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"shipping_amount"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	Total           decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total"`
	Currency        string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Notes           string          `gorm:"type:text" json:"notes"`
	ShippingAddress EmbeddableAddress `gorm:"embedded;embeddedPrefix:shipping_" json:"shipping_address"`
	BillingAddress  EmbeddableAddress `gorm:"embedded;embeddedPrefix:billing_" json:"billing_address"`
	OrderedAt       time.Time       `json:"ordered_at"`
	ShippedAt       *time.Time      `json:"shipped_at"`
	DeliveredAt     *time.Time      `json:"delivered_at"`
	CancelledAt     *time.Time      `json:"cancelled_at"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
	
	// Relationships
	User     User        `gorm:"foreignKey:UserID" json:"user"`
	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payments []Payment   `gorm:"foreignKey:OrderID" json:"payments"`
	Shipments []Shipment `gorm:"foreignKey:OrderID" json:"shipments"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	ProductID   uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	ProductName string          `gorm:"not null;type:varchar(255)" json:"product_name"` // snapshot
	ProductSKU  string          `gorm:"not null;type:varchar(100)" json:"product_sku"`  // snapshot
	Quantity    int             `gorm:"not null;check:quantity > 0" json:"quantity"`
	UnitPrice   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Total       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	
	// Relationships
	Order   Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Payment represents a payment transaction
type Payment struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID     `gorm:"type:uuid;not null" json:"order_id"`
	Amount          decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        string        `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Status          PaymentStatus `gorm:"not null;type:varchar(50)" json:"status"`
	Method          PaymentMethod `gorm:"not null;type:varchar(50)" json:"method"`
	TransactionID   string        `gorm:"type:varchar(255)" json:"transaction_id"`
	GatewayResponse string        `gorm:"type:text" json:"gateway_response"` // JSON response
	ProcessedAt     *time.Time    `json:"processed_at"`
	FailureReason   string        `gorm:"type:varchar(500)" json:"failure_reason"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	
	// Relationships
	Order Order `gorm:"foreignKey:OrderID" json:"-"`
}

// Shipment represents a shipment
type Shipment struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID       uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	TrackingNumber string        `gorm:"type:varchar(255)" json:"tracking_number"`
	Carrier       string        `gorm:"type:varchar(100)" json:"carrier"`
	Status        ShippingStatus `gorm:"not null;type:varchar(50)" json:"status"`
	ShippedAt     *time.Time     `json:"shipped_at"`
	DeliveredAt   *time.Time     `json:"delivered_at"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	
	// Relationships
	Order Order `gorm:"foreignKey:OrderID" json:"-"`
}

// Enums
type OrderStatus string
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string
const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string
const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPayPal     PaymentMethod = "paypal"
	PaymentMethodStripe     PaymentMethod = "stripe"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCash       PaymentMethod = "cash"
)

type ShippingStatus string
const (
	ShippingStatusPending   ShippingStatus = "pending"
	ShippingStatusPreparing ShippingStatus = "preparing"
	ShippingStatusShipped   ShippingStatus = "shipped"
	ShippingStatusInTransit ShippingStatus = "in_transit"
	ShippingStatusDelivered ShippingStatus = "delivered"
	ShippingStatusReturned  ShippingStatus = "returned"
)

// BeforeCreate hooks
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	// Calculate total
	ci.Total = ci.UnitPrice.Mul(decimal.NewFromInt(int64(ci.Quantity)))
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}
	if o.OrderedAt.IsZero() {
		o.OrderedAt = time.Now()
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	// Calculate total
	oi.Total = oi.UnitPrice.Mul(decimal.NewFromInt(int64(oi.Quantity)))
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (s *Shipment) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// Business logic methods
func (c *Cart) GetTotal() decimal.Decimal {
	total := decimal.Zero
	for _, item := range c.Items {
		total = total.Add(item.Total)
	}
	return total
}

func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

func (o *Order) CanBeShipped() bool {
	return o.Status == OrderStatusProcessing && o.PaymentStatus == PaymentStatusCompleted
}

func (o *Order) IsPaid() bool {
	return o.PaymentStatus == PaymentStatusCompleted
}

func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

// Helper function to generate order number
func generateOrderNumber() string {
	return "ORD-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:8]
}
