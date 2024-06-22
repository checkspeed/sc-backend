package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/checkspeed/sc-backend/internal/config"
// 	"github.com/checkspeed/sc-backend/internal/controllers"
// 	"github.com/checkspeed/sc-backend/internal/db"
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	cfg := config.LoadConfig()

// 	// init db
// 	store, err := db.NewStore(cfg.DBURL)
// 	if err != nil {
// 		log.Fatalf("unable to initialize database, %v \n", err.Error())
// 	}

// 	ctrl := controllers.NewController(cfg, store)

// 	// create channel to listen to shutdown signals
// 	shutdownChan := make(chan os.Signal, 1)
// 	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

// 	go func() {

// 		RunServer(ctrl, cfg)
// 	}()

// 	<-shutdownChan
// 	log.Println("Closing application")
// 	store.CloseConn(context.Background())
// }

// func RunServer(ctrl *controllers.Controller, cfg config.Config) {
// 	r := gin.Default()
	
// 	// add cors config
// 	corsConfig := cors.DefaultConfig()
// 	corsConfig.AllowCredentials = true
// 	corsConfig.AllowAllOrigins = true
// 	corsConfig.AddAllowMethods("OPTIONS")
// 	r.Use(cors.New(corsConfig))

// 	r.GET("/", welcome)
// 	r.GET("/network", ctrl.GetNetworkInfo)
// 	r.POST("/speed_test_result", ctrl.CreateSpeedtestResults)
// 	r.POST("/speed_test_result/list", ctrl.GetSpeedtestResults)

// 	r.Run(":" + cfg.Port)
// }

// func welcome(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "welcome to this draft api servive"})
// }
