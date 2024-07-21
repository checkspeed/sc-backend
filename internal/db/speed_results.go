package db

import (
	"context"

	"github.com/checkspeed/sc-backend/internal/models"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// const (
// 	SpeedtestResultCollection = "speed_test_results"
// 	usersCollection           = "users"
// )

type GetSpeedTestResultsFilter struct {
	CountryCode string `json:"country_code"` // 3 letter country code
}

type SpeedTestResults interface {
	Create(ctx context.Context, speedTestResult *models.SpeedTestResults) error
	Get(ctx context.Context, filters GetSpeedTestResultsFilter) ([]models.SpeedTestResults, error)
}

type speedTestResultsRepo struct {
	db *gorm.DB
}

func NewSpeedTestResultsRepo(store Store) (speedTestResultsRepo, error) {
	db := store.DB()

	return speedTestResultsRepo{db}, nil
}

func (s speedTestResultsRepo) Create(ctx context.Context, speedTestResult *models.SpeedTestResults) error {
	result := s.db.WithContext(ctx).Create(&speedTestResult)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s speedTestResultsRepo)  Get(ctx context.Context, filters GetSpeedTestResultsFilter) ([]models.SpeedTestResults, error) {
	var speedTestResult []models.SpeedTestResults

	query := s.db.WithContext(ctx)

	if filters.CountryCode != "" {
		query = query.Where("country_code = ?", filters.CountryCode)
	}

	result := query.Find(&speedTestResult)

	if result.Error != nil {
		return nil, result.Error
	}

	return speedTestResult, nil
}
