# ElectricityShop Go API - Authentication Testing Guide

## 🚀 Quick Start

### 1. Start the Application

```bash
# Set up environment
cp .env.example .env
# Edit .env file with your database credentials and JWT secret

# Run the application
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### 2. Test Authentication Endpoints

#### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": null
}
```

#### Login User

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "test@example.com",
    "role": "customer",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### Access Protected Endpoint

```bash
# Replace TOKEN with the JWT token from login response
curl -X GET http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer TOKEN"
```

## 🔐 Authentication Flow

### Public Endpoints (No Authentication Required)
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/products` - List products
- `GET /api/v1/products/:id` - Get product details
- `GET /api/v1/categories` - List categories
- `GET /api/v1/health` - Health check

### Protected Endpoints (Authentication Required)
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user profile
- `POST /api/v1/users/:id/addresses` - Add user address
- `GET /api/v1/users/:user_id/cart` - Get user's cart
- `POST /api/v1/orders` - Create order

### Admin-Only Endpoints (Admin Role Required)
- `GET /api/v1/admin/users` - List all users
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `POST /api/v1/categories` - Create category
- `PUT /api/v1/orders/:id/status` - Update order status

## 🧪 Test JWT Authentication System

Run the authentication test:

```bash
go run tests/auth_test.go
```

This will test:
- Password hashing and verification
- JWT token generation and validation
- Token refresh functionality
- Invalid token rejection

## 📋 API Error Responses

### Validation Errors (400)
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR"
}
```

### Authentication Errors (401)
```json
{
  "success": false,
  "message": "Invalid credentials",
  "error": "INVALID_CREDENTIALS"
}
```

### Authorization Errors (403)
```json
{
  "success": false,
  "message": "Insufficient permissions",
  "error": "INSUFFICIENT_PERMISSIONS"
}
```

### Resource Not Found (404)
```json
{
  "success": false,
  "message": "User not found",
  "error": "USER_NOT_FOUND"
}
```

### Conflict Errors (409)
```json
{
  "success": false,
  "message": "User already exists",
  "error": "USER_ALREADY_EXISTS"
}
```

## 🔧 Environment Configuration

Make sure your `.env` file includes:

```env
# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY_HOURS=24

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=electricity_shop
DB_SSLMODE=disable

# Application Configuration
APP_ENV=development
APP_PORT=8080
```

## 🚧 Current Implementation Status

### ✅ Completed
- JWT token generation and validation
- Password hashing with bcrypt
- User registration and login
- Protected route middleware
- Role-based authorization
- Comprehensive error handling
- CORS middleware

### 🔄 In Progress
- Complete DTO validation
- Enhanced error mapping
- API documentation

### 📝 Next Steps
1. Add missing DTOs for address and profile updates
2. Implement comprehensive input validation
3. Add Swagger/OpenAPI documentation
4. Set up automated testing
5. Add rate limiting
6. Implement refresh token rotation

## 🐛 Troubleshooting

### Common Issues

1. **"Database connection failed"**
   - Check your database is running
   - Verify connection details in `.env`
   - Ensure database exists

2. **"Invalid token"**
   - Check token format in Authorization header
   - Ensure token hasn't expired
   - Verify JWT_SECRET in environment

3. **"Validation failed"**
   - Check request body format
   - Ensure required fields are present
   - Verify email format and password length

4. **"Insufficient permissions"**
   - Check user role in JWT token
   - Verify endpoint requires correct role
   - Ensure user is authenticated

### Debug Tips

1. **Enable debug logging:**
   ```env
   LOG_LEVEL=debug
   ```

2. **Check JWT token contents:**
   Use [jwt.io](https://jwt.io) to decode and inspect tokens

3. **Test with curl:**
   Add `-v` flag for verbose output to see full HTTP request/response

4. **Check application logs:**
   Look for specific error messages in console output
