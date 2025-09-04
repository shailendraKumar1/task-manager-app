package utils

import "strconv"

var (
	TaskManagerUtils = &taskManagerUtils{}
)

type taskManagerUtils struct{}

func (t *taskManagerUtils) ParseStringToInt(secret string) int {
	if secret == "" {
		return 10 // default value for database connections
	}
	v, err := strconv.Atoi(secret)
	if err != nil {
		Sugar.Errorw(err.Error())
		return 10
	}
	return v
}

func (t *taskManagerUtils) GetStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
