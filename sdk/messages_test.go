package zaguansdk

import (
	"testing"
)

func TestMessagesRequest_Validation(t *testing.T) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.Model == "" {
		t.Error("Model should not be empty")
	}
	if req.MaxTokens == 0 {
		t.Error("MaxTokens should not be zero")
	}
	if len(req.Messages) == 0 {
		t.Error("Messages should not be empty")
	}
}

func TestMessagesRequest_WithSystem(t *testing.T) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		System:    "You are a helpful assistant.",
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.System != "You are a helpful assistant." {
		t.Errorf("System = %v, want 'You are a helpful assistant.'", req.System)
	}
}

func TestMessagesRequest_WithThinking(t *testing.T) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Solve this problem"},
		},
		Thinking: &AnthropicThinkingConfig{
			Type:         "enabled",
			BudgetTokens: 5000,
		},
	}

	if req.Thinking == nil {
		t.Error("Thinking should not be nil")
	}
	if req.Thinking.Type != "enabled" {
		t.Errorf("Thinking.Type = %v, want enabled", req.Thinking.Type)
	}
	if req.Thinking.BudgetTokens != 5000 {
		t.Errorf("Thinking.BudgetTokens = %d, want 5000", req.Thinking.BudgetTokens)
	}
}

func TestMessagesRequest_WithTemperature(t *testing.T) {
	temp := 0.7
	req := MessagesRequest{
		Model:       "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens:   1024,
		Temperature: &temp,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.Temperature == nil {
		t.Error("Temperature should not be nil")
	}
	if *req.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", *req.Temperature)
	}
}

func TestMessagesRequest_WithStopSequences(t *testing.T) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
		StopSequences: []string{"\n\n", "END"},
	}

	if len(req.StopSequences) != 2 {
		t.Errorf("len(StopSequences) = %d, want 2", len(req.StopSequences))
	}
}

func TestMessagesRequest_WithMetadata(t *testing.T) {
	req := MessagesRequest{
		Model:     "anthropic/claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
		Metadata: map[string]interface{}{
			"user_id": "123",
			"session": "abc",
		},
	}

	if len(req.Metadata) != 2 {
		t.Errorf("len(Metadata) = %d, want 2", len(req.Metadata))
	}
	if req.Metadata["user_id"] != "123" {
		t.Errorf("Metadata[user_id] = %v, want 123", req.Metadata["user_id"])
	}
}

func TestAnthropicMessage_StringContent(t *testing.T) {
	msg := AnthropicMessage{
		Role:    "user",
		Content: "Hello, Claude!",
	}

	if msg.Role != "user" {
		t.Errorf("Role = %v, want user", msg.Role)
	}

	// Content should be a string
	if str, ok := msg.Content.(string); ok {
		if str != "Hello, Claude!" {
			t.Errorf("Content = %v, want 'Hello, Claude!'", str)
		}
	} else {
		t.Error("Content should be a string")
	}
}

func TestAnthropicMessage_ArrayContent(t *testing.T) {
	msg := AnthropicMessage{
		Role: "user",
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": "Hello",
			},
			{
				"type": "image",
				"source": map[string]interface{}{
					"type": "url",
					"url":  "https://example.com/image.jpg",
				},
			},
		},
	}

	if msg.Role != "user" {
		t.Errorf("Role = %v, want user", msg.Role)
	}

	// Content should be an array
	if arr, ok := msg.Content.([]map[string]interface{}); ok {
		if len(arr) != 2 {
			t.Errorf("len(Content) = %d, want 2", len(arr))
		}
	} else {
		t.Error("Content should be an array")
	}
}

func TestMessagesResponse_WithThinking(t *testing.T) {
	resp := MessagesResponse{
		ID:   "msg_123",
		Type: "message",
		Role: "assistant",
		Content: []AnthropicContentBlock{
			{
				Type:     "thinking",
				Thinking: "Let me analyze this problem...",
			},
			{
				Type: "text",
				Text: "Here's my answer.",
			},
		},
		Model:      "anthropic/claude-3-5-sonnet-20241022",
		StopReason: "end_turn",
		Usage: AnthropicUsage{
			InputTokens:  10,
			OutputTokens: 50,
		},
	}

	if len(resp.Content) != 2 {
		t.Errorf("len(Content) = %d, want 2", len(resp.Content))
	}
	if resp.Content[0].Type != "thinking" {
		t.Errorf("Content[0].Type = %v, want thinking", resp.Content[0].Type)
	}
	if resp.Content[0].Thinking == "" {
		t.Error("Thinking content should not be empty")
	}
}

