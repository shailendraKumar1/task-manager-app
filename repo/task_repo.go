package repo

import (
	"sync"
	"task-manager-app/constants"
	"task-manager-app/exceptions"
	"task-manager-app/exceptions/errors"
	"task-manager-app/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) *errors.TaskManagerError
	GetByUUID(uuid string) (*models.Task, *errors.TaskManagerError)
	GetByUUIDForUpdate(uuid string) (*models.Task, *errors.TaskManagerError)
	Update(task *models.Task) *errors.TaskManagerError
	Delete(uuid string) *errors.TaskManagerError
	List(status string, userID string, priority string, limit, offset int) ([]models.Task, *errors.TaskManagerError)
}

type taskRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) *errors.TaskManagerError {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.db.Create(task).Error; err != nil {
		return exceptions.InternalServerException(constants.ErrFailedToCreateTask + ": " + err.Error())
	}
	return nil
}

// GetByUUID finds a task by its UUID
func (r *taskRepository) GetByUUID(uuid string) (*models.Task, *errors.TaskManagerError) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var task models.Task
	result := r.db.Where("uuid = ?", uuid).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exceptions.InternalServerException(constants.ErrFailedToGetTask + ": " + result.Error.Error())
	}
	return &task, nil
}

// GetByUUIDForUpdate finds a task by its UUID with row-level locking for updates
func (r *taskRepository) GetByUUIDForUpdate(uuid string) (*models.Task, *errors.TaskManagerError) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var task models.Task
	result := r.db.Set("gorm:query_option", "FOR UPDATE").Where("uuid = ?", uuid).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exceptions.InternalServerException(constants.ErrFailedToGetTask + ": " + result.Error.Error())
	}
	return &task, nil
}

// Update modifies an existing task
func (r *taskRepository) Update(task *models.Task) *errors.TaskManagerError {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.db.Save(task).Error; err != nil {
		return exceptions.InternalServerException(constants.ErrFailedToUpdateTask + ": " + err.Error())
	}
	return nil
}

// Delete removes a task by UUID
func (r *taskRepository) Delete(uuid string) *errors.TaskManagerError {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.db.Where("uuid = ?", uuid).Delete(&models.Task{}).Error; err != nil {
		return exceptions.InternalServerException(constants.ErrFailedToDeleteTask + ": " + err.Error())
	}
	return nil
}

// List fetches tasks with optional status, user_id, and priority filters + pagination
func (r *taskRepository) List(status string, userID string, priority string, limit, offset int) ([]models.Task, *errors.TaskManagerError) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tasks []models.Task
	query := r.db.Model(&models.Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, exceptions.InternalServerException(constants.ErrFailedToListTasks + ": " + err.Error())
	}
	return tasks, nil
}
