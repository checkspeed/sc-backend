package models

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID         string         `json:"id" gorm:"type:uuid;primaryKey"`
	Identifier string         `json:"identifier" gorm:"unique;not null"`
	UserID     string         `json:"user_id"`
	OS         string         `json:"os"`
	DeviceType string         `json:"device_type"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Belongs to Relationship
	User User `gorm:"foreignKey:UserID"`
}
