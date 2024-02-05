package migration

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migration(t *testing.T) {
	dbUrl := "postgres://uchkkzlm:5FIuFpVS2fd9bwqO3hfeODqBZ29_URlB@jelani.db.elephantsql.com/uchkkzlm"
	db, err := sql.Open("postgres", dbUrl)
	require.NoError(t, err)

	m := New(db)
	err = m.Up()
	assert.NoError(t, err)

	// err = m.Down()
	// assert.NoError(t, err)
}