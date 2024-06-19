package repositories

import (
	"context"
	"fmt"

	"github.com/checkspeed/sc-backend/internal/models"
	// "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	SpeedtestResultCollection = "speed_test_results"
	usersCollection           = "users"
)

type GetSpeedTestResultsFilter struct {
	CountryCode string `json:"country_code"` // 3 letter country code
}

type Datastore interface {
	CreateSpeedtestResult(ctx context.Context, speedTestResult *models.SpeedTestResult) error
	GetSpeedTestResults(ctx context.Context, filters GetSpeedTestResultsFilter) ([]models.SpeedTestResult, error)

	CloseConn(ctx context.Context) error
}

type store struct {
	db *gorm.DB
	// update stor to use gorm instead of sqlx
}

// sample from gorm documentation
// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

func NewStore(dbUrl string) (store, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return store{}, err
	}

	fmt.Println("database connected")

	return store{db}, nil
}

func (s *store) RunAutoMigrate() error {
	err := s.db.AutoMigrate(&models.User{}, &models.TestServer{}, &models.Device{}, &models.SpeedTestResult{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate models: %v", err)
	}
	fmt.Println("database migrated")
	return nil
}

func (s *store) RunDropTable() error {
	err := s.db.Migrator().DropTable(&models.User{}, &models.TestServer{}, &models.Device{}, &models.SpeedTestResult{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate models: %v", err)
	}
	fmt.Println("tables dropped")
	return nil
}

func (s store) CloseConn(ctx context.Context) error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (s store) CreateSpeedtestResult(ctx context.Context, speedTestResult *models.SpeedTestResult) error {
	result := s.db.WithContext(ctx).Create(&speedTestResult)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s store) GetSpeedTestResults(ctx context.Context, filters GetSpeedTestResultsFilter) ([]models.SpeedTestResult, error) {
	var speedTestResult []models.SpeedTestResult

	query := s.db.WithContext(ctx)
	if filters.CountryCode != "" {
		query.Where("name = ?", filters.CountryCode)
	}
	result := query.Find(&speedTestResult)

	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("Rows affected:", result.RowsAffected)

	return speedTestResult, nil

}

// func (s store) CreateSpeedtestResult(ctx context.Context, speedTestResult *models.SpeedtestResult) error {
// 	sqlStatement := `
// 		INSERT INTO speed_test_results (
// 			id,
// 			download_speed,
// 			max_download_speed,
// 			min_download_speed,
//         	total_download,

//         	upload_speed,
//         	max_upload_speed,
//         	min_upload_speed,
//         	total_upload,

//         	latency,
//         	loaded_latency,
//         	unloaded_latency,
//         	download_latency,
//         	upload_latency,

//         	client_id,
// 			client_ip,
// 			isp,
// 			isp_code,
// 			connection_type,
//         	connection_device,
// 			test_platform,

//         	city,
// 			longitude,
//         	latitude,
// 			country_code,
// 			country_name,
// 			server_location,
// 			server_name,
//         	location_access
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`

// 	_, err := s.db.Exec(sqlStatement,
// 		speedTestResult.ID,

// 		speedTestResult.DownloadSpeed,
// 		speedTestResult.MaxDownloadSpeed,
// 		speedTestResult.MinDownloadSpeed,
// 		speedTestResult.TotalDownload,

// 		speedTestResult.UploadSpeed,
// 		speedTestResult.MaxUploadSpeed,
// 		speedTestResult.MinUploadSpeed,
// 		speedTestResult.TotalUpload,

// 		speedTestResult.Latency,
// 		speedTestResult.LoadedLatency,
// 		speedTestResult.UnloadedLatency,
// 		speedTestResult.DownloadLatency,
// 		speedTestResult.UploadLatency,

// 		speedTestResult.ISP,
// 		speedTestResult.ISPCode,
// 		speedTestResult.ConnectionType,
// 		speedTestResult.ConnectionDevice,
// 		speedTestResult.TestPlatform,

// 		speedTestResult.City,
// 		speedTestResult.Longitude,
// 		speedTestResult.Latitude,
// 		speedTestResult.CountryCode,
// 		speedTestResult.CountryName,

// 		speedTestResult.LocationAccess,
// 	)

// 	return err
// }

// func (s store) GetSpeedtestResult(ctx context.Context, filters GetSpeedtestResultFilter) ([]models.SpeedtestResult, error) {
// 	var results = []models.SpeedtestResult{}
// 	sqlQuery := `
// 		SELECT *
// 		FROM speed_test_results
// 		`

// 	var args []interface{}
// 	if filters.CountryCode != "" {
// 		sqlQuery += " WHERE country_code = $"
// 		args = append(args, filters.CountryCode)
// 	}

// 	rows, err := s.db.Queryx(sqlQuery, args...)
// 	if err != nil {
// 		return []models.SpeedtestResult{}, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var sp models.SpeedtestResult
// 		err = rows.StructScan(&sp)
// 		if err != nil {
// 			return []models.SpeedtestResult{}, err
// 		}
// 		results = append(results, sp)
// 	}

// 	return results, nil
// }
