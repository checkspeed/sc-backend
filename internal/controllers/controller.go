package controllers

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jomei/notionapi"

	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/checkspeed/sc-backend/internal/config"
	"github.com/checkspeed/sc-backend/internal/db"
	"github.com/checkspeed/sc-backend/internal/models"
	"github.com/checkspeed/sc-backend/internal/utils"
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

func (ct *Controller) CreateFeedback(c *gin.Context) {
	var requestBody models.CreateFeedback

	// Track when request started (used later for logging duration & timestamp)
	startTime := time.Now()

	// 1. Parse JSON body
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("[%s] CreateFeedback - invalid request body: %s", startTime.Format(time.RFC3339), err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Status: models.StatusFail, Message: "Invalid request body",
			Code: "INVALID_BODY"})
		return
	}

	// 2. Validation (using utils)
	requestBody.Message = strings.TrimSpace(requestBody.Message)
	requestBody.Subject = strings.TrimSpace(requestBody.Subject)
	requestBody.Email = strings.TrimSpace(requestBody.Email)

	if !utils.IsValidText(requestBody.Message, true) {
		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "Message must contain letters/numbers and not only special characters",
			Code:    "INVALID_MESSAGE"})
		return
	}
	if len(requestBody.Message) > 5000 {
		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "Message too long (max 5000 characters)",
			Code:    "MESSAGE_TOO_LONG",
		})
		return
	}

	// Validate Subject (optional but must be valid if provided)
	if requestBody.Subject != "" && !utils.IsValidText(requestBody.Subject, false) {
		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "Subject must contain valid characters",
			Code:    "INVALID_SUBJECT"})
		return
	}

	if len(requestBody.Subject) > 200 {
		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "Subject too long (max 200 characters)",
			Code:    "SUBJECT_TOO_LONG",
		})
		return
	}

	//  Validate Email (optional but strict if provided)
	if requestBody.Email != "" && !utils.IsValidEmail(requestBody.Email) {
		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "Invalid email format (example: user@example.com)",
			Code:    "INVALID_EMAIL"})
		return
	}

	// 3. Prepare Notion payload : Initialize Notion client
	client := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_KEY")))
	databaseID := notionapi.DatabaseID(os.Getenv("NOTION_DATABASE_ID"))

	// Create short preview (first 100 chars)
	messagePreview := requestBody.Message
	if len(messagePreview) > 100 {
		messagePreview = messagePreview[:100] + "..."
	}

	// Map feedback to Notion properties
	properties := notionapi.Properties{
		"Message": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: messagePreview}},
			},
		},
	}

	// Optional field: Add Subject only if provided
	if requestBody.Subject != "" {
		properties["Subject"] = notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: requestBody.Subject}},
			},
		}
	}

	// Optional field: Add Email only if provided
	if requestBody.Email != "" {
		properties["Email"] = notionapi.EmailProperty{
			Email: requestBody.Email,
		}
	}

	dateCreated := notionapi.Date(time.Now().UTC())
	properties["Date created"] = notionapi.DateProperty{
		Date: &notionapi.DateObject{
			Start: &dateCreated,
		},
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// 4. Save to Notion
	_, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent:     notionapi.Parent{DatabaseID: databaseID},
		Properties: properties,
		//this will add the message to the page view
		Children: []notionapi.Block{
			&notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{Object: "block", Type: notionapi.BlockTypeParagraph},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{Text: &notionapi.Text{Content: requestBody.Message}},
					},
				},
			},
		},
	})

	if err != nil {
		// Here we use %v, which also calls err.Error() under the hood
		// (just showing both styles so you get familiar with them)
		log.Printf("[%s] CreateFeedback - Notion error: %v",
			startTime.Format(time.RFC3339), err)

		c.JSON(http.StatusInternalServerError, models.ApiResp{
			Status:  models.StatusError,
			Message: "Failed to save feedback",
			Code:    "NOTION_ERROR",
		})
		return
	}

	// 5. Log success (duration = how long the whole request took)
	log.Printf("[%s] CreateFeedback - success duration=%v",
		startTime.Format(time.RFC3339),
		time.Since(startTime),
	)

	c.JSON(http.StatusCreated, models.ApiResp{
		Status:  models.StatusSuccess,
		Message: "Feedback submitted successfully",
		Code:    "SUCCESS",
	})

}

