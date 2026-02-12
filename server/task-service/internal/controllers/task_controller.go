package controllers

import (
	"errors"
	"net/http"

	"prime-trade/task-service/internal/models"
	"prime-trade/task-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskController handles task endpoints
type TaskController struct {
	taskService services.TaskService
}

// NewTaskController creates a new task controller
func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateTaskRequest true "Task details"
// @Success 201 {object} models.APIResponse{data=models.TaskResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /tasks [post]
func (c *TaskController) Create(ctx *gin.Context) {
	// Get authenticated user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}
	authUser := user.(*models.AuthUser)

	var req models.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(formatValidationError(err)))
		return
	}

	response, err := c.taskService.Create(ctx.Request.Context(), &req, authUser.UserID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidStatus) {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to create task"))
		return
	}

	ctx.JSON(http.StatusCreated, models.NewSuccessResponse(response, "Task created successfully"))
}

// List godoc
// @Summary List tasks
// @Description List tasks for the authenticated user (or all tasks for admin)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(PENDING, IN_PROGRESS, DONE)
// @Success 200 {object} models.APIResponse{data=models.TaskListResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /tasks [get]
func (c *TaskController) List(ctx *gin.Context) {
	// Get authenticated user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}
	authUser := user.(*models.AuthUser)

	var req models.ListTasksRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(formatValidationError(err)))
		return
	}

	response, err := c.taskService.List(ctx.Request.Context(), &req, authUser)
	if err != nil {
		if errors.Is(err, services.ErrInvalidStatus) {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to list tasks"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(response, "Tasks retrieved successfully"))
}

// GetByID godoc
// @Summary Get task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} models.APIResponse{data=models.TaskResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /tasks/{id} [get]
func (c *TaskController) GetByID(ctx *gin.Context) {
	// Get authenticated user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}
	authUser := user.(*models.AuthUser)

	// Parse task ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid task ID"))
		return
	}

	response, err := c.taskService.GetByID(ctx.Request.Context(), id, authUser)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(err.Error()))
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get task"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(response, "Task retrieved successfully"))
}

// Update godoc
// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body models.UpdateTaskRequest true "Updated task details"
// @Success 200 {object} models.APIResponse{data=models.TaskResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /tasks/{id} [put]
func (c *TaskController) Update(ctx *gin.Context) {
	// Get authenticated user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}
	authUser := user.(*models.AuthUser)

	// Parse task ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid task ID"))
		return
	}

	var req models.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(formatValidationError(err)))
		return
	}

	response, err := c.taskService.Update(ctx.Request.Context(), id, &req, authUser)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(err.Error()))
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(err.Error()))
			return
		}
		if errors.Is(err, services.ErrInvalidStatus) {
			ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to update task"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(response, "Task updated successfully"))
}

// Delete godoc
// @Summary Delete a task
// @Description Delete an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /tasks/{id} [delete]
func (c *TaskController) Delete(ctx *gin.Context) {
	// Get authenticated user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}
	authUser := user.(*models.AuthUser)

	// Parse task ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid task ID"))
		return
	}

	err = c.taskService.Delete(ctx.Request.Context(), id, authUser)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, models.NewErrorResponse(err.Error()))
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, models.NewErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to delete task"))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccessResponse(nil, "Task deleted successfully"))
}

// formatValidationError formats validation errors into readable messages
func formatValidationError(err error) string {
	return "Validation failed: " + err.Error()
}
