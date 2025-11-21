package zaguansdk

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// Config holds the configuration for the Zaguan SDK client.
//
// At minimum, you must provide BaseURL and APIKey. Other fields are optional
// and will use sensible defaults if not specified.
type Config struct {
	// BaseURL is the base URL of your Zaguan CoreX instance.
	// Example: "https://api.zaguan.example.com"
	// Required.
	BaseURL string

	// APIKey is your Zaguan API key for authentication.
	// This will be sent as a Bearer token in the Authorization header.
	// Required.
	APIKey string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, http.DefaultClient will be used.
	// Optional.
	HTTPClient *http.Client

	// Timeout is the default timeout for all requests.
	// Individual requests can override this via RequestOptions.
	// If zero, no timeout is applied at the client level.
	// Optional.
	Timeout time.Duration

	// Logger is an optional logger for debugging and observability.
	// If nil, no logging will be performed.
	// Optional.
	Logger Logger
}

// Client is the main entry point for interacting with Zaguan CoreX.
//
// A Client is safe for concurrent use by multiple goroutines.
// You should create a single Client and reuse it throughout your application.
type Client struct {
	baseURL      string
	apiKey       string
	httpClient   *http.Client
	internalHTTP *internal.HTTPClient
	timeout      time.Duration
	logger       Logger
}

// NewClient creates a new Zaguan SDK client with the provided configuration.
//
// The client is safe for concurrent use and should be reused across requests.
//
// Example:
//
//	client := zaguansdk.NewClient(zaguansdk.Config{
//		BaseURL: "https://api.zaguan.example.com",
//		APIKey:  "your-api-key",
//	})
func NewClient(cfg Config) *Client {
	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		panic(fmt.Sprintf("zaguansdk: invalid configuration: %v", err))
	}

	// Use default HTTP client if none provided
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Trim trailing slash from base URL for consistency
	baseURL := strings.TrimRight(cfg.BaseURL, "/")

	// Create internal HTTP client
	internalHTTP := internal.NewHTTPClient(httpClient, baseURL, cfg.APIKey, Version)

	return &Client{
		baseURL:      baseURL,
		apiKey:       cfg.APIKey,
		httpClient:   httpClient,
		internalHTTP: internalHTTP,
		timeout:      cfg.Timeout,
		logger:       cfg.Logger,
	}
}

// BaseURL returns the base URL configured for this client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// log logs a message if a logger is configured.
func (c *Client) log(ctx context.Context, level LogLevel, msg string, keysAndValues ...interface{}) {
	if c.logger != nil {
		c.logger.Log(ctx, level, msg, keysAndValues...)
	}
}

// Logger is an interface for logging within the SDK.
//
// You can provide your own implementation to integrate with your logging framework.
// The SDK will never log sensitive information like API keys.
type Logger interface {
	// Log logs a message at the specified level with optional key-value pairs.
	Log(ctx context.Context, level LogLevel, msg string, keysAndValues ...interface{})
}

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	// LogLevelDebug is for detailed debugging information.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for general informational messages.
	LogLevelInfo
	// LogLevelWarn is for warning messages.
	LogLevelWarn
	// LogLevelError is for error messages.
	LogLevelError
)

// String returns the string representation of a log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", l)
	}
}

// Chat sends a chat completion request to Zaguan CoreX.
//
// This method supports all OpenAI chat completion parameters plus Zaguan extensions.
// For streaming responses, use ChatStream instead.
//
// Example:
//
//	resp, err := client.Chat(ctx, zaguansdk.ChatRequest{
//		Model: "openai/gpt-4o",
//		Messages: []zaguansdk.Message{
//			{Role: "user", Content: "Hello!"},
//		},
//	}, nil)
func (c *Client) Chat(ctx context.Context, req ChatRequest, opts *RequestOptions) (*ChatResponse, error) {
	// Validate request
	if err := validateChatRequest(&req); err != nil {
		return nil, err
	}

	// Ensure stream is false for non-streaming
	req.Stream = false

	c.log(ctx, LogLevelDebug, "sending chat completion request",
		"model", req.Model,
		"message_count", len(req.Messages))

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/chat/completions",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp ChatResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "chat completion request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "chat completion request succeeded",
		"response_id", resp.ID,
		"model", resp.Model)

	return &resp, nil
}

