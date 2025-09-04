package enums

// TaskPriority represents the priority levels for tasks
type TaskPriority string

const (
	PriorityLow    TaskPriority = "Low"
	PriorityMedium TaskPriority = "Medium"
	PriorityHigh   TaskPriority = "High"
	PriorityUrgent TaskPriority = "Urgent"
)

// IsValid checks if the task priority is valid
func (p TaskPriority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

// String returns the string representation of TaskPriority
func (p TaskPriority) String() string {
	return string(p)
}
