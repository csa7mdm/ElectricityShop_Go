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
	"github.com/yourusername/electricity-shop-go/internal/presentation/routes"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	
	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting ElectricityShop API...")
	
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
	
	// Seed initial data in development
	if os.Getenv("APP_ENV") == "development" {
		appLogger.Info("Seeding initial data...")
		if err := database.SeedData(db); err != nil {
			appLogger.Warnf("Failed to seed data: %v", err)
		}
	}
	
	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	// Initialize Gin router
	router := gin.New()
	
	// Setup routes
	routes.SetupRoutes(router, db, appLogger)
	
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
		appLogger.Info("API documentation will be available at /api/v1/health")
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
