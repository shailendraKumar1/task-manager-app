package taskManagerService

import (
	"task-manager-app/constants"
	"task-manager-app/constants/enums"
	"task-manager-app/exceptions"
	"task-manager-app/exceptions/errors"
	"task-manager-app/models"
	"task-manager-app/repo"
	"task-manager-app/request"
	"task-manager-app/response"
	"task-manager-app/services/validationService"
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
	repo              repo.TaskRepository
	validationService validationService.ValidationService
}

func NewTaskService(repository repo.TaskRepository, validationSvc validationService.ValidationService) TaskService {
	return &taskService{
		repo:              repository,
		validationService: validationSvc,
	}
}

func (s *taskService) CreateTask(req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError) {
	// Validate request
	if err := s.validationService.ValidateCreateTaskRequest(req); err != nil {
		return nil, err
	}

	// Convert request to model
	task := &models.Task{
		Title:       *req.Title,
		Description: utils.TaskManagerUtils.GetStringValue(req.Description),
		Status:      utils.TaskManagerUtils.GetStringValue(req.Status),
		Priority:    utils.TaskManagerUtils.GetStringValue(req.Priority),
		UserID:      req.UserID,
	}

	if task.Status == "" {
		task.Status = string(enums.StatusPending)
	}
	if task.Priority == "" {
		task.Priority = string(enums.PriorityMedium)
	}

	// Save
	if taskErr := s.repo.Create(task); taskErr != nil {
		return nil, taskErr
	}
	return s.toResponse(task), nil
}

func (s *taskService) GetTaskByUUID(uuid string) (*response.TaskResponse, *errors.TaskManagerError) {
	task, taskErr := s.repo.GetByUUID(uuid)
	if taskErr != nil {
		return nil, taskErr
	}
	if task == nil {
		return nil, exceptions.NotFoundException(constants.ErrTaskNotFound)
	}
	return s.toResponse(task), nil
}

func (s *taskService) DeleteTask(uuid string) *errors.TaskManagerError {
	task, taskErr := s.repo.GetByUUIDForUpdate(uuid)
	if taskErr != nil {
		return taskErr
	}
	if task == nil {
		return exceptions.NotFoundException(constants.ErrTaskNotFound)
	}
	return s.repo.Delete(uuid)
}

func (s *taskService) ListTasks(status string, userID string, priority string, page, pageSize int) (*response.TaskListResponse, *errors.TaskManagerError) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	// Validate filters
	if userID != "" {
		if err := s.validationService.ValidateUserID(userID); err != nil {
			return nil, err
		}
	}
	if priority != "" {
		if err := s.validationService.ValidateTaskPriority(priority); err != nil {
			return nil, err
		}
	}

	tasks, taskErr := s.repo.List(status, userID, priority, pageSize, offset)
	if taskErr != nil {
		return nil, taskErr
	}

	taskResponses := make([]response.TaskResponse, len(tasks))
	for i, t := range tasks {
		taskResponses[i] = *s.toResponse(&t)
	}

	return &response.TaskListResponse{
		Tasks:    taskResponses,
		Page:     page,
		PageSize: pageSize,
		Count:    len(tasks),
	}, nil
}

func (s *taskService) UpdateTask(uuid string, req *request.ReqCreateOrUpdateTasks) (*response.TaskResponse, *errors.TaskManagerError) {
	// Check if task exists
	task, taskErr := s.repo.GetByUUIDForUpdate(uuid)
	if taskErr != nil {
		return nil, taskErr
	}
	if task == nil {
		return nil, exceptions.NotFoundException(constants.ErrTaskNotFound)
	}

	// Apply updates in one place
	changed, err := s.applyUpdates(task, req)
	if err != nil {
		return nil, err
	}

	if !changed {
		return nil, exceptions.NewBadRequestException(constants.ErrNothingToChange)
	}

	if taskErr := s.repo.Update(task); taskErr != nil {
		return nil, taskErr
	}

	return s.toResponse(task), nil
}

func (s *taskService) applyUpdates(task *models.Task, req *request.ReqCreateOrUpdateTasks) (bool, *errors.TaskManagerError) {
	changed := false

	// Title
	if req.Title != nil && s.updateField(&task.Title, *req.Title) {
		changed = true
	}

	// Description
	if req.Description != nil && s.updateField(&task.Description, *req.Description) {
		changed = true
	}

	// Status
	if req.Status != nil {
		if err := s.validationService.ValidateTaskStatus(*req.Status); err != nil {
			return false, err
		}
		if s.updateField(&task.Status, *req.Status) {
			changed = true
		}
	}

	// Priority
	if req.Priority != nil {
		if err := s.validationService.ValidateTaskPriority(*req.Priority); err != nil {
			return false, err
		}
		if s.updateField(&task.Priority, *req.Priority) {
			changed = true
		}
	}

	// UserID
	if req.UserID != nil && *req.UserID != "" {
		if err := s.validationService.ValidateUserID(*req.UserID); err != nil {
			return false, err
		}
		if task.UserID == nil || *task.UserID != *req.UserID {
			task.UserID = req.UserID
			changed = true
		}
	}

	return changed, nil
}

func (s *taskService) updateField(field *string, newValue string) bool {
	if *field == newValue {
		return false
	}
	*field = newValue
	return true
}

func (s *taskService) toResponse(task *models.Task) *response.TaskResponse {
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
