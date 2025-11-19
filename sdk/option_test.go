package zaguansdk

import (
	"net/http"
	"testing"
	"time"
)

func TestWithRequestID(t *testing.T) {
	opts := WithRequestID("test-123")
	if opts.RequestID != "test-123" {
		t.Errorf("WithRequestID() RequestID = %v, want test-123", opts.RequestID)
	}
}

func TestWithTimeout(t *testing.T) {
	timeout := 30 * time.Second
	opts := WithTimeout(timeout)
	if opts.Timeout != timeout {
		t.Errorf("WithTimeout() Timeout = %v, want %v", opts.Timeout, timeout)
	}
}

func TestWithHeaders(t *testing.T) {
	headers := http.Header{
		"X-Custom": []string{"value"},
	}
	opts := WithHeaders(headers)
	if opts.Headers.Get("X-Custom") != "value" {
		t.Error("WithHeaders() did not set headers correctly")
	}
}

func TestWithRetries(t *testing.T) {
	opts := WithRetries(3, 5*time.Second)
	if opts.MaxRetries != 3 {
		t.Errorf("WithRetries() MaxRetries = %v, want 3", opts.MaxRetries)
	}
	if opts.RetryDelay != 5*time.Second {
		t.Errorf("WithRetries() RetryDelay = %v, want 5s", opts.RetryDelay)
	}
}

func TestRequestOptions_Merge(t *testing.T) {
	tests := []struct {
		name  string
		base  *RequestOptions
		other *RequestOptions
		want  *RequestOptions
	}{
		{
			name: "merge with nil other",
			base: &RequestOptions{
				RequestID: "base-id",
				Timeout:   10 * time.Second,
			},
			other: nil,
			want: &RequestOptions{
				RequestID: "base-id",
				Timeout:   10 * time.Second,
			},
		},
		{
			name: "other overrides base",
			base: &RequestOptions{
				RequestID: "base-id",
				Timeout:   10 * time.Second,
			},
			other: &RequestOptions{
				RequestID: "other-id",
				Timeout:   20 * time.Second,
			},
			want: &RequestOptions{
				RequestID: "other-id",
				Timeout:   20 * time.Second,
			},
		},
		{
			name: "partial override",
			base: &RequestOptions{
				RequestID: "base-id",
				Timeout:   10 * time.Second,
			},
			other: &RequestOptions{
				Timeout: 20 * time.Second,
			},
			want: &RequestOptions{
				RequestID: "base-id",
				Timeout:   20 * time.Second,
			},
		},
		{
			name: "merge headers",
			base: &RequestOptions{
				Headers: http.Header{
					"X-Base": []string{"base"},
				},
			},
			other: &RequestOptions{
				Headers: http.Header{
					"X-Other": []string{"other"},
				},
			},
			want: &RequestOptions{
				Headers: http.Header{
					"X-Base":  []string{"base"},
					"X-Other": []string{"other"},
				},
			},
		},
		{
			name: "other header overrides base",
			base: &RequestOptions{
				Headers: http.Header{
					"X-Header": []string{"base"},
				},
			},
			other: &RequestOptions{
				Headers: http.Header{
					"X-Header": []string{"other"},
				},
			},
			want: &RequestOptions{
				Headers: http.Header{
					"X-Header": []string{"other"},
				},
			},
		},
		{
			name: "merge retries",
			base: &RequestOptions{
				MaxRetries: 3,
				RetryDelay: 1 * time.Second,
			},
			other: &RequestOptions{
				MaxRetries: 5,
			},
			want: &RequestOptions{
				MaxRetries: 5,
				RetryDelay: 1 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Merge(tt.other)

			if got.RequestID != tt.want.RequestID {
				t.Errorf("Merge() RequestID = %v, want %v", got.RequestID, tt.want.RequestID)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("Merge() Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
			if got.MaxRetries != tt.want.MaxRetries {
				t.Errorf("Merge() MaxRetries = %v, want %v", got.MaxRetries, tt.want.MaxRetries)
			}
			if got.RetryDelay != tt.want.RetryDelay {
				t.Errorf("Merge() RetryDelay = %v, want %v", got.RetryDelay, tt.want.RetryDelay)
			}

			// Check headers
			if tt.want.Headers != nil {
				for key, wantValues := range tt.want.Headers {
					gotValues := got.Headers[key]
					if len(gotValues) != len(wantValues) {
						t.Errorf("Merge() Headers[%s] = %v, want %v", key, gotValues, wantValues)
					}
				}
			}
		})
	}
}

func TestRequestOptions_MergeNilBase(t *testing.T) {
	var base *RequestOptions
	other := &RequestOptions{
		RequestID: "other-id",
		Timeout:   10 * time.Second,
	}

	got := base.Merge(other)
	if got.RequestID != "other-id" {
		t.Errorf("Merge() with nil base RequestID = %v, want other-id", got.RequestID)
	}
}
