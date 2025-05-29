package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	appLogger "github.com/yourusername/electricity-shop-go/pkg/logger"
)

// Config holds database configuration
type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvAsInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "electricity_shop"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		MaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
	}
	return config
}

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection() (*gorm.DB, error) {
	config := LoadConfig()
	
	// Configure GORM logger
	var gormLogger logger.Interface
	if os.Getenv("APP_ENV") == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}
	
	// Connect to database
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)
	
	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	appLogger.Get().Info("Successfully connected to PostgreSQL database")
	return db, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}
