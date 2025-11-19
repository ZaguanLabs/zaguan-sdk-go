// Package zaguansdk provides the official Go client for the Zaguan CoreX API.
//
// Zaguan CoreX is an enterprise-grade AI gateway that unifies access to multiple AI providers
// (OpenAI, Anthropic, Google, etc.) through a single, consistent API. This SDK provides
// idiomatic Go access to all CoreX features, including:
//
//   - OpenAI-compatible Chat Completions
//   - Anthropic-native Messages API
//   - Provider-specific extensions (Reasoning, Caching, etc.)
//   - Credits and Usage tracking
//   - Advanced routing and virtual models
//
// # Basic Usage
//
// The entry point is the Client. You can create a new client with NewClient:
//
//	client := zaguansdk.NewClient(zaguansdk.Config{
//		BaseURL: "https://your-zaguan-instance.com",
//		APIKey:  "your-api-key",
//	})
//
// # Chat Completions (OpenAI Style)
//
// For most use cases, use the Chat method, which mirrors the OpenAI Chat Completion API.
// This method supports all standard features plus Zaguan extensions like provider-specific parameters.
//
//	resp, err := client.Chat(ctx, zaguansdk.ChatRequest{
//		Model: "openai/gpt-4o",
//		Messages: []zaguansdk.Message{
//			{Role: "user", Content: "Hello, world!"},
//		},
//	}, nil)
//
// # Anthropic Messages (Anthropic Style)
//
// If you prefer the Anthropic message format or need specific Anthropic features exposed
// via their native API shape, use the Messages method.
//
//	resp, err := client.Messages(ctx, zaguansdk.MessagesRequest{
//		Model: "anthropic/claude-3-5-sonnet",
//		Messages: []zaguansdk.AnthropicMessage{
//			{Role: "user", Content: "Explain quantum computing"},
//		},
//		MaxTokens: 1000,
//	}, nil)
//
// # Streaming
//
// Both Chat and Messages support streaming via ChatStream and MessagesStream methods.
// These return an iterator-style stream object.
//
//	stream, err := client.ChatStream(ctx, req, nil)
//	defer stream.Close()
//
//	for {
//		chunk, err := stream.Recv()
//		if errors.Is(err, io.EOF) {
//			break
//		}
//		// handle chunk
//	}
//
// # Configuration
//
// The Client can be customized with a custom http.Client, timeouts, and logging.
// RequestOptions can be passed to individual methods to override defaults per request.
//
// # Thread Safety
//
// The Client is safe for concurrent use by multiple goroutines. Configuration should
// not be modified after the client is created.
package zaguansdk
