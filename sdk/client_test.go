package zaguansdk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantPanic bool
	}{
		{
			name: "valid config",
			cfg: Config{
				BaseURL: "https://api.example.com",
				APIKey:  "test-key",
			},
			wantPanic: false,
		},
		{
			name: "missing base URL",
			cfg: Config{
				APIKey: "test-key",
			},
			wantPanic: true,
		},
		{
			name: "missing API key",
			cfg: Config{
				BaseURL: "https://api.example.com",
			},
			wantPanic: true,
		},
		{
			name: "with custom HTTP client",
			cfg: Config{
				BaseURL:    "https://api.example.com",
				APIKey:     "test-key",
				HTTPClient: &http.Client{Timeout: 10 * time.Second},
			},
			wantPanic: false,
		},
		{
			name: "with timeout",
			cfg: Config{
				BaseURL: "https://api.example.com",
				APIKey:  "test-key",
				Timeout: 30 * time.Second,
			},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("NewClient() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			client := NewClient(tt.cfg)
			if !tt.wantPanic && client == nil {
				t.Error("NewClient() returned nil")
			}
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com/",
		APIKey:  "test-key",
	})

	// Should trim trailing slash
	if got := client.BaseURL(); got != "https://api.example.com" {
		t.Errorf("BaseURL() = %v, want %v", got, "https://api.example.com")
	}
}

func TestClient_Chat(t *testing.T) {
	tests := []struct {
		name    string
		req     ChatRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing model",
			req: ChatRequest{
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing messages",
			req: ChatRequest{
				Model:    "openai/gpt-4o",
				Messages: []Message{},
			},
			wantErr: true,
		},
		{
			name: "invalid request - bad temperature",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: ptr(float32(5.0)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := testutil.NewMockServer(
				testutil.ChatCompletionHandler(testutil.ChatCompletionFixture()),
			)
			defer mockServer.Close()

			// Create client
			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			// Make request
			resp, err := client.Chat(context.Background(), tt.req, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Chat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("Chat() returned nil response")
					return
				}
				if resp.ID == "" {
					t.Error("Chat() response missing ID")
				}
			}
		})
	}
}

func TestClient_ChatWithTimeout(t *testing.T) {
	// Create mock server that delays
	mockServer := testutil.NewMockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		testutil.ChatCompletionHandler(testutil.ChatCompletionFixture())(w, r)
	}))
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	// Request with very short timeout should fail
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Chat(ctx, ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}, nil)

	if err == nil {
		t.Error("Chat() with timeout should have failed")
	}
}

func TestClient_Messages(t *testing.T) {
	tests := []struct {
		name    string
		req     MessagesRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing model",
			req: MessagesRequest{
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing max_tokens",
			req: MessagesRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := testutil.NewMockServer(
				testutil.MessagesHandler(testutil.MessagesFixture()),
			)
			defer mockServer.Close()

			// Create client
			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			// Make request
			resp, err := client.Messages(context.Background(), tt.req, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Messages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("Messages() returned nil response")
					return
				}
				if resp.ID == "" {
					t.Error("Messages() response missing ID")
				}
			}
		})
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorType  string
		message    string
	}{
		{
			name:       "bad request",
			statusCode: 400,
			errorType:  "invalid_request",
			message:    "Invalid request",
		},
		{
			name:       "unauthorized",
			statusCode: 401,
			errorType:  "authentication_error",
			message:    "Invalid API key",
		},
		{
			name:       "rate limit",
			statusCode: 429,
			errorType:  "rate_limit_exceeded",
			message:    "Rate limit exceeded",
		},
		{
			name:       "server error",
			statusCode: 500,
			errorType:  "server_error",
			message:    "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server that returns error
			mockServer := testutil.NewMockServer(
				testutil.ErrorHandler(tt.statusCode, tt.errorType, tt.message),
			)
			defer mockServer.Close()

			// Create client
			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			// Make request
			_, err := client.Chat(context.Background(), ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			}, nil)

			if err == nil {
				t.Error("Chat() should have returned error")
				return
			}

			// Error should contain status code information
			// The actual error type may be from internal package
			if err.Error() == "" {
				t.Error("Chat() returned empty error message")
			}
		})
	}
}
