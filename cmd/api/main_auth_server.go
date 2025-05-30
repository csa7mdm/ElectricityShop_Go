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
	"github.com/joho/godotenv"

	// Added/Verified Imports
	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/handlers"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database/repositories"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/messaging"
	"github.com/yourusername/electricity-shop-go/internal/presentation/controllers"
	"github.com/yourusername/electricity-shop-go/internal/presentation/middleware"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
	// gorm "gorm.io/gorm" // Implicitly used by database package
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.test"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("🚀 Starting ElectricityShop Authentication Server...")

	// Initialize Database
	db, err := database.NewPostgresConnection()
	if err != nil {
		appLogger.Fatalf("🚨 Failed to connect to database: %v", err)
	}
	appLogger.Info("🔗 Database connection established.")

	// Run Migrations
	if err := database.RunMigrations(db); err != nil {
		appLogger.Fatalf("🚨 Failed to run database migrations: %v", err)
	}
	appLogger.Info("🔄 Database migrations completed.")

	// Initialize auth service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		appLogger.Warn("⚠️  Using default JWT secret - change in production!")
	}
	authService := auth.NewAuthService(jwtSecret, 24*time.Hour)

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize Event Publisher
	eventPublisher := messaging.NewInMemoryEventPublisher(appLogger)

	// Initialize Command Handlers
	userCommandHandler := handlers.NewUserCommandHandler(userRepo, nil, eventPublisher, authService, appLogger)

	// Initialize Mediator
	mediatorInstance := mediator.NewConcreteMediator(appLogger) // Assuming NewConcreteMediator

	// Register Handlers with Mediator
	if err := mediatorInstance.RegisterCommandHandler(&commands.RegisterUserCommand{}, userCommandHandler); err != nil {
		appLogger.Fatalf("🚨 Failed to register RegisterUserCommand handler: %v", err)
	}
	if err := mediatorInstance.RegisterQueryHandler(&commands.LoginUserCommand{}, userCommandHandler); err != nil {
		appLogger.Fatalf("🚨 Failed to register LoginUserCommand handler: %v", err)
	}
	appLogger.Info("🔗 Mediator initialized and handlers registered.")

	// Initialize controllers
	userController := controllers.NewSimpleUserController(mediatorInstance, appLogger) // Updated to use mediator

	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	// Initialize router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS middleware
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
	
	setupRoutes(router, authService, userController, appLogger)
	startServer(router, appLogger)
}

func setupRoutes(router *gin.Engine, authService *auth.AuthService, userController *controllers.SimpleUserController, appLogger logger.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "electricity-shop-auth-api",
			"version": "1.0.0",
			"features": []string{"jwt-auth", "user-management", "protected-routes"},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// API info
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ElectricityShop Authentication API",
			"version": "1.0.0",
			"description": "JWT-based authentication system with user management",
			"endpoints": gin.H{
				"public": []string{
					"GET /health",
					"GET /api/v1/info", 
					"POST /api/v1/auth/register",
					"POST /api/v1/auth/login",
					"GET /api/v1/test/token",
				},
				"protected": []string{
					"GET /api/v1/users/:id",
					"PUT /api/v1/users/:id", 
					"GET /api/v1/test/auth",
				},
			},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// API routes
	api := router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.RegisterUser)
			auth.POST("/login", userController.Login)
		}
		
		// Protected user routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authService, appLogger))
		{
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUserProfile)
		}
		
		// Test endpoints
		api.GET("/test/token", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(200, gin.H{
					"message": "Token validation test endpoint",
					"status": "no_token_provided",
					"instructions": "Include 'Authorization: Bearer <token>' header to test validation",
				})
				return
			}
			
			c.JSON(200, gin.H{
				"message": "Token detected",
				"status": "token_provided",
				"note": "Use /api/v1/test/auth for full authentication test",
			})
		})
		
		api.GET("/test/auth", middleware.AuthMiddleware(authService, appLogger), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "🎉 Authentication successful!",
				"user": gin.H{
					"id":    c.GetString("user_id"),
					"email": c.GetString("user_email"),
					"role":  c.GetString("user_role"),
				},
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
	}
}

func startServer(router *gin.Engine, appLogger logger.Logger) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	go func() {
		appLogger.Infof("🌟 Authentication Server running on port %s", port)
		appLogger.Infof("🌐 Server URL: http://localhost:%s", port)
		appLogger.Info("")
		appLogger.Info("📋 Available endpoints:")
		appLogger.Info("  🏥 GET  /health")
		appLogger.Info("  ℹ️  GET  /api/v1/info")
		appLogger.Info("  📝 POST /api/v1/auth/register")
		appLogger.Info("  🔑 POST /api/v1/auth/login") 
		appLogger.Info("  🧪 GET  /api/v1/test/token")
		appLogger.Info("  👤 GET  /api/v1/users/:id (protected)")
		appLogger.Info("  ✏️  PUT  /api/v1/users/:id (protected)")
		appLogger.Info("  🔐 GET  /api/v1/test/auth (protected)")
		appLogger.Info("")
		appLogger.Info("🚀 Quick test: curl http://localhost:" + port + "/health")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close database connection
	if sqlDB, errDb := db.DB(); errDb == nil {
		appLogger.Info("Closing database connection...")
		if errClose := sqlDB.Close(); errClose != nil {
			appLogger.Errorf("Error closing database: %v", errClose)
		}
	}

	appLogger.Info("✅ Server exited gracefully")
}
