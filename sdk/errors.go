package zaguansdk

import (
	"fmt"
)

// APIError represents an error returned by the Zaguan CoreX API.
//
// It includes the HTTP status code, error message, request ID for debugging,
// and an optional error type/code for programmatic handling.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// Message is the human-readable error message.
	Message string

	// RequestID is the unique identifier for this request.
	// Include this when reporting issues to Zaguan support.
	RequestID string

	// Type is the error type/code for programmatic handling.
	// Examples: "insufficient_credits", "band_access_denied", "rate_limit_exceeded"
	Type string

	// Code is an optional error code (may be the same as Type).
	Code string

	// Param is the parameter that caused the error, if applicable.
	Param string

	// Details contains additional error details from the API.
	Details map[string]interface{}
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("zaguan API error (%d) [%s]: %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("zaguan API error (%d): %s", e.StatusCode, e.Message)
}

// IsInsufficientCredits returns true if this error is due to insufficient credits.
func (e *APIError) IsInsufficientCredits() bool {
	return e.Type == "insufficient_credits" || e.Code == "insufficient_credits"
}

// IsBandAccessDenied returns true if this error is due to band access restrictions.
func (e *APIError) IsBandAccessDenied() bool {
	return e.Type == "band_access_denied" || e.Code == "band_access_denied"
}

// IsRateLimitExceeded returns true if this error is due to rate limiting.
func (e *APIError) IsRateLimitExceeded() bool {
	return e.Type == "rate_limit_exceeded" || e.Code == "rate_limit_exceeded"
}

// IsAuthenticationError returns true if this error is due to authentication failure.
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == 401 || e.Type == "authentication_error"
}

// IsPermissionError returns true if this error is due to insufficient permissions.
func (e *APIError) IsPermissionError() bool {
	return e.StatusCode == 403 || e.Type == "permission_error"
}

// IsNotFoundError returns true if the requested resource was not found.
func (e *APIError) IsNotFoundError() bool {
	return e.StatusCode == 404
}

// IsServerError returns true if this is a server-side error (5xx).
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// InsufficientCreditsError represents an error when the user has insufficient credits.
//
// This is a specialized error type that includes credit balance information.
type InsufficientCreditsError struct {
	APIError
	CreditsRequired  int
	CreditsRemaining int
	ResetDate        string
}

// Error implements the error interface.
func (e *InsufficientCreditsError) Error() string {
	return fmt.Sprintf("insufficient credits: required %d, remaining %d (resets on %s)",
		e.CreditsRequired, e.CreditsRemaining, e.ResetDate)
}

// BandAccessError represents an error when the user's tier doesn't have access to a band.
//
// This is a specialized error type that includes tier and band information.
type BandAccessError struct {
	APIError
	Band         string
	RequiredTier string
	CurrentTier  string
}

// Error implements the error interface.
func (e *BandAccessError) Error() string {
	return fmt.Sprintf("band access denied: %s tier does not have access to band %s (requires %s tier)",
		e.CurrentTier, e.Band, e.RequiredTier)
}

// RateLimitError represents a rate limit error.
//
// This is a specialized error type that includes retry-after information.
type RateLimitError struct {
	APIError
	RetryAfter int // Seconds to wait before retrying
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded: retry after %d seconds", e.RetryAfter)
	}
	return "rate limit exceeded"
}
