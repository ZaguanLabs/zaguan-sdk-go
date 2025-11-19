package main

import (
	"context"
	"fmt"
	"log"
	"os"

	zaguansdk "github.com/ZaguanLabs/zaguan-sdk-go/sdk"
)

func main() {
	// Get configuration from environment
	baseURL := os.Getenv("ZAGUAN_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	apiKey := os.Getenv("ZAGUAN_API_KEY")
	if apiKey == "" {
		log.Fatal("ZAGUAN_API_KEY environment variable is required")
	}

	// Create a new Zaguan client
	client := zaguansdk.NewClient(zaguansdk.Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
	})

	// Create an Anthropic Messages request with extended thinking
	req := zaguansdk.MessagesRequest{
		Model: "anthropic/claude-3-5-sonnet-20241022",
		Messages: []zaguansdk.AnthropicMessage{
			{
				Role:    "user",
				Content: "Explain the concept of quantum entanglement in simple terms.",
			},
		},
		MaxTokens: 1000,
		Thinking: &zaguansdk.AnthropicThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 5000,
		},
	}

	// Send the request
	resp, err := client.Messages(context.Background(), req, nil)
	if err != nil {
		log.Fatalf("Messages request failed: %v", err)
	}

	// Print the response
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Stop Reason: %s\n\n", resp.StopReason)

	// Print content blocks
	for i, block := range resp.Content {
		fmt.Printf("Block %d (%s):\n", i+1, block.Type)

		switch block.Type {
		case "thinking":
			fmt.Printf("  [Thinking]: %s\n\n", block.Thinking)
		case "text":
			fmt.Printf("  %s\n\n", block.Text)
		}
	}

	// Print token usage
	fmt.Printf("Tokens: %d input + %d output = %d total\n",
		resp.Usage.InputTokens,
		resp.Usage.OutputTokens,
		resp.Usage.InputTokens+resp.Usage.OutputTokens)

	// Print cache information if available
	if resp.Usage.CacheReadInputTokens > 0 {
		fmt.Printf("Cache: %d tokens read from cache\n", resp.Usage.CacheReadInputTokens)
	}
}
