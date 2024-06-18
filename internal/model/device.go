package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID               string         `json:"id" gorm:"type:uuid;primaryKey"`
	DeviceIdentifier string         `json:"device_identifier" gorm:"type:varchar(100);unique;not null"`
	UserID           string         `json:"user_id" gorm:"type:uuid"`
	OS               string         `json:"os" gorm:"type:varchar(50)"`
	DeviceType       string         `json:"device_type" gorm:"type:varchar(50)"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at"` // time when record is created
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
