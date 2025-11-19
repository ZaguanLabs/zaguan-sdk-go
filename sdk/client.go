package zaguansdk

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	baseURL    string
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
	logger     Logger
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
	// Validate required fields
	if cfg.BaseURL == "" {
		panic("zaguansdk: BaseURL is required")
	}
	if cfg.APIKey == "" {
		panic("zaguansdk: APIKey is required")
	}

	// Use default HTTP client if none provided
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Trim trailing slash from base URL for consistency
	baseURL := strings.TrimRight(cfg.BaseURL, "/")

	return &Client{
		baseURL:    baseURL,
		apiKey:     cfg.APIKey,
		httpClient: httpClient,
		timeout:    cfg.Timeout,
		logger:     cfg.Logger,
	}
}

// BaseURL returns the base URL configured for this client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// buildURL constructs a full URL from the base URL and path.
func (c *Client) buildURL(path string) string {
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + path
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
	// Ensure stream is false for non-streaming
	req.Stream = false

	c.log(ctx, LogLevelDebug, "sending chat completion request",
		"model", req.Model,
		"message_count", len(req.Messages))

	// TODO: Implement HTTP request
	// This will be implemented using the internal HTTP client
	return nil, fmt.Errorf("not yet implemented")
}
