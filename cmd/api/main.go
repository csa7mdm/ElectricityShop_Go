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
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database"
	"github.com/yourusername/electricity-shop-go/internal/presentation/routes"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
)

func main() {
	// Initialize logger
	appLogger := logger.NewLogger()
	
	// Initialize database
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Setup routes
	routes.SetupRoutes(router, db, appLogger)
	
	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	
	// Start server in a goroutine
	go func() {
		appLogger.Info("Starting server on :8080")
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
	
	appLogger.Info("Server exited")
}
