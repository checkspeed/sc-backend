package models

import (
	"time"
)

type SpeedTestResult struct {
	ID string `json:"id" gorm:"type:uuid;primaryKey"`

	// Download
	DownloadSpeed    int `json:"download_speed" gorm:"not null"` // average | kbps
	MaxDownloadSpeed int `json:"max_download_speed"`             // kbps
	MinDownloadSpeed int `json:"min_download_speed"`             // kbps
	TotalDownload    int `json:"total_download"`                 // kbps

	// Upload
	UploadSpeed    int `json:"upload_speed" gorm:"not null"` // average | kbps
	MaxUploadSpeed int `json:"max_upload_speed"`             // kbps
	MinUploadSpeed int `json:"min_upload_speed"`             // kbps
	TotalUpload    int `json:"total_upload"`                 // kbps

	// Latency
	Latency         int `json:"latency" gorm:"not null"` // average | ms
	LoadedLatency   int `json:"loaded_latency"`          // ms
	UnloadedLatency int `json:"unloaded_latency"`        // ms
	DownloadLatency int `json:"download_latency"`        // ms
	UploadLatency   int `json:"upload_latency"`          // ms

	// Device and Server
	DeviceID         string `json:"device_id" gorm:"type:uuid;default:null"`
	ISP              string `json:"isp" gorm:"type:varchar(50)"`
	ISPCode          string `json:"isp_code" gorm:"type:varchar(15)"`
	ConnectionType   string `json:"connection_type" gorm:"type:varchar(50)"`
	ConnectionDevice string `json:"connection_device" gorm:"type:varchar(50)"`
	TestPlatform     string `json:"test_platform" gorm:"type:varchar(50)"`
	ServerID         string `json:"server_id" gorm:"type:uuid;default:null"`

	// Location
	City           string  `json:"city" gorm:"type:varchar(50)"`
	State          string  `json:"state" gorm:"type:varchar(50)"`
	CountryCode    string  `json:"country_code" gorm:"type:varchar(5)"`
	CountryName    string  `json:"country_name" gorm:"type:varchar(50)"`
	ContinentCode  string  `json:"continent_code" gorm:"type:varchar(5)"`
	ContinentName  string  `json:"continent_name" gorm:"type:varchar(50)"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	LocationAccess bool    `json:"location_access" gorm:"default:false"`

	// Timestamps
	TestTime  time.Time `json:"test_time" gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;not null"`

	// Belongs to Relationship
	Device Device     `gorm:"foreignKey:DeviceID"`
	Server TestServer `gorm:"foreignKey:ServerID"`
}
