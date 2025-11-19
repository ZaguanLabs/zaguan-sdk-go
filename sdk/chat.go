package zaguansdk

// ChatRequest represents a request to the chat completions endpoint.
//
// This follows the OpenAI chat completions API format with Zaguan extensions.
type ChatRequest struct {
	// Model is the model identifier to use.
	// Format: "provider/model-name" (e.g., "openai/gpt-4o", "anthropic/claude-3-5-sonnet")
	// Required.
	Model string `json:"model"`

	// Messages is the conversation history.
	// Required.
	Messages []Message `json:"messages"`

	// Temperature controls randomness in the response (0.0 - 2.0).
	// Lower values make output more focused and deterministic.
	// Optional.
	Temperature *float32 `json:"temperature,omitempty"`

	// MaxTokens is the maximum number of tokens to generate.
	// Optional.
	MaxTokens *int `json:"max_tokens,omitempty"`

	// TopP controls nucleus sampling (0.0 - 1.0).
	// Alternative to temperature for controlling randomness.
	// Optional.
	TopP *float32 `json:"top_p,omitempty"`

	// Stream enables streaming responses via Server-Sent Events.
	// Use ChatStream() method instead of Chat() when this is true.
	// Optional.
	Stream bool `json:"stream,omitempty"`

	// Stop sequences that will halt generation.
	// Optional.
	Stop []string `json:"stop,omitempty"`

	// PresencePenalty penalizes new tokens based on whether they appear in the text so far (-2.0 - 2.0).
	// Optional.
	PresencePenalty *float32 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty penalizes new tokens based on their frequency in the text so far (-2.0 - 2.0).
	// Optional.
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty"`

	// LogitBias modifies the likelihood of specified tokens appearing.
	// Optional.
	LogitBias map[string]float32 `json:"logit_bias,omitempty"`

	// User is an identifier for the end-user (for abuse monitoring).
	// Optional.
	User string `json:"user,omitempty"`

	// Tools available for the model to call.
	// Optional.
	Tools []Tool `json:"tools,omitempty"`

	// ToolChoice controls which tools the model can call.
	// Can be "none", "auto", "required", or a specific tool.
	// Optional.
	ToolChoice interface{} `json:"tool_choice,omitempty"`

	// ParallelToolCalls enables parallel tool execution.
	// Optional.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// ResponseFormat specifies the output format.
	// Can be {"type": "text"}, {"type": "json_object"}, or {"type": "json_schema", "json_schema": {...}}
	// Optional.
	ResponseFormat interface{} `json:"response_format,omitempty"`

	// Modalities specifies input/output modalities (e.g., ["text", "audio"]).
	// Optional.
	Modalities []string `json:"modalities,omitempty"`

	// Audio configuration for audio output (GPT-4o Audio).
	// Optional.
	Audio *AudioConfig `json:"audio,omitempty"`

	// ReasoningEffort controls reasoning for o1/o3 models.
	// Values: "minimal", "low", "medium", "high"
	// Optional.
	ReasoningEffort string `json:"reasoning_effort,omitempty"`

	// --- Zaguan Extensions ---

	// ProviderOptions contains provider-specific parameters.
	// This is the primary extension mechanism for accessing provider features.
	// Optional.
	ProviderOptions map[string]interface{} `json:"provider_specific_params,omitempty"`

	// ExtraBody is an alternative name for ProviderOptions (OpenAI SDK compatibility).
	// Optional.
	ExtraBody map[string]interface{} `json:"extra_body,omitempty"`

	// VirtualModelID specifies a virtual model alias.
	// Optional.
	VirtualModelID string `json:"virtual_model_id,omitempty"`

	// Metadata for application-specific tracking (not interpreted by CoreX).
	// Optional.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Store enables conversation storage (if supported).
	// Optional.
	Store *bool `json:"store,omitempty"`

	// Thinking controls DeepSeek thinking output.
	// Optional.
	Thinking *bool `json:"thinking,omitempty"`
}

