# Zaguan Go SDK Outline

This document outlines the design and structure for the official Zaguan Go SDK. It is based on the architectural principles defined in `docs/SDK/` and tailored for idiomatic Go usage.

## 1. Overview & Goals

The Zaguan Go SDK aims to provide a robust, type-safe, and idiomatic client for interacting with Zaguan CoreX.

**Key Goals:**
*   **Idiomatic Go:** Use `context.Context`, strong typing, and proper error handling.
*   **OpenAI Compatibility:** Mirror OpenAI's API shape for easy migration.
*   **Zaguan Native:** First-class support for Zaguan features like `provider_specific_params`, routing, and credits.
*   **Performance:** Efficient streaming and memory usage.

## 2. Package Structure

We will follow a modular package layout to keep the codebase clean and maintainable.

```text
zaguan-sdk-go/
└── sdk/
    ├── client.go           // Main Client struct, constructor, and configuration
    ├── option.go           // RequestOptions and functional options pattern
    ├── chat.go             // Chat completions (req/resp types) and methods
    ├── messages.go         // Anthropic-style Messages API
    ├── models.go           // Model definitions and listing methods
    ├── capabilities.go     // Model capability discovery
    ├── credits.go          // Credits, balance, and history
    ├── errors.go           // Custom error types
    ├── stream.go           // Streaming iterator/channel handling
    └── internal/           // Internal utilities (HTTP helpers, etc.)
```

## 3. Core Configuration & Client

### Configuration (`Config`)

```go
type Config struct {
    BaseURL    string        // Defaults to standard Zaguan endpoint if empty
    APIKey     string        // Required: Bearer token
    HTTPClient *http.Client  // Optional: Defaults to http.DefaultClient
    Timeout    time.Duration // Global timeout
    Logger     Logger        // Optional interface for logging
}
```

### Client (`Client`)

The main entry point.

```go
type Client struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
    timeout    time.Duration
    logger     Logger
}

// NewClient creates a new Zaguan SDK client.
func NewClient(cfg Config) *Client
```

### Request Options (`RequestOptions`)

Allow per-request overrides.

```go
type RequestOptions struct {
    RequestID string        // X-Request-Id header
    Timeout   time.Duration // Context timeout override
    Headers   http.Header   // Extra headers
}
```

## 4. API Types

### Chat Types (`chat.go`)

Mirroring OpenAI with Zaguan extensions.

```go
type ChatRequest struct {
    Model       string            `json:"model"`
    Messages    []Message         `json:"messages"`
    Temperature *float32          `json:"temperature,omitempty"`
    MaxTokens   *int              `json:"max_tokens,omitempty"`
    TopP        *float32          `json:"top_p,omitempty"`
    Stream      bool              `json:"stream,omitempty"`
    Tools       []Tool            `json:"tools,omitempty"`
    ToolChoice  any               `json:"tool_choice,omitempty"`
    ResponseFmt any               `json:"response_format,omitempty"`

    // Zaguan Extensions
    ProviderOptions map[string]any `json:"provider_specific_params,omitempty"`
    VirtualModelID  string         `json:"virtual_model_id,omitempty"`
    Metadata        map[string]any `json:"metadata,omitempty"`
}

type Message struct {
    Role       string         `json:"role"`
    Content    any            `json:"content"` // string or []ContentPart
    Name       string         `json:"name,omitempty"`
    ToolCalls  []ToolCall     `json:"tool_calls,omitempty"`
    ToolCallID string         `json:"tool_call_id,omitempty"`
}

type ContentPart struct {
    Type      string     `json:"type"`
    Text      string     `json:"text,omitempty"`
    ImageURL  *ImageURL  `json:"image_url,omitempty"`
    InputAudio *InputAudio `json:"input_audio,omitempty"`
}

type ChatResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}

type Choice struct {
    Index        int      `json:"index"`
    Message      Message  `json:"message,omitempty"`
    FinishReason string   `json:"finish_reason,omitempty"`
    Delta        *Message `json:"delta,omitempty"` // For streaming
}

type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`

    PromptTokensDetails     *TokenDetails `json:"prompt_tokens_details,omitempty"`
    CompletionTokensDetails *TokenDetails `json:"completion_tokens_details,omitempty"`
}

type TokenDetails struct {
    ReasoningTokens int `json:"reasoning_tokens,omitempty"`
    CachedTokens    int `json:"cached_tokens,omitempty"`
    AudioTokens     int `json:"audio_tokens,omitempty"`
}

### Anthropic Messages Types (`messages.go`)

Native support for Anthropic's API shape.

```go
type MessagesRequest struct {
    Model         string                   `json:"model"`
    Messages      []AnthropicMessage       `json:"messages"`
    System        string                   `json:"system,omitempty"`
    MaxTokens     int                      `json:"max_tokens,omitempty"`
    Temperature   *float64                 `json:"temperature,omitempty"`
    TopP          *float64                 `json:"top_p,omitempty"`
    Stream        bool                     `json:"stream,omitempty"`
    StopSequences []string                 `json:"stop_sequences,omitempty"`
    Thinking      *AnthropicThinkingConfig `json:"thinking,omitempty"`
}

type AnthropicMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type AnthropicThinkingConfig struct {
    Type         string `json:"type"` // "enabled" or "disabled"
    BudgetTokens int    `json:"budget_tokens,omitempty"`
}

type MessagesResponse struct {
    ID           string                  `json:"id"`
    Type         string                  `json:"type"`
    Role         string                  `json:"role"`
    Content      []AnthropicContentBlock `json:"content"`
    Model        string                  `json:"model"`
    StopReason   string                  `json:"stop_reason,omitempty"`
    StopSequence string                  `json:"stop_sequence,omitempty"`
    Usage        AnthropicUsage          `json:"usage"`
}

type AnthropicContentBlock struct {
    Type      string `json:"type"`
    Text      string `json:"text,omitempty"`
    Thinking  string `json:"thinking,omitempty"`
    Signature string `json:"signature,omitempty"`
}

type AnthropicUsage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}
```

