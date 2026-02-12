package routes

import (
	"net/http"

	"prime-trade/auth-service/internal/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup configures all routes for the Auth Service
func Setup(router *gin.Engine, authController *controllers.AuthController) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.GET("/validate", authController.Validate)
		}
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "auth-service",
			"status":  "healthy",
		},
		"message": "Auth Service is running",
	})
}
