package controllers

import (
	"errors"
	"net/http"
	"strings"

	"prime-trade/auth-service/internal/models"
	"prime-trade/auth-service/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(formatValidationError(err)))
		return
	}

	response, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to register user"))
		return
	}

	ctx.JSON(http.StatusCreated, models.NewSuccessResponse(response, "User registered successfully"))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(formatValidationError(err)))
		return
	}

	response, err := c.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to login"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(response, "Login successful"))
}

// Validate godoc
// @Summary Validate JWT token
// @Description Validate JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.ValidateResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/validate [get]
func (c *AuthController) Validate(ctx *gin.Context) {
	// Extract token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required"))
		return
	}

	// Remove "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid authorization format"))
		return
	}

	response, err := c.authService.ValidateToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid or expired token"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(response, "Token is valid"))
}

// formatValidationError formats validation errors into readable messages
func formatValidationError(err error) string {
	return "Validation failed: " + err.Error()
}
