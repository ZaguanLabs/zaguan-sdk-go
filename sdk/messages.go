package zaguansdk

// MessagesRequest represents a request to Anthropic's native Messages API.
//
// This follows the Anthropic Messages API format and is exposed via the
// /v1/messages endpoint in Zaguan CoreX.
type MessagesRequest struct {
	// Model is the model identifier to use.
	// Format: "anthropic/model-name" (e.g., "anthropic/claude-3-5-sonnet-20241022")
	// Required.
	Model string `json:"model"`

	// Messages is the conversation history.
	// Required.
	Messages []AnthropicMessage `json:"messages"`

	// System is the system prompt.
	// Optional.
	System string `json:"system,omitempty"`

	// MaxTokens is the maximum number of tokens to generate.
	// Required for Anthropic API.
	MaxTokens int `json:"max_tokens"`

	// Temperature controls randomness (0.0 - 1.0).
	// Optional.
	Temperature *float64 `json:"temperature,omitempty"`

	// TopP controls nucleus sampling (0.0 - 1.0).
	// Optional.
	TopP *float64 `json:"top_p,omitempty"`

	// TopK controls top-k sampling.
	// Optional.
	TopK *int `json:"top_k,omitempty"`

	// Stream enables streaming responses.
	// Use MessagesStream() method instead of Messages() when this is true.
	// Optional.
	Stream bool `json:"stream,omitempty"`

	// StopSequences are sequences that will halt generation.
	// Optional.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Thinking configures extended thinking (Beta feature).
	// Optional.
	Thinking *AnthropicThinkingConfig `json:"thinking,omitempty"`

	// Metadata for application-specific tracking.
	// Optional.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AnthropicMessage represents a message in Anthropic's format.
type AnthropicMessage struct {
	// Role is the message role.
	// Values: "user", "assistant"
	// Required.
	Role string `json:"role"`

	// Content is the message content.
	// Can be a string or an array of content blocks for multimodal.
	// Required.
	Content interface{} `json:"content"`
}

// AnthropicThinkingConfig configures extended thinking (Beta).
type AnthropicThinkingConfig struct {
	// Type controls thinking behavior.
	// Values: "enabled", "disabled"
	// Required.
	Type string `json:"type"`

	// BudgetTokens is the maximum number of thinking tokens (1,000 - 10,000).
	// Optional.
	BudgetTokens int `json:"budget_tokens,omitempty"`
}

// MessagesResponse represents a response from Anthropic's Messages API.
type MessagesResponse struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`

	// Type is the object type (always "message").
	Type string `json:"type"`

	// Role is the message role (always "assistant").
	Role string `json:"role"`

	// Content is the message content blocks.
	Content []AnthropicContentBlock `json:"content"`

	// Model is the model used for the completion.
	Model string `json:"model"`

	// StopReason indicates why the generation stopped.
	// Values: "end_turn", "max_tokens", "stop_sequence", "tool_use"
	StopReason string `json:"stop_reason,omitempty"`

	// StopSequence is the stop sequence that was matched (if any).
	StopSequence string `json:"stop_sequence,omitempty"`

	// Usage contains token usage information.
	Usage AnthropicUsage `json:"usage"`
}

// AnthropicContentBlock represents a content block in the response.
type AnthropicContentBlock struct {
	// Type is the content block type.
	// Values: "text", "thinking", "tool_use"
	Type string `json:"type"`

	// Text content (for type="text").
	Text string `json:"text,omitempty"`

	// Thinking content (for type="thinking").
	// This is the extended thinking output when thinking is enabled.
	Thinking string `json:"thinking,omitempty"`

	// Signature is the cryptographic signature for thinking verification.
	Signature string `json:"signature,omitempty"`

	// ID is the tool use ID (for type="tool_use").
	ID string `json:"id,omitempty"`

	// Name is the tool name (for type="tool_use").
	Name string `json:"name,omitempty"`

	// Input is the tool input (for type="tool_use").
	Input interface{} `json:"input,omitempty"`
}

// AnthropicUsage represents token usage in Anthropic's format.
type AnthropicUsage struct {
	// InputTokens is the number of tokens in the input.
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of tokens in the output.
	OutputTokens int `json:"output_tokens"`

	// CacheCreationInputTokens is the number of tokens used to create the cache.
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`

	// CacheReadInputTokens is the number of tokens read from the cache.
	CacheReadInputTokens int `json:"cache_read_input_tokens,omitempty"`
}

// CountTokensRequest represents a request to count tokens.
type CountTokensRequest struct {
	// Model is the model identifier.
	Model string `json:"model"`

	// Messages is the conversation to count tokens for.
	Messages []AnthropicMessage `json:"messages"`

	// System is the system prompt.
	System string `json:"system,omitempty"`
}

// CountTokensResponse represents the response from token counting.
type CountTokensResponse struct {
	// InputTokens is the number of tokens in the input.
	InputTokens int `json:"input_tokens"`
}

// MessagesBatchRequest represents a request to create a message batch.
type MessagesBatchRequest struct {
	// Requests is the array of message requests to process in batch.
	Requests []MessagesBatchItem `json:"requests"`
}

// MessagesBatchItem represents a single item in a batch.
type MessagesBatchItem struct {
	// CustomID is a user-provided identifier for this request.
	CustomID string `json:"custom_id"`

	// Params is the messages request parameters.
	Params MessagesRequest `json:"params"`
}

// MessagesBatchResponse represents a batch processing response.
type MessagesBatchResponse struct {
	// ID is the unique identifier for this batch.
	ID string `json:"id"`

	// Type is the object type (always "message_batch").
	Type string `json:"type"`

	// ProcessingStatus is the current processing status.
	// Values: "in_progress", "canceling", "ended"
	ProcessingStatus string `json:"processing_status"`

	// RequestCounts contains counts of requests by status.
	RequestCounts MessagesBatchRequestCounts `json:"request_counts"`

	// EndedAt is when the batch processing ended.
	EndedAt string `json:"ended_at,omitempty"`

	// CreatedAt is when the batch was created.
	CreatedAt string `json:"created_at"`

	// ExpiresAt is when the batch will expire.
	ExpiresAt string `json:"expires_at"`

	// ArchivedAt is when the batch was archived.
	ArchivedAt string `json:"archived_at,omitempty"`

	// CancelInitiatedAt is when cancellation was initiated.
	CancelInitiatedAt string `json:"cancel_initiated_at,omitempty"`

	// ResultsURL is the URL to download results.
	ResultsURL string `json:"results_url,omitempty"`
}

// MessagesBatchRequestCounts contains counts of batch requests.
type MessagesBatchRequestCounts struct {
	// Processing is the number of requests currently processing.
	Processing int `json:"processing"`

	// Succeeded is the number of successful requests.
	Succeeded int `json:"succeeded"`

	// Errored is the number of failed requests.
	Errored int `json:"errored"`

	// Canceled is the number of canceled requests.
	Canceled int `json:"canceled"`

	// Expired is the number of expired requests.
	Expired int `json:"expired"`
}
