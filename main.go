package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/controllers"
	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// init db
	store, err := db.NewStore(cfg.DBURL)
	if err != nil {
		log.Fatalf("unable to initialize database, %v \n", err.Error())
	}

	ctrl, err := controllers.NewController(cfg, store)
	if err != nil {
		log.Fatalf("unable to initialize controller, %v \n", err.Error())
	}

	// Initialize rate limiter
	clientLimiter := middleware.NewClientLimiter()

	// Start cleanup routine for rate limiter
	go func() {
		for {
			time.Sleep(10 * time.Minute) // Cleanup every 10 minutes
			clientLimiter.CleanupStaleIPs()
		}
	}()

	// create channel to listen to shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		RunServer(ctrl, cfg, clientLimiter)
	}()

	<-shutdownChan
	log.Println("Closing application")
	store.CloseConn(context.Background())
}

func RunServer(ctrl *controllers.Controller, cfg config.Config, clientLimiter *middleware.ClientLimiter) {
	r := gin.Default()

	// add cors config
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}
	r.Use(cors.New(corsConfig))

	r.GET("/", welcome)
	r.GET("/network", ctrl.GetNetworkInfo)
	r.GET("/geolocation", ctrl.GetNetworkInfo)
	r.POST("/speed_test_result", middleware.RateLimit(clientLimiter), ctrl.CreateSpeedtestResults)
	r.POST("/speed_test_result/list", ctrl.GetSpeedtestResults)
	r.POST("/feedback", ctrl.CreateFeedback)
	r.Run(":" + cfg.Port)

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "welcome to this SpeedCheck api server"})
}
