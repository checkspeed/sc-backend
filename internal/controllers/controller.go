package controllers

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/models"
)

type Controller struct {
	cfg         config.Config
	devicesRepo db.Devices
	speedTRepo  db.SpeedTestResults
	testSrvRepo db.TestServers
}

const Timelayout = "Mon, 02 Jan 2006 15:04:05 MST"

func NewController(cfg config.Config, store db.Store) (*Controller, error) {
	devicesRepo, err := db.NewDevicesRepo(store)
	if err != nil {
		return nil, err
	}
	testSrvRepo, err := db.NewTestServerRepo(store)
	if err != nil {
		return nil, err
	}
	speedTRepo, err := db.NewSpeedTestResultsRepo(store)
	if err != nil {
		return nil, err
	}
	return &Controller{
		cfg:         cfg,
		devicesRepo: devicesRepo,
		speedTRepo:  speedTRepo,
		testSrvRepo: testSrvRepo,
	}, nil
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

func (ct *Controller) GetGeoLocationInfo(c *gin.Context) {
	ipAddr := c.Request.URL.Query().Get("ip")

	geoUrl := fmt.Sprintf("https://get.geojs.io/v1/ip/geo.json?ip=%s", ipAddr)
	resp, err := http.Get(geoUrl)
	if err != nil {
		log.Println("error calling ip-geo endpoint", err)
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
	}

	defer resp.Body.Close()

	var respBody models.GeoLocationData
	json.NewDecoder(resp.Body).Decode(&respBody)

	apiResp := models.ApiResp{
		Message: "success",
		Data:    respBody,
	}

	// send response to user
	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) CreateSpeedtestResults(c *gin.Context) {
	var requestBody models.CreateSpeedTestResult
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
		return
	}

	// Get or create device if deviceID is not provided in request body
	if requestBody.DeviceID == "" || requestBody.DeviceID == "undefined" {
		deviceIdentifier := Hash([]string{requestBody.Device.OS, requestBody.Device.ScreenResolution, requestBody.Device.DeviceIP})
		device := models.Device{
			ID:         uuid.NewString(),
			Identifier: deviceIdentifier,
			OS:               requestBody.Device.OS,
			DeviceType:       "Desktop",
			Manufacturer:     requestBody.Device.Manufacturer,
			Model:            requestBody.Device.Model,
			ScreenResolution: requestBody.Device.ScreenResolution,
			IsPlatformDevice: true,
		}

		// Get device by identifier if it exists
		deviceID, err := ct.devicesRepo.GetIDByIdentifier(c.Request.Context(), device.Identifier)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Println("failed to get or create device: ", err.Error())
				c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
				return
			}
		}

		// Create device if it doesn't exist
		if deviceID == "" {
			err := ct.devicesRepo.Create(c.Request.Context(), device)
			if err != nil {
				log.Println("failed to get or create device: ", err.Error())
				c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
				return
			}
			deviceID = device.ID
		}

		requestBody.DeviceID = deviceID
	}

	// TODO: Validate provided device id

	speedTestResult, err := transformSpeedTestResult(requestBody)
	if err != nil {
		log.Println("failed to transform input: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
		return
	}
	if err := ct.speedTRepo.Create(c.Request.Context(), &speedTestResult); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ApiResp{Error: err.Error()})
		return
	}

	apiResp := models.CreateSpeedTestResultResponse{
		Message:  "success",
		DeviceID: speedTestResult.DeviceID,
	}

	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) GetSpeedtestResults(c *gin.Context) {
	var filters db.GetSpeedTestResultsFilter

	if err := c.BindJSON(&filters); err != nil {
		log.Println("invalid request body error: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Error: err.Error()})
		return
	}
	results, err := ct.speedTRepo.Get(c.Request.Context(), filters)
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

func transformSpeedTestResult(input models.CreateSpeedTestResult) (models.SpeedTestResults, error) {
	// testTime, err := time.Parse(time.RFC1123, input.TestTime)
	// if err != nil {
	// 	return models.SpeedTestResult{}, err
	// }
	return models.SpeedTestResults{
		ID:               uuid.NewString(),
		DownloadSpeed:    input.DownloadSpeed,
		MaxDownloadSpeed: input.MaxDownloadSpeed,
		MinDownloadSpeed: input.MinDownloadSpeed,
		TotalDownload:    input.TotalDownload,
		UploadSpeed:      input.UploadSpeed,
		MaxUploadSpeed:   input.MaxUploadSpeed,
		MinUploadSpeed:   input.MinUploadSpeed,
		TotalUpload:      input.TotalUpload,
		Latency:          input.Latency,
		LoadedLatency:    input.LoadedLatency,
		UnloadedLatency:  input.UnloadedLatency,
		DownloadLatency:  input.DownloadLatency,
		UploadLatency:    input.UploadLatency,
		DeviceID:         input.DeviceID,
		ISP:              input.ISP,
		ISPCode:          input.ISPCode,
		ConnectionType:   input.ConnectionType,
		ConnectionDevice: input.ConnectionDevice,
		TestPlatform:     input.TestPlatform,
		ServerName:     input.ServerName,
		State:          input.State,
		CountryCode:    input.CountryCode,
		CountryName:    input.CountryName,
		ContinentCode:  input.ContinentCode,
		ContinentName:  input.ContinentName,
		Longitude:      input.Longitude,
		Latitude:       input.Latitude,
		LocationAccess: input.LocationAccess,
		// TestTime:         testTime,
		TestTime:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Hash returns a hash formed from a concat slice of string input
func Hash(input []string) string {
	s := strings.Join(input, "")
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
