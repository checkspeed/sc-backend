package models

type ApiResp struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type NetworkData struct {
	Isp           string `json:"isp,omitempty"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	CountryCode   string `json:"country_code2,omitempty"` // 3 letter country code
	CountryName   string `json:"country_name,omitempty"`
	ConitnentName string `json:"continent_name,omitempty"`
	ContinentCode string `json:"continent_code,omitempty"`
	State         string `json:"state_prov,omitempty"`
}
