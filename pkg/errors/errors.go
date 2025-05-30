package errors

import (
	"fmt"
	"strings"
)

// AppError represents a custom application error.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"`
}

// Error returns the message of the AppError.
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
		Status:  e.Status,
	}
}

// Standard application errors.
var (
	// User errors
	ErrUserNotFound       = &AppError{Code: "USER_NOT_FOUND", Message: "User not found", Status: 404}
	ErrUserAlreadyExists  = &AppError{Code: "USER_ALREADY_EXISTS", Message: "User already exists", Status: 409}
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials", Status: 401}
	ErrUserInactive       = &AppError{Code: "USER_INACTIVE", Message: "User account is inactive", Status: 403}
	ErrUnauthorized       = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized access", Status: 403}
	
	// Product errors
	ErrProductNotFound      = &AppError{Code: "PRODUCT_NOT_FOUND", Message: "Product not found", Status: 404}
	ErrProductAlreadyExists = &AppError{Code: "PRODUCT_ALREADY_EXISTS", Message: "Product already exists", Status: 409}
	ErrInsufficientStock    = &AppError{Code: "INSUFFICIENT_STOCK", Message: "Insufficient stock", Status: 400}
	
	// Category errors
	ErrCategoryNotFound      = &AppError{Code: "CATEGORY_NOT_FOUND", Message: "Category not found", Status: 404}
	ErrCategoryAlreadyExists = &AppError{Code: "CATEGORY_ALREADY_EXISTS", Message: "Category already exists", Status: 409}
	
	// Address errors
	ErrAddressNotFound = &AppError{Code: "ADDRESS_NOT_FOUND", Message: "Address not found", Status: 404}
	
	// Cart errors
	ErrCartNotFound = &AppError{Code: "CART_NOT_FOUND", Message: "Cart not found", Status: 404}
	ErrCartEmpty    = &AppError{Code: "CART_EMPTY", Message: "Cart is empty", Status: 400}
	
	// Order errors
	ErrOrderNotFound = &AppError{Code: "ORDER_NOT_FOUND", Message: "Order not found", Status: 404}
	ErrOrderCannotBeCancelled = &AppError{Code: "ORDER_CANNOT_BE_CANCELLED", Message: "Order cannot be cancelled", Status: 400}
	
	// Payment errors
	ErrPaymentNotFound = &AppError{Code: "PAYMENT_NOT_FOUND", Message: "Payment not found", Status: 404}
	ErrPaymentFailed   = &AppError{Code: "PAYMENT_FAILED", Message: "Payment processing failed", Status: 400}
	ErrDuplicateOrderNumber = &AppError{Code: "DUPLICATE_ORDER_NUMBER", Message: "Duplicate order number", Status: 409}
	
	// Access errors
	ErrForbidden = &AppError{Code: "FORBIDDEN", Message: "Access forbidden", Status: 403}
	
	// Generic errors
	ErrResourceInUse     = &AppError{Code: "RESOURCE_IN_USE", Message: "Resource is in use", Status: 400}
	ErrValidationFailed  = &AppError{Code: "VALIDATION_FAILED", Message: "Validation failed", Status: 400}
	ErrInternalError     = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", Status: 500}
)

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap wraps an existing error with AppError context
func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: err.Error(),
		Status:  status,
	}
}

// IsErrorType checks if an error is of a specific type
func IsErrorType(err error, errorType string) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == errorType
	}
	return false
}

// NewBusinessLogicError creates a business logic error
func NewBusinessLogicError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  400,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// Database error utility functions

// IsUniqueConstraintError checks if the error is a unique constraint violation
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// For PostgreSQL
	errorStr := err.Error()
	return strings.Contains(errorStr, "duplicate key") || 
		   strings.Contains(errorStr, "UNIQUE constraint")
}
