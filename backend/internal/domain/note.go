package domain

import "time"

type Note struct {
	ID        uint      `json:"id"         gorm:"primarykey"`
	Title     string    `json:"title"      gorm:"not null"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"    gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
