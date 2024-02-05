package main

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	defaultPort  = "8080"
	defaultDbUrl = "postgresql://localhost:5432"
)

// Config contain all the config that this application needs
type Config struct {
	Port      string
	GeoAPIKey string
	DBURL     string
}

// LoadConfig loads Config from the environment and returns it
// if a .env file is present, it would be loaded first
// default values are also set
func LoadConfig(filename ...string) Config {
	f := ".env"
	if len(filename) > 0 {
		f = filename[0]
	}
	_ = godotenv.Load(f)
	config := Config{}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	config.Port = port

	dbUrl, ok := os.LookupEnv("DB_URL")
	if !ok {
		dbUrl = defaultDbUrl
	}
	config.DBURL = dbUrl

	geoAPIKey, ok := os.LookupEnv("GEO_API_KEY")
	if !ok {
		geoAPIKey = ""
	}
	config.GeoAPIKey = geoAPIKey

	return config
}
