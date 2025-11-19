package zaguansdk

import "context"

// ModelCapabilities represents the capabilities of a specific model.
//
// This provides detailed information about what features a model supports,
// including vision, tools, reasoning, audio, and provider-specific features.
type ModelCapabilities struct {
	// ModelID is the model identifier.
	// Format: "provider/model-name"
	ModelID string `json:"model_id"`

	// Provider is the provider name.
	// Example: "openai", "anthropic", "google"
	Provider string `json:"provider,omitempty"`

	// SupportsVision indicates if the model supports image inputs.
	SupportsVision bool `json:"supports_vision"`

	// SupportsTools indicates if the model supports tool/function calling.
	SupportsTools bool `json:"supports_tools"`

	// SupportsReasoning indicates if the model supports reasoning/thinking.
	// This includes models like o1, o3, Gemini with reasoning, Claude with extended thinking.
	SupportsReasoning bool `json:"supports_reasoning"`

	// SupportsAudioInput indicates if the model supports audio inputs.
	SupportsAudioInput bool `json:"supports_audio_input,omitempty"`

	// SupportsAudioOutput indicates if the model supports audio outputs.
	SupportsAudioOutput bool `json:"supports_audio_output,omitempty"`

	// SupportsStreaming indicates if the model supports streaming responses.
	SupportsStreaming bool `json:"supports_streaming,omitempty"`

	// SupportsSystemMessages indicates if the model supports system messages.
	SupportsSystemMessages bool `json:"supports_system_messages,omitempty"`

	// MaxContextTokens is the maximum context window size in tokens.
	MaxContextTokens int `json:"max_context_tokens,omitempty"`

	// MaxOutputTokens is the maximum output length in tokens.
	MaxOutputTokens int `json:"max_output_tokens,omitempty"`

	// InputCostPer1M is the cost per 1M input tokens (in USD).
	InputCostPer1M float64 `json:"input_cost_per_1m,omitempty"`

	// OutputCostPer1M is the cost per 1M output tokens (in USD).
	OutputCostPer1M float64 `json:"output_cost_per_1m,omitempty"`

	// ReasoningCostPer1M is the cost per 1M reasoning tokens (in USD).
	ReasoningCostPer1M float64 `json:"reasoning_cost_per_1m,omitempty"`

	// ProviderSpecific contains provider-specific capability information.
	ProviderSpecific map[string]interface{} `json:"provider_specific,omitempty"`

	// Features is a list of supported features.
	// Examples: "json_mode", "structured_outputs", "prompt_caching", "extended_thinking"
	Features []string `json:"features,omitempty"`

	// Modalities is a list of supported modalities.
	// Examples: "text", "image", "audio"
	Modalities []string `json:"modalities,omitempty"`
}

// CapabilitiesResponse represents the response from GET /v1/capabilities.
type CapabilitiesResponse struct {
	// Capabilities is a map of model ID to capabilities.
	Capabilities map[string]ModelCapabilities `json:"capabilities,omitempty"`

	// Models is an array of model capabilities (alternative format).
	Models []ModelCapabilities `json:"models,omitempty"`
}

// GetCapabilities retrieves capability information for all models.
//
// This endpoint provides detailed information about what each model supports,
// including vision, tools, reasoning, context limits, and pricing.
//
// Example:
//
//	caps, err := client.GetCapabilities(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, cap := range caps {
//		fmt.Printf("%s: vision=%v, tools=%v, reasoning=%v\n",
//			cap.ModelID, cap.SupportsVision, cap.SupportsTools, cap.SupportsReasoning)
//	}
func (c *Client) GetCapabilities(ctx context.Context, opts *RequestOptions) ([]ModelCapabilities, error) {
	c.log(ctx, LogLevelDebug, "getting model capabilities")

	// TODO: Implement HTTP request
	return nil, nil
}

// GetModelCapabilities retrieves capability information for a specific model.
//
// Example:
//
//	cap, err := client.GetModelCapabilities(ctx, "openai/gpt-4o", nil)
//	if cap.SupportsVision {
//		fmt.Println("Model supports vision!")
//	}
func (c *Client) GetModelCapabilities(ctx context.Context, modelID string, opts *RequestOptions) (*ModelCapabilities, error) {
	c.log(ctx, LogLevelDebug, "getting model capabilities", "model_id", modelID)

	// Get all capabilities
	caps, err := c.GetCapabilities(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Find the specific model
	for _, cap := range caps {
		if cap.ModelID == modelID {
			return &cap, nil
		}
	}

	return nil, &APIError{
		StatusCode: 404,
		Message:    "model not found",
		Type:       "not_found",
	}
}

// SupportsVision checks if a model supports vision/image inputs.
//
// Example:
//
//	if client.SupportsVision(ctx, "openai/gpt-4o", nil) {
//		// Use vision features
//	}
func (c *Client) SupportsVision(ctx context.Context, modelID string, opts *RequestOptions) bool {
	cap, err := c.GetModelCapabilities(ctx, modelID, opts)
	if err != nil {
		return false
	}
	return cap.SupportsVision
}

// SupportsTools checks if a model supports tool/function calling.
func (c *Client) SupportsTools(ctx context.Context, modelID string, opts *RequestOptions) bool {
	cap, err := c.GetModelCapabilities(ctx, modelID, opts)
	if err != nil {
		return false
	}
	return cap.SupportsTools
}

// SupportsReasoning checks if a model supports reasoning/thinking.
func (c *Client) SupportsReasoning(ctx context.Context, modelID string, opts *RequestOptions) bool {
	cap, err := c.GetModelCapabilities(ctx, modelID, opts)
	if err != nil {
		return false
	}
	return cap.SupportsReasoning
}
