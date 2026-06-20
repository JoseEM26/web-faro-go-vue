package domain

import "time"

type User struct {
	ID           uint      `json:"id"         gorm:"primarykey"`
	Email        string    `json:"email"      gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-"          gorm:"not null"` // nunca se serializa en JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
