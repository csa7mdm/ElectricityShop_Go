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
	"github.com/yourusername/electricity-shop-go/pkg/logger"
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
	appLogger.Info("🚀 Starting ElectricityShop API (Minimal Version)...")
	
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
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"features": []string{"health-check", "cors", "basic-routing"},
		})
	})
	
	// API info endpoint
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ElectricityShop API - Minimal Working Version",
			"version": "1.0.0",
			"status": "operational",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"endpoints": gin.H{
				"health": "GET /health",
				"info": "GET /api/v1/info",
				"auth": gin.H{
					"register": "POST /api/v1/auth/register (coming soon)",
					"login": "POST /api/v1/auth/login (coming soon)",
				},
			},
			"next_steps": []string{
				"Database connection integration",
				"User registration endpoint",
				"User login endpoint", 
				"JWT token authentication",
				"Protected routes",
			},
		})
	})
	
	// Authentication endpoints (placeholders for now)
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "User registration endpoint",
					"status": "placeholder",
					"note": "Implementation in progress - authentication system is ready",
					"expected_payload": gin.H{
						"email": "user@example.com",
						"password": "password123",
					},
				})
			})
			
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "User login endpoint", 
					"status": "placeholder",
					"note": "Implementation in progress - JWT system is ready",
					"expected_payload": gin.H{
						"email": "user@example.com",
						"password": "password123",
					},
					"expected_response": gin.H{
						"id": "uuid",
						"email": "user@example.com",
						"role": "customer",
						"token": "jwt_token_here",
					},
				})
			})
		}
		
		// Test endpoint to verify routing works
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "API routing is working correctly!",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"headers": c.Request.Header,
			})
		})
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
		appLogger.Infof("🌟 Server starting on port %s", port)
		appLogger.Infof("🔧 Environment: %s", os.Getenv("APP_ENV"))
		appLogger.Info("📋 Available endpoints:")
		appLogger.Info("   📍 GET  /health - Health check")
		appLogger.Info("   📍 GET  /api/v1/info - API information")
		appLogger.Info("   📍 GET  /api/v1/test - Test endpoint")
		appLogger.Info("   🔐 POST /api/v1/auth/register - User registration (placeholder)")
		appLogger.Info("   🔐 POST /api/v1/auth/login - User login (placeholder)")
		appLogger.Info("")
		appLogger.Infof("🌐 Server ready at: http://localhost:%s", port)
		appLogger.Info("✅ Try: curl http://localhost:" + port + "/health")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("🛑 Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	appLogger.Info("✅ Server exited gracefully")
}
