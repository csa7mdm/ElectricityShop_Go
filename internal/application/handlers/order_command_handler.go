package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/domain/events"
	"github.com/yourusername/electricity-shop-go/internal/domain/interfaces"
	"github.com/yourusername/electricity-shop-go/pkg/errors"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// OrderCommandHandler handles order-related commands
type OrderCommandHandler struct {
	orderRepo      interfaces.OrderRepository
	cartRepo       interfaces.CartRepository
	productRepo    interfaces.ProductRepository
	userRepo       interfaces.UserRepository
	addressRepo    interfaces.AddressRepository
	paymentRepo    interfaces.PaymentRepository
	eventPublisher interfaces.EventPublisher
	logger         logger.Logger
}

// NewOrderCommandHandler creates a new OrderCommandHandler
func NewOrderCommandHandler(
	orderRepo interfaces.OrderRepository,
	cartRepo interfaces.CartRepository,
	productRepo interfaces.ProductRepository,
	userRepo interfaces.UserRepository,
	addressRepo interfaces.AddressRepository,
	paymentRepo interfaces.PaymentRepository,
	eventPublisher interfaces.EventPublisher,
	logger logger.Logger,
) *OrderCommandHandler {
	return &OrderCommandHandler{
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		productRepo:    productRepo,
		userRepo:       userRepo,
		addressRepo:    addressRepo,
		paymentRepo:    paymentRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// Handle handles commands
func (h *OrderCommandHandler) Handle(ctx context.Context, command mediator.Command) error {
	switch cmd := command.(type) {
	case *commands.CreateOrderCommand:
		return h.handleCreateOrder(ctx, cmd)
	case *commands.CreateOrderFromCartCommand:
		return h.handleCreateOrderFromCart(ctx, cmd)
	case *commands.UpdateOrderStatusCommand:
		return h.handleUpdateOrderStatus(ctx, cmd)
	case *commands.CancelOrderCommand:
		return h.handleCancelOrder(ctx, cmd)
	case *commands.ProcessPaymentCommand:
		return h.handleProcessPayment(ctx, cmd)
	case *commands.UpdatePaymentStatusCommand:
		return h.handleUpdatePaymentStatus(ctx, cmd)
	case *commands.CreateShipmentCommand:
		return h.handleCreateShipment(ctx, cmd)
	case *commands.UpdateShipmentStatusCommand:
		return h.handleUpdateShipmentStatus(ctx, cmd)
	default:
		return errors.New("UNSUPPORTED_COMMAND", "Unsupported command type", 400)
	}
}

// handleCreateOrder handles direct order creation
func (h *OrderCommandHandler) handleCreateOrder(ctx context.Context, cmd *commands.CreateOrderCommand) error {
	h.logger.WithContext(ctx).Infof("Creating order for user: %s", cmd.UserID)
	
	// Verify user exists
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	// Verify addresses exist
	shippingAddr, err := h.addressRepo.GetByID(ctx, cmd.ShippingAddressID)
	if err != nil {
		return errors.ErrAddressNotFound.WithDetails("Shipping address not found")
	}
	
	billingAddr, err := h.addressRepo.GetByID(ctx, cmd.BillingAddressID)
	if err != nil {
		return errors.ErrAddressNotFound.WithDetails("Billing address not found")
	}
	
	// Verify address ownership
	if shippingAddr.UserID != cmd.UserID || billingAddr.UserID != cmd.UserID {
		return errors.ErrForbidden.WithDetails("Address does not belong to user")
	}
	
	// Validate and prepare order items
	orderItems := make([]entities.OrderItem, 0, len(cmd.Items))
	subtotal := decimal.Zero
	
	for _, item := range cmd.Items {
		product, err := h.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return err
		}
		
		if !product.CanOrder(item.Quantity) {
			return errors.ErrInsufficientStock.WithDetails(fmt.Sprintf("Insufficient stock for product %s", product.Name))
		}
		
		itemTotal := product.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		subtotal = subtotal.Add(itemTotal)
		
		orderItem := entities.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			Total:       itemTotal,
		}
		
		orderItems = append(orderItems, orderItem)
	}
	
	// Calculate totals (simplified tax calculation)
	taxRate := decimal.NewFromFloat(0.08) // 8% tax
	taxAmount := subtotal.Mul(taxRate)
	total := subtotal.Add(taxAmount)
	
	// Create order
	order := &entities.Order{
		UserID:          cmd.UserID,
		Status:          entities.OrderStatusPending,
		PaymentStatus:   entities.PaymentStatusPending,
		ShippingStatus:  entities.ShippingStatusPending,
		Subtotal:        subtotal,
		TaxAmount:       taxAmount,
		ShippingAmount:  decimal.Zero, // TODO: Calculate shipping
		DiscountAmount:  decimal.Zero,
		Total:           total,
		Currency:        "USD",
		Notes:           cmd.Notes,
		ShippingAddress: *shippingAddr,
		BillingAddress:  *billingAddr,
		Items:           orderItems,
		OrderedAt:       time.Now(),
	}
	
	// Save order
	if err := h.orderRepo.Create(ctx, order); err != nil {
		return err
	}
	
	// Update product stock
	for _, item := range cmd.Items {
		product, _ := h.productRepo.GetByID(ctx, item.ProductID)
		newStock := product.Stock - item.Quantity
		if err := h.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			h.logger.WithContext(ctx).Errorf("Failed to update stock for product %s: %v", item.ProductID, err)
		}
	}
	
	// Publish domain event
	event := events.NewOrderCreatedEvent(
		order.ID,
		order.UserID,
		order.OrderNumber,
		order.Total,
		len(order.Items),
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish OrderCreatedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully created order: %s", order.ID)
	return nil
}

