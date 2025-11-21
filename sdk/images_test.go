package zaguansdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateImage(t *testing.T) {
	tests := []struct {
		name           string
		request        ImageGenerationRequest
		mockResponse   ImageResponse
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "successful generation - DALL-E 3",
			request: ImageGenerationRequest{
				Prompt:  "A cute baby sea otter",
				Model:   "openai/dall-e-3",
				Size:    "1024x1024",
				Quality: "hd",
			},
			mockResponse: ImageResponse{
				Created: 1234567890,
				Data: []ImageData{
					{
						URL:           "https://example.com/image.png",
						RevisedPrompt: "A cute baby sea otter swimming",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful generation - DALL-E 2",
			request: ImageGenerationRequest{
				Prompt: "A sunset over mountains",
				Model:  "openai/dall-e-2",
				Size:   "512x512",
				N:      intPtr(2),
			},
			mockResponse: ImageResponse{
				Created: 1234567890,
				Data: []ImageData{
					{URL: "https://example.com/image1.png"},
					{URL: "https://example.com/image2.png"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "with base64 response format",
			request: ImageGenerationRequest{
				Prompt:         "Test image",
				ResponseFormat: "b64_json",
			},
			mockResponse: ImageResponse{
				Created: 1234567890,
				Data: []ImageData{
					{B64JSON: "base64encodeddata"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "missing prompt",
			request: ImageGenerationRequest{
				Model: "openai/dall-e-3",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "API error",
			request: ImageGenerationRequest{
				Prompt: "Test",
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
				if r.URL.Path != "/v1/images/generations" {
					t.Errorf("Expected path /v1/images/generations, got %s", r.URL.Path)
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

			resp, err := client.CreateImage(context.Background(), tt.request, nil)

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

			if len(resp.Data) != len(tt.mockResponse.Data) {
				t.Errorf("Expected %d images, got %d", len(tt.mockResponse.Data), len(resp.Data))
			}

			if len(resp.Data) > 0 {
				if resp.Data[0].URL != tt.mockResponse.Data[0].URL {
					t.Errorf("Expected URL %s, got %s", tt.mockResponse.Data[0].URL, resp.Data[0].URL)
				}
			}
		})
	}
}

func TestEditImage(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	req := ImageEditRequest{
		Image:  "test.png",
		Prompt: "Add a hat",
	}

	_, err := client.EditImage(context.Background(), req, nil)
	if err == nil {
		t.Error("Expected not implemented error, got nil")
	}

	// Check that it's an API error with status 501
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.StatusCode != 501 {
			t.Errorf("Expected status code 501, got %d", apiErr.StatusCode)
		}
	} else {
		t.Error("Expected APIError type")
	}
}

func TestCreateImageVariation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	req := ImageVariationRequest{
		Image: "test.png",
	}

	_, err := client.CreateImageVariation(context.Background(), req, nil)
	if err == nil {
		t.Error("Expected not implemented error, got nil")
	}

	// Check that it's an API error with status 501
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.StatusCode != 501 {
			t.Errorf("Expected status code 501, got %d", apiErr.StatusCode)
		}
	} else {
		t.Error("Expected APIError type")
	}
}

func TestValidateImageGenerationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ImageGenerationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ImageGenerationRequest{
				Prompt: "A test image",
			},
			wantErr: false,
		},
		{
			name: "with all parameters",
			req: ImageGenerationRequest{
				Prompt:         "Test",
				Model:          "dall-e-3",
				Size:           "1024x1024",
				Quality:        "hd",
				Style:          "vivid",
				ResponseFormat: "url",
				N:              intPtr(1),
			},
			wantErr: false,
		},
		{
			name:    "missing prompt",
			req:     ImageGenerationRequest{},
			wantErr: true,
		},
		{
			name: "empty prompt",
			req: ImageGenerationRequest{
				Prompt: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageGenerationRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateImageGenerationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImageEditRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ImageEditRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ImageEditRequest{
				Image:  "test.png",
				Prompt: "Edit this",
			},
			wantErr: false,
		},
		{
			name: "missing image",
			req: ImageEditRequest{
				Prompt: "Edit this",
			},
			wantErr: true,
		},
		{
			name: "missing prompt",
			req: ImageEditRequest{
				Image: "test.png",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageEditRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateImageEditRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImageVariationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ImageVariationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: ImageVariationRequest{
				Image: "test.png",
			},
			wantErr: false,
		},
		{
			name:    "missing image",
			req:     ImageVariationRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageVariationRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateImageVariationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function for tests
func intPtr(i int) *int {
	return &i
}
