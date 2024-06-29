package db

// import (
// 	// "database/sql"
// 	// "fmt"
// 	// "log"
// 	"context"
// 	"testing"
// 	// "time"

// 	"github.com/checkspeed/sc-backend/internal/db"
// 	_ "github.com/lib/pq"

// 	// "github.com/ory/dockertest/v3"
// 	// "github.com/ory/dockertest/v3/docker"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // var databaseUrl string

// // func TestMain(m *testing.M) {
// // 	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
// // 	pool, err := dockertest.NewPool("")
// // 	if err != nil {
// // 		log.Fatalf("Could not construct pool: %s", err)
// // 	}

// // 	err = pool.Client.Ping()
// // 	if err != nil {
// // 		log.Fatalf("Could not connect to Docker: %s", err)
// // 	}

// // 	// pulls an image, creates a container based on it and runs it
// // 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// // 		Repository: "postgres",
// // 		Tag:        "11",
// // 		Env: []string{
// // 			"POSTGRES_PASSWORD=secret",
// // 			"POSTGRES_USER=user_name",
// // 			"POSTGRES_DB=dbname",
// // 			"listen_addresses = '*'",
// // 		},
// // 	}, func(config *docker.HostConfig) {
// // 		// set AutoRemove to true so that stopped container goes away by itself
// // 		config.AutoRemove = true
// // 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// // 	})
// // 	if err != nil {
// // 		log.Fatalf("Could not start resource: %s", err)
// // 	}

// // 	hostAndPort := resource.GetHostPort("5432/tcp")
// // 	databaseUrl = fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

// // 	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

// // 	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// // 	pool.MaxWait = 120 * time.Second
// // 	if err = pool.Retry(func() error {
// // 		db, err := sql.Open("postgres", databaseUrl)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		return db.Ping()
// // 	}); err != nil {
// // 		log.Fatalf("Could not connect to docker: %s", err)
// // 	}

// // 	log.Println("postgres started on: ")
// // 	defer func() {
// // 		if err := pool.Purge(resource); err != nil {
// // 			log.Fatalf("Could not purge resource: %s", err)
// // 		}
// // 	}()

// // 	// run tests
// // 	m.Run()
// // }

// // postgres://aesuxgse:26nyPHEXNXX0xBJ_oeu1eq61l8BGNI3P@silly.db.elephantsql.com/aesuxgse

// func Test_RunUpMigration(t *testing.T) {
// 	ctx := context.Background()
// 	store, err := db.NewStore(databaseUrl)
// 	require.NoError(t, err)
// 	t.Run("up migration", func(t *testing.T) {
// 		migrator, err := db.NewMigrator(store)
// 		require.NoError(t, err)
// 		err = migrator.Up(ctx)
// 		assert.NoError(t, err)
// 	})
// 	t.Run("Down migration", func(t *testing.T) {
// 		migrator, err := db.NewMigrator(store)
// 		require.NoError(t, err)
// 		err = migrator.Down(ctx)
// 		assert.NoError(t, err)
// 	})
// }
