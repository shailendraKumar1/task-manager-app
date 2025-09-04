package enums

type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusInProgress TaskStatus = "InProgress"
	StatusCompleted  TaskStatus = "Completed"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	}
	return false
}
