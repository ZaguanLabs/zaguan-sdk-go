package zaguansdk

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// ChatStream represents a streaming chat completion response.
//
// Use Recv() to read events from the stream and Close() to clean up resources.
type ChatStream struct {
	reader *bufio.Reader
	resp   *http.Response
	ctx    context.Context
	closed bool
}

// Recv reads the next event from the chat stream.
//
// Returns io.EOF when the stream is complete.
// Returns an error if the stream encounters an error.
//
// Example:
//
//	for {
//		event, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatal(err)
//		}
//		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
//			fmt.Print(event.Choices[0].Delta.Content)
//		}
//	}
func (s *ChatStream) Recv() (*ChatStreamEvent, error) {
	if s.closed {
		return nil, errors.New("stream is closed")
	}

	// Check context
	if err := s.ctx.Err(); err != nil {
		_ = s.Close() // Explicitly ignore error in cleanup
		return nil, err
	}

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				_ = s.Close() // Explicitly ignore error in cleanup
			}
			return nil, err
		}

		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for SSE data prefix
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Extract data
		data := strings.TrimPrefix(line, "data: ")

		// Check for stream end
		if data == "[DONE]" {
			_ = s.Close() // Explicitly ignore error in cleanup
			return nil, io.EOF
		}

		// Parse JSON event
		var event ChatStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, fmt.Errorf("failed to parse stream event: %w", err)
		}

		return &event, nil
	}
}

// Close closes the stream and releases resources.
func (s *ChatStream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	if s.resp != nil && s.resp.Body != nil {
		return s.resp.Body.Close()
	}
	return nil
}

// ChatStreamEvent represents a single event in a chat stream.
type ChatStreamEvent struct {
	// ID is the unique identifier for this completion.
	ID string `json:"id"`

	// Object is the object type (always "chat.completion.chunk").
	Object string `json:"object"`

	// Created is the Unix timestamp of when the completion was created.
	Created int64 `json:"created"`

	// Model is the model used for the completion.
	Model string `json:"model"`

	// Choices is the list of completion choices.
	Choices []ChatStreamChoice `json:"choices"`

	// Usage contains token usage information (only in final event).
	Usage *Usage `json:"usage,omitempty"`
}

// ChatStreamChoice represents a choice in a streaming response.
type ChatStreamChoice struct {
	// Index is the index of this choice.
	Index int `json:"index"`

	// Delta contains the incremental content.
	Delta ChatStreamDelta `json:"delta"`

	// FinishReason indicates why the generation stopped.
	// Values: "stop", "length", "tool_calls", "content_filter", null
	FinishReason *string `json:"finish_reason"`

	// Logprobs contains log probabilities (if requested).
	Logprobs interface{} `json:"logprobs,omitempty"`
}

// ChatStreamDelta represents incremental content in a streaming response.
type ChatStreamDelta struct {
	// Role is the message role (only in first chunk).
	Role string `json:"role,omitempty"`

	// Content is the incremental text content.
	Content string `json:"content,omitempty"`

	// ToolCalls contains incremental tool call information.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatStream sends a streaming chat completion request to Zaguan CoreX.
//
// The stream must be closed when done to release resources.
//
// Example:
//
//	stream, err := client.ChatStream(ctx, zaguansdk.ChatRequest{
//		Model: "openai/gpt-4o",
//		Messages: []zaguansdk.Message{
//			{Role: "user", Content: "Tell me a story"},
//		},
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer stream.Close()
//
//	for {
//		event, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatal(err)
//		}
//		if len(event.Choices) > 0 {
//			fmt.Print(event.Choices[0].Delta.Content)
//		}
//	}
func (c *Client) ChatStream(ctx context.Context, req ChatRequest, opts *RequestOptions) (*ChatStream, error) {
	// Validate request
	if err := validateChatRequest(&req); err != nil {
		return nil, err
	}

	// Ensure stream is true
	req.Stream = true

	c.log(ctx, LogLevelDebug, "sending streaming chat completion request",
		"model", req.Model,
		"message_count", len(req.Messages))

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/chat/completions",
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
	resp, err := c.internalHTTP.Do(ctx, reqCfg)
	if err != nil {
		c.log(ctx, LogLevelError, "streaming chat completion request failed", "error", err)
		return nil, err
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, internal.ParseErrorResponse(resp)
	}

	c.log(ctx, LogLevelDebug, "streaming chat completion request started")

	// Create stream
	stream := &ChatStream{
		reader: bufio.NewReader(resp.Body),
		resp:   resp,
		ctx:    ctx,
		closed: false,
	}

	return stream, nil
}

// MessagesStream represents a streaming Anthropic Messages response.
type MessagesStream struct {
	reader *bufio.Reader
	resp   *http.Response
	ctx    context.Context
	closed bool
}

// Recv reads the next event from the messages stream.
//
// Returns io.EOF when the stream is complete.
func (s *MessagesStream) Recv() (*MessagesStreamEvent, error) {
	if s.closed {
		return nil, errors.New("stream is closed")
	}

	// Check context
	if err := s.ctx.Err(); err != nil {
		_ = s.Close() // Explicitly ignore error in cleanup
		return nil, err
	}

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				_ = s.Close() // Explicitly ignore error in cleanup
			}
			return nil, err
		}

		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for SSE event prefix
		if strings.HasPrefix(line, "event: ") {
			// Event type line - we'll read the data on the next iteration
			continue
		}

		// Check for SSE data prefix
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Extract data
		data := strings.TrimPrefix(line, "data: ")

		// Parse JSON event
		var event MessagesStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, fmt.Errorf("failed to parse stream event: %w", err)
		}

		// Check for stream end
		if event.Type == "message_stop" {
			_ = s.Close() // Explicitly ignore error in cleanup
			return &event, io.EOF
		}

		return &event, nil
	}
}

