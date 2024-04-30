package reddit

// Config holds the configuration details for Reddit API access
type Config struct {
	UserAgent    string
	ClientID     string
	ClientSecret string
	AccessToken  string // You might dynamically obtain this, so consider structuring your application to refresh it as needed
}

// NewConfig creates a new Config instance
func NewConfig(userAgent, clientID, clientSecret, accessToken string) *Config {
	return &Config{
		UserAgent:    userAgent,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	}
}

// WithOAuth initializes the Config with OAuth tokens
func (c *Config) WithOAuth(clientID, clientSecret, username, password string) error {
	tokenResponse, err := GetOAuthToken(clientID, clientSecret, username, password)
	if err != nil {
		return err
	}
	c.AccessToken = tokenResponse.AccessToken
	return nil
}
