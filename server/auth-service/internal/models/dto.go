package models

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     Role   `json:"role,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ValidateResponse represents the token validation response
type ValidateResponse struct {
	Valid  bool         `json:"valid"`
	User   UserResponse `json:"user,omitempty"`
	UserID string       `json:"user_id,omitempty"`
	Role   string       `json:"role,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Data:    nil,
		Message: message,
	}
}
