package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/checkspeed/sc-backend/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var databaseUrl string

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl = fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	log.Println("postgres started on: ", databaseUrl)
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	store, err := NewStore(databaseUrl)
	if err != nil {
		log.Fatalf("Could not init store: %s", err)
	}
	migrator, err := NewMigrator(store)
	if err != nil {
		log.Fatalf("Could not run migation: %s", err)
	}
	err = migrator.Up(context.Background())
	if err != nil {
		log.Fatalf("Could not run migation: %s", err)
	}

	// run tests
	m.Run()
}

func Test_GetOrCreate(t *testing.T) {
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)

	repo, err := NewDevicesRepo(store)
	require.NoError(t, err)

	t.Run("Get or Create device", func(t *testing.T) {
		ctx := context.Background()

		// define a test device
		testDevice := models.Device{
			ID:         uuid.NewString(),
			Identifier: "unique_device_identifier",
			OS:         "Android",
			DeviceType: "Mobile",
		}

		// perform GetOrCreate operation
		id, rowsAffected, err := repo.GetOrCreate(ctx, testDevice)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, int64(1), rowsAffected)

		// Perform GetOrCreate operation again to test the "get" case
		id2, rowsAffected2, err := repo.GetOrCreate(ctx, testDevice)
		assert.NoError(t, err)
		assert.Equal(t, id, id2)
		assert.Equal(t, int64(0), rowsAffected2)

		// Perform GetByID to validate the created fields
		dbDevice, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, testDevice.Identifier, dbDevice.Identifier)
		assert.Equal(t, testDevice.OS, dbDevice.OS)
		assert.Equal(t, testDevice.DeviceType, dbDevice.DeviceType)

	})
}

func Test_GetIDByIdentifier(t *testing.T) {
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)

	repo, err := NewDevicesRepo(store)
	require.NoError(t, err)

	t.Run("OK -Get by identifier", func(t *testing.T) {
		ctx := context.Background()

		// define a test device
		testDevice := models.Device{
			ID:         uuid.NewString(),
			Identifier: "unique_device_identifier",
			OS:         "Android",
			DeviceType: "Mobile",
		}

		// perform GetOrCreate operation
		err := repo.Create(ctx, testDevice)
		assert.NoError(t, err)

		// Perform GetOrCreate operation again to test the "get" case
		id, err := repo.GetIDByIdentifier(ctx, testDevice.Identifier)
		assert.NoError(t, err)
		assert.Equal(t, testDevice.ID, id)

		// Perform GetByID to validate the created fields
		dbDevice, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, testDevice.Identifier, dbDevice.Identifier)
		assert.Equal(t, testDevice.OS, dbDevice.OS)
		assert.Equal(t, testDevice.DeviceType, dbDevice.DeviceType)

	})
}

func Test_RunManualUpMigration(t *testing.T) {
	ctx := context.Background()
	err := godotenv.Load("../../.env")
	require.NoError(t, err)
	databaseUrl := os.Getenv("TEST_DB_URL")
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)
	t.Run("up migration", func(t *testing.T) {
		migrator, err := NewMigrator(store)
		require.NoError(t, err)
		err = migrator.Up(ctx)
		assert.NoError(t, err)
	})
}