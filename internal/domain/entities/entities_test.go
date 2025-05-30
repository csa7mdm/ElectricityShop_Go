package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestUser_FullName(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
	}
	
	expected := "John Doe"
	actual := user.FullName()
	
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Admin user", RoleAdmin, true},
		{"Customer user", RoleCustomer, false},
		{"Employee user", RoleEmployee, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			if got := user.IsAdmin(); got != tt.expected {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_IsInStock(t *testing.T) {
	tests := []struct {
		name     string
		stock    int
		expected bool
	}{
		{"In stock", 10, true},
		{"Out of stock", 0, false},
		{"Low stock", 1, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{Stock: tt.stock}
			if got := product.IsInStock(); got != tt.expected {
				t.Errorf("IsInStock() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProduct_CanOrder(t *testing.T) {
	product := &Product{
		Stock:    10,
		IsActive: true,
	}
	
	tests := []struct {
		name     string
		quantity int
		expected bool
	}{
		{"Can order available quantity", 5, true},
		{"Can order exact stock", 10, true},
		{"Cannot order more than stock", 15, false},
		{"Cannot order zero", 0, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := product.CanOrder(tt.quantity); got != tt.expected {
				t.Errorf("CanOrder(%d) = %v, want %v", tt.quantity, got, tt.expected)
			}
		})
	}
}

func TestCart_GetTotal(t *testing.T) {
	cart := &Cart{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Items: []CartItem{
			{
				Quantity:  2,
				UnitPrice: decimal.NewFromFloat(10.50),
				Total:     decimal.NewFromFloat(21.00),
			},
			{
				Quantity:  1,
				UnitPrice: decimal.NewFromFloat(15.00),
				Total:     decimal.NewFromFloat(15.00),
			},
		},
	}
	
	expected := decimal.NewFromFloat(36.00)
	actual := cart.GetTotal()
	
	if !actual.Equal(expected) {
		t.Errorf("Expected %s, got %s", expected.String(), actual.String())
	}
}

func TestCart_GetItemCount(t *testing.T) {
	cart := &Cart{
		Items: []CartItem{
			{Quantity: 2},
			{Quantity: 3},
			{Quantity: 1},
		},
	}
	
	expected := 6
	actual := cart.GetItemCount()
	
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestAddress_FullAddress(t *testing.T) {
	address := &Address{
		AddressLine1: "123 Main St",
		AddressLine2: "Apt 4B",
		City:         "New York",
		State:        "NY",
		ZipCode:      "10001",
		Country:      "USA",
	}
	
	expected := "123 Main St, Apt 4B, New York, NY 10001"
	actual := address.FullAddress()
	
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAddress_FullAddressWithoutLine2(t *testing.T) {
	address := &Address{
		AddressLine1: "123 Main St",
		City:         "New York",
		State:        "NY",
		ZipCode:      "10001",
		Country:      "USA",
	}
	
	expected := "123 Main St, New York, NY 10001"
	actual := address.FullAddress()
	
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestCategory_IsRoot(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uuid.UUID
		expected bool
	}{
		{"Root category", nil, true},
		{"Child category", &uuid.UUID{}, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := &Category{ParentID: tt.parentID}
			if got := category.IsRoot(); got != tt.expected {
				t.Errorf("IsRoot() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOrder_CanBeCancelled(t *testing.T) {
	tests := []struct {
		name     string
		status   OrderStatus
		expected bool
	}{
		{"Pending order", OrderStatusPending, true},
		{"Confirmed order", OrderStatusConfirmed, true},
		{"Processing order", OrderStatusProcessing, false},
		{"Shipped order", OrderStatusShipped, false},
		{"Delivered order", OrderStatusDelivered, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{Status: tt.status}
			if got := order.CanBeCancelled(); got != tt.expected {
				t.Errorf("CanBeCancelled() = %v, want %v", got, tt.expected)
			}
		})
	}
}
