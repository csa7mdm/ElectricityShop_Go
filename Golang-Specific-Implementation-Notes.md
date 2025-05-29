# Golang Implementation Specific Notes

## Package Dependencies

### Core Dependencies
```go
// Database & ORM
"gorm.io/gorm"
"gorm.io/driver/postgres"

// Web Framework
"github.com/gin-gonic/gin" // Or "github.com/labstack/echo/v4"

// Authentication
"github.com/golang-jwt/jwt/v5"

// Validation
"github.com/go-playground/validator/v10"

// UUID
"github.com/google/uuid"

// Decimal handling
"github.com/shopspring/decimal"

// Configuration
"github.com/spf13/viper"

// Logging
"github.com/sirupsen/logrus" // Or "go.uber.org/zap"

// Testing
"github.com/stretchr/testify"
"github.com/golang/mock" // For gomock

// Background Jobs
"github.com/hibiken/asynq"

// Redis
"github.com/go-redis/redis/v9"

// RabbitMQ
"github.com/streadway/amqp" // Note: This library is deprecated. Consider "github.com/rabbitmq/amqp091-go"
```

## Code Organization Patterns

### Error Handling Strategy
```go
// pkg/errors/errors.go
// (Content will be added in a later step by the main agent)
```

### Response Patterns
```go
// internal/presentation/responses/responses.go
// (Content will be added in a later step by the main agent)
```
