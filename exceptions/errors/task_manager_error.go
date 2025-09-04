package errors

type TaskManagerError struct {
	ErrorTimestamp int64  `json:"timestamp"`
	Message        string `json:"message"`
	ResponseCode   int    `json:"response_code"`
}
