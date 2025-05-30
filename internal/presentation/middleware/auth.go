package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(authService *auth.AuthService, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Authorization header required", "MISSING_AUTH_HEADER"))
			c.Abort()
			return
		}

		// Check for Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid authorization header format", "INVALID_AUTH_HEADER"))
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			logger.Errorf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid token", "INVALID_TOKEN"))
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RequireRole creates middleware that requires specific user roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("User role not found in context", "MISSING_USER_ROLE"))
			c.Abort()
			return
		}

		userRoleStr := string(userRole.(entities.UserRole))
		
		// Check if user has required role
		for _, requiredRole := range roles {
			if userRoleStr == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, responses.NewErrorResponse("Insufficient permissions", "INSUFFICIENT_PERMISSIONS"))
		c.Abort()
	}
}

// OptionalAuth middleware that extracts user info if token is present but doesn't require it
func OptionalAuth(authService *auth.AuthService, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			logger.Warnf("Optional auth token validation failed: %v", err)
			c.Next()
			return
		}

		// Set user context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}
