package models

type ApiStatus string

const (
	StatusSuccess ApiStatus = "success"
	StatusFail    ApiStatus = "fail"
	StatusError   ApiStatus = "error"
)

type ApiResp struct {
	Status  ApiStatus `json:"status"`            // "success", "fail", or "error"
	Message string    `json:"message,omitempty"` // Human-readable message
	Data    any       `json:"data,omitempty"`    // Response payload
	Code    string    `json:"code,omitempty"`    // Machine-readable code
	Error   string    `json:"error,omitempty"`   // Human-readable message

}

type NetworkData struct {
	Isp           string `json:"isp,omitempty"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	CountryCode   string `json:"country_code2,omitempty"` // 2 letter country code
	CountryName   string `json:"country_name,omitempty"`
	ConitnentName string `json:"continent_name,omitempty"`
	ContinentCode string `json:"continent_code,omitempty"`
	State         string `json:"state_prov,omitempty"`
}

type GeoLocationData struct {
	OrganizationName string `json:"organization_name,omitempty"`
	Organization     string `json:"organization,omitempty"`
	Longitude        string `json:"longitude"`
	Latitude         string `json:"latitude"`
	CountryCode      string `json:"country_code,omitempty"` // 2 letter country code
	Country          string `json:"country,omitempty"`
	ConitnentName    string `json:"continent_name,omitempty"`
	ContinentCode    string `json:"continent_code,omitempty"`
	Region           string `json:"region,omitempty"`
	City             string `json:"city,omitempty"`
	Timezone         string `json:"timezone,omitempty"`
}
