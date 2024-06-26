package db

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	_ "embed"

// 	"github.com/checkspeed/sc-backend/internal/models"
// 	"github.com/google/uuid"
// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func Test_CreateSpeedtestResults(t *testing.T) {
// 	err := godotenv.Load()
// 	require.NoError(t, err)
// 	dbUrl := os.Getenv("DB_URL")
// 	ctx := context.Background()
// 	store, err := NewStore(dbUrl)
// 	require.NoError(t, err)

// 	defer store.CloseConn(ctx)

// 	sampleResult := models.SpeedTestResult{
// 		ID:            uuid.NewString(),
// 		DownloadSpeed: 15000,
// 		UploadSpeed:   8000,
// 		Latency:       27,
// 		ISP:           "test",
// 	}
// 	err = store.CreateSpeedtestResult(ctx, &sampleResult)
// 	assert.NoError(t, err)

// 	_, err = store.GetSpeedTestResults(ctx, GetSpeedTestResultsFilter{})
// 	assert.NoError(t, err)
// 	// fmt.Println(resp)
// }

// func Test_GetSpeedtestResults(t *testing.T) {
// 	err := godotenv.Load()
// 	require.NoError(t, err)
// 	dbUrl := os.Getenv("DB_URL")
// 	ctx := context.Background()
// 	store, err := NewStore(dbUrl)
// 	require.NoError(t, err)

// 	defer store.CloseConn(ctx)

// 	resp, err := store.GetSpeedTestResults(ctx, GetSpeedTestResultsFilter{})
// 	assert.NoError(t, err)
// 	fmt.Println(resp)
// }
