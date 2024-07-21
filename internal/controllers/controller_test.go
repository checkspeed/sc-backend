package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"

	"testing"
	"time"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/controllers"
	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/models"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var store db.Store

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
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

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

	fmt.Println(databaseUrl)
	store, err = db.NewStore(databaseUrl)
	if err != nil {
		log.Fatalf("Could not init store: %s", err)
	}
	migrator, err := db.NewMigrator(store)
	if err != nil {
		log.Fatalf("Could not run migation: %s", err)
	}

	// Get the directory of this test file
	_, testFilePath, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFilePath)
	migrationsPath := filepath.Join(testDir, "../db/migrations")
	fmt.Println("mig path t", migrationsPath, testFilePath, testDir)
	err = migrator.Down(context.Background(), migrationsPath)
	if err != nil {
		log.Fatalf("Could not run migation: %s", err)
	}
	err = migrator.Up(context.Background(), migrationsPath)
	if err != nil {
		log.Fatalf("Could not run migation: %s", err)
	}

	// run tests
	m.Run()
}

func Test_CreateSpeedtestResults(t *testing.T) {
	// ok speed test results saved
	// ok empty device id

	// setup
	// init db\
	// init repos
	// init controller
	cfg := config.Config{}
	ctrl, err := controllers.NewController(cfg, store)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		requestJson string
	}{
		{
			name: "test",
			requestJson: `{
				"download_speed":320000,
				"total_downlad":490000,
				"upload_speed":78000,
				"total_upload":210000,
				"latency":35,
				"loaded_latency":52,
				"unloaded_latency":18,
				"isp":"MTN Nigeria",
				"isp_code":"MTN",
				"test_platform":"fast.com",
				"connection_type":"4g",
				"server_name":"Ojota, NG&nbsp;&nbsp;|&nbsp;&nbsp;Secaucus, US",
				"city":"Lagos, NG",
				"longitude":3.37921,
				"latitude":6.52438,
				"test_time":"2024-07-03T11:10:25.223Z"
			}`,
		},
		{
			name: "test 2",
			requestJson: `{
				"download_speed":19000,
				"total_downlad":20000,
				"upload_speed":7200,
				"total_upload":30000,
				"latency":46,
				"loaded_latency":58,
				"unloaded_latency":34,
				"device_id":"undefined",
				"isp":"Starlink Internet Services Nigeria Ltd",
				"isp_code":"STARLINK",
				"test_platform":"fast.com",
				"connection_type":"4g",
				"server_name":"Harare, ZW - Lagos, NG - nairobi, KE",
				"state":"Lagos",
				"country_code":"NG",
				"country_name":"Nigeria",
				"continent_code":"AF",
				"continent_name":"Africa",
				"longitude":3.38876,
				"latitude":6.4547,
				"test_time":"2024-07-08T21:46:42.279Z",
				"device":{
					"id":"undefined",
					"identifier":"4bfae97e456bb2db53537f95f4bf117773f0578424edd34006ee3f0822225cd4",
					"os":"Windows",
					"screen_resolution":"3440x1440"
					}
				}`,
		},
		{
			name: "test 2",
			requestJson: `{
				"download_speed":19000,
				"total_downlad":20000,
				"upload_speed":7200,
				"total_upload":30000,
				"latency":46,
				"loaded_latency":58,
				"unloaded_latency":34,
				"device_id":"undefined",
				"isp":"Starlink Internet Services Nigeria Ltd",
				"isp_code":"STARLINK",
				"test_platform":"fast.com",
				"connection_type":"4g",
				"server_name":"Harare, ZW - Lagos, NG - nairobi, KE",
				"state":"Lagos",
				"country_code":"NG",
				"country_name":"Nigeria",
				"continent_code":"AF",
				"continent_name":"Africa",
				"longitude":3.38876,
				"latitude":6.4547,
				"test_time":"2024-07-08T21:46:42.279Z",
				"device":{
					"id":"undefined",
					"device_ip":"some-ip-addr",
					"os":"Windows",
					"screen_resolution":"3440x1440"
					}
				}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the router
			router := gin.Default()
			router.POST("/speedtest", ctrl.CreateSpeedtestResults)

			// Create a request
			req, err := http.NewRequest(http.MethodPost, "/speedtest", bytes.NewBuffer([]byte(tc.requestJson)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response models.CreateSpeedTestResultResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			fmt.Println("response: ", response)
			require.NoError(t, err)
			assert.Nil(t, nil)

		})
	}
}