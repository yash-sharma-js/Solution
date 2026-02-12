package services

import (
	"context"
	"errors"
	"time"

	"prime-trade/auth-service/internal/config"
	"prime-trade/auth-service/internal/models"
	"prime-trade/auth-service/internal/repositories"
	"prime-trade/auth-service/internal/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Common errors
var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*models.ValidateResponse, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	config      *config.Config
	redisClient *redis.Client
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		config:      cfg,
		redisClient: redisClient,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = models.RoleUser
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, err
	}

	// Cache token in Redis if available
	if s.redisClient != nil {
		cacheKey := "token:" + token
		s.redisClient.Set(ctx, cacheKey, user.ID.String(), s.config.JWTExpiry)
	}

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns user information
func (s *authService) ValidateToken(ctx context.Context, token string) (*models.ValidateResponse, error) {
	// Check cache first if Redis is available
	if s.redisClient != nil {
		cacheKey := "validated:" + token
		cached, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			// Token was previously validated, get user info
			userID, _ := uuid.Parse(cached)
			user, err := s.userRepo.FindByID(ctx, userID)
			if err == nil && user != nil {
				return &models.ValidateResponse{
					Valid:  true,
					User:   user.ToResponse(),
					UserID: user.ID.String(),
					Role:   string(user.Role),
				}, nil
			}
		}
	}

	// Validate token
	claims, err := utils.ValidateToken(token, s.config.JWTSecret)
	if err != nil {
		return &models.ValidateResponse{Valid: false}, ErrInvalidToken
	}

	// Get user from database
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return &models.ValidateResponse{Valid: false}, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return &models.ValidateResponse{Valid: false}, ErrUserNotFound
	}

	// Cache validated token if Redis is available
	if s.redisClient != nil {
		cacheKey := "validated:" + token
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			s.redisClient.Set(ctx, cacheKey, user.ID.String(), ttl)
		}
	}

	return &models.ValidateResponse{
		Valid:  true,
		User:   user.ToResponse(),
		UserID: user.ID.String(),
		Role:   string(user.Role),
	}, nil
}
