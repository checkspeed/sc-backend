package migration

import (
	"database/sql"
	"os"
	"testing"

	_ "embed"

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

	m := New(db)
	err = m.Up()
	assert.NoError(t, err)

	// err = m.Down()
	// assert.NoError(t, err)
}