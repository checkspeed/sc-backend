package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateDevice struct {
	DeviceIP         string `json:"device_ip,omitempty"` // Used to generate hash, not saved in db
	UserID           string `json:"user_id,omitempty"`
	OS               string `json:"os,omitempty"`
	DeviceType       string `json:"device_type,omitempty"`
	Manufacturer     string `json:"manufacturer,omitempty"`
	Model            string `json:"model,omitempty"`
	ScreenResolution string `json:"screen_resolution,omitempty"`
}

type Device struct {
	ID string `json:"id"`

	UserID           *string `json:"user_id"`
	DeviceID         *string `json:"device_id"` // self referenced field

	Identifier       string  `json:"identifier"`
	OS               string  `json:"os"`
	Manufacturer     string  `json:"manufacturer"`
	Model            string  `json:"model"`
	ScreenResolution string  `json:"screen_resolution"`
	DeviceType       string  `json:"device_type"` // Mobile, Desktop
	IsPlatformDevice bool    `json:"is_platform_device"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
