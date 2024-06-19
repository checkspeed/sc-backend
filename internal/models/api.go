package models

type ApiResp struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type NetworkData struct {
	Isp         string `json:"isp,omitempty"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	CountryCode string `json:"country_code3,omitempty"` // 3 letter country code
	CountryName string `json:"country_name,omitempty"`
}
