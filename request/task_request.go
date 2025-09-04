package request

type ReqCreateOrUpdateTasks struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
}

