package repositories

import (
	"context"
	"errors"

	"prime-trade/task-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int, status models.TaskStatus) ([]models.Task, int64, error)
	FindAll(ctx context.Context, page, limit int, status models.TaskStatus) ([]models.Task, int64, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByIDAndUserID(ctx context.Context, id, userID uuid.UUID) error
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create creates a new task in the database
func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// FindByID finds a task by ID
func (r *taskRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// FindByUserID finds tasks by user ID with pagination
func (r *taskRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int, status models.TaskStatus) ([]models.Task, int64, error) {
	var tasks []models.Task
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("user_id = ?", userID)

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, totalCount, nil
}

// FindAll finds all tasks with pagination (for admin)
func (r *taskRepository) FindAll(ctx context.Context, page, limit int, status models.TaskStatus) ([]models.Task, int64, error) {
	var tasks []models.Task
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Task{})

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, totalCount, nil
}

// Update updates a task in the database
func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete soft deletes a task from the database
func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByIDAndUserID deletes a task by ID and user ID
func (r *taskRepository) DeleteByIDAndUserID(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
