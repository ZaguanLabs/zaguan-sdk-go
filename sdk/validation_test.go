package zaguansdk

import (
	"testing"
)

func TestValidateChatRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ChatRequest
		wantErr bool
		errMsg  string
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
			name: "missing model",
			req: ChatRequest{
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "model is required",
		},
		{
			name: "missing messages",
			req: ChatRequest{
				Model:    "openai/gpt-4o",
				Messages: []Message{},
			},
			wantErr: true,
			errMsg:  "at least one message is required",
		},
		{
			name: "temperature too low",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: ptr(float32(-1.0)),
			},
			wantErr: true,
			errMsg:  "temperature must be between 0 and 2",
		},
		{
			name: "temperature too high",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: ptr(float32(3.0)),
			},
			wantErr: true,
			errMsg:  "temperature must be between 0 and 2",
		},
		{
			name: "valid temperature",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: ptr(float32(0.7)),
			},
			wantErr: false,
		},
		{
			name: "top_p too low",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				TopP: ptr(float32(-0.1)),
			},
			wantErr: true,
			errMsg:  "top_p must be between 0 and 1",
		},
		{
			name: "top_p too high",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				TopP: ptr(float32(1.5)),
			},
			wantErr: true,
			errMsg:  "top_p must be between 0 and 1",
		},
		{
			name: "invalid max_tokens",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: ptr(0),
			},
			wantErr: true,
			errMsg:  "max_tokens must be at least 1",
		},
		{
			name: "invalid presence_penalty",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				PresencePenalty: ptr(float32(-3.0)),
			},
			wantErr: true,
			errMsg:  "presence_penalty must be between -2 and 2",
		},
		{
			name: "invalid frequency_penalty",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				FrequencyPenalty: ptr(float32(3.0)),
			},
			wantErr: true,
			errMsg:  "frequency_penalty must be between -2 and 2",
		},
		{
			name: "invalid reasoning_effort",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				ReasoningEffort: "invalid",
			},
			wantErr: true,
			errMsg:  "reasoning_effort must be one of",
		},
		{
			name: "valid reasoning_effort",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				ReasoningEffort: "medium",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateChatRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateChatRequest() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateMessagesRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     MessagesRequest
		wantErr bool
		errMsg  string
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
			name: "missing model",
			req: MessagesRequest{
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "model is required",
		},
		{
			name: "missing messages",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages:  []AnthropicMessage{},
			},
			wantErr: true,
			errMsg:  "at least one message is required",
		},
		{
			name: "missing max_tokens",
			req: MessagesRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "max_tokens is required",
		},
		{
			name: "invalid temperature",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: ptr(1.5),
			},
			wantErr: true,
			errMsg:  "temperature must be between 0 and 1",
		},
		{
			name: "invalid thinking type",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
				Thinking: &AnthropicThinkingConfig{
					Type: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "thinking.type must be",
		},
		{
			name: "invalid thinking budget",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
				Thinking: &AnthropicThinkingConfig{
					Type:         "enabled",
					BudgetTokens: 500,
				},
			},
			wantErr: true,
			errMsg:  "thinking.budget_tokens must be between 1000 and 10000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessagesRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMessagesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateMessagesRequest() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: Config{
				BaseURL: "https://api.example.com",
				APIKey:  "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing base URL",
			cfg: Config{
				APIKey: "test-key",
			},
			wantErr: true,
			errMsg:  "BaseURL is required",
		},
		{
			name: "missing API key",
			cfg: Config{
				BaseURL: "https://api.example.com",
			},
			wantErr: true,
			errMsg:  "APIKey is required",
		},
		{
			name: "invalid base URL format",
			cfg: Config{
				BaseURL: "not-a-url",
				APIKey:  "test-key",
			},
			wantErr: true,
			errMsg:  "BaseURL must start with",
		},
		{
			name: "http URL (allowed but not recommended)",
			cfg: Config{
				BaseURL: "http://localhost:8080",
				APIKey:  "test-key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateConfig() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

// Helper functions

func ptr[T any](v T) *T {
	return &v
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
