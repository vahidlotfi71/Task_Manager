package Models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

// define the possible statuses of task
const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_pending"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Status      TaskStatus     `json:"status" gorm:"not null;default:'pending'"`
	Assignee    string         `json:"assignee"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Task) TableName() string {
	return "tasks"
}
