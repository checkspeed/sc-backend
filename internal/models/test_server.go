package models

import (
	"time"

	"gorm.io/gorm"
)

type TestServer struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string         `json:"name"`
	City      string         `json:"city"`
	Country   string         `json:"country" gorm:"not null"`
	URL       string         `json:"url"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
