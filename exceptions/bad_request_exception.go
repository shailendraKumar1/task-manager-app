package exceptions

import (
	"task-manager-app/exceptions/errors"
	"net/http"

	"time"
)

func NewBadRequestException(message string) *errors.TaskManagerError {
	return &errors.TaskManagerError{
		ErrorTimestamp: time.Now().UnixMilli(),
		Message:        message,
		ResponseCode:   http.StatusBadRequest,
	}
}
