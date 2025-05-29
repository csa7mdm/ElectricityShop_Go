package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/electricity-shop-go/internal/presentation/controllers"
	// "github.com/yourusername/electricity-shop-go/internal/presentation/middleware" // For auth middleware later
)

func SetupUserRoutes(apiGroup *gin.RouterGroup, userController *controllers.UserController) {
	authRoutes := apiGroup.Group("/auth")
	{
		authRoutes.POST("/register", userController.Register)
		authRoutes.POST("/login", userController.Login)
		// Add other auth routes here: /refresh-token, /me, /logout etc.
	}
}
