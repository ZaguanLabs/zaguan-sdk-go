// Package zaguansdk provides the official Go client for the Zaguan CoreX API.
//
// Zaguan CoreX is an enterprise-grade AI gateway that unifies access to 15+ AI providers
// (OpenAI, Anthropic, Google, DeepSeek, Groq, Perplexity, xAI, and more) through a single,
// consistent API. This SDK provides idiomatic Go access to all CoreX features, including:
//
//   - OpenAI-compatible Chat Completions (streaming & non-streaming)
//   - Anthropic-native Messages API with extended thinking
//   - Embeddings for semantic search and clustering
//   - Audio transcription, translation, and speech synthesis
//   - Image generation with DALL-E
//   - Content moderation and safety checks
//   - Batch processing with 50% cost reduction
//   - Provider-specific extensions (Reasoning, Caching, etc.)
//   - Credits and Usage tracking
//   - Model discovery and capability queries
//
// # Basic Usage
//
// The entry point is the Client. You can create a new client with NewClient:
//
//	client := zaguansdk.NewClient(zaguansdk.Config{
//		BaseURL: "https://api.zaguanai.com",
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
//		Model: "anthropic/claude-3-5-sonnet-20241022",
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
// # Embeddings
//
// Create vector embeddings for semantic search, clustering, and recommendations:
//
//	resp, err := client.CreateEmbeddings(ctx, zaguansdk.EmbeddingsRequest{
//		Model: "openai/text-embedding-3-small",
//		Input: []string{"Hello", "World"},
//	}, nil)
//
// # Audio
//
// Transcribe audio to text, translate to English, or generate speech:
//
//	// Transcription
//	resp, err := client.CreateTranscription(ctx, zaguansdk.AudioTranscriptionRequest{
//		File:  "/path/to/audio.mp3",
//		Model: "openai/whisper-1",
//	}, nil)
//
//	// Speech synthesis
//	audio, err := client.CreateSpeech(ctx, zaguansdk.AudioSpeechRequest{
//		Model: "openai/tts-1",
//		Input: "Hello, world!",
//		Voice: "alloy",
//	}, nil)
//	defer audio.Close()
//
// # Images
//
// Generate images from text prompts using DALL-E:
//
//	resp, err := client.CreateImage(ctx, zaguansdk.ImageGenerationRequest{
//		Prompt: "A cute baby sea otter",
//		Model:  "openai/dall-e-3",
//		Size:   "1024x1024",
//	}, nil)
//
// # Moderations
//
// Check content for policy violations:
//
//	resp, err := client.CreateModeration(ctx, zaguansdk.ModerationRequest{
//		Input: "Content to check",
//	}, nil)
//	if resp.Results[0].Flagged {
//		// Handle flagged content
//	}
//
// # Batches
//
// Process multiple requests asynchronously with 50% cost reduction:
//
//	resp, err := client.CreateBatch(ctx, zaguansdk.BatchRequest{
//		InputFileID:      "file-abc123",
//		Endpoint:         "/v1/chat/completions",
//		CompletionWindow: "24h",
//	}, nil)
//
// # Credits Tracking
//
// Monitor usage and billing:
//
//	balance, err := client.GetCreditsBalance(ctx, nil)
//	history, err := client.GetCreditsHistory(ctx, nil, nil)
//	stats, err := client.GetCreditsStats(ctx, nil, nil)
//
// # Models & Capabilities
//
// Discover available models and their capabilities:
//
//	models, err := client.ListModels(ctx, nil)
//	caps, err := client.GetCapabilities(ctx, nil)
//	supportsVision := client.SupportsVision(ctx, "openai/gpt-4o", nil)
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
//
// # Error Handling
//
// All methods return structured errors that can be type-asserted for detailed information:
//
//	if err != nil {
//		if apiErr, ok := err.(*zaguansdk.APIError); ok {
//			log.Printf("API error: %d - %s", apiErr.StatusCode, apiErr.Message)
//		}
//	}
//
// # Multi-Provider Support
//
// Access 15+ AI providers through a single API by changing the model prefix:
//
//	// OpenAI
//	Model: "openai/gpt-4o"
//	// Anthropic
//	Model: "anthropic/claude-3-5-sonnet-20241022"
//	// Google
//	Model: "google/gemini-2.0-flash"
//	// DeepSeek
//	Model: "deepseek/deepseek-reasoner"
//	// And more...
//
// For complete documentation, visit https://zaguanai.com/docs
package zaguansdk