// Message represents a single message in a conversation.
type Message struct {
	// Role is the message role.
	// Values: "system", "user", "assistant", "tool", "function", "developer"
	// Required.
	Role string `json:"role"`

	// Content is the message content.
	// Can be a string or an array of ContentPart for multimodal messages.
	// Required for most roles.
	Content interface{} `json:"content,omitempty"`

	// Name is the name of the user/assistant/function.
	// Optional.
	Name string `json:"name,omitempty"`

	// ToolCalls are the tool calls made by the assistant.
	// Only for assistant role.
	// Optional.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// ToolCallID is the ID of the tool call being responded to.
	// Only for tool role.
	// Optional.
	ToolCallID string `json:"tool_call_id,omitempty"`

	// FunctionCall is the legacy function call format.
	// Deprecated: Use ToolCalls instead.
	// Optional.
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// ContentPart represents a part of multimodal content.
type ContentPart struct {
	// Type is the content type.
	// Values: "text", "image_url", "input_audio"
	Type string `json:"type"`

	// Text content (for type="text").
	Text string `json:"text,omitempty"`

	// ImageURL content (for type="image_url").
	ImageURL *ImageURL `json:"image_url,omitempty"`

	// InputAudio content (for type="input_audio").
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

// ImageURL represents an image URL or base64-encoded image.
type ImageURL struct {
	// URL is the image URL or data URI.
	// Examples: "https://example.com/image.jpg" or "data:image/jpeg;base64,..."
	URL string `json:"url"`

	// Detail controls image processing detail.
	// Values: "auto", "low", "high"
	// Optional.
	Detail string `json:"detail,omitempty"`
}

// InputAudio represents audio input.
type InputAudio struct {
	// Data is the base64-encoded audio data.
	Data string `json:"data"`

	// Format is the audio format.
	// Values: "wav", "mp3"
	Format string `json:"format"`
}

// AudioConfig represents audio output configuration.
type AudioConfig struct {
	// Voice is the voice to use for audio output.
	// Values: "alloy", "echo", "fable", "onyx", "nova", "shimmer"
	Voice string `json:"voice,omitempty"`

	// Format is the audio output format.
	// Values: "wav", "mp3", "opus", "aac", "flac", "pcm"
	Format string `json:"format,omitempty"`
}

// Tool represents a tool/function available to the model.
type Tool struct {
	// Type is the tool type (currently only "function").
	Type string `json:"type"`

	// Function is the function definition.
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition defines a function that can be called.
type FunctionDefinition struct {
	// Name is the function name.
	Name string `json:"name"`

	// Description explains what the function does.
	Description string `json:"description,omitempty"`

	// Parameters is the JSON Schema for the function parameters.
	Parameters interface{} `json:"parameters,omitempty"`

	// Strict enables strict schema adherence.
	Strict *bool `json:"strict,omitempty"`
}

// ToolCall represents a tool call made by the model.
type ToolCall struct {
	// ID is the unique identifier for this tool call.
	ID string `json:"id"`

	// Type is the tool call type (currently only "function").
	Type string `json:"type"`

	// Function is the function call details.
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call.
type FunctionCall struct {
	// Name is the function name.
	Name string `json:"name"`

	// Arguments is the JSON-encoded function arguments.
	Arguments string `json:"arguments"`
}

// ChatResponse represents a response from the chat completions endpoint.
type ChatResponse struct {
	// ID is the unique identifier for this completion.
	ID string `json:"id"`

	// Object is the object type (always "chat.completion").
	Object string `json:"object"`

	// Created is the Unix timestamp of when the completion was created.
	Created int64 `json:"created"`

	// Model is the model used for the completion.
	Model string `json:"model"`

	// Choices are the completion choices.
	Choices []Choice `json:"choices"`

	// Usage contains token usage information.
	Usage Usage `json:"usage"`

	// SystemFingerprint is a unique identifier for the backend configuration.
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
}

// Choice represents a completion choice.
type Choice struct {
	// Index is the choice index.
	Index int `json:"index"`

	// Message is the generated message (for non-streaming).
	Message *Message `json:"message,omitempty"`

	// Delta is the message delta (for streaming).
	Delta *Message `json:"delta,omitempty"`

	// FinishReason indicates why the generation stopped.
	// Values: "stop", "length", "tool_calls", "content_filter", "function_call"
	FinishReason string `json:"finish_reason,omitempty"`

	// Logprobs contains log probability information.
	Logprobs interface{} `json:"logprobs,omitempty"`
}

// Usage represents token usage information.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens.
	TotalTokens int `json:"total_tokens"`

	// PromptTokensDetails provides detailed breakdown of prompt tokens.
	PromptTokensDetails *TokenDetails `json:"prompt_tokens_details,omitempty"`

	// CompletionTokensDetails provides detailed breakdown of completion tokens.
	CompletionTokensDetails *TokenDetails `json:"completion_tokens_details,omitempty"`
}

// TokenDetails provides detailed token usage breakdown.
type TokenDetails struct {
	// ReasoningTokens is the number of reasoning/thinking tokens.
	// Populated by: OpenAI o1/o3, Google Gemini, Anthropic, DeepSeek, Alibaba Qwen
	// NOT populated by: Perplexity (uses <think> tags in content instead)
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`

	// CachedTokens is the number of tokens served from cache.
	// Populated by: Anthropic (prompt caching), OpenAI (prompt caching)
	CachedTokens int `json:"cached_tokens,omitempty"`

	// AudioTokens is the number of audio tokens.
	// Populated by: OpenAI GPT-4o Audio
	AudioTokens int `json:"audio_tokens,omitempty"`

	// AcceptedPredictionTokens is the number of accepted predicted tokens.
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`

	// RejectedPredictionTokens is the number of rejected predicted tokens.
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

// HasReasoningTokens returns true if reasoning tokens are present.
func (u *Usage) HasReasoningTokens() bool {
	return u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0
}

// HasCachedTokens returns true if cached tokens are present.
func (u *Usage) HasCachedTokens() bool {
	return u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0
}
