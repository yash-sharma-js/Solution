package routes

import (
	"net/http"

	"prime-trade/task-service/internal/controllers"
	"prime-trade/task-service/internal/middleware"
	"prime-trade/task-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup configures all routes for the Task Service
func Setup(router *gin.Engine, taskController *controllers.TaskController, authClient *services.AuthClient) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.Auth(authClient))
		{
			tasks.POST("", taskController.Create)
			tasks.GET("", taskController.List)
			tasks.GET("/:id", taskController.GetByID)
			tasks.PUT("/:id", taskController.Update)
			tasks.DELETE("/:id", taskController.Delete)
		}
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "task-service",
			"status":  "healthy",
		},
		"message": "Task Service is running",
	})
}
