package iplookup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IPApiResponse represents the JSON response structure from IP-API
type IPApiResponse struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

// GetCountryFromIP takes an IP address and returns detailed location information
// apiKey parameter is your IP-API Pro key
func GetCountryFromIP(ip string, apiKey string) (IPApiResponse, error) {
	// Adjust the URL to use the pro version and include the API key
	url := fmt.Sprintf("https://pro.ip-api.com/json/%s?key=%s", ip, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return IPApiResponse{}, fmt.Errorf("error making request to IP-API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPApiResponse{}, fmt.Errorf("error reading response from IP-API: %v", err)
	}

	var apiResponse IPApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return IPApiResponse{}, fmt.Errorf("error parsing JSON response from IP-API: %v", err)
	}

	if apiResponse.Status != "success" {
		return IPApiResponse{}, fmt.Errorf("IP-API error: %s", apiResponse.Message)
	}

	return apiResponse, nil
}
