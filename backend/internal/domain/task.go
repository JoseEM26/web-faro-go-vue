package domain

import "time"

type Task struct {
	ID          uint      `json:"id"          gorm:"primarykey"`
	Title       string    `json:"title"       gorm:"not null"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"   gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
