package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"encoding/json"
	"log"
	"net/http"
)

type controller struct {
	cfg   Config
	store Datastore
}

const Timelayout = "Mon, 02 Jan 2006 15:04:05 MST"

func NewController(cfg Config, store Datastore) *controller {
	return &controller{
		cfg,
		store,
	}
}

func (ct *controller) GetNetworkInfo(c *gin.Context) {
	ipAddr := c.Request.URL.Query().Get("ip")

	geoUrl := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", ct.cfg.GeoAPIKey, ipAddr)
	resp, err := http.Get(geoUrl)
	if err != nil {
		log.Println("error calling ip-api endpoint", err)
		c.JSON(http.StatusBadRequest, apiResp{Error: err.Error()})
	}

	defer resp.Body.Close()

	var respBody NetworkData
	json.NewDecoder(resp.Body).Decode(&respBody)

	apiResp := apiResp{
		Message: "success",
		Data:    respBody,
	}

	// send response to user
	c.JSON(http.StatusOK, apiResp)
}

func (ct *controller) CreateSpeedtestResults(c *gin.Context) {
	var requestBody SpeedtestResults
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, apiResp{Error: err.Error()})
		return
	}

	if err := ct.store.CreateSpeedtestResults(c.Request.Context(), &requestBody); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, apiResp{Error: err.Error()})
		return
	}

	apiResp := apiResp{
		Message: "success",
	}

	c.JSON(http.StatusOK, apiResp)
}

func (ct *controller) GetSpeedtestResults(c *gin.Context) {
	results, err := ct.store.GetSpeedtestResults(c.Request.Context())
	if err != nil {
		log.Println("failed to retrieve speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, apiResp{Error: err.Error()})
	}

	apiResp := apiResp{
		Message: "success",
		Data:    results,
	}

	c.JSON(http.StatusOK, apiResp)
}
