package database

import (
	"fmt"
	"log" // Or a proper logger
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities" // For migrations
)

var DB *gorm.DB

func InitDatabase() error {
	// TODO: Replace with proper configuration management (e.g., Viper)
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=electricity_shop port=5432 sslmode=disable TimeZone=UTC"
		log.Println("DATABASE_DSN not set, using default:", dsn)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established.")

	// Auto-migrate entities
	log.Println("Running auto-migrations...")
	err = DB.AutoMigrate(&entities.User{}, &entities.Category{}, &entities.Product{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	log.Println("Auto-migrations completed.")

	return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		// This scenario should ideally not happen if InitDatabase is called correctly at startup.
		// Depending on the application's needs, this could panic, return an error, or attempt to re-initialize.
		log.Fatal("Database not initialized. Call InitDatabase first.")
	}
	return DB
}
