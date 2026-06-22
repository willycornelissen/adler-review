package config_test

import (
	"os"
	"testing"

	"adler-review-cli/internal/config"
)

func TestLoadConfig(t *testing.T) {
	// Test prioritizing CLI key with Gemini
	cfg, err := config.LoadConfig("gemini", "cli-key", "cli-model")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Provider != "gemini" {
		t.Errorf("Expected Provider 'gemini', got '%s'", cfg.Provider)
	}

	if cfg.APIKey != "cli-key" {
		t.Errorf("Expected APIKey 'cli-key', got '%s'", cfg.APIKey)
	}

	if cfg.ModelOverride != "cli-model" {
		t.Errorf("Expected ModelOverride 'cli-model', got '%s'", cfg.ModelOverride)
	}

	// Test environment variable fallback for Gemini
	os.Setenv("GEMINI_API_KEY", "env-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg, err = config.LoadConfig("", "", "")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey != "env-key" {
		t.Errorf("Expected APIKey 'env-key', got '%s'", cfg.APIKey)
	}

	// Test Groq provider CLI override
	cfg, err = config.LoadConfig("groq", "groq-key", "gemma2-9b-it")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Provider != "groq" {
		t.Errorf("Expected Provider 'groq', got '%s'", cfg.Provider)
	}

	if cfg.APIKey != "groq-key" {
		t.Errorf("Expected APIKey 'groq-key', got '%s'", cfg.APIKey)
	}

	if cfg.ModelOverride != "gemma2-9b-it" {
		t.Errorf("Expected ModelOverride 'gemma2-9b-it', got '%s'", cfg.ModelOverride)
	}

	// Test Groq default model fallback
	cfg, err = config.LoadConfig("groq", "groq-key", "")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ModelOverride != "llama-3.3-70b-versatile" {
		t.Errorf("Expected ModelOverride 'llama-3.3-70b-versatile', got '%s'", cfg.ModelOverride)
	}
}
