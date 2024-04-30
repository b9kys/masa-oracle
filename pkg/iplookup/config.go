package iplookup

// Config holds any configuration that the IP lookup functionality might need.
type Config struct {
	// APIKey is used for services that require an API key for IP lookups, such as IP-API Pro.
	APIKey string
	// Fields to specify which fields to include in the IP-API response. This is optional.
	// Use the field numbers as specified in the IP-API documentation. For example, "61439" for all fields.
	Fields string
	// Lang specifies the language of the IP-API response. This is optional.
	// Use ISO 639-1 language codes. For example, "en" for English.
	Lang string
}

// NewConfig creates a new instance of Config with the provided settings.
// This function can be expanded to include more parameters as your configuration needs grow.
func NewConfig(apiKey, fields, lang string) *Config {
	return &Config{
		APIKey: apiKey,
		Fields: fields,
		Lang:   lang,
	}
}
