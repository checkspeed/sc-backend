package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	_ "embed"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateSpeedtestResults(t *testing.T) {
	dbUrl := os.Getenv("TEST_DB_URL")
	ctx := context.Background()
	store, err := NewStore(dbUrl)
	require.NoError(t, err)

	defer store.CloseConn(ctx)
	
	sampleResult := SpeedtestResults{
		ID: uuid.NewString(),
		DownloadSpeed: 15000,
		UploadSpeed: 8000,
		Latency: 27,
		ISP: "test",
		ServerName: "test server",
		ClientID: uuid.NewString(),
	}
	err = store.CreateSpeedtestResults(ctx, &sampleResult)
	assert.NoError(t, err)

	resp, err := store.GetSpeedtestResults(ctx)
	assert.NoError(t, err)
	fmt.Println(resp)
}

func Test_GetSpeedtestResults(t *testing.T) {
	dbUrl := os.Getenv("TEST_DB_URL")
	ctx := context.Background()
	store, err := NewStore(dbUrl)
	require.NoError(t, err)

	defer store.CloseConn(ctx)

	resp, err := store.GetSpeedtestResults(ctx)
	assert.NoError(t, err)
	fmt.Println(resp)
}
