package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TokenResponse represents the response from Reddit OAuth token request
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// GetOAuthToken fetches the OAuth token using client credentials
func GetOAuthToken(clientID, clientSecret, username, password string) (*TokenResponse, error) {
	client := &http.Client{}
	url := "https://www.reddit.com/api/v1/access_token"

	data := fmt.Sprintf("grant_type=password&username=%s&password=%s", username, password)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("User-Agent", "YourApp/0.1 by YourRedditUsername")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

// RefreshOAuthToken fetches a new OAuth token using a refresh token
func RefreshOAuthToken(clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	client := &http.Client{}
	url := "https://www.reddit.com/api/v1/access_token"

	// Data for the refresh token request
	data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", refreshToken)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("User-Agent", "YourApp/0.1 by YourRedditUsername")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
