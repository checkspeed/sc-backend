package db

import (
	"context"

	// "github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/checkspeed/sc-backend/internal/models"
)

type Devices interface {
	GetOrCreate(ctx context.Context, device models.Device) (string, int64, error)
}

type devices struct {
	db *gorm.DB
}

func NewDevicesRepo(store Store) (*devices, error) {
	return &devices{
		db: store.DB(),
	}, nil
}


func (d *devices) GetOrCreate(ctx context.Context, device models.Device) (string, int64, error) {
resp := d.db.WithContext(ctx).
		Where("identifier = ?", device.Identifier).
		// Attrs(models.Device{
		// 	ID: device.ID,
		// 	OS:              device.OS,
		// 	Manufacturer:    device.Manufacturer,
		// 	Model:           device.Model,
		// 	ScreenResolution: device.ScreenResolution,
		// 	DeviceType:      device.DeviceType,
		// 	UserID:          device.UserID,
		// }).
		Attrs(&device).
		FirstOrCreate(&device).
		Select("id")

	if resp.Error != nil {
		return "", 0, resp.Error
	}

	return device.ID, resp.RowsAffected, nil
}
