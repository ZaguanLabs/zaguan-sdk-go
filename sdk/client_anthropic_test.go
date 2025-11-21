package zaguansdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCountTokens(t *testing.T) {
	tests := []struct {
		name           string
		request        CountTokensRequest
		mockResponse   CountTokensResponse
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "successful token count",
			request: CountTokensRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello, world!"},
				},
			},
			mockResponse: CountTokensResponse{
				InputTokens: 10,
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "with system prompt",
			request: CountTokensRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Test"},
				},
				System: "You are a helpful assistant",
			},
			mockResponse: CountTokensResponse{
				InputTokens: 15,
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "missing model",
			request: CountTokensRequest{
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Test"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "missing messages",
			request: CountTokensRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/messages/count_tokens" {
					t.Errorf("Expected path /v1/messages/count_tokens, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
			})

			resp, err := client.CountTokens(context.Background(), tt.request, nil)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.InputTokens != tt.mockResponse.InputTokens {
				t.Errorf("Expected %d input tokens, got %d", tt.mockResponse.InputTokens, resp.InputTokens)
			}
		})
	}
}

func TestCreateMessagesBatch(t *testing.T) {
	tests := []struct {
		name           string
		request        MessagesBatchRequest
		mockResponse   MessagesBatchResponse
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "successful batch creation",
			request: MessagesBatchRequest{
				Requests: []MessagesBatchItem{
					{
						CustomID: "req-1",
						Params: MessagesRequest{
							Model:     "anthropic/claude-3-5-sonnet-20241022",
							MaxTokens: 1024,
							Messages: []AnthropicMessage{
								{Role: "user", Content: "Hello"},
							},
						},
					},
				},
			},
			mockResponse: MessagesBatchResponse{
				ID:               "msgbatch-123",
				Type:             "message_batch",
				ProcessingStatus: "in_progress",
				RequestCounts: MessagesBatchRequestCounts{
					Processing: 1,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "multiple requests",
			request: MessagesBatchRequest{
				Requests: []MessagesBatchItem{
					{
						CustomID: "req-1",
						Params: MessagesRequest{
							Model:     "anthropic/claude-3-5-sonnet-20241022",
							MaxTokens: 1024,
							Messages: []AnthropicMessage{
								{Role: "user", Content: "Hello"},
							},
						},
					},
					{
						CustomID: "req-2",
						Params: MessagesRequest{
							Model:     "anthropic/claude-3-5-sonnet-20241022",
							MaxTokens: 1024,
							Messages: []AnthropicMessage{
								{Role: "user", Content: "Goodbye"},
							},
						},
					},
				},
			},
			mockResponse: MessagesBatchResponse{
				ID:               "msgbatch-456",
				Type:             "message_batch",
				ProcessingStatus: "in_progress",
				RequestCounts: MessagesBatchRequestCounts{
					Processing: 2,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty requests",
			request: MessagesBatchRequest{
				Requests: []MessagesBatchItem{},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/messages/batches" {
					t.Errorf("Expected path /v1/messages/batches, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
			})

			resp, err := client.CreateMessagesBatch(context.Background(), tt.request, nil)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.ID != tt.mockResponse.ID {
				t.Errorf("Expected ID %s, got %s", tt.mockResponse.ID, resp.ID)
			}
		})
	}
}

func TestGetMessagesBatch(t *testing.T) {
	mockResponse := MessagesBatchResponse{
		ID:               "msgbatch-123",
		Type:             "message_batch",
		ProcessingStatus: "ended",
		RequestCounts: MessagesBatchRequestCounts{
			Succeeded: 10,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/messages/batches/msgbatch-123" {
			t.Errorf("Expected path /v1/messages/batches/msgbatch-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	resp, err := client.GetMessagesBatch(context.Background(), "msgbatch-123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != mockResponse.ID {
		t.Errorf("Expected ID %s, got %s", mockResponse.ID, resp.ID)
	}

	if resp.ProcessingStatus != mockResponse.ProcessingStatus {
		t.Errorf("Expected status %s, got %s", mockResponse.ProcessingStatus, resp.ProcessingStatus)
	}
}

func TestGetMessagesBatchEmptyID(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	_, err := client.GetMessagesBatch(context.Background(), "", nil)
	if err == nil {
		t.Error("Expected error for empty batch ID, got nil")
	}
}

func TestCancelMessagesBatch(t *testing.T) {
	mockResponse := MessagesBatchResponse{
		ID:               "msgbatch-123",
		Type:             "message_batch",
		ProcessingStatus: "canceling",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/messages/batches/msgbatch-123/cancel" {
			t.Errorf("Expected path /v1/messages/batches/msgbatch-123/cancel, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	resp, err := client.CancelMessagesBatch(context.Background(), "msgbatch-123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ProcessingStatus != "canceling" {
		t.Errorf("Expected status canceling, got %s", resp.ProcessingStatus)
	}
}

func TestCancelMessagesBatchEmptyID(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	_, err := client.CancelMessagesBatch(context.Background(), "", nil)
	if err == nil {
		t.Error("Expected error for empty batch ID, got nil")
	}
}
