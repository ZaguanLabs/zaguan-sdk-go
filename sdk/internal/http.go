package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// HTTPClient is an internal wrapper around http.Client with Zaguan-specific functionality.
type HTTPClient struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	userAgent string
}

// NewHTTPClient creates a new internal HTTP client.
func NewHTTPClient(client *http.Client, baseURL, apiKey, sdkVersion string) *HTTPClient {
	return &HTTPClient{
		client:    client,
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: fmt.Sprintf("zaguan-go-sdk/%s", sdkVersion),
	}
}

// RequestConfig holds configuration for an HTTP request.
type RequestConfig struct {
	Method      string
	Path        string
	Body        interface{}
	Headers     http.Header
	RequestID   string
	Timeout     time.Duration
	QueryParams map[string]string
}

// Do executes an HTTP request and returns the response.
func (c *HTTPClient) Do(ctx context.Context, cfg RequestConfig) (*http.Response, error) {
	// Build URL
	url := c.baseURL + cfg.Path
	if len(cfg.QueryParams) > 0 {
		url += "?"
		first := true
		for k, v := range cfg.QueryParams {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	// Marshal body if present
	var bodyReader io.Reader
	if cfg.Body != nil {
		bodyBytes, err := json.Marshal(cfg.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, cfg.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	// Set request ID
	requestID := cfg.RequestID
	if requestID == "" {
		requestID = uuid.New().String()
	}
	req.Header.Set("X-Request-Id", requestID)

	// Merge custom headers
	if cfg.Headers != nil {
		for k, v := range cfg.Headers {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
	}

	// Apply timeout if specified
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// DoJSON executes an HTTP request and unmarshals the JSON response.
func (c *HTTPClient) DoJSON(ctx context.Context, cfg RequestConfig, result interface{}) error {
	resp, err := c.Do(ctx, cfg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return ParseErrorResponse(resp)
	}

	// Decode response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// ErrorResponse represents the error response format from the API.
type ErrorResponse struct {
	Error struct {
		Message string                 `json:"message"`
		Type    string                 `json:"type"`
		Code    string                 `json:"code"`
		Param   string                 `json:"param"`
		Details map[string]interface{} `json:"details"`
	} `json:"error"`
}

// ParseErrorResponse parses an error response from the API.
func ParseErrorResponse(resp *http.Response) error {
	requestID := resp.Header.Get("X-Request-Id")

	// Try to parse as structured error
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		// If we can't parse the error, return a generic one
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			RequestID:  requestID,
		}
	}

	// Create base API error
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Message:    errResp.Error.Message,
		Type:       errResp.Error.Type,
		Code:       errResp.Error.Code,
		Param:      errResp.Error.Param,
		RequestID:  requestID,
		Details:    errResp.Error.Details,
	}

	// Check for specialized error types
	if apiErr.Type == "insufficient_credits" || apiErr.Code == "insufficient_credits" {
		return parseInsufficientCreditsError(apiErr)
	}

	if apiErr.Type == "band_access_denied" || apiErr.Code == "band_access_denied" {
		return parseBandAccessError(apiErr)
	}

	if apiErr.Type == "rate_limit_exceeded" || apiErr.Code == "rate_limit_exceeded" {
		return parseRateLimitError(apiErr, resp)
	}

	return apiErr
}

// APIError represents a generic API error.
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
	Type       string
	Code       string
	Param      string
	Details    map[string]interface{}
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("zaguan API error (%d) [%s]: %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("zaguan API error (%d): %s", e.StatusCode, e.Message)
}

func parseInsufficientCreditsError(base *APIError) error {
	err := &InsufficientCreditsError{APIError: *base}

	if base.Details != nil {
		if v, ok := base.Details["credits_required"].(float64); ok {
			err.CreditsRequired = int(v)
		}
		if v, ok := base.Details["credits_remaining"].(float64); ok {
			err.CreditsRemaining = int(v)
		}
		if v, ok := base.Details["reset_date"].(string); ok {
			err.ResetDate = v
		}
	}

	return err
}

func parseBandAccessError(base *APIError) error {
	err := &BandAccessError{APIError: *base}

	if base.Details != nil {
		if v, ok := base.Details["band"].(string); ok {
			err.Band = v
		}
		if v, ok := base.Details["required_tier"].(string); ok {
			err.RequiredTier = v
		}
		if v, ok := base.Details["current_tier"].(string); ok {
			err.CurrentTier = v
		}
	}

	return err
}

func parseRateLimitError(base *APIError, resp *http.Response) error {
	err := &RateLimitError{APIError: *base}

	// Try to get retry-after from header
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		_, _ = fmt.Sscanf(retryAfter, "%d", &err.RetryAfter) // Ignore parse error, will use default
	}

	// Also check details
	if base.Details != nil {
		if v, ok := base.Details["retry_after"].(float64); ok {
			err.RetryAfter = int(v)
		}
	}

	return err
}

// InsufficientCreditsError represents an insufficient credits error.
type InsufficientCreditsError struct {
	APIError
	CreditsRequired  int
	CreditsRemaining int
	ResetDate        string
}

func (e *InsufficientCreditsError) Error() string {
	return fmt.Sprintf("insufficient credits: required %d, remaining %d (resets on %s)",
		e.CreditsRequired, e.CreditsRemaining, e.ResetDate)
}

// BandAccessError represents a band access denied error.
type BandAccessError struct {
	APIError
	Band         string
	RequiredTier string
	CurrentTier  string
}

func (e *BandAccessError) Error() string {
	return fmt.Sprintf("band access denied: %s tier does not have access to band %s (requires %s tier)",
		e.CurrentTier, e.Band, e.RequiredTier)
}

// RateLimitError represents a rate limit error.
type RateLimitError struct {
	APIError
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded: retry after %d seconds", e.RetryAfter)
	}
	return "rate limit exceeded"
}
