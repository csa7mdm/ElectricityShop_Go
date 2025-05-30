package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims represents the claims in our JWT token
type JWTClaims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   entities.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	secretKey []byte
	tokenTTL  time.Duration
}

// NewAuthService creates a new AuthService
func NewAuthService(secretKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
		tokenTTL:  tokenTTL,
	}
}

// HashPassword hashes a plain text password
func (s *AuthService) HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// VerifyPassword verifies a password against its hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userID uuid.UUID, email string, role entities.UserRole) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID.String(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "electricity-shop",
			Subject:   userID.String(),
			Audience:  []string{"electricity-shop-api"},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token from an existing valid token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Parse user ID back to UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", err
	}

	return s.GenerateToken(userID, claims.Email, claims.Role)
}
