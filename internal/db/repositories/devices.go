package repositories

import (
	"context"

	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/models"
)

type Devices interface {
	GetOrCreate(ctx context.Context, device models.Device) (string, int64, error)
}

type devices struct {
	db *gorm.DB
}

func NewDevicesRepo(store db.Store) (devices, error) {
	return devices{
		db: store.DB(),
	}, nil
}

// GetOrCreate returns existing device id of device connected to hash proviced or creates a new one with the details provided
func (d *devices) GetOrCreate(ctx context.Context, device models.Device) (string, int64, error) {
	resp := d.db.WithContext(ctx).
		Select("id").
		Where("identifier = ?", device.Identifier).
		Attrs(&device).
		FirstOrCreate(&device)

	if resp.Error != nil {
		return "", 0, resp.Error
	}

	return device.ID, resp.RowsAffected, nil
}
