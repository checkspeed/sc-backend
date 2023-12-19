package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type apiResp struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
}

type NetworkData struct {
	Isp string `json:"isp,omitempty"`
}

func main() {
	cfg := LoadConfig()

	// init db
	store, client, err := NewMongoStore(cfg.DBURL, cfg.DBName)
	if err != nil {
		log.Fatalf("unable to initialize database, %v \n", err.Error())
	}

	ctrl := NewController(store)
	

	// create channel to listen to shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		 RunServer(ctrl, cfg)
	}()

	<-shutdownChan
	log.Println("Closing application")
	client.Disconnect(context.Background())
}

func RunServer(ctrl *controller, cfg Config) {
	r := gin.Default()
	
	r.GET("/", welcome)
	r.GET("/network", getNetworkInfo)
	r.POST("/speed_test_result", ctrl.CreateSpeedtestResults)
	r.GET("/speed_test_result/list", ctrl.CreateSpeedtestResults)


	r.Run(":"+cfg.Port)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "welcome to this draft api servive"})
}
