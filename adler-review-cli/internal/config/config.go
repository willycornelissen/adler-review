package config

import (
	"os"
	"github.com/joho/godotenv"
)

// Config holds the configuration loaded from flags and environment.
type Config struct {
	APIKey        string
	ModelOverride string
}

// LoadConfig loads variables from .env (if present) and returns a Config.
// It prioritizes CLI arguments over environment variables.
func LoadConfig(cliKey, cliModel string) (*Config, error) {
	// Ignore error if .env file doesn't exist
	_ = godotenv.Load()

	apiKey := cliKey
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	modelOverride := cliModel
	if modelOverride == "" {
		modelOverride = os.Getenv("GEMINI_MODEL")
	}

	return &Config{
		APIKey:        apiKey,
		ModelOverride: modelOverride,
	}, nil
}
