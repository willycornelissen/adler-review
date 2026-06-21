package config_test

import (
	"os"
	"testing"

	"adler-review-cli/internal/config"
)

func TestLoadConfig(t *testing.T) {
	// Test prioritizing CLI key
	cfg, err := config.LoadConfig("cli-key", "cli-model")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey != "cli-key" {
		t.Errorf("Expected APIKey 'cli-key', got '%s'", cfg.APIKey)
	}

	if cfg.ModelOverride != "cli-model" {
		t.Errorf("Expected ModelOverride 'cli-model', got '%s'", cfg.ModelOverride)
	}

	// Test environment variable fallback
	os.Setenv("GEMINI_API_KEY", "env-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg, err = config.LoadConfig("", "")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey != "env-key" {
		t.Errorf("Expected APIKey 'env-key', got '%s'", cfg.APIKey)
	}
}
