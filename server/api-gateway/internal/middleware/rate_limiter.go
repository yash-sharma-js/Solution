package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"prime-trade/api-gateway/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter handles rate limiting logic
type RateLimiter struct {
	config      *config.Config
	redisClient *redis.Client
	localStore  *localRateLimiter
	useRedis    bool
}

type localRateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*clientRequests
}

type clientRequests struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg *config.Config) *RateLimiter {
	rl := &RateLimiter{
		config: cfg,
		localStore: &localRateLimiter{
			requests: make(map[string]*clientRequests),
		},
	}

	// Try to connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err == nil {
		rl.redisClient = client
		rl.useRedis = true
	}

	return rl
}

// Limit returns the rate limiting middleware
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		var allowed bool
		if rl.useRedis {
			allowed = rl.checkRedis(clientIP)
		} else {
			allowed = rl.checkLocal(clientIP)
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"data":    nil,
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) checkRedis(clientIP string) bool {
	ctx := context.Background()
	key := "rate_limit:" + clientIP
	window := time.Duration(rl.config.RateLimitWindow) * time.Second

	// Increment counter
	count, err := rl.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return true // Allow on error
	}

	// Set expiry on first request
	if count == 1 {
		rl.redisClient.Expire(ctx, key, window)
	}

	return count <= int64(rl.config.RateLimit)
}

func (rl *RateLimiter) checkLocal(clientIP string) bool {
	rl.localStore.mu.Lock()
	defer rl.localStore.mu.Unlock()

	now := time.Now()
	window := time.Duration(rl.config.RateLimitWindow) * time.Second

	client, exists := rl.localStore.requests[clientIP]
	if !exists || now.After(client.resetTime) {
		rl.localStore.requests[clientIP] = &clientRequests{
			count:     1,
			resetTime: now.Add(window),
		}
		return true
	}

	client.count++
	return client.count <= rl.config.RateLimit
}
