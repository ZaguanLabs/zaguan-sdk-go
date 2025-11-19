package zaguansdk

import (
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     APIError
		wantMsg string
	}{
		{
			name: "with request ID",
			err: APIError{
				StatusCode: 400,
				Message:    "Invalid request",
				RequestID:  "req_123",
			},
			wantMsg: "zaguan API error (400) [req_123]: Invalid request",
		},
		{
			name: "without request ID",
			err: APIError{
				StatusCode: 500,
				Message:    "Internal error",
			},
			wantMsg: "zaguan API error (500): Internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestAPIError_IsInsufficientCredits(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "insufficient credits type",
			err:  APIError{Type: "insufficient_credits"},
			want: true,
		},
		{
			name: "insufficient credits code",
			err:  APIError{Code: "insufficient_credits"},
			want: true,
		},
		{
			name: "other error",
			err:  APIError{Type: "rate_limit_exceeded"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsInsufficientCredits(); got != tt.want {
				t.Errorf("IsInsufficientCredits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsBandAccessDenied(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "band access denied type",
			err:  APIError{Type: "band_access_denied"},
			want: true,
		},
		{
			name: "band access denied code",
			err:  APIError{Code: "band_access_denied"},
			want: true,
		},
		{
			name: "other error",
			err:  APIError{Type: "authentication_error"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsBandAccessDenied(); got != tt.want {
				t.Errorf("IsBandAccessDenied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsRateLimitExceeded(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "rate limit type",
			err:  APIError{Type: "rate_limit_exceeded"},
			want: true,
		},
		{
			name: "rate limit code",
			err:  APIError{Code: "rate_limit_exceeded"},
			want: true,
		},
		{
			name: "other error",
			err:  APIError{Type: "invalid_request"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRateLimitExceeded(); got != tt.want {
				t.Errorf("IsRateLimitExceeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsAuthenticationError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "401 status",
			err:  APIError{StatusCode: 401},
			want: true,
		},
		{
			name: "authentication_error type",
			err:  APIError{Type: "authentication_error"},
			want: true,
		},
		{
			name: "other error",
			err:  APIError{StatusCode: 400},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsAuthenticationError(); got != tt.want {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsPermissionError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "403 status",
			err:  APIError{StatusCode: 403},
			want: true,
		},
		{
			name: "permission_error type",
			err:  APIError{Type: "permission_error"},
			want: true,
		},
		{
			name: "other error",
			err:  APIError{StatusCode: 400},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsPermissionError(); got != tt.want {
				t.Errorf("IsPermissionError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "404 status",
			err:  APIError{StatusCode: 404},
			want: true,
		},
		{
			name: "other status",
			err:  APIError{StatusCode: 400},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsNotFoundError(); got != tt.want {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want bool
	}{
		{
			name: "500 status",
			err:  APIError{StatusCode: 500},
			want: true,
		},
		{
			name: "502 status",
			err:  APIError{StatusCode: 502},
			want: true,
		},
		{
			name: "599 status",
			err:  APIError{StatusCode: 599},
			want: true,
		},
		{
			name: "400 status",
			err:  APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "600 status",
			err:  APIError{StatusCode: 600},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsServerError(); got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsufficientCreditsError_Error(t *testing.T) {
	err := InsufficientCreditsError{
		APIError: APIError{
			StatusCode: 402,
			Message:    "Insufficient credits",
		},
		CreditsRequired:  100,
		CreditsRemaining: 50,
		ResetDate:        "2025-12-01",
	}

	expected := "insufficient credits: required 100, remaining 50 (resets on 2025-12-01)"
	if got := err.Error(); got != expected {
		t.Errorf("InsufficientCreditsError.Error() = %v, want %v", got, expected)
	}
}

func TestBandAccessError_Error(t *testing.T) {
	err := BandAccessError{
		APIError: APIError{
			StatusCode: 403,
			Message:    "Band access denied",
		},
		Band:         "premium",
		RequiredTier: "pro",
		CurrentTier:  "free",
	}

	expected := "band access denied: free tier does not have access to band premium (requires pro tier)"
	if got := err.Error(); got != expected {
		t.Errorf("BandAccessError.Error() = %v, want %v", got, expected)
	}
}

func TestRateLimitError_Error(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter int
		wantMsg    string
	}{
		{
			name:       "with retry after",
			retryAfter: 60,
			wantMsg:    "rate limit exceeded: retry after 60 seconds",
		},
		{
			name:       "without retry after",
			retryAfter: 0,
			wantMsg:    "rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RateLimitError{
				APIError: APIError{
					StatusCode: 429,
					Message:    "Rate limit exceeded",
				},
				RetryAfter: tt.retryAfter,
			}

			if got := err.Error(); got != tt.wantMsg {
				t.Errorf("RateLimitError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "temperature",
		Message: "must be between 0 and 2",
	}

	expected := "validation error: temperature: must be between 0 and 2"
	if got := err.Error(); got != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", got, expected)
	}
}
