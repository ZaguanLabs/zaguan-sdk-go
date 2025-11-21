package zaguansdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBatch(t *testing.T) {
	tests := []struct {
		name           string
		request        BatchRequest
		mockResponse   BatchResponse
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "successful batch creation",
			request: BatchRequest{
				InputFileID:      "file-abc123",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			mockResponse: BatchResponse{
				ID:               "batch-123",
				Object:           "batch",
				Endpoint:         "/v1/chat/completions",
				InputFileID:      "file-abc123",
				CompletionWindow: "24h",
				Status:           "validating",
				CreatedAt:        1234567890,
				RequestCounts: BatchRequestCounts{
					Total:     100,
					Completed: 0,
					Failed:    0,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "with metadata",
			request: BatchRequest{
				InputFileID:      "file-xyz",
				Endpoint:         "/v1/embeddings",
				CompletionWindow: "24h",
				Metadata: map[string]string{
					"project": "test",
					"user":    "alice",
				},
			},
			mockResponse: BatchResponse{
				ID:            "batch-456",
				Object:        "batch",
				Status:        "validating",
				CreatedAt:     1234567890,
				RequestCounts: BatchRequestCounts{},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "missing input file ID",
			request: BatchRequest{
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "missing endpoint",
			request: BatchRequest{
				InputFileID:      "file-abc",
				CompletionWindow: "24h",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "missing completion window",
			request: BatchRequest{
				InputFileID: "file-abc",
				Endpoint:    "/v1/chat/completions",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "API error",
			request: BatchRequest{
				InputFileID:      "file-abc",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/batches" {
					t.Errorf("Expected path /v1/batches, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": map[string]interface{}{
							"message": "Bad request",
						},
					})
				}
			}))
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
			})

			resp, err := client.CreateBatch(context.Background(), tt.request, nil)

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

			if resp.Status != tt.mockResponse.Status {
				t.Errorf("Expected status %s, got %s", tt.mockResponse.Status, resp.Status)
			}
		})
	}
}

func TestGetBatch(t *testing.T) {
	mockResponse := BatchResponse{
		ID:          "batch-123",
		Object:      "batch",
		Status:      "completed",
		CreatedAt:   1234567890,
		CompletedAt: 1234567990,
		RequestCounts: BatchRequestCounts{
			Total:     100,
			Completed: 100,
			Failed:    0,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/batches/batch-123" {
			t.Errorf("Expected path /v1/batches/batch-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	resp, err := client.GetBatch(context.Background(), "batch-123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != mockResponse.ID {
		t.Errorf("Expected ID %s, got %s", mockResponse.ID, resp.ID)
	}

	if resp.Status != mockResponse.Status {
		t.Errorf("Expected status %s, got %s", mockResponse.Status, resp.Status)
	}
}

func TestListBatches(t *testing.T) {
	mockResponse := BatchListResponse{
		Object: "list",
		Data: []BatchResponse{
			{
				ID:     "batch-1",
				Status: "completed",
			},
			{
				ID:     "batch-2",
				Status: "in_progress",
			},
		},
		HasMore: false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/batches" {
			t.Errorf("Expected path /v1/batches, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	resp, err := client.ListBatches(context.Background(), nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(resp.Data) != len(mockResponse.Data) {
		t.Errorf("Expected %d batches, got %d", len(mockResponse.Data), len(resp.Data))
	}
}

func TestCancelBatch(t *testing.T) {
	mockResponse := BatchResponse{
		ID:            "batch-123",
		Status:        "cancelling",
		CancellingAt:  1234567890,
		RequestCounts: BatchRequestCounts{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/batches/batch-123/cancel" {
			t.Errorf("Expected path /v1/batches/batch-123/cancel, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	resp, err := client.CancelBatch(context.Background(), "batch-123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Status != "cancelling" {
		t.Errorf("Expected status cancelling, got %s", resp.Status)
	}
}

func TestBatchResponseHelpers(t *testing.T) {
	t.Run("IsCompleted", func(t *testing.T) {
		tests := []struct {
			name   string
			status string
			want   bool
		}{
			{"completed", "completed", true},
			{"in_progress", "in_progress", false},
			{"failed", "failed", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				batch := BatchResponse{Status: tt.status}
				if got := batch.IsCompleted(); got != tt.want {
					t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("IsFailed", func(t *testing.T) {
		tests := []struct {
			name   string
			status string
			want   bool
		}{
			{"failed", "failed", true},
			{"completed", "completed", false},
			{"in_progress", "in_progress", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				batch := BatchResponse{Status: tt.status}
				if got := batch.IsFailed(); got != tt.want {
					t.Errorf("IsFailed() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("IsInProgress", func(t *testing.T) {
		tests := []struct {
			name   string
			status string
			want   bool
		}{
			{"in_progress", "in_progress", true},
			{"validating", "validating", true},
			{"finalizing", "finalizing", true},
			{"completed", "completed", false},
			{"failed", "failed", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				batch := BatchResponse{Status: tt.status}
				if got := batch.IsInProgress(); got != tt.want {
					t.Errorf("IsInProgress() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestValidateBatchRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     BatchRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: BatchRequest{
				InputFileID:      "file-123",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			wantErr: false,
		},
		{
			name: "missing input file ID",
			req: BatchRequest{
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			wantErr: true,
		},
		{
			name: "missing endpoint",
			req: BatchRequest{
				InputFileID:      "file-123",
				CompletionWindow: "24h",
			},
			wantErr: true,
		},
		{
			name: "missing completion window",
			req: BatchRequest{
				InputFileID: "file-123",
				Endpoint:    "/v1/chat/completions",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBatchRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBatchRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
