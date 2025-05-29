package dtos

import "time" // Ensure time is imported

// RegisterUserRequest is the DTO for user registration requests.
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	// ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"` // Optional: if confirm password is needed
}

// LoginUserRequest is the DTO for user login requests.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginUserResponse is the DTO for user login responses.
type LoginUserResponse struct {
	ID    string `json:"id"`         // Using string for UUID representation
	Email string `json:"email"`
	Role  string `json:"role"`        // Using string for UserRole representation
	Token string `json:"token"`       // Placeholder for JWT token
}

// UserResponse is the DTO for returning user details.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Add Profile information here later if UserProfile is implemented
}
