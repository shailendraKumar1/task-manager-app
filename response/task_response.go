package response

import "time"

type TaskResponse struct {
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	UserID      *string   `json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks    []TaskResponse `json:"tasks"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Count    int            `json:"count"`
}
