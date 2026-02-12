package routes

import (
	"net/http"

	"prime-trade/api-gateway/internal/config"
	"prime-trade/api-gateway/internal/proxy"

	"github.com/gin-gonic/gin"
)

// Setup configures all routes for the API Gateway
func Setup(router *gin.Engine, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Create proxies for each service
	authProxy := proxy.NewReverseProxy(cfg.AuthServiceURL)
	taskProxy := proxy.NewReverseProxy(cfg.TaskServiceURL)

	// Auth service routes
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.Any("/*path", authProxy.Handler())
	}

	// Task service routes
	taskRoutes := router.Group("/api/v1/tasks")
	{
		taskRoutes.Any("", taskProxy.Handler())
		taskRoutes.Any("/*path", taskProxy.Handler())
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "api-gateway",
			"status":  "healthy",
		},
		"message": "API Gateway is running",
	})
}
