package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        string    `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Status      string    `gorm:"type:varchar(20);not null" json:"status"`
	Priority    string    `gorm:"type:varchar(20);not null;default:'Medium'" json:"priority"`
	UserID      *string   `gorm:"index" json:"user_id,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Hook to generate UUID before creating a record
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == "" {
		t.UUID = uuid.New().String()
	}
	return
}
