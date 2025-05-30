package routes

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/handlers"
	"github.com/yourusername/electricity-shop-go/internal/application/queries"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/database/repositories"
	"github.com/yourusername/electricity-shop-go/internal/infrastructure/messaging"
	"github.com/yourusername/electricity-shop-go/internal/presentation/controllers"
	"github.com/yourusername/electricity-shop-go/internal/presentation/middleware"
	"github.com/yourusername/electricity-shop-go/pkg/auth"
	"github.com/yourusername/electricity-shop-go/pkg/logger"
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, appLogger logger.Logger) {
	// Initialize auth service
	authService := auth.NewAuthService(
		os.Getenv("JWT_SECRET"),
		24*time.Hour, // Token TTL
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	addressRepo := repositories.NewAddressRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	
	// Initialize event publisher
	eventPublisher := messaging.NewInMemoryEventPublisher(appLogger)
	// Setup default event handlers
	if inMemoryPublisher, ok := eventPublisher.(*messaging.InMemoryEventPublisher); ok {
		inMemoryPublisher.SetupDefaultHandlers()
	}
	
	// Initialize mediator
	mediatorInstance := mediator.NewEnhancedMediator(appLogger)
	
	// Register command handlers
	userCommandHandler := handlers.NewUserCommandHandler(userRepo, addressRepo, eventPublisher, authService, appLogger)
	productCommandHandler := handlers.NewProductCommandHandler(productRepo, categoryRepo, eventPublisher, appLogger)
	cartCommandHandler := handlers.NewCartCommandHandler(cartRepo, productRepo, userRepo, eventPublisher, appLogger)
	orderCommandHandler := handlers.NewOrderCommandHandler(orderRepo, cartRepo, productRepo, userRepo, addressRepo, paymentRepo, eventPublisher, appLogger)
	
	// Register query handlers
	userQueryHandler := handlers.NewUserQueryHandler(userRepo, addressRepo, appLogger)
	productQueryHandler := handlers.NewProductQueryHandler(productRepo, categoryRepo, appLogger)
	cartQueryHandler := handlers.NewCartQueryHandler(cartRepo, appLogger)
	orderQueryHandler := handlers.NewOrderQueryHandler(orderRepo, paymentRepo, appLogger)
	
	// Register handlers with mediator
	registerUserHandlers(mediatorInstance, userCommandHandler, userQueryHandler)
	registerProductHandlers(mediatorInstance, productCommandHandler, productQueryHandler)
	registerCartHandlers(mediatorInstance, cartCommandHandler, cartQueryHandler)
	registerOrderHandlers(mediatorInstance, orderCommandHandler, orderQueryHandler)
	
	// Initialize controllers
	userController := controllers.NewUserController(mediatorInstance, appLogger)
	productController := controllers.NewProductController(mediatorInstance, appLogger)
	categoryController := controllers.NewCategoryController(mediatorInstance, appLogger)
	cartController := controllers.NewCartController(mediatorInstance, appLogger)
	orderController := controllers.NewOrderController(mediatorInstance, appLogger)
	
	// Setup API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "electricity-shop-api"})
		})
		
		// Public authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.RegisterUser)
			auth.POST("/login", userController.Login)
			// TODO: Add refresh token endpoint
			// auth.POST("/refresh", userController.RefreshToken)
		}
		
		// Protected user routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authService, appLogger))
		{
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUserProfile)
			users.DELETE("/:id", userController.DeleteUser)
			
			// Address routes
			users.GET("/:id/addresses", userController.GetUserAddresses)
			users.POST("/:id/addresses", userController.AddAddress)
			users.PUT("/:id/addresses/:address_id", userController.UpdateAddress)
			users.DELETE("/:id/addresses/:address_id", userController.DeleteAddress)
			
			// Cart routes (protected)
			users.GET("/:user_id/cart", cartController.GetCart)
			users.GET("/:user_id/cart/summary", cartController.GetCartSummary)
			users.POST("/:user_id/cart/items", cartController.AddToCart)
			users.PUT("/:user_id/cart/items/:product_id", cartController.UpdateCartItem)
			users.DELETE("/:user_id/cart/items/:product_id", cartController.RemoveFromCart)
			users.POST("/:user_id/cart/clear", cartController.ClearCart)
		}
		
		// Admin-only user management routes
		adminUsers := api.Group("/admin/users")
		adminUsers.Use(middleware.AuthMiddleware(authService, appLogger))
		adminUsers.Use(middleware.RequireRole("admin"))
		{
			adminUsers.GET("/", userController.ListUsers)
		}
		
		// Product routes (public read, admin write)
		products := api.Group("/products")
		{
			// Public product routes
			products.GET("/", productController.ListProducts)
			products.GET("/search", productController.SearchProducts)
			products.GET("/:id", productController.GetProduct)
			products.GET("/sku/:sku", productController.GetProductBySKU)
			
			// Protected admin routes
			adminProducts := products.Group("/")
			adminProducts.Use(middleware.AuthMiddleware(authService, appLogger))
			adminProducts.Use(middleware.RequireRole("admin"))
			{
				adminProducts.POST("/", productController.CreateProduct)
				adminProducts.PUT("/:id", productController.UpdateProduct)
				adminProducts.PUT("/:id/stock", productController.UpdateProductStock)
				adminProducts.DELETE("/:id", productController.DeleteProduct)
				adminProducts.GET("/low-stock", productController.GetLowStockProducts)
			}
		}
		
		// Category routes (public read, admin write)
		categories := api.Group("/categories")
		{
			// Public category routes
			categories.GET("/", categoryController.ListCategories)
			categories.GET("/root", categoryController.GetRootCategories)
			categories.GET("/:id", categoryController.GetCategory)
			categories.GET("/slug/:slug", categoryController.GetCategoryBySlug)
			categories.GET("/:id/children", categoryController.GetCategoryChildren)
			
			// Protected admin routes
			adminCategories := categories.Group("/")
			adminCategories.Use(middleware.AuthMiddleware(authService, appLogger))
			adminCategories.Use(middleware.RequireRole("admin"))
			{
				adminCategories.POST("/", categoryController.CreateCategory)
				adminCategories.PUT("/:id", categoryController.UpdateCategory)
				adminCategories.DELETE("/:id", categoryController.DeleteCategory)
			}
		}
		
		// Order routes (protected)
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(authService, appLogger))
		{
			orders.POST("/", orderController.CreateOrder)
			orders.POST("/from-cart", orderController.CreateOrderFromCart)
			orders.GET("/", orderController.ListOrders)
			orders.GET("/summary", orderController.GetOrderSummary)
			orders.GET("/:id", orderController.GetOrder)
			orders.GET("/number/:number", orderController.GetOrderByNumber)
			orders.POST("/:id/cancel", orderController.CancelOrder)
			orders.POST("/:id/payment", orderController.ProcessPayment)
			orders.GET("/:id/payments", orderController.GetOrderPayments)
			
			// Admin-only order routes
			adminOrders := orders.Group("/")
			adminOrders.Use(middleware.RequireRole("admin"))
			{
				adminOrders.GET("/to-process", orderController.GetOrdersToProcess)
				adminOrders.PUT("/:id/status", orderController.UpdateOrderStatus)
			}
		}
	}
	
	// Setup middleware
	setupMiddleware(router, appLogger)
}

