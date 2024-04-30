// reddit/redditprofile.go
package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RedditUserProfile represents a user's profile on Reddit
type RedditUserProfile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Karma    int    `json:"total_karma"`
	AboutURL string `json:"subreddit"`
}

// GetUserProfile fetches the Reddit user profile using OAuth
func GetUserProfile(username string, config *Config) (*RedditUserProfile, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://oauth.reddit.com/user/%s/about", username), nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header to use the Bearer token
	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user profile, status code: %d", resp.StatusCode)
	}

	var data struct {
		Data RedditUserProfile `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Data, nil
}
