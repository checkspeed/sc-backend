package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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
	databaseUrl := os.Getenv("MIGRATE_DB_URL")
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)
	t.Run("up migration", func(t *testing.T) {
		migrator, err := NewMigrator(store)
		require.NoError(t, err)
		err = migrator.Up(ctx)
		assert.NoError(t, err)
	})
}

func TestXxx(t *testing.T) {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)
	databaseUrl := os.Getenv("OLD_DB_URL")
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)

	var speedTestResult []models.SpeedtestResultsOld
	result := store.db.Table("speed_test_results").Find(&speedTestResult)
	require.NoError(t, result.Error)

	fmt.Println(len(speedTestResult))
	fmt.Println(speedTestResult[23].Latitude, ReduceFloat(speedTestResult[23].Latitude), ReduceFloat(speedTestResult[23].Latitude) == 6.5)
	fmt.Println(speedTestResult[45].Latitude, ReduceFloat(speedTestResult[45].Latitude) == 6.4)

	// Victor details
	deviceID1 := ""
	deviceID2 := ""

	var usCount, kCount, vCount, tCount int
	var newData []models.SpeedTestResults
	for _, input := range speedTestResult {
		if input.ISP == "test" {
			tCount++
			continue
		}

		deviceID := deviceID2
		kCount++
		if ReduceFloat(input.Latitude) == 6.5 {
			kCount--
			vCount++
			deviceID = deviceID1
		}

		state := "Lagos"
		countryCode := "NG"
		CountryName := "Nigeria"
		ContinentCode := "AF"
		ContinentName := "Africa"

		if input.ISP == "Google LLC" || input.ISP == "Google" {
			state = "Florida"
			countryCode = "US"
			CountryName = "United States of America"
			ContinentCode = "NA"
			ContinentName = "North America"

			usCount++
		}

		s := models.SpeedTestResults{
			ID:               uuid.NewString(),
			DownloadSpeed:    input.DownloadSpeed,
			MaxDownloadSpeed: input.MaxDownloadSPeed,
			MinDownloadSpeed: input.MinDownloadSpeed,
			TotalDownload:    input.TotalDownload,
			UploadSpeed:      input.UploadSpeed,
			MaxUploadSpeed:   input.MaxUploadSpeed,
			MinUploadSpeed:   input.MinUploadSpeed,
			TotalUpload:      input.TotalUpload,
			Latency:          input.Latency,
			LoadedLatency:    input.LoadedLatency,
			UnloadedLatency:  input.UnloadedLatency,
			DownloadLatency:  input.DownloadLatency,
			UploadLatency:    input.UploadLatency,
			DeviceID:         deviceID,
			ISP:              input.ISP,
			ISPCode:          input.ISPCode,
			ConnectionType:   input.ConnectionType,
			ConnectionDevice: input.ConnectionDevice,
			TestPlatform:     input.TestPlatform,
			// ServerID:         input.TestServer.ID,
			ServerName:     input.ServerName,
			State:          state,
			CountryCode:    countryCode,
			CountryName:    CountryName,
			ContinentCode:  ContinentCode,
			ContinentName:  ContinentName,
			Longitude:      input.Longitude,
			Latitude:       input.Latitude,
			LocationAccess: input.LocationAccess,
			TestTime:       input.TestTime,
			CreatedAt:      input.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		newData = append(newData, s)
	}

	fmt.Println(len(newData), usCount, vCount, kCount, tCount)
	fmt.Println(newData[0])
	// databaseUrl = os.Getenv("DB_URL")
	// newStore, err := NewStore(databaseUrl)
	// require.NoError(t, err)

	// res := newStore.db.Create(&newData)
	// require.NoError(t, res.Error)

}

func ReduceFloat(value float64) float64 {
	// Format the float to 1 decimal place and then parse it back to a float
	formattedValue := fmt.Sprintf("%.1f", value)
	reducedValue, _ := strconv.ParseFloat(formattedValue, 64)
	return reducedValue
}
