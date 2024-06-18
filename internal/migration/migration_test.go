package migration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migration(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)
	dbUrl := os.Getenv("TEST_DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	require.NoError(t, err)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	assert.NoError(t, err)

	err = m.Down()
	assert.NoError(t, err)
}

func Test_RunUpMigration(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)
	dbUrl := os.Getenv("TEST_DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	require.NoError(t, err)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	assert.NoError(t, err)
}