func (ct *Controller) GetNetworkInfo(c *gin.Context) {

	startTime := time.Now()

	ipAddr := c.Request.URL.Query().Get("ip")

	if ipAddr == "" {
		log.Printf("GetNetworkInfo - missing IP parameter from client: %s", c.ClientIP())

		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "IP parameter is required",
			Code:    "MISSING_IP",
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	geoUrl := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", ct.cfg.GeoAPIKey, ipAddr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoUrl, nil)

	if err != nil {
		log.Printf("GetNetworkInfo - failed to create request for IP %s: %v", ipAddr, err)
		c.JSON(http.StatusInternalServerError, models.ApiResp{
			Status:  models.StatusError,
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
		})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("GetNetworkInfo - API call failed for IP %s: %v", ipAddr, err)
		c.JSON(http.StatusBadGateway, models.ApiResp{
			Status:  models.StatusError,
			Message: "Geolocation service unavailable",
			Code:    "SERVICE_ERROR",
		})
		return
	}

	defer resp.Body.Close()

	var respBody models.NetworkData
	json.NewDecoder(resp.Body).Decode(&respBody)

	// Log success with duration
	duration := time.Since(startTime)
	log.Printf("GetNetworkInfo - success for IP %s (took %v)", ipAddr, duration)
	c.JSON(http.StatusOK, models.ApiResp{
		Status:  models.StatusSuccess,
		Message: "Success",
		Data:    respBody,
	})
}

func (ct *Controller) GetGeoLocationInfo(c *gin.Context) {
	startTime := time.Now()

	ipAddr := c.Request.URL.Query().Get("ip")

	if ipAddr == "" {
		log.Printf("GetGeoLocationInfo - missing IP parameter from client: %s", c.ClientIP())

		c.JSON(http.StatusBadRequest, models.ApiResp{
			Status:  models.StatusFail,
			Message: "IP parameter is required",
			Code:    "MISSING_IP",
		})
		return
	}

	geoUrl := fmt.Sprintf("https://get.geojs.io/v1/ip/geo.json?ip=%s", ipAddr)
	resp, err := http.Get(geoUrl)

	if err != nil {
		log.Printf("GetGeoLocationInfo - API call failed: %v", err)

		c.JSON(http.StatusBadGateway, models.ApiResp{
			Status:  models.StatusError,
			Message: "GetGeoLocationInfo service unavailable",
			Code:    "SERVICE_ERROR",
		})
		return
	}

	defer resp.Body.Close()

	var respBody models.GeoLocationData
	json.NewDecoder(resp.Body).Decode(&respBody)

	// send response to user
	duration := time.Since(startTime)

	log.Printf("GetGeoLocationInfo - success for IP %s (took %v)", ipAddr, duration)
	c.JSON(http.StatusOK, models.ApiResp{
		Status:  models.StatusSuccess,
		Message: "Success",
		Data:    respBody,
	})
}

func (ct *Controller) CreateSpeedtestResults(c *gin.Context) {

	var requestBody models.CreateSpeedTestResult

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("CreateSpeedTestResult - invalid request body: ", err.Error())
		c.JSON(http.StatusBadRequest, models.ApiResp{Status: models.StatusFail, Message: "Invalid request body",
			Code: "INVALID_BODY"})
		return
	}

	// Get or create device if deviceID is not provided in request body
	if requestBody.DeviceID == "" || requestBody.DeviceID == "undefined" {
		deviceIdentifier := Hash([]string{requestBody.Device.OS, requestBody.Device.ScreenResolution, requestBody.Device.DeviceIP})
		device := models.Device{
			ID:               uuid.NewString(),
			Identifier:       deviceIdentifier,
			OS:               requestBody.Device.OS,
			DeviceType:       "Desktop",
			Manufacturer:     requestBody.Device.Manufacturer,
			Model:            requestBody.Device.Model,
			ScreenResolution: requestBody.Device.ScreenResolution,
			IsPlatformDevice: true,
		}

		// Get device by identifier if it exists
		deviceID, err := ct.devicesRepo.GetIDByIdentifier(ctx, device.Identifier)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Println("CreateSpeedTestResult - failed to get or create device: ", err.Error())
				c.JSON(http.StatusInternalServerError, models.ApiResp{
					Status:  models.StatusError,
					Message: "failed to get or create device",
					Code:    "INTERNAL_ERROR"})
				return
			}
		}

		// Create device if it doesn't exist
		if deviceID == "" {
			err := ct.devicesRepo.Create(ctx, device)
			if err != nil {
				log.Println("CreateSpeedTestResult - failed to get or create device: ", err.Error())
				c.JSON(http.StatusInternalServerError, models.ApiResp{
					Status:  models.StatusError,
					Message: "failed to get or create device",
					Code:    "INTERNAL_ERROR"})
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
		c.JSON(http.StatusInternalServerError, models.ApiResp{Message: err.Error()})
		return
	}
	if err := ct.speedTRepo.Create(ctx, &speedTestResult); err != nil {
		log.Println("failed to store speed test results: ", err.Error())
		c.JSON(http.StatusInternalServerError, models.ApiResp{Message: err.Error()})
		return
	}

	apiResp := models.CreateSpeedTestResultResponse{
		Message:  "success",
		DeviceID: speedTestResult.DeviceID,
	}

	c.JSON(http.StatusOK, apiResp)
}

func (ct *Controller) GetSpeedtestResults(c *gin.Context) {
	startTime := time.Now()

	var filters db.GetSpeedTestResultsFilter

	if err := c.BindJSON(&filters); err != nil {
		log.Printf("GetSpeedTestResultsFilter - invalid request body: %s", err.Error())

		c.JSON(http.StatusBadRequest, models.ApiResp{Status: models.StatusFail, Message: "Invalid request body",
			Code: "INVALID_BODY"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	results, err := ct.speedTRepo.Get(ctx, filters)
	if err != nil {
		log.Printf("GetSpeedTestResults: failed to retrieve speed test results: %s", err.Error())

		c.JSON(http.StatusInternalServerError, models.ApiResp{Status: models.StatusError, Message: "Internal server error",
			Code: "INTERNAL_ERROR"})
		return
	}

	// send response to user
	duration := time.Since(startTime)

	log.Printf("[%s] GetSpeedTestResults - retrieved success duration=%v",
		startTime.Format(time.RFC3339),
		duration,
	)

	c.JSON(http.StatusOK, models.ApiResp{
		Status: models.StatusSuccess,
		Data:   results,
	})
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
		ServerName:       input.ServerName,
		State:            input.State,
		CountryCode:      input.CountryCode,
		CountryName:      input.CountryName,
		ContinentCode:    input.ContinentCode,
		ContinentName:    input.ContinentName,
		Longitude:        input.Longitude,
		Latitude:         input.Latitude,
		LocationAccess:   input.LocationAccess,
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
