package repositories

import (
	"context"

	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/models"
)

type TestServers interface {
	GetOrCreate(ctx context.Context, device models.Device) (string, int64, error)
}

type testServers struct {
	db *gorm.DB
}

func New(store db.Store) (testServers, error) {
	return testServers{
		db: store.DB(),
	}, nil
}

// GetOrCreate returns the existing test server id of the test server with the provided identifier creates a new one with the details provided
func (d *testServers) GetOrCreate(ctx context.Context, testServer models.TestServer) (string, int64, error) {
	resp := d.db.WithContext(ctx).
		Select("id").
		Where("identifier = ?", testServer.Identifier).
		Attrs(&testServer).
		FirstOrCreate(&testServer)

	if resp.Error != nil {
		return "", 0, resp.Error
	}

	return testServer.ID, resp.RowsAffected, nil
}