// Messages sends a request to Anthropic's native Messages API.
//
// This method supports all Anthropic Messages API parameters including extended thinking.
// For streaming responses, use MessagesStream instead.
//
// Example:
//
//	resp, err := client.Messages(ctx, zaguansdk.MessagesRequest{
//		Model: "anthropic/claude-3-5-sonnet-20241022",
//		MaxTokens: 1024,
//		Messages: []zaguansdk.AnthropicMessage{
//			{Role: "user", Content: "Hello!"},
//		},
//	}, nil)
func (c *Client) Messages(ctx context.Context, req MessagesRequest, opts *RequestOptions) (*MessagesResponse, error) {
	// Validate request
	if err := validateMessagesRequest(&req); err != nil {
		return nil, err
	}

	// Ensure stream is false for non-streaming
	req.Stream = false

	c.log(ctx, LogLevelDebug, "sending messages request",
		"model", req.Model,
		"message_count", len(req.Messages))

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/messages",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp MessagesResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "messages request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "messages request succeeded",
		"response_id", resp.ID,
		"model", resp.Model)

	return &resp, nil
}

// CountTokens counts the number of tokens in a Messages request.
//
// This is useful for estimating costs before making a request.
//
// Example:
//
//	resp, err := client.CountTokens(ctx, zaguansdk.CountTokensRequest{
//		Model: "anthropic/claude-3-5-sonnet-20241022",
//		Messages: []zaguansdk.AnthropicMessage{
//			{Role: "user", Content: "Hello, world!"},
//		},
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Input tokens: %d\n", resp.InputTokens)
func (c *Client) CountTokens(ctx context.Context, req CountTokensRequest, opts *RequestOptions) (*CountTokensResponse, error) {
	// Validate request
	if req.Model == "" {
		return nil, &ValidationError{Field: "model", Message: "model is required"}
	}
	if len(req.Messages) == 0 {
		return nil, &ValidationError{Field: "messages", Message: "at least one message is required"}
	}

	c.log(ctx, LogLevelDebug, "counting tokens", "model", req.Model)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/messages/count_tokens",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp CountTokensResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "count tokens request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "count tokens request succeeded", "input_tokens", resp.InputTokens)

	return &resp, nil
}

// CreateMessagesBatch creates a batch of Anthropic Messages requests.
//
// Example:
//
//	resp, err := client.CreateMessagesBatch(ctx, zaguansdk.MessagesBatchRequest{
//		Requests: []zaguansdk.MessagesBatchItem{
//			{
//				CustomID: "request-1",
//				Params: zaguansdk.MessagesRequest{
//					Model:     "anthropic/claude-3-5-sonnet-20241022",
//					MaxTokens: 1024,
//					Messages: []zaguansdk.AnthropicMessage{
//						{Role: "user", Content: "Hello!"},
//					},
//				},
//			},
//		},
//	}, nil)
func (c *Client) CreateMessagesBatch(ctx context.Context, req MessagesBatchRequest, opts *RequestOptions) (*MessagesBatchResponse, error) {
	// Validate request
	if len(req.Requests) == 0 {
		return nil, &ValidationError{Field: "requests", Message: "at least one request is required"}
	}

	c.log(ctx, LogLevelDebug, "creating messages batch", "count", len(req.Requests))

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/messages/batches",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp MessagesBatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create messages batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create messages batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}

// GetMessagesBatch retrieves information about a Messages batch.
//
// Example:
//
//	batch, err := client.GetMessagesBatch(ctx, "msgbatch_abc123", nil)
func (c *Client) GetMessagesBatch(ctx context.Context, batchID string, opts *RequestOptions) (*MessagesBatchResponse, error) {
	if batchID == "" {
		return nil, &ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}

	c.log(ctx, LogLevelDebug, "getting messages batch", "batch_id", batchID)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "GET",
		Path:   "/v1/messages/batches/" + batchID,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp MessagesBatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "get messages batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "get messages batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}

// CancelMessagesBatch cancels a Messages batch that is in progress.
//
// Example:
//
//	batch, err := client.CancelMessagesBatch(ctx, "msgbatch_abc123", nil)
func (c *Client) CancelMessagesBatch(ctx context.Context, batchID string, opts *RequestOptions) (*MessagesBatchResponse, error) {
	if batchID == "" {
		return nil, &ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}

	c.log(ctx, LogLevelDebug, "cancelling messages batch", "batch_id", batchID)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   fmt.Sprintf("/v1/messages/batches/%s/cancel", batchID),
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp MessagesBatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "cancel messages batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "cancel messages batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}
