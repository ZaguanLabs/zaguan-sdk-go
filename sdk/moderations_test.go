package zaguansdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateModeration(t *testing.T) {
	tests := []struct {
		name           string
		request        ModerationRequest
		mockResponse   ModerationResponse
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "safe content",
			request: ModerationRequest{
				Input: "Hello, how are you?",
			},
			mockResponse: ModerationResponse{
				ID:    "modr-123",
				Model: "text-moderation-latest",
				Results: []ModerationResult{
					{
						Flagged: false,
						Categories: ModerationCategories{
							Sexual:     false,
							Hate:       false,
							Harassment: false,
							SelfHarm:   false,
							Violence:   false,
						},
						CategoryScores: ModerationCategoryScores{
							Sexual:     0.001,
							Hate:       0.001,
							Harassment: 0.001,
							SelfHarm:   0.001,
							Violence:   0.001,
						},
					},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "flagged content - violence",
			request: ModerationRequest{
				Input: "I want to hurt someone",
			},
			mockResponse: ModerationResponse{
				ID:    "modr-456",
				Model: "text-moderation-latest",
				Results: []ModerationResult{
					{
						Flagged: true,
						Categories: ModerationCategories{
							Sexual:     false,
							Hate:       false,
							Harassment: false,
							SelfHarm:   false,
							Violence:   true,
						},
						CategoryScores: ModerationCategoryScores{
							Sexual:     0.001,
							Hate:       0.002,
							Harassment: 0.003,
							SelfHarm:   0.001,
							Violence:   0.95,
						},
					},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "multiple inputs",
			request: ModerationRequest{
				Input: []string{"Hello", "Goodbye"},
			},
			mockResponse: ModerationResponse{
				ID:    "modr-789",
				Model: "text-moderation-latest",
				Results: []ModerationResult{
					{
						Flagged:    false,
						Categories: ModerationCategories{},
						CategoryScores: ModerationCategoryScores{
							Sexual: 0.001,
						},
					},
					{
						Flagged:    false,
						Categories: ModerationCategories{},
						CategoryScores: ModerationCategoryScores{
							Sexual: 0.001,
						},
					},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "with specific model",
			request: ModerationRequest{
				Input: "Test content",
				Model: "text-moderation-stable",
			},
			mockResponse: ModerationResponse{
				ID:    "modr-abc",
				Model: "text-moderation-stable",
				Results: []ModerationResult{
					{
						Flagged:        false,
						Categories:     ModerationCategories{},
						CategoryScores: ModerationCategoryScores{},
					},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "missing input",
			request: ModerationRequest{
				Model: "text-moderation-latest",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name: "API error",
			request: ModerationRequest{
				Input: "Test",
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
				if r.URL.Path != "/v1/moderations" {
					t.Errorf("Expected path /v1/moderations, got %s", r.URL.Path)
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

			resp, err := client.CreateModeration(context.Background(), tt.request, nil)

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

			if len(resp.Results) != len(tt.mockResponse.Results) {
				t.Errorf("Expected %d results, got %d", len(tt.mockResponse.Results), len(resp.Results))
			}

			if len(resp.Results) > 0 {
				if resp.Results[0].Flagged != tt.mockResponse.Results[0].Flagged {
					t.Errorf("Expected flagged=%v, got %v", tt.mockResponse.Results[0].Flagged, resp.Results[0].Flagged)
				}
			}
		})
	}
}

func TestModerationResultHelpers(t *testing.T) {
	t.Run("IsSafe", func(t *testing.T) {
		tests := []struct {
			name     string
			result   ModerationResult
			wantSafe bool
		}{
			{
				name: "safe content",
				result: ModerationResult{
					Flagged: false,
				},
				wantSafe: true,
			},
			{
				name: "unsafe content",
				result: ModerationResult{
					Flagged: true,
				},
				wantSafe: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.result.IsSafe(); got != tt.wantSafe {
					t.Errorf("IsSafe() = %v, want %v", got, tt.wantSafe)
				}
			})
		}
	})

	t.Run("GetViolatedCategories", func(t *testing.T) {
		tests := []struct {
			name        string
			result      ModerationResult
			wantCount   int
			wantContain string
		}{
			{
				name: "no violations",
				result: ModerationResult{
					Categories: ModerationCategories{},
				},
				wantCount: 0,
			},
			{
				name: "single violation - violence",
				result: ModerationResult{
					Categories: ModerationCategories{
						Violence: true,
					},
				},
				wantCount:   1,
				wantContain: "violence",
			},
			{
				name: "multiple violations",
				result: ModerationResult{
					Categories: ModerationCategories{
						Violence:   true,
						Hate:       true,
						Harassment: true,
					},
				},
				wantCount: 3,
			},
			{
				name: "all categories flagged",
				result: ModerationResult{
					Categories: ModerationCategories{
						Sexual:                true,
						Hate:                  true,
						Harassment:            true,
						SelfHarm:              true,
						SexualMinors:          true,
						HateThreatening:       true,
						ViolenceGraphic:       true,
						SelfHarmIntent:        true,
						SelfHarmInstructions:  true,
						HarassmentThreatening: true,
						Violence:              true,
					},
				},
				wantCount: 11,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.result.GetViolatedCategories()
				if len(got) != tt.wantCount {
					t.Errorf("GetViolatedCategories() returned %d categories, want %d", len(got), tt.wantCount)
				}

				if tt.wantContain != "" {
					found := false
					for _, cat := range got {
						if cat == tt.wantContain {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("GetViolatedCategories() should contain %q, got %v", tt.wantContain, got)
					}
				}
			})
		}
	})
}

func TestValidateModerationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     ModerationRequest
		wantErr bool
	}{
		{
			name: "valid string input",
			req: ModerationRequest{
				Input: "test content",
			},
			wantErr: false,
		},
		{
			name: "valid array input",
			req: ModerationRequest{
				Input: []string{"test1", "test2"},
			},
			wantErr: false,
		},
		{
			name: "with model",
			req: ModerationRequest{
				Input: "test",
				Model: "text-moderation-stable",
			},
			wantErr: false,
		},
		{
			name:    "missing input",
			req:     ModerationRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModerationRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModerationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
