package responses

// SuccessResponse wraps a successful API response.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse wraps an error API response.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// Pagination represents pagination details for a list response.
type Pagination struct {
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

// PaginatedResponse wraps a successful API response that includes pagination.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// NewSuccessResponse creates a new success response.
func NewSuccessResponse(data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse creates a new error response.
func NewErrorResponse(err string, code string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
		Code:    code,
	}
}

// NewPaginatedResponse creates a new paginated response.
func NewPaginatedResponse(data interface{}, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
}
