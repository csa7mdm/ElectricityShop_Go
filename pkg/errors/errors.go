package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with additional context
type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	HTTPStatus int               `json:"-"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// New creates a new AppError
func New(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    err.Error(),
		HTTPStatus: httpStatus,
	}
}

// Pre-defined error types for common scenarios

// Authentication and Authorization Errors
var (
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "Authentication token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidToken = &AppError{
		Code:       "INVALID_TOKEN",
		Message:    "Invalid authentication token",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Access denied",
		HTTPStatus: http.StatusForbidden,
	}

	ErrInsufficientPermissions = &AppError{
		Code:       "INSUFFICIENT_PERMISSIONS",
		Message:    "You don't have permission to perform this action",
		HTTPStatus: http.StatusForbidden,
	}
)

// Validation Errors
var (
	ErrValidationFailed = &AppError{
		Code:       "VALIDATION_FAILED",
		Message:    "Validation failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidInput = &AppError{
		Code:       "INVALID_INPUT",
		Message:    "Invalid input provided",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrRequiredField = &AppError{
		Code:       "REQUIRED_FIELD",
		Message:    "Required field is missing",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidFormat = &AppError{
		Code:       "INVALID_FORMAT",
		Message:    "Invalid format provided",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidEmail = &AppError{
		Code:       "INVALID_EMAIL",
		Message:    "Invalid email format",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrWeakPassword = &AppError{
		Code:       "WEAK_PASSWORD",
		Message:    "Password does not meet security requirements",
		HTTPStatus: http.StatusBadRequest,
	}
)

// Resource Not Found Errors
var (
	ErrUserNotFound = &AppError{
		Code:       "USER_NOT_FOUND",
		Message:    "User not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrProductNotFound = &AppError{
		Code:       "PRODUCT_NOT_FOUND",
		Message:    "Product not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrCategoryNotFound = &AppError{
		Code:       "CATEGORY_NOT_FOUND",
		Message:    "Category not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrOrderNotFound = &AppError{
		Code:       "ORDER_NOT_FOUND",
		Message:    "Order not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrCartNotFound = &AppError{
		Code:       "CART_NOT_FOUND",
		Message:    "Cart not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrPaymentNotFound = &AppError{
		Code:       "PAYMENT_NOT_FOUND",
		Message:    "Payment not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrAddressNotFound = &AppError{
		Code:       "ADDRESS_NOT_FOUND",
		Message:    "Address not found",
		HTTPStatus: http.StatusNotFound,
	}
)

// Business Logic Errors
var (
	ErrInsufficientStock = &AppError{
		Code:       "INSUFFICIENT_STOCK",
		Message:    "Insufficient stock available",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrProductOutOfStock = &AppError{
		Code:       "PRODUCT_OUT_OF_STOCK",
		Message:    "Product is out of stock",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrCartEmpty = &AppError{
		Code:       "CART_EMPTY",
		Message:    "Cart is empty",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrOrderAlreadyProcessed = &AppError{
		Code:       "ORDER_ALREADY_PROCESSED",
		Message:    "Order has already been processed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrOrderCannotBeCancelled = &AppError{
		Code:       "ORDER_CANNOT_BE_CANCELLED",
		Message:    "Order cannot be cancelled at this stage",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPaymentFailed = &AppError{
		Code:       "PAYMENT_FAILED",
		Message:    "Payment processing failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPaymentAlreadyProcessed = &AppError{
		Code:       "PAYMENT_ALREADY_PROCESSED",
		Message:    "Payment has already been processed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidQuantity = &AppError{
		Code:       "INVALID_QUANTITY",
		Message:    "Invalid quantity specified",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPriceChanged = &AppError{
		Code:       "PRICE_CHANGED",
		Message:    "Product price has changed since it was added to cart",
		HTTPStatus: http.StatusBadRequest,
	}
)

// Conflict Errors
var (
	ErrUserAlreadyExists = &AppError{
		Code:       "USER_ALREADY_EXISTS",
		Message:    "User with this email already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrProductAlreadyExists = &AppError{
		Code:       "PRODUCT_ALREADY_EXISTS",
		Message:    "Product with this SKU already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrCategoryAlreadyExists = &AppError{
		Code:       "CATEGORY_ALREADY_EXISTS",
		Message:    "Category with this slug already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrDuplicateOrderNumber = &AppError{
		Code:       "DUPLICATE_ORDER_NUMBER",
		Message:    "Order number already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrResourceInUse = &AppError{
		Code:       "RESOURCE_IN_USE",
		Message:    "Resource is currently in use and cannot be deleted",
		HTTPStatus: http.StatusConflict,
	}
)

// System Errors
var (
	ErrInternalServer = &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An internal server error occurred",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrDatabaseConnection = &AppError{
		Code:       "DATABASE_CONNECTION_ERROR",
		Message:    "Database connection failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrCacheUnavailable = &AppError{
		Code:       "CACHE_UNAVAILABLE",
		Message:    "Cache service is unavailable",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrExternalServiceUnavailable = &AppError{
		Code:       "EXTERNAL_SERVICE_UNAVAILABLE",
		Message:    "External service is unavailable",
		HTTPStatus: http.StatusServiceUnavailable,
	}

	ErrRateLimitExceeded = &AppError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    "Rate limit exceeded",
		HTTPStatus: http.StatusTooManyRequests,
	}

	ErrTimeout = &AppError{
		Code:       "TIMEOUT",
		Message:    "Request timeout",
		HTTPStatus: http.StatusRequestTimeout,
	}
)

// Helper functions

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError tries to convert an error to AppError
func GetAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if appErr, ok := GetAppError(err); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// NewValidationError creates a new validation error with field details
func NewValidationError(field, message string) *AppError {
	return ErrValidationFailed.
		WithDetails(fmt.Sprintf("Field '%s': %s", field, message)).
		WithMetadata("field", field)
}

// NewNotFoundError creates a new not found error for a specific resource
func NewNotFoundError(resource, identifier string) *AppError {
	return New(
		fmt.Sprintf("%s_NOT_FOUND", resource),
		fmt.Sprintf("%s not found", resource),
		http.StatusNotFound,
	).WithDetails(fmt.Sprintf("%s with identifier '%s' was not found", resource, identifier))
}

// NewConflictError creates a new conflict error for a specific resource
func NewConflictError(resource, field, value string) *AppError {
	return New(
		fmt.Sprintf("%s_ALREADY_EXISTS", resource),
		fmt.Sprintf("%s already exists", resource),
		http.StatusConflict,
	).WithDetails(fmt.Sprintf("%s with %s '%s' already exists", resource, field, value))
}

// NewBusinessLogicError creates a new business logic error
func NewBusinessLogicError(code, message string) *AppError {
	return New(code, message, http.StatusBadRequest)
}

// ErrorResponse represents the structure for error responses
type ErrorResponse struct {
	Success   bool              `json:"success"`
	Error     string            `json:"error"`
	Code      string            `json:"code"`
	Details   string            `json:"details,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp string            `json:"timestamp"`
}

// ToErrorResponse converts an AppError to an ErrorResponse
func ToErrorResponse(err *AppError) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Error:     err.Message,
		Code:      err.Code,
		Details:   err.Details,
		Metadata:  err.Metadata,
		Timestamp: fmt.Sprintf("%d", 1234567890), // You should use actual timestamp here
	}
}
