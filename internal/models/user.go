package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateUser struct {
	ID        string         `json:"id,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty"`
}

type User struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
