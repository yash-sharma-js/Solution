package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"prime-trade/task-service/internal/models"
	"prime-trade/task-service/internal/repositories"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Common errors
var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrUnauthorized  = errors.New("unauthorized to perform this action")
	ErrInvalidStatus = errors.New("invalid task status")
)

// TaskService defines the interface for task operations
type TaskService interface {
	Create(ctx context.Context, req *models.CreateTaskRequest, userID uuid.UUID) (*models.TaskResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, user *models.AuthUser) (*models.TaskResponse, error)
	List(ctx context.Context, req *models.ListTasksRequest, user *models.AuthUser) (*models.TaskListResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateTaskRequest, user *models.AuthUser) (*models.TaskResponse, error)
	Delete(ctx context.Context, id uuid.UUID, user *models.AuthUser) error
}

type taskService struct {
	taskRepo    repositories.TaskRepository
	redisClient *redis.Client
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repositories.TaskRepository, redisClient *redis.Client) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		redisClient: redisClient,
	}
}

// Create creates a new task
func (s *taskService) Create(ctx context.Context, req *models.CreateTaskRequest, userID uuid.UUID) (*models.TaskResponse, error) {
	// Validate status if provided
	if req.Status != "" && !models.IsValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	status := req.Status
	if status == "" {
		status = models.StatusPending
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		UserID:      userID,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Invalidate cache
	s.invalidateUserTaskCache(ctx, userID)

	response := task.ToResponse()
	return &response, nil
}

// GetByID gets a task by ID
func (s *taskService) GetByID(ctx context.Context, id uuid.UUID, user *models.AuthUser) (*models.TaskResponse, error) {
	// Try to get from cache
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("task:%s", id.String())
		cached, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var task models.Task
			if json.Unmarshal([]byte(cached), &task) == nil {
				// Check authorization
				if user.Role != "ADMIN" && task.UserID != user.UserID {
					return nil, ErrUnauthorized
				}
				response := task.ToResponse()
				return &response, nil
			}
		}
	}

	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Check authorization: USER can only access own tasks, ADMIN can access all
	if user.Role != "ADMIN" && task.UserID != user.UserID {
		return nil, ErrUnauthorized
	}

	// Cache the task
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("task:%s", id.String())
		data, _ := json.Marshal(task)
		s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	response := task.ToResponse()
	return &response, nil
}

// List lists tasks based on user role
func (s *taskService) List(ctx context.Context, req *models.ListTasksRequest, user *models.AuthUser) (*models.TaskListResponse, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// Validate status if provided
	if req.Status != "" && !models.IsValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	var tasks []models.Task
	var totalCount int64
	var err error

	// ADMIN can see all tasks, USER can only see own tasks
	if user.Role == "ADMIN" {
		tasks, totalCount, err = s.taskRepo.FindAll(ctx, req.Page, req.Limit, req.Status)
	} else {
		tasks, totalCount, err = s.taskRepo.FindByUserID(ctx, user.UserID, req.Page, req.Limit, req.Status)
	}

	if err != nil {
		return nil, err
	}

	// Convert to response
	taskResponses := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = task.ToResponse()
	}

	totalPages := int(totalCount) / req.Limit
	if int(totalCount)%req.Limit > 0 {
		totalPages++
	}

	return &models.TaskListResponse{
		Tasks:      taskResponses,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Update updates a task
func (s *taskService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateTaskRequest, user *models.AuthUser) (*models.TaskResponse, error) {
	// Validate status if provided
	if req.Status != "" && !models.IsValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Check authorization: USER can only update own tasks, ADMIN can update all
	if user.Role != "ADMIN" && task.UserID != user.UserID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Invalidate cache
	s.invalidateTaskCache(ctx, id)
	s.invalidateUserTaskCache(ctx, task.UserID)

	response := task.ToResponse()
	return &response, nil
}

// Delete deletes a task
func (s *taskService) Delete(ctx context.Context, id uuid.UUID, user *models.AuthUser) error {
	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Check authorization: USER can only delete own tasks, ADMIN can delete all
	if user.Role == "ADMIN" {
		err = s.taskRepo.Delete(ctx, id)
	} else {
		if task.UserID != user.UserID {
			return ErrUnauthorized
		}
		err = s.taskRepo.DeleteByIDAndUserID(ctx, id, user.UserID)
	}

	if err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateTaskCache(ctx, id)
	s.invalidateUserTaskCache(ctx, task.UserID)

	return nil
}

func (s *taskService) invalidateTaskCache(ctx context.Context, id uuid.UUID) {
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("task:%s", id.String())
		s.redisClient.Del(ctx, cacheKey)
	}
}

func (s *taskService) invalidateUserTaskCache(ctx context.Context, userID uuid.UUID) {
	if s.redisClient != nil {
		pattern := fmt.Sprintf("tasks:user:%s:*", userID.String())
		keys, _ := s.redisClient.Keys(ctx, pattern).Result()
		if len(keys) > 0 {
			s.redisClient.Del(ctx, keys...)
		}
	}
}
