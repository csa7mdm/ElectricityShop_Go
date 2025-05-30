package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter UserFilter) ([]*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProductFilter) ([]*entities.Product, error)
	Search(ctx context.Context, query string, filter ProductFilter) ([]*entities.Product, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID, filter ProductFilter) ([]*entities.Product, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error
	GetLowStockProducts(ctx context.Context, threshold int) ([]*entities.Product, error)
	ExistsBySKU(ctx context.Context, sku string) (bool, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CategoryFilter) ([]*entities.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entities.Category, error)
	GetRootCategories(ctx context.Context) ([]*entities.Category, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	Create(ctx context.Context, cart *entities.Cart) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Cart, error)
	GetBySessionID(ctx context.Context, sessionID string) (*entities.Cart, error)
	Update(ctx context.Context, cart *entities.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddItem(ctx context.Context, cartItem *entities.CartItem) error
	UpdateItem(ctx context.Context, cartItem *entities.CartItem) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	ClearItems(ctx context.Context, cartID uuid.UUID) error
	GetItems(ctx context.Context, cartID uuid.UUID) ([]*entities.CartItem, error)
	GetItemByProductID(ctx context.Context, cartID, productID uuid.UUID) (*entities.CartItem, error)
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID, filter OrderFilter) ([]*entities.Order, error)
	List(ctx context.Context, filter OrderFilter) ([]*entities.Order, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error
	GetOrdersToProcess(ctx context.Context) ([]*entities.Order, error)
	GetOrdersByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Order, error)
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(ctx context.Context, payment *entities.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error)
	Update(ctx context.Context, payment *entities.Payment) error
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entities.PaymentStatus) error
	List(ctx context.Context, filter PaymentFilter) ([]*entities.Payment, error)
}

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	Create(ctx context.Context, address *entities.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Address, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error)
	Update(ctx context.Context, address *entities.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetAsDefault(ctx context.Context, addressID uuid.UUID, addressType entities.AddressType) error
	UnsetDefaultForUser(ctx context.Context, userID uuid.UUID) error
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
	PublishBatch(ctx context.Context, events []interface{}) error
}

// CacheService defines the interface for caching
type CacheService interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl int) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// EmailService defines the interface for email notifications
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
	SendOrderConfirmation(ctx context.Context, email string, order *entities.Order) error
	SendOrderStatusUpdate(ctx context.Context, email string, order *entities.Order) error
	SendPasswordReset(ctx context.Context, email, resetToken string) error
	SendLowStockAlert(ctx context.Context, products []*entities.Product) error
}

// Filter structs for various queries
type UserFilter struct {
	Page     int
	PageSize int
	Role     entities.UserRole
	IsActive *bool
	Search   string
	SortBy   string
	SortDesc bool
}

type ProductFilter struct {
	Page       int
	PageSize   int
	CategoryID *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
	InStock    *bool
	IsActive   *bool
	IsFeatured *bool
	Brand      string
	Search     string
	SortBy     string
	SortDesc   bool
}

type CategoryFilter struct {
	Page     int
	PageSize int
	ParentID *uuid.UUID
	IsActive *bool
	SortBy   string
	SortDesc bool
}

type OrderFilter struct {
	Page          int
	PageSize      int
	UserID        *uuid.UUID
	Status        entities.OrderStatus
	PaymentStatus entities.PaymentStatus
	StartDate     *string
	EndDate       *string
	MinTotal      *float64
	MaxTotal      *float64
	SortBy        string
	SortDesc      bool
}

type PaymentFilter struct {
	Page      int
	PageSize  int
	OrderID   *uuid.UUID
	UserID    *uuid.UUID
	Status    entities.PaymentStatus
	Method    entities.PaymentMethod
	StartDate *string
	EndDate   *string
	SortBy    string
	SortDesc  bool
}

type ReviewFilter struct {
	Page       int
	PageSize   int
	Rating     *int
	IsApproved *bool
	IsVerified *bool
	SortBy     string
	SortDesc   bool
}

// UnitOfWork defines the interface for unit of work pattern
type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	UserRepository() UserRepository
	ProductRepository() ProductRepository
	CategoryRepository() CategoryRepository
	CartRepository() CartRepository
	OrderRepository() OrderRepository
	PaymentRepository() PaymentRepository
	AddressRepository() AddressRepository
}
