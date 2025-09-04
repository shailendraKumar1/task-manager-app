package taskManagerService

import (
	"fmt"
	"task-manager-app/constants"
	"task-manager-app/constants/enums"
	"task-manager-app/exceptions"
	"task-manager-app/exceptions/errors"
	"task-manager-app/models"
	"task-manager-app/repo"
	"task-manager-app/request"
	"task-manager-app/response"
	"task-manager-app/services/userManagerServices"
	"task-manager-app/utils"
)

type TaskService interface {
	CreateTask(req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError)
	UpdateTask(uuid string, req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError)
	GetTaskByUUID(uuid string) (*response.TaskResponse, *errors.TaskManagerError)
	DeleteTask(uuid string) *errors.TaskManagerError
	ListTasks(status string, userID string, priority string, page, pageSize int) (*response.TaskListResponse, *errors.TaskManagerError)
}

type taskService struct {
	repo        repo.TaskRepository
	userService userManagerServices.UserService
}

func NewTaskService(repository repo.TaskRepository) TaskService {
	return &taskService{
		repo:        repository,
		userService: userManagerServices.NewUserService(),
	}
}

func (s *taskService) CreateTask(req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError) {
	// Validate for create
	if err := s.validateReq(req); err != nil {
		return nil, err
	}

	// Validate user ID if provided
	if req.UserID != nil && *req.UserID != "" {
		isValid, userErr := s.userService.ValidateUser(*req.UserID)
		if userErr != nil {
			return nil, exceptions.InternalServerException(fmt.Sprintf("Failed to validate user: %v", userErr))
		}
		if !isValid {
			return nil, exceptions.NotFoundException(constants.ErrUserNotFound)
		}
	}

	// Convert request to model
	task := &models.Task{
		Title:       *req.Title,
		Description: utils.TaskManagerUtils.GetStringValue(req.Description),
		Status:      utils.TaskManagerUtils.GetStringValue(req.Status),
		Priority:    utils.TaskManagerUtils.GetStringValue(req.Priority),
		UserID:      req.UserID,
	}

	// Set default status if empty
	if task.Status == "" {
		task.Status = string(enums.StatusPending)
	}

	// Set default priority if empty
	if task.Priority == "" {
		task.Priority = string(enums.PriorityMedium)
	}

	if taskErr := s.repo.Create(task); taskErr != nil {
		return nil, taskErr
	}

	return s.taskDTOMapperResponse(task), nil
}

func (s *taskService) UpdateTask(uuid string, req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError) {
	// Validate for update
	if err := s.validateReq(req); err != nil {
		return nil, err
	}

	// Get existing task with row-level locking for update
	task, taskErr := s.repo.GetByUUIDForUpdate(uuid)
	if taskErr != nil {
		return nil, taskErr
	}
	if task == nil {
		return nil, exceptions.NotFoundException(constants.ErrTaskNotFound)
	}

	// Apply updates from request
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.UserID != nil && *req.UserID != "" {
		// Validate user ID before updating
		isValid, userErr := s.userService.ValidateUser(*req.UserID)
		if userErr != nil {
			return nil, exceptions.InternalServerException(fmt.Sprintf("Failed to validate user: %v", userErr))
		}
		if !isValid {
			return nil, exceptions.NotFoundException(constants.ErrUserNotFound)
		}
		task.UserID = req.UserID
	}

	// Update in database
	if taskErr := s.repo.Update(task); taskErr != nil {
		return nil, taskErr
	}

	return s.taskDTOMapperResponse(task), nil
}

// GetTaskByUUID fetches a task
func (s *taskService) GetTaskByUUID(uuid string) (*response.TaskResponse, *errors.TaskManagerError) {
	task, taskErr := s.repo.GetByUUID(uuid)
	if taskErr != nil {
		return nil, taskErr
	}
	if task == nil {
		return nil, exceptions.NotFoundException(constants.ErrTaskNotFound)
	}
	return s.taskDTOMapperResponse(task), nil
}

// DeleteTask removes a task
func (s *taskService) DeleteTask(uuid string) *errors.TaskManagerError {
	// Check if task exists with row-level locking before deletion
	task, taskErr := s.repo.GetByUUIDForUpdate(uuid)
	if taskErr != nil {
		return taskErr
	}
	if task == nil {
		return exceptions.NotFoundException(constants.ErrTaskNotFound)
	}

	return s.repo.Delete(uuid)
}

// ListTasks with pagination and filtering
func (s *taskService) ListTasks(status string, userID string, priority string, page, pageSize int) (*response.TaskListResponse, *errors.TaskManagerError) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	
	// Validate user ID if provided for filtering
	if userID != "" {
		isValid, userErr := s.userService.ValidateUser(userID)
		if userErr != nil {
			return nil, exceptions.InternalServerException(fmt.Sprintf("Failed to validate user: %v", userErr))
		}
		if !isValid {
			return nil, exceptions.NotFoundException(constants.ErrUserNotFound)
		}
	}
	
	// Validate priority if provided for filtering
	if priority != "" {
		if !enums.TaskPriority(priority).IsValid() {
			return nil, exceptions.NewBadRequestException(constants.ErrInvalidTaskPriority)
		}
	}
	
	tasks, taskErr := s.repo.List(status, userID, priority, pageSize, offset)
	if taskErr != nil {
		return nil, taskErr
	}

	// Convert to response objects
	taskResponses := make([]response.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *s.taskDTOMapperResponse(&task)
	}

	return &response.TaskListResponse{
		Tasks:    taskResponses,
		Page:     page,
		PageSize: pageSize,
		Count:    len(tasks),
	}, nil
}

func (s *taskService) validateReq(req *request.ReqCreateOrUpdateTasks) *errors.TaskManagerError {
	if req.Title == nil || *req.Title == "" {
		return exceptions.NewBadRequestException(constants.ErrTaskNotFound)
	}

	if req.Status != nil && !enums.TaskStatus(*req.Status).IsValid() {
		return exceptions.NewBadRequestException(constants.ErrInvalidTaskStatus)
	}

	if req.Priority != nil && !enums.TaskPriority(*req.Priority).IsValid() {
		return exceptions.NewBadRequestException(constants.ErrInvalidTaskPriority)
	}

	// Validate user ID if provided
	if req.UserID != nil {
		if err := s.validateUserID(*req.UserID); err != nil {
			return exceptions.NewBadRequestException(err.Error())
		}
	}

	return nil
}

func (s *taskService) taskDTOMapperResponse(task *models.Task) *response.TaskResponse {
	return &response.TaskResponse{
		UUID:        task.UUID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

// validateUserID validates if the user ID exists using user service
func (s *taskService) validateUserID(userID string) error {
	exists, err := s.userService.ValidateUser(userID)
	if err != nil {
		return fmt.Errorf("failed to validate user ID: %s", err.Error())
	}

	if !exists {
		return fmt.Errorf("user ID %s does not exist", userID)
	}

	return nil
}