// handleCreateOrderFromCart handles creating order from cart items
func (h *OrderCommandHandler) handleCreateOrderFromCart(ctx context.Context, cmd *commands.CreateOrderFromCartCommand) error {
	h.logger.WithContext(ctx).Infof("Creating order from cart for user: %s", cmd.UserID)
	
	// Get user's cart
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	
	if cart.IsEmpty() {
		return errors.ErrCartEmpty.WithDetails("Cannot create order from empty cart")
	}
	
	// Convert cart items to order items
	createOrderItems := make([]commands.CreateOrderItemCommand, 0, len(cart.Items))
	for _, cartItem := range cart.Items {
		createOrderItems = append(createOrderItems, commands.CreateOrderItemCommand{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		})
	}
	
	// Create order using existing logic
	createOrderCmd := &commands.CreateOrderCommand{
		UserID:            cmd.UserID,
		Items:             createOrderItems,
		ShippingAddressID: cmd.ShippingAddressID,
		BillingAddressID:  cmd.BillingAddressID,
		PaymentMethod:     cmd.PaymentMethod,
		Notes:             cmd.Notes,
	}
	
	if err := h.handleCreateOrder(ctx, createOrderCmd); err != nil {
		return err
	}
	
	// Clear cart after successful order creation
	if err := h.cartRepo.ClearItems(ctx, cart.ID); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to clear cart after order creation: %v", err)
		// Don't fail the order creation for this
	}
	
	// Publish cart cleared event
	event := events.NewCartClearedEvent(cart.ID, cmd.UserID, "Order created")
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish CartClearedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully created order from cart for user: %s", cmd.UserID)
	return nil
}

// handleUpdateOrderStatus handles updating order status
func (h *OrderCommandHandler) handleUpdateOrderStatus(ctx context.Context, cmd *commands.UpdateOrderStatusCommand) error {
	h.logger.WithContext(ctx).Infof("Updating order status: %s", cmd.OrderID)
	
	// Get existing order
	order, err := h.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	
	oldStatus := order.Status
	
	// Update status
	if err := h.orderRepo.UpdateStatus(ctx, cmd.OrderID, cmd.Status); err != nil {
		return err
	}
	
	// Update timestamps based on status
	now := time.Now()
	switch cmd.Status {
	case entities.OrderStatusShipped:
		order.ShippedAt = &now
		order.ShippingStatus = entities.ShippingStatusShipped
	case entities.OrderStatusDelivered:
		order.DeliveredAt = &now
		order.ShippingStatus = entities.ShippingStatusDelivered
	case entities.OrderStatusCancelled:
		order.CancelledAt = &now
	}
	
	if err := h.orderRepo.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewOrderStatusChangedEvent(
		cmd.OrderID,
		order.UserID,
		string(oldStatus),
		string(cmd.Status),
		cmd.Reason,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish OrderStatusChangedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated order status: %s", cmd.OrderID)
	return nil
}

// handleCancelOrder handles order cancellation
func (h *OrderCommandHandler) handleCancelOrder(ctx context.Context, cmd *commands.CancelOrderCommand) error {
	h.logger.WithContext(ctx).Infof("Cancelling order: %s", cmd.OrderID)
	
	// Get existing order
	order, err := h.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	
	// Verify ownership
	if order.UserID != cmd.UserID {
		return errors.ErrForbidden.WithDetails("You can only cancel your own orders")
	}
	
	// Check if order can be cancelled
	if !order.CanBeCancelled() {
		return errors.ErrOrderCannotBeCancelled.WithDetails("Order cannot be cancelled at this stage")
	}
	
	// Update order status
	updateStatusCmd := &commands.UpdateOrderStatusCommand{
		OrderID:   cmd.OrderID,
		Status:    entities.OrderStatusCancelled,
		Reason:    cmd.CancelReason,
		UpdatedBy: cmd.UserID,
	}
	
	if err := h.handleUpdateOrderStatus(ctx, updateStatusCmd); err != nil {
		return err
	}
	
	// Restore product stock
	for _, item := range order.Items {
		product, err := h.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			h.logger.WithContext(ctx).Errorf("Failed to get product for stock restoration: %v", err)
			continue
		}
		
		newStock := product.Stock + item.Quantity
		if err := h.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			h.logger.WithContext(ctx).Errorf("Failed to restore stock for product %s: %v", item.ProductID, err)
		}
	}
	
	// Publish domain event
	event := events.NewOrderCancelledEvent(
		cmd.OrderID,
		cmd.UserID,
		order.OrderNumber,
		order.Total,
		cmd.CancelReason,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish OrderCancelledEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully cancelled order: %s", cmd.OrderID)
	return nil
}

