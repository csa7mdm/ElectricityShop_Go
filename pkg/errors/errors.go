package errors

// AppError represents a custom application error.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error returns the message of the AppError.
func (e *AppError) Error() string {
	return e.Message
}

// Standard application errors.
var (
	ErrUserNotFound     = &AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrProductNotFound  = &AppError{Code: "PRODUCT_NOT_FOUND", Message: "Product not found"}
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	// Add other generic errors here as needed
)
