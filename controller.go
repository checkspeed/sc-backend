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
	var requestBody CreateSpeedtestResultsAPIRequest
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, apiResp{Error: err.Error()})
		return
	}

	stRes, err := transformReq(requestBody)
	if err != nil {
		log.Println("failed to store transform test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, apiResp{Error: err.Error()})
		return
	}
	if err := ct.store.CreateSpeedtestResults(c.Request.Context(), &stRes); err != nil {
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

func transformReq(req CreateSpeedtestResultsAPIRequest) (SpeedtestResults, error) {
	var stRes SpeedtestResults
	stRes.ID = uuid.NewString()

	// download
	arr := strings.Split(req.DownloadSpeed, " ")
	downloadSpeedValue, err := strconv.Atoi(arr[0])
	if err != nil {
		return SpeedtestResults{}, err
	}
	downloadSpeedUnit := strings.ToLower(arr[1])
	// check unit
	switch downloadSpeedUnit {
	case "mbps":
		// convert to kbps
		stRes.DownloadSpeed = downloadSpeedValue * 1000
	case "kbps":
		// convert to kbps
		stRes.DownloadSpeed = downloadSpeedValue
	default:
		return SpeedtestResults{}, fmt.Errorf("invalid download speed unit: %s", downloadSpeedUnit)
	}

	// upload
	arr = strings.Split(req.UploadSpeed, " ")
	uploadSpeedValue, err := strconv.Atoi(arr[0])
	if err != nil {
		return SpeedtestResults{}, err
	}
	uploadSpeedUnit := strings.ToLower(arr[1])
	// check unit
	switch uploadSpeedUnit {
	case "mbps":
		// convert to kbps
		stRes.UploadSpeed = uploadSpeedValue * 1000
	case "kbps":
		// convert to kbps
		stRes.UploadSpeed = uploadSpeedValue
	default:
		return SpeedtestResults{}, fmt.Errorf("invalid upload speed unit: %s", uploadSpeedUnit)
	}

	// latency //
	arr = strings.Split(req.Latency, " ")
	latencyValue, err := strconv.Atoi(arr[0])
	if err != nil {
		return SpeedtestResults{}, err
	}
	latencyUnit := arr[1]
	// check unit
	switch latencyUnit {
	case "ms":
		stRes.UploadSpeed = latencyValue
	
	default:
		return SpeedtestResults{}, fmt.Errorf("invalid latency unit: %s", latencyUnit)
	}

	// loaded latency
	arr = strings.Split(req.LoadedLatency, " ")
	loadedLatencyValue, err := strconv.Atoi(arr[0])
	if err != nil {
		return SpeedtestResults{}, err
	}
	loadedLatencyUnit := arr[1]
	// check unit
	switch loadedLatencyUnit {
	case "ms":
		stRes.UploadSpeed = loadedLatencyValue
	
	default:
		return SpeedtestResults{}, fmt.Errorf("invalid latency unit: %s", loadedLatencyUnit)
	}

	// unloaded latency
	arr = strings.Split(req.UnloadedLatency, " ")
	unloadedLatencyValue, err := strconv.Atoi(arr[0])
	if err != nil {
		return SpeedtestResults{}, err
	}
	unloadedLatencyUnit := arr[1]
	// check unit
	switch unloadedLatencyUnit {
	case "ms":
		stRes.UploadSpeed = unloadedLatencyValue
	
	default:
		return SpeedtestResults{}, fmt.Errorf("invalid latency unit: %s", unloadedLatencyUnit)
	}

	stRes.ISP = req.ISP
	stRes.ServerLocation = req.ServerLocation
	stRes.TestPlatform = req.TestPlatform
	stRes.ServerName = req.ServerName
	stRes.ClientIP = req.ClientIP
	stRes.ClientID = req.ClientID
	stRes.ConnectionDevice = req.ConnectionDevice
	stRes.ConnectionType = req.ConnectionType

	long, err := strconv.ParseFloat(req.Longitude, 64)
	if err != nil {
        return SpeedtestResults{}, err
    }
	stRes.Longitude = long

	lat, err := strconv.ParseFloat(req.Latitude, 64)
	if err != nil {
        return SpeedtestResults{}, err
    }
	stRes.Latitude = lat
	stRes.LocationAccess = req.LocationAccess

	parsedTime, err := time.Parse(Timelayout, req.TestTime)
    if err != nil {
        return SpeedtestResults{}, err
    }
	stRes.TestTime = parsedTime
	stRes.CreatedAt = time.Now()

	return stRes, nil
}
