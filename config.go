package main

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	defaultPort   = "8080"
	defaultSecret = "secret"
	defaultDbUrl  = "mongodb://localhost:27017"
	defaultDbName = "speed_check_db"
)

// Config contain all the config that this application needs
type Config struct {
	Port      string `json:"port"`
	DBName    string `json:"db_name"`
	DBURL     string `json:"dburl"`
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

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dbUrl = defaultDbUrl
	}
	config.DBURL = dbUrl

	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		dbName = defaultDbName
	}
	config.DBName = dbName

	return config
}