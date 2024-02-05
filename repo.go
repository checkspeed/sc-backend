package main

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

const (
	SpeedtestResultsCollection = "speed_test_results"
	usersCollection            = "users"
)

type SpeedtestResults struct {
	ID string `json:"id"`

	// download
	DownloadSpeed    int `json:"download_speed" db:"download_speed"` // average | kbps
	MaxDownloadSPeed int `json:"max_download_speed" db:"max_download_speed"` // kbps
	MinDownloadSpeed int `json:"min_download_speed" db:"min_download_speed"` // kbps
	TotalDownload    int `json:"total_download" db:"total_download"`     // kbp

	// upload
	UploadSpeed    int `json:"upload_speed" db:"upload_speed"`     // average | kbps
	MaxUploadSpeed int `json:"max_upload_speed" db:"max_upload_speed"` // kbps
	MinUploadSpeed int `json:"min_upload_speed" db:"min_upload_speed"` // kbps
	TotalUpload    int `json:"total_upload" db:"total_upload"`     // kbps

	// latency
	Latency         int `json:"latency" db:"latency"`          // average | ms
	LoadedLatency   int `json:"loaded_latency" db:"loaded_latency"`   // ms
	UnloadedLatency int `json:"unloaded_latency" db:"unloaded_latency"` // ms
	DownloadLatency int `json:"download_latency" db:"download_latency"` // ms
	UploadLatency   int `json:"upload_latency" db:"upload_latency"`   // ms

	ConnectionType   string `json:"connection_type" db:"connection_type"`   // "DSL," "Cable," "Fiber," or "Wireless."
	ConnectionDevice string `json:"connection_device" db:"connection_device"` // "5G Router," "Mobile," "Fiber," or "Wireless."
	ISP              string `json:"isp" db:"isp"`
	ClientIP         string `json:"client_ip" db:"client_ip"` // Is this really needed for storage?
	ClientID         string `json:"client_id" db:"client_id"` // unique way to identiy the client device
	City             string `json:"city" db:"city"`
	ServerName       string `json:"server_name" db:"server_name"`
	TestServerID     string `json:"test_server_id" db:"test_server_id"`

	Long           float64 `json:"longitude" db:"longitude"`
	Lat            float64 `json:"latitude" db:"latitude"`
	LocationAccess bool    `json:"location_access" db:"location_access"`
	// there should be another field to indicate how accurate

	CreatedAt time.Time `json:"created_at" db:"created_at"` // time when record is created
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // time when record is created
	TestTime  time.Time `json:"test_time" db:"test_time"`  // time when the internet test was taken
}

type Datastore interface {
	CreateSpeedtestResults(ctx context.Context, speedTestResult *SpeedtestResults) error
	GetSpeedtestResults(ctx context.Context) ([]SpeedtestResults, error)

	CloseConn(ctx context.Context) error
}

type store struct {
	db *sqlx.DB
}

func NewStore(dbUrl string) (store, error) {
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return store{}, err
	}

	return store{db}, nil
}

func (s store) CreateSpeedtestResults(ctx context.Context, speedTestResult *SpeedtestResults) error {
	sqlStatement := `
		INSERT INTO speed_test_results (
			id,
			download_speed,
			max_download_speed,
			min_download_speed,
        	total_download,

        	upload_speed,
        	max_upload_speed,
        	min_upload_speed,
        	total_upload,

        	latency,
        	loaded_latency,
        	unloaded_latency,
        	download_latency,
        	upload_latency,

        	connection_type,
        	connection_device,
        	isp,
        	client_ip,
        	client_id,
        	city,
        	server_name,

        	longitude,
        	latitude,

        	location_access
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`

	_, err := s.db.Exec(sqlStatement,
		speedTestResult.ID,

		speedTestResult.DownloadSpeed,
		speedTestResult.MaxDownloadSPeed,
		speedTestResult.MinDownloadSpeed,
		speedTestResult.TotalDownload,

		speedTestResult.UploadSpeed,
		speedTestResult.MaxUploadSpeed,
		speedTestResult.MinUploadSpeed,
		speedTestResult.TotalUpload,

		speedTestResult.Latency,
		speedTestResult.LoadedLatency,
		speedTestResult.UnloadedLatency,
		speedTestResult.DownloadLatency,
		speedTestResult.UnloadedLatency,

		speedTestResult.ConnectionType,
		speedTestResult.ConnectionDevice,
		speedTestResult.ISP,
		speedTestResult.ClientIP,
		speedTestResult.ClientID,
		speedTestResult.City,
		speedTestResult.ServerName,

		speedTestResult.Long,
		speedTestResult.Lat,
		speedTestResult.LocationAccess,
	)

	return err
}

func (s store) GetSpeedtestResults(ctx context.Context) ([]SpeedtestResults, error) {
	var results = []SpeedtestResults{}
	sqlQuery := `
		SELECT *
		FROM speed_test_results
		`

		rows, err := s.db.Queryx(sqlQuery)
		if err != nil {
			return []SpeedtestResults{}, err
		}
		defer rows.Close()
		for rows.Next() {
			var sp SpeedtestResults
			err = rows.StructScan(&sp)
			if err != nil {
				return []SpeedtestResults{}, err
			}
			results = append(results, sp)
		}

	return results, nil
}

func (s store) CloseConn(ctx context.Context) error {

	return nil
}
