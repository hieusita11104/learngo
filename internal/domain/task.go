package domain

import (
	"fmt"
)

// Task model
type Task struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
	CreatedBy uint   `json:"created_by"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy uint   `json:"updated_by"`
}

type ValidateTask struct {
	Title string
	Done  bool
}

type UserTask struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
}

func (t *Task) Validate() error {
	if len(t.Title) == 0 {
		return fmt.Errorf("title is required")
	}

	if len(t.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters long")
	}
	return nil
}
