package googleai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client wraps the genai.Client and implements our high-level features.
type Client struct {
	*genai.Client
}

// NewClient initializes a new GenAI Client using the provided API Key.
func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	c, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create GenAI client: %w", err)
	}
	return &Client{c}, nil
}

// ListAndSelectBestProModel discovers active models and prioritizes newer, highly available versions.
// Falls back to "gemini-2.5-flash" if listing fails.
func (c *Client) ListAndSelectBestProModel(ctx context.Context) string {
	fallback := "gemini-2.5-flash"
	iter := c.ListModels(ctx)

	has25Flash := false
	has15Flash := false
	has25Pro := false
	has15Pro := false

	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Fail gracefully, print warning to stderr and return fallback.
			fmt.Fprintf(os.Stderr, "Warning: Failed to list models: %v. Using fallback %s.\n", err, fallback)
			return fallback
		}

		name := strings.ToLower(m.Name)

		supportsGenerate := false
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				supportsGenerate = true
				break
			}
		}

		if supportsGenerate {
			if strings.Contains(name, "gemini-2.5-flash") {
				has25Flash = true
			} else if strings.Contains(name, "gemini-1.5-flash") {
				has15Flash = true
			} else if strings.Contains(name, "gemini-2.5-pro") {
				has25Pro = true
			} else if strings.Contains(name, "gemini-1.5-pro") {
				has15Pro = true
			}
		}
	}

	if has25Flash {
		return "gemini-2.5-flash"
	}
	if has15Flash {
		return "gemini-1.5-flash"
	}
	if has25Pro {
		return "gemini-2.5-pro"
	}
	if has15Pro {
		return "gemini-1.5-pro"
	}

	return fallback
}

// CountTokens returns the token count for a given text and model.
func (c *Client) CountTokens(ctx context.Context, modelName string, text string) (int, error) {
	model := c.GenerativeModel(modelName)
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}
	return int(resp.TotalTokens), nil
}

// GenerateContentWithRetry attempts content generation with an exponential backoff retry logic (up to 3 times) for 429/RESOURCE_EXHAUSTED errors.
func (c *Client) GenerateContentWithRetry(ctx context.Context, modelName string, systemInstruction string, prompt string) (string, error) {
	model := c.GenerativeModel(modelName)
	if systemInstruction != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemInstruction)},
		}
	}

	var resp *genai.GenerateContentResponse
	var err error
	delay := 2 * time.Second

	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = model.GenerateContent(ctx, genai.Text(prompt))
		if err == nil {
			break
		}

		// Check if it's a rate limit or resource exhausted error
		isRateLimit := false
		if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 429 {
			isRateLimit = true
		} else if strings.Contains(strings.ToLower(err.Error()), "rate limit") ||
			strings.Contains(strings.ToLower(err.Error()), "resource_exhausted") ||
			strings.Contains(strings.ToLower(err.Error()), "429") {
			isRateLimit = true
		}

		if isRateLimit && attempt < 3 {
			fmt.Fprintf(os.Stderr, "Rate limit hit (attempt %d/3). Retrying in %v...\n", attempt, delay)
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
			}
			delay *= 2
			continue
		}

		return "", fmt.Errorf("generation failed after retries: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("empty response received from Gemini")
	}

	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			result.WriteString(string(textPart))
		}
	}

	return result.String(), nil
}
