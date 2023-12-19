package main

import (
	"github.com/gin-gonic/gin"

	"encoding/json"
	"log"
	"net/http"
)

type controller struct {
	store Datastore
}

func NewController(store Datastore) *controller {
	return &controller{
		store: store,
	}
}

// get network info (isp) from 3rd party
func getNetworkInfo(c *gin.Context) {
	ipAddr := c.Request.URL.Query().Get("ip")

	log.Println("received request ip:  ", ipAddr)

	resp, err := http.Get("http://ip-api.com/json/" + ipAddr)
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
	}

	// TODO: create seperate request body struct
	// handle request body validations
	if err := ct.store.CreateSpeedtestResults(c.Request.Context(), &requestBody); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, apiResp{Error: err.Error()})
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
		Data: results,
	}

	c.JSON(http.StatusOK, apiResp)
}
