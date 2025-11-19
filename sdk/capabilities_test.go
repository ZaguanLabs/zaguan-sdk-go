package zaguansdk

import (
	"context"
	"net/http"
	"testing"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func TestClient_GetCapabilities(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/capabilities" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("unexpected method: %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"models": [
					{
						"model_id": "openai/gpt-4o",
						"provider": "openai",
						"supports_vision": true,
						"supports_tools": true,
						"supports_reasoning": false,
						"max_context_tokens": 128000,
						"max_output_tokens": 4096
					},
					{
						"model_id": "anthropic/claude-3-5-sonnet-20241022",
						"provider": "anthropic",
						"supports_vision": true,
						"supports_tools": true,
						"supports_reasoning": true,
						"max_context_tokens": 200000,
						"max_output_tokens": 8192
					}
				]
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	caps, err := client.GetCapabilities(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetCapabilities() error = %v", err)
	}

	if len(caps) != 2 {
		t.Errorf("GetCapabilities() returned %d capabilities, want 2", len(caps))
	}

	// Check first capability
	if caps[0].ModelID != "openai/gpt-4o" {
		t.Errorf("first capability ModelID = %s, want openai/gpt-4o", caps[0].ModelID)
	}
	if !caps[0].SupportsVision {
		t.Error("gpt-4o should support vision")
	}
	if !caps[0].SupportsTools {
		t.Error("gpt-4o should support tools")
	}
	if caps[0].SupportsReasoning {
		t.Error("gpt-4o should not support reasoning")
	}
}

func TestClient_GetModelCapabilities(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		wantErr bool
	}{
		{
			name:    "existing model",
			modelID: "openai/gpt-4o",
			wantErr: false,
		},
		{
			name:    "non-existent model",
			modelID: "nonexistent/model",
			wantErr: true,
		},
		{
			name:    "empty model ID",
			modelID: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutil.NewMockServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"models": [
							{
								"model_id": "openai/gpt-4o",
								"provider": "openai",
								"supports_vision": true,
								"supports_tools": true,
								"supports_reasoning": false
							}
						]
					}`))
				}),
			)
			defer mockServer.Close()

			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			cap, err := client.GetModelCapabilities(context.Background(), tt.modelID, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetModelCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cap == nil {
					t.Error("GetModelCapabilities() returned nil")
					return
				}
				if cap.ModelID != tt.modelID {
					t.Errorf("capability ModelID = %s, want %s", cap.ModelID, tt.modelID)
				}
			}
		})
	}
}

func TestClient_SupportsVision(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"models": [
					{
						"model_id": "openai/gpt-4o",
						"supports_vision": true
					}
				]
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	if !client.SupportsVision(context.Background(), "openai/gpt-4o", nil) {
		t.Error("SupportsVision() should return true for gpt-4o")
	}

	if client.SupportsVision(context.Background(), "nonexistent", nil) {
		t.Error("SupportsVision() should return false for nonexistent model")
	}
}

func TestClient_SupportsTools(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"models": [
					{
						"model_id": "openai/gpt-4o",
						"supports_tools": true
					}
				]
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	if !client.SupportsTools(context.Background(), "openai/gpt-4o", nil) {
		t.Error("SupportsTools() should return true for gpt-4o")
	}
}

func TestClient_SupportsReasoning(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"models": [
					{
						"model_id": "anthropic/claude-3-5-sonnet-20241022",
						"supports_reasoning": true
					}
				]
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	if !client.SupportsReasoning(context.Background(), "anthropic/claude-3-5-sonnet-20241022", nil) {
		t.Error("SupportsReasoning() should return true for claude")
	}
}

func TestCapabilitiesResponse_MapFormat(t *testing.T) {
	// Test that we can handle the map format response
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"capabilities": {
					"openai/gpt-4o": {
						"model_id": "openai/gpt-4o",
						"supports_vision": true
					}
				}
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	caps, err := client.GetCapabilities(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetCapabilities() error = %v", err)
	}

	if len(caps) != 1 {
		t.Errorf("GetCapabilities() returned %d capabilities, want 1", len(caps))
	}
}
