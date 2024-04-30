package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SearchResult represents a single search result
type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// SearchResponse represents the Google Search API response
type SearchResponse struct {
	Items []SearchResult `json:"items"`
}

// PerformSearch performs a search query using Google Custom Search Engine
func PerformSearch(query string, config *Config) ([]SearchResult, error) {
	// Construct the search URL
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		config.APIKey, config.CSEID, query)

	// Make the HTTP request
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and decode the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, err
	}

	return searchResponse.Items, nil
}
