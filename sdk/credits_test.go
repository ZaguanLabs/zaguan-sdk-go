package zaguansdk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func TestClient_GetCreditsBalance(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/credits/balance" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("unexpected method: %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"credits_remaining": 1000,
				"credits_total": 2000,
				"credits_used": 1000,
				"credits_percent": 50.0,
				"tier": "pro",
				"bands": ["A", "B", "C"],
				"reset_date": "2025-12-01T00:00:00Z"
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	balance, err := client.GetCreditsBalance(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetCreditsBalance() error = %v", err)
	}

	if balance.CreditsRemaining != 1000 {
		t.Errorf("CreditsRemaining = %d, want 1000", balance.CreditsRemaining)
	}
	if balance.Tier != "pro" {
		t.Errorf("Tier = %s, want pro", balance.Tier)
	}
	if len(balance.Bands) != 3 {
		t.Errorf("Bands length = %d, want 3", len(balance.Bands))
	}
}

func TestClient_GetCreditsHistory(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/credits/history" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			// Check query parameters
			query := r.URL.Query()
			if query.Get("limit") != "50" {
				t.Errorf("limit = %s, want 50", query.Get("limit"))
			}
			if query.Get("model") != "openai/gpt-4o" {
				t.Errorf("model = %s, want openai/gpt-4o", query.Get("model"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"entries": [
					{
						"id": "entry_123",
						"timestamp": "2025-11-19T12:00:00Z",
						"request_id": "req_123",
						"model": "openai/gpt-4o",
						"provider": "openai",
						"band": "A",
						"prompt_tokens": 10,
						"completion_tokens": 20,
						"total_tokens": 30,
						"credits_debited": 5,
						"status": "success"
					}
				],
				"total": 1,
				"has_more": false
			}`))
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	history, err := client.GetCreditsHistory(context.Background(), &CreditsHistoryOptions{
		Limit: 50,
		Model: "openai/gpt-4o",
	}, nil)
	if err != nil {
		t.Fatalf("GetCreditsHistory() error = %v", err)
	}

	if len(history.Entries) != 1 {
		t.Errorf("Entries length = %d, want 1", len(history.Entries))
	}
	if history.Total != 1 {
		t.Errorf("Total = %d, want 1", history.Total)
	}
	if history.HasMore {
		t.Error("HasMore should be false")
	}

	entry := history.Entries[0]
	if entry.Model != "openai/gpt-4o" {
		t.Errorf("Entry.Model = %s, want openai/gpt-4o", entry.Model)
	}
	if entry.CreditsDebited != 5 {
		t.Errorf("Entry.CreditsDebited = %d, want 5", entry.CreditsDebited)
	}
}

func TestClient_GetCreditsStats(t *testing.T) {
	mockServer := testutil.NewMockServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/credits/stats" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			query := r.URL.Query()
			if query.Get("period") != "month" {
				t.Errorf("period = %s, want month", query.Get("period"))
			}
			if query.Get("group_by") != "provider,model" {
				t.Errorf("group_by = %s, want provider,model", query.Get("group_by"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"period": "month",
				"total_credits_used": 1000,
				"total_requests": 100,
				"total_tokens": 50000,
				"by_provider": {
					"openai": {
						"provider": "openai",
						"credits_used": 600,
						"requests": 60,
						"tokens": 30000
					},
					"anthropic": {
						"provider": "anthropic",
						"credits_used": 400,
						"requests": 40,
						"tokens": 20000
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

	stats, err := client.GetCreditsStats(context.Background(), &CreditsStatsOptions{
		Period:  "month",
		GroupBy: []string{"provider", "model"},
	}, nil)
	if err != nil {
		t.Fatalf("GetCreditsStats() error = %v", err)
	}

	if stats.TotalCreditsUsed != 1000 {
		t.Errorf("TotalCreditsUsed = %d, want 1000", stats.TotalCreditsUsed)
	}
	if stats.TotalRequests != 100 {
		t.Errorf("TotalRequests = %d, want 100", stats.TotalRequests)
	}
	if len(stats.ByProvider) != 2 {
		t.Errorf("ByProvider length = %d, want 2", len(stats.ByProvider))
	}

	openaiStats := stats.ByProvider["openai"]
	if openaiStats.CreditsUsed != 600 {
		t.Errorf("OpenAI CreditsUsed = %d, want 600", openaiStats.CreditsUsed)
	}
}

func TestCreditsBalance_ParseResetDate(t *testing.T) {
	balance := CreditsBalance{
		ResetDate: "2025-12-01T00:00:00Z",
	}

	resetTime, err := balance.ParseResetDate()
	if err != nil {
		t.Fatalf("ParseResetDate() error = %v", err)
	}

	if resetTime.IsZero() {
		t.Error("ParseResetDate() returned zero time")
	}

	expected := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	if !resetTime.Equal(expected) {
		t.Errorf("ParseResetDate() = %v, want %v", resetTime, expected)
	}
}

func TestCreditsBalance_DaysUntilReset(t *testing.T) {
	// Set reset date to 30 days from now
	futureDate := time.Now().Add(30 * 24 * time.Hour)
	balance := CreditsBalance{
		ResetDate: futureDate.Format(time.RFC3339),
	}

	days, err := balance.DaysUntilReset()
	if err != nil {
		t.Fatalf("DaysUntilReset() error = %v", err)
	}

	// Should be approximately 30 days (allow some tolerance)
	if days < 29 || days > 31 {
		t.Errorf("DaysUntilReset() = %d, want ~30", days)
	}
}

func TestCreditsBalance_IsLowCredits(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		want    bool
	}{
		{
			name:    "low credits",
			percent: 5.0,
			want:    true,
		},
		{
			name:    "sufficient credits",
			percent: 50.0,
			want:    false,
		},
		{
			name:    "exactly at threshold",
			percent: 10.0,
			want:    false,
		},
		{
			name:    "just below threshold",
			percent: 9.9,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance := CreditsBalance{
				CreditsPercent: tt.percent,
			}

			if got := balance.IsLowCredits(); got != tt.want {
				t.Errorf("IsLowCredits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditsHistoryOptions(t *testing.T) {
	opts := CreditsHistoryOptions{
		Limit:     100,
		Cursor:    "cursor_123",
		StartDate: "2025-01-01",
		EndDate:   "2025-12-31",
		Model:     "openai/gpt-4o",
		Provider:  "openai",
		Band:      "A",
		Status:    "success",
	}

	if opts.Limit != 100 {
		t.Errorf("Limit = %d, want 100", opts.Limit)
	}
	if opts.Model != "openai/gpt-4o" {
		t.Errorf("Model = %s, want openai/gpt-4o", opts.Model)
	}
}

func TestCreditsStatsOptions(t *testing.T) {
	opts := CreditsStatsOptions{
		Period:    "month",
		StartDate: "2025-01-01",
		EndDate:   "2025-12-31",
		GroupBy:   []string{"provider", "model", "band"},
	}

	if opts.Period != "month" {
		t.Errorf("Period = %s, want month", opts.Period)
	}
	if len(opts.GroupBy) != 3 {
		t.Errorf("GroupBy length = %d, want 3", len(opts.GroupBy))
	}
}
