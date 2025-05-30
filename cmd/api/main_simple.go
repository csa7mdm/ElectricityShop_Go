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
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	
	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting ElectricityShop API (Authentication Only)...")
	
	// Initialize database
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Run database migrations
	appLogger.Info("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	// Initialize Gin router
	router := gin.New()
	
	// Basic middleware
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
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok", 
			"service": "electricity-shop-api",
			"version": "1.0.0",
			"features": []string{"authentication", "user-management"},
		})
	})
	
	// API info endpoint
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ElectricityShop API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/health",
				"auth": gin.H{
					"register": "POST /api/v1/auth/register",
					"login": "POST /api/v1/auth/login",
				},
			},
		})
	})
	
	// Placeholder authentication endpoints
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Registration endpoint - Implementation in progress",
					"status": "coming_soon",
				})
			})
			
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Login endpoint - Implementation in progress", 
					"status": "coming_soon",
				})
			})
		}
	}
	
	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	
	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		appLogger.Infof("Starting server on port %s", port)
		appLogger.Infof("Environment: %s", os.Getenv("APP_ENV"))
		appLogger.Info("API endpoints:")
		appLogger.Info("  - GET  /health")
		appLogger.Info("  - GET  /api/v1/info")
		appLogger.Info("  - POST /api/v1/auth/register (placeholder)")
		appLogger.Info("  - POST /api/v1/auth/login (placeholder)")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
	
	appLogger.Info("Server exited gracefully")
}