func TestAnthropicUsage_WithCaching(t *testing.T) {
	usage := AnthropicUsage{
		InputTokens:              100,
		OutputTokens:             50,
		CacheCreationInputTokens: 80,
		CacheReadInputTokens:     20,
	}

	if usage.InputTokens != 100 {
		t.Errorf("InputTokens = %d, want 100", usage.InputTokens)
	}
	if usage.CacheCreationInputTokens != 80 {
		t.Errorf("CacheCreationInputTokens = %d, want 80", usage.CacheCreationInputTokens)
	}
	if usage.CacheReadInputTokens != 20 {
		t.Errorf("CacheReadInputTokens = %d, want 20", usage.CacheReadInputTokens)
	}
}

func TestCountTokensRequest(t *testing.T) {
	req := CountTokensRequest{
		Model: "anthropic/claude-3-5-sonnet-20241022",
		Messages: []AnthropicMessage{
			{Role: "user", Content: "Hello"},
		},
		System: "You are helpful.",
	}

	if req.Model == "" {
		t.Error("Model should not be empty")
	}
	if len(req.Messages) == 0 {
		t.Error("Messages should not be empty")
	}
}

func TestCountTokensResponse(t *testing.T) {
	resp := CountTokensResponse{
		InputTokens: 42,
	}

	if resp.InputTokens != 42 {
		t.Errorf("InputTokens = %d, want 42", resp.InputTokens)
	}
}

func TestMessagesBatchRequest(t *testing.T) {
	req := MessagesBatchRequest{
		Requests: []MessagesBatchItem{
			{
				CustomID: "req-1",
				Params: MessagesRequest{
					Model:     "anthropic/claude-3-5-sonnet-20241022",
					MaxTokens: 1024,
					Messages: []AnthropicMessage{
						{Role: "user", Content: "Hello"},
					},
				},
			},
			{
				CustomID: "req-2",
				Params: MessagesRequest{
					Model:     "anthropic/claude-3-5-sonnet-20241022",
					MaxTokens: 1024,
					Messages: []AnthropicMessage{
						{Role: "user", Content: "Hi"},
					},
				},
			},
		},
	}

	if len(req.Requests) != 2 {
		t.Errorf("len(Requests) = %d, want 2", len(req.Requests))
	}
	if req.Requests[0].CustomID != "req-1" {
		t.Errorf("Requests[0].CustomID = %v, want req-1", req.Requests[0].CustomID)
	}
}

func TestMessagesBatchResponse(t *testing.T) {
	resp := MessagesBatchResponse{
		ID:               "batch_123",
		Type:             "message_batch",
		ProcessingStatus: "in_progress",
		RequestCounts: MessagesBatchRequestCounts{
			Processing: 5,
			Succeeded:  10,
			Errored:    1,
			Canceled:   0,
			Expired:    0,
		},
		CreatedAt: "2025-11-19T12:00:00Z",
		ExpiresAt: "2025-11-20T12:00:00Z",
	}

	if resp.ID != "batch_123" {
		t.Errorf("ID = %v, want batch_123", resp.ID)
	}
	if resp.ProcessingStatus != "in_progress" {
		t.Errorf("ProcessingStatus = %v, want in_progress", resp.ProcessingStatus)
	}
	if resp.RequestCounts.Succeeded != 10 {
		t.Errorf("RequestCounts.Succeeded = %d, want 10", resp.RequestCounts.Succeeded)
	}
}

func TestMessagesBatchRequestCounts(t *testing.T) {
	counts := MessagesBatchRequestCounts{
		Processing: 5,
		Succeeded:  10,
		Errored:    2,
		Canceled:   1,
		Expired:    0,
	}

	total := counts.Processing + counts.Succeeded + counts.Errored + counts.Canceled + counts.Expired
	if total != 18 {
		t.Errorf("total count = %d, want 18", total)
	}
}
