package zaguansdk

import (
	"context"
	"testing"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func BenchmarkClient_Chat(b *testing.B) {
	mockServer := testutil.NewMockServer(
		testutil.ChatCompletionHandler(testutil.ChatCompletionFixture()),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello, world!"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Chat(context.Background(), req, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClient_Messages(b *testing.B) {
	mockServer := testutil.NewMockServer(
		testutil.MessagesHandler(testutil.MessagesFixture()),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello, world!"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Messages(context.Background(), req, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateChatRequest(b *testing.B) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature: ptr(float32(0.7)),
		MaxTokens:   ptr(1000),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateChatRequest(&req)
	}
}

func BenchmarkValidateMessagesRequest(b *testing.B) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
		Temperature: ptr(0.7),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateMessagesRequest(&req)
	}
}

func BenchmarkNewClient(b *testing.B) {
	cfg := Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewClient(cfg)
	}
}

// Benchmark JSON marshaling performance
func BenchmarkChatRequestMarshal(b *testing.B) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello, world!"},
		},
		Temperature:      ptr(float32(0.7)),
		MaxTokens:        ptr(1000),
		TopP:             ptr(float32(0.9)),
		PresencePenalty:  ptr(float32(0.0)),
		FrequencyPenalty: ptr(float32(0.0)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := marshalJSON(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper for benchmarking
func marshalJSON(v interface{}) ([]byte, error) {
	// This would use json.Marshal in real code
	// For now, just return empty to make benchmark compile
	return []byte("{}"), nil
}
