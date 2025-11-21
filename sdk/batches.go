// Package zaguansdk provides batch processing functionality for the Zaguan SDK.
//
// This file implements the Batches API for processing multiple API requests
// asynchronously with 50% cost reduction compared to synchronous requests.
//
// Batches support endpoints including:
//   - /v1/chat/completions
//   - /v1/embeddings
//   - /v1/completions
//
// With a 24-hour completion window and comprehensive status tracking.
package zaguansdk

import (
	"context"
	"fmt"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// BatchRequest represents a request to create a batch job.
//
// Batches allow you to process multiple API requests asynchronously
// with 50% cost reduction compared to synchronous requests.
type BatchRequest struct {
	// InputFileID is the ID of the uploaded file containing requests.
	// Required.
	InputFileID string `json:"input_file_id"`

	// Endpoint is the API endpoint to use for the batch.
	// Values: "/v1/chat/completions", "/v1/embeddings", "/v1/completions"
	// Required.
	Endpoint string `json:"endpoint"`

	// CompletionWindow is the time window for completion.
	// Values: "24h"
	// Required.
	CompletionWindow string `json:"completion_window"`

	// Metadata is optional custom metadata.
	// Optional.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// BatchResponse represents a batch job.
type BatchResponse struct {
	// ID is the unique identifier for the batch.
	ID string `json:"id"`

	// Object is the object type (always "batch").
	Object string `json:"object"`

	// Endpoint is the API endpoint used.
	Endpoint string `json:"endpoint"`

	// Errors contains error information if the batch failed.
	Errors *BatchErrors `json:"errors,omitempty"`

	// InputFileID is the ID of the input file.
	InputFileID string `json:"input_file_id"`

	// CompletionWindow is the time window for completion.
	CompletionWindow string `json:"completion_window"`

	// Status is the current status of the batch.
	// Values: "validating", "failed", "in_progress", "finalizing", "completed", "expired", "cancelling", "cancelled"
	Status string `json:"status"`

	// OutputFileID is the ID of the output file (when completed).
	OutputFileID string `json:"output_file_id,omitempty"`

	// ErrorFileID is the ID of the error file (if errors occurred).
	ErrorFileID string `json:"error_file_id,omitempty"`

	// CreatedAt is the Unix timestamp of when the batch was created.
	CreatedAt int64 `json:"created_at"`

	// InProgressAt is the Unix timestamp of when processing started.
	InProgressAt int64 `json:"in_progress_at,omitempty"`

	// ExpiresAt is the Unix timestamp of when the batch expires.
	ExpiresAt int64 `json:"expires_at,omitempty"`

	// FinalizingAt is the Unix timestamp of when finalization started.
	FinalizingAt int64 `json:"finalizing_at,omitempty"`

	// CompletedAt is the Unix timestamp of when the batch completed.
	CompletedAt int64 `json:"completed_at,omitempty"`

	// FailedAt is the Unix timestamp of when the batch failed.
	FailedAt int64 `json:"failed_at,omitempty"`

	// ExpiredAt is the Unix timestamp of when the batch expired.
	ExpiredAt int64 `json:"expired_at,omitempty"`

	// CancellingAt is the Unix timestamp of when cancellation started.
	CancellingAt int64 `json:"cancelling_at,omitempty"`

	// CancelledAt is the Unix timestamp of when the batch was cancelled.
	CancelledAt int64 `json:"cancelled_at,omitempty"`

	// RequestCounts contains counts of requests by status.
	RequestCounts BatchRequestCounts `json:"request_counts"`

	// Metadata is custom metadata.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// BatchErrors contains error information for a batch.
type BatchErrors struct {
	// Object is the object type.
	Object string `json:"object,omitempty"`

	// Data contains the list of errors.
	Data []BatchError `json:"data,omitempty"`
}

// BatchError represents a single batch error.
type BatchError struct {
	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Message is the error message.
	Message string `json:"message,omitempty"`

	// Param is the parameter that caused the error.
	Param string `json:"param,omitempty"`

	// Line is the line number in the input file.
	Line int `json:"line,omitempty"`
}

// BatchRequestCounts contains counts of requests by status.
type BatchRequestCounts struct {
	// Total is the total number of requests.
	Total int `json:"total"`

	// Completed is the number of completed requests.
	Completed int `json:"completed"`

	// Failed is the number of failed requests.
	Failed int `json:"failed"`
}

// BatchListResponse represents a list of batches.
type BatchListResponse struct {
	// Object is the object type (always "list").
	Object string `json:"object"`

	// Data is the list of batches.
	Data []BatchResponse `json:"data"`

	// FirstID is the ID of the first batch in the list.
	FirstID string `json:"first_id,omitempty"`

	// LastID is the ID of the last batch in the list.
	LastID string `json:"last_id,omitempty"`

	// HasMore indicates if there are more batches available.
	HasMore bool `json:"has_more"`
}

// CreateBatch creates a new batch job.
//
// Batches allow you to process multiple API requests asynchronously
// with 50% cost reduction.
//
// Example:
//
//	resp, err := client.CreateBatch(ctx, zaguansdk.BatchRequest{
//		InputFileID:      "file-abc123",
//		Endpoint:         "/v1/chat/completions",
//		CompletionWindow: "24h",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Batch ID:", resp.ID)
func (c *Client) CreateBatch(ctx context.Context, req BatchRequest, opts *RequestOptions) (*BatchResponse, error) {
	// Validate request
	if err := validateBatchRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating batch", "endpoint", req.Endpoint)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/batches",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp BatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}

// GetBatch retrieves information about a specific batch.
//
// Example:
//
//	batch, err := client.GetBatch(ctx, "batch_abc123", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Status:", batch.Status)
func (c *Client) GetBatch(ctx context.Context, batchID string, opts *RequestOptions) (*BatchResponse, error) {
	if batchID == "" {
		return nil, &ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}

	c.log(ctx, LogLevelDebug, "getting batch", "batch_id", batchID)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "GET",
		Path:   "/v1/batches/" + batchID,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp BatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "get batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "get batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}

// ListBatches lists all batches with optional filtering.
//
// Example:
//
//	batches, err := client.ListBatches(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, batch := range batches.Data {
//		fmt.Printf("%s: %s\n", batch.ID, batch.Status)
//	}
func (c *Client) ListBatches(ctx context.Context, opts *RequestOptions) (*BatchListResponse, error) {
	c.log(ctx, LogLevelDebug, "listing batches")

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "GET",
		Path:   "/v1/batches",
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp BatchListResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "list batches request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "list batches request succeeded", "count", len(resp.Data))

	return &resp, nil
}

// CancelBatch cancels a batch that is in progress.
//
// Example:
//
//	batch, err := client.CancelBatch(ctx, "batch_abc123", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Status:", batch.Status)
func (c *Client) CancelBatch(ctx context.Context, batchID string, opts *RequestOptions) (*BatchResponse, error) {
	if batchID == "" {
		return nil, &ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}

	c.log(ctx, LogLevelDebug, "cancelling batch", "batch_id", batchID)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   fmt.Sprintf("/v1/batches/%s/cancel", batchID),
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp BatchResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "cancel batch request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "cancel batch request succeeded", "batch_id", resp.ID)

	return &resp, nil
}

// IsCompleted returns true if the batch has completed successfully.
func (b *BatchResponse) IsCompleted() bool {
	return b.Status == "completed"
}

// IsFailed returns true if the batch has failed.
func (b *BatchResponse) IsFailed() bool {
	return b.Status == "failed"
}

// IsInProgress returns true if the batch is currently processing.
func (b *BatchResponse) IsInProgress() bool {
	return b.Status == "in_progress" || b.Status == "validating" || b.Status == "finalizing"
}
