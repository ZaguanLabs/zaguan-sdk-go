package zaguansdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEmbeddings(t *testing.T) {
	tests := []struct {
		name           string
		request        EmbeddingsRequest
		mockResponse   EmbeddingsResponse
		mockStatusCode int
		wantErr        bool
		errType        string
	}{
		{
			name: "successful single text embedding",
			request: EmbeddingsRequest{
				Model: "openai/text-embedding-3-small",
				Input: "Hello, world!",
			},
			mockResponse: EmbeddingsResponse{
				Object: "list",
				Model:  "text-embedding-3-small",
				Data: []Embedding{
					{
						Object:    "embedding",
						Embedding: []interface{}{0.1, 0.2, 0.3},
						Index:     0,
					},
				},
				Usage: EmbeddingsUsage{
					PromptTokens: 3,
					TotalTokens:  3,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful multiple text embeddings",
			request: EmbeddingsRequest{
				Model: "openai/text-embedding-3-small",
				Input: []string{"Hello", "World"},
			},
			mockResponse: EmbeddingsResponse{
				Object: "list",
				Model:  "text-embedding-3-small",
				Data: []Embedding{
					{
						Object:    "embedding",
						Embedding: []interface{}{0.1, 0.2, 0.3},
						Index:     0,
					},
					{
						Object:    "embedding",
						Embedding: []interface{}{0.4, 0.5, 0.6},
						Index:     1,
					},
				},
				Usage: EmbeddingsUsage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "with dimensions parameter",
			request: EmbeddingsRequest{
				Model:      "openai/text-embedding-3-small",
				Input:      "Hello",
				Dimensions: 256,
			},
			mockResponse: EmbeddingsResponse{
				Object: "list",
				Model:  "text-embedding-3-small",
				Data: []Embedding{
					{
						Object:    "embedding",
						Embedding: make([]interface{}, 256),
						Index:     0,
					},
				},
				Usage: EmbeddingsUsage{
					PromptTokens: 1,
					TotalTokens:  1,
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "missing model",
			request: EmbeddingsRequest{
				Input: "Hello",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
			errType:        "validation",
		},
		{
			name: "missing input",
			request: EmbeddingsRequest{
				Model: "openai/text-embedding-3-small",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
			errType:        "validation",
		},
		{
			name: "empty string input",
			request: EmbeddingsRequest{
				Model: "openai/text-embedding-3-small",
				Input: "",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
			errType:        "validation",
		},
		{
			name: "API error",
			request: EmbeddingsRequest{
				Model: "openai/text-embedding-3-small",
				Input: "Hello",
			},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
			errType:        "api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/embeddings" {
					t.Errorf("Expected path /v1/embeddings, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": map[string]interface{}{
							"message": "Bad request",
							"type":    "invalid_request_error",
						},
					})
				}
			}))
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
			})

			resp, err := client.CreateEmbeddings(context.Background(), tt.request, nil)

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

			if resp.Model != tt.mockResponse.Model {
				t.Errorf("Expected model %s, got %s", tt.mockResponse.Model, resp.Model)
			}

			if len(resp.Data) != len(tt.mockResponse.Data) {
				t.Errorf("Expected %d embeddings, got %d", len(tt.mockResponse.Data), len(resp.Data))
			}
		})
	}
}

func TestGetEmbeddingVector(t *testing.T) {
	tests := []struct {
		name      string
		embedding Embedding
		wantLen   int
		wantErr   bool
	}{
		{
			name: "valid float vector",
			embedding: Embedding{
				Embedding: []interface{}{0.1, 0.2, 0.3, 0.4, 0.5},
			},
			wantLen: 5,
			wantErr: false,
		},
		{
			name: "invalid type",
			embedding: Embedding{
				Embedding: "not a vector",
			},
			wantErr: true,
		},
		{
			name: "mixed types in vector",
			embedding: Embedding{
				Embedding: []interface{}{0.1, "invalid", 0.3},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vec, err := tt.embedding.GetEmbeddingVector()

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

			if len(vec) != tt.wantLen {
				t.Errorf("Expected vector length %d, got %d", tt.wantLen, len(vec))
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name    string
		a       []float64
		b       []float64
		want    float64
		wantErr bool
	}{
		{
			name: "identical vectors",
			a:    []float64{1.0, 0.0, 0.0},
			b:    []float64{1.0, 0.0, 0.0},
			want: 1.0,
		},
		{
			name: "orthogonal vectors",
			a:    []float64{1.0, 0.0},
			b:    []float64{0.0, 1.0},
			want: 0.0,
		},
		{
			name: "opposite vectors",
			a:    []float64{1.0, 0.0},
			b:    []float64{-1.0, 0.0},
			want: -1.0,
		},
		{
			name:    "different lengths",
			a:       []float64{1.0, 0.0},
			b:       []float64{1.0, 0.0, 0.0},
			wantErr: true,
		},
		{
			name:    "zero vector",
			a:       []float64{0.0, 0.0},
			b:       []float64{1.0, 0.0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CosineSimilarity(tt.a, tt.b)

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

			// Allow small floating point differences
			if abs(got-tt.want) > 0.0001 {
				t.Errorf("Expected similarity %f, got %f", tt.want, got)
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestValidateEmbeddingsRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     EmbeddingsRequest
		wantErr bool
	}{
		{
			name: "valid string input",
			req: EmbeddingsRequest{
				Model: "test-model",
				Input: "test input",
			},
			wantErr: false,
		},
		{
			name: "valid array input",
			req: EmbeddingsRequest{
				Model: "test-model",
				Input: []string{"test1", "test2"},
			},
			wantErr: false,
		},
		{
			name: "missing model",
			req: EmbeddingsRequest{
				Input: "test",
			},
			wantErr: true,
		},
		{
			name: "missing input",
			req: EmbeddingsRequest{
				Model: "test-model",
			},
			wantErr: true,
		},
		{
			name: "empty string input",
			req: EmbeddingsRequest{
				Model: "test-model",
				Input: "",
			},
			wantErr: true,
		},
		{
			name: "empty array input",
			req: EmbeddingsRequest{
				Model: "test-model",
				Input: []string{},
			},
			wantErr: true,
		},
		{
			name: "invalid encoding format",
			req: EmbeddingsRequest{
				Model:          "test-model",
				Input:          "test",
				EncodingFormat: "invalid",
			},
			wantErr: true,
		},
		{
			name: "negative dimensions",
			req: EmbeddingsRequest{
				Model:      "test-model",
				Input:      "test",
				Dimensions: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmbeddingsRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmbeddingsRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
