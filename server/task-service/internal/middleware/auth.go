package middleware

import (
	"net/http"
	"strings"

	"prime-trade/task-service/internal/models"
	"prime-trade/task-service/internal/services"

	"github.com/gin-gonic/gin"
)

// Auth creates an authentication middleware
func Auth(authClient *services.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required"))
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid authorization format. Use: Bearer <token>"))
			c.Abort()
			return
		}

		// Validate token with auth service
		user, err := authClient.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}
