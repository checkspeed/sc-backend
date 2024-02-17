package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type apiResp struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
}

type NetworkData struct {
	Isp string `json:"isp,omitempty"`
	Longitude string `json:"longitude"`
	Latitude string `json:"latitude"`
}

func main() {
	cfg := LoadConfig()

	// init db
	store, err := NewStore(cfg.DBURL)
	if err != nil {
		log.Fatalf("unable to initialize database, %v \n", err.Error())
	}

	ctrl := NewController(cfg, store)
	

	// create channel to listen to shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		 RunServer(ctrl, cfg)
	}()

	<-shutdownChan
	log.Println("Closing application")
	store.CloseConn(context.Background())
}

func RunServer(ctrl *controller, cfg Config) {
	r := gin.Default()
	
	// add cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))

	r.GET("/", welcome)
	r.GET("/network", ctrl.GetNetworkInfo)
	r.POST("/speed_test_result", ctrl.CreateSpeedtestResults)
	r.GET("/speed_test_result/list", ctrl.GetSpeedtestResults)


	r.Run(":"+cfg.Port)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "welcome to this draft api servive"})
}
