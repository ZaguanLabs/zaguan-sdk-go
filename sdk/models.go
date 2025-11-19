package zaguansdk

import "context"

// Model represents a model available in Zaguan CoreX.
//
// Models are returned from the GET /v1/models endpoint and include
// provider-prefixed IDs (e.g., "openai/gpt-4o", "anthropic/claude-3-5-sonnet").
type Model struct {
	// ID is the model identifier.
	// Format: "provider/model-name"
	// Example: "openai/gpt-4o-mini"
	ID string `json:"id"`

	// Object is the object type (always "model").
	Object string `json:"object"`

	// Created is the Unix timestamp of when the model was created.
	Created int64 `json:"created,omitempty"`

	// OwnedBy is the organization that owns the model.
	// Example: "openai", "anthropic", "google"
	OwnedBy string `json:"owned_by,omitempty"`

	// Description is a human-readable description of the model.
	Description string `json:"description,omitempty"`

	// Metadata contains additional model information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Permission information (OpenAI compatibility).
	Permission []ModelPermission `json:"permission,omitempty"`
}

// ModelPermission represents permissions for a model.
type ModelPermission struct {
	ID                 string `json:"id"`
	Object             string `json:"object"`
	Created            int64  `json:"created"`
	AllowCreateEngine  bool   `json:"allow_create_engine"`
	AllowSampling      bool   `json:"allow_sampling"`
	AllowLogprobs      bool   `json:"allow_logprobs"`
	AllowSearchIndices bool   `json:"allow_search_indices"`
	AllowView          bool   `json:"allow_view"`
	AllowFineTuning    bool   `json:"allow_fine_tuning"`
	Organization       string `json:"organization"`
	Group              string `json:"group"`
	IsBlocking         bool   `json:"is_blocking"`
}

// ModelsResponse represents the response from GET /v1/models.
type ModelsResponse struct {
	// Object is the object type (always "list").
	Object string `json:"object"`

	// Data is the list of models.
	Data []Model `json:"data"`
}

// ListModels retrieves all available models from Zaguan CoreX.
//
// This includes models from all configured providers with their provider-prefixed IDs.
//
// Example:
//
//	models, err := client.ListModels(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, model := range models {
//		fmt.Printf("%s - %s\n", model.ID, model.Description)
//	}
func (c *Client) ListModels(ctx context.Context, opts *RequestOptions) ([]Model, error) {
	c.log(ctx, LogLevelDebug, "listing models")

	// TODO: Implement HTTP request
	return nil, nil
}

// GetModel retrieves information about a specific model.
//
// Example:
//
//	model, err := client.GetModel(ctx, "openai/gpt-4o", nil)
func (c *Client) GetModel(ctx context.Context, modelID string, opts *RequestOptions) (*Model, error) {
	c.log(ctx, LogLevelDebug, "getting model", "model_id", modelID)

	// TODO: Implement HTTP request
	return nil, nil
}

// DeleteModel deletes a fine-tuned model (if supported).
//
// Example:
//
//	err := client.DeleteModel(ctx, "ft:gpt-3.5-turbo:org:model:id", nil)
func (c *Client) DeleteModel(ctx context.Context, modelID string, opts *RequestOptions) error {
	c.log(ctx, LogLevelDebug, "deleting model", "model_id", modelID)

	// TODO: Implement HTTP request
	return nil
}
