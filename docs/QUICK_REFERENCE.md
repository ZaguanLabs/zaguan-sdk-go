# Zaguan Go SDK - Quick Reference

A concise reference for the Zaguan Go SDK's implemented features.

## Installation

```bash
go get github.com/ZaguanLabs/zaguan-sdk-go/sdk
```

## Client Initialization

```go
import zaguansdk "github.com/ZaguanLabs/zaguan-sdk-go/sdk"

client := zaguansdk.NewClient(zaguansdk.Config{
    BaseURL: "https://api.zaguanai.com",
    APIKey:  "your-api-key",
    Timeout: 30 * time.Second,  // Optional
    Logger:  myLogger,           // Optional
})
```

## Chat Completions (OpenAI-style)

### Non-Streaming

```go
resp, err := client.Chat(ctx, zaguansdk.ChatRequest{
    Model: "openai/gpt-4o",
    Messages: []zaguansdk.Message{
        {Role: "system", Content: "You are a helpful assistant."},
        {Role: "user", Content: "Hello!"},
    },
    Temperature: ptr(0.7),
    MaxTokens:   ptr(1000),
}, nil)

if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Choices[0].Message.Content)
fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
```

### Streaming

```go
stream, err := client.ChatStream(ctx, zaguansdk.ChatRequest{
    Model: "openai/gpt-4o",
    Messages: []zaguansdk.Message{
        {Role: "user", Content: "Tell me a story"},
    },
}, nil)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    event, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    if len(event.Choices) > 0 {
        fmt.Print(event.Choices[0].Delta.Content)
    }
}
```

## Anthropic Messages

### Non-Streaming

```go
resp, err := client.Messages(ctx, zaguansdk.MessagesRequest{
    Model:     "anthropic/claude-3-5-sonnet-20241022",
    MaxTokens: 1024,
    Messages: []zaguansdk.AnthropicMessage{
        {Role: "user", Content: "Hello!"},
    },
    System: "You are a helpful assistant.",
}, nil)

if err != nil {
    log.Fatal(err)
}

for _, block := range resp.Content {
    if block.Type == "text" {
        fmt.Println(block.Text)
    }
}
```

### With Extended Thinking

```go
resp, err := client.Messages(ctx, zaguansdk.MessagesRequest{
    Model:     "anthropic/claude-3-5-sonnet-20241022",
    MaxTokens: 2048,
    Messages: []zaguansdk.AnthropicMessage{
        {Role: "user", Content: "Solve this complex problem..."},
    },
    Thinking: &zaguansdk.AnthropicThinkingConfig{
        Type:         "enabled",
        BudgetTokens: 5000,
    },
}, nil)
```

### Streaming

```go
stream, err := client.MessagesStream(ctx, zaguansdk.MessagesRequest{
    Model:     "anthropic/claude-3-5-sonnet-20241022",
    MaxTokens: 1024,
    Messages: []zaguansdk.AnthropicMessage{
        {Role: "user", Content: "Tell me a story"},
    },
}, nil)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    event, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    if event.Delta != nil && event.Delta.Text != "" {
        fmt.Print(event.Delta.Text)
    }
}
```

## Models

### List All Models

```go
models, err := client.ListModels(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, model := range models {
    fmt.Printf("%s - %s\n", model.ID, model.Description)
}
```

### Get Specific Model

```go
model, err := client.GetModel(ctx, "openai/gpt-4o", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Model: %s\nOwner: %s\n", model.ID, model.OwnedBy)
```

### Delete Model

```go
err := client.DeleteModel(ctx, "ft:gpt-3.5-turbo:org:model:id", nil)
if err != nil {
    log.Fatal(err)
}
```

## Capabilities

### Get All Capabilities

```go
caps, err := client.GetCapabilities(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, cap := range caps {
    fmt.Printf("%s:\n", cap.ModelID)
    fmt.Printf("  Vision: %v\n", cap.SupportsVision)
    fmt.Printf("  Tools: %v\n", cap.SupportsTools)
    fmt.Printf("  Reasoning: %v\n", cap.SupportsReasoning)
    fmt.Printf("  Max Context: %d tokens\n", cap.MaxContextTokens)
}
```

### Get Model Capabilities

```go
cap, err := client.GetModelCapabilities(ctx, "openai/gpt-4o", nil)
if err != nil {
    log.Fatal(err)
}

if cap.SupportsVision {
    fmt.Println("This model supports vision!")
}
```

### Helper Methods

