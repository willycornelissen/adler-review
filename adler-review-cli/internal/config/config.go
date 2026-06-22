package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the configuration loaded from flags and environment.
type Config struct {
	Provider      string // "gemini" or "groq"
	APIKey        string
	ModelOverride string
}

// LoadConfig loads variables from .env (if present) and returns a Config.
// It prioritizes CLI arguments over environment variables.
func LoadConfig(cliProvider, cliKey, cliModel string) (*Config, error) {
	// Ignore error if .env file doesn't exist
	_ = godotenv.Load()

	provider := strings.ToLower(cliProvider)
	if provider == "" {
		provider = strings.ToLower(os.Getenv("PROVIDER"))
		if provider == "" {
			provider = "gemini"
		}
	}

	apiKey := cliKey
	modelOverride := cliModel

	if provider == "groq" {
		if apiKey == "" {
			apiKey = os.Getenv("GROQ_API_KEY")
		}
		if modelOverride == "" {
			modelOverride = os.Getenv("GROQ_MODEL")
			if modelOverride == "" {
				modelOverride = "llama-3.3-70b-versatile"
			}
		}
	} else {
		provider = "gemini"
		if apiKey == "" {
			apiKey = os.Getenv("GEMINI_API_KEY")
		}
		if modelOverride == "" {
			modelOverride = os.Getenv("GEMINI_MODEL")
		}
	}

	return &Config{
		Provider:      provider,
		APIKey:        apiKey,
		ModelOverride: modelOverride,
	}, nil
}
