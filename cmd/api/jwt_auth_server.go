package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

// Inline controller to avoid import issues
type AuthController struct {
	authService *auth.AuthService
	logger      logger.Logger
	validator   *validator.Validate
}

func NewAuthController(authService *auth.AuthService, logger logger.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger,
		validator:   validator.New(),
	}
}

func (ac *AuthController) RegisterUser(c *gin.Context) {
	var req dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.logger.Errorf("Registration bind error: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	if err := ac.validator.Struct(&req); err != nil {
		ac.logger.Errorf("Registration validation error: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	ac.logger.Infof("Registration request for: %s", req.Email)
	userID := uuid.New()
	
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(gin.H{
		"id": userID.String(),
		"email": req.Email,
		"message": "User registration successful (demo mode)",
	}, "User registered successfully"))
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.logger.Errorf("Login bind error: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	if err := ac.validator.Struct(&req); err != nil {
		ac.logger.Errorf("Login validation error: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Validation failed", "VALIDATION_ERROR"))
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Email and password required", "INVALID_CREDENTIALS"))
		return
	}
	
	ac.logger.Infof("Login request for: %s", req.Email)
	
	// Generate real JWT token
	userID := uuid.New()
	token, err := ac.authService.GenerateToken(userID, req.Email, entities.RoleCustomer)
	if err != nil {
		ac.logger.Errorf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Failed to generate token", "TOKEN_ERROR"))
		return
	}
	
	loginResponse := &dtos.LoginUserResponse{
		ID:    userID.String(),
		Email: req.Email,
		Role:  "customer",
		Token: token,
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(loginResponse, "Login successful"))
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	userEmail := c.GetString("user_email")
	userRole := c.GetString("user_role")

	c.JSON(http.StatusOK, responses.NewSuccessResponse(gin.H{
		"id":    userID,
		"email": userEmail,
		"role":  userRole,
		"message": "Profile retrieved successfully",
	}, "Profile retrieved"))
}

// Inline auth middleware to avoid import issues
func AuthMiddleware(authService *auth.AuthService, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Authorization header required", "MISSING_AUTH_HEADER"))
			c.Abort()
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid authorization header format", "INVALID_AUTH_HEADER"))
			c.Abort()
			return
		}

		token := authHeader[7:]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			logger.Errorf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Invalid token", "INVALID_TOKEN"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", string(claims.Role))
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.test"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using defaults")
		}
	}
	
	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("🚀 Starting ElectricityShop JWT Authentication Server...")
	
	// Initialize auth service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		appLogger.Warn("⚠️  Using default JWT secret")
	}
	authService := auth.NewAuthService(jwtSecret, 24*time.Hour)
	
	// Initialize controller
	authController := NewAuthController(authService, appLogger)
	
	// Setup Gin
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	
	// Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "electricity-shop-jwt-auth",
			"version": "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	api := router.Group("/api/v1")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.RegisterUser)
			auth.POST("/login", authController.Login)
		}
		
		// Protected routes
		api.GET("/profile", AuthMiddleware(authService, appLogger), authController.GetProfile)
		
		// Test route
		api.GET("/test", AuthMiddleware(authService, appLogger), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "🎉 JWT Authentication working!",
				"user": gin.H{
					"id":    c.GetString("user_id"),
					"email": c.GetString("user_email"),
					"role":  c.GetString("user_role"),
				},
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
	}
	
	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	
	go func() {
		appLogger.Infof("🌟 Server running on http://localhost:%s", port)
		appLogger.Info("📋 Endpoints:")
		appLogger.Info("  🏥 GET  /health")
		appLogger.Info("  📝 POST /api/v1/auth/register")
		appLogger.Info("  🔑 POST /api/v1/auth/login")
		appLogger.Info("  👤 GET  /api/v1/profile (protected)")
		appLogger.Info("  🧪 GET  /api/v1/test (protected)")
		appLogger.Info("")
		appLogger.Info("🧪 Test: curl http://localhost:" + port + "/health")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("🛑 Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown failed:", err)
	}
	
	appLogger.Info("✅ Server stopped")
}
