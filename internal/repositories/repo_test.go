package repositories

// "context"
// "fmt"
// "os"
// "testing"

// _ "embed"

// "github.com/checkspeed/sc-backend/internal/model"
// "github.com/google/uuid"
// "github.com/joho/godotenv"
// _ "github.com/lib/pq"
// "github.com/stretchr/testify/assert"
// "github.com/stretchr/testify/require"

// func Test_CreateSpeedtestResults(t *testing.T) {
// 	err := godotenv.Load()
// 	require.NoError(t, err)
// 	dbUrl := os.Getenv("TEST_DB_URL")
// 	ctx := context.Background()
// 	store, err := NewStore(dbUrl)
// 	require.NoError(t, err)

// 	defer store.CloseConn(ctx)

// 	sampleResult := model.SpeedtestResults{
// 		ID:            uuid.NewString(),
// 		DownloadSpeed: 15000,
// 		UploadSpeed:   8000,
// 		Latency:       27,
// 		ISP:           "test",
// 	}
// 	err = store.CreateSpeedtestResults(ctx, &sampleResult)
// 	assert.NoError(t, err)

// 	_, err = store.GetSpeedtestResults(ctx, GetSpeedtestResultsFilter{})
// 	assert.NoError(t, err)
// 	// fmt.Println(resp)
// }

// func Test_GetSpeedtestResults(t *testing.T) {
// 	err := godotenv.Load()
// 	require.NoError(t, err)
// 	dbUrl := os.Getenv("TEST_DB_URL")
// 	ctx := context.Background()
// 	store, err := NewStore(dbUrl)
// 	require.NoError(t, err)

// 	defer store.CloseConn(ctx)

// 	resp, err := store.GetSpeedtestResults(ctx, GetSpeedtestResultsFilter{})
// 	assert.NoError(t, err)
// 	fmt.Println(resp)
// }
