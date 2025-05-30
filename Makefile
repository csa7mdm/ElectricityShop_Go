# ElectricityShop Go - Makefile for development tasks

.PHONY: help build run test clean docker-build docker-run docker-stop logs deps fmt vet lint

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install/update dependencies"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run golangci-lint (requires golangci-lint installation)"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run with Docker Compose"
	@echo "  docker-stop - Stop Docker Compose"
	@echo "  docker-logs - View Docker logs"
	@echo "  db-up       - Start only database with Docker"
	@echo "  db-down     - Stop database"

# Application name and version
APP_NAME := electricity-shop-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) cmd/api/main.go
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application locally
run:
	@echo "Starting $(APP_NAME)..."
	@go run cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean ./...

# Install/update dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go mod verify

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

# Start all services with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started. API available at http://localhost:8080"
	@echo "Adminer (DB UI) available at http://localhost:8081"

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# View Docker logs
docker-logs:
	@docker-compose logs -f

# Start only database
db-up:
	@echo "Starting database..."
	@docker-compose up -d postgres redis
	@echo "Database services started"

# Stop database
db-down:
	@echo "Stopping database..."
	@docker-compose stop postgres redis

# Development workflow
dev: deps fmt vet test build
	@echo "Development build complete"

# Production build
prod: deps test
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(LDFLAGS) -w -s" -o bin/$(APP_NAME) cmd/api/main.go
	@echo "Production build complete"

# Create .env file from example
env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
		echo "Please update .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Database operations (requires running database)
db-migrate:
	@echo "Running database migrations..."
	@go run cmd/api/main.go --migrate-only

# Generate API documentation (if you add swagger later)
docs:
	@echo "Generating API documentation..."
	@echo "Swagger documentation generation not implemented yet"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

# Check if everything is working
check: fmt vet test
	@echo "All checks passed!"

# Setup development environment
setup: env deps install-tools db-up
	@echo "Development environment setup complete!"
	@echo "Run 'make run' to start the application"
