package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prime-trade/task-service/internal/config"
	"prime-trade/task-service/internal/controllers"
	"prime-trade/task-service/internal/middleware"
	"prime-trade/task-service/internal/models"
	"prime-trade/task-service/internal/repositories"
	"prime-trade/task-service/internal/routes"
	"prime-trade/task-service/internal/services"

	_ "prime-trade/task-service/docs"

	"github.com/gin-gonic/gin"
)

// @title Prime Trade Task Service API
// @version 1.0
// @description Task Management Service for Prime Trade
// @host localhost:8082
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Redis (optional)
	redisClient := config.InitRedis(cfg)

	// Initialize HTTP client for auth service
	authClient := services.NewAuthClient(cfg.AuthServiceURL)

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(db)

	// Initialize services
	taskService := services.NewTaskService(taskRepo, redisClient)

	// Initialize controllers
	taskController := controllers.NewTaskController(taskService)

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Setup routes
	routes.Setup(router, taskController, authClient)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Task Service starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Task Service...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close Redis connection
	if redisClient != nil {
		redisClient.Close()
	}

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Task Service stopped")
}