// registerUserHandlers registers user command and query handlers with the mediator
func registerUserHandlers(med *mediator.EnhancedMediator, cmdHandler *handlers.UserCommandHandler, queryHandler *handlers.UserQueryHandler) {
	// Register command handlers
	med.RegisterCommandHandler(&commands.RegisterUserCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateUserProfileCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.DeleteUserCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.AddAddressCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateAddressCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.DeleteAddressCommand{}, cmdHandler)
	
	// Register query handlers
	med.RegisterQueryHandler(&queries.GetUserByIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetUserByEmailQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.ListUsersQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetUserAddressesQuery{}, queryHandler)
}

// registerProductHandlers registers product command and query handlers with the mediator
func registerProductHandlers(med *mediator.EnhancedMediator, cmdHandler *handlers.ProductCommandHandler, queryHandler *handlers.ProductQueryHandler) {
	// Register command handlers
	med.RegisterCommandHandler(&commands.CreateProductCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateProductCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateProductStockCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.DeleteProductCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.CreateCategoryCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateCategoryCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.DeleteCategoryCommand{}, cmdHandler)
	
	// Register query handlers
	med.RegisterQueryHandler(&queries.GetProductByIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetProductBySKUQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.ListProductsQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.SearchProductsQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetProductsByCategoryQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetLowStockProductsQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCategoryByIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCategoryBySlugQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.ListCategoriesQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCategoryChildrenQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetRootCategoriesQuery{}, queryHandler)
}

// registerCartHandlers registers cart command and query handlers with the mediator
func registerCartHandlers(med *mediator.EnhancedMediator, cmdHandler *handlers.CartCommandHandler, queryHandler *handlers.CartQueryHandler) {
	// Register command handlers
	med.RegisterCommandHandler(&commands.AddToCartCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateCartItemCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.RemoveFromCartCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.ClearCartCommand{}, cmdHandler)
	
	// Register query handlers
	med.RegisterQueryHandler(&queries.GetCartByUserIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCartByIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCartItemsQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetCartSummaryQuery{}, queryHandler)
}

// registerOrderHandlers registers order command and query handlers with the mediator
func registerOrderHandlers(med *mediator.EnhancedMediator, cmdHandler *handlers.OrderCommandHandler, queryHandler *handlers.OrderQueryHandler) {
	// Register command handlers
	med.RegisterCommandHandler(&commands.CreateOrderCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.CreateOrderFromCartCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.UpdateOrderStatusCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.CancelOrderCommand{}, cmdHandler)
	med.RegisterCommandHandler(&commands.ProcessPaymentCommand{}, cmdHandler)

	// Register query handlers
	med.RegisterQueryHandler(&queries.GetOrderByIDQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetOrderByNumberQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.ListOrdersQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetOrderSummaryQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetOrdersToProcessQuery{}, queryHandler)
	med.RegisterQueryHandler(&queries.GetOrderPaymentsQuery{}, queryHandler)
}

// setupMiddleware configures middleware for the application
func setupMiddleware(router *gin.Engine, appLogger logger.Logger) {
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Request logging middleware
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/v1/health"},
	}))
	
	// Recovery middleware
	router.Use(gin.Recovery())
	
	// Request ID middleware
	router.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	})
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// In production, you might want to use a proper UUID library
	return "req-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[len(charset)/2] // Simple implementation for demo
	}
	return string(b)
}
