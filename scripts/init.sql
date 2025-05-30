-- Initialize database for ElectricityShop
-- This script is automatically executed when the PostgreSQL container starts

-- Create database if it doesn't exist (already handled by POSTGRES_DB)
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a test user for development
-- This is handled by the application's seed data, but included here as backup
-- The application will handle all table creation via GORM migrations
