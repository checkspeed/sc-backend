package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type apiResp struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    NetworkData `json:"data,omitempty"`
}

type NetworkData struct {
	Isp string `json:"isp,omitempty"`
}

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/network", getNetworkInfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, nil)
}

func welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "welcome to this draft api servive")
}

// get network info (isp) from 3rd party
func getNetworkInfo(w http.ResponseWriter, req *http.Request) {
	ipAddr := req.URL.Query().Get("ip")

	log.Println("received request ip:  ", ipAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := http.Get("http://ip-api.com/json/" + ipAddr)
	if err != nil {
		log.Println("error calling ip-api endpoint", err)

		json.NewEncoder(w).Encode(apiResp{Error: err.Error()})
	}
	defer resp.Body.Close()

	var respBody NetworkData
	json.NewDecoder(resp.Body).Decode(&respBody)

	apiResp := apiResp{
		Message: "success",
		Data:    respBody,
	}

	log.Println("api response:  ", apiResp)

	// send response to user
	json.NewEncoder(w).Encode(apiResp)
}
