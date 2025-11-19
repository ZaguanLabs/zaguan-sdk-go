package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantType   string
	}{
		{
			name:       "insufficient credits error",
			statusCode: 402,
			body:       `{"error": {"type": "insufficient_credits", "message": "Not enough credits", "credits_required": 100, "credits_remaining": 50}}`,
			wantType:   "*internal.InsufficientCreditsError",
		},
		{
			name:       "band access error",
			statusCode: 403,
			body:       `{"error": {"type": "band_access_denied", "message": "Access denied", "band": "premium", "required_tier": "pro"}}`,
			wantType:   "*internal.BandAccessError",
		},
		{
			name:       "rate limit error",
			statusCode: 429,
			body:       `{"error": {"type": "rate_limit_exceeded", "message": "Rate limit exceeded"}}`,
			wantType:   "*internal.RateLimitError",
		},
		{
			name:       "generic API error",
			statusCode: 400,
			body:       `{"error": {"type": "invalid_request", "message": "Bad request"}}`,
			wantType:   "*internal.APIError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
				Header:     http.Header{},
			}

			err := ParseErrorResponse(resp)
			if err == nil {
				t.Fatal("ParseErrorResponse() returned nil error")
			}

			// Check error type by type assertion
			switch tt.wantType {
			case "*internal.InsufficientCreditsError":
				if _, ok := err.(*InsufficientCreditsError); !ok {
					t.Errorf("error type = %T, want %s", err, tt.wantType)
				}
			case "*internal.BandAccessError":
				if _, ok := err.(*BandAccessError); !ok {
					t.Errorf("error type = %T, want %s", err, tt.wantType)
				}
			case "*internal.RateLimitError":
				if _, ok := err.(*RateLimitError); !ok {
					t.Errorf("error type = %T, want %s", err, tt.wantType)
				}
			case "*internal.APIError":
				if _, ok := err.(*APIError); !ok {
					t.Errorf("error type = %T, want %s", err, tt.wantType)
				}
			}
		})
	}
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
		Header:     http.Header{},
	}

	err := ParseErrorResponse(resp)
	if err == nil {
		t.Fatal("ParseErrorResponse() should return error for invalid JSON")
	}
}

func TestParseInsufficientCreditsError(t *testing.T) {
	body := `{"error": {"type": "insufficient_credits", "message": "Not enough credits", "details": {"credits_required": 100, "credits_remaining": 50, "reset_date": "2025-12-01"}}}`
	resp := &http.Response{
		StatusCode: 402,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     http.Header{},
	}

	err := ParseErrorResponse(resp)
	if err == nil {
		t.Fatal("ParseErrorResponse() returned nil")
	}

	credErr, ok := err.(*InsufficientCreditsError)
	if !ok {
		t.Fatalf("error type = %T, want *InsufficientCreditsError", err)
	}

	if credErr.CreditsRequired != 100 {
		t.Errorf("CreditsRequired = %d, want 100", credErr.CreditsRequired)
	}
	if credErr.CreditsRemaining != 50 {
		t.Errorf("CreditsRemaining = %d, want 50", credErr.CreditsRemaining)
	}
}

func TestParseBandAccessError(t *testing.T) {
	body := `{"error": {"type": "band_access_denied", "message": "Access denied", "details": {"band": "premium", "required_tier": "pro", "current_tier": "free"}}}`
	resp := &http.Response{
		StatusCode: 403,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     http.Header{},
	}

	err := ParseErrorResponse(resp)
	if err == nil {
		t.Fatal("ParseErrorResponse() returned nil")
	}

	bandErr, ok := err.(*BandAccessError)
	if !ok {
		t.Fatalf("error type = %T, want *BandAccessError", err)
	}

	if bandErr.Band != "premium" {
		t.Errorf("Band = %s, want premium", bandErr.Band)
	}
	if bandErr.RequiredTier != "pro" {
		t.Errorf("RequiredTier = %s, want pro", bandErr.RequiredTier)
	}
	if bandErr.CurrentTier != "free" {
		t.Errorf("CurrentTier = %s, want free", bandErr.CurrentTier)
	}
}