// Close closes the stream and releases resources.
func (s *MessagesStream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	if s.resp != nil && s.resp.Body != nil {
		return s.resp.Body.Close()
	}
	return nil
}

// MessagesStreamEvent represents a single event in an Anthropic Messages stream.
type MessagesStreamEvent struct {
	// Type is the event type.
	// Values: "message_start", "content_block_start", "content_block_delta",
	//         "content_block_stop", "message_delta", "message_stop"
	Type string `json:"type"`

	// Message contains the initial message (for message_start).
	Message *MessagesResponse `json:"message,omitempty"`

	// Index is the content block index (for content_block_* events).
	Index int `json:"index,omitempty"`

	// ContentBlock contains the content block (for content_block_start).
	ContentBlock *AnthropicContentBlock `json:"content_block,omitempty"`

	// Delta contains incremental content (for content_block_delta, message_delta).
	Delta *MessagesStreamDelta `json:"delta,omitempty"`

	// Usage contains token usage updates (for message_delta).
	Usage *AnthropicUsage `json:"usage,omitempty"`
}

// MessagesStreamDelta represents incremental content in a Messages stream.
type MessagesStreamDelta struct {
	// Type is the delta type.
	// Values: "text_delta", "thinking_delta", "input_json_delta"
	Type string `json:"type,omitempty"`

	// Text is the incremental text content.
	Text string `json:"text,omitempty"`

	// Thinking is the incremental thinking content.
	Thinking string `json:"thinking,omitempty"`

	// PartialJSON is the incremental JSON for tool inputs.
	PartialJSON string `json:"partial_json,omitempty"`

	// StopReason indicates why the generation stopped (for message_delta).
	StopReason string `json:"stop_reason,omitempty"`

	// StopSequence is the stop sequence that was matched (for message_delta).
	StopSequence string `json:"stop_sequence,omitempty"`
}

// MessagesStream sends a streaming Anthropic Messages request.
//
// The stream must be closed when done to release resources.
//
// Example:
//
//	stream, err := client.MessagesStream(ctx, zaguansdk.MessagesRequest{
//		Model: "anthropic/claude-3-5-sonnet-20241022",
//		MaxTokens: 1024,
//		Messages: []zaguansdk.AnthropicMessage{
//			{Role: "user", Content: "Tell me a story"},
//		},
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer stream.Close()
//
//	for {
//		event, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatal(err)
//		}
//		if event.Delta != nil && event.Delta.Text != "" {
//			fmt.Print(event.Delta.Text)
//		}
//	}
func (c *Client) MessagesStream(ctx context.Context, req MessagesRequest, opts *RequestOptions) (*MessagesStream, error) {
	// Validate request
	if err := validateMessagesRequest(&req); err != nil {
		return nil, err
	}

	// Ensure stream is true
	req.Stream = true

	c.log(ctx, LogLevelDebug, "sending streaming messages request",
		"model", req.Model,
		"message_count", len(req.Messages))

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/messages",
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
	resp, err := c.internalHTTP.Do(ctx, reqCfg)
	if err != nil {
		c.log(ctx, LogLevelError, "streaming messages request failed", "error", err)
		return nil, err
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, internal.ParseErrorResponse(resp)
	}

	c.log(ctx, LogLevelDebug, "streaming messages request started")

	// Create stream
	stream := &MessagesStream{
		reader: bufio.NewReader(resp.Body),
		resp:   resp,
		ctx:    ctx,
		closed: false,
	}

	return stream, nil
}
