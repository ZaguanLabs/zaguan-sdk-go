package zaguansdk

import (
	"context"
	"net/http"
	"testing"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func TestClient_ListModels(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/models" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("unexpected method: %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"object": "list",
				"data": [
					{
						"id": "openai/gpt-4o",
						"object": "model",
						"created": 1677652288,
						"owned_by": "openai"
					},
					{
						"id": "anthropic/claude-3-5-sonnet-20241022",
						"object": "model",
						"created": 1677652288,
						"owned_by": "anthropic"
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

	models, err := client.ListModels(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}

	if len(models) != 2 {
		t.Errorf("ListModels() returned %d models, want 2", len(models))
	}

	if models[0].ID != "openai/gpt-4o" {
		t.Errorf("first model ID = %s, want openai/gpt-4o", models[0].ID)
	}
}

func TestClient_GetModel(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		wantErr bool
	}{
		{
			name:    "valid model ID",
			modelID: "openai/gpt-4o",
			wantErr: false,
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
						"id": "openai/gpt-4o",
						"object": "model",
						"created": 1677652288,
						"owned_by": "openai"
					}`))
				}),
			)
			defer mockServer.Close()

			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			model, err := client.GetModel(context.Background(), tt.modelID, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if model == nil {
					t.Error("GetModel() returned nil model")
					return
				}
				if model.ID == "" {
					t.Error("GetModel() returned model with empty ID")
				}
			}
		})
	}
}

func TestClient_DeleteModel(t *testing.T) {
	tests := []struct {
		name       string
		modelID    string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			modelID:    "ft:gpt-3.5-turbo:org:model:id",
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "model not found",
			modelID:    "nonexistent",
			statusCode: 404,
			wantErr:    true,
		},
		{
			name:       "empty model ID",
			modelID:    "",
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := testutil.NewMockServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodDelete {
						t.Errorf("unexpected method: %s", r.Method)
					}

					w.WriteHeader(tt.statusCode)
					if tt.statusCode >= 400 {
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(`{"error": {"message": "Not found", "type": "not_found"}}`))
					}
				}),
			)
			defer mockServer.Close()

			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			err := client.DeleteModel(context.Background(), tt.modelID, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModel_Fields(t *testing.T) {
	model := Model{
		ID:          "openai/gpt-4o",
		Object:      "model",
		Created:     1677652288,
		OwnedBy:     "openai",
		Description: "GPT-4o model",
	}

	if model.ID != "openai/gpt-4o" {
		t.Errorf("Model.ID = %s, want openai/gpt-4o", model.ID)
	}
	if model.Object != "model" {
		t.Errorf("Model.Object = %s, want model", model.Object)
	}
	if model.OwnedBy != "openai" {
		t.Errorf("Model.OwnedBy = %s, want openai", model.OwnedBy)
	}
}
