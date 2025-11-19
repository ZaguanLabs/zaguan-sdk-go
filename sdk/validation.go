package zaguansdk

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}

// validateChatRequest validates a ChatRequest before sending to the API.
func validateChatRequest(req *ChatRequest) error {
	// Model is required
	if req.Model == "" {
		return &ValidationError{Field: "model", Message: "model is required"}
	}

	// Messages is required and must not be empty
	if len(req.Messages) == 0 {
		return &ValidationError{Field: "messages", Message: "at least one message is required"}
	}

	// Validate temperature range
	if req.Temperature != nil {
		if *req.Temperature < 0 || *req.Temperature > 2 {
			return &ValidationError{
				Field:   "temperature",
				Message: "temperature must be between 0 and 2",
			}
		}
	}

	// Validate top_p range
	if req.TopP != nil {
		if *req.TopP < 0 || *req.TopP > 1 {
			return &ValidationError{
				Field:   "top_p",
				Message: "top_p must be between 0 and 1",
			}
		}
	}

	// Validate max_tokens
	if req.MaxTokens != nil && *req.MaxTokens < 1 {
		return &ValidationError{
			Field:   "max_tokens",
			Message: "max_tokens must be at least 1",
		}
	}

	// Validate presence_penalty range
	if req.PresencePenalty != nil {
		if *req.PresencePenalty < -2 || *req.PresencePenalty > 2 {
			return &ValidationError{
				Field:   "presence_penalty",
				Message: "presence_penalty must be between -2 and 2",
			}
		}
	}

	// Validate frequency_penalty range
	if req.FrequencyPenalty != nil {
		if *req.FrequencyPenalty < -2 || *req.FrequencyPenalty > 2 {
			return &ValidationError{
				Field:   "frequency_penalty",
				Message: "frequency_penalty must be between -2 and 2",
			}
		}
	}

	// Validate reasoning_effort
	if req.ReasoningEffort != "" {
		validEfforts := map[string]bool{
			"minimal": true,
			"low":     true,
			"medium":  true,
			"high":    true,
		}
		if !validEfforts[req.ReasoningEffort] {
			return &ValidationError{
				Field:   "reasoning_effort",
				Message: "reasoning_effort must be one of: minimal, low, medium, high",
			}
		}
	}

	return nil
}

// validateMessagesRequest validates a MessagesRequest before sending to the API.
func validateMessagesRequest(req *MessagesRequest) error {
	// Model is required
	if req.Model == "" {
		return &ValidationError{Field: "model", Message: "model is required"}
	}

	// Messages is required and must not be empty
	if len(req.Messages) == 0 {
		return &ValidationError{Field: "messages", Message: "at least one message is required"}
	}

	// MaxTokens is required for Anthropic API
	if req.MaxTokens < 1 {
		return &ValidationError{
			Field:   "max_tokens",
			Message: "max_tokens is required and must be at least 1",
		}
	}

	// Validate temperature range (Anthropic uses 0-1)
	if req.Temperature != nil {
		if *req.Temperature < 0 || *req.Temperature > 1 {
			return &ValidationError{
				Field:   "temperature",
				Message: "temperature must be between 0 and 1 for Anthropic models",
			}
		}
	}

	// Validate top_p range
	if req.TopP != nil {
		if *req.TopP < 0 || *req.TopP > 1 {
			return &ValidationError{
				Field:   "top_p",
				Message: "top_p must be between 0 and 1",
			}
		}
	}

	// Validate top_k
	if req.TopK != nil && *req.TopK < 1 {
		return &ValidationError{
			Field:   "top_k",
			Message: "top_k must be at least 1",
		}
	}

	// Validate thinking config
	if req.Thinking != nil {
		if req.Thinking.Type != "enabled" && req.Thinking.Type != "disabled" {
			return &ValidationError{
				Field:   "thinking.type",
				Message: "thinking.type must be 'enabled' or 'disabled'",
			}
		}
		if req.Thinking.BudgetTokens > 0 {
			if req.Thinking.BudgetTokens < 1000 || req.Thinking.BudgetTokens > 10000 {
				return &ValidationError{
					Field:   "thinking.budget_tokens",
					Message: "thinking.budget_tokens must be between 1000 and 10000",
				}
			}
		}
	}

	return nil
}

// validateConfig validates the client configuration.
func validateConfig(cfg *Config) error {
	if cfg.BaseURL == "" {
		return errors.New("BaseURL is required")
	}

	if cfg.APIKey == "" {
		return errors.New("APIKey is required")
	}

	// Validate base URL format
	baseURL := strings.ToLower(cfg.BaseURL)
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return errors.New("BaseURL must start with http:// or https://")
	}

	// Warn about http (not https) but don't fail
	// This is just basic validation, not security enforcement

	return nil
}

// validateModelID validates a model ID format.
func validateModelID(modelID string) error {
	if modelID == "" {
		return &ValidationError{Field: "model_id", Message: "model_id is required"}
	}

	// Model IDs should typically be in format "provider/model-name"
	// But we'll be lenient and just check it's not empty
	// The API will validate the actual format

	return nil
}