func TestParseRateLimitError(t *testing.T) {
	body := `{"error": {"type": "rate_limit_exceeded", "message": "Rate limit exceeded"}}`
	resp := &http.Response{
		StatusCode: 429,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header: http.Header{
			"Retry-After": []string{"60"},
		},
	}

	err := ParseErrorResponse(resp)
	if err == nil {
		t.Fatal("ParseErrorResponse() returned nil")
	}

	rateErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("error type = %T, want *RateLimitError", err)
	}

	if rateErr.RetryAfter != 60 {
		t.Errorf("RetryAfter = %d, want 60", rateErr.RetryAfter)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	err := &ErrorResponse{
		Error: struct {
			Message string                 `json:"message"`
			Type    string                 `json:"type"`
			Code    string                 `json:"code"`
			Param   string                 `json:"param"`
			Details map[string]interface{} `json:"details"`
		}{
			Type:    "test_error",
			Message: "Test message",
		},
	}

	if err.Error.Message != "Test message" {
		t.Errorf("Error.Message = %v, want Test message", err.Error.Message)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Type:       "invalid_request",
		Message:    "Bad request",
		RequestID:  "req_123",
	}

	expected := "zaguan API error (400) [req_123]: Bad request"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestInsufficientCreditsError_Error(t *testing.T) {
	err := &InsufficientCreditsError{
		APIError: APIError{
			StatusCode: 402,
			Message:    "Insufficient credits",
		},
		CreditsRequired:  100,
		CreditsRemaining: 50,
		ResetDate:        "2025-12-01",
	}

	expected := "insufficient credits: required 100, remaining 50 (resets on 2025-12-01)"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestBandAccessError_Error(t *testing.T) {
	err := &BandAccessError{
		APIError: APIError{
			StatusCode: 403,
			Message:    "Band access denied",
		},
		Band:         "premium",
		RequiredTier: "pro",
		CurrentTier:  "free",
	}

	expected := "band access denied: free tier does not have access to band premium (requires pro tier)"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestRateLimitError_Error(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter int
		want       string
	}{
		{
			name:       "with retry after",
			retryAfter: 60,
			want:       "rate limit exceeded: retry after 60 seconds",
		},
		{
			name:       "without retry after",
			retryAfter: 0,
			want:       "rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &RateLimitError{
				APIError: APIError{
					StatusCode: 429,
					Message:    "Rate limit exceeded",
				},
				RetryAfter: tt.retryAfter,
			}

			if err.Error() != tt.want {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}

func TestHTTPClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Authorization header not set")
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent header not set")
		}
		if r.Header.Get("X-Request-Id") == "" {
			t.Error("X-Request-Id header not set")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := NewHTTPClient(&http.Client{}, server.URL, "test-key", "test-version")

	cfg := RequestConfig{
		Method: "GET",
		Path:   "/",
	}
	resp, err := client.Do(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHTTPClient_DoJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(&http.Client{}, server.URL, "test-key", "test-version")

	var result struct {
		Message string `json:"message"`
	}

	cfg := RequestConfig{
		Method: "GET",
		Path:   "/",
	}

	err := client.DoJSON(context.Background(), cfg, &result)
	if err != nil {
		t.Fatalf("DoJSON() error = %v", err)
	}

	if result.Message != "success" {
		t.Errorf("Message = %v, want success", result.Message)
	}
}

func TestHTTPClient_DoJSON_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"type": "invalid_request", "message": "Bad request"}}`))
	}))
	defer server.Close()

	client := NewHTTPClient(&http.Client{}, server.URL, "test-key", "test-version")

	var result map[string]interface{}
	cfg := RequestConfig{
		Method: "GET",
		Path:   "/",
	}

	err := client.DoJSON(context.Background(), cfg, &result)
	if err == nil {
		t.Fatal("DoJSON() should have returned error")
	}
}
