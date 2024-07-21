package db

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Migration(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(databaseUrl)
	require.NoError(t, err)

	t.Run("down migration", func(t *testing.T) {
		migrator, err := NewMigrator(store)
		require.NoError(t, err)
		err = migrator.Down(ctx)
		assert.NoError(t, err)
	})

	t.Run("up migration", func(t *testing.T) {
		migrator, err := NewMigrator(store)
		require.NoError(t, err)
		err = migrator.Up(ctx)
		assert.NoError(t, err)
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
