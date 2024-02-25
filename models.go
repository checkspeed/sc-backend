package main

type CreateSpeedtestResultsAPIRequest struct {

	// download
	DownloadSpeed    string `json:"download_speed"` // kbps | required
	MaxDownloadSPeed string `json:"max_download_speed"` // kbps
	MinDownloadSpeed string `json:"min_download_speed"` // kbps
	TotalDownload    string `json:"total_download"`     // kb

	// upload
	UploadSpeed    string `json:"upload_speed"`     // average | kbps
	MaxUploadSpeed string `json:"max_upload_speed"` // kbps
	MinUploadSpeed string `json:"min_upload_speed"` // kbps
	TotalUpload    string `json:"total_upload"`     // kbps

	// latency
	Latency         string `json:"latency"`          // average | ms
	LoadedLatency   string `json:"loaded_latency"`   // ms
	UnloadedLatency string `json:"unloaded_latency"` // ms
	DownloadLatency string `json:"download_latency"` // ms
	UploadLatency   string `json:"upload_latency"`   // ms

	ConnectionType   string `json:"connection_type" db:"connection_type"`   // "DSL," "Cable," "Fiber," or "Wireless."
	ConnectionDevice string `json:"connection_device" db:"connection_device"` // "5G Router," "Mobile," "Fiber," or "Wireless."
	ISP              string `json:"isp" db:"isp"`
	ClientIP         string `json:"client_ip" db:"client_ip"` // Is this really needed for storage?
	ClientID         string `json:"client_id" db:"client_id"` // unique way to identiy the client device
	City             string `json:"city" db:"city"` // ?
	ServerName       string `json:"server_name" db:"server_name"`
	TestServerID     string `json:"test_server_id" db:"test_server_id"`
	ServerLocation   string `json:"server_location" db:"server_location"`
	TestPlatform     string `json:"test_platform" db:"test_platform"`

	Longitude           string `json:"longitude" db:"longitude"`
	Latitude            string `json:"latitude" db:"latitude"`
	LocationAccess bool    `json:"location_access" db:"location_access"`

	TestTime  string `json:"test_time" db:"test_time"`  // time when the internet test was taken
}