```go
// Check specific capabilities
if client.SupportsVision(ctx, "openai/gpt-4o", nil) {
    // Use vision features
}

if client.SupportsTools(ctx, "openai/gpt-4o", nil) {
    // Use tool calling
}

if client.SupportsReasoning(ctx, "openai/o1", nil) {
    // Use reasoning features
}
```

## Credits

### Get Balance

```go
balance, err := client.GetCreditsBalance(ctx, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Credits: %d/%d\n", balance.CreditsRemaining, balance.CreditsTotal)
fmt.Printf("Tier: %s\n", balance.Tier)
fmt.Printf("Percent: %.1f%%\n", balance.CreditsPercent)

if balance.IsLowCredits() {
    fmt.Println("Warning: Low credits!")
}

days, _ := balance.DaysUntilReset()
fmt.Printf("Resets in %d days\n", days)
```

### Get History

```go
history, err := client.GetCreditsHistory(ctx, &zaguansdk.CreditsHistoryOptions{
    Limit:     50,
    StartDate: "2025-01-01",
    Model:     "openai/gpt-4o",
}, nil)
if err != nil {
    log.Fatal(err)
}

for _, entry := range history.Entries {
    fmt.Printf("%s: %d credits (%d tokens)\n",
        entry.Timestamp,
        entry.CreditsDebited,
        entry.TotalTokens)
}

// Pagination
if history.HasMore {
    nextPage, _ := client.GetCreditsHistory(ctx, &zaguansdk.CreditsHistoryOptions{
        Cursor: history.NextCursor,
        Limit:  50,
    }, nil)
}
```

### Get Statistics

```go
stats, err := client.GetCreditsStats(ctx, &zaguansdk.CreditsStatsOptions{
    Period:  "month",
    GroupBy: []string{"provider", "model"},
}, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Credits: %d\n", stats.TotalCreditsUsed)
fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
fmt.Printf("Total Tokens: %d\n", stats.TotalTokens)

// By provider
for provider, pstats := range stats.ByProvider {
    fmt.Printf("%s: %d credits\n", provider, pstats.CreditsUsed)
}

// By model
for model, mstats := range stats.ByModel {
    fmt.Printf("%s: %d credits\n", model, mstats.CreditsUsed)
}
```

## Request Options

### Custom Timeout

```go
resp, err := client.Chat(ctx, req, &zaguansdk.RequestOptions{
    Timeout: 60 * time.Second,
})
```

### Custom Headers

```go
headers := http.Header{}
headers.Set("X-Custom-Header", "value")

resp, err := client.Chat(ctx, req, &zaguansdk.RequestOptions{
    Headers: headers,
})
```

### Custom Request ID

```go
resp, err := client.Chat(ctx, req, &zaguansdk.RequestOptions{
    RequestID: "my-custom-id",
})
```

## Error Handling

### Basic Error Handling

```go
resp, err := client.Chat(ctx, req, nil)
if err != nil {
    // Check for specific error types
    var apiErr *zaguansdk.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)
        fmt.Printf("Request ID: %s\n", apiErr.RequestID)
    }
    return err
}
```

### Specialized Error Types

```go
resp, err := client.Chat(ctx, req, nil)
if err != nil {
    // Check for insufficient credits
    var creditsErr *internal.InsufficientCreditsError
    if errors.As(err, &creditsErr) {
        fmt.Printf("Insufficient credits: need %d, have %d\n",
            creditsErr.CreditsRequired,
            creditsErr.CreditsRemaining)
        fmt.Printf("Resets on: %s\n", creditsErr.ResetDate)
        return err
    }
    
    // Check for band access denied
    var bandErr *internal.BandAccessError
    if errors.As(err, &bandErr) {
        fmt.Printf("Band %s requires %s tier (you have %s)\n",
            bandErr.Band,
            bandErr.RequiredTier,
            bandErr.CurrentTier)
        return err
    }
    
    // Check for rate limiting
    var rateLimitErr *internal.RateLimitError
    if errors.As(err, &rateLimitErr) {
        fmt.Printf("Rate limited. Retry after %d seconds\n",
            rateLimitErr.RetryAfter)
        time.Sleep(time.Duration(rateLimitErr.RetryAfter) * time.Second)
        // Retry request
    }
}
```

## Advanced Features

### Multimodal (Vision)

```go
resp, err := client.Chat(ctx, zaguansdk.ChatRequest{
    Model: "openai/gpt-4o",
    Messages: []zaguansdk.Message{
        {
            Role: "user",
            Content: []zaguansdk.ContentPart{
                {
                    Type: "text",
                    Text: "What's in this image?",
                },
                {
                    Type: "image_url",
                    ImageURL: &zaguansdk.ImageURL{
                        URL: "https://example.com/image.jpg",
                    },
                },
            },
        },
    },
}, nil)
```

