package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection() (*gorm.DB, error) {
	dsn := buildDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Println("Database connection established successfully")
	return db, nil
}

// buildDSN builds the database connection string from environment variables
func buildDSN() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("DB_NAME", "electricity_shop")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")
	
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, password, dbname, sslmode)
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Legacy functions for backward compatibility

func InitDatabase() error {
	var err error
	DB, err = NewPostgresConnection()
	return err
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDatabase first.")
	}
	return DB
}
