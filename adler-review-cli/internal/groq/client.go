package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client wraps HTTP communication with the Groq API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient initializes a new Groq API Client.
func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key cannot be empty")
	}
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}, nil
}

// ChatCompletionRequest represents the payload for the Groq chat completions endpoint.
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Message represents a single message in the chat context.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response received from the Groq API.
type ChatCompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// GenerateContentWithRetry attempts content generation with an exponential backoff retry logic (up to 3 times) for 429 errors.
func (c *Client) GenerateContentWithRetry(ctx context.Context, modelName string, systemInstruction string, prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	messages := []Message{}
	if systemInstruction != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: systemInstruction,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})

	reqBody := ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		Temperature: 0.3,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	delay := 2 * time.Second

	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt < 3 {
				fmt.Fprintf(os.Stderr, "Network error (attempt %d/3). Retrying in %v...\n", attempt, delay)
				select {
				case <-ctx.Done():
					return "", ctx.Err()
				case <-time.After(delay):
				}
				delay *= 2
				continue
			}
			return "", fmt.Errorf("http request failed after retries: %w", err)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 429 {
			if attempt < 3 {
				fmt.Fprintf(os.Stderr, "Groq rate limit hit (attempt %d/3). Retrying in %v...\n", attempt, delay)
				select {
				case <-ctx.Done():
					return "", ctx.Err()
				case <-time.After(delay):
				}
				delay *= 2
				continue
			}
			return "", fmt.Errorf("rate limit hit and failed after retries: %s", string(bodyBytes))
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("groq api returned error (status %d): %s", resp.StatusCode, string(bodyBytes))
		}

		var chatResp ChatCompletionResponse
		if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
			return "", fmt.Errorf("failed to parse groq response: %w", err)
		}

		if len(chatResp.Choices) == 0 {
			return "", fmt.Errorf("empty response received from Groq")
		}

		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("unexpected error in retry loop")
}
