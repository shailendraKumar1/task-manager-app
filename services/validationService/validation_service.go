package validationService

import (
	"fmt"
	"task-manager-app/constants"
	"task-manager-app/constants/enums"
	"task-manager-app/exceptions"
	"task-manager-app/exceptions/errors"
	"task-manager-app/repo"
	"task-manager-app/request"
	"task-manager-app/services/userManagerServices"
)

type ValidationService interface {
	ValidateCreateTaskRequest(req *request.ReqCreateOrUpdateTasks) *errors.TaskManagerError
	ValidateUpdateTaskRequest(req *request.ReqCreateOrUpdateTasks) *errors.TaskManagerError
	ValidateUserID(userID string) *errors.TaskManagerError
	ValidateTaskStatus(status string) *errors.TaskManagerError
	ValidateTaskPriority(priority string) *errors.TaskManagerError
	ValidateTaskTitle(title *string) *errors.TaskManagerError
	CheckTaskDuplicateByTitle(title, userID string) *errors.TaskManagerError
}

type validationService struct {
	userService userManagerServices.UserService
	taskRepo    repo.TaskRepository
}

func NewValidationService(userService userManagerServices.UserService, taskRepo repo.TaskRepository) ValidationService {
	return &validationService{
		userService: userService,
		taskRepo:    taskRepo,
	}
}

func (v *validationService) ValidateCreateTaskRequest(req *request.ReqCreateOrUpdateTasks) *errors.TaskManagerError {
	if err := v.validateCommonFields(req, true); err != nil {
		return err
	}

	if req.UserID != nil && *req.UserID != "" {
		if err := v.CheckTaskDuplicateByTitle(*req.Title, *req.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (v *validationService) ValidateUpdateTaskRequest(req *request.ReqCreateOrUpdateTasks) *errors.TaskManagerError {
	return v.validateCommonFields(req, false)
}

func (v *validationService) validateCommonFields(req *request.ReqCreateOrUpdateTasks, isCreate bool) *errors.TaskManagerError {
	if isCreate || req.Title != nil {
		if err := v.ValidateTaskTitle(req.Title); err != nil {
			return err
		}
	}

	if req.Status != nil {
		if err := v.ValidateTaskStatus(*req.Status); err != nil {
			return err
		}
	}

	if req.Priority != nil {
		if err := v.ValidateTaskPriority(*req.Priority); err != nil {
			return err
		}
	}

	if req.UserID != nil && *req.UserID != "" {
		if err := v.ValidateUserID(*req.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (v *validationService) ValidateTaskTitle(title *string) *errors.TaskManagerError {
	if title == nil || *title == "" {
		return exceptions.NewBadRequestException(constants.ErrInvalidTaskTitle)
	}
	return nil
}

func (v *validationService) ValidateTaskStatus(status string) *errors.TaskManagerError {
	if !enums.TaskStatus(status).IsValid() {
		return exceptions.NewBadRequestException(constants.ErrInvalidTaskStatus)
	}
	return nil
}

func (v *validationService) ValidateTaskPriority(priority string) *errors.TaskManagerError {
	if !enums.TaskPriority(priority).IsValid() {
		return exceptions.NewBadRequestException(constants.ErrInvalidTaskPriority)
	}
	return nil
}

func (v *validationService) ValidateUserID(userID string) *errors.TaskManagerError {
	valid, err := v.userService.ValidateUser(userID)
	if err != nil {
		return exceptions.InternalServerException(fmt.Sprintf("Failed to validate user: %v", err))
	}
	if !valid {
		return exceptions.NotFoundException(constants.ErrUserNotFound)
	}
	return nil
}

func (v *validationService) CheckTaskDuplicateByTitle(title, userID string) *errors.TaskManagerError {
	exists, err := v.taskRepo.ExistsByTitleAndUser(title, userID)
	if err != nil {
		return err
	}
	if exists {
		return exceptions.NewBadRequestException(constants.ErrTaskAlreadyExists)
	}
	return nil
}
