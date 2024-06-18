package model

import (
	"time"

	"gorm.io/gorm"
)

type TestServer struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100)"`
	City      string         `json:"city" gorm:"type:varchar(100)"`
	Country   string         `json:"country" gorm:"type:varchar(100);not null"`
	URL       string         `json:"url" gorm:"type:varchar(255)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
