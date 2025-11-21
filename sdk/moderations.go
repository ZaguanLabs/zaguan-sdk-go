// Package zaguansdk provides content moderation functionality for the Zaguan SDK.
//
// This file implements the Moderations API for checking whether content
// complies with usage policies. It provides classification across 11 categories
// including sexual content, hate speech, harassment, violence, and self-harm,
// with confidence scores for each category.
package zaguansdk

import (
	"context"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// ModerationRequest represents a request to classify content.
//
// The moderation endpoint checks whether content complies with usage policies.
type ModerationRequest struct {
	// Input is the text to classify.
	// Can be a string or array of strings.
	// Required.
	Input interface{} `json:"input"`

	// Model is the moderation model to use.
	// Examples: "text-moderation-latest", "text-moderation-stable"
	// Optional (default: "text-moderation-latest").
	Model string `json:"model,omitempty"`
}

// ModerationResponse represents the response from a moderation request.
type ModerationResponse struct {
	// ID is the unique identifier for the moderation request.
	ID string `json:"id"`

	// Model is the model used for moderation.
	Model string `json:"model"`

	// Results contains the moderation results for each input.
	Results []ModerationResult `json:"results"`
}

// ModerationResult represents the moderation result for a single input.
type ModerationResult struct {
	// Flagged indicates if the content violates usage policies.
	Flagged bool `json:"flagged"`

	// Categories contains flags for each category.
	Categories ModerationCategories `json:"categories"`

	// CategoryScores contains confidence scores for each category.
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

// ModerationCategories contains boolean flags for content categories.
type ModerationCategories struct {
	// Sexual indicates sexual content.
	Sexual bool `json:"sexual"`

	// Hate indicates hateful content.
	Hate bool `json:"hate"`

	// Harassment indicates harassment content.
	Harassment bool `json:"harassment"`

	// SelfHarm indicates self-harm content.
	SelfHarm bool `json:"self-harm"`

	// SexualMinors indicates sexual content involving minors.
	SexualMinors bool `json:"sexual/minors"`

	// HateThreatening indicates hateful content with threats.
	HateThreatening bool `json:"hate/threatening"`

	// ViolenceGraphic indicates graphic violence.
	ViolenceGraphic bool `json:"violence/graphic"`

	// SelfHarmIntent indicates intent to self-harm.
	SelfHarmIntent bool `json:"self-harm/intent"`

	// SelfHarmInstructions indicates instructions for self-harm.
	SelfHarmInstructions bool `json:"self-harm/instructions"`

	// HarassmentThreatening indicates harassment with threats.
	HarassmentThreatening bool `json:"harassment/threatening"`

	// Violence indicates violent content.
	Violence bool `json:"violence"`
}

// ModerationCategoryScores contains confidence scores for content categories.
type ModerationCategoryScores struct {
	// Sexual is the confidence score for sexual content (0-1).
	Sexual float64 `json:"sexual"`

	// Hate is the confidence score for hateful content (0-1).
	Hate float64 `json:"hate"`

	// Harassment is the confidence score for harassment (0-1).
	Harassment float64 `json:"harassment"`

	// SelfHarm is the confidence score for self-harm content (0-1).
	SelfHarm float64 `json:"self-harm"`

	// SexualMinors is the confidence score for sexual content involving minors (0-1).
	SexualMinors float64 `json:"sexual/minors"`

	// HateThreatening is the confidence score for hateful threats (0-1).
	HateThreatening float64 `json:"hate/threatening"`

	// ViolenceGraphic is the confidence score for graphic violence (0-1).
	ViolenceGraphic float64 `json:"violence/graphic"`

	// SelfHarmIntent is the confidence score for self-harm intent (0-1).
	SelfHarmIntent float64 `json:"self-harm/intent"`

	// SelfHarmInstructions is the confidence score for self-harm instructions (0-1).
	SelfHarmInstructions float64 `json:"self-harm/instructions"`

	// HarassmentThreatening is the confidence score for threatening harassment (0-1).
	HarassmentThreatening float64 `json:"harassment/threatening"`

	// Violence is the confidence score for violent content (0-1).
	Violence float64 `json:"violence"`
}

// CreateModeration classifies text content for policy violations.
//
// The moderation endpoint checks whether content complies with usage policies
// and provides detailed category scores.
//
// Example:
//
//	resp, err := client.CreateModeration(ctx, zaguansdk.ModerationRequest{
//		Input: "I want to hurt someone",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if resp.Results[0].Flagged {
//		fmt.Println("Content flagged for policy violations")
//		if resp.Results[0].Categories.Violence {
//			fmt.Println("- Violence detected")
//		}
//	}
func (c *Client) CreateModeration(ctx context.Context, req ModerationRequest, opts *RequestOptions) (*ModerationResponse, error) {
	// Validate request
	if err := validateModerationRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating moderation")

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/moderations",
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
	var resp ModerationResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create moderation request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create moderation request succeeded")

	return &resp, nil
}

// IsSafe is a helper method that returns true if the content is not flagged.
func (r *ModerationResult) IsSafe() bool {
	return !r.Flagged
}

// GetViolatedCategories returns a list of category names that were flagged.
func (r *ModerationResult) GetViolatedCategories() []string {
	var violated []string
	if r.Categories.Sexual {
		violated = append(violated, "sexual")
	}
	if r.Categories.Hate {
		violated = append(violated, "hate")
	}
	if r.Categories.Harassment {
		violated = append(violated, "harassment")
	}
	if r.Categories.SelfHarm {
		violated = append(violated, "self-harm")
	}
	if r.Categories.SexualMinors {
		violated = append(violated, "sexual/minors")
	}
	if r.Categories.HateThreatening {
		violated = append(violated, "hate/threatening")
	}
	if r.Categories.ViolenceGraphic {
		violated = append(violated, "violence/graphic")
	}
	if r.Categories.SelfHarmIntent {
		violated = append(violated, "self-harm/intent")
	}
	if r.Categories.SelfHarmInstructions {
		violated = append(violated, "self-harm/instructions")
	}
	if r.Categories.HarassmentThreatening {
		violated = append(violated, "harassment/threatening")
	}
	if r.Categories.Violence {
		violated = append(violated, "violence")
	}
	return violated
}
