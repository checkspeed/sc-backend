package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"encoding/json"
	"log"
	"net/http"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/model"
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
		c.JSON(http.StatusBadRequest, model.ApiResp{Error: err.Error()})
	}

	defer resp.Body.Close()

	var respBody model.NetworkData
	json.NewDecoder(resp.Body).Decode(&respBody)

	apiResp := model.ApiResp{
		Message: "success",
		Data:    respBody,
	}

	// send response to user
	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) CreateSpeedtestResults(c *gin.Context) {
	var requestBody model.SpeedtestResults
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, model.ApiResp{Error: err.Error()})
		return
	}

	requestBody.ID = uuid.NewString()
	if err := ct.store.CreateSpeedtestResults(c.Request.Context(), &requestBody); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, model.ApiResp{Error: err.Error()})
		return
	}

	apiResp := model.ApiResp{
		Message: "success",
	}

	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) GetSpeedtestResults(c *gin.Context) {
	var filters repositories.GetSpeedtestResultsFilter
	if err := c.BindJSON(&filters); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, model.ApiResp{Error: err.Error()})
		return
	}
	results, err := ct.store.GetSpeedtestResults(c.Request.Context(), filters)
	if err != nil {
		log.Println("failed to retrieve speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, model.ApiResp{Error: err.Error()})
	}

	apiResp := model.ApiResp{
		Message: "success",
		Data:    results,
	}

	c.JSON(http.StatusOK, apiResp)
}
