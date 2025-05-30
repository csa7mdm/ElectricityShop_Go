package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("🔐 Testing JWT Authentication System...")
	
	// Initialize auth service
	authService := auth.NewAuthService("test-secret-key", 24*time.Hour)
	
	// Test password hashing
	fmt.Println("\n1. Testing Password Hashing...")
	password := "testPassword123"
	hashedPassword, err := authService.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Printf("✅ Password hashed successfully")
	
	// Test password verification
	fmt.Println("\n\n2. Testing Password Verification...")
	err = authService.VerifyPassword(hashedPassword, password)
	if err != nil {
		log.Fatalf("Password verification failed: %v", err)
	}
	fmt.Printf("✅ Password verification successful")
	
	// Test wrong password
	err = authService.VerifyPassword(hashedPassword, "wrongPassword")
	if err == nil {
		log.Fatalf("Password verification should have failed")
	}
	fmt.Printf("✅ Wrong password correctly rejected")
	
	// Test JWT token generation
	fmt.Println("\n\n3. Testing JWT Token Generation...")
	userID := uuid.New()
	email := "test@example.com"
	role := entities.RoleCustomer
	
	token, err := authService.GenerateToken(userID, email, role)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	fmt.Printf("✅ JWT token generated: %s...", token[:50])
	
	// Test JWT token validation
	fmt.Println("\n\n4. Testing JWT Token Validation...")
	claims, err := authService.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.UserID != userID.String() {
		log.Fatalf("Token validation failed: UserID mismatch")
	}
	if claims.Email != email {
		log.Fatalf("Token validation failed: Email mismatch")
	}
	if claims.Role != role {
		log.Fatalf("Token validation failed: Role mismatch")
	}
	fmt.Printf("✅ JWT token validation successful")
	fmt.Printf("\n   - User ID: %s", claims.UserID)
	fmt.Printf("\n   - Email: %s", claims.Email)
	fmt.Printf("\n   - Role: %s", claims.Role)
	
	// Test token refresh
	fmt.Println("\n\n5. Testing JWT Token Refresh...")
	newToken, err := authService.RefreshToken(token)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	fmt.Printf("✅ JWT token refreshed: %s...", newToken[:50])
	
	// Test invalid token
	fmt.Println("\n\n6. Testing Invalid Token...")
	_, err = authService.ValidateToken("invalid.token.here")
	if err == nil {
		log.Fatalf("Invalid token validation should have failed")
	}
	fmt.Printf("✅ Invalid token correctly rejected")
	
	fmt.Println("\n\n🎉 All JWT authentication tests passed!")
	fmt.Println("\n📋 Test Summary:")
	fmt.Println("   ✅ Password hashing")
	fmt.Println("   ✅ Password verification") 
	fmt.Println("   ✅ Wrong password rejection")
	fmt.Println("   ✅ JWT token generation")
	fmt.Println("   ✅ JWT token validation")
	fmt.Println("   ✅ JWT token refresh")
	fmt.Println("   ✅ Invalid token rejection")
	fmt.Println("\n🚀 Authentication system is ready for use!")
}
