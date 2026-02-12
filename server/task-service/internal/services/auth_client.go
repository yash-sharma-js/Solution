package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"prime-trade/task-service/internal/models"

	"github.com/google/uuid"
)

// AuthClient handles communication with the auth service
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

// ValidateResponse represents the response from auth service validation
type ValidateResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		User   struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

// NewAuthClient creates a new auth client
func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateToken validates a JWT token with the auth service
func (c *AuthClient) ValidateToken(token string) (*models.AuthUser, error) {
	url := fmt.Sprintf("%s/api/v1/auth/validate", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid or expired token")
	}

	var validateResp ValidateResponse
	if err := json.Unmarshal(body, &validateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !validateResp.Success || !validateResp.Data.Valid {
		return nil, errors.New("token validation failed")
	}

	userID, err := uuid.Parse(validateResp.Data.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return &models.AuthUser{
		UserID: userID,
		Email:  validateResp.Data.User.Email,
		Role:   validateResp.Data.Role,
	}, nil
}
