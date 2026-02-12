package models

import "github.com/google/uuid"

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=5000"`
	Status      TaskStatus `json:"status,omitempty"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       string     `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description string     `json:"description,omitempty" binding:"omitempty,max=5000"`
	Status      TaskStatus `json:"status,omitempty"`
}

// ListTasksRequest represents the query parameters for listing tasks
type ListTasksRequest struct {
	Page   int        `form:"page,default=1" binding:"min=1"`
	Limit  int        `form:"limit,default=10" binding:"min=1,max=100"`
	Status TaskStatus `form:"status,omitempty"`
}

// TaskListResponse represents a paginated list of tasks
type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// AuthUser represents the authenticated user from JWT validation
type AuthUser struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
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
