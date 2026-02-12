package config

import (
	"os"
	"strconv"
)

// Config holds the API Gateway configuration
type Config struct {
	Port            string
	AuthServiceURL  string
	TaskServiceURL  string
	RedisURL        string
	GinMode         string
	RateLimit       int
	RateLimitWindow int
}

// Load loads configuration from environment variables
func Load() *Config {
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))

	return &Config{
		Port:            getEnv("PORT", "8080"),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		TaskServiceURL:  getEnv("TASK_SERVICE_URL", "http://localhost:8082"),
		RedisURL:        getEnv("REDIS_URL", "localhost:6379"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		RateLimit:       rateLimit,
		RateLimitWindow: rateLimitWindow,
	}
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
