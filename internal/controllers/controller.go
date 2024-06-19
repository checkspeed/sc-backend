package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"encoding/json"
	"log"
	"net/http"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/models"
	"github.com/checkspeed/sc-backend/internal/repositories"
)

type Controller struct {
	cfg   config.Config
	store repositories.Datastore
}

const Timelayout = "Mon, 02 Jan 2006 15:04:05 MST"

func NewController(cfg config.Config, store repositories.Datastore) *Controller {
	return &Controller{
		cfg,
		store,
	}
}

func (ct *Controller) GetNetworkInfo(c *gin.Context) {
	ipAddr := c.Request.URL.Query().Get("ip")

	geoUrl := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", ct.cfg.GeoAPIKey, ipAddr)
	resp, err := http.Get(geoUrl)
	if err != nil {
		log.Println("error calling ip-api endpoint", err)
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
	}

	defer resp.Body.Close()

	var respBody models.NetworkData
	json.NewDecoder(resp.Body).Decode(&respBody)

	apiResp := models.ApiResp{
		Message: "success",
		Data:    respBody,
	}

	// send response to user
	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) CreateSpeedtestResults(c *gin.Context) {
	var requestBody models.SpeedTestResult
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
		return
	}

	requestBody.ID = uuid.NewString()
	if err := ct.store.CreateSpeedtestResult(c.Request.Context(), &requestBody); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
		return
	}

	apiResp := models.ApiResp{
		Message: "success",
	}

	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) GetSpeedtestResults(c *gin.Context) {
	var filters repositories.GetSpeedTestResultsFilter

	if err := c.BindJSON(&filters); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
		return
	}
	results, err := ct.store.GetSpeedTestResults(c.Request.Context(), filters)
	if err != nil {
		log.Println("failed to retrieve speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
	}

	apiResp := models.ApiResp{
		Message: "success",
		Data:    results,
	}

	c.JSON(http.StatusOK, apiResp)
}
