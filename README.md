# ElectricityShop Go API

A modern, scalable e-commerce API built with Go, following Clean Architecture principles and implementing CQRS pattern with domain-driven design.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/api/                 # Application entry point
├── internal/
│   ├── domain/             # Domain layer (entities, interfaces, events)
│   ├── application/        # Application layer (commands, queries, handlers)
│   ├── infrastructure/     # Infrastructure layer (database, messaging, external services)
│   └── presentation/       # Presentation layer (controllers, routes, middleware)
├── pkg/                    # Shared utilities (logger, errors, mediator)
└── web/                    # Frontend assets (future)
```

## 🚀 Features

### Core Features
- **User Management**: Registration, authentication, profile management
- **Product Catalog**: Products, categories, inventory tracking
- **Shopping Cart**: Add/remove items, quantity management
- **Order Management**: Order creation, status tracking, history
- **Address Management**: Multiple shipping/billing addresses
- **Inventory Tracking**: Stock levels, low stock alerts

### Technical Features
- **Clean Architecture**: Domain-driven design with clear boundaries
- **CQRS Pattern**: Separate read and write operations
- **Event-Driven**: Domain events for loose coupling
- **Database**: PostgreSQL with GORM ORM
- **API**: RESTful JSON API with Gin framework
- **Logging**: Structured logging with Logrus
- **Validation**: Request validation and error handling
- **Graceful Shutdown**: Proper server lifecycle management

## 🛠️ Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Logging**: Logrus
- **Architecture**: Clean Architecture + CQRS
- **Patterns**: Repository, Mediator, Event-Driven

## 📋 Prerequisites

- Go 1.23 or higher
- PostgreSQL 13+ 
- Git

## 🔧 Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/electricity-shop-go.git
cd electricity-shop-go
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Set up PostgreSQL Database
```bash
# Create database
createdb electricity_shop

# Or using psql
psql -U postgres -c "CREATE DATABASE electricity_shop;"
```

### 4. Configure Environment Variables
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your database credentials
nano .env
```

Required environment variables:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=electricity_shop

# Application
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=info
```

### 5. Run Database Migrations & Seed Data
The application automatically runs migrations and seeds initial data on startup in development mode.

### 6. Start the Application
```bash
# Development mode
go run cmd/api/main.go

# Or build and run
go build -o bin/api cmd/api/main.go
./bin/api
```

The API will be available at: `http://localhost:8080`

## 📚 API Documentation

### Health Check
```bash
GET /api/v1/health
```

### User Endpoints
```bash
# Register user
POST /api/v1/users/register
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}

# Get user
GET /api/v1/users/{id}

# List users (with filtering)
GET /api/v1/users?page=1&page_size=10&search=john

# Update user profile
PUT /api/v1/users/{id}
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}

# User addresses
GET /api/v1/users/{id}/addresses
POST /api/v1/users/{id}/addresses
PUT /api/v1/users/{id}/addresses/{address_id}
DELETE /api/v1/users/{id}/addresses/{address_id}
```

### Product Endpoints
```bash
# Create product
POST /api/v1/products
{
  "name": "LED Light Bulb",
  "description": "Energy efficient LED bulb",
  "sku": "LED-001",
  "price": "12.99",
  "category_id": "uuid",
  "stock": 100
}

# Get product
GET /api/v1/products/{id}
GET /api/v1/products/sku/{sku}

# List products (with filtering)
GET /api/v1/products?page=1&page_size=10&category_id=uuid&min_price=10&max_price=100

# Search products
GET /api/v1/products/search?q=LED

# Update product
PUT /api/v1/products/{id}

# Update stock
PUT /api/v1/products/{id}/stock
{
  "quantity": 50,
  "reason": "Stock adjustment"
}

# Get low stock products
GET /api/v1/products/low-stock?threshold=10
```

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/application/handlers/...
```

### Test with curl
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Register user
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"Test","last_name":"User"}'
```

## 📁 Project Structure Details

### Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (User, Product, Order, etc.)
- **Interfaces**: Repository and service contracts
- **Events**: Domain events for business logic
- **Value Objects**: Enums and business rules

### Application Layer (`internal/application/`)
- **Commands**: Write operations (Create, Update, Delete)
- **Queries**: Read operations (Get, List, Search)
- **Handlers**: Command and query processors
- **DTOs**: Data transfer objects

### Infrastructure Layer (`internal/infrastructure/`)
- **Database**: Repositories, migrations, connections
- **Messaging**: Event publishers and handlers
- **External**: Third-party service integrations

### Presentation Layer (`internal/presentation/`)
- **Controllers**: HTTP request handlers
- **Routes**: API route definitions
- **Middleware**: Cross-cutting concerns
- **Responses**: API response structures

## 🔒 Security Considerations

- Input validation on all endpoints
- SQL injection prevention with GORM
- Password hashing with bcrypt
- CORS configuration
- Request rate limiting (planned)
- JWT authentication (planned)

## 🚀 Deployment

### Docker (Coming Soon)
```bash
# Build image
docker build -t electricity-shop-api .

# Run with docker-compose
docker-compose up -d
```

### Production Environment Variables
```env
APP_ENV=production
DB_SSLMODE=require
LOG_LEVEL=warn
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Failed**
```bash
# Check PostgreSQL is running
pg_isready -U postgres

# Verify database exists
psql -U postgres -l | grep electricity_shop
```

**Port Already in Use**
```bash
# Check what's using port 8080
lsof -i :8080

# Use different port
export APP_PORT=8081
```

**Migration Errors**
```bash
# Check database permissions
psql -U postgres -d electricity_shop -c "\\dp"
```

## 📞 Support

For support, email support@electricityshop.com or create an issue on GitHub.

## 🗺️ Roadmap

- [ ] JWT Authentication & Authorization
- [ ] Shopping Cart & Checkout
- [ ] Order Management System
- [ ] Payment Integration (Stripe, PayPal)
- [ ] Email Notifications
- [ ] File Upload for Product Images
- [ ] Advanced Search & Filtering
- [ ] Admin Dashboard
- [ ] API Rate Limiting
- [ ] Docker Containerization
- [ ] CI/CD Pipeline
- [ ] Comprehensive Test Suite
- [ ] API Documentation (Swagger)
- [ ] Monitoring & Metrics
- [ ] Caching Layer (Redis)
