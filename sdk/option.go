package zaguansdk

import (
	"net/http"
	"time"
)

// RequestOptions allows per-request configuration overrides.
//
// All fields are optional. If not specified, the client's default configuration
// will be used.
type RequestOptions struct {
	// RequestID is a unique identifier for this request.
	// If empty, a UUID will be automatically generated.
	// This ID is sent in the X-Request-Id header and can be used for debugging.
	RequestID string

	// Timeout overrides the client's default timeout for this request.
	// If zero, the client's default timeout is used.
	Timeout time.Duration

	// Headers are additional HTTP headers to include in the request.
	// These will be merged with the default headers (Authorization, Content-Type, etc.).
	Headers http.Header

	// MaxRetries specifies the maximum number of retry attempts for this request.
	// If zero, no retries will be attempted.
	// If negative, the client's default retry policy is used.
	MaxRetries int

	// RetryDelay is the initial delay between retry attempts.
	// Subsequent retries use exponential backoff.
	// If zero, a default of 1 second is used.
	RetryDelay time.Duration
}

// WithRequestID returns a new RequestOptions with the specified request ID.
func WithRequestID(id string) *RequestOptions {
	return &RequestOptions{RequestID: id}
}

// WithTimeout returns a new RequestOptions with the specified timeout.
func WithTimeout(timeout time.Duration) *RequestOptions {
	return &RequestOptions{Timeout: timeout}
}

// WithHeaders returns a new RequestOptions with the specified headers.
func WithHeaders(headers http.Header) *RequestOptions {
	return &RequestOptions{Headers: headers}
}

// WithRetries returns a new RequestOptions with the specified retry configuration.
func WithRetries(maxRetries int, delay time.Duration) *RequestOptions {
	return &RequestOptions{
		MaxRetries: maxRetries,
		RetryDelay: delay,
	}
}

// Merge merges this RequestOptions with another, with the other taking precedence.
func (o *RequestOptions) Merge(other *RequestOptions) *RequestOptions {
	if other == nil {
		return o
	}

	merged := &RequestOptions{}

	// Request ID
	if other.RequestID != "" {
		merged.RequestID = other.RequestID
	} else if o != nil {
		merged.RequestID = o.RequestID
	}

	// Timeout
	if other.Timeout > 0 {
		merged.Timeout = other.Timeout
	} else if o != nil {
		merged.Timeout = o.Timeout
	}

	// Headers
	merged.Headers = make(http.Header)
	if o != nil && o.Headers != nil {
		for k, v := range o.Headers {
			merged.Headers[k] = v
		}
	}
	if other.Headers != nil {
		for k, v := range other.Headers {
			merged.Headers[k] = v
		}
	}

	// Retries
	if other.MaxRetries != 0 {
		merged.MaxRetries = other.MaxRetries
	} else if o != nil {
		merged.MaxRetries = o.MaxRetries
	}

	if other.RetryDelay > 0 {
		merged.RetryDelay = other.RetryDelay
	} else if o != nil {
		merged.RetryDelay = o.RetryDelay
	}

	return merged
}