### Tool Calling

```go
resp, err := client.Chat(ctx, zaguansdk.ChatRequest{
    Model: "openai/gpt-4o",
    Messages: []zaguansdk.Message{
        {Role: "user", Content: "What's the weather in Paris?"},
    },
    Tools: []zaguansdk.Tool{
        {
            Type: "function",
            Function: zaguansdk.FunctionDefinition{
                Name:        "get_weather",
                Description: "Get the current weather",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "location": map[string]interface{}{
                            "type":        "string",
                            "description": "City name",
                        },
                    },
                    "required": []string{"location"},
                },
            },
        },
    },
}, nil)

// Check for tool calls
if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
    for _, toolCall := range resp.Choices[0].Message.ToolCalls {
        fmt.Printf("Tool: %s\n", toolCall.Function.Name)
        fmt.Printf("Args: %s\n", toolCall.Function.Arguments)
    }
}
```

### Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Chat(ctx, req, nil)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Request timed out")
    }
}
```

### Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Log(ctx context.Context, level zaguansdk.LogLevel, msg string, keysAndValues ...interface{}) {
    fmt.Printf("[%s] %s", level, msg)
    for i := 0; i < len(keysAndValues); i += 2 {
        if i+1 < len(keysAndValues) {
            fmt.Printf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
        }
    }
    fmt.Println()
}

client := zaguansdk.NewClient(zaguansdk.Config{
    BaseURL: "https://api.zaguanai.com",
    APIKey:  "your-api-key",
    Logger:  &MyLogger{},
})
```

## Helper Functions

### Pointer Helpers

```go
func ptr[T any](v T) *T {
    return &v
}

// Usage
req := zaguansdk.ChatRequest{
    Temperature: ptr(0.7),
    MaxTokens:   ptr(1000),
    TopP:        ptr(0.9),
}
```

## Common Patterns

### Retry Logic

```go
func chatWithRetry(client *zaguansdk.Client, ctx context.Context, req zaguansdk.ChatRequest) (*zaguansdk.ChatResponse, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        resp, err := client.Chat(ctx, req, nil)
        if err == nil {
            return resp, nil
        }
        
        var rateLimitErr *internal.RateLimitError
        if errors.As(err, &rateLimitErr) {
            time.Sleep(time.Duration(rateLimitErr.RetryAfter) * time.Second)
            continue
        }
        
        return nil, err
    }
    return nil, fmt.Errorf("max retries exceeded")
}
```

### Stream to String

```go
func streamToString(stream *zaguansdk.ChatStream) (string, error) {
    var builder strings.Builder
    
    for {
        event, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", err
        }
        
        if len(event.Choices) > 0 {
            builder.WriteString(event.Choices[0].Delta.Content)
        }
    }
    
    return builder.String(), nil
}
```

## API Endpoints Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `Chat()` | POST /v1/chat/completions | OpenAI-style chat completion |
| `ChatStream()` | POST /v1/chat/completions | Streaming chat completion |
| `Messages()` | POST /v1/messages | Anthropic Messages API |
| `MessagesStream()` | POST /v1/messages | Streaming Messages API |
| `ListModels()` | GET /v1/models | List all models |
| `GetModel()` | GET /v1/models/{id} | Get model details |
| `DeleteModel()` | DELETE /v1/models/{id} | Delete model |
| `GetCapabilities()` | GET /v1/capabilities | Get all capabilities |
| `GetModelCapabilities()` | - | Get specific model capabilities |
| `GetCreditsBalance()` | GET /v1/credits/balance | Get credit balance |
| `GetCreditsHistory()` | GET /v1/credits/history | Get usage history |
| `GetCreditsStats()` | GET /v1/credits/stats | Get usage statistics |

## Environment Variables

```bash
# Recommended setup
export ZAGUAN_API_KEY="your-api-key"
export ZAGUAN_BASE_URL="https://api.zaguanai.com"  # or EU endpoint
```

```go
client := zaguansdk.NewClient(zaguansdk.Config{
    BaseURL: os.Getenv("ZAGUAN_BASE_URL"),
    APIKey:  os.Getenv("ZAGUAN_API_KEY"),
})
```

## Regional Endpoints

```go
// US endpoint (default)
baseURL := "https://api.zaguanai.com"

// EU endpoint (Finland)
baseURL := "https://api-eu-fi-01.zaguanai.com"
```

---

For more details, see:
- [README.md](../README.md) - Full documentation
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Implementation details
- [Examples](../examples/) - Working code examples
