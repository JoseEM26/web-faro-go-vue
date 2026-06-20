package domain

import "time"

type Category struct {
	ID        uint      `json:"id"         gorm:"primarykey"`
	Name      string    `json:"name"       gorm:"uniqueIndex;not null"`
	Color     string    `json:"color"      gorm:"not null;default:'#007A6E'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
