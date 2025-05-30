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
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database/repositories"
	"github.com/yourusername/electricity-shop-go/internal/application/handlers"
	"github.com/yourusername/electricity-shop-go/internal/presentation/controllers"
	"github.com/yourusername/electricity-shop-go/internal/presentation/middleware"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/messaging"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.test"); err != nil {
		log.Println("No .env.test file found, trying .env...")
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}
	
	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("🚀 Starting ElectricityShop API with Authentication...")
	
	// Initialize database
	appLogger.Info("📊 Connecting to database...")
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("🔄 Running in API-only mode (no database)")
		runAPIOnlyMode(appLogger)
		return
	}
	
	// Run database migrations
	appLogger.Info("🔧 Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Printf("⚠️  Migration failed: %v", err)
		log.Println("🔄 Running in API-only mode")
		runAPIOnlyMode(appLogger)
		return
	}
	
	// Initialize auth service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "development-secret-key-change-in-production"
		appLogger.Warn("JWT_SECRET not set, using default development key")
	}
	authService := auth.NewAuthService(jwtSecret, 24*time.Hour)
	
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	addressRepo := repositories.NewAddressRepository(db)
	
	// Initialize event publisher
	eventPublisher := messaging.NewInMemoryEventPublisher(appLogger)
	
	// Initialize mediator
	mediatorInstance := mediator.NewEnhancedMediator(appLogger)
	
	// Initialize handlers
	userCommandHandler := handlers.NewUserCommandHandler(userRepo, addressRepo, eventPublisher, authService, appLogger)
	
	// Register handlers with mediator
	// Note: We'll add a simplified registration for now
	
	// Initialize controllers
	userController := controllers.NewSimpleUserController(authService, appLogger)
	
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	// Initialize Gin router
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
	
	// Setup routes
	setupRoutes(router, authService, userController, appLogger)
	
	// Start server
	startServer(router, appLogger)
}

func runAPIOnlyMode(appLogger logger.Logger) {
	appLogger.Info("🔄 Starting API-only mode (no database required)")
	
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
	
	// Basic routes for API-only mode
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"mode": "api-only",
			"message": "Database not connected, running in demo mode",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ElectricityShop API - Demo Mode",
			"mode": "api-only",
			"note": "Database connection required for full functionality",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// Demo authentication endpoints
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Demo registration - database required for full functionality",
					"status": "demo_mode",
					"note": "Connect database to enable user registration",
				})
			})
			
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Demo login - database required for full functionality",
					"status": "demo_mode", 
					"note": "Connect database to enable user authentication",
				})
			})
		}
	}
	
	startServer(router, appLogger)
}

func setupRoutes(router *gin.Engine, authService *auth.AuthService, userController *controllers.SimpleUserController, appLogger logger.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "electricity-shop-api",
			"mode": "full",
			"database": "connected",
			"authentication": "enabled",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// API info
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ElectricityShop API - Full Mode",
			"version": "1.0.0",
			"features": []string{"authentication", "user-management", "database"},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// API routes
	api := router.Group("/api/v1")
	{
		// Public authentication routes
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
		
		// Test routes
		api.GET("/test/auth", middleware.AuthMiddleware(authService, appLogger), func(c *gin.Context) {
			userID := c.GetString("user_id")
			userEmail := c.GetString("user_email")
			userRole := c.GetString("user_role")
			
			c.JSON(200, gin.H{
				"message": "Authentication test successful!",
				"user_id": userID,
				"user_email": userEmail,
				"user_role": userRole,
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
		appLogger.Infof("🌟 Server starting on port %s", port)
		appLogger.Infof("🔧 Environment: %s", os.Getenv("APP_ENV"))
		appLogger.Info("📋 Available endpoints:")
		appLogger.Info("   📍 GET  /health")
		appLogger.Info("   📍 GET  /api/v1/info")
		appLogger.Info("   🔐 POST /api/v1/auth/register")
		appLogger.Info("   🔐 POST /api/v1/auth/login")
		appLogger.Info("   👤 GET  /api/v1/users/:id (protected)")
		appLogger.Info("   🧪 GET  /api/v1/test/auth (protected)")
		appLogger.Info("")
		appLogger.Infof("🌐 Server ready at: http://localhost:%s", port)
		
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
	
	appLogger.Info("✅ Server exited gracefully")
}