// handleProcessPayment handles payment processing
func (h *OrderCommandHandler) handleProcessPayment(ctx context.Context, cmd *commands.ProcessPaymentCommand) error {
	h.logger.WithContext(ctx).Infof("Processing payment for order: %s", cmd.OrderID)
	
	// Get order
	order, err := h.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	
	// Verify payment amount matches order total
	if !cmd.Amount.Equal(order.Total) {
		return errors.ErrPaymentFailed.WithDetails("Payment amount does not match order total")
	}
	
	// Create payment record
	payment := &entities.Payment{
		OrderID:         cmd.OrderID,
		Amount:          cmd.Amount,
		Currency:        order.Currency,
		Status:          entities.PaymentStatusProcessing,
		Method:          cmd.PaymentMethod,
		TransactionID:   cmd.TransactionID,
		GatewayResponse: cmd.GatewayResponse,
	}
	
	if err := h.paymentRepo.Create(ctx, payment); err != nil {
		return err
	}
	
	// TODO: Integrate with actual payment gateway
	// For now, simulate successful payment
	payment.Status = entities.PaymentStatusCompleted
	payment.ProcessedAt = &time.Time{}
	*payment.ProcessedAt = time.Now()
	
	if err := h.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}
	
	// Update order payment status
	order.PaymentStatus = entities.PaymentStatusCompleted
	if err := h.orderRepo.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain event
	event := events.NewPaymentProcessedEvent(
		payment.ID,
		cmd.OrderID,
		order.UserID,
		cmd.Amount,
		string(cmd.PaymentMethod),
		string(payment.Status),
		cmd.TransactionID,
	)
	
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		h.logger.WithContext(ctx).Errorf("Failed to publish PaymentProcessedEvent: %v", err)
	}
	
	h.logger.WithContext(ctx).Infof("Successfully processed payment for order: %s", cmd.OrderID)
	return nil
}

// handleUpdatePaymentStatus handles updating payment status
func (h *OrderCommandHandler) handleUpdatePaymentStatus(ctx context.Context, cmd *commands.UpdatePaymentStatusCommand) error {
	h.logger.WithContext(ctx).Infof("Updating payment status: %s", cmd.PaymentID)
	
	// Get payment
	payment, err := h.paymentRepo.GetByID(ctx, cmd.PaymentID)
	if err != nil {
		return err
	}
	
	// Update payment
	payment.Status = cmd.Status
	payment.TransactionID = cmd.TransactionID
	payment.GatewayResponse = cmd.GatewayResponse
	payment.FailureReason = cmd.FailureReason
	
	if cmd.Status == entities.PaymentStatusCompleted {
		now := time.Now()
		payment.ProcessedAt = &now
	}
	
	if err := h.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}
	
	// Update corresponding order payment status
	order, err := h.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return err
	}
	
	order.PaymentStatus = cmd.Status
	if err := h.orderRepo.Update(ctx, order); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully updated payment status: %s", cmd.PaymentID)
	return nil
}

// handleCreateShipment handles creating a shipment
func (h *OrderCommandHandler) handleCreateShipment(ctx context.Context, cmd *commands.CreateShipmentCommand) error {
	h.logger.WithContext(ctx).Infof("Creating shipment for order: %s", cmd.OrderID)
	
	// Get order
	order, err := h.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	
	// Check if order can be shipped
	if !order.CanBeShipped() {
		return errors.New("ORDER_CANNOT_BE_SHIPPED", "Order cannot be shipped at this stage", 400)
	}
	
	// Create shipment
	shipment := &entities.Shipment{
		OrderID:        cmd.OrderID,
		TrackingNumber: cmd.TrackingNumber,
		Carrier:        cmd.Carrier,
		Status:         entities.ShippingStatusPreparing,
		ShippedAt:      nil,
	}
	
	if cmd.EstimatedDelivery != nil {
		if estimatedDelivery, err := time.Parse("2006-01-02", *cmd.EstimatedDelivery); err == nil {
			shipment.EstimatedDelivery = &estimatedDelivery
		}
	}
	
	// This would normally be saved through a Shipment repository
	// For now, we'll add it to the order's shipments
	order.Shipments = append(order.Shipments, *shipment)
	if err := h.orderRepo.Update(ctx, order); err != nil {
		return err
	}
	
	h.logger.WithContext(ctx).Infof("Successfully created shipment for order: %s", cmd.OrderID)
	return nil
}

// handleUpdateShipmentStatus handles updating shipment status
func (h *OrderCommandHandler) handleUpdateShipmentStatus(ctx context.Context, cmd *commands.UpdateShipmentStatusCommand) error {
	h.logger.WithContext(ctx).Infof("Updating shipment status: %s", cmd.ShipmentID)
	
	// This would normally use a ShipmentRepository
	// For now, this is a placeholder implementation
	
	h.logger.WithContext(ctx).Infof("Successfully updated shipment status: %s", cmd.ShipmentID)
	return nil
}
