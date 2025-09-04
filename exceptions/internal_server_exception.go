package exceptions

import (
	"task-manager-app/exceptions/errors"
	"task-manager-app/utils"
	"net/http"
	"time"
)

func InternalServerException(message string) *errors.TaskManagerError {
	utils.Sugar.Error(message)
	return &errors.TaskManagerError{
		ErrorTimestamp: time.Now().UnixMilli(),
		Message:        message,
		ResponseCode:   http.StatusInternalServerError,
	}
}
