// Package zaguansdk provides image generation functionality for the Zaguan SDK.
//
// This file implements the Images API for:
//   - Image Generation: Creating images from text prompts (DALL-E support)
//   - Image Editing: Modifying existing images (placeholder)
//   - Image Variations: Creating variations of existing images (placeholder)
//
// Supports DALL-E 2 and DALL-E 3 models with various sizes, quality levels, and styles.
package zaguansdk

import (
	"context"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// ImageGenerationRequest represents a request to generate images.
//
// Supports DALL-E and other image generation models.
type ImageGenerationRequest struct {
	// Prompt is the text description of the desired image(s).
	// Required.
	Prompt string `json:"prompt"`

	// Model is the model identifier to use.
	// Examples: "openai/dall-e-3", "openai/dall-e-2"
	// Optional (default: "dall-e-2").
	Model string `json:"model,omitempty"`

	// N is the number of images to generate (1-10).
	// Optional (default: 1). Only supported by dall-e-2.
	N *int `json:"n,omitempty"`

	// Quality specifies the quality of the image.
	// Values: "standard", "hd"
	// Optional (default: "standard"). Only supported by dall-e-3.
	Quality string `json:"quality,omitempty"`

	// ResponseFormat specifies the format of the generated images.
	// Values: "url", "b64_json"
	// Optional (default: "url").
	ResponseFormat string `json:"response_format,omitempty"`

	// Size specifies the size of the generated images.
	// dall-e-2: "256x256", "512x512", "1024x1024"
	// dall-e-3: "1024x1024", "1792x1024", "1024x1792"
	// Optional (default: "1024x1024").
	Size string `json:"size,omitempty"`

	// Style specifies the style of the generated images.
	// Values: "vivid", "natural"
	// Optional (default: "vivid"). Only supported by dall-e-3.
	Style string `json:"style,omitempty"`

	// User is an optional unique identifier for the end-user.
	// Optional.
	User string `json:"user,omitempty"`
}

// ImageEditRequest represents a request to edit an image.
//
// Creates an edited or extended image given an original image and a prompt.
type ImageEditRequest struct {
	// Image is the image to edit.
	// Must be a valid PNG file, less than 4MB, and square.
	// Can be a file path (string) or io.Reader.
	// Required.
	Image interface{}

	// ImageFileName is the name of the image file (required if Image is io.Reader).
	ImageFileName string

	// Prompt is the text description of the desired edits.
	// Required.
	Prompt string

	// Mask is an optional mask image.
	// Must be a valid PNG file, less than 4MB, and same dimensions as image.
	// Can be a file path (string) or io.Reader.
	// Optional.
	Mask interface{}

	// MaskFileName is the name of the mask file (required if Mask is io.Reader).
	MaskFileName string

	// Model is the model identifier to use.
	// Example: "openai/dall-e-2"
	// Optional (default: "dall-e-2").
	Model string

	// N is the number of images to generate (1-10).
	// Optional (default: 1).
	N *int

	// Size specifies the size of the generated images.
	// Values: "256x256", "512x512", "1024x1024"
	// Optional (default: "1024x1024").
	Size string

	// ResponseFormat specifies the format of the generated images.
	// Values: "url", "b64_json"
	// Optional (default: "url").
	ResponseFormat string

	// User is an optional unique identifier for the end-user.
	// Optional.
	User string
}

// ImageVariationRequest represents a request to create image variations.
//
// Creates variations of a given image.
type ImageVariationRequest struct {
	// Image is the image to create variations of.
	// Must be a valid PNG file, less than 4MB, and square.
	// Can be a file path (string) or io.Reader.
	// Required.
	Image interface{}

	// ImageFileName is the name of the image file (required if Image is io.Reader).
	ImageFileName string

	// Model is the model identifier to use.
	// Example: "openai/dall-e-2"
	// Optional (default: "dall-e-2").
	Model string

	// N is the number of images to generate (1-10).
	// Optional (default: 1).
	N *int

	// Size specifies the size of the generated images.
	// Values: "256x256", "512x512", "1024x1024"
	// Optional (default: "1024x1024").
	Size string

	// ResponseFormat specifies the format of the generated images.
	// Values: "url", "b64_json"
	// Optional (default: "url").
	ResponseFormat string

	// User is an optional unique identifier for the end-user.
	// Optional.
	User string
}

// ImageResponse represents the response from image generation/edit/variation.
type ImageResponse struct {
	// Created is the Unix timestamp of when the images were created.
	Created int64 `json:"created"`

	// Data is the list of generated images.
	Data []ImageData `json:"data"`
}

// ImageData represents a single generated image.
type ImageData struct {
	// URL is the URL of the generated image (if response_format is "url").
	URL string `json:"url,omitempty"`

	// B64JSON is the base64-encoded JSON of the generated image (if response_format is "b64_json").
	B64JSON string `json:"b64_json,omitempty"`

	// RevisedPrompt is the revised prompt that was used (DALL-E 3 only).
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// CreateImage generates images from a text prompt.
//
// Supports DALL-E 2 and DALL-E 3 models.
//
// Example:
//
//	resp, err := client.CreateImage(ctx, zaguansdk.ImageGenerationRequest{
//		Prompt: "A cute baby sea otter",
//		Model:  "openai/dall-e-3",
//		Size:   "1024x1024",
//		Quality: "hd",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Image URL:", resp.Data[0].URL)
func (c *Client) CreateImage(ctx context.Context, req ImageGenerationRequest, opts *RequestOptions) (*ImageResponse, error) {
	// Validate request
	if err := validateImageGenerationRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating image", "model", req.Model)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/images/generations",
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
	var resp ImageResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create image request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create image request succeeded", "count", len(resp.Data))

	return &resp, nil
}

// EditImage creates an edited or extended image given an original image and a prompt.
//
// Example:
//
//	resp, err := client.EditImage(ctx, zaguansdk.ImageEditRequest{
//		Image:  "/path/to/image.png",
//		Prompt: "Add a party hat to the otter",
//		Size:   "1024x1024",
//	}, nil)
func (c *Client) EditImage(ctx context.Context, req ImageEditRequest, opts *RequestOptions) (*ImageResponse, error) {
	// Validate request
	if err := validateImageEditRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "editing image", "model", req.Model)

	// Note: Image editing requires multipart form data
	// This is a simplified implementation - full implementation would handle file uploads
	return nil, &APIError{
		StatusCode: 501,
		Message:    "image editing not yet implemented - requires multipart form support",
		Type:       "not_implemented",
	}
}

// CreateImageVariation creates variations of a given image.
//
// Example:
//
//	resp, err := client.CreateImageVariation(ctx, zaguansdk.ImageVariationRequest{
//		Image: "/path/to/image.png",
//		N:     intPtr(2),
//		Size:  "1024x1024",
//	}, nil)
func (c *Client) CreateImageVariation(ctx context.Context, req ImageVariationRequest, opts *RequestOptions) (*ImageResponse, error) {
	// Validate request
	if err := validateImageVariationRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating image variation", "model", req.Model)

	// Note: Image variations require multipart form data
	// This is a simplified implementation - full implementation would handle file uploads
	return nil, &APIError{
		StatusCode: 501,
		Message:    "image variations not yet implemented - requires multipart form support",
		Type:       "not_implemented",
	}
}
