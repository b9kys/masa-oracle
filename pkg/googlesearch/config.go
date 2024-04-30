package google

// Config stores the Google API key and Custom Search Engine ID
type Config struct {
	APIKey string
	CSEID  string
}

// NewConfig creates a new instance of Config
func NewConfig(apiKey, cseID string) *Config {
	return &Config{
		APIKey: apiKey,
		CSEID:  cseID,
	}
}
