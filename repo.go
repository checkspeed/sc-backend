package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"context"
	"time"
)

const (
	SpeedtestResultsCollection = "speed_test_results"
	usersCollection            = "users"
)

type SpeedtestResults struct {
	ID string `json:"id"`

	// download
	DownloadSpeed    int `json:"avg_download_speed"` // average | kbps
	MaxDownloadSPeed int `json:"max_download_speed"` // kbps
	MinDownloadSpeed int `json:"min_download_speed"` // kbps
	TotalDownload    int `json:"total_download"`     // kbps

	// upload
	UploadSpeed    int `json:"upload_speed"`     // average | kbps
	MaxUploadSpeed int `json:"max_upload_speed"` // kbps
	MinUploadSpeed int `json:"min_upload_speed"` // kbps
	TotalUpload    int `json:"total_upload"`     // kbps

	// latency
	Latency         int `json:"latency"`          // average | ms
	LoadedLatency   int `json:"loaded_latency"`   // ms
	UnloadedLatency int `json:"unloaded_latency"` // ms
	DownloadLatency int `json:"download_latency"` // ms
	UploadLatency   int `json:"upload_latency"`   // ms

	ConnectionType   string `json:"connection_type"`   // "DSL," "Cable," "Fiber," or "Wireless."
	ConnectionDevice string `json:"connection_device"` // "5G Router," "Mobile," "Fiber," or "Wireless."
	ISP              string `json:"isp"`
	ClientIP         string `json:"client_ip"` // Is this really needed for storage?
	ClientID         string `json:"client_id"` // unique way to identiy the client device
	City             string `json:"city"`
	ServerName       string `json:"server_name"`
	TestServerID     string `json:"test_server_id"`

	Long           float64 `json:"long"`
	Lat            float64 `json:"lat"`
	LocationAccess bool    `json:"location_access"`
	// there should be another field to indicate how accurate

	CreatedAt time.Time `json:"created_at"` // time when record is created
	TestTime  time.Time `json:"test_time"`  // time when the internet test was taken
}

type Datastore interface {
	CreateSpeedtestResults(ctx context.Context, speedTestResult *SpeedtestResults) error
	GetSpeedtestResults(ctx context.Context) ([]SpeedtestResults, error)
}

type mongoStore struct {
	client *mongo.Client
	dbName string
}

// ensure mongostore implements the datastore interface
var _ Datastore = &mongoStore{}

func NewMongoStore(dbUrl, dbName string) (*mongoStore, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return &mongoStore{
		client: client,
		dbName: dbName,
	}, client, nil
}

func (m *mongoStore) col(collectionName string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(collectionName)
}

func (m *mongoStore) CreateSpeedtestResults(ctx context.Context, speedtestResults *SpeedtestResults) error {
	_, err := m.col(SpeedtestResultsCollection).
		InsertOne(ctx, speedtestResults)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoStore) GetSpeedtestResults(ctx context.Context) ([]SpeedtestResults, error) {
	var speedtestResults []SpeedtestResults
	query := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"ts": -1})
	cursor, err := m.col(SpeedtestResultsCollection).Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &speedtestResults); err != nil {
		return nil, err
	}

	return speedtestResults, nil
}
