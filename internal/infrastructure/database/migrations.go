package database

import (
	"gorm.io/gorm"
	"github.com/yourusername/electricity-shop-go/internal/domain/entities"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}
	
	// Auto-migrate all entities
	return db.AutoMigrate(
		// User-related entities
		&entities.User{},
		&entities.Address{},
		
		// Product-related entities
		&entities.Category{},
		&entities.Product{},
		
		// Cart-related entities
		&entities.Cart{},
		&entities.CartItem{},
		
		// Order-related entities
		&entities.Order{},
		&entities.OrderItem{},
		&entities.Payment{},
		&entities.Shipment{},
	)
}

// SeedData creates initial seed data for development
func SeedData(db *gorm.DB) error {
	// Create admin user if not exists
	var adminCount int64
	db.Model(&entities.User{}).Where("role = ?", entities.RoleAdmin).Count(&adminCount)
	
	if adminCount == 0 {
		adminUser := &entities.User{
			Email:    "admin@electricityshop.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: password
			Role:     entities.RoleAdmin,
			IsActive: true,
		}
		if err := db.Create(adminUser).Error; err != nil {
			return err
		}
	}
	
	// Create sample categories
	var categoryCount int64
	db.Model(&entities.Category{}).Count(&categoryCount)
	
	if categoryCount == 0 {
		categories := []*entities.Category{
			{
				Name:        "Electrical Components",
				Slug:        "electrical-components",
				Description: "Basic electrical components and parts",
				IsActive:    true,
				SortOrder:   1,
			},
			{
				Name:        "Cables & Wires",
				Slug:        "cables-wires",
				Description: "Various types of cables and wires",
				IsActive:    true,
				SortOrder:   2,
			},
			{
				Name:        "Tools & Equipment",
				Slug:        "tools-equipment",
				Description: "Electrical tools and equipment",
				IsActive:    true,
				SortOrder:   3,
			},
			{
				Name:        "Lighting",
				Slug:        "lighting",
				Description: "LED lights, bulbs, and fixtures",
				IsActive:    true,
				SortOrder:   4,
			},
			{
				Name:        "Safety Equipment",
				Slug:        "safety-equipment",
				Description: "Safety gear and protective equipment",
				IsActive:    true,
				SortOrder:   5,
			},
		}
		
		for _, category := range categories {
			if err := db.Create(category).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}
