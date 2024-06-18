package model

import "time"

type SpeedtestResults struct {
	ID string `json:"id"`

	// download
	DownloadSpeed    int `json:"download_speed" db:"download_speed"`         // average | kbps
	MaxDownloadSPeed int `json:"max_download_speed" db:"max_download_speed"` // kbps
	MinDownloadSpeed int `json:"min_download_speed" db:"min_download_speed"` // kbps
	TotalDownload    int `json:"total_download" db:"total_download"`         // kbp

	// upload
	UploadSpeed    int `json:"upload_speed" db:"upload_speed"`         // average | kbps
	MaxUploadSpeed int `json:"max_upload_speed" db:"max_upload_speed"` // kbps
	MinUploadSpeed int `json:"min_upload_speed" db:"min_upload_speed"` // kbps
	TotalUpload    int `json:"total_upload" db:"total_upload"`         // kbps

	// latency
	Latency         int `json:"latency" db:"latency"`                   // average | ms
	LoadedLatency   int `json:"loaded_latency" db:"loaded_latency"`     // ms
	UnloadedLatency int `json:"unloaded_latency" db:"unloaded_latency"` // ms
	DownloadLatency int `json:"download_latency" db:"download_latency"` // ms
	UploadLatency   int `json:"upload_latency" db:"upload_latency"`     // ms

	// client
	ClientID         string `json:"client_id" db:"client_id"` // unique way to identiy the client device
	ClientIP         string `json:"client_ip" db:"client_ip"` // Is this really needed for storage?
	ISP              string `json:"isp" db:"isp"`
	ISPCode          string `json:"isp_code" db:"isp_code"`
	ConnectionType   string `json:"connection_type" db:"connection_type"`     // "DSL," "Cable," "Fiber," or "Wireless."
	ConnectionDevice string `json:"connection_device" db:"connection_device"` // "5G Router," "Mobile," "Fiber," or "Wireless."
	TestPlatform     string `json:"test_platform" db:"test_platform"`

	// location
	City           string  `json:"city" db:"city"`
	Longitude      float64 `json:"longitude" db:"longitude"` // note: consider using a field to depict how accurate
	Latitude       float64 `json:"latitude" db:"latitude"`
	CountryCode    string  `json:"country_code" db:"country_code"` // 3 letter country code
	CountryName    string  `json:"country_name" db:"country_name"`
	ServerLocation string  `json:"server_location" db:"server_location"`
	ServerName     string  `json:"server_name" db:"server_name"`
	// ServerID       string  `json:"server_id" db:"server_id"`
	LocationAccess bool `json:"location_access" db:"location_access"`
	// there should be another field to indicate how accurate

	CreatedAt time.Time `json:"created_at" db:"created_at"` // time when record is created
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // time when record is created
	TestTime  time.Time `json:"test_time" db:"test_time"`   // time when the internet test was taken
}
