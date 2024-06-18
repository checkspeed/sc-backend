package repositories

import (
	"context"

	"github.com/checkspeed/sc-backend/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	SpeedtestResultsCollection = "speed_test_results"
	usersCollection            = "users"
)

type GetSpeedtestResultsFilter struct {
	CountryCode string `json:"country_code"` // 3 letter country code
}

type Datastore interface {
	CreateSpeedtestResults(ctx context.Context, speedTestResult *model.SpeedtestResults) error
	GetSpeedtestResults(ctx context.Context, filters GetSpeedtestResultsFilter) ([]model.SpeedtestResults, error)

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

func (s store) CreateSpeedtestResults(ctx context.Context, speedTestResult *model.SpeedtestResults) error {
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

        	client_id,
			client_ip,
			isp,
			isp_code,
			connection_type,
        	connection_device,
			test_platform,

        	city,
			longitude,
        	latitude,
			country_code,
			country_name,
			server_location,
			server_name,
        	location_access
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`

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
		speedTestResult.UploadLatency,

		speedTestResult.ClientID,
		speedTestResult.ClientIP,
		speedTestResult.ISP,
		speedTestResult.ISPCode,
		speedTestResult.ConnectionType,
		speedTestResult.ConnectionDevice,
		speedTestResult.TestPlatform,

		speedTestResult.City,
		speedTestResult.Longitude,
		speedTestResult.Latitude,
		speedTestResult.CountryCode,
		speedTestResult.CountryName,
		speedTestResult.ServerLocation,
		speedTestResult.ServerName,
		speedTestResult.LocationAccess,
	)

	return err
}

func (s store) GetSpeedtestResults(ctx context.Context, filters GetSpeedtestResultsFilter) ([]model.SpeedtestResults, error) {
	var results = []model.SpeedtestResults{}
	sqlQuery := `
		SELECT *
		FROM speed_test_results
		`

	var args []interface{}
	if filters.CountryCode != "" {
		sqlQuery += " WHERE country_code = $"
		args = append(args, filters.CountryCode)
	}

	rows, err := s.db.Queryx(sqlQuery, args...)
	if err != nil {
		return []model.SpeedtestResults{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var sp model.SpeedtestResults
		err = rows.StructScan(&sp)
		if err != nil {
			return []model.SpeedtestResults{}, err
		}
		results = append(results, sp)
	}

	return results, nil
}

func (s store) CloseConn(ctx context.Context) error {

	return nil
}
