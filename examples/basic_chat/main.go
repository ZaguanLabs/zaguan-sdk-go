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

	// Create a chat completion request
	req := zaguansdk.ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []zaguansdk.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "What is the capital of France?",
			},
		},
	}

	// Send the request
	resp, err := client.Chat(context.Background(), req, nil)
	if err != nil {
		log.Fatalf("Chat request failed: %v", err)
	}

	// Print the response
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens: %d prompt + %d completion = %d total\n",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens)

	// Check for reasoning tokens
	if resp.Usage.HasReasoningTokens() {
		fmt.Printf("Reasoning tokens: %d\n",
			resp.Usage.CompletionTokensDetails.ReasoningTokens)
	}
}
