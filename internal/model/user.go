package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"` // UUID primary key
	Username  string         `json:"username" gorm:"type:varchar(50);unique;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"` // time when record is created
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // Soft delete field
}
