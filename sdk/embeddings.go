// Package zaguansdk provides embeddings functionality for the Zaguan SDK.
//
// This file implements the Embeddings API for creating vector representations
// of text that can be used for semantic search, clustering, recommendations,
// and anomaly detection.
package zaguansdk

import (
	"context"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// EmbeddingsRequest represents a request to create embeddings.
//
// Embeddings are vector representations of text that can be used for
// semantic search, clustering, recommendations, and anomaly detection.
type EmbeddingsRequest struct {
	// Model is the model identifier to use for embeddings.
	// Examples: "openai/text-embedding-3-small", "cohere/embed-english-v3.0"
	// Required.
	Model string `json:"model"`

	// Input is the text or array of texts to embed.
	// Can be a string or []string.
	// Required.
	Input interface{} `json:"input"`

	// EncodingFormat specifies the format of the returned embeddings.
	// Values: "float", "base64"
	// Optional (default: "float").
	EncodingFormat string `json:"encoding_format,omitempty"`

	// Dimensions is the number of dimensions for the embedding.
	// Only supported by some models (e.g., text-embedding-3-*).
	// Optional.
	Dimensions int `json:"dimensions,omitempty"`

	// User is an optional unique identifier for the end-user.
	// Optional.
	User string `json:"user,omitempty"`

	// ProviderSpecificParams contains provider-specific parameters.
	// For Cohere: input_type ("search_document", "search_query", "classification", "clustering")
	// Optional.
	ProviderSpecificParams map[string]interface{} `json:"provider_specific_params,omitempty"`
}

// EmbeddingsResponse represents the response from an embeddings request.
type EmbeddingsResponse struct {
	// Object is the object type (always "list").
	Object string `json:"object"`

	// Data is the list of embeddings.
	Data []Embedding `json:"data"`

	// Model is the model used for the embeddings.
	Model string `json:"model"`

	// Usage contains token usage information.
	Usage EmbeddingsUsage `json:"usage"`
}

// Embedding represents a single embedding vector.
type Embedding struct {
	// Object is the object type (always "embedding").
	Object string `json:"object"`

	// Embedding is the vector representation.
	// Type depends on encoding_format: []float64 for "float", string for "base64"
	Embedding interface{} `json:"embedding"`

	// Index is the index of this embedding in the input array.
	Index int `json:"index"`
}

// EmbeddingsUsage represents token usage for embeddings.
type EmbeddingsUsage struct {
	// PromptTokens is the number of tokens in the input.
	PromptTokens int `json:"prompt_tokens"`

	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
}

// CreateEmbeddings creates embeddings for the given input text(s).
//
// Embeddings are useful for semantic search, clustering, recommendations,
// and anomaly detection. The returned vectors can be compared using
// cosine similarity or other distance metrics.
//
// Example:
//
//	resp, err := client.CreateEmbeddings(ctx, zaguansdk.EmbeddingsRequest{
//		Model: "openai/text-embedding-3-small",
//		Input: []string{
//			"The quick brown fox",
//			"jumps over the lazy dog",
//		},
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, emb := range resp.Data {
//		fmt.Printf("Embedding %d: %d dimensions\n", emb.Index, len(emb.Embedding.([]float64)))
//	}
func (c *Client) CreateEmbeddings(ctx context.Context, req EmbeddingsRequest, opts *RequestOptions) (*EmbeddingsResponse, error) {
	// Validate request
	if err := validateEmbeddingsRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating embeddings", "model", req.Model)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/embeddings",
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
	var resp EmbeddingsResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create embeddings request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create embeddings request succeeded",
		"model", resp.Model,
		"count", len(resp.Data))

	return &resp, nil
}

// GetEmbeddingVector is a helper that extracts the float64 vector from an Embedding.
//
// Returns an error if the embedding is not in float format.
func (e *Embedding) GetEmbeddingVector() ([]float64, error) {
	vec, ok := e.Embedding.([]interface{})
	if !ok {
		return nil, &APIError{
			StatusCode: 0,
			Message:    "embedding is not in float format",
			Type:       "invalid_format",
		}
	}

	result := make([]float64, len(vec))
	for i, v := range vec {
		f, ok := v.(float64)
		if !ok {
			return nil, &APIError{
				StatusCode: 0,
				Message:    "embedding contains non-float value",
				Type:       "invalid_format",
			}
		}
		result[i] = f
	}

	return result, nil
}

// CosineSimilarity calculates the cosine similarity between two embedding vectors.
//
// Returns a value between -1 and 1, where 1 means identical, 0 means orthogonal,
// and -1 means opposite.
func CosineSimilarity(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, &APIError{
			StatusCode: 0,
			Message:    "vectors must have the same length",
			Type:       "invalid_input",
		}
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, &APIError{
			StatusCode: 0,
			Message:    "cannot compute similarity with zero vector",
			Type:       "invalid_input",
		}
	}

	return dotProduct / (normA * normB), nil
}