### Models & Capabilities (`models.go`, `capabilities.go`)

```go
type Model struct {
    ID          string         `json:"id"`
    Object      string         `json:"object"`
    OwnedBy     string         `json:"owned_by,omitempty"`
    Description string         `json:"description,omitempty"`
    Metadata    map[string]any `json:"metadata,omitempty"`
}

type ModelCapabilities struct {
    ModelID           string         `json:"model_id"`
    SupportsVision    bool           `json:"supports_vision"`
    SupportsTools     bool           `json:"supports_tools"`
    SupportsReasoning bool           `json:"supports_reasoning"`
    MaxContextTokens  int            `json:"max_context_tokens,omitempty"`
    ProviderSpecific  map[string]any `json:"provider_specific,omitempty"`
}
```

### Credits (`credits.go`)

```go
type CreditsBalance struct {
    CreditsRemaining int      `json:"credits_remaining"`
    Tier             string   `json:"tier"`
    Bands            []string `json:"bands"`
    ResetDate        *string  `json:"reset_date,omitempty"`
}
```

## 5. Methods Interface

The `Client` will expose these primary methods:

```go
// Chat
func (c *Client) Chat(ctx context.Context, req ChatRequest, opts *RequestOptions) (*ChatResponse, error)
func (c *Client) ChatStream(ctx context.Context, req ChatRequest, opts *RequestOptions) (*ChatStream, error)

// Anthropic Messages
func (c *Client) Messages(ctx context.Context, req MessagesRequest, opts *RequestOptions) (*MessagesResponse, error)
func (c *Client) MessagesStream(ctx context.Context, req MessagesRequest, opts *RequestOptions) (*MessagesStream, error)
func (c *Client) CountTokens(ctx context.Context, req MessagesRequest, opts *RequestOptions) (int, error)

// Models
func (c *Client) ListModels(ctx context.Context, opts *RequestOptions) ([]Model, error)
func (c *Client) GetCapabilities(ctx context.Context, opts *RequestOptions) ([]ModelCapabilities, error)

// Credits
func (c *Client) GetCreditsBalance(ctx context.Context, opts *RequestOptions) (*CreditsBalance, error)
// ... GetCreditsHistory, GetCreditsStats
```

## 6. Streaming Implementation

For `ChatStream`, we will return a struct that facilitates reading chunks.

```go
type ChatStream struct {
    stream *sse.Stream // Internal SSE stream handler
}

// Recv returns the next response chunk or error (io.EOF when done).
func (s *ChatStream) Recv() (*ChatResponse, error)

// Close closes the stream connection.
func (s *ChatStream) Close() error
```

## 7. Provider-Specific Helpers

To make `ProviderOptions` easier to use, we can provide helper structs/methods for common providers.

```go
// Example Helper
func WithGoogleReasoning(effort string, budget int) map[string]any {
    return map[string]any{
        "reasoning_effort": effort,
        "thinking_budget":  budget,
    }
}
```

## 8. Error Handling

We will use a structured error type to expose HTTP status codes and API messages.

```go
type APIError struct {
    StatusCode int
    Message    string
    RequestID  string
    Type       string // e.g., "insufficient_credits"
}

func (e *APIError) Error() string
```

## 9. Testing Strategy

1.  **Unit Tests:** Test struct marshaling/unmarshaling and request building.
2.  **Mock Server:** Use `httptest.Server` to simulate CoreX responses (success, errors, streaming events).
3.  **Integration Tests:** Optional suite that runs against a live CoreX instance (requires `ZAGUAN_API_KEY`).

## 10. Completed Work

### Phase 1: Foundation
- [x] Created `version.go` with SDK version constant
- [x] Created comprehensive `doc.go` with package documentation
- [x] Implemented `client.go` with Config and Client structs
- [x] Implemented `errors.go` with structured error types
- [x] Implemented `option.go` with RequestOptions
- [x] Added specialized error types (InsufficientCreditsError, BandAccessError, RateLimitError)
- [x] Added Logger interface for observability

### Documentation Created
- [x] `SDK_OUTLINE.md` - Complete SDK design document
- [x] `API_ENDPOINTS.md` - Comprehensive endpoint catalog
- [x] `IMPLEMENTATION_PLAN.md` - Phased implementation roadmap
- [x] `README.md` - SDK overview

## 11. Next Steps

### Immediate (Phase 1 Continuation)
1.  Initialize Go module with `go mod init`
2.  Implement HTTP request/response handling in `client.go`
3.  Implement `Chat()` method in `chat.go`
4.  Implement `ChatStream()` method in `stream.go`
5.  Add unit tests with `httptest.Server`

### Short-term (Phase 2)
1.  Implement Anthropic Messages API
2.  Implement token counting
3.  Add streaming support for Messages
4.  Create working examples

### Medium-term (Phase 3-4)
1.  Implement Models and Capabilities endpoints
2.  Implement Credits tracking (balance, history, stats)
3.  Add comprehensive test coverage
4.  Create migration guide from OpenAI SDK

### Long-term (Phase 5+)
1.  Implement extended OpenAI features (embeddings, audio, images)
2.  Add provider-specific helpers
3.  Implement admin endpoints
4.  Prepare for beta release
