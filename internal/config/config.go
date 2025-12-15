package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey string
}

// Load reads configuration from environment variables
// It first tries to load from .env file, then falls back to system env vars
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("EXCHANGE_RATE_API_KEY environment variable is not set.\n" +
			"Please create a .env file with your API key or set the environment variable.\n" +
			"Get your free API key at: https://www.exchangerate-api.com/")
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}
