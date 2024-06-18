package model

import (
	"time"
)

type SpeedtestResults struct {
	ID string `json:"id" gorm:"primaryKey"` // Assuming ID is your primary key

	// Download
	DownloadSpeed    int `json:"download_speed" gorm:"column:download_speed"`         // average | kbps
	MaxDownloadSpeed int `json:"max_download_speed" gorm:"column:max_download_speed"` // kbps
	MinDownloadSpeed int `json:"min_download_speed" gorm:"column:min_download_speed"` // kbps
	TotalDownload    int `json:"total_download" gorm:"column:total_download"`         // kbps

	// Upload
	UploadSpeed    int `json:"upload_speed" gorm:"column:upload_speed"`         // average | kbps
	MaxUploadSpeed int `json:"max_upload_speed" gorm:"column:max_upload_speed"` // kbps
	MinUploadSpeed int `json:"min_upload_speed" gorm:"column:min_upload_speed"` // kbps
	TotalUpload    int `json:"total_upload" gorm:"column:total_upload"`         // kbps

	// Latency
	Latency         int `json:"latency" gorm:"column:latency"`                   // average | ms
	LoadedLatency   int `json:"loaded_latency" gorm:"column:loaded_latency"`     // ms
	UnloadedLatency int `json:"unloaded_latency" gorm:"column:unloaded_latency"` // ms
	DownloadLatency int `json:"download_latency" gorm:"column:download_latency"` // ms
	UploadLatency   int `json:"upload_latency" gorm:"column:upload_latency"`     // ms

	// Client
	ClientID         string `json:"client_id" gorm:"column:client_id"` // unique way to identify the client device
	ClientIP         string `json:"client_ip" gorm:"column:client_ip"` // Consider if this is necessary for storage
	ISP              string `json:"isp" gorm:"column:isp"`
	ISPCode          string `json:"isp_code" gorm:"column:isp_code"`
	ConnectionType   string `json:"connection_type" gorm:"column:connection_type"`     // "DSL," "Cable," "Fiber," or "Wireless."
	ConnectionDevice string `json:"connection_device" gorm:"column:connection_device"` // "5G Router," "Mobile," "Fiber," or "Wireless."
	TestPlatform     string `json:"test_platform" gorm:"column:test_platform"`

	// Location
	City           string  `json:"city" gorm:"column:city"`
	Longitude      float64 `json:"longitude" gorm:"column:longitude"` // Consider accuracy field for location
	Latitude       float64 `json:"latitude" gorm:"column:latitude"`
	CountryCode    string  `json:"country_code" gorm:"column:country_code"` // 3-letter country code
	CountryName    string  `json:"country_name" gorm:"column:country_name"`
	ServerLocation string  `json:"server_location" gorm:"column:server_location"`
	ServerName     string  `json:"server_name" gorm:"column:server_name"`
	LocationAccess bool    `json:"location_access" gorm:"column:location_access"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"` // time when record is created
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"` // time when record is updated
	TestTime  time.Time `json:"test_time" gorm:"column:test_time"`   // time when the internet test was taken

	// Optional fields for GORM
	// ServerID       string  `json:"server_id" gorm:"column:server_id"`
}

// Optional:
